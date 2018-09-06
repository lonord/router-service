package netutil

import (
	"time"

	"github.com/lonord/router-service/base"
)

type NetSpeed struct {
	DevName       string `json:"name"`
	ReceiveSpeed  uint64 `json:"recv"`
	TransmitSpeed uint64 `json:"send"`
}

type NetSpeedReader struct {
	targetDevs   map[string]bool
	fileReaderFn ba.FileReaderFn
	recordTime   time.Time
	statusMap    map[string]DevStatus
}

func NewNetSpeedReader(fileReaderFn ba.FileReaderFn, targetDevs ...string) *NetSpeedReader {
	m := make(map[string]bool)
	for _, dev := range targetDevs {
		m[dev] = true
	}
	return &NetSpeedReader{
		fileReaderFn: fileReaderFn,
		targetDevs:   m,
	}
}

func (n *NetSpeedReader) Init() error {
	s, err := ReadDevStatus(n.fileReaderFn)
	if err != nil {
		return err
	}
	n.recordTime = time.Now()
	n.statusMap = convertToDevStatusMapWithFilter(s, n.targetDevs)
	return nil
}

func (n *NetSpeedReader) Read() ([]NetSpeed, error) {
	lt := n.recordTime
	ls := n.statusMap
	t := time.Now()
	s, err := ReadDevStatus(n.fileReaderFn)
	if err != nil {
		return nil, err
	}
	sMap := convertToDevStatusMapWithFilter(s, n.targetDevs)
	speedList := cal(ls, sMap, uint64((t.UnixNano()-lt.UnixNano())/1000000))
	n.recordTime = t
	n.statusMap = sMap
	return speedList, nil
}

func MeasureNetSpeed(timeSpanMilli uint64) ([]NetSpeed, error) {
	status1, err := readDevStatusDefault()
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Duration(timeSpanMilli) * time.Millisecond)
	status2, err := readDevStatusDefault()
	if err != nil {
		return nil, err
	}
	map1 := convertToDevStatusMap(status1)
	map2 := convertToDevStatusMap(status2)
	return cal(map1, map2, timeSpanMilli), nil
}

func readDevStatusDefault() ([]DevStatus, error) {
	return ReadDevStatus(ba.DefaultFileReader)
}

func cal(map1, map2 map[string]DevStatus, timeSpanMilli uint64) []NetSpeed {
	speedList := []NetSpeed{}
	for name := range map2 {
		s1, ok1 := map1[name]
		s2, ok2 := map2[name]
		if ok1 && ok2 {
			rSpeed := (s2.Receive.Bytes - s1.Receive.Bytes) * 1000 / uint64(timeSpanMilli)
			tSpeed := (s2.Transmit.Bytes - s1.Transmit.Bytes) * 1000 / uint64(timeSpanMilli)
			speedList = append(speedList, NetSpeed{
				DevName:       name,
				ReceiveSpeed:  rSpeed,
				TransmitSpeed: tSpeed,
			})
		}
	}
	return speedList
}

func convertToDevStatusMap(status []DevStatus) map[string]DevStatus {
	m := make(map[string]DevStatus)
	for _, s := range status {
		m[s.Name] = s
	}
	return m
}

func convertToDevStatusMapWithFilter(status []DevStatus, targetDevs map[string]bool) map[string]DevStatus {
	m := make(map[string]DevStatus)
	for _, s := range status {
		if _, has := targetDevs[s.Name]; has {
			m[s.Name] = s
		}
	}
	return m
}

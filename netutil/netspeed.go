package netutil

import (
	"time"
)

type NetSpeed struct {
	DevName       string
	ReceiveSpeed  uint64
	TransmitSpeed uint64
}

type NetSpeedStatusContext struct {
	recordTime time.Time
	statusMap  map[string]DevStatus
}

func MeasureNetSpeed(timeSpanMilli int) ([]NetSpeed, error) {
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

func CreateNetSpeedStatusContext() (*NetSpeedStatusContext, error) {
	s, err := readDevStatusDefault()
	if err != nil {
		return nil, err
	}
	return &NetSpeedStatusContext{
		recordTime: time.Now(),
		statusMap:  convertToDevStatusMap(s),
	}, nil
}

func CalculateNetSpeed(c *NetSpeedStatusContext) ([]NetSpeed, error) {
	lt := c.recordTime
	ls := c.statusMap
	t := time.Now()
	s, err := readDevStatusDefault()
	if err != nil {
		return nil, err
	}
	sMap := convertToDevStatusMap(s)
	speedList := cal(sMap, ls, int((t.UnixNano()-lt.UnixNano())/1000))
	c.recordTime = t
	c.statusMap = sMap
	return speedList, nil
}

func readDevStatusDefault() ([]DevStatus, error) {
	return ReadDevStatus(DefaultFileReader)
}

func cal(map1, map2 map[string]DevStatus, timeSpanMilli int) []NetSpeed {
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

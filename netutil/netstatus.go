package netutil

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type DevReceiveStatus struct {
	Bytes      uint64 `json:"bytes"`
	Packets    uint64 `json:"packets"`
	Errs       uint64 `json:"errs"`
	Drop       uint64 `json:"drop"`
	Fifo       uint64 `json:"fifo"`
	Frame      uint64 `json:"frame"`
	Compressed uint64 `json:"compressed"`
	Multicast  uint64 `json:"multicast"`
}

type DevTransmitStatus struct {
	Bytes      uint64 `json:"bytes"`
	Packets    uint64 `json:"packets"`
	Errs       uint64 `json:"errs"`
	Drop       uint64 `json:"drop"`
	Fifo       uint64 `json:"fifo"`
	Colls      uint64 `json:"colls"`
	Carrier    uint64 `json:"carrier"`
	Compressed uint64 `json:"compressed"`
}

type DevStatus struct {
	Name     string
	Receive  DevReceiveStatus
	Transmit DevTransmitStatus
}

type FileReaderFn func(path string) (string, error)

func DefaultFileReader(path string) (string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

func ReadDevStatus(read FileReaderFn) ([]DevStatus, error) {
	content, err := read("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(content, "\n")
	if len(lines) < 3 {
		return nil, errors.New("Could not read net device proc file")
	}
	devStatusList := make([]DevStatus, 0)
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) > 1 {
			devName := strings.TrimSpace(parts[0])
			re, _ := regexp.Compile("[ ]+")
			dataItems := strings.Split(re.ReplaceAllString(strings.TrimSpace(parts[1]), " "), " ")
			if len(dataItems) == 16 {
				rBytes, _ := strconv.ParseUint(dataItems[0], 10, 64)
				rPackets, _ := strconv.ParseUint(dataItems[1], 10, 64)
				rErrs, _ := strconv.ParseUint(dataItems[2], 10, 64)
				rDrop, _ := strconv.ParseUint(dataItems[3], 10, 64)
				rFifo, _ := strconv.ParseUint(dataItems[4], 10, 64)
				rFrame, _ := strconv.ParseUint(dataItems[5], 10, 64)
				rCompressed, _ := strconv.ParseUint(dataItems[6], 10, 64)
				rMulticast, _ := strconv.ParseUint(dataItems[7], 10, 64)
				tBytes, _ := strconv.ParseUint(dataItems[8], 10, 64)
				tPackets, _ := strconv.ParseUint(dataItems[9], 10, 64)
				tErrs, _ := strconv.ParseUint(dataItems[10], 10, 64)
				tDrop, _ := strconv.ParseUint(dataItems[11], 10, 64)
				tFifo, _ := strconv.ParseUint(dataItems[12], 10, 64)
				tColls, _ := strconv.ParseUint(dataItems[13], 10, 64)
				tCarrier, _ := strconv.ParseUint(dataItems[14], 10, 64)
				tCompressed, _ := strconv.ParseUint(dataItems[15], 10, 64)
				devStatusList = append(devStatusList, DevStatus{
					Name: devName,
					Receive: DevReceiveStatus{
						Bytes:      rBytes,
						Packets:    rPackets,
						Errs:       rErrs,
						Drop:       rDrop,
						Fifo:       rFifo,
						Frame:      rFrame,
						Compressed: rCompressed,
						Multicast:  rMulticast,
					},
					Transmit: DevTransmitStatus{
						Bytes:      tBytes,
						Packets:    tPackets,
						Errs:       tErrs,
						Drop:       tDrop,
						Fifo:       tFifo,
						Colls:      tColls,
						Carrier:    tCarrier,
						Compressed: tCompressed,
					},
				})
			}
		}
	}
	return devStatusList, nil
}

package bridge

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"../base"
)

type Bridge struct {
	cfg    *ba.Config
	execFn ba.CmdExecutorFn
}

func NewBridge(fn ba.CmdExecutorFn, c *ba.Config) *Bridge {
	return &Bridge{
		cfg:    c,
		execFn: fn,
	}
}

func (b *Bridge) SetupBridge() error {
	bridgeList, err := readBridge(b.execFn)
	if err != nil {
		return err
	}
	bridge := findBridge(bridgeList, b.cfg.BridgeName)
	if bridge != nil {
		if isCurrentSettingValid(bridge, b.cfg) {
			return nil
		}
		dealDeleteBridge(b.execFn, bridge)
	}
	err = dealCreateBridge(b.execFn, b.cfg)
	if err != nil {
		return err
	}
	log.Println("bridge setted up")
	return nil
}

func (b *Bridge) ClearBridge() error {
	bridgeList, err := readBridge(b.execFn)
	if err != nil {
		return err
	}
	bridge := findBridge(bridgeList, b.cfg.BridgeName)
	if bridge != nil {
		dealDeleteBridge(b.execFn, bridge)
	}
	log.Println("bridge cleared")
	return nil
}

type BridgeInfo struct {
	name   string
	ifList []string
}

func findBridge(bridgeList []BridgeInfo, brName string) *BridgeInfo {
	for _, b := range bridgeList {
		if b.name == brName {
			return &b
		}
	}
	return nil
}

func dealCreateBridge(execFn ba.CmdExecutorFn, c *ba.Config) error {
	_, err := execFn(fmt.Sprint("brctl addbr ", c.BridgeName))
	if err != nil {
		return err
	}
	_, err = execFn(fmt.Sprint("ifconfig ", c.BridgeName, " ", c.BridgeAddr, " netmask 255.255.255.0 up"))
	if err != nil {
		return err
	}
	for _, ifName := range c.LanNames {
		_, err := execFn(fmt.Sprint("brctl addif ", c.BridgeName, " ", ifName))
		if err != nil {
			return err
		}
	}
	return nil
}

func dealDeleteBridge(execFn ba.CmdExecutorFn, info *BridgeInfo) error {
	_, err := execFn(fmt.Sprint("ifconfig ", info.name, " down"))
	if err != nil {
		return err
	}
	_, err = execFn(fmt.Sprint("brctl delbr ", info.name))
	if err != nil {
		return err
	}
	return nil
}

func isCurrentSettingValid(info *BridgeInfo, c *ba.Config) bool {
	if info.name != c.BridgeName {
		return false
	}
	if len(info.ifList) != len(c.LanNames) {
		return false
	}
	for i := 0; i < len(info.ifList); i++ {
		if info.ifList[i] != c.LanNames[i] {
			return false
		}
	}
	return true
}

func readBridge(execFn ba.CmdExecutorFn) ([]BridgeInfo, error) {
	content, err := execFn("brctl show")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(content, "\n")
	if len(lines) < 3 {
		return []BridgeInfo{}, nil
	}
	lines = lines[1 : len(lines)-1]
	re, _ := regexp.Compile("\\s+")
	bridgeList := []BridgeInfo{}
	for _, line := range lines {
		cols := strings.Split(re.ReplaceAllString(line, " "), " ")
		if len(cols) == 4 {
			bridgeList = append(bridgeList, BridgeInfo{
				name:   cols[0],
				ifList: []string{cols[3]},
			})
		} else if len(cols) == 2 && len(bridgeList) > 0 {
			lastBridge := &bridgeList[len(bridgeList)-1]
			lastBridge.ifList = append(lastBridge.ifList, cols[1])
		}
	}
	return bridgeList, nil
}

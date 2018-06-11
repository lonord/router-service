package bridge

import (
	"fmt"
	"regexp"
	"strings"

	"../context"
	"../util"
)

func SetupBridge(execFn util.CmdExecutorFn, c *context.MainContext) error {
	bridgeList, err := readBridge(execFn)
	if err != nil {
		return err
	}
	bridge := findBridge(bridgeList, c.Cfg.BridgeName)
	if bridge != nil {
		if isCurrentSettingValid(bridge, c) {
			return nil
		}
		dealDeleteBridge(execFn, bridge)
	}
	err = dealCreateBridge(execFn, c)
	if err != nil {
		return err
	}
	return nil
}

func ClearBridge(execFn util.CmdExecutorFn, c *context.MainContext) error {
	bridgeList, err := readBridge(execFn)
	if err != nil {
		return err
	}
	bridge := findBridge(bridgeList, c.Cfg.BridgeName)
	if bridge != nil {
		dealDeleteBridge(execFn, bridge)
	}
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

func dealCreateBridge(execFn util.CmdExecutorFn, c *context.MainContext) error {
	_, err := execFn(fmt.Sprint("brctl addbr ", c.Cfg.BridgeName))
	if err != nil {
		return err
	}
	_, err = execFn(fmt.Sprint("ifconfig ", c.Cfg.BridgeName, " ", c.Cfg.BridgeAddr, " netmask 255.255.255.0 up"))
	if err != nil {
		return err
	}
	for _, ifName := range c.Cfg.LanNames {
		_, err := execFn(fmt.Sprint("brctl addif ", c.Cfg.BridgeName, " ", ifName))
		if err != nil {
			return err
		}
	}
	return nil
}

func dealDeleteBridge(execFn util.CmdExecutorFn, info *BridgeInfo) error {
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

func isCurrentSettingValid(info *BridgeInfo, c *context.MainContext) bool {
	if info.name != c.Cfg.BridgeName {
		return false
	}
	if len(info.ifList) != len(c.Cfg.LanNames) {
		return false
	}
	for i := 0; i < len(info.ifList); i++ {
		if info.ifList[i] != c.Cfg.LanNames[i] {
			return false
		}
	}
	return true
}

func readBridge(execFn util.CmdExecutorFn) ([]BridgeInfo, error) {
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

package bridge

import (
	"reflect"
	"testing"
)

func TestReadBridge(t *testing.T) {
	bridgeList, err := readBridge(executeCmd)
	expect(t, err, nil)
	expect(t, len(bridgeList), 2)
	expect(t, bridgeList[0].name, "br-lan")
	expect(t, bridgeList[1].name, "docker0")
	expect(t, len(bridgeList[0].ifList), 2)
	expect(t, bridgeList[0].ifList[0], "enx000ec6d7e74f")
	expect(t, bridgeList[0].ifList[1], "enx00ed4d680193")
	expect(t, len(bridgeList[1].ifList), 1)
	expect(t, bridgeList[1].ifList[0], "veth2648ffe")
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(b, a) {
		t.Errorf("Expected %#v (type %v) - Got %#v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func executeCmd(cmd string) (string, error) {
	if cmd == "brctl show" {
		return outContent, nil
	}
	return "", nil
}

const outContent = `bridge name	bridge id		STP enabled	interfaces
br-lan		8000.000ec6d7e74f	no		enx000ec6d7e74f
							enx00ed4d680193
docker0		8000.024213f65672	no		veth2648ffe
`

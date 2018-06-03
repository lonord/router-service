package netutil

import (
	"reflect"
	"testing"
)

func TestReadDevStatus(t *testing.T) {
	result, err := ReadDevStatus(func(path string) (string, error) {
		return getTestFileContent(), nil
	})
	expect(t, err, nil)
	expect(t, len(result), 7)
	resultItem := result[6]
	expect(t, resultItem.Name, "enx00ed4d680193")
	expect(t, resultItem.Receive.Bytes, uint64(1654931333))
	expect(t, resultItem.Transmit.Bytes, uint64(1285412852))
	expect(t, resultItem.Transmit.Packets, uint64(16952848))
	expect(t, resultItem.Transmit.Compressed, uint64(0))
}

func getTestFileContent() string {
	return `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
br-lan: 18811839434 11783734    0    0    0     0          0         0 13511448597 9449594    0    0    0     0       0          0
 wlan0:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0
    lo: 657147266  579787    0    0    0     0          0         0 657147266  579787    0    0    0     0       0          0
docker0: 12078457115 10240430    0    0    0     0          0         0 9486769101 10571301    0    0    0     0       0          0
enxb827eb1f34f3: 1788010491 27006046   58 2927199    0    17          0         0 69735478 18145712    0    0    0     0       0          0
veth2648ffe: 145800525 1362490    0    0    0     0          0         0 2014120660 1617866    0    0    0     0       0          0
enx00ed4d680193: 1654931333 29387889    0    0    0     0          0         0 1285412852 16952848    0    0    0     0       0          0
`
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(b, a) {
		t.Errorf("Expected %#v (type %v) - Got %#v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

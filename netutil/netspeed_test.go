package netutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCal(t *testing.T) {
	s1, err := ReadDevStatus(func(p string) (string, error) {
		return data1, nil
	})
	assert.NoError(t, err)
	s2, err := ReadDevStatus(func(p string) (string, error) {
		return data2, nil
	})
	assert.NoError(t, err)
	map1 := convertToDevStatusMap(s1)
	map2 := convertToDevStatusMap(s2)
	speeds := cal(map1, map2, 1000)
	assert.Equal(t, len(speeds), 3)
	var speed NetSpeed
	for _, ss := range speeds {
		if ss.DevName == "br-lan" {
			speed = ss
		}
	}
	assert.Equal(t, speed.DevName, "br-lan")
	assert.Equal(t, speed.ReceiveSpeed, uint64(1000))
	assert.Equal(t, speed.TransmitSpeed, uint64(2000))
}

const data1 = `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
br-lan: 88232992879 78760799    0    0    0     0          0         0 99798877563 67202989    0    0    0     0       0          0
 wlan0:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0
    lo: 2312260588 3252792    0    0    0     0          0         0 2312260588 3252792    0    0    0     0       0          0
`

const data2 = `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
br-lan: 88232993879 78760799    0    0    0     0          0         0 99798879563 67202989    0    0    0     0       0          0
 wlan0:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0
    lo: 2312260588 3252792    0    0    0     0          0         0 2312260588 3252792    0    0    0     0       0          0
`

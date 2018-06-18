package ba

import "strings"

func GetSubnetPrefix(ipAddr string) string {
	ipChunks := strings.Split(ipAddr, ".")
	return strings.Join([]string{ipChunks[0], ipChunks[1], ipChunks[2]}, ".")
}

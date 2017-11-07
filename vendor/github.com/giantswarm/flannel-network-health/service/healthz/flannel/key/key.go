package key

import (
	"github.com/vishvananda/netlink"
)

func GetInterfaceIP(ipList []netlink.Addr) string {
	if len(ipList) == 0 {
		return ""
	}
	return ipList[0].IP.String()
}

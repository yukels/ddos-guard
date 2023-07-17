package system

import (
	"net"

	"github.com/yukels/util/context"
	"github.com/yukels/util/log"
)

// OutboundIP return preferred outbound ip of this machine
func OutboundIP() string {
	ctx := context.Background()
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Log(ctx).WithError(err).Errorf("Can't get outbound server IP")
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

package common

import (
	"net"
	"strings"

	"github.com/fananchong/v-micro/common/log"
	maddr "github.com/fananchong/v-micro/internal/addr"
	"github.com/fananchong/v-micro/server"
)

// ExtractAdvertiseIP When Advertise is "0.0.0.0", "[::]", "::", the network segment in PrivateIPBlocks is preferentially matched. PrivateIPBlocks default value is 192.168.0.0/16
func ExtractAdvertiseIP(config server.Options) (addr string, port string, err error) {
	var advt, host string
	// check the advertise address first
	// if it exists then use it, otherwise
	// use the address
	if len(config.Advertise) > 0 {
		advt = config.Advertise
	} else {
		advt = config.Address
	}

	if cnt := strings.Count(advt, ":"); cnt >= 1 {
		// ipv6 address in format [host]:port or ipv4 host:port
		host, port, err = net.SplitHostPort(advt)
		if err != nil {
			log.Error(err)
			return
		}
	} else {
		host = advt
	}

	ipBlocks := "192.168.0.0/16"
	if len(config.PrivateIPBlocks) != 0 {
		ipBlocks = config.PrivateIPBlocks
	}
	addr, err = maddr.Extract(host, ipBlocks)
	if err != nil {
		log.Error(err)
	}

	return
}

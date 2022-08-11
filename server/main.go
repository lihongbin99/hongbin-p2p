package main

import (
	"hongbin-p2p/common/log"
	"net"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:13520")
	if err != nil {
		logger.Err("resolve error", err)
		return
	}
	logger.Info("resolve success")

	conn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		logger.Err("listen error", err)
		return
	}
	logger.Info("listen success")

	buf := make([]byte, 1024*64)
	var lastUserAddr *net.UDPAddr
	for {
		_, userAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			logger.Err("read from error", err)
			continue
		}
		if lastUserAddr == nil {
			lastUserAddr = userAddr
			logger.Info("save addr: %s", lastUserAddr.String())
		} else {
			if _, err := conn.WriteToUDP([]byte(userAddr.String()), lastUserAddr); err != nil {
				logger.Err("write to error, addr: %s", lastUserAddr.String(), err)
			}
			if _, err := conn.WriteToUDP([]byte(lastUserAddr.String()), userAddr); err != nil {
				logger.Err("write to error, addr: %s", userAddr.String(), err)
			}
			logger.Info("exchange [%s] -> [%s]", lastUserAddr.String(), userAddr.String())
			lastUserAddr = nil
		}
	}
}

package main

import (
	logger "hongbin-p2p/common/log"
	"net"
)

func main() {
	//udpServer()
	udpClient()
}

func udpServer() {
	udpAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:3390")
	if err != nil {
		logger.Err("resolve server error", err)
		return
	}
	logger.Info("resolve server success")

	conn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		logger.Err("listen server error", err)
		return
	}
	logger.Info("listen server success")

	buf := make([]byte, 1024*64)
	for {
		readLength, userAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			logger.Err("read from client error", err)
			continue
		}
		logger.Info("read from client success: %s", string(buf[:readLength]))
		writeLength, err := conn.WriteToUDP([]byte("Im server"), userAddr)
		if err != nil {
			logger.Err("write to client error", err)
			continue
		}
		logger.Err("write to client success: %d", writeLength)
	}
}

func udpClient() {
	remoteAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:9090")
	if err != nil {
		logger.Err("resolve target error", err)
		return
	}
	logger.Info("resolve target success")
	targetConn, err := net.DialUDP("udp4", nil, remoteAddr)
	if err != nil {
		logger.Err("dial target error", err)
		return
	}

	writeLength, err := targetConn.Write([]byte("Im client"))
	if err != nil {
		logger.Err("write to server error", err)
		return
	}
	logger.Err("write to server success: %d", writeLength)

	buf := make([]byte, 1024*64)
	readLength, _, err := targetConn.ReadFromUDP(buf)
	if err != nil {
		logger.Err("read from server error", err)
		return
	}
	logger.Info("read from server success: %s", string(buf[:readLength]))
}

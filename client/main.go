package main

import (
	"hongbin-p2p/common/convert"
	"hongbin-p2p/common/ioutil"
	"hongbin-p2p/common/log"
	"hongbin-p2p/common/transfer"
	"net"
	"time"
)

var (
	runType = "server"
	//runType = "client"
	prefix = 5
)

func main() {
	natMessageChan := make(chan ioutil.Message)
	transferMessageChan := make(chan ioutil.Message)
	if runType == "client" {
		udpAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:9090")
		if err != nil {
			logger.Err("resolve client error", err)
			return
		}
		logger.Info("resolve client success")

		conn, err := net.ListenUDP("udp4", udpAddr)
		if err != nil {
			logger.Err("listen client error", err)
			return
		}
		logger.Info("listen client success")

		clientMap := make(map[int32]*net.UDPAddr)
		clientAddrMap := make(map[string]int32)

		go func() {
			for {
				message := <-transferMessageChan
				if userAddr, ok := clientMap[message.ConnId]; ok {
					writeLength, err := conn.WriteToUDP(message.Buf[:message.ReadLength], userAddr)
					if err != nil {
						logger.Err("write to client error, addr: %s", userAddr.String(), err)
					} else {
						logger.Info("write to client success, write length: %d->%d", message.ReadLength, writeLength)
					}
				} else {
					logger.Error("no find coinId: %d", message.ConnId)
				}
				transferMessageChan <- ioutil.Message{}
			}
		}()

		go func() {
			buf := make([]byte, 1024*64)
			buf[0] = transfer.Transfer
			var connIdIndex int32 = 1
			for {
				readLength, userAddr, err := conn.ReadFromUDP(buf[prefix:])
				if err != nil {
					logger.Err("read from client error", err)
					continue
				}
				totalLength := readLength + prefix
				if connId, ok := clientAddrMap[userAddr.String()]; ok {
					copy(buf[1:], convert.I2b(connId))
					natMessageChan <- ioutil.Message{ConnId: connId, Buf: buf[:totalLength], ReadLength: totalLength}
				} else {
					connIdIndex++
					clientAddrMap[userAddr.String()] = connIdIndex
					clientMap[connIdIndex] = userAddr
					copy(buf[1:], convert.I2b(connIdIndex))
					natMessageChan <- ioutil.Message{ConnId: connIdIndex, Buf: buf[:totalLength], ReadLength: totalLength}
				}
				<-natMessageChan
			}
		}()
	} else if runType == "server" {
		targetMap := make(map[int32]*net.UDPConn)
		go func() {
			for {
				message := <-transferMessageChan
				logger.Info("read message from remote, connId: %d, message length: %d", message.ConnId, message.ReadLength)
				targetConn, ok := targetMap[message.ConnId]
				if !ok {
					logger.Info("no find connIdL %d, start dial to target", message.ConnId)
					remoteAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:3389")
					if err != nil {
						logger.Err("resolve target error", err)
						return
					}
					logger.Info("resolve target success")
					targetConn, err = net.DialUDP("udp4", nil, remoteAddr)
					if err != nil {
						logger.Err("dial target error", err)
						return
					}
					logger.Info("dial target success")
					targetMap[message.ConnId] = targetConn
					go func(connId int32, tConn *net.UDPConn) {
						buf := make([]byte, 1024*64)
						buf[0] = transfer.Transfer
						copy(buf[1:], convert.I2b(connId))
						for {
							readLength, _, err := tConn.ReadFromUDP(buf[prefix:])
							if err != nil {
								logger.Err("read from target error", err)
								continue
							}
							totalLength := readLength + prefix
							natMessageChan <- ioutil.Message{ConnId: connId, Buf: buf[:totalLength], ReadLength: totalLength}
							<-natMessageChan
						}
					}(message.ConnId, targetConn)
				}
				writeLength, err := targetConn.Write(message.Buf[:message.ReadLength])
				if err != nil {
					logger.Err("write to target error", err)
				} else {
					logger.Info("write to target success, write length: %d->%d", message.ReadLength, writeLength)
				}
				transferMessageChan <- ioutil.Message{}
			}
		}()
	} else {
		logger.Error("run type error, run type: %s", runType)
		return
	}

	buf := make([]byte, 1024*64)

	//serverAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:13520")
	serverAddr, err := net.ResolveUDPAddr("udp4", "43.128.70.137:13520")
	if err != nil {
		logger.Err("resolve server error", err)
		return
	}
	logger.Info("resolve server success")

	conn, err := net.DialUDP("udp4", nil, serverAddr)
	if err != nil {
		logger.Err("dial error", err)
		return
	}
	logger.Info("dial success")

	localAddrS := conn.LocalAddr().String()
	logger.Info("local addr: %s", localAddrS)

	localAddr, err := net.ResolveUDPAddr("udp", localAddrS)
	if err != nil {
		logger.Err("resolve local error", err)
		return
	}
	logger.Info("resolve local success")

	_, err = conn.Write(buf[:1])
	if err != nil {
		logger.Err("write to server error", err)
		return
	}
	logger.Info("write to server success")

	readLength, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		logger.Err("read from server error", err)
		return
	}
	if readLength <= 0 {
		logger.Err("read from server error, readLength: %d", readLength, err)
		return
	}
	logger.Info("read from server success, readLength: %d", readLength)

	remoteAddrS := string(buf[:readLength])
	logger.Info("remote addr: %s", remoteAddrS)

	remoteAddr, err := net.ResolveUDPAddr("udp", remoteAddrS)
	if err != nil {
		logger.Err("resolve remote error", err)
		return
	}
	logger.Info("resolve remote success")

	conn, err = net.DialUDP("udp4", localAddr, remoteAddr)
	if err != nil {
		logger.Err("dial error", err)
		return
	}
	logger.Info("dial success")

	messageChan := make(chan ioutil.Message)

	go func() {
		for {
			readLength, _, err = conn.ReadFromUDP(buf)
			if err != nil {
				messageChan <- ioutil.Message{Err: err}
				break
			}
			if readLength >= prefix {
				messageChan <- ioutil.Message{ReadLength: readLength - prefix, TransferType: buf[0], ConnId: convert.B2i(buf[1:prefix]), Buf: buf[prefix:]}
			} else {
				messageChan <- ioutil.Message{ReadLength: readLength, TransferType: transfer.Error, Buf: buf}
			}
			<-messageChan
		}
	}()

	pingTicker := time.NewTicker(3 * time.Second)
	pingBuf := make([]byte, 19+prefix)
	pingBuf[0] = transfer.Ping
	lastPingTime := time.Now()

	for {
		select {
		case pingTime := <-pingTicker.C:
			_ = time.Now().Sub(lastPingTime) // TODO KeepAlice
			copy(pingBuf[prefix:], []byte(pingTime.String())[:19])
			if _, err = conn.Write(pingBuf); err != nil {
				logger.Err("write ping to remote error", err)
				return
			}
			logger.Info("write ping to remote success")
		case message := <-messageChan:
			if message.Err != nil {
				logger.Err("read from remote error", message.Err)
				return
			}
			switch message.TransferType {
			case transfer.Error:
				logger.Error("transfer type error, readLength: %d, buf[:prefix]: %v", message.ReadLength, message.Buf[:prefix])
			case transfer.Ping:
				lastPingTime = time.Now()
				logger.Info("read ping from remote success, time: %s", string(message.Buf[:message.ReadLength]))
			case transfer.Transfer:
				transferMessageChan <- message
				<-transferMessageChan
			default:
				logger.Error("read from remote transfer type error, transfer type: %d", message.TransferType)
			}
			messageChan <- ioutil.Message{}
		case message := <-natMessageChan:
			writeLength, err := conn.Write(message.Buf)
			if err != nil {
				logger.Err("write message to remote error, connId: %d", message.ConnId, err)
				return
			}
			logger.Info("write message to remote success, connId: %d, write length: %d->%d", message.ConnId, message.ReadLength, writeLength)
			natMessageChan <- ioutil.Message{}
		}
	}
}

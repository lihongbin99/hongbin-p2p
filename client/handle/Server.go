package handle

import (
	"hongbin-p2p/common/ioutil"
	"net"
)

type Server struct {
	clientMap map[int32]net.UDPConn
}

func (server *Server) Start() {

}

func (server *Server) HandleTransfer(message ioutil.Message) {

}

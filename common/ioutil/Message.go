package ioutil

type Message struct {
	ReadLength   int
	TransferType byte
	ConnId       int32
	Buf          []byte
	Err          error
}

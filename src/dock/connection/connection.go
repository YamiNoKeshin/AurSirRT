package connection


type Connection interface {
	Init() (err error)
	Send(msgtype int64, codec string,msg []byte) (err error)
	Close() error
}


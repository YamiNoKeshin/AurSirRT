package dockzmq

import (
	"strconv"
	zmq "github.com/pebbe/zmq4"
	"fmt"
)


type ConnectionZmq struct {
	myport int64
	port int64
	myip string
	ip string
	id string
	skt *zmq.Socket
}

func NewConnection(homeport, port int64, localip,targetip,id string) ConnectionZmq{
	return ConnectionZmq{homeport, port,localip,targetip,id, nil}
}



func (cz *ConnectionZmq) Init() (err error) {
	cz.skt, _ = zmq.NewSocket(zmq.DEALER)
	cz.skt.SetIdentity(cz.id)

	cz.skt.SetSndtimeo(1000)
	port := strconv.FormatInt(cz.port,10)
	printDebug("ZMQAppDocker opening channel to ip "+cz.ip)
	printDebug("ZMQAppDocker opening channel to port "+port)

	err = cz.skt.Connect(fmt.Sprintf("tcp://%s:%d",cz.ip, cz.port))

	return

}

func (cz ConnectionZmq) Send(msgtype int64, codec string,msg []byte) (err error){
	_,err = cz.skt.SendMessage(
		[]string{strconv.FormatInt(msgtype,10),codec,string(msg), strconv.FormatInt(cz.myport,10),cz.myip},0)

	if err != nil {
		mprint(fmt.Sprintf("Error on zqm port %d, closing:",cz.port,err))
	}
	return
}
func (cz ConnectionZmq) Close() (err error){
	err = cz.skt.Close()
	return
}

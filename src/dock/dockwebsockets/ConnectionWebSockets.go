package dockwebsockets

import (
	"code.google.com/p/go.net/websocket"
	"strconv"
	"log"
)

type ConnectionWebSockets struct {

	ws *websocket.Conn

}

func NewConnection(ws *websocket.Conn) ConnectionWebSockets {
	return ConnectionWebSockets{ws}
}

func (cw ConnectionWebSockets) Init() (err error) {

	return
}

func (cw ConnectionWebSockets) Send(msgtype int64, codec string,msg []byte) (err error) {
	err = websocket.Message.Send(cw.ws,strconv.FormatInt(msgtype,10))
	log.Println("WEBSOCET1",err)
	if err != nil {
		return
	}
	err = websocket.Message.Send(cw.ws,codec)
	log.Println("WEBSOCET2",err)
	if err != nil {
		return
	}
	err = websocket.Message.Send(cw.ws,string(msg))
	log.Println("WEBSOCET3",err)
	return
}

func (ConnectionWebSockets) Close() (err error){
	return
	
}

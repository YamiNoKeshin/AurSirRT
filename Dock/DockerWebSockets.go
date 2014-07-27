package dock

import (
	"github.com/joernweissenborn/aursirrt/core"
	"log"
	"code.google.com/p/go.net/websocket"
	"net/http"
	"io"
	"strconv"
	"github.com/joernweissenborn/AurSir4Go"
)

type DockerWebSockets struct {
	appInChan chan core.AppMessage
	regChan chan registerDockedApp
	port string
}

func (dws DockerWebSockets)	Launch(appInChan chan core.AppMessage, regApp chan registerDockedApp){

	log.Println("DockerWebsockets Launching")

	dws.appInChan = appInChan
	dws.regChan = regApp
	dws.port = "8086"

	go dws.server()
}

func (dws DockerWebSockets) server(){
	http.Handle("/aursirrt", websocket.Handler(dws.onConnect))
	//http.Handle("/", http.FileServer(http.Dir(".")))
	log.Println("DockerWebsockets starting to listen on port:",dws.port)
	err := http.ListenAndServe(":"+dws.port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}

func (dws *DockerWebSockets) onConnect(ws *websocket.Conn) {
	log.Println("DockerWebsockets openening ws from:",ws.RemoteAddr())
	defer log.Println("Closing Websocket to",ws.RemoteAddr())
	defer ws.Close()
	dws.listen(ws)
}

func (dws DockerWebSockets) listen(ws *websocket.Conn){

	log.Println("DockerWebsockets starting listening to:",ws.RemoteAddr())

	closed := make(chan struct{})

	for {
		eba := []byte("{}")

		msgtype, err := receiveMsg(ws)

		senderId :=ws.RemoteAddr().String()
		if err == io.EOF {
			log.Println("DockerWebsockets got EOF on client", senderId)
			dws.appInChan <- core.AppMessage{senderId,AurSir4Go.AppMessage{AurSir4Go.LEAVE,"JSON",&eba}}
			return
		} else if err != nil {
			log.Println("DockerWebsockets got error on client %x:", ws.RemoteAddr(), err)
			dws.appInChan <- core.AppMessage{senderId,AurSir4Go.AppMessage{AurSir4Go.LEAVE,"JSON",&eba}}
			return
		}
		msgCodec, err := receiveMsg(ws)
		if err == io.EOF {
			log.Println("DockerWebsockets got EOF on client", ws.RemoteAddr())
			dws.appInChan <- core.AppMessage{senderId,AurSir4Go.AppMessage{AurSir4Go.LEAVE,"JSON",&eba}}
			closed <- struct{}{}
			return
		} else if err != nil {
			log.Println("DockerWebsockets got error on client %x:", ws.RemoteAddr(), err)
			dws.appInChan <- core.AppMessage{senderId,AurSir4Go.AppMessage{AurSir4Go.LEAVE,"JSON",&eba}}
			return
		}
		msgBytes, err := receiveMsg(ws)
		if err == io.EOF {
			log.Println("DockerWebsockets got EOF on client", ws.RemoteAddr())
			dws.appInChan <- core.AppMessage{senderId,AurSir4Go.AppMessage{AurSir4Go.LEAVE,"JSON",&eba}}
			closed <- struct{}{}
			return
		} else if err != nil {
			closed <- struct{}{}
			log.Println("DockerWebsockets got error on client %x:", ws.RemoteAddr(), err)
			dws.appInChan <- core.AppMessage{senderId,AurSir4Go.AppMessage{AurSir4Go.LEAVE,"JSON",&eba}}

			return
		}
		msgType, err := strconv.ParseInt(string((*msgtype)[0]),10,8)
		if err != nil {
			log.Println("DockerWebsockets got invalid Message on client:", ws.RemoteAddr())
			return
		}

		if  msgType==AurSir4Go.DOCK{
			go dws.openConnection(ws,closed)
		}

		go dws.processMsg(ws.RemoteAddr().String(),msgType,msgCodec,msgBytes)
	}

}

func receiveMsg( ws *websocket.Conn) (*[]byte,error){
	var msg []byte
	err := websocket.Message.Receive(ws, &msg)
	return &msg,err
}

func (dws DockerWebSockets) processMsg(senderId string,msgType int64,msgCodec *[]byte,msgBytes *[]byte){

	codec := string((*msgCodec))

	dws.appInChan <- core.AppMessage{senderId,AurSir4Go.AppMessage{msgType,codec,msgBytes}}

	log.Println(msgBytes)

}

func (dws DockerWebSockets) openConnection(ws *websocket.Conn, closed chan struct{}){
	c := make(chan core.AppMessage )

	dws.regChan <- registerDockedApp{ws.RemoteAddr().String(), c}
	log.Println("DockerWebsockets opening outgoing channel to:", ws.RemoteAddr())

	for {

		select {
		case msg, _ := <- c:
			appmsg := (msg.AppMsg)
			log.Println(*appmsg.Msg)
			websocket.Message.Send(ws,strconv.FormatInt(appmsg.MsgType,10))
			websocket.Message.Send(ws,appmsg.MsgCodec)
			websocket.Message.Send(ws,*appmsg.Msg)

		case <- closed:
			return
		}


	}

}

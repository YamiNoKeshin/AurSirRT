package dock

import (
	uuid "github.com/nu7hatch/gouuid"

	"github.com/joernweissenborn/aursirrt/core"
	"log"
	"code.google.com/p/go.net/websocket"
	"net/http"
	"io"
	"strconv"
	"github.com/joernweissenborn/aursir4go"
)

type DockerWebSockets struct {
	appInChan chan core.AppMessage
	regChan chan interface {}
	port string
	uuid string
}

func (dws DockerWebSockets)	Launch(appInChan chan core.AppMessage, regApp chan interface {}){

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
	dws.uuid = generateUuid()
	log.Println("DockerWebsockets assinging uuid ",dws.uuid)

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
		senderId :=   dws.uuid
		if err == io.EOF {
			log.Println("DockerWebsockets got EOF on client", senderId)
			dws.appInChan <- core.AppMessage{senderId,aursir4go.AppMessage{aursir4go.LEAVE,"JSON",eba}}
			return
		} else if err != nil {
			log.Println("DockerWebsockets got error on client %x:", ws.RemoteAddr(), err)
			dws.appInChan <- core.AppMessage{senderId,aursir4go.AppMessage{aursir4go.LEAVE,"JSON",eba}}
			return
		}
		msgCodec, err := receiveMsg(ws)
		if err == io.EOF {
			log.Println("DockerWebsockets got EOF on client", ws.RemoteAddr())
			dws.appInChan <- core.AppMessage{senderId,aursir4go.AppMessage{aursir4go.LEAVE,"JSON",eba}}
			closed <- struct{}{}
			return
		} else if err != nil {
			log.Println("DockerWebsockets got error on client %x:", ws.RemoteAddr(), err)
			dws.appInChan <- core.AppMessage{senderId,aursir4go.AppMessage{aursir4go.LEAVE,"JSON",eba}}
			return
		}
		msgBytes, err := receiveMsg(ws)
		if err == io.EOF {
			log.Println("DockerWebsockets got EOF on client", ws.RemoteAddr())
			dws.appInChan <- core.AppMessage{senderId,aursir4go.AppMessage{aursir4go.LEAVE,"JSON",eba}}
			closed <- struct{}{}
			return
		} else if err != nil {
			closed <- struct{}{}
			log.Println("DockerWebsockets got error on client %x:", ws.RemoteAddr(), err)
			dws.appInChan <- core.AppMessage{senderId,aursir4go.AppMessage{aursir4go.LEAVE,"JSON",eba}}

			return
		}
		msgType, err := strconv.ParseInt(string((*msgtype)),10,64)


		if err != nil {
			log.Println("DockerWebsockets got invalid Message on client:", ws.RemoteAddr())
			return
		}

		if  msgType==aursir4go.DOCK{
			go dws.openConnection(ws,closed)
		}

		go dws.processMsg(senderId,msgType,msgCodec,msgBytes)
	}

}

func receiveMsg( ws *websocket.Conn) (*[]byte,error){
	var msg []byte
	err := websocket.Message.Receive(ws, &msg)
	return &msg,err
}

func (dws DockerWebSockets) processMsg(senderId string,msgType int64,msgCodec *[]byte,msgBytes *[]byte){

	codec := string((*msgCodec))

	dws.appInChan <- core.AppMessage{senderId,aursir4go.AppMessage{msgType,codec,*msgBytes}}


}

func (dws DockerWebSockets) openConnection(ws *websocket.Conn, closed chan struct{}){
	c := make(chan core.AppMessage )

	dws.regChan <- registerDockedApp{dws.uuid, c}
	log.Println("DockerWebsockets opening outgoing channel to:",dws.uuid, ws.RemoteAddr())

	for {

		select {
		case msg, _ := <- c:
			appmsg := (msg.AppMsg)
			websocket.Message.Send(ws,strconv.FormatInt(appmsg.MsgType,10))
			websocket.Message.Send(ws,appmsg.MsgCodec)
			websocket.Message.Send(ws,string(appmsg.Msg))

		case <- closed:
			return
		}


	}

}

func generateUuid() string {
	Uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal("Failed to generate UUID")
		return ""
	}
	return Uuid.String()
}

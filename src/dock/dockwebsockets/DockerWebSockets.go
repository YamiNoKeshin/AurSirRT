package dockwebsockets

import (
	"log"
	"code.google.com/p/go.net/websocket"
	"net/http"
	"strconv"
	"aursirrt/src/dock"
	"fmt"
	"github.com/joernweissenborn/aursir4go/messages"
)
type DockerWebSockets struct {

	agent dock.DockAgent
	port string

}

func (dws DockerWebSockets)	Launch(agent dock.DockAgent, id string)(err error){

	mprint("Launching")

	dws.agent = agent
	dws.port = "8086"

	go dws.server()
	return
}

func (dws DockerWebSockets) server(){
	http.Handle("/aursirrt", websocket.Handler(dws.onConnect))
	mprint(fmt.Sprint("starting to listen on port:",dws.port))
	err := http.ListenAndServe(":"+dws.port, nil)
	if err != nil {
		mprint("Error starting server " + err.Error())
	}

}

func (dws *DockerWebSockets) onConnect(ws *websocket.Conn) {
	printDebug(fmt.Sprintf("DockerWebsockets openening ws from:",ws.RemoteAddr()))
	defer printDebug(fmt.Sprintf("Closing Websocket to",ws.RemoteAddr()))
	defer ws.Close()
	dws.listen(ws)
}

func (dws DockerWebSockets) listen(ws *websocket.Conn){

	printDebug(fmt.Sprintf("DockerWebsockets starting listening to:",ws.RemoteAddr()))

	for {
		sndId, err :=receiveMsg(ws)
		senderId := string(sndId)
		if err != nil {
			dws.remove(senderId)
			return
		}
		msgtype, err := receiveMsg(ws)

		if err != nil {
			dws.remove(senderId)
			return
		}
		msgCodec, err := receiveMsg(ws)
		if err != nil {
			dws.remove(senderId)
			return
		}
		msgBytes, err := receiveMsg(ws)
		if err != nil {
			dws.remove(senderId)
			return
		}
		msgType, err := strconv.ParseInt(string(msgtype),10,64)
		
		if err != nil {
			printDebug(fmt.Sprint("DockerWebsockets got invalid Message on client:", ws.RemoteAddr(),err))
			return
		}
		if  msgType==messages.DOCK{
			conn := NewConnection(ws)
			dws.agent.InitDocking(senderId,string(msgCodec),msgBytes,conn)
		} else {
			dws.agent.ProcessMsg(senderId, msgType, string(msgCodec), msgBytes)
		}
	}

}

func receiveMsg( ws *websocket.Conn) ([]byte,error){
	var msg []byte
	err := websocket.Message.Receive(ws, &msg)
	return msg,err
}

func (dws DockerWebSockets) remove(senderId string){
	printDebug("DockerWebsockets got EOF or error on client" + senderId)
	dws.agent.ProcessMsg(senderId,messages.LEAVE,"JSON",[]byte("{}"))
	return
}

func mprint(msg string){
	log.Println("DOCKER WEBSOCKETS", msg)
}

func printDebug(msg string){
	log.Println("DEBUG","DOCKER WEBSOCKETS", msg)
}

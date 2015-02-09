package dockzmq

import (
	"log"
	zmq "github.com/pebbe/zmq4"
	"strconv"
	"aursirrt/src/dock"
	"github.com/joernweissenborn/aursir4go/messages"
	"time"
	"encoding/json"
	"fmt"
	"aursirrt/src/config"
)

type DockerZmq struct {
	agent dock.DockAgent
	skt *zmq.Socket
	id string
	appLastPing map[string]time.Time
	interfaceip string
	incomingport int64
	broadcast bool
	udpport int64
}

func (dzmq DockerZmq) Launch(agent dock.DockAgent, id string) (err error) {
	    zsocket := fmt.Sprintf("tcp://%s:%d",dzmq.interfaceip,dzmq.incomingport)
	mprint(fmt.Sprint("Launching on ",zsocket))
	dzmq.agent = agent
	dzmq.id = id
	dzmq.appLastPing = map[string]time.Time{}
	dzmq.skt, err = zmq.NewSocket(zmq.ROUTER)

	if err != nil {
		mprint("Failed to start ZMQ Dock")
		return
	}

	dzmq.skt.Bind(zsocket)
	if err != nil {
		return
	}
	go dzmq.listen()
	dzmq.launchUdp()

	return
}

func (dzmq *DockerZmq) SetBroadcast(bcast bool){
	mprint(fmt.Sprint("P2P active: ",bcast))
	dzmq.broadcast = bcast
}
func (dzmq *DockerZmq) SetUDPPort(udpport int64){
	mprint(fmt.Sprint("udpport port is: ",udpport))
	dzmq.udpport = udpport
}
func (dzmq *DockerZmq) SetIncomingPort(incomingport int64){
	mprint(fmt.Sprint("Incoming port is: ",incomingport))
	dzmq.incomingport = incomingport
}
func (dzmq *DockerZmq) SetIp(ip string){
	mprint(fmt.Sprint("Broadcastinterfaceadress is: ",ip))
	dzmq.interfaceip = ip
}


func (dzmq *DockerZmq) listen() {

	mprint("Listening for incoming")

	for {

		msg, _ := dzmq.skt.RecvMessage(0)
		printDebug("got message from", msg)

		if len(msg) > 4 {

			senderId := msg[0]

			msgtype, err := strconv.ParseInt(msg[1], 10, 64)
			printDebug("ZMQAppDocker got message from", msg)
			codec := msg[2]
			if err == nil {

				switch msgtype {

				case messages.DOCK:
					encmsg := []byte(msg[3])
					IP := "localhost"
					if len(msg) > 5{
						IP = msg[5]
					}
					printDebug(IP)
					port, err := strconv.ParseInt(msg[4], 10, 64)
					if err == nil {
						conn := NewConnection(dzmq.incomingport, port, dzmq.interfaceip,IP, dzmq.id)
						dzmq.addAppPing(senderId)
						dzmq.agent.InitDocking(senderId, codec, encmsg, &conn)
					}
				default:
					encmsg := []byte(msg[3])
					dzmq.agent.ProcessMsg(senderId,msgtype,codec,encmsg)

				}
			}
		}
	}
}


func (dzmq *DockerZmq)launchUdp() (err error) {
	b := createBeacon(dzmq.broadcast,dzmq.id,dzmq.udpport,dzmq.incomingport,dzmq.interfaceip,dzmq)
	err = b.launch()
	if err != nil {
		return
	}
	go dzmq.checkAppLiveliness()
	return
}

func (dzmq *DockerZmq) addAppPing(AppId string) {
	dzmq.appLastPing[AppId] = time.Now()
}

func (dzmq *DockerZmq) checkInAppPing(AppId string, ip string, Port string) {
	if _,f:= dzmq.appLastPing[AppId];f {
		dzmq.appLastPing[AppId] = time.Now()
	} else if dzmq.broadcast && AppId!= dzmq.id && Port != ""{
		port,_ := strconv.ParseInt(Port,10,64)
		conn := NewConnection( dzmq.incomingport, port, dzmq.interfaceip, ip, dzmq.id)
		err := conn.Init()
		if err != nil {
			return
		}
		defer conn.Close()
		m,_ := json.Marshal(messages.DockMessage{"runtime@"+ip,[]string{"JSON"},true})
		conn.Send(messages.DOCK,"JSON",m)
	}
}

func (dzmq *DockerZmq) removeAppPing(AppId string) {
	delete(dzmq.appLastPing,AppId)
}

func (dzmq *DockerZmq) checkAppLiveliness() {

	t := time.NewTimer(3 * time.Second)
	for _ = range t.C {
		for id, lastCheckIn := range dzmq.appLastPing {
			if time.Since(lastCheckIn) > 3*time.Second {
				mprint("Apptimeout: "+id)
				dzmq.closeConnection(id)
			}
		}
		t.Reset(10*time.Second)
	}
}

func (dzmq *DockerZmq) closeConnection(id string){

	dzmq.removeAppPing(id)
	dzmq.agent.ProcessMsg(id,messages.LEAVE,"JSON",[]byte{})
}




func mprint(msg ...interface {}){
	log.Println("DOCKERZMQ", msg)
}

func printDebug(msg ...interface {}){
	if config.Debug{
		log.Println("DEBUG","DOCKERZMQ", msg)
	}
}

package dockzmq

import (
	"time"
	"log"
	"net"
	"strings"
	"fmt"
)

var pingtime = 1 * time.Second


type udpBeacon struct {
	broadcast bool
	id string
	udpport int64
	zmqport int64
	interfaceip string
	sock *net.UDPConn
	dzmq *DockerZmq
	kill chan struct{}
	outsock *net.UDPConn
}

func createBeacon(broadcast bool, id string,udpport, zmqport int64, interfaceip string,dzmq *DockerZmq) udpBeacon{
	return udpBeacon{broadcast,id,udpport,zmqport,interfaceip,nil,dzmq, nil,nil}
}

func (b udpBeacon) launch() (err error) {

	if b.broadcast {
		err = b.setupBroadcastListener()
		if err != nil {
			return
		}
		err = b.setupBeacon()
		err = b.setupBroadcastListener()
		if err != nil {
			return
		}
		go b.listen()
		go b.ping()
	} else {
		err = b.setupLocalListener()
		if err != nil {
			return
		}
		go b.listen()
	}

	return
}
func (b *udpBeacon) setupLocalListener() (err error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d",b.udpport))
	if err != nil {
		return
	}
	b.sock, err = net.ListenUDP("udp", addr)
	if err != nil {
		return
	}
	mprint(fmt.Sprint("Startet to listen for UDP ping on port", b.udpport))
	return
}

func (b *udpBeacon) setupBroadcastListener() (err error) {

	ifaceaddresses,_ := net.InterfaceAddrs()
	var Interface net.Interface
	for i,iface := range ifaceaddresses {
		addr := strings.Split(iface.String(),"/")[0]
		printDebug(fmt.Sprint("found networkinterface ",addr))
		if addr == b.interfaceip {
			ifaces,_ := net.Interfaces()
			Interface = ifaces[i]
			break
		}
	}

	mcip, err := net.ResolveUDPAddr("udp", fmt.Sprintf("224.0.0.251:%d",b.udpport))
	if err != nil {
		return
	}

	b.sock, err = net.ListenMulticastUDP("udp4", &Interface, mcip)
	if err != nil {
		return
	}
	mprint(fmt.Sprint("Setup listen for multicast UDP ping on port", b.udpport))
	return
}

func (b udpBeacon) listen(){
	var buf [1024]byte
	for {
		rlen, Ip, err := b.sock.ReadFromUDP(buf[:])
		if err != nil {
			log.Fatal("DOCKERZMQ",err)
		}
		beaconstring := strings.Split(string(buf[:rlen]),":")
		appid := beaconstring[0]
		port := ""
		if len(beaconstring) >1{
			port = beaconstring[1]
		}
		ip := strings.Split(Ip.String(),":")[0]

		go b.dzmq.checkInAppPing(appid, ip,port)
	}
}



func (b *udpBeacon) setupBeacon() (err error) {
	b.kill = make(chan struct{})

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:0",b.interfaceip))
	if err != nil {
		return
	}
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("224.0.0.251:%d",b.udpport))
	if err != nil {
		return
	}
	b.outsock, err = net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		return
	}
	return
}

func (b udpBeacon) ping() {
	t := time.NewTimer(pingtime)
	mprint(fmt.Sprintf("Beginning UDP Broadcast with %s:%d",b.id,b.zmqport))
	for {
		select {
		case <-b.kill:
			mprint("Stopping UDP")
			return

		case <-t.C:
			b.outsock.Write([]byte(fmt.Sprintf("%s:%d", b.id, b.zmqport)))
			t.Reset(pingtime)
		}
	}

}

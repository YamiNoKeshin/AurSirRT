package boot

import (
	"log"
	"aursirrt/src/processor"
	"aursirrt/src/storage"
	"aursirrt/src/dock"
	"aursirrt/src/dock/dockzmq"
	"aursirrt/src/dock/dockwebsockets"
	"aursirrt/src/cmdlineinterface"
	"flag"
	"aursirrt/src/config"
	"net"
	"strings"
	"strconv"
)

const (
	MAX_PROCESSORS = 8
)

func Boot(){


	bootCmdlineinterface(bootFunctionalCore())
}

func bootFunctionalCore() chan processor.Processor{
	mprint("AurSir RT starting")
	flag.Parse()

	//Boot the storage core

	a := bootStorage()

	//get id
	id:= a.GetId()
	mprint("Nodeid is "+id)

	//Boot the processing core
	p := bootCore(a)



	//if *config.Broadcast {

	//}
	bootZeromqDocker(p,id,"127.0.0.1",5555,5557,false)

	for _,conn :=range (config.Zconnections) {
		split := strings.Split(conn, ":")
		if len(split) == 1 {
			bootZeromqDocker(p,id,conn,getRandomPort(),5556,true)
		} else {
			port, _ := strconv.ParseInt(split[1],10,64)
			bootZeromqDocker(p,id,split[0],port,5557,false)
		}
	}

	var w dockwebsockets.DockerWebSockets
	bootDocker(p,w,id)
	return p

}
func BootWithoutCmdlineinterface(){
	bootFunctionalCore()
}

func bootStorage() storage.StorageAgent {
	mprint("Launching Storage")

	return storage.NewAgent()
}

func bootZeromqDocker(p chan processor.Processor, id, ip string,port int64,udpport int64, broadcast bool)  {
	mprint("Launching Docker")
	var nz dockzmq.DockerZmq
	nz.SetIp(ip)
	nz.SetBroadcast(broadcast)
	nz.SetIncomingPort(port)
	nz.SetUDPPort(udpport)
	bootDocker(p, nz, id)
}



func bootCore(a storage.StorageAgent) (processingChan chan processor.Processor){

	mprint("Launching Core")

	processingChan = make(chan processor.Processor)

	go processor.Process(processingChan, a, MAX_PROCESSORS)

	return
}

func bootDocker(p chan processor.Processor, d dock.Docker, id string) {
	mprint("Launching Dock")
	agent := dock.NewAgent(p)
	err := d.Launch(agent,id)
	if err != nil {
		log.Fatal("BOOT", err)
	}
}


func bootCmdlineinterface(p chan processor.Processor) {
	mprint("Launching Cmdlineinterface")
	cli := cmdlineinterface.CmdLineInterface{}
	cli.Run()

}


func mprint(msg string){

	log.Println("BOOT",msg)

}


func getRandomPort() int64 {
	l, err := net.Listen("tcp", "127.0.0.1:0") // listen on localhost
	if err != nil {
		panic("Could not find a free port")
	}
	defer l.Close()
	return int64(l.Addr().(*net.TCPAddr).Port)
}

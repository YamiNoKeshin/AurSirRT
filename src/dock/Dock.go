package dock


import (
	"aursirrt/src/processor/processors"
	"aursirrt/src/dock/connection"
	"aursirrt/src/processor"
	"log"
)

func NewAgent(c chan processor.Processor) DockAgent{
	 return DockAgent{c}
}

type DockAgent struct {
	procchan chan processor.Processor
}

func (da DockAgent) ProcessMsg(appid string,msgtype int64, codec string, msg []byte){
	var p processors.ParseMessageProccesor
	p.AppId = appid
	p.Codec = codec
	p.Msg = msg
	p.Type = msgtype
	p.GenericProcessor = processor.GetGenericProcessor()
	go da.launchProcess(p)
}

func (da DockAgent) InitDocking(appid string, codec string, msg []byte, connection connection.Connection){
	debugPrint("initializing docking")
	var p processors.DockProcessor
	p.AppId = appid
	p.Codec = codec
	p.DockMessage = msg
	p.Connection = connection
	p.GenericProcessor = processor.GetGenericProcessor()

	go da.launchProcess(p)

}

func (da DockAgent) launchProcess(p processor.Processor){
	da.procchan <- p
}


func debugPrint(msg string){

	if true {
		log.Println("DEBUG DOCK",msg)

	}

}

package processor

import (
	"aursirrt/src/storage"
)

type Processor interface {
	Process()
	Init(chan Processor, storage.StorageAgent)
	SpawnProcess(Processor)
}

func GetGenericProcessor() (gp *GenericProcessor) {
	var proc GenericProcessor
	gp = &proc
	return
}

type GenericProcessor struct {
	procchan chan Processor
	agent storage.StorageAgent
}

func (gp GenericProcessor) Process() {}
func (gp GenericProcessor) GetAgent() storage.StorageAgent{
	return gp.agent
}

func (gp GenericProcessor) SpawnProcess(p Processor) {
	go func(){
		gp.procchan <- p
	}()
}

func (p *GenericProcessor) Init(c chan Processor, a storage.StorageAgent) {
	p.procchan = c
	p.agent = a
}



func Testprocessor() chan Processor {
	pc := make(chan Processor)
	go Process(pc,storage.NewAgent(), 1)
	return pc
}

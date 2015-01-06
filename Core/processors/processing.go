package processors

import (
	"github.com/joernweissenborn/aursirrt/core/storage/agent"
	"github.com/joernweissenborn/aursirrt/config"
)

func StartProcessing(cfg config.RtConfig) (procChan chan Processor) {

	procChan = make(chan Processor)
	go processingListener(procChan, agent.NewAgent(cfg))
	return 
}

func processingListener(procChan chan Processor, storeAgent agent.StorageAgent) {

	for proc := range procChan {

		proc.StorageAgent = storeAgent

		go proc.Process()

	}

}

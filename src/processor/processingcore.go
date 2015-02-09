package processor

import (
	"log"
	"aursirrt/src/storage"
	"fmt"
	"aursirrt/src/config"
)



func Process(procchan chan Processor, storageagent storage.StorageAgent, maxprocesses int64) {

	print("Initialized")

	procslots := make(chan struct {},maxprocesses)

	for proc := range procchan {
		procslots <- struct{}{}
		if proc != nil {
			debugPrint("Processing ")
			debugPrint(fmt.Sprint(proc))

			proc.Init(procchan, storageagent)


			go process(proc,procslots)
		}
	}

}

func process(p Processor, ps chan struct{}){
	p.Process()
	debugPrint("finished ")
	<- ps

}

func print(msg string){
	log.Println("PROCESSOR",msg)

}

func debugPrint(msg string){

	if config.Debug {
		log.Println("DEBUG PROCESSOR",msg)

	}

}

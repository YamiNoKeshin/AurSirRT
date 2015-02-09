package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type StartListenProcessor struct {

	*processor.GenericProcessor

	AppId string

	StartListenMsg messages.ListenMessage

}

func (p StartListenProcessor) Process() {

	Import := types.GetImportById(p.StartListenMsg.ImportId,p.GetAgent())

	Import.StartListenToFunction(p.StartListenMsg.FunctionName)
	
	if !Import.GetApp().IsNode(){
		for _,n := range types.GetNodes(p.GetAgent()){
				var smp SendMessageProcessor
				smp.App = n
				smp.Msg = p.StartListenMsg
				smp.GenericProcessor = processor.GetGenericProcessor()
				p.SpawnProcess(smp)

			
		}		
		
	}



}


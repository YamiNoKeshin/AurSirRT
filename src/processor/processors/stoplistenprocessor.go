package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type StopListenProcessor struct {

	*processor.GenericProcessor

	AppId string

	StopListenMsg messages.StopListenMessage

}

func (p StopListenProcessor) Process() {

	Import := types.GetImportById(p.StopListenMsg.ImportId,p.GetAgent())

	Import.StopListenToFunction(p.StopListenMsg.FunctionName)
	
	if !Import.GetApp().IsNode(){
		for _,n := range types.GetNodes(p.GetAgent()){
				var smp SendMessageProcessor
				smp.App = n
				smp.Msg = p.StopListenMsg
				smp.GenericProcessor = processor.GetGenericProcessor()
				p.SpawnProcess(smp)
			
			
		}		
		
	}



}


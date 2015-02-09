package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type RemoveExportProcessor struct {

	*processor.GenericProcessor

	AppId string

	RemoveExportMsg messages.RemoveExportMessage

}

func (p RemoveExportProcessor) Process() {
	printDebug("REMOVEEXPORT",p.RemoveExportMsg)
	Export := types.GetExportById(p.RemoveExportMsg.ExportId,p.GetAgent())
	isapp := !Export.GetApp().IsNode()
	Export.Remove()
	var uesp ExportedStateProcessor
	uesp.AppKey = Export.GetAppKey()
	uesp .GenericProcessor = processor.GetGenericProcessor()

	p.SpawnProcess(uesp)


	if isapp {
		for _, node := range types.GetNodes(p.GetAgent()){
			node.Lock()
			var smp SendMessageProcessor
			smp.App = node
			smp.Msg = p.RemoveExportMsg
			smp.GenericProcessor = processor.GetGenericProcessor()
			p.SpawnProcess(smp)
			node.Unlock()
		}
	}
}


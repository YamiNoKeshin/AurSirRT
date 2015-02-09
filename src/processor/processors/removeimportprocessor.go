package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type RemoveImportProcessor struct {

	*processor.GenericProcessor

	AppId string

	RemoveImportMsg messages.RemoveImportMessage

}

func (p RemoveImportProcessor) Process() {
	printDebug("REMOVEIMPORT",p.RemoveImportMsg)
	Import := types.GetImportById(p.RemoveImportMsg.ImportId,p.GetAgent())
	isapp := !Import.GetApp().IsNode()
	Import.Remove()
	var uesp ExportedStateProcessor
	uesp.AppKey = Import.GetAppKey()
	uesp .GenericProcessor = processor.GetGenericProcessor()

	p.SpawnProcess(uesp)


	if isapp {
		for _, node := range types.GetNodes(p.GetAgent()){
			node.Lock()
			var smp SendMessageProcessor
			smp.App = node
			smp.Msg = p.RemoveImportMsg
			smp.GenericProcessor = processor.GetGenericProcessor()
			p.SpawnProcess(smp)
			node.Unlock()
		}
	}
}


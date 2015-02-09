package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type AddImportProcessor struct {

	*processor.GenericProcessor

	AppId string

	AddImportMsg messages.AddImportMessage

}

func (p AddImportProcessor) Process() {

	Import := types.GetImport(p.AppId,p.AddImportMsg.AppKey, p.AddImportMsg.Tags ,p.AddImportMsg.ImportId,p.GetAgent())
	Import.Add()
	app := Import.GetApp()
	if !app.IsNode() {
		var smp SendMessageProcessor
		smp.App = app
		smp.Msg = messages.ImportAddedMessage{Import.GetId(),Import.HasExporter(), p.AddImportMsg.AppKey,p.AddImportMsg.Tags}
		smp.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(smp)
		p.AddImportMsg.ImportId = Import.GetId()
		for _, node := range types.GetNodes(p.GetAgent()){
			node.Lock()
			var smp SendMessageProcessor
			smp.App = node
			smp.Msg = p.AddImportMsg
			smp.GenericProcessor = processor.GetGenericProcessor()
			p.SpawnProcess(smp)
			node.Unlock()
		}
	}
}


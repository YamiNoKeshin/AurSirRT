package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type AddExportProcessor struct {

	*processor.GenericProcessor

	AppId string

	AddExportMsg messages.AddExportMessage

}

func (p AddExportProcessor) Process() {
	printDebug("AddExport")
	Export := types.GetExport(p.AppId,p.AddExportMsg.AppKey,p.AddExportMsg.Tags,p.AddExportMsg.ExportId,p.GetAgent())
	Export.Add()
	app := Export.GetApp()

	if !app.IsNode() {
		var smp SendMessageProcessor
		smp.App = app
		smp.Msg = messages.ExportAddedMessage{Export.GetId(),p.AddExportMsg.AppKey,p.AddExportMsg.Tags}
		smp.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(smp)
	}

	var pjp PendingJobProcessor
	pjp.Appkey = Export.GetAppKey()
	pjp.GenericProcessor = processor.GetGenericProcessor()
	p.SpawnProcess(pjp)
	var uesp ExportedStateProcessor
	uesp.AppKey = Export.GetAppKey()
	uesp.GenericProcessor = processor.GetGenericProcessor()
	p.SpawnProcess(uesp)

	if !app.IsNode() {
		p.AddExportMsg.ExportId = Export.GetId()
		for _, node := range types.GetNodes(p.GetAgent()){
			var smp SendMessageProcessor
			smp.App = node
			smp.Msg = p.AddExportMsg
			smp.GenericProcessor = processor.GetGenericProcessor()
			p.SpawnProcess(smp)
		}
	}
}


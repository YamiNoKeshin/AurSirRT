package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type UpdateExportProcessor struct {

	*processor.GenericProcessor

	AppId string

	UpdateExportMsg messages.UpdateExportMessage

}

func (p UpdateExportProcessor) Process() {
	printDebug("UPDATEEXPORT",p.UpdateExportMsg)
	Export := types.GetExportById(p.UpdateExportMsg.ExportId,p.GetAgent())
	Export.UpdateTags(p.UpdateExportMsg.Tags)
	var pjp PendingJobProcessor
	pjp.Appkey = Export.GetAppKey()
	pjp .GenericProcessor = processor.GetGenericProcessor()
	p.SpawnProcess(pjp)
	var uesp ExportedStateProcessor
	uesp.AppKey = Export.GetAppKey()
	uesp .GenericProcessor = processor.GetGenericProcessor()

	p.SpawnProcess(uesp)


	if !Export.GetApp().IsNode() {
		for _, node := range types.GetNodes(p.GetAgent()){
			node.Lock()
			var smp SendMessageProcessor
			smp.App = node
			smp.Msg = p.UpdateExportMsg
			smp.GenericProcessor = processor.GetGenericProcessor()
			p.SpawnProcess(smp)
			node.Unlock()
		}
	}
}




package processors

import (
	"processor"
	"storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type UpdateImportProcessor struct {

	*processor.GenericProcessor

	AppId string

	UpdateImportMsg messages.UpdateImportMessage

}

func (p UpdateImportProcessor) Process() {

	Import := types.GetImportById(p.UpdateImportMsg.ImportId,p.GetAgent())
	Import.UpdateTags(p.UpdateImportMsg.Tags)

	p.SpawnProcess(smp)
	if !Import.GetApp().IsNode() {
		var smp SendMessageProcessor
		smp.App = Import.GetApp()
		smp.Msg = messages.ImportUpdatedMessage{Import.GetId(),Import.HasExporter()}
		smp.GenericProcessor = processor.GetGenericProcessor()
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


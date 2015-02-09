package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"

	"github.com/joernweissenborn/aursir4go/messages"
)

type ExportedStateProcessor struct {

	*processor.GenericProcessor

	AppKey types.AppKey


}

func (p ExportedStateProcessor) Process() {
	for _, imp := range p.AppKey.GetImporter() {
		app := imp.GetApp()
		ok := app.Lock()
		if ok {
			if !app.IsNode() {
				var smp SendMessageProcessor
				smp.App = app
				printDebug(imp.HasExporter())
				smp.Msg = messages.ImportUpdatedMessage{imp.GetId(),imp.HasExporter()}
				smp.GenericProcessor = processor.GetGenericProcessor()
				p.SpawnProcess(smp)
			}
		app.Unlock()
		}
	}
}


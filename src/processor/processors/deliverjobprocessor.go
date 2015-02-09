package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"
	"github.com/joernweissenborn/aursir4go/calltypes"
)


type DeliverJobProcessor struct {

	*processor.GenericProcessor

	Job types.Job

}

func (p DeliverJobProcessor) Process() {
	printDebug("Deliverjob")

	if p.Job.Exists() {

		imp := p.Job.GetImport()
		if imp.HasExporter() {
			req := p.Job.GetRequest()
			switch req.CallType{
			case calltypes.MANY2MANY, calltypes.MANY2ONE:
			for _, exp := range imp.GetExporter() {
				p.Job.Assign(exp)
				var smp SendMessageProcessor
				smp.App = exp.GetApp()
				req.ExportId = exp.GetId()
				smp.Msg = req
				smp.GenericProcessor = processor.GetGenericProcessor()

				p.SpawnProcess(smp)
			}

			case calltypes.ONE2MANY, calltypes.ONE2ONE:
				exp := imp.GetExporter()[0]
				p.Job.Assign( exp)
				var smp SendMessageProcessor
				smp.App = exp.GetApp()
				req.ExportId = exp.GetId()
				smp.Msg = req
				smp.GenericProcessor = processor.GetGenericProcessor()
				p.SpawnProcess(smp)
			}
		}
	}
}

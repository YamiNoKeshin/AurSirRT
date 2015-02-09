package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"
)


type PendingJobProcessor struct {

	*processor.GenericProcessor

	Appkey types.AppKey

}

func (p PendingJobProcessor) Process() {
	printDebug("pendingjob")

	for _, imp := range p.Appkey.GetImporter(){
		for _, j := range imp.GetJobs() {
			if !j.IsAssigned() {
				if imp.HasExporter() {
					exp := imp.GetExporter()[0]
					j.Assign(exp)
					var djp DeliverJobProcessor
					djp.Job = j
					djp.GenericProcessor = processor.GetGenericProcessor()

					p.SpawnProcess(djp)
				}
			}
		}
	}

}

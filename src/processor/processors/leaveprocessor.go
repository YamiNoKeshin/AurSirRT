package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type LeaveProcessor struct {

	*processor.GenericProcessor

	AppId string

}

func (p LeaveProcessor) Process() {

	app := types.GetApp(p.AppId, p.GetAgent())
	if app.Exists(){
		keys := []types.AppKey{}
		for _, e:= range app.GetExports() {
			keys = append(keys,e.GetAppKey())
		}
		conn,_ := app.GetConnection()
		conn.Close()
		printDebug("LEAVE locking")
		app.Lock()
		printDebug("LEAVE unlocking")
		app.Unlock()
		if !app.IsNode() {
			for _, exp := range app.GetExports() {
				printDebug("LEAVE unlocking")
				var m messages.RemoveExportMessage
				m.ExportId = exp.GetId()
				var smp SendMessageProcessor
				smp.Msg = m
				smp.GenericProcessor = processor.GetGenericProcessor()
				for _, n := range types.GetNodes(p.GetAgent()){
					smp.App = n
					p.SpawnProcess(smp)
				}
			}
			for _, imp := range app.GetImports() {
				var m messages.RemoveImportMessage
				m.ImportId = imp.GetId()
				var smp SendMessageProcessor
				smp.Msg = m
				smp.GenericProcessor = processor.GetGenericProcessor()
				for _, n := range types.GetNodes(p.GetAgent()){
					smp.App = n
					p.SpawnProcess(smp)
				}
			}

			printDebug("LEAVE removing")
			app.Remove()
			printDebug("LEAVE removed")

			for _, k:= range keys {
				if k.Exists() {
					var esp ExportedStateProcessor
					esp.AppKey = k
					esp.GenericProcessor = processor.GetGenericProcessor()
					p.SpawnProcess(esp)
				}
			}
		}
	}

}


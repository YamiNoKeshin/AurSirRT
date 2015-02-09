package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"

	"github.com/joernweissenborn/aursir4go/messages"
	"github.com/joernweissenborn/aursir4go/util"
	"aursirrt/src/dock/connection"
)

type DockProcessor struct {

	*processor.GenericProcessor

	AppId string

	Codec string

	DockMessage []byte

	Connection connection.Connection

}

func (p DockProcessor) Process() {
	if p.Connection != nil {
		decoder := util.GetCodec(p.Codec)
		if decoder == nil {
			return
		}
		var dmsg messages.DockMessage
		err := decoder.Decode(p.DockMessage, &dmsg)
		if err != nil {
			return
		}
		app := types.GetApp(p.AppId, p.GetAgent())
		ok := app.Create(dmsg, p.Connection)
		if ok {
			app.Lock()
			defer app.Unlock()
			conn,_ := app.GetConnection()
			err = conn.Init()
			if err != nil {
				return
			}
			if app.IsNode() {
				var sp SendMessageProcessor
				sp.App = app
				sp.Msg = messages.DockMessage{"runtime",[]string{"JSON"},true}
				sp.GenericProcessor = processor.GetGenericProcessor()
				p.SpawnProcess(sp)
			}

		}

		var sp SendMessageProcessor
		sp.App = app
		sp.Msg = messages.DockedMessage{ok}
		sp.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(sp)

		if ok && app.IsNode() {
			for _, localapp := range types.GetApps(p.GetAgent()){
			 	for _, imp := range localapp.GetImports() {
					var m messages.AddImportMessage
					m.AppKey = imp.GetAppKey().GetKey()
					m.Tags = imp.GetTagNames()
					m.ImportId = imp.GetId()
					var sp SendMessageProcessor
					sp.App = app
					sp.Msg = m
					sp.GenericProcessor = processor.GetGenericProcessor()
					p.SpawnProcess(sp)
				}
				for _, exp := range localapp.GetExports() {
					var m messages.AddExportMessage
					m.AppKey = exp.GetAppKey().GetKey()
					m.Tags = exp.GetTagNames()
					m.ExportId = exp.GetId()
					var sp SendMessageProcessor
					sp.App = app
					sp.Msg = m
					sp.GenericProcessor = processor.GetGenericProcessor()
					p.SpawnProcess(sp)
				}
		}
		}
	}

}

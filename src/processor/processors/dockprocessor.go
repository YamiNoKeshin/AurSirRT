package processors

import (
	"processor"
	"storage/types"
	"dock"
	"github.com/joernweissenborn/aursir4go/messages"
	"github.com/joernweissenborn/aursir4go/util"
)

type DockProcessor struct {

	*processor.GenericProcessor

	AppId string

	Codec string

	DockMessage []byte

	Connection dock.Connection

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
		conn := app.GetConnection()
		err = conn.Init()
		if err != nil {
			return
		}
		var sp SendMessageProcessor
		sp.App = app
		sp.Msg = messages.DockedMessage{ok}
		p.SpawnProcess(sp)

	}

}


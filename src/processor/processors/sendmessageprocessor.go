package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
	"github.com/joernweissenborn/aursir4go/util"
	"aursirrt/src/dock/connection"
	"fmt"
)

type SendMessageProcessor struct {

	*processor.GenericProcessor

	App types.App

	Msg messages.AurSirMessage
}

func (p SendMessageProcessor) Process() {
	printDebug(fmt.Sprint("sendmsg",p.Msg) )
	conn, ok := p.App.GetConnection()
	id := p.App.Id
	if ok {

		var msgtype int64
		var codec string = "JSON"

		switch p.Msg.(type) {

		case messages.DockedMessage:
			msgtype = messages.DOCKED

		case messages.AddExportMessage:
			msgtype = messages.ADD_EXPORT
		case messages.ExportAddedMessage:
			msgtype = messages.EXPORT_ADDED
		case messages.RemoveExportMessage:
			msgtype = messages.REMOVE_EXPORT


		case messages.AddImportMessage:
			msgtype = messages.ADD_IMPORT

		case messages.ImportAddedMessage:
			msgtype = messages.IMPORT_ADDED

		case messages.ImportUpdatedMessage:
			msgtype = messages.IMPORT_UPDATED

		case messages.RemoveImportMessage:
			msgtype = messages.REMOVE_IMPORT

		case *messages.Request:
			msgtype = messages.REQUEST

		case messages.Result:
			msgtype = messages.RESULT

		case messages.ListenMessage:
			msgtype = messages.LISTEN

		case messages.StopListenMessage:
			msgtype = messages.STOP_LISTEN

		}

		decoder := util.GetCodec(codec)
		if decoder == nil {
			return
		}
		msg, err := decoder.Encode(p.Msg)
		if err != nil {
			return
		}
		go p.send(id,conn,msgtype,codec,msg)
	}

}

func (p SendMessageProcessor) send(id string,conn connection.Connection,msgtype int64, codec string,msg []byte) {
	err := conn.Send(msgtype,codec,msg)
	if nil != err {
		var lp LeaveProcessor
		lp.AppId = id
		lp.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(lp)
	} else {
		printDebug(fmt.Sprint("msg send"))

	}

}


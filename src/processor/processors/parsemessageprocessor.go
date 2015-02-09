package processors

import (
	"aursirrt/src/processor"
	"github.com/joernweissenborn/aursir4go/util"
	"github.com/joernweissenborn/aursir4go/messages"
)


type ParseMessageProccesor struct {

	*processor.GenericProcessor

	AppId string
	Codec string
	Type int64
	Msg []byte
}

func (p ParseMessageProccesor) Process() {

	decoder := util.GetCodec(p.Codec)
	if decoder == nil {
		return
	}

	switch p.Type {
	case messages.LEAVE:
		var np LeaveProcessor
			np.AppId = p.AppId
		np.GenericProcessor = processor.GetGenericProcessor()

		p.SpawnProcess(np)

	case messages.ADD_EXPORT:
		var m messages.AddExportMessage
		decoder.Decode(p.Msg, &m)
		var np AddExportProcessor
		np.AppId = p.AppId
		np.AddExportMsg = m
		np.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(np)

	case messages.UPDATE_EXPORT:
		var m messages.UpdateExportMessage
		decoder.Decode(p.Msg, &m)
		var np UpdateExportProcessor
		np.AppId = p.AppId
		np.UpdateExportMsg = m
		np.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(np)

	case messages.REMOVE_EXPORT:
		var m messages.RemoveExportMessage
		decoder.Decode(p.Msg, &m)
		var np RemoveExportProcessor
		np.AppId = p.AppId
		np.RemoveExportMsg = m
		np.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(np)

	case messages.REMOVE_IMPORT:
		var m messages.RemoveImportMessage
		decoder.Decode(p.Msg, &m)
		var np RemoveImportProcessor
		np.AppId = p.AppId
		np.RemoveImportMsg = m
		np.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(np)

	case messages.ADD_IMPORT:
		var m messages.AddImportMessage
		decoder.Decode(p.Msg, &m)
		var np AddImportProcessor
		np.AppId = p.AppId
		np.AddImportMsg = m
		np.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(np)

	case messages.UPDATE_IMPORT:
		var m messages.UpdateImportMessage
		decoder.Decode(p.Msg, &m)
		var np UpdateImportProcessor
		np.AppId = p.AppId
		np.UpdateImportMsg = m
		np.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(np)
case messages.LISTEN:
		var m messages.ListenMessage
		decoder.Decode(p.Msg, &m)
		var np StartListenProcessor
		np.AppId = p.AppId
		np.StartListenMsg = m
		np.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(np)
case messages.STOP_LISTEN:
		var m messages.StopListenMessage
		decoder.Decode(p.Msg, &m)
		var np StopListenProcessor
		np.AppId = p.AppId
		np.StopListenMsg = m
		np.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(np)

	case messages.REQUEST:
		var m messages.Request
		decoder.Decode(p.Msg, &m)
		var np RequestProcessor
		np.AppId = p.AppId
		np.Request = m
		np.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(np)
	case messages.RESULT:
		var m messages.Result
		decoder.Decode(p.Msg, &m)
		var np ResultProcessor
		np.AppId = p.AppId
		np.Result = m
		np.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(np)
	}

}


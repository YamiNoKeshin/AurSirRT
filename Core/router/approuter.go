package router
/*

import (
	"log"
	"github.com/joernweissenborn/aursirrt/core"
	"github.com/joernweissenborn/aursir4go"
	"github.com/joernweissenborn/aursirrt/core/processors"
	"github.com/joernweissenborn/aursirrt/core/processors/dockprocessor"
)

func RouteIncomingAppMsg(appInChannel chan core.AppMessage, procChan chan processors.Processor) {

	for AppMessage := range appInChannel {

		aursirMessage, err := AppMessage.AppMsg.Decode()
		log.Println("DEBUG",aursirMessage)
		log.Println("DEBUG",string(AppMessage.AppMsg.Msg),err)


		if err == nil {

			switch aursirMessage := aursirMessage.(type) {

			case aursir4go.AurSirDockMessage:
				procChan <- dockprocessor.DockProcessor{AppMessage.SenderUUID,aursirMessage}

			case aursir4go.AurSirLeaveMessage:
				go c.leave(AppMessage.SenderUUID)

			case aursir4go.AurSirAddExportMessage:
				go c.addExport(AppMessage.SenderUUID, aursirMessage)

			case  aursir4go.AurSirUpdateExportMessage:
				go c.updateExport(AppMessage.SenderUUID,aursirMessage)

			case aursir4go.AurSirAddImportMessage:
				go c.addImport(AppMessage.SenderUUID, aursirMessage)
			case  aursir4go.AurSirUpdateImportMessage:
				go c.updateImport(AppMessage.SenderUUID,aursirMessage)

			case aursir4go.AurSirRequest:
				go c.request(AppMessage.SenderUUID,aursirMessage)

			case aursir4go.AurSirResult:
				go c.result(AppMessage.SenderUUID,aursirMessage)
			case aursir4go.AurSirListenMessage:
				go c.listen(AppMessage.SenderUUID,aursirMessage)
			case aursir4go.AurSirStopListenMessage:
				go c.stopListen(AppMessage.SenderUUID,aursirMessage)

			case aursir4go.AurSirCallChain:
				go c.callChain(AppMessage.SenderUUID,aursirMessage)

			default:
				log.Println("unknown message")
			}
		}
	}
}
*/

package processors

import (
	"aursirrt/src/processor"

	"github.com/joernweissenborn/aursir4go/messages"
)

type DockedProcessor struct {

	*processor.GenericProcessor

	AppId string

	DockedMessage messages.DockedMessage



}

func (p DockedProcessor) Process() {
	if p.DockedMessage.Ok {


	}

}

package core

import (
	//"github.com/joernweissenborn/aursirrt/core/processors"
	//"github.com/joernweissenborn/aursirrt/core/router"
	"github.com/joernweissenborn/aursirrt/config"
)

func Boot(cfg config.RtConfig) (AppInChan chan AppMessage){

//	procChan := processors.StartProcessing(cfg)

	AppInChan = make(chan AppMessage)

	//go router.RouteIncomingAppMsg(AppInChan,procChan)

	return
}

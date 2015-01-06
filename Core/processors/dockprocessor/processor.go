package dockprocessor

import "github.com/joernweissenborn/aursir4go"
import "github.com/joernweissenborn/aursirrt/core/processors"

type DockProcessor struct {

	processors.Processor

	AppId string

	DockMessage aursir4go.AurSirDockMessage

}

func (p DockProcessor) Process() {

	app := p.GetApp()

	if !app.Exists("") {
		app.Create(p.AppId,p.DockMessage)
	}

}

func registerApp(){}

package processors

import (
	"github.com/joernweissenborn/aursirrt/core/storage/agent"
	"github.com/joernweissenborn/aursirrt/core/storage/types"
)


type Processor struct {
	StorageAgent agent.StorageAgent
}


func (p Processor) Process() {}
func (p Processor) GetApp() (app types.App) {
	app.Agent = p.StorageAgent
	return
}

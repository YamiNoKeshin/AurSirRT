package types

import (
	"log"
	"aursirrt/src/storage"
	"github.com/joernweissenborn/aursir4go/messages"
	"aursirrt/src/dock/connection"
	"fmt"
	"sync"
)

type appproperties struct {
	dockmsg    messages.DockMessage
	connection connection.Connection
	*sync.Mutex
}

type App struct {
	agent storage.StorageAgent
	Id string
}

func GetApp(Id string, Agent storage.StorageAgent) App {
	return App{Agent,Id}
}

func GetNodes(agent storage.StorageAgent) []App {

	c := make(chan string)
	agent.Read(func (sc *storage.StorageCore){
		for _, nodeedge := range sc.Root.Outgoing() {
			if nodeedge.Label() == KNOWN_NODE_EDGE {
				c <- nodeedge.Head().Id()
			}
		}
		close(c)
	})
	nodes := []App{}
	nodeids := []string{}
	for nodeid := range c {
		nodeids = append(nodeids,nodeid)
	}
	for _,nodeid := range nodeids {
		nodes = append(nodes,GetApp(nodeid,agent))
	}
	return nodes
}

func GetApps(agent storage.StorageAgent) []App {

	c := make(chan string)
	agent.Read(func (sc *storage.StorageCore){
		for _, nodeedge := range sc.Root.Outgoing() {
			if nodeedge.Label() == KNOWN_APP_EDGE {
				c <- nodeedge.Head().Id()
			}
		}
		close(c)
	})
	apps := []App{}
	appids := []string{}
	for appid := range c {
		appids = append(appids,appid)
	}
	for _, appid := range appids {
		apps = append(apps,GetApp(appid,agent))
	}
	return apps

}



func (app App) Exists() bool {

	c := make(chan bool)
	defer close(c)
	app.agent.Read(func (sc *storage.StorageCore){
		if av := sc.GetVertex(app.Id); av!=nil {
			c <- av.Id() == app.Id
		} else {
			c<-false
		}
	})

	return <- c
}

func (app App) IsNode() bool {
	if !app.Exists() {
		return false
	}
	c := make(chan bool)
	defer close(c)
	app.agent.Read(func (sc *storage.StorageCore){
		c <- sc.GetVertex(app.Id).Properties().(appproperties).dockmsg.Node
	})

	return <- c
}


func (app App) Lock() bool {
	props, ok := app.getProperties()
	if ok {
		props.Lock()
	}
	return ok
}

func (app App) Unlock() {
	props, ok := app.getProperties()
	if ok {
		props.Unlock()
	}
}

func (app App) GetConnection() (connection.Connection, bool) {

	props, ok := app.getProperties()
	return props.connection,ok
}

func (app App) getProperties() (appproperties,bool) {

	c := make(chan appproperties)
	defer close(c)
	fail := make(chan struct{})
	defer close(fail)
	app.agent.Read(func (sc *storage.StorageCore){
		av := sc.GetVertex(app.Id)
		if av != nil {
			c <- av.Properties().(appproperties)
		} else {
			fail <- struct{}{}
		}
	})

	select {
	case aprops,_ := <-c:
		return aprops, true

	case <-fail:
		return appproperties{},false

	}

}

func (app App) Create(DockMessage messages.DockMessage, Connection connection.Connection) bool{
	if app.Exists() {
		printDebug(fmt.Sprint("Cannot add app",app.Id))
		return false
	}
	c := make(chan struct{})
	defer close(c)
	mutex := &sync.Mutex{}
	edge := KNOWN_APP_EDGE
	if DockMessage.Node {
		edge = KNOWN_NODE_EDGE
	}
	app.agent.Write(func (sc *storage.StorageCore){
		log.Println("STORAGECORE","Assinging ID to App",DockMessage.AppName,"->", app.Id)
		av := sc.CreateVertex(app.Id, appproperties{DockMessage,Connection,mutex})
			sc.CreateEdge(storage.GenerateUuid(), edge,av, sc.Root,nil)

		c <- struct{}{}
	})

	<-c
	return true
}

func (app App) GetExports() (exports []Export){
	exportids := []string{}
	exports = []Export{}
	c := make(chan string)

	app.agent.Read(func (sc *storage.StorageCore){
		av := sc.GetVertex(app.Id)
		for _, expedge := range av.Outgoing() {
			if expedge.Label() == EXPORT_EDGE {
				c<- expedge.Head().Id()
			}
		}
		close(c)
	})
	for id := range c {
		exportids = append(exportids, id)
	}

	for _, id := range exportids {
		exports = append(exports, GetExportById(id,app.agent))
	}


	return
}

func (app App) GetImports() (imports []Import){
	importids := []string{}
	imports = []Import{}
	c := make(chan string)

	app.agent.Read(func (sc *storage.StorageCore){
		av := sc.GetVertex(app.Id)
		for _, expedge := range av.Outgoing() {
			if expedge.Label() == IMPORT_EDGE {
				c<- expedge.Head().Id()
			}
		}
		close(c)
	})
	for id := range c {
		importids = append(importids, id)
	}
	for _,id := range importids {
		imports = append(imports, GetImportById(id,app.agent))
	}
	return
}


func (app App) Remove(){
	c := make(chan struct{})
	if !app.Exists() {
		log.Println("STORAGECORE", "Error removing app %d, does not exist", app.Id)
		return
	}
	for _, i := range app.GetImports() {
		i.Remove()
	}
	for _, e := range app.GetExports() {
		e.Remove()
	}
	app.agent.Write(func (sc *storage.StorageCore){
		log.Println("STORAGECORE","Removing App",app.Id)
		sc.RemoveVertex(app.Id)
		c <- struct{}{}
	})

	<-c
}



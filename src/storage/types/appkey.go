package types

import (
	"aursirrt/src/storage"
	"github.com/joernweissenborn/aursir4go/appkey"
	"log"
)

type AppKey struct {
	agent storage.StorageAgent
	appkey appkey.AppKey
}

func GetAppKey(appkey appkey.AppKey, agent storage.StorageAgent) AppKey{
	return AppKey{agent,appkey}
}

func (a AppKey) GetKey() appkey.AppKey {
	return a.appkey
}
func (a AppKey) Exists() bool {

	c := make(chan bool)

	a.agent.Read(func (sc *storage.StorageCore){
		for _, ke := range sc.Root.Outgoing() {
			if ke.Label() == KNOWN_APPKEY_EDGE {
				key, _ := ke.Head().Properties().(appkey.AppKey)
				if key.ApplicationKeyName == a.appkey.ApplicationKeyName {
					c<-true
					return
				}
			}
		}
		c<-false
	})

	return <- c
}

func (a AppKey) GetId() string {

	c := make(chan string)

	a.agent.Read(func (sc *storage.StorageCore){
		for _, ke := range sc.Root.Outgoing() {
			if ke.Label() == KNOWN_APPKEY_EDGE {
				key, _ := ke.Head().Properties().(appkey.AppKey)
				if key.ApplicationKeyName == a.appkey.ApplicationKeyName {
					c<-ke.Head().Id()
					return
				}
			}
		}
		c<-""
	})

	return <- c
}
func (a AppKey) GetExporter() (exporter []Export) {
	exporter = []Export{}
	exporterids := []string{}
	c := make(chan string,1)
	kid := a.GetId()
	a.agent.Read(func (sc *storage.StorageCore){
		for _, ke := range sc.GetVertex(kid).Incoming() {
			if ke.Label() == EXPORT_EDGE {
					c<-ke.Tail().Id()
				}
			}

		close(c)
	})
	for eid := range c {
		exporterids = append(exporterids,eid)
	}
	for _, eid := range exporterids {
		exporter = append(exporter,GetExportById(eid,a.agent))
	}
	return
}
func (a AppKey) GetImporter() (importer []Import) {
	importer = []Import{}
	importerids := []string{}
	c := make(chan string,1)
	kid := a.GetId()
	a.agent.Read(func (sc *storage.StorageCore){
		for _, ke := range sc.GetVertex(kid).Incoming() {
			if ke.Label() == IMPORT_EDGE {
					c<-ke.Tail().Id()
				}
			}

		close(c)
	})
	for iid := range c {
		importerids = append(importerids,iid)
	}
	for _,iid := range importerids {
		importer = append(importer,GetImportById(iid,a.agent))
	}
	return
}

func (a AppKey) GetListener(function string,export Export) (importer []Import) {
	importer = []Import{}
	importerids := []string{}
	c := make(chan string,1)
	kid := a.GetId()
	a.agent.Read(func (sc *storage.StorageCore){
		for _, ke := range sc.GetVertex(kid).Incoming() {
			if ke.Label() == LISTEN_EDGE  {
				if ke.Properties().(string) == function {
					c<-ke.Tail().Id()
				}
				}
			}

		close(c)
	})
	for iid := range c {
		importerids = append(importerids,iid)
	}
	for _, iid := range importerids {
		imp := GetImportById(iid,a.agent)
		if export.HasTags(imp.GetTagNames()) {
			importer = append(importer,imp)
		}
	}
	return
}

func (a AppKey) Create() {
	c := make(chan bool)
	if !a.Exists() {
		a.agent.Write(func(sc *storage.StorageCore) {
			kv := sc.InMemoryGraph.CreateVertex(storage.GenerateUuid(), a.appkey)

			sc.InMemoryGraph.CreateEdge(storage.GenerateUuid(), KNOWN_APPKEY_EDGE, kv, sc.Root, nil)

			log.Println("StorageCore Key registered")
			c<-false
		})
		<-c
	}
	return
}

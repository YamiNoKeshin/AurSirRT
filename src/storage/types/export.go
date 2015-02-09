package types

import (
	"aursirrt/src/storage"
	"github.com/joernweissenborn/aursir4go/appkey"
	"log"
)




type Export struct {
	agent storage.StorageAgent
	appid string
	key appkey.AppKey
	tags []string
	id string
}


func GetExport(appid string, key appkey.AppKey, tags []string,id string, agent storage.StorageAgent) Export {
	e :=  Export{agent,appid,key,tags,id}
	return e
}
func GetExportById(id string, agent storage.StorageAgent) Export {
	var e Export
	e.id = id
	e.agent = agent
	c := make(chan Export)
	e.agent.Read(func (sc *storage.StorageCore) {
		iv := sc.GetVertex(id)
		for _,appedge := range iv.Incoming() {
			if appedge.Label() == EXPORT_EDGE {
				e.appid = appedge.Tail().Id()
				break
			}
		}
		for _,keyedge := range iv.Outgoing() {
			if keyedge.Label() == EXPORT_EDGE {
				e.key = keyedge.Head().Properties().(appkey.AppKey)
				break
			}
		}
		c<-e
	})

	e = <-c
	e.tags = []string{}
	for _, tag := range e.GetTags() {
		e.tags = append(e.tags,tag.name)
	}
	return e
}

func (e *Export) Exists() bool {
	if e.id != "" {
		c := make(chan bool)
		defer close(c)
		e.agent.Read(func(sc *storage.StorageCore) {
			c <- sc.GetVertex(e.id)!=nil
		})
		return <-c
	}
	return false
}
func (e *Export) Add() {
	id := make(chan string)
	defer close(id)
	log.Println("STORAGECORE", "Adding export from", e.appid)

	a := GetApp(e.appid, e.agent)

	if !a.Exists() {
		log.Println("STORAGECORE", "Adding exporter failed, app does not exist:", e.appid)
		return
	}
	k := GetAppKey(e.key, e.agent)

	k.Create()
	keyid := k.GetId()
	if e.id == "" {
		e.id = storage.GenerateUuid()
	}
	e.agent.Write(func(sc *storage.StorageCore) {
		av := sc.InMemoryGraph.GetVertex(a.Id)
		kv := sc.InMemoryGraph.GetVertex(keyid)
		ev := sc.InMemoryGraph.CreateVertex(e.id, nil)


		sc.InMemoryGraph.CreateEdge(storage.GenerateUuid(), EXPORT_EDGE, kv, ev, nil)
		sc.InMemoryGraph.CreateEdge(storage.GenerateUuid(), EXPORT_EDGE, ev, av, nil)

		id <- ev.Id()
	})

	e.id = <-id

	if e.id != "" {
		for _, tag := range e.tags {

			t := GetTag(k,tag,e.agent)

			t.Create()

			t.LinkExport(e)


		}
	}

}

func (e *Export) GetApp() App{
	return GetApp(e.appid,e.agent)
}
func (e *Export) GetAppKey() AppKey{
	return GetAppKey(e.key,e.agent)
}
func (e *Export) setId() {
	keyid := e.GetAppKey().GetId()
	a := GetApp(e.appid, e.agent)
	if !a.Exists() {
		log.Println("STORAGECORE", "Setting exporterid failed, app does not exist:", e.appid)
		return
	}
	c := make(chan string)
	defer close(c)
	e.agent.Read(func (sc *storage.StorageCore){
		av := sc.GetVertex(e.appid)

		i := 0
		//app - EXPORTEDGE > Export
		for _,exportedge := range av.Outgoing(){
			if exportedge.Label() == EXPORT_EDGE {

				//Export - EXPORTEDGE > Key
				export := exportedge.Head()
				for _,exportkeyedge := range export.Outgoing() {
					if exportkeyedge.Label() == EXPORT_EDGE {
						if keyid == exportkeyedge.Head().Id() {

							//Export - TAGEDGE > Tag
							for _, tagedge := range export.Outgoing() {
								if tagedge.Label() == TAG_EDGE {
									tagname := tagedge.Head().Properties().(string)
									for _, tn := range e.tags {

										if tn == tagname {
											i++
											break
										}
									}
								}
							}

							if len(e.tags) == i {
								c <- export.Id()
								return
							}

						}
					}
				}
			}
		}
		c  <- ""
	})
	e.id = <- c
	return
}
func (e *Export) GetId() string {
	if e.id == "" {
		e.setId()
	}
	return e.id
}



	func (e *Export) UpdateTags(tags []string){
	e.ClearTags()
	e.tags = tags
	k := e.GetAppKey()
	for _, tag := range e.tags {
		t := GetTag(k,tag,e.agent)
		t.Create()
		t.LinkExport(e)
	}
}

func (e Export) ClearTags(){
	//key := GetAppKey(e.key,e.agent)
	for _, tag := range e.GetTags() {
		tag.UnlinkExport(e)
	}
}
func (e Export) GetTags() ([]Tag){
	tags := []Tag{}
	if e.GetId() == "" {
		return tags
	}
	k := GetAppKey(e.key, e.agent)
	c := make (chan []Tag)
	defer close(c)
	e.agent.Read(func (sc *storage.StorageCore){
		ev := sc.GetVertex(e.GetId())
		for _,tagedge := range ev.Outgoing(){
			if tagedge.Label() == TAG_EDGE {
				tagname,_ := tagedge.Head().Properties().(string)
				tags = append(tags,Tag{e.agent,k,tagname,tagedge.Head().Id()})
			}
		}
		c <- tags
	})
	return <-c
}


func (e Export) GetTagNames() []string {
	return e.tags
}
func (e Export) HasTags(tags []string) bool{
	mytags := e.GetTags()
	i:=0
	for _,tag := range tags {
		for _,mytag := range mytags {
			if mytag.name == tag {
				i++
				break
			}
		}
	}

	return i == len(tags)
}

func (e Export) Remove()  {
	c := make(chan bool)
	defer close(c)
	e.agent.Write(func (sc *storage.StorageCore){
		sc.RemoveVertex(e.id)
		c<-true
		return
	})
	 <- c
	return
}

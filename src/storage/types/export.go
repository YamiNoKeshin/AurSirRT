package types

import (
	"storage"
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


func GetExport(appid string, key appkey.AppKey, tags []string, agent storage.StorageAgent) Export {
	e :=  Export{agent,appid,key,tags,""}
	return e
}
func GetExportById(id string, agent storage.StorageAgent) Export {
	var e Export
	e.id = id
	e.agent = agent
	return e
}

func (e *Export) Exists() bool {
	if e.id != "" {
		c := make(chan bool)
		defer close(c)
		e.agent.Read(func(sc *storage.StorageCore) {
			c <- sc.GetVertex(e.id) != nil
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
	e.agent.Write(func(sc *storage.StorageCore) {
		av := sc.InMemoryGraph.GetVertex(a.Id)
		kv := sc.InMemoryGraph.GetVertex(keyid)
		ev := sc.InMemoryGraph.CreateVertex(storage.GenerateUuid(), nil)


		sc.InMemoryGraph.CreateEdge(storage.GenerateUuid(), storage.EXPORT_EDGE, kv, ev, nil)
		sc.InMemoryGraph.CreateEdge(storage.GenerateUuid(), storage.EXPORT_EDGE, ev, av, nil)

		id <- ev.Id
	})

	e.id = <-id

	if e.id != "" {
		for _, tag := range e.tags {

			t := GetTag(k,tag,e.agent)

			t.Create()

			t.LinkExport(*e)


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
		for _,exportedge := range av.Outgoing{
			if exportedge.Label == storage.EXPORT_EDGE {

				//Export - EXPORTEDGE > Key
				export := exportedge.Head
				for _,exportkeyedge := range export.Outgoing {
					if exportkeyedge.Label == storage.EXPORT_EDGE {
						if keyid == exportkeyedge.Head.Id {
							log.Println("STORAGECORE",len(export.Outgoing))

							//Export - TAGEDGE > Tag
							for _, tagedge := range export.Outgoing {
								if tagedge.Label == storage.TAG_EDGE {
									tagname := tagedge.Head.Properties.(string)
									for _, tn := range e.tags {

										if tn == tagname {
											i++
											break
										}
									}
								}
							}

							if len(e.tags) == i {
								c <- export.Id
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



	func (e Export) UpdateTags(tags []string){
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
		for _,tagedge := range ev.Outgoing{
			if tagedge.Label == storage.TAG_EDGE {
				tagname,_ := tagedge.Head.Properties.(string)
				tags = append(tags,Tag{e.agent,k,tagname,tagedge.Head.Id})
			}
		}
		c <- tags
	})
	return <-c
}
func (e Export) HasTags(tags []string) bool{
	mytags := e.GetTags()
	i:=0
	for tag := range tags {
		for mytag := range mytags {
			if mytag == tag {
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

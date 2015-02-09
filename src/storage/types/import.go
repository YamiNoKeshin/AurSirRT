package types

import (
	"aursirrt/src/storage"
	"log"

	"github.com/joernweissenborn/aursir4go/appkey"
)




type Import struct {
	agent storage.StorageAgent
	appid string
	key appkey.AppKey
	tags []string
	id string
}


func GetImport(appid string, key appkey.AppKey, tags []string, id string, agent storage.StorageAgent) Import {
	i :=  Import{agent,appid,key,tags,id}
	return i
}

func GetImportById(id string, agent storage.StorageAgent) Import {
	var i Import
	i.id = id
	i.agent = agent
	c := make(chan Import)
	defer close(c)
	i.agent.Read(func (sc *storage.StorageCore) {
		iv := sc.GetVertex(id)
		if iv == nil {
			c <-i
			return
		}
		for _,appedge := range iv.Incoming() {
			if appedge.Label() == IMPORT_EDGE {
				i.appid = appedge.Tail().Id()
				break
			}
		}	
		for _,keyedge := range iv.Outgoing() {
			if keyedge.Label() == IMPORT_EDGE {
				i.key = keyedge.Head().Properties().(appkey.AppKey)
				break
			}
		}
		c<-i
	})
	imp := <-c
	imp.tags = []string{}
	for _, tag := range imp.GetTags() {
		imp.tags = append(imp.tags,tag.name)
	}
	return imp
}


func (i *Import) GetTagNames() []string {
	 return i.tags
}
func (i *Import) Exists() bool {
	if i.id != "" {
		c := make(chan bool)
		defer close(c)
		i.agent.Read(func(sc *storage.StorageCore) {
			c <- sc.GetVertex(i.id) != nil
		})
		return <-c
	}
	return false

}


func (i *Import) Add() {
	id := make(chan string)
	defer close(id)
	a := GetApp(i.appid, i.agent)
	if !a.Exists() {
		log.Println("STORAGECORE", "Adding Importer failed, app does not exist:", i.appid)
		return
	}
	k := GetAppKey(i.key, i.agent)
	if i.id == "" {
		i.id = storage.GenerateUuid()
	}
	k.Create()
	keyid := k.GetId()
	i.agent.Write(func(sc *storage.StorageCore) {
		av := sc.InMemoryGraph.GetVertex(a.Id)
		kv := sc.InMemoryGraph.GetVertex(keyid)
		iv := sc.InMemoryGraph.CreateVertex(i.id, nil)

		sc.InMemoryGraph.CreateEdge(storage.GenerateUuid(), IMPORT_EDGE, kv, iv, nil)
		sc.InMemoryGraph.CreateEdge(storage.GenerateUuid(), IMPORT_EDGE, iv, av, nil)

		id <- iv.Id()
	})
	i.id = <-id

	if i.id != "" {
		for _, tag := range i.tags {
			t := GetTag(k,tag,i.agent)
			t.Create()
			t.LinkImport(i)
		}
	}
}

func (i *Import) GetApp() App{
	return GetApp(i.appid,i.agent)
}

func (i *Import) StartListenToFunction(Function string) {

	kid := i.GetAppKey().GetId()
	iid := i.GetId()
	i.agent.Write(func(sc *storage.StorageCore) {

		kv := sc.GetVertex(kid)
		iv := sc.GetVertex(iid)
		sc.InMemoryGraph.CreateEdge(storage.GenerateUuid(), LISTEN_EDGE, kv, iv, Function)

	})

}

func (i *Import) StopListenToFunction(Function string) {

	iid := i.GetId()
	i.agent.Write(func(sc *storage.StorageCore) {

		iv := sc.GetVertex(iid)
		for _,listenedge := range iv.Outgoing() {
			if listenedge.Label() == LISTEN_EDGE {
				if listenedge.Properties().(string) == Function {
					sc.RemoveEdge(listenedge.Id())
				}
			}
		}

	})

}

func (e *Import) GetAppKey() AppKey{
	return GetAppKey(e.key,e.agent)
}

func (e *Import) setId() {
	keyid := e.GetAppKey().GetId()
	a := GetApp(e.appid, e.agent)
	if !a.Exists() {
		log.Println("STORAGECORE", "Setting Importerid failed, app does not exist:", e.appid)
		return
	}
	c := make(chan string)
	defer close(c)
	e.agent.Read(func (sc *storage.StorageCore){
		av := sc.GetVertex(e.appid)
		i := 0
		//app - ImportEDGE > Import
		for _,Importedge := range av.Outgoing(){
			if Importedge.Label() == IMPORT_EDGE {
				//Import - ImportEDGE > Key
				Import := Importedge.Head()
				for _,Importkeyedge := range Import.Outgoing() {

					if Importkeyedge.Label() == IMPORT_EDGE {
						if keyid == Importkeyedge.Head().Id() {

							//Import - TAGEDGE > Tag
							for _, tagedge := range Import.Outgoing() {

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
								c <- Import.Id()
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
func (i *Import) GetId() string {
	if i.id == ""{
		i.setId()
	}
	return i.id
}

func (i *Import) UpdateTags(tags []string){
	i.ClearTags()
	i.tags = tags
	printDebug(" imp is ",i.id)
	printDebug(" imp is ",i.GetId())
	k := i.GetAppKey()
	for _, tag := range i.tags {
		t := GetTag(k,tag,i.agent)
		printDebug("tag imp is ",i.GetId())

		t.Create()
		t.LinkImport(i)
	}
}

func (i *Import) ClearTags(){
	//key := GetAppKey(e.key,e.agent)
	for _, tag := range i.GetTags() {
		tag.UnlinkImport(i)
	}
}

func (e Import) GetTags() ([]Tag){
	tags := []Tag{}
	if e.GetId() == "" {
		return tags
	}
	k := GetAppKey(e.key, e.agent)
	c := make (chan []Tag)
	defer close(c)
	iid := e.GetId()
	e.agent.Read(func (sc *storage.StorageCore){
		ev := sc.GetVertex(iid)
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
func (i Import) GetExporter() (exporter[]Export){
	exporter = []Export{}
	key := i.GetAppKey()

	for _,exp := range key.GetExporter(){
		if exp.HasTags(i.tags) {
			exporter = append(exporter,exp)
		}
	}
	return exporter
}
func (i Import) HasExporter() (bool){

	return len(i.GetExporter())!=0
}
func (i Import) GetJobs() (jobs []Job){
	jobs = []Job{}
	jobids := []string{}
	c := make(chan string)
	i.agent.Read(func (sc *storage.StorageCore){
		ev := sc.GetVertex(i.id)
		for _,tagedge := range ev.Outgoing(){
			if tagedge.Label() == AWAITING_JOB_EDGE {
				c<-tagedge.Head().Id()

			}
		}
		close(c)
	})

	for jid := range c {
		jobids = append(jobids, jid)
	}
		for _, jid := range jobids {
		jobs = append(jobs, GetJobById(jid,i.agent))
	}
	return
}
func (i Import) Remove()  {
	jobs := i.GetJobs()
	for _, job := range jobs {
		job.Remove()
	}
	c := make(chan bool)
	defer close(c)
	i.agent.Write(func (sc *storage.StorageCore){
		sc.RemoveVertex(i.id)
		c<-true
		return
	})
	 <- c
	return
}

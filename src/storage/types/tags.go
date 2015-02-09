package types

import (
	"aursirrt/src/storage"
	"fmt"
)

type Tag struct {
	agent storage.StorageAgent
	key AppKey
	name string
	id string
}

func GetTag(key AppKey, name string, agent storage.StorageAgent)Tag{
	tag := Tag{agent,key,name,""}
	tag.setId()

	return tag
}

func (t *Tag) setId() {
	c := make(chan string)
	defer close(c)
	keyid := t.key.GetId()
	t.agent.Read(func (sc *storage.StorageCore){
		kv := sc.GetVertex(keyid)
		if kv == nil {
			c <-""
			return
		}
		for _,e := range kv.Incoming(){
			if e.Label() == TAG_EDGE {
				tagname,_ := e.Tail().Properties().(string)
				if tagname == t.name {
					c<-e.Tail().Id()
					return
				}
			}
		}
		c <- ""
	})
	t.id = <- c
}
func (t *Tag) GetId() string{

	if t.id == ""{
		t.setId()
	}
	if !t.Exists() {
		return ""
	}

	return t.id



}

func (t *Tag) Create(){
	if !t.Exists() {
		printDebug(fmt.Sprint("Creating tag",t.name, t.key))

		c := make(chan string)
		defer close(c)
		keyid := t.key.GetId()
		printDebug(fmt.Sprint("tag keyid is ",keyid))

		t.agent.Write(func (sc *storage.StorageCore){

			tv := sc.CreateVertex(storage.GenerateUuid(), t.name)
			kv := sc.GetVertex(keyid)
			sc.CreateEdge(storage.GenerateUuid(), TAG_EDGE, kv,tv, nil)

			c <- tv.Id()
		})
		t.id = <-c

	}
}



func (t *Tag) LinkExport(e *Export){
	if t.Exists() {
		c := make(chan string)
		defer close(c)
		eid := e.GetId()
		printDebug(fmt.Sprint("linking tag and key ",t.id,eid))
		       tagid := t.GetId()
		t.agent.Write(func (sc *storage.StorageCore){
			ev := sc.GetVertex(eid)

			tv := sc.GetVertex(tagid)
			sc.CreateEdge(storage.GenerateUuid(), TAG_EDGE, tv, ev, nil)

			c <- ""

		})
		<-c
		printDebug(fmt.Sprint("linking tag and key sucess"))

	}
}

func (t *Tag) LinkImport(i *Import){
	if t.Exists() {
		c := make(chan string)
		defer close(c)
		iid := i.GetId()
		printDebug("linking tag and import ",t.id,iid)
		tagid := t.GetId()
		t.agent.Write(func (sc *storage.StorageCore){
			iv := sc.GetVertex(iid)

			tv := sc.GetVertex(tagid)
			sc.CreateEdge(storage.GenerateUuid(), TAG_EDGE, tv, iv, nil)

			c <- ""

		})
		<-c
		printDebug(fmt.Sprint("linking tag and key sucess"))

	}
}
func (t *Tag) UnlinkImport(i *Import){
	if t.Exists() {
		c := make(chan string)
		defer close(c)
		iid := i.GetId()
        tagid := t.GetId()
		delete := true
		t.agent.Write(func (sc *storage.StorageCore){

			tv := sc.GetVertex(tagid)
			for _, tagedge := range tv.Outgoing() {
				if tagedge.Label() == TAG_EDGE {
					if tagedge.Tail().Id() == iid {
						sc.RemoveEdge(tagedge.Id())
					} else {
						delete = false
					}
				}
			}
			if delete {
				sc.RemoveVertex(tagid)
			}
			c <- ""

		})
		<-c

	}
}
func (t *Tag) UnlinkExport(e Export){
	if t.Exists() {
		c := make(chan string)
		defer close(c)
		eid := e.GetId()
        tagid := t.GetId()
		delete := true
		t.agent.Write(func (sc *storage.StorageCore){

			tv := sc.GetVertex(tagid)
			for _, tagedge := range tv.Incoming() {
				if tagedge.Label() == TAG_EDGE {
					if tagedge.Tail().Id() == eid {
						sc.RemoveEdge(tagedge.Id())
					} else {
						delete = false
					}
				}
			}
			if delete {
				sc.RemoveVertex(tagid)
			}
			c <- ""

		})
		<-c

	}
}

func (t *Tag) Exists() bool {
	c := make(chan bool)
	defer close(c)
	keyid := t.key.GetId()

	t.agent.Read(func (sc *storage.StorageCore){
		kv := sc.GetVertex(keyid)

		if kv == nil {
			c <-false
			return
		}

		for _,e := range kv.Incoming(){

			if e.Label() == TAG_EDGE {

				tagname,_ := e.Tail().Properties().(string)
				if tagname == t.name {

					c<-true
					return
				}
			}
		}

		c <- false

	})
	return <- c
}

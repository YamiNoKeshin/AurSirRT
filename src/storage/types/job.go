package types

import (
	"aursirrt/src/storage"
	"fmt"
	"github.com/joernweissenborn/aursir4go/messages"
)

type jobproperties struct {
	request *messages.Request
	result *messages.Result
}
type Job struct {
	agent storage.StorageAgent
	request *messages.Request
	result *messages.Result
}

func GetJobFromRequest(request messages.Request, agent storage.StorageAgent)(Job Job){
	Job.agent = agent
	Job.request = &request
	return
}

func GetJobFromResult(result messages.Result, agent storage.StorageAgent)(Job Job){
	Job = GetJobById(result.Uuid,agent)

	Job.result = &result
	return
}
func GetJobById(id string, agent storage.StorageAgent)(Job Job){
	Job.agent = agent
	c := make(chan *messages.Request)
	c2 := make(chan *messages.Result)
	defer close(c)
	defer close(c2)
	Job.agent.Read(func (sc *storage.StorageCore){
		jv := sc.GetVertex(id)
		if jv == nil {
			c <- nil
			c2 <- nil
			return
		}
		c<-jv.Properties().(jobproperties).request
		c2<-jv.Properties().(jobproperties).result
		return
	})


	Job.request = <- c
	Job.result = <- c2
	return
}

func (j *Job) GetId() string{

	return j.request.Uuid

}

func (j Job) GetImport() Import{

	return GetImportById(j.request.ImportId,j.agent)

}
func (j Job) GetRequest() *messages.Request{

	return j.request

}
func (j Job) GetResult() *messages.Result{

	return j.result

}

func (j *Job) Assign(e Export) {
	eid := e.GetId()
	c := make(chan bool)
	j.agent.Write(func (sc *storage.StorageCore){

		jv := sc.GetVertex(j.request.Uuid)
		ev := sc.GetVertex(eid)
		sc.CreateEdge(storage.GenerateUuid(), DOING_JOB_EDGE, jv,ev, nil)

		c <- true
	})
	<-c
	return

}

func (j *Job) IsAssigned() bool{
	c := make(chan bool)
	defer close(c)
	j.agent.Read(func (sc *storage.StorageCore){
		jv := sc.GetVertex(j.request.Uuid)
		for _,doingedge := range jv.Incoming(){
			if doingedge.Label() == DOING_JOB_EDGE {
				c <- true
				return
			}
		}
		c<-false
	})

	return <-c

}
func (j *Job) GetAssignedExport() Export{
	c := make(chan string)
	defer close(c)
	j.agent.Read(func (sc *storage.StorageCore){
		jv := sc.GetVertex(j.request.Uuid)
		for _,doingedge := range jv.Incoming(){
			if doingedge.Label() == DOING_JOB_EDGE {
				c <- doingedge.Tail().Id()
				return
			}
		}
		c<-""
	})

	return GetExportById(<-c,j.agent)

}

func (j *Job) Create(){
	imp := GetImportById(j.request.ImportId,j.agent)
	if imp.Exists() {
		printDebug(fmt.Sprint("Creating Job"))

		c := make(chan string)
		defer close(c)

		j.agent.Write(func (sc *storage.StorageCore){

			jv := sc.CreateVertex(j.request.Uuid, jobproperties{j.request,j.result})
			iv := sc.GetVertex(j.request.ImportId)
			sc.CreateEdge(storage.GenerateUuid(), AWAITING_JOB_EDGE, jv,iv, nil)

			c <- ""
		})
		<-c

	}
}


func (j Job) Exists() bool {
	if j.request == nil {
		return false
	}
	c := make(chan bool)
	defer close(c)
	j.agent.Read(func (sc *storage.StorageCore){
		kv := sc.GetVertex(j.request.Uuid)
		if kv == nil {
			c <-false
			return
		}
		c<-true
		return
	})
	return <- c
}


func (j Job) Remove() bool {
	if !j.Exists(){
		return false
		
	}
	c := make(chan bool)
	defer close(c)
	j.agent.Write(func (sc *storage.StorageCore){
		sc.RemoveVertex(j.request.Uuid)
		c<-true
		return
	})
	return <- c
}

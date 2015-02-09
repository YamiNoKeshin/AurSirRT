package processor

import (
	"testing"
	"storage"
)

type testprocessor struct {
	*GenericProcessor
	c chan bool
}

func (tp testprocessor) Process(){
	tp.c <- true
}

type emptytestprocessor struct {
	*GenericProcessor
}

func (tp emptytestprocessor) Process(){
}

func TestBootProcessingCore(t *testing.T){
	var tp testprocessor
	tp.GenericProcessor = GetGenericProcessor()
	tp.c = make(chan bool)
	pc := make(chan Processor)
	go Process(pc,storage.NewAgent(), 1)
	defer close(pc)
	pc <- tp

	<-tp.c


}

func Test2Processors(t *testing.T){

	pc := Testprocessor()
	defer close(pc)
	var tp testprocessor
	tp.GenericProcessor = GetGenericProcessor()
	tp.c = make(chan bool)
	pc <- tp

	var tp2 testprocessor
	tp2.GenericProcessor = GetGenericProcessor()
	tp2.c = make(chan bool)
	pc <- tp2

	<-tp.c
	<-tp2.c


}

func Test3Processors(t *testing.T){

	pc :=Testprocessor()
	defer close(pc)
	var ep emptytestprocessor
	ep.GenericProcessor = GetGenericProcessor()
	pc <- ep
	pc <- ep
	pc <- ep
	pc <- ep

	var tp testprocessor
	tp.GenericProcessor = GetGenericProcessor()
	tp.c = make(chan bool)
	pc <- tp
	<-tp.c


}

type spawntestprocessor struct {
	*GenericProcessor
	c chan bool
}

func (tsp spawntestprocessor) Process(){
	var tp testprocessor
	tp.GenericProcessor = GetGenericProcessor()
	tp.c = tsp.c
	tsp.SpawnProcess(tp)
}

func TestSpawnProcess(t *testing.T){
	var tp spawntestprocessor
	tp.GenericProcessor = GetGenericProcessor()
	tp.c = make(chan bool)
	pc := Testprocessor()
	defer close(pc)
	pc <- &tp

	<-tp.c


}


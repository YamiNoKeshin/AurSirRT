package processors

import (
	"testing"
	"processor"
	"github.com/joernweissenborn/aursir4go/messages"
	"encoding/json"
	"dock/connection"
)

var Testdockmsg = messages.DockMessage{"testapp",[]string{"JSON"}}


func TestDockProcessor(t *testing.T){

	pc := processor.Testprocessor()
	defer close(pc)

	c := make(chan bool)
	defer close(c)
	conn := testdockconnection{c}

	pc <- GetTestDockProcessor("testid", conn)
	if !(<-c) {
		t.Error("Dockmessages says not ok")
	}


	pc <- GetTestDockProcessor("", conn)
	if (<-c) {
		t.Error("Dockmessages says  ok")
	}


}

type testdockconnection struct {
	c chan bool
}

func (tc testdockconnection) Init() (err error) {
	return
}
func (tc testdockconnection) Send(msgtype int64, codec string,msg []byte) (err error) {
	var m messages.DockedMessage
	json.Unmarshal(msg, &m)
	tc.c <- m.Ok
	return
}



func GetTestDockProcessor(appid string, conn connection.Connection) DockProcessor{
	var dp DockProcessor
	dp.GenericProcessor = processor.GetGenericProcessor()
	dp.AppId = appid
	dp.Codec = "JSON"
	dp.DockMessage,_ = json.Marshal(Testdockmsg)
	dp.Connection = conn
	return dp
}

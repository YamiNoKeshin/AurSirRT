package processors

import (
	"testing"
	"processor"
	"encoding/json"
	"github.com/joernweissenborn/aursir4go/messages"
	"github.com/joernweissenborn/aursir4go/Example/keys"
)


var Testaddimpmsg = messages.AddImportMessage{keys.HelloAurSirAppKey, []string{"one","two"}}



func TestAddImportProcessor(t *testing.T){

	pc := processor.Testprocessor()
	defer close(pc)
	c := make(chan bool)
	defer close(c)
	conn := testexpconnection{c}
	var dp DockProcessor
	dp.GenericProcessor = processor.GetGenericProcessor()
	dp.AppId = "testid"
	dp.Connection = &conn
	dp.DockMessage,_ = json.Marshal(Testdockmsg)
	 func (){pc <- dp}()

	var ap AddImportProcessor
	ap.GenericProcessor = processor.GetGenericProcessor()
	ap.AppId = "testid"
	ap.AddImportMsg = Testaddimpmsg

	 func (){pc <- ap}()
	<-c




}


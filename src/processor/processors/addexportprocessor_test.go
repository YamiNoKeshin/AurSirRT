package processors

import (
	"testing"
	"processor"
	"github.com/joernweissenborn/aursir4go/messages"
	"github.com/joernweissenborn/aursir4go/Example/keys"
)

var Testaddexpmsg = messages.AddExportMessage{keys.HelloAurSirAppKey, []string{"one","two"}}


func TestAddExportProcessor(t *testing.T){
	c := make(chan bool)
	defer close(c)
	conn := testexpconnection{c}

	pc := processor.Testprocessor()
	defer close(pc)

	dp := GetTestDockProcessor("testid", &conn)
	 func (){pc <- dp}()

	var ap AddExportProcessor
	ap.GenericProcessor = processor.GetGenericProcessor()
	ap.AppId = "testid"
	ap.AddExportMsg = Testaddexpmsg

	 func (){pc <- ap}()
	if !(<-c) {
		t.Error("Failed to create export")
	}


}


type testexpconnection struct {
	c chan bool
}

func (*testexpconnection) Init() (err error) {
	return
}
func (tc *testexpconnection) Send(msgtype int64, codec string,msg []byte) (err error) {

	tc.c <- true
	return
}

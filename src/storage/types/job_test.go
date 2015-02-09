package types
import (
	"testing"
	"aursirrt/src/storage"
	"github.com/joernweissenborn/aursir4go/messages"
	"time"
	"github.com/joernweissenborn/aursir4go/calltypes"
	"github.com/joernweissenborn/aursir4go/Example/keys"
)

func TestJob(t *testing.T){
	agent := storage.NewAgent()
	app := GetApp("testid",agent)
	dockmsg := messages.DockMessage{"HelloWorld",[]string{"JSON"},false}
	app.Create(dockmsg,testconn{})

	Import := GetImport("testid",keys.HelloAurSirAppKey, []string{"one","two"},"",agent)
	Import.Add()
	testrequest.ImportId = Import.GetId()
	job := GetJobFromRequest(testrequest,agent)
	job.Create()
	if !job.Exists() {
		t.Error("Could not create job")
	}

	if job.IsAssigned() {
		t.Error("Job is assigned w/o export")
	}
	eapp := GetApp("testexp",agent)
	eapp.Create(dockmsg,testconn{})

	export := GetExport("testexp",keys.HelloAurSirAppKey, []string{"one","two"},"",agent)
	export.Add()
	job.Assign(export)

	if !job.IsAssigned() {
		t.Error("Job is not assigned")
	}
	assignee := job.GetAssignedExport()
	if assignee.GetId()!=export.GetId() {
		t.Error("Could not retrieve assignee")
	}



}

var testrequest = messages.Request{
	"org.aursir.helloaursir",
	"SayHello",
	calltypes.ONE2ONE,
	[]string{"one","two"},
	"testjob",
	"",
	"",
	time.Now(),
	"JSON",
	[]byte{},
	false,
	true,
}

package types
import (
	"testing"
	"storage"
	"github.com/joernweissenborn/aursir4go"
	"github.com/joernweissenborn/aursir4go/messages"
)

func TestImportCreation(t *testing.T){
	agent := storage.NewAgent()
	app := GetApp("testid",agent)
	dockmsg := messages.DockMessage{"HelloWorld",[]string{"JSON"}}
	app.Create(dockmsg,testconn{})

	Import := GetImport("",aursir4go.HelloAurSirAppKey, []string{"one","two"},agent)
	Import.Add()

	if Import.GetId() != "" {
		t.Error("Created Import for non existing app")
	}
	Import = GetImport("testid",aursir4go.HelloAurSirAppKey, []string{"one","two"},agent)
	Import.Add()

	if Import.GetId() == "" {
		t.Error("Could not add Import")
	}
	Import = GetImport("testid",aursir4go.HelloAurSirAppKey, []string{"one","two"},agent)

	if Import.GetId() == "" {
		t.Error("Could not retrieve Import")
	}
}

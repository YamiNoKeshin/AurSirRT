package types
import (
	"testing"
	"aursirrt/src/storage"
	"github.com/joernweissenborn/aursir4go/messages"
	"github.com/joernweissenborn/aursir4go/Example/keys"
)

func TestImport(t *testing.T){
	agent := storage.NewAgent()
	app := GetApp("testid",agent)
	dockmsg := messages.DockMessage{"HelloWorld",[]string{"JSON"},false}
	app.Create(dockmsg,testconn{})

	Import := GetImport("",keys.HelloAurSirAppKey, []string{"one","two"},"",agent)
	Import.Add()

	if Import.GetId() != "" {
		t.Error("Created Import for non existing app")
	}
	Import = GetImport("testid",keys.HelloAurSirAppKey, []string{"one","two"},"",agent)
	Import.Add()

	if Import.GetId() == "" {
		t.Error("Could not add Import")
	}
	Import = GetImport("testid",keys.HelloAurSirAppKey, []string{"one","two"},"",agent)

	if Import.GetId() == "" {
		t.Error("Could not retrieve Import")
	}
	key := Import.GetAppKey()
	if !key.Exists(){
		t.Error("Could not create key")
	}

	if len(key.GetImporter()) == 0 {
		t.Error("Could not retrieve import from key")

	}

	if Import.HasExporter() {
		t.Error("Exporter should not be present")
	}
	eapp := GetApp("testexp",agent)
	eapp.Create(dockmsg,testconn{})

	export := GetExport("testexp",keys.HelloAurSirAppKey, []string{"one","two"},"",agent)
	export.Add()
	if !Import.HasExporter() {
		t.Error("Exporter should be present")
	}
	Import.UpdateTags([]string{"hi"})
	if Import.GetTagNames()[0] != "hi"{

		t.Error("Could not retrieve tag from key")

	}

}

func TestListen(t *testing.T){
	agent := storage.NewAgent()
	app := GetApp("testid",agent)
	dockmsg := messages.DockMessage{"HelloWorld",[]string{"JSON"},false}
	app.Create(dockmsg,testconn{})


	Import := GetImport("testid",keys.HelloAurSirAppKey, []string{"one","two"},"",agent)
	Import.Add()

	if Import.GetId() == "" {
		t.Error("Could not add Import")
	}

	eapp := GetApp("testexp",agent)
	eapp.Create(dockmsg,testconn{})

	export := GetExport("testexp",keys.HelloAurSirAppKey, []string{"one","two"},"",agent)
	export.Add()

	Import.StartListenToFunction("testfun")

	key := export.GetAppKey()
	if key.GetListener("testfun",export)[0].GetId() != Import.GetId() {
		t.Error("failed to listen")
	}
}


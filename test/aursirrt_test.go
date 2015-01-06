package test


import (
"log"
"testing"
"time"
	"github.com/joernweissenborn/aursir4go"
)

func TestInitCloseIface(t *testing.T) {
	iface := aursir4go.NewInterface("test")
	defer iface.Close()
	for !iface.Connected() {
		time.Sleep(1 * time.Millisecond)
	}

}
func TestExportKey(t *testing.T) {
	iface, exp := testexporter()
	defer iface.Close()
	log.Println("ExportId", exp)
	time.Sleep(1 * time.Second)

}

func TestImportKey(t *testing.T) {
	iface, imp := testimporter()
	defer iface.Close()
	time.Sleep(1 * time.Second)
	log.Println("Import", imp)

}

func TestKeyAvailable(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()
	exporter, _ := testexporter()
	time.Sleep(100 * time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}
	exporter.Close()
	time.Sleep(100 * time.Millisecond)
	if imp.Connected == true {
		T.Error("could not disconnect from appkey")
	}

}

func TestFunCall121(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()
	exporter, exp := testexporter()
	defer exporter.Close()

	res, _ := imp.CallFunction(aursir4go.HelloAurSirAppKey.Functions[0].Name, aursir4go.SayHelloReq{"AHOI"}, aursir4go.ONE2ONE)
	req := <-exp.Request
	var SayHelloReq aursir4go.SayHelloReq
	req.Decode(&SayHelloReq)
	if SayHelloReq.Greeting != "AHOI" {
		T.Error("got wrong request parameter")
	}
	err := exp.Reply(&req, aursir4go.SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var result aursir4go.SayHelloRes
	asres := <-res
	err = asres.Decode(&result)
	if err != nil {
		T.Error(err)
	}
	log.Println(result)
	if result.Answer != "MOINSEN" {
		T.Error("got wrong result parameter")
	}
}
func Test10FunCall121(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()
	exporter, exp := testexporter()
	defer exporter.Close()
	for i:=0;i<10;i++ {
	res, _ := imp.CallFunction(aursir4go.HelloAurSirAppKey.Functions[0].Name, aursir4go.SayHelloReq{"AHOI"}, aursir4go.ONE2ONE)
	req := <-exp.Request
	var SayHelloReq aursir4go.SayHelloReq
	req.Decode(&SayHelloReq)
	if SayHelloReq.Greeting != "AHOI" {
		T.Error("got wrong request parameter")
	}
	err := exp.Reply(&req, aursir4go.SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var result aursir4go.SayHelloRes
	asres := <-res
	err = asres.Decode(&result)
	if err != nil {
		T.Error(err)
	}
	log.Println(result)
	if result.Answer != "MOINSEN" {
		T.Error("got wrong result parameter")
	}
	}
}

func TestFunCallN2N(T *testing.T) {
	importer1, imp1 := testimporter()
	defer importer1.Close()
	imp1.ListenToFunction("SayHello")
	importer2, imp2 := testimporter()
	imp2.ListenToFunction("SayHello")
	defer importer2.Close()
	exporter1, exp1 := testexporter()
	defer exporter1.Close()
exporter2, exp2 := testexporter()
	defer exporter2.Close()

	imp2.CallFunction(aursir4go.HelloAurSirAppKey.Functions[0].Name, aursir4go.SayHelloReq{"AHOI"}, aursir4go.MANY2MANY)
	req := <-exp1.Request
	var SayHelloReq aursir4go.SayHelloReq
	req.Decode(&SayHelloReq)
	if SayHelloReq.Greeting != "AHOI" {
		T.Error("got wrong request parameter")
	}
	err := exp1.Reply(&req, aursir4go.SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	req = <-exp2.Request
	req.Decode(&SayHelloReq)
	if SayHelloReq.Greeting != "AHOI" {
		T.Error("got wrong request parameter")
	}
	err = exp2.Reply(&req, aursir4go.SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var res1 aursir4go.SayHelloRes
	imp1.Listen().Decode(&res1)
	log.Println("res1", res1)
	if res1.Answer != "MOINSEN" {
		T.Error("got wrong result parameter")
	}
	var res2 aursir4go.SayHelloRes
	imp2.Listen().Decode(&res2)
	log.Println("res2", res2)
	if res2.Answer != "MOINSEN" {
		T.Error("got wrong result parameter")
	}
}

func TestFunCallN21(T *testing.T) {
	importer1, imp1 := testimporter()
	defer importer1.Close()
	imp1.ListenToFunction("SayHello")
	importer2, imp2 := testimporter()
	imp2.ListenToFunction("SayHello")
	defer importer2.Close()
	exporter, exp := testexporter()
	defer exporter.Close()

	imp2.CallFunction(aursir4go.HelloAurSirAppKey.Functions[0].Name, aursir4go.SayHelloReq{"AHOI"}, aursir4go.MANY2ONE)
	req := <-exp.Request
	var SayHelloReq aursir4go.SayHelloReq
	req.Decode(&SayHelloReq)
	if SayHelloReq.Greeting != "AHOI" {
		T.Error("got wrong request parameter")
	}
	err := exp.Reply(&req, aursir4go.SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var res1 aursir4go.SayHelloRes
	imp1.Listen().Decode(&res1)
	log.Println("res1", res1)
	if res1.Answer != "MOINSEN" {
		T.Error("got wrong result parameter")
	}
	var res2 aursir4go.SayHelloRes
	imp2.Listen().Decode(&res2)
	log.Println("res2", res2)
	if res2.Answer != "MOINSEN" {
		T.Error("got wrong result parameter")
	}
}

func TestDelayedExporter(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()

	res, _ := imp.CallFunction(aursir4go.HelloAurSirAppKey.Functions[0].Name, aursir4go.SayHelloReq{"AHOI"}, aursir4go.ONE2ONE)
	exporter, exp := testexporter()
	defer exporter.Close()
	req := <-exp.Request
	var SayHelloReq aursir4go.SayHelloReq
	req.Decode(&SayHelloReq)
	log.Println(SayHelloReq)

	if SayHelloReq.Greeting != "AHOI" {
		T.Error("got wrong request parameter")
	}
	err := exp.Reply(&req, aursir4go.SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var result aursir4go.SayHelloRes
	(<-res).Decode(&result)
	log.Println(result)
	if result.Answer != "MOINSEN" {
		T.Error("got wrong result parameter")
	}
}
//
//func TestExporterImporterIfaceBlocking(T *testing.T) {
//	importer, imp := testimporter()
//	defer importer.Close()
//	exporter1, exp1 := testexporter()
//	defer exporter1.Close()
//	exporter2, exp2 := testexporterctrstr()
//	defer exporter2.Close()
//
//	imp2:= exporter1.AddImport(aursir4go.CountStringKey, nil)
//
//	res1, _ := imp2.CallFunction(aursir4go.CountStringKey.Functions[0].Name, aursir4go.CountStringReq{"AHOI1"}, aursir4go.ONE2ONE)
//	res2, _ := imp.CallFunction(aursir4go.HelloAurSirAppKey.Functions[0].Name, aursir4go.SayHelloReq{"AHOI2"}, aursir4go.ONE2ONE)
//	req := <-exp2.Request
//
//	exp2.Reply(&req, aursir4go.CountStringRes{1})
//	<-res1
//	req = <-exp1.Request
//	exp1.Reply(&req, aursir4go.SayHelloRes{"MOINSEN"})
//
//	<-res2
//
//}



func TestTagging(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()
	exporter, exp := testexporter()
	defer exporter.Close()
	time.Sleep(100 * time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}

	imp.UpdateTags([]string{"testtag"})
	time.Sleep(300 * time.Millisecond)
	if imp.Connected == true {
		T.Error("could not disconnect from appkey")
	}
	exp.UpdateTags([]string{"testtag"})

	time.Sleep(300 * time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}
	exp.UpdateTags([]string{"testtag", "anothertag"})

	time.Sleep(300 * time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}

	exp.UpdateTags([]string{"anothertag"})
	time.Sleep(300 * time.Millisecond)
	if imp.Connected == true {
		T.Error("could not disconnect from appkey")
	}
	imp.UpdateTags([]string{})
	time.Sleep(300 * time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}
}



func TestCallChain(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()

	cc, _ := imp.NewCallChain(aursir4go.HelloAurSirAppKey.Functions[0].Name, aursir4go.SayHelloReq{"AHOI"}, aursir4go.ONE2ONE)
	paramap := map[string]string{}
	paramap["String"] = "Answer"

	cc.AddCall("org.aursir.countstring", "CountString", paramap, aursir4go.ONE2ONE, []string{})
	err := cc.Finalize()
	log.Println(err)
//	if err == nil {
//		T.Error("Finalize should have thrown err now")
//	}
	exporter1, exp1 := testexporterctrstr()
	defer exporter1.Close()

	err = cc.Finalize()
	if err != nil {
		T.Error("Finalize should not have thrown err now")
	}
	exporter2, exp2 := testexporter()
	defer exporter2.Close()

	req := <-exp2.Request
	var SayHelloReq aursir4go.SayHelloReq
	req.Decode(&SayHelloReq)
	log.Println(SayHelloReq)

	if SayHelloReq.Greeting != "AHOI" {
		T.Error("got wrong request parameter")
	}

	err = exp2.Reply(&req, aursir4go.SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var csr aursir4go.CountStringReq
	req = <-exp1.Request

	req.Decode(&csr)
	log.Println(csr)

	if csr.String != "MOINSEN" {
		T.Error("got wrong request parameter")
	}
	err = exp2.Reply(&req, aursir4go.CountStringRes{int64(len([]byte(csr.String)))})
	if err != nil {
		T.Error(err)
	}
}

func TestCallChainFinalize(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()
	imp2 := importer.AddImport(aursir4go.CountStringKey, []string{})
	cc, _ := imp.NewCallChain(aursir4go.HelloAurSirAppKey.Functions[0].Name, aursir4go.SayHelloReq{"AHOI"}, aursir4go.ONE2ONE)
	paramap := map[string]string{}
	paramap["String"] = "Answer"

	exporter1, exp1 := testexporterctrstr()
	defer exporter1.Close()

	exporter2, exp2 := testexporter()
	defer exporter2.Close()

	rep, err := imp2.FinalizeCallChain(aursir4go.CountStringKey.Functions[0].Name, paramap, aursir4go.ONE2ONE, cc)

	if err != nil {
		T.Error("Finalize should not have thrown err now")
	}
	req := <-exp2.Request
	var SayHelloReq aursir4go.SayHelloReq
	req.Decode(&SayHelloReq)
	log.Println(SayHelloReq)

	if SayHelloReq.Greeting != "AHOI" {
		T.Error("got wrong request parameter")
	}

	err = exp2.Reply(&req, aursir4go.SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var csr aursir4go.CountStringReq
	req = <-exp1.Request

	req.Decode(&csr)
	log.Println(csr)

	if csr.String != "MOINSEN" {
		T.Error("got wrong request parameter")
	}

	err = exp2.Reply(&req, aursir4go.CountStringRes{int64(len([]byte(csr.String)))})
	if err != nil {
		T.Error(err)
	}

	rply := <-rep
	var csrep aursir4go.CountStringRes
	rply.Decode(&csrep)

	if csrep.Size != int64(len([]byte(csr.String))) {
		T.Error("got wrong result parameter")
	}
}


func TestPersitenceLogging(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()
	exporter, exp := testexporter()
	defer exporter.Close()
	exp.SetLogging("SayHello")
	res, _ := imp.CallFunction(aursir4go.HelloAurSirAppKey.Functions[0].Name,
		aursir4go.SayHelloReq{"AHOI"}, aursir4go.ONE2ONE)
	req := <-exp.Request

	exp.Reply(&req, aursir4go.SayHelloRes{"MOINSEN"})
	<-res
}

func TestStream(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()
	exporter, exp := testexporter()
	defer exporter.Close()

	res, _ := imp.CallFunction(aursir4go.HelloAurSirAppKey.Functions[0].Name, aursir4go.SayHelloReq{"AHOI"}, aursir4go.ONE2ONE)
	req := <-exp.Request
	var SayHelloReq aursir4go.SayHelloReq
	req.Decode(&SayHelloReq)
	if SayHelloReq.Greeting != "AHOI" {
		T.Error("got wrong request parameter")
	}
	err := exp.StreamingReply(&req, aursir4go.SayHelloRes{"MOINSEN"},false)
	if err != nil {
		T.Error(err)
	}
	err = exp.StreamingReply(&req, aursir4go.SayHelloRes{"ZUMZWEITEN"},true)
	if err != nil {
		T.Error(err)
	}
	var result aursir4go.SayHelloRes
	asres := <-res
	asres.Decode(&result)
	log.Println(result)
	if result.Answer != "MOINSEN" {
		T.Error("got wrong result parameter")
	}
	if  !asres.Stream || asres.Finished {
		T.Error("got wrong stream flags")
	}
	asres = <-res

	asres.Decode(&result)
	if result.Answer != "ZUMZWEITEN" {
		T.Error("got wrong result parameter")
	}
	if  !asres.Stream || !asres.Finished {
		T.Error("got wrong stream flags")
	}
}



func testexporter() (*aursir4go.AurSirInterface, *aursir4go.ExportedAppKey) {
	iface := aursir4go.NewInterface("testex")

	exp := iface.AddExport(aursir4go.HelloAurSirAppKey, nil)
	return iface, exp

}
func testimporter() (*aursir4go.AurSirInterface, *aursir4go.ImportedAppKey) {
	iface := aursir4go.NewInterface("testimp")
	imp := iface.AddImport(aursir4go.HelloAurSirAppKey, nil)
	return iface, imp
}

func testexporterctrstr() (*aursir4go.AurSirInterface, *aursir4go.ExportedAppKey) {
	iface := aursir4go.NewInterface("testex")

	exp := iface.AddExport(aursir4go.CountStringKey, nil)
	return iface, exp

}

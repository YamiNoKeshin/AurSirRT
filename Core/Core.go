/*
	Package Core is responsible for processing incoming messages
*/
package core

import (
	"github.com/joernweissenborn/aursir4go"
	"github.com/joernweissenborn/aursirrt/storagecore"
	"log"
	"github.com/joernweissenborn/aursirrt/config"
	"time"
	"github.com/joernweissenborn/aursirrt/datastorage"
)

func Launch(AppInChannel, AppOutChannel chan AppMessage,cfg config.RtConfig) {

	log.Println("Core Launching")

	var c core

	c.datastorage = make(chan interface {})

	go datastorage.Open(cfg,c.datastorage)

	c.storageAgent.Launch(cfg)

	c.appOutChannel = AppOutChannel

	go c.routeIncomingAppMsg(AppInChannel)

}

type core struct {
	storageAgent storagecore.StorageCoreAgent
	appOutChannel chan AppMessage
	datastorage chan interface {}
}

func (c core) routeIncomingAppMsg(appInChannel chan AppMessage) {

	for AppMessage := range appInChannel {

		aursirMessage, err := AppMessage.AppMsg.Decode()
		//log.Println("DEBUG",aursirMessage)
		//log.Println("DEBUG",string(AppMessage.AppMsg.Msg),err)


		if err == nil {

			switch aursirMessage := aursirMessage.(type) {

			case aursir4go.AurSirDockMessage:
				go c.dock(AppMessage.SenderUUID, aursirMessage)

			case aursir4go.AurSirLeaveMessage:
				go c.leave(AppMessage.SenderUUID)

			case aursir4go.AurSirAddExportMessage:
				go c.addExport(AppMessage.SenderUUID, aursirMessage)

			case  aursir4go.AurSirUpdateExportMessage:
				go c.updateExport(AppMessage.SenderUUID,aursirMessage)

			case aursir4go.AurSirAddImportMessage:
				go c.addImport(AppMessage.SenderUUID, aursirMessage)
			case  aursir4go.AurSirUpdateImportMessage:
				go c.updateImport(AppMessage.SenderUUID,aursirMessage)

			case aursir4go.AurSirRequest:
				go c.request(AppMessage.SenderUUID,aursirMessage)

			case aursir4go.AurSirResult:
				go c.result(AppMessage.SenderUUID,aursirMessage)
			case aursir4go.AurSirListenMessage:
				go c.listen(AppMessage.SenderUUID,aursirMessage)
			case aursir4go.AurSirStopListenMessage:
				go c.stopListen(AppMessage.SenderUUID,aursirMessage)

			case aursir4go.AurSirCallChain:
				go c.callChain(AppMessage.SenderUUID,aursirMessage)

			default:
				log.Println("unknown message")
			}
		}
	}
}

func (c core) callChain(senderId string,ccmsg aursir4go.AurSirCallChain) {
	log.Println("Processing CALL_CHAIN request from", senderId)


	oak := ccmsg.OriginRequest.AppKeyName
	ofn := ccmsg.OriginRequest.FunctionName
	ok := true
	insane := []int64{}
	for i,call := range ccmsg.CallChain {
		tak := call.AppKeyName
		tfn := call.FunctionName
		if !c.chainChecker(oak,ofn,tak,tfn,call.ArgumentMap){
			ok = false
			insane = append(insane,int64(i))
		}

		oak = tak
		ofn = tfn
	}
	log.Println("CALLCHAIN is vailid:",ok, insane)

	var ccm aursir4go.AppMessage
	ccm.Encode(aursir4go.AurSirCallChainAddedMessage{ok,insane},"JSON")
	c.appOutChannel <-AppMessage{senderId,ccm}

	if ok {
		exports, ok := (c.storageAgent.Write(storagecore.AddCallChainRequest{senderId,ccmsg})).([]string)
		if ok {
			for _, export := range exports {
				var reqmsg aursir4go.AppMessage
				reqmsg.Encode(ccmsg.OriginRequest,"JSON")
				c.appOutChannel <-AppMessage{export,reqmsg}
			}
		}
	}
}

func (c core) chainChecker(orgAppKey, orgFun, tarAppKey, tarFun string, paramap map[string]string) bool {
	oak, f := (c.storageAgent.Read(storagecore.GetAppKey{orgAppKey})).(aursir4go.AppKey)
	log.Println(orgFun)
	if !f {
		return false
	}
	tak, f := c.storageAgent.Read(storagecore.GetAppKey{tarAppKey}).(aursir4go.AppKey)
	if !f  {
		return false
	}
	log.Println(tak)

	var ofn aursir4go.Function
	f = false
	for _,fkt := range oak.Functions{
		if fkt.Name == orgFun {
			ofn = fkt
			f = true
		}
	}
	log.Println(ofn)

	if !f {return false}
	tmp := map[string]int{}
	log.Println(len(paramap))
	if len(paramap) == 0 {
		f =true
	} else {
		for input, output := range paramap {
			f = false
			for _, out := range ofn.Output {
				//log.Println(output)
				//log.Println(out.Name)
				if out.Name == output {
					f = true
					tmp[input] = out.Type
				}
			}
		}
	}
	if !f {return false}

	var tfn aursir4go.Function
	f = false
	for _,fkt := range tak.Functions{
		if fkt.Name == tarFun {
			tfn = fkt
			f = true
		}
	}
	//log.Println(tfn,f)

	if !f {return false}
	for _,in :=range tfn.Input {
		t,f := tmp[in.Name]
		//log.Println(t,f)
		//log.Println(tmp)
		//log.Println(in.Name)
		if !f || t != in.Type {
			return false
		}
	}
	//log.Println(f)

	return true
}


func (c core) updateImport(senderId string,uimsg aursir4go.AurSirUpdateImportMessage) {
	log.Println("Processing UPDATE_IMPORT request from", senderId)
	reply := c.storageAgent.Write(storagecore.UpdateImportRequest{uimsg})
	imp, ok := reply.(storagecore.ImportAdded)

	if ok {
		var em aursir4go.AppMessage
		em.Encode(aursir4go.AurSirImportUpdatedMessage{imp.ImportId, imp.Exported}, "JSON")
		c.appOutChannel <- AppMessage{senderId, em}
	}
}

func (c core) updateExport(senderId string,uemsg aursir4go.AurSirUpdateExportMessage){
	log.Println("Processing UPDATE_EXPORT request from", senderId)
	reply := c.storageAgent.Write(storagecore.UpdateExportRequest{uemsg})
	export, ok := reply.(storagecore.ExportAdded)
	//log.Println(export)
	if ok {
		for imp, appid := range export.ConnectedImports {
			var em aursir4go.AppMessage
			em.Encode(aursir4go.AurSirImportUpdatedMessage{imp, true}, "JSON")
			c.appOutChannel <- AppMessage{appid, em}
			log.Println("sending to",appid)
		}


		for imp, appid := range export.DisconnectedImports {
			var em aursir4go.AppMessage
			em.Encode(aursir4go.AurSirImportUpdatedMessage{imp, false}, "JSON")
			c.appOutChannel <- AppMessage{appid, em}
		}
		for _, r := range export.PendingJobs{
			var rm aursir4go.AppMessage
			rm.Encode(r,"JSON")
			c.appOutChannel <-AppMessage{senderId,rm}
		}
	}
}

func (c core) listen(senderId string, lmsg aursir4go.AurSirListenMessage){
	log.Println("Processing LISTEN request from", senderId)

	c.storageAgent.Write(storagecore.ListenRequest{senderId,lmsg.FunctionName,lmsg.ImportId})
}

func (c core) stopListen(senderId string, slmsg aursir4go.AurSirStopListenMessage){
	log.Println("Processing STOP_LISTEN request from", senderId)

}

func (c core) result(senderId string, rmsg aursir4go.AurSirResult) {
	log.Println("Processing RESULT request from", senderId)
	if rmsg.Persistent {
		var request aursir4go.AurSirRequest
		getReq := c.storageAgent.Read(storagecore.GetRequest{rmsg.Uuid})
		if _,f := getReq.(storagecore.ReadFail);!f{
			request = getReq.(aursir4go.AurSirRequest)
		}
		go c.persistResult(request,rmsg)
	}

	reply := c.storageAgent.Write(storagecore.AddResRequest{senderId,rmsg})
	resReg, ok := reply.(storagecore.ResRegistered)

	if ok && len(resReg.Importer)!=0{
		for imp, codecs := range resReg.Importer{

			//check if we need to recode
			recode := true
			for _, codec := range codecs {
				if codec == rmsg.Codec {
					recode = false
					break
				}
			}
			var rm aursir4go.AppMessage

			if recode {
				log.Println("CORE","Recoding")
				var tmp map[string]interface {}
				log.Println(codecs)
				src := aursir4go.GetCodec(rmsg.Codec)
				target := aursir4go.GetCodec(codecs[0])
				src.Decode(rmsg.Result, &tmp)
				newres := rmsg
				newres.Codec = codecs[0]
				newres.Result,_ = target.Encode(tmp)
				rm.Encode(newres,"JSON")
			} else {
				rm.Encode(rmsg,"JSON")
			}

			c.appOutChannel <-AppMessage{imp,rm}
		}
	}

	if ok && resReg.IsChainCall{
		log.Println("CORE","Creating next call in chain")

		c.createChainCall(senderId,rmsg,resReg.ChainCall,resReg.ChainCallImportId)
	}
}


func (c core) persistResult(req aursir4go.AurSirRequest, res aursir4go.AurSirResult) {
	answer := make(chan string)
	c.datastorage <- datastorage.CommitRequest{answer,datastorage.CommitData{&req,&res}}
	path := <-answer
	res.IsFile = true
	res.Result = []byte(path)
	c.storageAgent.Write(storagecore.AddPersistentResultRequest{res})
}

func (c core) createChainCall(senderId string, prevResult aursir4go.AurSirResult,cc aursir4go.ChainCall,ccImportId string){

	codec := aursir4go.GetCodec(prevResult.Codec)
	if codec==nil {
		log.Println("CORE","Error creating next call in chain, missing Codec in PrevResult")

		return
	}
	var tmp interface {}

	codec.Decode(prevResult.Result,&tmp)
	requestParameter := map[string]interface {}{}

	resultParameter, ok := tmp.(map[interface{}]interface {})
	if !ok {
		resultParameter2, _ := tmp.(map[string]interface {})
		for target, origin := range cc.ArgumentMap {
			requestParameter[target] = resultParameter2[origin]
		}
	}else{
		for target, origin := range cc.ArgumentMap {
			requestParameter[target] = resultParameter[origin]
		}
	}

	//log.Println(requestParameter)
	req , err:=codec.Encode(requestParameter)
	if err != nil {
		return
	}

	if cc.Tags == nil {
		cc.Tags = []string{}
	}
	request := aursir4go.AurSirRequest{
		cc.AppKeyName,
		cc.FunctionName,
		cc.CallType,
		cc.Tags,
		cc.ChainCallId,
		ccImportId,
		time.Now(),
		prevResult.Codec,
		false,
		false,
		"",
		req,
		prevResult.Stream,
		prevResult.StreamFinished}
	//	log.Println("ChainCalling",request)
	//	log.Println("ChainCalling",string(request.Request))
	reqReg, ok := c.storageAgent.Write(storagecore.AddReqRequest{senderId,request}).(storagecore.ReqRegistered)
	//log.Println(ok,reqReg)
	if ok{
		log.Println("CORE","Sucessfully created callchain, transmitting")

		c.transmitRequests(request,reqReg)
	}
}

func (c core) transmitRequests(request aursir4go.AurSirRequest, targets storagecore.ReqRegistered){
	for exp, codecs := range targets.Exporter {
		//check if we need to recode
		recode := true
		for _, codec := range codecs {
			if codec == request.Codec {
				recode = false
				break
			}
		}
		var rm aursir4go.AppMessage

		if recode{
			var err error
			log.Println("CORE","recoding for transmission to",codecs[0])
			var tmp map[string]interface {}
			src := aursir4go.GetCodec(request.Codec)
			target := aursir4go.GetCodec(codecs[0])
			src.Decode(request.Request, &tmp)
			newreq := request
			newreq.Codec = codecs[0]
			if newreq.Tags == nil {
				newreq.Tags = []string{}
			}
			//log.Println("CORE","Error transmitting, recoding failed", string(request.Request))
			//log.Println("CORE","Error transmitting, recoding failed", tmp)

			newreq.Request,err  = target.Encode(tmp)
			if err != nil{
				log.Println("CORE","Error transmitting, recoding failed", err)

			}
			log.Println("CORE","Sucessfully recoded request")

			rm.Encode(newreq,"JSON")
		} else {
			rm.Encode(request,"JSON")
		}
		log.Println("CORE","Transmitting")

		c.appOutChannel <-AppMessage{exp, rm}
	}
}

func (c core) request(senderId string, rmsg aursir4go.AurSirRequest) {
	log.Println("Processing REQUEST request from", senderId)
	reply := c.storageAgent.Write(storagecore.AddReqRequest{senderId,rmsg})
	reqreg, ok := reply.(storagecore.ReqRegistered)
	if ok && len(reqreg.Exporter)!=0{
		c.transmitRequests(rmsg,reqreg)
	}
}


func (c core) dock(senderId string, dmsg aursir4go.AurSirDockMessage) {
	log.Println("Processing DOCK request from", senderId)
	c.storageAgent.Write(storagecore.RegisterAppRequest{
	senderId,
	dmsg})
	var dm aursir4go.AppMessage
	dm.Encode(aursir4go.AurSirDockedMessage{}, "JSON")
	c.appOutChannel <- AppMessage{senderId, dm}
}

func (c core) leave(senderId string) {
	log.Println("Processing LEAVE request from", senderId)
	reply := c.storageAgent.Write(storagecore.RemoveAppRequest{
	senderId})
	leave, ok := reply.(storagecore.AppRemoved)
	if ok {

		for imp, appid := range leave.DisconnectedImports {
			var em aursir4go.AppMessage
			em.Encode(aursir4go.AurSirImportUpdatedMessage{imp, false}, "JSON")
			c.appOutChannel <- AppMessage{appid, em}
		}
	}

}

func (c core) addExport(senderId string, expMsg aursir4go.AurSirAddExportMessage) {
	log.Println("Processing ADD_EXPORT request from", senderId)
	reply := c.storageAgent.Write(storagecore.AddExportRequest{senderId, expMsg.AppKey, expMsg.Tags})
	export, ok := reply.(storagecore.ExportAdded)

	if ok {
		var em aursir4go.AppMessage
		em.Encode(aursir4go.AurSirExportAddedMessage{export.ExportId}, "JSON")
		c.appOutChannel <- AppMessage{senderId, em}
		for imp, appid := range export.ConnectedImports {
			var em aursir4go.AppMessage
			em.Encode(aursir4go.AurSirImportUpdatedMessage{imp, true}, "JSON")
			c.appOutChannel <- AppMessage{appid, em}
		}

		for _, r := range export.PendingJobs{
			var rm aursir4go.AppMessage

			rm.Encode(r,"JSON")


			c.appOutChannel <-AppMessage{senderId,rm}
		}
	}
}

func (c core) addImport(senderId string, expMsg aursir4go.AurSirAddImportMessage) {
	log.Println("Processing ADD_IMPORT request from", senderId)
	reply := c.storageAgent.Write(storagecore.AddImportRequest{senderId, expMsg.AppKey, expMsg.Tags})
	imp, ok := reply.(storagecore.ImportAdded)

	if ok {
		var em aursir4go.AppMessage
		em.Encode(aursir4go.AurSirImportAddedMessage{imp.ImportId, imp.Exported, expMsg.AppKey.ApplicationKeyName,expMsg.Tags}, "JSON")
		c.appOutChannel <- AppMessage{senderId, em}
	}
}

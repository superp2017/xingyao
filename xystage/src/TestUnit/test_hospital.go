package main

import (
	"JsLib/JsDispatcher"
	"cache/cacheIO"
	"constant"

	// . "JsLib/JsLogger"
	"JsLib/JsNet"
	// . "cache/cacheLib"
	"common"
	// "constant"
	. "util"
)

func Init_hospital() {
	JsDispatcher.Http("/getglobalhospital", GetGlobalHospotal)                 //新建订单
	JsDispatcher.Http("/gethospitalsimpleinfo", common.GetHosSimpleInfo)       //获取医院的简短信息
	JsDispatcher.Http("/refreshhospitalsimpleinfo", RefreshHospitalSimpleInfo) //刷新医院简短信息

	JsDispatcher.Http("/changedocvisit", Changedoc) ////
	JsDispatcher.Http("/changehosvisit", Changehos) ////

	JsDispatcher.Http("/ChangeHosStatistic", ChangeHosStatistic) //更新医院统计
}

func  ChangeHosStatistic(session *JsNet.StSession)  {
	list := common.GlobalHospitalList()
	for _,v:= range list{
		d:=&common.ST_Hospital{}
		if err:=WriteLock(constant.Hash_Hospital,v,d);err!=nil{
			continue
		}
		order := []string{}
		TradeNum:=0
		if err := ShareLock(constant.Hash_HospitalOrder, v, &order); err == nil {
			orderlist := common.QueryMoreOrders(order)
			for _, o := range orderlist {
				TradeNum += o.HosPayRealPrice
			}
		}
		doc:=[]string{}
		ShareLock(constant.Hash_DoctorCache,v,&doc)
		pro :=[]string{}
		ShareLock(constant.Hash_HosProduct,v,&pro)
		d.Orderquantity=len(order)
		d.Orderquota = TradeNum
		d.Reservepeople =len(order)
		d.ProductNum   =len(pro)
		d.DoctorNum     =len(doc)
		d.Evaluatepeople=len(order)
		d.ConsultNum   =len(order)
		if err:=WriteBack(constant.Hash_Hospital,v,d);err!=nil{
			continue
		}
	}

	Forward(session,"0","ok")
}

func Changedoc(session *JsNet.StSession) {
	ca, err := cacheIO.GetGlobalDocPtr("")
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	for _, v := range ca.Full {
		go func(id string) {
			data := &common.ST_Doctor{}
			if err := WriteLock(constant.Hash_Doctor, id, data); err != nil {
				return
			}
			data.VisitNum = 0
			if err := WriteBack(constant.Hash_Doctor, id, data); err != nil {
				return
			}
		}(v)
	}
	Forward(session, "0", nil)
}

func Changehos(session *JsNet.StSession) {
	ca, err := cacheIO.GetGlobalHosPtr("")
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	for _, v := range ca.Full {
		go func(id string) {
			data := &common.ST_Hospital{}
			if err := WriteLock(constant.Hash_Hospital, id, data); err != nil {
				return
			}
			data.VisitNum = 0
			if err := WriteBack(constant.Hash_Hospital, id, data); err != nil {
				return
			}
		}(v)
	}
	Forward(session, "0", nil)
}

func GetGlobalHospotal(session *JsNet.StSession) {
	Forward(session, "0", gl_hospital())
}



func RefreshHospitalSimpleInfo(session *JsNet.StSession) {

	list := gl_hospital()

	data := common.QueryMoreHosInfo(list)
	for _, v := range data {
		//ErrorLog("HosID:%s,Name:%s\n", v.HosID,v.HosName)
		common.AddHosNameList(v.HosID, v.HosName)
	}
	Forward(session, "0", data)
}

func gl_hospital() []string {

	return common.GlobalHospitalList()
}

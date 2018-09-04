package main

import (
	"JsLib/JsDispatcher"
	// . "JsLib/JsLogger"
	"JsLib/JsNet"
	"cache/cacheLib"
	"common"
	"constant"
	"strconv"
	. "util"
)

func hospital_init() {

	JsDispatcher.Http("/gethospitals", GetDedicateHospital) //所有医院列表
	JsDispatcher.Http("/hosp", GenerateHos)
	JsDispatcher.Http("/reviewnewhospital", ReviewNewHospital)           //审核新创建的医院
	JsDispatcher.Http("/reviewmodifiedhospital", ReviewModifiedHospital) //审核修改后的的医院
	JsDispatcher.Http("/forcehospitaloffline", ForceHospitalOffLine)     //医院下线

	// JsDispatcher.Http("/RevieweNewHospital", RevieweNewHospital) //新建医院审核
	// JsDispatcher.Http("/QueryHospitalInfo", QueryHospitalInfo)   //查询医院详情
}

//Get dedicate product
func GetDedicateHospital(session *JsNet.StSession) {
	common.GetDedicateHospital(session)
}

func initialHospital(t *common.ST_Hospital, i int) {
	info := common.ST_Opreat{
		OpreatPart:   "OpreatPart",
		OpreatAction: "OpreatAction",
		OpreatStatus: "OpreatStatus",
		OpreatReason: "OpreatReason",
		OpreatJobNum: "OpreatJobNum",
		OpreatName:   "OpreatName",
		OpreatTime:   CurTime(),
		ApplyJobNum:  "ApplyJobNum",
		ApplyName:    "ApplyName",
		ApplyCell:    "ApplyCell",
	}
	t.OpreatInfo = append(t.OpreatInfo, info)
	t.HosID = "Hos" + strconv.Itoa(i)
	t.HosName = "HosName" + strconv.Itoa(i)
	t.Current.OpreatName = "tian"
	t.Current.OpreatTime = "2010-12-01"
	t.Current.OpreatReason = "Wrong Product"
}

//Get all the new product (Test)--tianfeng
func GenerateHos(session *JsNet.StSession) {
	//Info("**************Enter GetProList\n")

	globalCashProduct := cacheLib.ST_ONUMCache{}

	for i := 1; i < 20; i++ {
		t := &common.ST_Hospital{}
		initialHospital(t, i)
		t.Current.OpreatStatus = constant.OperatingStatus_new
		globalCashProduct.New = append(globalCashProduct.New, t.HosID)

		DirectWrite(constant.Hash_Hospital, t.HosID, t)

	}

	for i := 20; i < 40; i++ {
		t := &common.ST_Hospital{}
		initialHospital(t, i)
		t.Current.OpreatStatus = constant.OperatingStatus_modify

		globalCashProduct.Modify = append(globalCashProduct.Modify, t.HosID)
		DirectWrite(constant.Hash_Hospital, t.HosID, t)
	}

	for i := 40; i < 60; i++ {
		t := &common.ST_Hospital{}
		initialHospital(t, i)

		t.Current.OpreatStatus = constant.OperatingStatus_online
		globalCashProduct.OnLine = append(globalCashProduct.OnLine, t.HosID)
		DirectWrite(constant.Hash_Hospital, t.HosID, t)
	}

	for i := 60; i < 80; i++ {
		t := &common.ST_Hospital{}
		initialHospital(t, i)

		t.Current.OpreatStatus = constant.OperatingStatus_Offline_self

		globalCashProduct.OffLine = append(globalCashProduct.OffLine, t.HosID)
		DirectWrite(constant.Hash_Hospital, t.HosID, t)
	}

	for i := 80; i < 100; i++ {
		t := &common.ST_Hospital{}
		initialHospital(t, i)

		t.Current.OpreatStatus = constant.OperatingStatus_Reviewer_NotPass
		globalCashProduct.UnPass = append(globalCashProduct.UnPass, t.HosID)
		DirectWrite(constant.Hash_Hospital, t.HosID, t)
	}

	DirectWrite(constant.Hash_Global, constant.Hash_HospitalCache, globalCashProduct)

	ForwardEx(session, "0", nil, "")
}

func ReviewNewHospital(session *JsNet.StSession) {
	common.RevieweNewHospital(session)
}

func ReviewModifiedHospital(session *JsNet.StSession) {
	common.RevieweModifyHospital(session)
}

func ForceHospitalOffLine(session *JsNet.StSession) {
	common.HosOfflineOnForce(session)
}

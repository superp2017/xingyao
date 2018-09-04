package common

import (
	"JsLib/JsNet"

	"constant"

	. "util"
)

func GetDedicateProduct(session *JsNet.StSession) {
	lsProductSend := []*ST_Product{}

	//Get the request info
	st := &ST_RequestPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	listID, err := GetDedicateListID(constant.Hash_ProductCache, st, constant.ItemAccountPerPage_Product)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}
	lsProd := QueryMoreProducts(listID)
	lsProductSend = append(lsProductSend, lsProd...)
	ForwardEx(session, "0", lsProductSend, "sucess")

}

func GetDedicateHospital(session *JsNet.StSession) {

	lsHospitalSend := []*ST_Hospital{}

	//Get the request info
	st := &ST_RequestPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}
	listID, err := GetDedicateListID(constant.Hash_HospitalCache, st, constant.ItemAccountPerPage_Product)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}
	lsProd := QueryMoreHosInfo(listID)
	lsHospitalSend = append(lsHospitalSend, lsProd...)
	ForwardEx(session, "0", lsHospitalSend, "sucess")
}

func GetDedicateOrder(session *JsNet.StSession) {

	//Get the request info
	st := &ST_RequestPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}
	listID, err := GetDedicateListID(constant.KEY_OrderStatusList, st, constant.ItemAccountPerPage_Product)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}
	lsProd := QueryMoreOrders(listID)
	ForwardEx(session, "0", lsProd, "sucess")
}

func GetDedicateDoctor(session *JsNet.StSession) {
	lsDoctorSend := []*ST_Doctor{}
	//Get the request info
	st := &ST_RequestPar{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	listID, err := GetDedicateListID(constant.Hash_DoctorCache, st, constant.ItemAccountPerPage_Product)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}
	lsProd := QueryMoreDocInfo(listID)
	lsDoctorSend = append(lsDoctorSend, lsProd...)
	ForwardEx(session, "0", lsDoctorSend, "sucess")
}

package main

import (
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"cache/cacheIO"
	// . "cache/cacheLib"
	"common"
	"constant"
	"math/rand"
	. "util"
)

func init_onum() {
	JsDispatcher.Http("/GetGlobalHospital", GetGlobalHospital)
	JsDispatcher.Http("/GetGlobalProduct", GetGlobalProduct)
	JsDispatcher.Http("/GetGlobalDoctor", GetGlobalDoctor)
	JsDispatcher.Http("/GetHoslDoctor", GetHoslDoctor)
	JsDispatcher.Http("/GetHosProduct", GetHosProduct)
	JsDispatcher.Http("/GetCityHos", GetCityHos)
	JsDispatcher.Http("/GetCityItemHos", GetCityItemHos)

	JsDispatcher.Http("/GetCityDoc", GetCityDoc)

	JsDispatcher.Http("/GetCityItemDoc", GetCityItemDoc)

	JsDispatcher.Http("/GetFirstSecondItemMap", GetFirstSecondItemMap)

	JsDispatcher.Http("/GetFirstSecondItemPro", GetFirstSecondItemPro)

	JsDispatcher.Http("/GetGlobalCityPro", GetGlobalCityPro)

	JsDispatcher.Http("/GetItemProduct", GetItemProduct)

	JsDispatcher.Http("/GetCityPriceProduct", GetCityPriceProduct)

	JsDispatcher.Http("/GetPriceProduct", GetPriceProduct)

	JsDispatcher.Http("/GetGlobalPrice", GetGlobalPrice)

	JsDispatcher.Http("/RebuildDocPro", RebuildDocPro)

	JsDispatcher.Http("/ChangeAllProSaleNums", ChangeAllProSaleNums)

	//

}

///获取全局的医院
func GetGlobalHospital(session *JsNet.StSession) {
	type st_Get struct {
		SortType string
	}
	st := &st_Get{}
	session.GetPara(st)

	hos, err := cacheIO.GetGlobalHosPtr(st.SortType)
	if err != nil {
		ForwardEx(session, "1", nil, "GetGlobalHospital failed\n")
		return
	}
	Info("global hos =%v", hos)
	Forward(session, "0", hos)
}

///获取全局的医生
func GetGlobalDoctor(session *JsNet.StSession) {
	type st_Get struct {
		SortType string
	}
	st := &st_Get{}
	session.GetPara(st)

	doc, err := cacheIO.GetGlobalDocPtr(st.SortType)
	if err != nil {
		ForwardEx(session, "1", nil, "GetGlobalDoctor failed\n")
		return
	}
	Info("global doc =%v", doc)
	Forward(session, "0", doc)
}

///获取全局的产品
func GetGlobalProduct(session *JsNet.StSession) {
	type st_Get struct {
		SortType string
	}
	st := &st_Get{}
	session.GetPara(st)

	pro, err := cacheIO.GetGlobalProPtr(st.SortType)
	if err != nil {
		ForwardEx(session, "1", nil, "GetGlobalProduct failed\n")
		return
	}
	Info("global pro =%v", pro)
	Forward(session, "0", pro)
}

///获取全局的医生
func GetHoslDoctor(session *JsNet.StSession) {
	type st_get struct {
		HosID string
	}
	st := &st_get{}
	session.GetPara(st)

	doc, err := cacheIO.GetHosDocPtr(st.HosID)
	if err != nil {
		ForwardEx(session, "1", nil, "GetGlobalDoctor failed\n")
		return
	}
	Info("global doc =%v", doc)
	Forward(session, "0", doc)
}

///获取全局的产品
func GetHosProduct(session *JsNet.StSession) {

	type st_get struct {
		HosID string
	}
	st := &st_get{}
	session.GetPara(st)

	pro, err := cacheIO.GetHosProPtr(st.HosID)
	if err != nil {
		ForwardEx(session, "1", nil, "GetGlobalProduct failed\n")
		return
	}
	Info("global pro =%v", pro)
	Forward(session, "0", pro)
}

func GetCityHos(session *JsNet.StSession) {
	type st_get struct {
		City string
	}
	st := &st_get{}
	session.GetPara(st)

	data, err := cacheIO.GetCityHos(st.City)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func GetCityItemHos(session *JsNet.StSession) {
	type st_get struct {
		City      string
		FirstItem string
	}
	st := &st_get{}
	session.GetPara(st)

	data, err := cacheIO.GetCityItemHos(st.City, st.FirstItem)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func GetCityDoc(session *JsNet.StSession) {
	type st_get struct {
		City string
	}
	st := &st_get{}
	session.GetPara(st)

	data, err := cacheIO.GetCityDoc(st.City)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func GetCityItemDoc(session *JsNet.StSession) {
	type st_get struct {
		City      string
		FirstItem string
	}
	st := &st_get{}
	session.GetPara(st)

	data, err := cacheIO.GetCityItemDoc(st.City, st.FirstItem)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func GetFirstSecondItemMap(session *JsNet.StSession) {
	data, err := cacheIO.GetFirstSecondItemMap()
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func GetFirstSecondItemPro(session *JsNet.StSession) {
	type st_Get struct {
		FirstItem string
	}
	st := &st_Get{}
	session.GetPara(st)
	data, err := cacheIO.GetFirstSecondItemPro(st.FirstItem)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func GetGlobalCityPro(session *JsNet.StSession) {
	type st_get struct {
		City string
	}
	st := &st_get{}
	session.GetPara(st)
	data, err := cacheIO.GetGlobalCityPro(st.City)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "1", data)
}

func GetItemProduct(session *JsNet.StSession) {
	type st_Get struct {
		FirstItem  string
		SecondItem string
	}
	st := &st_Get{}
	session.GetPara(st)
	data, err := cacheIO.GetItemProduct(st.FirstItem, st.SecondItem)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "1", data)

}

func GetCityPriceProduct(session *JsNet.StSession) {
	type st_Get struct {
		City string
	}
	st := &st_Get{}
	session.GetPara(st)
	data, err := cacheIO.GetCityPriceProduct(st.City)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "1", data)
}

func GetPriceProduct(session *JsNet.StSession) {
	data, err := cacheIO.GetPriceProduct()
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func GetGlobalPrice(session *JsNet.StSession) {
	data, err := cacheIO.GetGlobalPrice()
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func RebuildDocPro(session *JsNet.StSession) {

	hoslist := common.GlobalHospitalList()
	for _, v := range hoslist {
		pros := common.QueryMoreProducts(common.GetHosProductList(v))
		for _, pro := range pros {
			common.AddProjectToDoc(pro.Doctors, pro.ProID)
		}
	}
	Forward(session, "0", nil)
}

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

func ChangeAllProSaleNums(session *JsNet.StSession) {
	pro, err := cacheIO.GetGlobalProPtr("")
	if err != nil {
		ForwardEx(session, "1", nil, "ChangeAllProSaleNums failed\n")
		return
	}

	for _, v := range pro.Full {

		pro := &common.ST_Product{}
		Update(constant.Hash_HosProduct, v, pro, func() {
			pro.Salesvolumes = (int)(RandInt64(25, 35))
			pro.AppointmentNums = pro.Salesvolumes
		})
	}

	Forward(session, "0", nil)
}

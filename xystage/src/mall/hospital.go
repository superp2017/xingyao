package main

import (
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"cache/cacheIO"
	"common"
	"constant"
	// "strconv"
	// "errors"

	"cache/cacheLib"
	"sort"
	. "util"
)

//产品详情
type ST_RequestPar struct {
	City        string
	RequestArg  string
	RequestPage int
	SortType    string
}

type ST_HospitalNet struct {
	HospitalInfo common.ST_Hospital
	LsProduct    cacheLib.ST_ONUMCache
	LsDoctor     cacheLib.ST_ONUMCache
}

func hospital_init() {

	JsDispatcher.Http("/gethospitals", GetCityHospital)
	JsDispatcher.Http("/gethospitalsn", GetCityHospitaln)
	JsDispatcher.Http("/gethospital", GetNetHospital)
	JsDispatcher.Http("/getlshospital", GetNetLsHospital)              //所有医院列表
	JsDispatcher.Http("/gettotalpagenumhospital", GettotalNumHospital) //所有医院列表
}

func GettotalNumHospital(session *JsNet.StSession) {

	type ST_PageNumTotal struct {
		PageNum int
		City    string
	}

	st := &ST_RequestPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	lsHospitalID := GetHospitalList(st.SortType, st.City, st.RequestArg)
	totalNum := &ST_PageNumTotal{}
	totalNum.PageNum = getCeilNum(len(lsHospitalID), constant.ItemAccountPerPage_Hospital)
	// totalNum.HmTotalPageNum = make(map[string]int)
	// totalNum.HmTotalPageNum[st.RequireType] = getCeilNum(len(lsHospitalID), constant.ItemAccountPerPage_Hospital)
	totalNum.City = st.City
	ForwardEx(session, "0", totalNum, st.RequestArg)
}

func getCeilNum(a int, b int) int {
	c := 0.1
	c = float64(a) / float64(b)
	d := int(c)
	e := float64(d)

	if c > e {
		return d + 1
	} else {
		return d
	}
}

func GetCityHospital(session *JsNet.StSession) {

	lsHospitalSend := []*common.ST_Hospital{}

	//Get the request info
	st := &ST_RequestPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	listID := GetDedicateListIDHospital(st, constant.ItemAccountPerPage_Hospital)

	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}
	lsHospitalSend = common.QueryMoreHosInfo(listID)
	// superhospital.HmHospital[st.RequestArg] = lsHospitalSend
	ForwardEx(session, "0", lsHospitalSend, st.RequestArg)
}

func GetCityHospitaln(session *JsNet.StSession) {

	// type SuperHospital struct {
	// 	HmHospital map[string]([]*common.ST_Hospital)
	// }

	// superhospital := SuperHospital{}
	// superhospital.HmHospital = make(map[string]([]*common.ST_Hospital))

	// lsHospitalSend := []*common.ST_Hospital{}

	//Get the request info
	st := &ST_RequestPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	listID := GetDedicateListIDHospital(st, constant.ItemAccountPerPage_Hospital)

	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}

	lsNetHospital := []*ST_HospitalNet{}

	for _, v := range listID {
		hosST_HospitalNet, err := getNetHospitalL(v)
		if err == nil {
			lsNetHospital = append(lsNetHospital, hosST_HospitalNet)
		}

	}

	// superhospital.HmHospital[st.RequestArg] = lsHospitalSend
	ForwardEx(session, "0", lsNetHospital, st.RequestArg)
}

func GetDedicateListIDHospital(st *ST_RequestPar, itemPerPage int) []string {

	Info("Request Par=%v\n", st)

	listPageID := []string{}
	listID := GetHospitalList(st.SortType, st.City, st.RequestArg)
	listStartDex := (st.RequestPage - 1) * itemPerPage

	if listStartDex+itemPerPage > len(listID) {
		listStartDex = len(listID) - itemPerPage

	}

	if len(listID)-listStartDex <= itemPerPage {
		listPageID = listID
		return listPageID
	} else {
		listPageID = listID[:listStartDex+itemPerPage]
	}
	return listPageID
}

func GetHospitalList(SortType, cityName string, bodypart string) []string {
	lsID := []string{}

	if cityName == "" {
		Error("City name is null\n")
		return lsID
	}

	if cityName == "全部" || cityName == "全国" {
		cityName = constant.All
	}

	var isE error = nil

	if bodypart == "" || bodypart == " " {
		lsID, isE = cacheIO.GetCityHos(cityName)
	} else {
		lsID, isE = cacheIO.GetCityItemHos(cityName, bodypart)
	}

	if isE != nil {
		Error("GetHospitalList err:\n", isE.Error())
		return lsID
	}
	lsIDBackHos := []string{}
	lsHos, err := cacheIO.GetGlobalHosPtr(SortType)

	if err != nil {
		Error("cacheIO.GetGlobalHosPtr() err:\n", err.Error())
		return lsID
	}

	lsHosOnline := lsHos.OnLine

	//lsID, _ := common.GetCityItemDoc(cityName, bodypart)
	for _, v := range lsID {
		for _, t := range lsHosOnline {
			if v == t {
				lsIDBackHos = append(lsIDBackHos, t)
			}
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(lsIDBackHos)))

	return lsIDBackHos
}

func GetNetHospital(session *JsNet.StSession) {
	//get the request
	//hosST_HospitalNet := &ST_HospitalNet{}
	type st_query struct {
		HosID string //医院id
	}

	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" {
		ForwardEx(session, "1", nil, "医院id为空,QueryHospitalInfo()，查询失败!")
		return
	}
	hosST_HospitalNet, err := getNetHospitalL(st.HosID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	// if hosST_HospitalNet.HospitalInfo.Current.OpreatStatus != constant.OperatingStatus_online {
	// 	Forward(session, "1", nil)
	// 	return

	// }
	Forward(session, "0", hosST_HospitalNet)
}

func GetNetLsHospital(session *JsNet.StSession) {
	type st_query struct {
		HosID []string //医院id
	}

	lsNetHospital := []*ST_HospitalNet{}

	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	for _, v := range st.HosID {
		hosST_HospitalNet, err := getNetHospitalL(v)
		if err == nil {
			lsNetHospital = append(lsNetHospital, hosST_HospitalNet)
		}

	}
	Forward(session, "0", lsNetHospital)
}

func getNetHospitalL(hospitalID string) (*ST_HospitalNet, error) {

	hosST_HospitalNet := &ST_HospitalNet{}
	//get the hospital
	hospital, err := common.GetHospitalInfo(hospitalID)
	if err != nil {
		return nil, err
	}

	// if hospital.Current.OpreatStatus != constant.OperatingStatus_online {
	// 	return nil, errors.New("The hospital status is not online")

	// }

	hosST_HospitalNet.HospitalInfo = *hospital

	//get the doctor's information
	hosdoc, err := cacheIO.GetHosDocPtr(hospitalID)
	if err != nil {

		return nil, err
	}
	hosST_HospitalNet.LsDoctor = *hosdoc

	//get the product's information

	hospro, err := cacheIO.GetHosProPtr(hospitalID)
	if err != nil {
		return nil, err
	}
	hosST_HospitalNet.LsProduct = *hospro
	return hosST_HospitalNet, nil
}

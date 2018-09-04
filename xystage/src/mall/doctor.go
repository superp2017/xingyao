package main

import (
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"cache/cacheIO"
	"common"
	"constant"
	// "errors"
	"sort"
	// "strconv"
	. "util"
)

type ST_DoctorNet struct {
	DoctorInfo common.ST_Doctor
}

func doctor_init() {
	JsDispatcher.Http("/getdoctors", GetCitydoctor)
	JsDispatcher.Http("/getdoctorsn", GetCitydoctorn)
	JsDispatcher.Http("/getdoctor", GetNetDoctor)                  //所有医院列表
	JsDispatcher.Http("/gettotalpagenumdoctor", GettotalNumDoctor) //所有医院列表
	// JsDispatcher.Http("/getexpertdoctors", GetExpertDoctor)        ////获取大牌名医列表

}

func GettotalNumDoctor(session *JsNet.StSession) {

	type ST_PageNumTotal struct {
		PageNum    int
		RequestArg string
		City       string
	}

	st := &ST_RequestPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	lsdoctorID := GetdoctorList(st.SortType, st.City, st.RequestArg)
	totalNum := &ST_PageNumTotal{}
	totalNum.RequestArg = st.RequestArg
	totalNum.PageNum = getCeilNum(len(lsdoctorID), constant.ItemAccountPerPage_Doctor)
	// totalNum.HmTotalPageNum = make(map[string]int)
	// totalNum.HmTotalPageNum[st.RequireType] = getCeilNum(len(lsdoctorID), constant.ItemAccountPerPage_doctor)
	totalNum.City = st.City
	ForwardEx(session, "0", totalNum, st.RequestArg)
}

func GetCitydoctor(session *JsNet.StSession) {

	lsdoctorSend := []*common.ST_Doctor{}

	//Get the request info
	st := &ST_RequestPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}
	listID := GetDedicateListIDDoctor(st, constant.ItemAccountPerPage_Doctor)

	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}
	lsdoctorSend = common.QueryMoreDocInfo(listID)
	// superdoctor.Hmdoctor[st.doctorType] = lsdoctorSend
	ForwardEx(session, "0", lsdoctorSend, st.RequestArg)
}

func GetCitydoctorn(session *JsNet.StSession) {

	//Get the request info
	st := &ST_RequestPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	listID := GetDedicateListIDDoctor(st, constant.ItemAccountPerPage_Doctor)

	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}

	lsNetDoctor := []*ST_DoctorNet{}

	for _, v := range listID {
		hosST_DoctorNet, err := getNetDoctorL(v)
		if err == nil {
			lsNetDoctor = append(lsNetDoctor, hosST_DoctorNet)
		}

	}
	// superdoctor.Hmdoctor[st.doctorType] = lsdoctorSend
	ForwardEx(session, "0", lsNetDoctor, st.RequestArg)
}

func GetDedicateListIDDoctor(st *ST_RequestPar, itemPerPage int) []string {
	listPageID := []string{}
	listID := GetdoctorList(st.SortType, st.City, st.RequestArg)
	listStartDex := (st.RequestPage - 1) * itemPerPage
	if len(listID)-listStartDex <= itemPerPage {
		listPageID = listID
		return listPageID
	} else {
		listPageID = listID[:listStartDex+itemPerPage]
	}
	//OutPutDocInfo("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@Get result\n")
	return listPageID
}

func GetdoctorList(SortType, cityName string, bodypart string) []string {
	lsID := []string{}

	var isE error = nil
	if bodypart == "" || bodypart == " " {
		lsID, isE = cacheIO.GetCityDoc(cityName)
	} else {
		lsID, isE = cacheIO.GetCityItemDoc(cityName, bodypart)
	}
	if isE != nil {
		Error("GetdoctorList , err:", isE.Error())
		return lsID
	}
	//add new one

	lsIDBackDoc := []string{}
	lsDoc, err := cacheIO.GetGlobalDocPtr(SortType)
	if err != nil {
		Error("GetGlobalDocPtr , err:", err.Error())
		return lsID
	}

	lsDocOnline := lsDoc.OnLine

	//lsID, _ := common.GetCityItemDoc(cityName, bodypart)
	for _, v := range lsID {
		for _, t := range lsDocOnline {
			if v == t {
				lsIDBackDoc = append(lsIDBackDoc, t)
			}
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(lsIDBackDoc)))
	return lsIDBackDoc
}

func GetNetDoctor(session *JsNet.StSession) {
	//get the request
	//hosST_DoctorNet := &ST_DoctorNet{}
	type st_query struct {
		DocID string //医院id
	}

	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.DocID == "" {
		ForwardEx(session, "1", nil, "医院id为空,QueryDoctorInfo()，查询失败!")
		return
	}
	docST_DoctorNet, err := getNetDoctorL(st.DocID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", docST_DoctorNet)
}

func getNetDoctorL(DoctorID string) (*ST_DoctorNet, error) {

	docST_DoctorNet := &ST_DoctorNet{}
	//get the Doctor
	Doctor, err := common.QueryDoctor(DoctorID)
	if err != nil {
		return nil, err
	}

	docST_DoctorNet.DoctorInfo = *Doctor
	return docST_DoctorNet, nil
}

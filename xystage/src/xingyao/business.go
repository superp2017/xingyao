package main

import (
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"common"
	"constant"

	. "util"
)

//产品详情
type ST_RequestBusinessPar struct {
	UID          string
	BusinessType string
	RequestPage  int
}

type ST_BusinessNet struct {
	BusinessInfo common.ST_User
}

func Business_init() {

	JsDispatcher.Http("/getbusinesss", GetBusinesss)
	JsDispatcher.Http("/reviewnewbusiness", common.RevieweAgent)

	JsDispatcher.Http("/getbusinesssn", GetBusinesssN)
	JsDispatcher.Http("/getbusiness", GetNetBusiness)
	JsDispatcher.Http("/getlsbusiness", GetNetLsBusiness)              //所有医院列表
	JsDispatcher.Http("/gettotalpagenumBusiness", GettotalNumBusiness) //所有医院列表                    //所有医院列表

}

func GettotalNumBusiness(session *JsNet.StSession) {

	type ST_PageNumTotal struct {
		PageNum      int
		BusinessType string
		UID          string
	}

	st := &ST_RequestBusinessPar{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if st == nil {
		ForwardEx(session, "1", nil, "GettotalNumBusiness get param is nil")
		return
	}

	lsBusinessID := GetBusinessList(st)
	totalNum := &ST_PageNumTotal{}
	totalNum.BusinessType = st.BusinessType
	totalNum.PageNum = getCeilNum(len(lsBusinessID), constant.ItemAccountPerPage_Business)
	totalNum.UID = st.UID
	ForwardEx(session, "0", totalNum, st.BusinessType)
}

func GetBusinesss(session *JsNet.StSession) {

	lsBusinessSend := []*common.ST_User{}
	//Get the request info
	st := &ST_RequestBusinessPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	listID := GetDedicateListIDBusiness(st, constant.ItemAccountPerPage_Business)
	Info("List ID=%v\n", listID)
	// if err != nil {
	// 	ForwardEx(session, "1", nil, err.Error())
	// 	return
	// }

	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}
	lsBusinessSend = common.GetMoreUserInfo(listID)

	// superBusiness.HmBusiness[st.BusinessType] = lsBusinessSend
	ForwardEx(session, "0", lsBusinessSend, st.BusinessType)
}

func GetBusinesssN(session *JsNet.StSession) {

	lsBusinessSend := []*common.ST_User{}
	//Get the request info
	st := &ST_RequestBusinessPar{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if st == nil {
		ForwardEx(session, "1", nil, "GetBusinesssN param is nil")
		return
	}

	listID := GetDedicateListIDBusiness(st, constant.ItemAccountPerPage_Business)
	Info("List ID=%v\n", listID)
	// if err != nil {
	// 	ForwardEx(session, "1", nil, err.Error())
	// 	return
	// }

	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}
	lsBusinessSend = common.GetMoreUserInfo(listID)

	lsBusinessSuper := []*ST_BusinessNet{}
	for _, v := range lsBusinessSend {
		superBusiness := &ST_BusinessNet{}
		superBusiness.BusinessInfo = *v
		lsBusinessSuper = append(lsBusinessSuper, superBusiness)
	}

	// superBusiness.HmBusiness[st.BusinessType] = lsBusinessSend
	ForwardEx(session, "0", lsBusinessSend, st.BusinessType)
}

func GetDedicateListIDBusiness(st *ST_RequestBusinessPar, itemPerPage int) []string {

	if st == nil {
		return []string{}
	}

	if itemPerPage < 0 {
		return []string{}
	}

	listPageID := []string{}
	listID := GetBusinessList(st)
	if len(listID) <= 0 {
		return []string{}
	}

	listStartDex := (st.RequestPage - 1) * itemPerPage
	if len(listID)-listStartDex <= itemPerPage {
		listPageID = listID[listStartDex:]
		Info("List ID1=%v\n", listPageID)
		return listPageID
	} else {
		Info("List ID2=%v\n", listPageID)
		listPageID = listID[listStartDex : listStartDex+itemPerPage]
	}
	Info("List ID=%v\n", listPageID)
	return listPageID
}

// type ST_AgentStatusID struct {
// 	Apply   []*string //申请中
// 	WaitPay []*string //待支付
// 	NoPass  []*string //审核不通过
// 	Online  []*string //在线的
// 	OfflineOnForce []*string //解约的
// 	OfflineSelf []*string //解约的
// }

// const (
// 	Agent_Apply         = "Agent_Apply"         //申请成为代理
// 	Agent_ReApply       = "Agent_ReApply"       //重新申请成为代理
// 	Agent_PassReviewe   = "Agent_PassReviewe"   //资料审核通过
// 	Agent_NoPassReviewe = "Agent_NoPassReviewe" //资料审核不通过
// 	Agent_Online        = "Agent_Online"        //代理上线
// 	Agent_Offline_force = "Agent_Offline_force" //平台强制下线
// 	Agent_Offline_self  = "Agent_Offline_self"  //代理自己下线
// )

// for _, v := range UserList {
// 	if v.Agent.Current.OpreatStatus == constant.Agent_Apply ||
// 		v.Agent.Current.OpreatStatus == constant.Agent_ReApply {
// 		info.Apply = append(info.Apply, v)
// 	} else if v.Agent.Current.OpreatStatus == constant.Agent_PassReviewe {
// 		info.WaitPay = append(info.WaitPay, v)
// 	} else if v.Agent.Current.OpreatStatus == constant.Agent_NoPassReviewe {
// 		info.NoPass = append(info.NoPass, v)
// 	} else if v.Agent.Current.OpreatStatus == constant.Agent_Offline_force ||
// 		v.Agent.Current.OpreatStatus == constant.Agent_Offline_self {
// 		info.Offline = append(info.Offline, v)
// 	} else {
// 		Error("v.Agent.Current.OpreatStatus=%s\n", v.Agent.Current.OpreatStatus)
// 		continue
// 	}
// }

func GetBusinessList(st *ST_RequestBusinessPar) []string {
	if st == nil {
		return []string{}
	}

	Info("The Request Info=%v\n", st)
	lsID := []string{}
	res := common.GetGlobalAgentID()

	if res == nil {

		ErrorLog("GetBusinessList   common.GetGlobalAgentID is nil ")
		return []string{}
	}

	Info("-----------------------------Status =%v\n", res)
	//"待预约"
	// Status_Business_PenddingPay      = "Business_PenddingPay"      //"待支付"//用户、代理、系统
	if st.BusinessType == constant.Agent_Apply ||
		st.BusinessType == constant.Agent_ReApply {
		lsID = res.Apply

	} else if st.BusinessType == constant.Agent_PassReviewe {

		lsID = res.WaitPay

	} else if st.BusinessType == constant.Agent_NoPassReviewe {

		lsID = res.NoPass

	} else if st.BusinessType == constant.Agent_Offline_force {

		lsID = res.OfflineOnForce

	} else if st.BusinessType == constant.Agent_Offline_self {

		lsID = res.OfflineSelf

	} else if st.BusinessType == constant.Agent_Online {

		lsID = res.Online

	} else {

	}

	Info("+++++++++++++++++++++++++++++++++Serve FB ID List=%v\n", lsID)
	return lsID

}

func GetNetBusiness(session *JsNet.StSession) {
	//get the request
	//hosST_HospitalNet := &ST_HospitalNet{}
	type st_query struct {
		BusinessID string //医院id
	}
	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.BusinessID == "" {
		ForwardEx(session, "1", nil, "Business id为空,QueryHospitalInfo()，查询失败!")
		return
	}
	BusinessST_BusinessNet, err := getNetBusinessL(st.BusinessID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", BusinessST_BusinessNet)
}

func GetNetLsBusiness(session *JsNet.StSession) {
	type st_queryBusiness struct {
		BusinessID []string
	}

	lsNetBusiness := []*ST_BusinessNet{}

	st := &st_queryBusiness{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	for _, v := range st.BusinessID {
		BusinessST_BusinessNet, err := getNetBusinessL(v)
		if err == nil {
			lsNetBusiness = append(lsNetBusiness, BusinessST_BusinessNet)
		}
	}
	Forward(session, "0", lsNetBusiness)
}

func getNetBusinessL(BusinessID string) (*ST_BusinessNet, error) {
	BusinessNet := &ST_BusinessNet{}
	//get the common
	BusinessInfo, err := common.GetUserInfo(BusinessID)

	if err != nil {
		return nil, err

	}
	BusinessNet.BusinessInfo = *BusinessInfo
	return BusinessNet, nil
}

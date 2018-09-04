package main

import (
	"JsLib/JsDispatcher"
	"JsLib/JsNet"
	"cache/cacheIO"
	"common"
	"constant"
	. "util"
)

func init_router() {

	JsDispatcher.Http("/getuserfromopenid", common.GetUserFormOpenID) /////通过openID，查询用户信息,

	JsDispatcher.Http("/getuserfromunionid", common.GetUserFormUnionID) /////通过unionID，查询用户信息,

	JsDispatcher.Http("/getuserinfo", GetUserInfo) //通过UID查询用户信息

	JsDispatcher.Http("/modifyuser", ModifyUser) ///修改用户

	JsDispatcher.Http("/bindusercell", BindUserCell) //绑定手机

	JsDispatcher.Http("/tochangecell", TochangeCell) //换绑手机

	JsDispatcher.Http("/gethomeproducts", GetHomeProducts) ////获取首页产品列表

	JsDispatcher.Http("/getexpertdoctors", GetExpertDoctor) ////获取大牌名医列表

	JsDispatcher.Http("/getcitylist", GetGlobalCityList) ////获取城市

	JsDispatcher.Http("/getcityhospital", GetGlobalCityHosList) ////获取分城市的医院列表

	JsDispatcher.Http("/getcityitemhospital", GetCityItemHosList) ////获取分城市分项目的医院列表

	JsDispatcher.Http("/getcitydoctor", GetGlobalCityDocList) ////获取分城市的医生列表

	JsDispatcher.Http("/gethomeexpertdoctor", GetHomeExpertDoctor) ////获取首页大牌名医

	JsDispatcher.Http("/getcityitemdoctor", GetCityItemDocList) ////获取分城市分项目的城市列表

	JsDispatcher.Http("/followhospitial", FollowHospitial) ///关注医院

	JsDispatcher.Http("/followdocto", FollowDocto) ///关注医生

	JsDispatcher.Http("/collectionporduct", CollectionPorduct) ///收藏产品

	JsDispatcher.Http("/followhosview", FollowHosView) ///关注医院查看

	JsDispatcher.Http("/followdocView", FollowDocView) ////关注医生查看

	JsDispatcher.Http("/collectionproView", CollectionProView) ////收藏产品查看

	JsDispatcher.Http("/appendproductcomment", AppendProductComment) //新的产品评论

	JsDispatcher.Http("/changeprovisitnum", common.ChangeProVisitNum) //更新产品访问量
	JsDispatcher.Http("/changedocvisitnum", common.ChangeDocVisitNum) //更新医生访问量
	JsDispatcher.Http("/changehosvisitnum", common.ChangeHosVisitNum) //更新医院访问量


	JsDispatcher.Http("/changeproconsultnum", common.ChangeProConsultNum) //更新产品咨询量
	JsDispatcher.Http("/changedocconsultnum", common.ChangeDocConsultNum) //更新医生咨询量
	JsDispatcher.Http("/changehosConsultnum", common.ChangeHosConsultNum) //更新医生咨询量



	JsDispatcher.Http("/usefav", UseFav) //用户的喜好

	JsDispatcher.Http("/querywithdrawinfo", common.QueryWithDrawInfo) ///获取提现信息

	JsDispatcher.Http("/querymorewithdrawrecord", common.QueryMoreRecord) ///获取提现信息

	////////////////////////多个查询接口//////////////////////////////////////
	JsDispatcher.Http("/querymorehospital", common.QueryMoreHospital)   //查询多个医院信息
	JsDispatcher.Http("/querymoredoctor", common.QueryMoreDoctors)      //查询多个医生信息
	JsDispatcher.Http("/querymoreproduct", common.QueryMoreProductInfo) //查询多个产品信息
	JsDispatcher.Http("/querymoreorder", common.QueryMoreOrderInfo)     //查询多个订单信息

	/////////////////////订单接口/////////////////////////////////
	JsDispatcher.Http("/ordercancle", OrderCancle)           //订单取消
	JsDispatcher.Http("/queryorderinfo", QueryOrderInfo)     //查看订单详情
	JsDispatcher.Http("/getuserallorders", GetUserAllOrders) //获取某个用户的所有订单
	////////////////////////代理人接口///////////////////////////////
	JsDispatcher.Http("/applyagent", common.ApplyAgent)                   //申请为代理人
	JsDispatcher.Http("/reapplyagent", common.ReApplyAgent)               //修改后重新为代理人
	JsDispatcher.Http("/modifyagentbaseinfo", common.ModifyAgentBaseInfo) //修改代理的基本信息
	JsDispatcher.Http("/userbindagent", common.UserBindAgent)             ///用户绑定一个上级代理
	JsDispatcher.Http("/trytouseagent", common.TryToUseAgent)             ///试用小B

	JsDispatcher.Http("/getsearchcontent", GetSearcheContent) //获取搜索Hash

	JsDispatcher.Http("/getagentinvitestatistic", common.GetAgentInviteStatistic) //获取小B邀请统计
	JsDispatcher.Http("/getuserincrease", common.GetUserIncrease)                 //获取用户每天的增加量
	///////////////////////////
	// JsDispatcher.Http("/noticelist", NoticeList)
}

///获取首页的产品列表
func GetHomeProducts(session *JsNet.StSession) {
	type st_get struct {
		City     string //城市
		BodyPart string //身体部位
		Item     string //二级菜单
		MinPrice int    //最小价格
		MaxPrice int    //最大价格
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.City == "" && st.BodyPart == "" && st.Item == "" && st.MinPrice <= 0 && st.MaxPrice <= 0 {
		ForwardEx(session, "1", nil, "GetHomeProducts failed , param failed,City=%s,BodyPart=%s,Item=%s,MinPrice=%d,MaxPrice=%d\n",
			st.City, st.BodyPart, st.Item, st.MinPrice, st.MaxPrice)
		return
	}
	city := st.City
	if city == "全部" || city == "全国" {
		city = constant.All
	}
	data, err := getProducts(city, st.BodyPart, st.Item, st.MinPrice, st.MaxPrice)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", common.QueryMoreProducts(data))
}

//获取所有城市列表
func GetGlobalCityList(session *JsNet.StSession) {
	common.GetCityMap(session)
}

//获取所有的分城市的医院列表
func GetGlobalCityHosList(session *JsNet.StSession) {
	type st_get struct {
		City string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.City == "" {
		ForwardEx(session, "1", nil, "GetGlobalCityHosList param City is empty\n")
		return
	}
	city := st.City
	if city == "全部" || city == "全国" {
		city = constant.All
	}

	data, err := cacheIO.GetCityHos(city)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Forward(session, "0", common.QueryMoreHosInfo(data))
}

////获取分城市分项目的医院列表
func GetCityItemHosList(session *JsNet.StSession) {
	type st_get struct {
		City      string
		FirstItem string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.City == "" || st.FirstItem == "" {
		ForwardEx(session, "1", nil, "GetCityItemHosList , param City=%s,FirstItem=%s\n", st.City, st.FirstItem)
		return
	}
	city := st.City
	if city == "全部" || city == "全国" {
		city = constant.All
	}

	data, err := cacheIO.GetCityItemHos(city, st.FirstItem)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", common.QueryMoreHosInfo(data))
}

//获取所有的分城市的医生列表
func GetGlobalCityDocList(session *JsNet.StSession) {
	type st_get struct {
		City string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.City == "" {
		ForwardEx(session, "1", nil, "GetGlobalCityDocList param City is empty\n")
		return
	}
	city := st.City
	if city == "全部" || city == "全国" {
		city = constant.All
	}
	data, err := cacheIO.GetCityDoc(city)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", common.QueryMoreDocInfo(data))
}

//获取所有的分城市的医生列表
func GetCityItemDocList(session *JsNet.StSession) {
	type st_get struct {
		City      string
		FirstItem string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.City == "" {
		ForwardEx(session, "1", nil, "GetGlobalCityDocList param City is empty\n")
		return
	}
	city := st.City
	if city == "全部" || city == "全国" {
		city = constant.All
	}
	data, err := cacheIO.GetCityItemDoc(city, st.FirstItem)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", common.QueryMoreDocInfo(data))
}

//获取分城市的专家
func GetExpertDoctor(session *JsNet.StSession) {
	common.GetExpertDoctor(session)
}

///获取 首页大牌名医
func GetHomeExpertDoctor(session *JsNet.StSession) {
	common.GetHomeExpertDoctor(session)
}

///关注医院
func FollowHospitial(session *JsNet.StSession) {
	common.FollowHospitial(session)
}

///关注医生
func FollowDocto(session *JsNet.StSession) {
	common.FollowDocto(session)
}

///收藏产品
func CollectionPorduct(session *JsNet.StSession) {
	common.CollectionPorduct(session)
}

///关注医院查看
func FollowHosView(session *JsNet.StSession) {
	common.FollowHosView(session)
}

////关注医生查看
func FollowDocView(session *JsNet.StSession) {
	common.FollowDocView(session)
}

////收藏产品查看
func CollectionProView(session *JsNet.StSession) {
	common.CollectionProView(session)
}

///修改用户信息
func ModifyUser(session *JsNet.StSession) {
	common.ModifyUser(session)
}

///绑定手机
func BindUserCell(session *JsNet.StSession) {
	common.BindUserCell(session)
}

///换绑手机
func TochangeCell(session *JsNet.StSession) {
	common.TochangeCell(session)
}

///新的产品评论
func AppendProductComment(session *JsNet.StSession) {
	common.AppendProductComment(session)
}

///用户的喜好
func UseFav(session *JsNet.StSession) {
	common.UseFav(session)

}

////////////////////////////////订单接口///////////////////////////////////////////

///查询订单详情
func QueryOrderInfo(session *JsNet.StSession) {
	common.QueryOrderInfo(session)
}

//校验前取消
func OrderCancle(session *JsNet.StSession) {
	type st_new struct {
		OrderID  string
		UserName string
		UID      string
	}
	st := &st_new{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.OrderID == "" || st.UID == "" {
		ForwardEx(session, "1", nil, "param is empty,OrderID=%s,UID=%s\n", st.OrderID, st.UID)
		return
	}
	data, err := common.OrderCancle(st.OrderID, "用户取消", st.UID, st.UserName)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

///获取某个用户的所有的订单
func GetUserAllOrders(session *JsNet.StSession) {
	type st_get struct {
		UID string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" {
		ForwardEx(session, "1", nil, "GetUserAllOrder failed,UID=%s\n", st.UID)
		return
	}
	list, err := common.GetUserAllOrderList(st.UID)
	if err != nil {
		ForwardEx(session, "1", nil, "GetUserContinueOrderList failed,UID=%s,err:%s\n", st.UID, err.Error())
		return
	}
	Forward(session, "0", common.QueryMoreOrders(list))
}

func GetSearcheContent(session *JsNet.StSession) {

	data, err := cacheIO.GetSearchHash()
	if err != nil {
		Forward(session, "1", err.Error())
		return
	}
	Forward(session, "0", data.Data)
}

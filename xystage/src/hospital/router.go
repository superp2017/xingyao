package main

import (
	"JsGo/JsMobile"
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"article"
	"common"
	"constant"
	"sort"
	"time"
	. "util"
)

func init_router() {
	////////////////////////医院注册修改查询，账号/////////////
	JsDispatcher.Http("/registerhospital", common.RegisterHospital) //医生注册
	JsDispatcher.WhiteList("/registerhospital")
	JsDispatcher.Http("/ModifyHospitalInfo", common.ModifyHospitalInfo) //医院修改
	JsDispatcher.Http("/QueryHospitalInfo", common.QueryHospitalInfo)   //查询医院信息
	///////////////////////////医院登录、账号修改///////////////
	JsDispatcher.Http("/login", common.HosLogin) //医院登录
	JsDispatcher.WhiteList("/login")
	JsDispatcher.Http("/HosResetLoginCode", common.HosResetLoginCode)           ///重置医院登录账号
	JsDispatcher.Http("/ModifyAdminInfoAccount", common.ModifyAdminInfoAccount) //修改管理员信息

	////////////////////////////医生//////////////////////////////////////////
	JsDispatcher.Http("/NewDoc", common.NewDoc)                 //新建医生
	JsDispatcher.Http("/GetHosDocFullList", GetHosDocFullList)  //查询所有医生列表
	JsDispatcher.Http("/ModifyDocInfo", common.ModifyDocInfo)   //修改医生
	JsDispatcher.Http("/DocOfflineSelf", common.DocOfflineSelf) //医院将医生下线

	///////////////////////////////产品创建，查询，修改////////////////////////////
	JsDispatcher.Http("/ApplyNewProduct", common.ApplyNewProduct)       //新建一个产品
	JsDispatcher.Http("/GetHosProFullList", GetHosProFullList)          //查询所有产品列表
	JsDispatcher.Http("/ModifyProductInfo", common.ModifyProductInfo)   //修改产品信息
	JsDispatcher.Http("/OfflineProductSelf", common.OfflineProductSelf) //产品修改下线

	/////////////////////////////订单////////////////////////////////////

	JsDispatcher.Http("/verifyqueryorder", VerifyQueryOrder) //校验码获取订单号
	JsDispatcher.Http("/orderhosverify", OrderHosVerify)     //医院校验
	JsDispatcher.Http("/hosconfirmopreat", HosConfirmOpreat) //医院确认手术，并且追加金额
	JsDispatcher.Http("/hosorderglist", HosOrdergList)       //医院所有订单列表

	////////////////////其他///////////////////////////
	JsDispatcher.Http("/mobileverify", mobile_veirfy) //手机短息验证

	JsDispatcher.Http("/hospitalnotify", article.GetFrontPageShowHospital) //医院公告列表

	JsDispatcher.Http("/deldoctor", common.DelDoctor)     //删除医生
	JsDispatcher.Http("/delproduct", common.DelProduct)   //删除产品
	JsDispatcher.Http("/delhospital", common.DelHospital) //删除医院

}

//查询所有医生列表

func GetHosDocFullList(session *JsNet.StSession) {
	type st_get struct {
		HosID string
	}
	st := &st_get{}
	if err := session.GetPara(&st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" {
		ForwardEx(session, "1", nil,
			"GetHosDocFullList param is empty,HosID=%s\n", st.HosID)
		return
	}

	type all_doc struct {
		New     []*common.ST_Doctor
		OnLine  []*common.ST_Doctor
		OffLine []*common.ST_Doctor
		Modify  []*common.ST_Doctor
		UnPass  []*common.ST_Doctor
	}
	all := &all_doc{}

	data := common.QueryMoreDocInfo(common.GetHosDoctorList(st.HosID))

	for _, v := range data {
		if v.Current.OpreatStatus == constant.OperatingStatus_new {
			all.New = append(all.New, v)
		}
		if v.Current.OpreatStatus == constant.OperatingStatus_modify {
			all.Modify = append(all.Modify, v)
		}
		if v.Current.OpreatStatus == constant.OperatingStatus_online {
			all.OnLine = append(all.OnLine, v)
		}
		if v.Current.OpreatStatus == constant.OperatingStatus_Reviewer_NotPass {
			all.UnPass = append(all.UnPass, v)
		}
		if v.Current.OpreatStatus == constant.OperatingStatus_Offline_self ||
			v.Current.OpreatStatus == constant.OperatingStatus_Offline_onforce {
			all.OffLine = append(all.OffLine, v)
		}
	}

	Forward(session, "0", all)
}

//获取所有产品列表
func GetHosProFullList(session *JsNet.StSession) {
	type st_get struct {
		HosID string
	}
	st := &st_get{}
	if err := session.GetPara(&st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" {
		ForwardEx(session, "1", nil,
			"GetHosProFullList param is empty,HosID=%s\n", st.HosID)
		return
	}

	type all_pro struct {
		New     []*common.ST_Product
		OnLine  []*common.ST_Product
		OffLine []*common.ST_Product
		Modify  []*common.ST_Product
		UnPass  []*common.ST_Product
	}
	all := &all_pro{}

	data := common.QueryMoreProducts(common.GetHosProductList(st.HosID))
	for _, v := range data {
		if v.Current.OpreatStatus == constant.OperatingStatus_new {
			all.New = append(all.New, v)
		}
		if v.Current.OpreatStatus == constant.OperatingStatus_modify {
			all.Modify = append(all.Modify, v)
		}
		if v.Current.OpreatStatus == constant.OperatingStatus_online {
			all.OnLine = append(all.OnLine, v)
		}
		if v.Current.OpreatStatus == constant.OperatingStatus_Reviewer_NotPass {
			all.UnPass = append(all.UnPass, v)
		}
		if v.Current.OpreatStatus == constant.OperatingStatus_Offline_self ||
			v.Current.OpreatStatus == constant.OperatingStatus_Offline_onforce {
			all.OffLine = append(all.OffLine, v)
		}
	}
	Forward(session, "0", all)
}

/////医院订单列表
func HosOrdergList(session *JsNet.StSession) {
	type st_get struct {
		HosID string //
	}
	st := &st_get{}
	if err := session.GetPara(&st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" {
		ForwardEx(session, "1", nil,
			"HosOrdergList param is empty,HosID=%s\n", st.HosID)
		return
	}
	data, err := common.GetHosOrder(st.HosID)

	sort.Slice(data.Check, func(i, j int) bool {
		t1, e1 := time.Parse("2006-01-02 15:04:05", data.Check[i].CheckDate)
		t2, e2 := time.Parse("2006-01-02 15:04:05", data.Check[j].CheckDate)
		if e1 == nil && e2 == nil {
			return t1.Unix() > t2.Unix()
		}
		return false
	})

	sort.Slice(data.Success, func(i, j int) bool {
		t1, e1 := time.Parse("2006-01-02 15:04:05", data.Success[i].CheckDate)
		t2, e2 := time.Parse("2006-01-02 15:04:05", data.Success[j].CheckDate)
		if e1 == nil && e2 == nil {
			return t1.Unix() > t2.Unix()
		}
		return false
	})

	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Forward(session, "0", data)
}

///医院校验
func OrderHosVerify(session *JsNet.StSession) {
	type st_new struct {
		HosID     string
		VerifCode string
	}
	st := &st_new{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" || st.VerifCode == "" {
		ForwardEx(session, "1", nil, "param is empty,HosID=%s,VerifCode=%s\n", st.HosID, st.VerifCode)
		return
	}
	data, err := common.OrderHosVerify(st.HosID, st.VerifCode)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

////校验码查询订单
func VerifyQueryOrder(session *JsNet.StSession) {
	type st_get struct {
		HosID     string
		VerifCode string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" || st.VerifCode == "" {
		ForwardEx(session, "1", nil, "param is empty")
		return
	}
	data, err := common.VerifyQueryOrder(st.HosID, st.VerifCode)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

//医院确认手术，并追加金额
func HosConfirmOpreat(session *JsNet.StSession) {
	type info struct {
		HosID   string //医院ID
		OrderID string //订单id
		Price   int    //追加的钱（可以为0）
		Note    string //备注
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if err := common.HosConfirmOperation(st.HosID, st.OrderID, st.Note, st.Price); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", nil)
}

////手机短信验证
func mobile_veirfy(session *JsNet.StSession) {
	type MobilePara struct {
		SignName string
		Mobile   string
		Expire   int
		SmsCode  string
	}
	type Return struct {
		Ret string
		Msg string
	}
	ret := &Return{}
	para := &MobilePara{}
	e := session.GetPara(para)
	if e != nil {
		Error(e.Error())
		ret.Ret = "1"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}

	JsMobile.ComJsMobileVerify(para.SignName, para.Mobile, para.SmsCode, "a", 300, nil)
	//RegisterAuth(para.Mobile, "星喜医美", para.Expire)

	ret.Ret = "0"
	ret.Msg = "success"
	session.Forward(ret)
}

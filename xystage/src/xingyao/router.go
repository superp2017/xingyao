package main

import (
	"JsGo/JsMobile"
	"JsLib/JsDispatcher"
	"JsLib/JsNet"
	"cache/cacheIO"
	"common"
	"constant"
	"sort"
	. "util"
)

func Init_Router() {

	JsDispatcher.Http("/setcitymap", SetCityMap) //设置全局的城市映射
	/////////////////////////////////订单接口///////////////////////////////////////////////
	JsDispatcher.Http("/queryorderinfo", QueryOrderInfo)     //查看订单详情
	JsDispatcher.Http("/orderappointment", OrderAppointment) //系统预约
	JsDispatcher.Http("/ordercancle", OrderCancle)           //订单取消

	/////////////////////////////////代理人接口////////////////////////////////////////////
	JsDispatcher.Http("/getcityagentinfo", common.GetCityAgent) //获取城市代理人信息
	JsDispatcher.Http("/reviewagent", common.RevieweAgent)      //审核代理
	JsDispatcher.Http("/offlineagent", common.OfflineAgent)     //下线代理
	JsDispatcher.Http("/onlineagent", common.OnlineAgent)       //下线代理

	JsDispatcher.Http("/getsimplehospital", common.GetHosSimpleInfo) //获取全局的简短的医院信息
	////////////////////////多个查询接口//////////////////////////////////////
	JsDispatcher.Http("/querymorehospital", common.QueryMoreHospital)   //查询多个医院信息
	JsDispatcher.Http("/querymoredoctor", common.QueryMoreDoctors)      //查询多个医生信息
	JsDispatcher.Http("/querymoreproduct", common.QueryMoreProductInfo) //查询多个产品信息
	JsDispatcher.Http("/querymoreorder", common.QueryMoreOrderInfo)     //查询多个订单信息

	JsDispatcher.Http("/queryproductinfo", common.QueryProduct)

	JsDispatcher.Http("/ModifyProductRatio", common.ModifyProductRatio)   //修改产品结算比例
	JsDispatcher.Http("/ModifyAgentProRatio", common.ModifyAgentProRatio) //修改代理小B的产品结算比例

	JsDispatcher.Http("/getsearchcontent", GetSearcheContent) //获取搜索Hash

	JsDispatcher.Http("/getglobalwithdrawrecord", common.GetGlobalWithDrawRecord) //获取全局的提现记录
	JsDispatcher.Http("/querywithdrawinfo", common.QueryWithDrawInfo)             //查询提现详情
	JsDispatcher.Http("/confirmwithdraw", common.ComfirmWithDraw)                 ///后台确认提现
	JsDispatcher.Http("/getwithdrawmonth", common.GetWithDrawMonth)               //获取每个月的提现记录
	JsDispatcher.Http("/getagentinvitestatistic", common.GetAgentInviteStatistic) //获取小B邀请统计
	JsDispatcher.Http("/getuserincrease", common.GetUserIncrease)                 //获取用户每天的增加量

	JsDispatcher.Http("/getglobalonlieagents", GetGlobalonlineAgent) //获取全局在线的小B列表

	JsDispatcher.Http("/getdoctors", common.GetDedicateDoctor)

	JsDispatcher.Http("/reviewnewdoctor", common.RevieweNewDoctor)
	JsDispatcher.Http("/reviewmodifieddoctor", common.RevieweModifyDoctor)
	JsDispatcher.Http("/forcedoctoroffline", common.DocOfflineOnForce)

	JsDispatcher.Http("/uptoexpertdoc", common.UpToExpertDoc)             //提升为大牌名医
	JsDispatcher.Http("/downdocexpert", common.DownExpertDoc)             //大牌名医下架
	JsDispatcher.Http("/getexpertdoctors", common.GetExpertDoctor)        ////获取大牌名医列表
	JsDispatcher.Http("/gethomeexpertdoctor", common.GetHomeExpertDoctor) ////获取首页大牌名医
	JsDispatcher.Http("/sethomedoctor", common.SetHomeExpertDoctor)       ////设置首页大牌名医

	JsDispatcher.Http("/getglobalorderinfo", GetGlobalOrderinfo) ////获取全局的所有订单信息

	JsDispatcher.Http("/manualsettlement", ManualSettlement) ////手动结算
	JsDispatcher.Http("/manualcollection", ManualCollection) ////确认收款

	JsDispatcher.Http("/deldoctor", common.DelDoctor)     //删除医生
	JsDispatcher.Http("/delproduct", common.DelProduct)   //删除产品
	JsDispatcher.Http("/delhospital", common.DelHospital) //删除医院



}








///设置城市map
func SetCityMap(session *JsNet.StSession) {
	common.SetCityMap(session)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
///查询订单详情
func QueryOrderInfo(session *JsNet.StSession) {
	common.QueryOrderInfo(session)
}

/////系统后台预约订单
func OrderAppointment(session *JsNet.StSession) {
	type st_new struct {
		OpreatJobNum    string
		OrderID         string
		OpreatName      string
		AppointmentDate string
		AppointmentDes  string
	}

	st := &st_new{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if st.OrderID == "" || st.OpreatJobNum == "" || st.OpreatName == "" || st.AppointmentDate == "" {
		ForwardEx(session, "1", nil, "param is empty,OrderID=%s,OpreatJobNum=%s,OpreatName=%s,AppointmentDate=%s\n",
			st.OrderID, st.OpreatJobNum, st.OpreatName, st.AppointmentDate)
		return
	}
	data, err := common.OrderAppointment(st.OrderID, st.OpreatJobNum, st.OpreatName, st.AppointmentDate, st.AppointmentDes)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if data != nil {
		//VerifyPass_ex(data.UserCell, data.UserName, data.AppointmentDate, data.HosName, data.ProFirstItem, data.VerifCode, "021-80392781-998", 3600)
		m := make(map[string]string)
		m["name"] = data.UserName
		m["date"] = data.AppointmentDate
		m["hosp"] = data.HosName
		m["prod"] = data.ProFirstItem
		m["code"] = data.VerifCode
		JsMobile.ComJsMobileVerify("喜妹儿", data.UserCell, "SMS_126782106", "b", 3600, m)
	}

	Forward(session, "0", data)
}

//手动生成账单
func ManualSettlement(session *JsNet.StSession) {

	type info struct {
		OrderID string //订单id
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.OrderID == "" {
		ForwardEx(session, "1", nil, "ManualSettlement failed,OrderID is empty\n")
		return
	}
	data, err := common.SysGenBill(st.OrderID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

//系统确认收款
func ManualCollection(session *JsNet.StSession) {
	type info struct {
		OrderID string
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.OrderID == "" {
		ForwardEx(session, "1", nil, "ManualCollection  OrderID is empty\n")
		return
	}
	if err := common.OrderSysConfirmCollection(st.OrderID, "确认收款", "admin", "admin"); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", nil)
}

///后台取消订单
func OrderCancle(session *JsNet.StSession) {
	type st_new struct {
		OrderID string
		Msg     string
		JobNum  string
		Name    string
	}
	st := &st_new{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.OrderID == "" || st.JobNum == "" {
		ForwardEx(session, "1", nil, "param is empty,OrderID=%s,JobNum=%s\n", st.OrderID, st.JobNum)
		return
	}
	data, err := common.OrderCancle(st.OrderID, st.Msg, st.JobNum, st.Name)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func GetSearcheContent(session *JsNet.StSession) {

	data, err := cacheIO.GetSearchHash()
	if err != nil {
		Forward(session, "1", err.Error())
		return
	}
	Forward(session, "0", data.Data)
}

func GetGlobalonlineAgent(session *JsNet.StSession) {

	data, err := common.GetGlobalAgent()
	if err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}

	var list = make([]*common.ST_User, 0)
	UserList := common.GetMoreUserInfo(data.Ids)
	for _, v := range UserList {
		if v.Agent.Current.OpreatStatus == constant.Agent_Offline_self || v.Agent.Current.OpreatStatus == constant.Agent_Online {
			list = append(list, v)
		}
	}
	Forward(session, "0", data)

}

func GetGlobalOrderinfo(session *JsNet.StSession) {

	type ST_GlobalOrder struct {
		Appointment []*common.ST_Order //待预约
		Verify      []*common.ST_Order //待校验
		Confirm     []*common.ST_Order //待确认
		Settlement  []*common.ST_Order //待结算
		Collection  []*common.ST_Order //待收款
		Complete    []*common.ST_Order //已完成
		Cancel      []*common.ST_Order //已取消
	}
	data := ST_GlobalOrder{}
	ids := common.GetGlobalOrderList()
	sort.Sort(sort.Reverse(sort.StringSlice(ids)))
	list := common.QueryMoreOrders(ids)
	for _, v := range list {
		if v.Current.OpreatStatus == constant.Status_Order_PenddingAppointment {
			data.Appointment = append(data.Appointment, v)
		}
		if v.Current.OpreatStatus == constant.Status_Order_PenddingVerify {
			data.Verify = append(data.Verify, v)
		}
		if v.Current.OpreatStatus == constant.Status_Order_PendingConfirm {
			data.Confirm = append(data.Confirm, v)
		}
		if v.Current.OpreatStatus == constant.Status_Order_PendingStatements {
			data.Settlement = append(data.Settlement, v)
		}
		if v.Current.OpreatStatus == constant.Status_OrderPenddingCollection {
			data.Collection = append(data.Collection, v)
		}
		if v.Current.OpreatStatus == constant.Status_Order_Succeed {
			data.Complete = append(data.Complete, v)
		}
		if v.Current.OpreatStatus == constant.Status_Order_CancleBeforeVerfy ||
			v.Current.OpreatStatus == constant.Status_Order_CancleAfterVerfy {
			data.Cancel = append(data.Cancel, v)
		}
	}
	Forward(session, "0", data)
}

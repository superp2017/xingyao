package main

import (
	"JsLib/JsDispatcher"
	// . "JsLib/JsLogger"
	"JsLib/JsNet"
	"common"
	"constant"
	// "encoding/csv"
	"ider"
	// "io"
	"sort"
	. "util"
)

func init_testOrder() {

	JsDispatcher.Http("/NewOrder", NewOrder)                                   //新建订单
	JsDispatcher.Http("/OrderUserPaid", OrderUserPaid)                         //用户支付
	JsDispatcher.Http("/OrderAppointment", OrderAppointment)                   //系统预约
	JsDispatcher.Http("/OrderHosVerify", OrderHosVerify)                       //医院校验
	JsDispatcher.Http("/OrderInvalid", OrderInvalid)                           //订单无效
	JsDispatcher.Http("/OrderCancle", OrderCancle)                             //订单取消
	JsDispatcher.Http("/SysGenBill", SysGenBill)                               //系统生成账单
	JsDispatcher.Http("/OrderHosReconcile", OrderHosReconcile)                 //确认对账
	JsDispatcher.Http("/OrderSysConfirmCollection", OrderSysConfirmCollection) //系统确认收款
	JsDispatcher.Http("/ordercomment", OrderComment)                           //用户评论
	JsDispatcher.Http("/GetGlobalOrder", GetGlobalOrder)                       //获取全局的订单
	JsDispatcher.Http("/GetHosOrder", GetHosOrder)                             //获取某个医院的订单
	JsDispatcher.Http("/VerifyQueryOrder", VerifyQueryOrder)                   //校验码获取订单号

	JsDispatcher.Http("/UniqueHosOrder", UniqueHosOrder) //
	JsDispatcher.Http("/UpdateGlobalOrder", UpdateGlobalOrder)

	//JsDispatcher.Http("/moveHosOrder", moveHosOrder) //

}

func NewOrder(session *JsNet.StSession) {
	type st_new struct {
		UID   string
		ProID string
	}
	st := &st_new{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ProID == "" || st.UID == "" {
		ForwardEx(session, "1", nil, "param is empty")
		return
	}
	order := &common.ST_Order{}
	id, err := ider.GenOrderID()
	if err != nil {
		ForwardEx(session, "1", nil, "ider.GenOrderID failed....\n")
		return
	}
	order.OrderID = id
	order.UID = st.UID
	order.ProID = st.ProID
	if err := common.SubmitOrder(order); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", order)
}

func OrderUserPaid(session *JsNet.StSession) {
	type st_new struct {
		OrderID string
		Cb      common.WxST_PayCb
	}
	st := &st_new{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.OrderID == "" {
		ForwardEx(session, "1", nil, "param is empty")
		return
	}
	data, err := common.OrderUserPaid(st.OrderID, &st.Cb)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

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
		ForwardEx(session, "1", nil, "param is empty")
		return
	}
	data, err := common.OrderAppointment(st.OrderID, st.OpreatJobNum, st.OpreatName, st.AppointmentDate, st.AppointmentDes)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func OrderHosVerify(session *JsNet.StSession) {
	type st_new struct {
		HosID        string
		VerifCode    string
		OpreatJobNum string
		OpreatName   string
	}
	st := &st_new{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" || st.VerifCode == "" || st.OpreatJobNum == "" {
		ForwardEx(session, "1", nil, "param is empty")
		return
	}
	data, err := common.OrderHosVerify(st.HosID, st.VerifCode, st.OpreatJobNum, st.OpreatName)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

///限时的产品到期，状态：待付款->已失效（已取消）
func OrderInvalid(session *JsNet.StSession) {
	type st_new struct {
		OrderID string
	}
	st := &st_new{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.OrderID == "" {
		ForwardEx(session, "1", nil, "param is empty")
		return
	}
	data, err := common.OrderInvalid(st.OrderID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

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
		ForwardEx(session, "1", nil, "param is empty")
		return
	}
	data, err := common.OrderCancle(st.OrderID, st.Msg, st.JobNum, st.Name)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

//系统生成订单,状态：待计算
func SysGenBill(session *JsNet.StSession) {
	type st_new struct {
		OrderID string
	}
	st := &st_new{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.OrderID == "" {
		ForwardEx(session, "1", nil, "param is empty")
		return
	}
	data, err := common.SysGenBill(st.OrderID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func OrderHosReconcile(session *JsNet.StSession) {
	type st_new struct {
		HosID   string
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
	if st.HosID == "" || st.OrderID == "" || st.JobNum == "" {
		ForwardEx(session, "1", nil, "param is empty")
		return
	}
	err := common.OrderHosReconcile(st.HosID, st.OrderID, st.Msg, st.JobNum, st.Name)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", nil)
}

func OrderSysConfirmCollection(session *JsNet.StSession) {
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
		ForwardEx(session, "1", nil, "param is empty")
		return
	}
	err := common.OrderSysConfirmCollection(st.OrderID, st.Msg, st.JobNum, st.Name)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", nil)
}

//订单评论
func OrderComment(session *JsNet.StSession) {
	type st_get struct {
		OrderID string
		UID     string
		Name    string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.OrderID == "" || st.UID == "" || st.Name == "" {
		ForwardEx(session, "1", nil, "OrderComment failed,OrderID=%s,UID=%s,Name=%s\n", st.OrderID, st.UID, st.Name)
		return
	}
	if err := common.OrderComment(st.OrderID, st.UID, st.Name); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", nil)
}

func GetGlobalOrder(session *JsNet.StSession) {

	type ST_GlobalOrder struct {
		Check      []string //待确认
		Settlement []string //待结算
		Verify     []string
		Complete   []string //已完成
		Cancle     []string //已取消
	}
	data := ST_GlobalOrder{}
	ids := common.GetGlobalOrderList()
	sort.Sort(sort.Reverse(sort.StringSlice(ids)))
	list := common.QueryMoreOrders(ids)
	for _, v := range list {
		if v.Current.OpreatStatus == constant.Status_Order_PenddingAppointment {
			data.Check = append(data.Check, v.OrderID)
		}
		if v.Current.OpreatStatus == constant.Status_Order_PenddingVerify {
			data.Verify = append(data.Verify, v.OrderID)
		}
		if v.Current.OpreatStatus == constant.Status_Order_PendingStatements {
			data.Settlement = append(data.Settlement, v.OrderID)
		}
		if v.Current.OpreatUserStatus == constant.Status_Order_Succeed {
			data.Complete = append(data.Complete, v.OrderID)
		}
		if v.Current.OpreatStatus == constant.Status_Order_CancleBeforeVerfy ||
			v.Current.OpreatStatus == constant.Status_Order_CancleAfterVerfy {
			data.Cancle = append(data.Cancle, v.OrderID)
		}
	}

	Forward(session, "0", data)
}

func GetHosOrder(session *JsNet.StSession) {
	type st_get struct {
		HosID string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" {
		ForwardEx(session, "1", nil, "param is empty")
		return
	}
	data, err := common.GetHosOrder(st.HosID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

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

func UniqueHosOrder(session *JsNet.StSession) {
	hos := gl_hospital()

	for _, v := range hos {
		data := []string{}
		if err := WriteLock(constant.Hash_HospitalOrder, v, &data); err != nil {
			ForwardEx(session, "1", nil, err.Error())
			return
		}
		cur := []string{}
		for _, v := range data {
			exist := false
			for _, v1 := range cur {
				if v == v1 {
					exist = true
					break
				}
			}
			if !exist {
				cur = append(cur, v)
			}
		}
		if err := WriteBack(constant.Hash_HospitalOrder, v, &cur); err != nil {
			ForwardEx(session, "1", nil, err.Error())
			return
		}
	}
	Forward(session, "0", nil)
}

// func moveHosOrder(session *JsNet.StSession) {
// 	list, err := gl_hospital()
// 	if err != nil {
// 		ForwardEx(session, "1", nil, err.Error())
// 		return
// 	}

// 	for _, v := range list {
// 		cache, err := cacheIO.GetHosOrderList(v)
// 		if err == nil {
// 			data := []string{}
// 			if cache.Pay != nil {
// 				data = append(data, cache.Pay...)
// 			}
// 			if cache.Appointment != nil {
// 				data = append(data, cache.Appointment...)
// 			}
// 			if cache.Check != nil {
// 				data = append(data, cache.Check...)
// 			}
// 			if cache.Settlement != nil {
// 				data = append(data, cache.Settlement...)
// 			}
// 			if cache.Reconciliate != nil {
// 				data = append(data, cache.Reconciliate...)
// 			}
// 			if cache.Receive != nil {
// 				data = append(data, cache.Receive...)
// 			}
// 			if cache.Cancle != nil {
// 				data = append(data, cache.Cancle...)
// 			}
// 			if cache.Complete != nil {
// 				data = append(data, cache.Complete...)
// 			}

// 			if err := DirectWrite(constant.Hash_HospitalOrder, v, &data); err != nil {
// 				ErrorLog("moveHosOrder DirectWrite failed,HosID=%s\n", v)
// 			}

// 		}
// 	}
// 	Forward(session, "0", nil)
// }

// Pay          []string   //待支付
// Appointment  []string   //待预约
// Check        []string   //待校验
// Settlement   []string   //待结算
// Reconciliate []string   //待对账
// Receive      []string   //待收款
// Refund       []string   //已退款
// Cancle       []string   //已取消
// Complete     []string   //已完成

func UpdateGlobalOrder(session *JsNet.StSession) {
	hosList := gl_hospital()

	// Error("hosList:=%v", hosList)
	// sort.Sort(sort.Reverse(sort.StringSlice(hosList)))
	data := []string{}
	for _, v := range hosList {
		ids := []string{}
		if err := ShareLock(constant.Hash_HospitalOrder, v, &ids); err != nil {
			if e1 := ShareLock(constant.Hash_HospitalOrder, v, &ids); e1 != nil {
				DirectWrite(constant.Hash_HospitalOrder, v, &ids)
			}
		}
		if len(ids) > 0 {
			data = append(data, ids...)
		}
	}
	sort.Sort(sort.StringSlice(data))
	if err := DirectWrite(constant.Hash_Global, constant.KEY_GlobalOrderList, &data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", nil)
}

func ExportYMJ(session *JsNet.StSession) {

}

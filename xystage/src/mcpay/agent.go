package main

import (
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"common"
	"constant"
	"ider"
	"time"

	"fmt"
	"strings"
)

func init_agent() {
	JsDispatcher.Http("/newwxbondorder", new_wx_bond_order)
	// JsDispatcher.Http("/paybondorder", pay_wx_bond_success_cb)
	JsDispatcher.Http("/refundbond", refund_wx_bond_order)
	JsDispatcher.Http("/offlineagent", offline_agent)
}

func new_wx_bond_order(session *JsNet.StSession) {
	Info("Enter new bond order\n")
	para := &common.ST_Order{}
	type order_ret struct {
		Ret   string
		Msg   string
		Order *common.ST_Order
	}
	ret := &order_ret{"1", "", nil}

	e := session.GetPara(para)
	if e != nil {
		g_log.Error(e.Error())
		ret.Ret = "2"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}
	if para.UID == "" || para.PlatformFee < 0 || para.Bond < 0 || para.FranchiseFee < 0 {
		ret.Ret = "2"
		ret.Msg = fmt.Sprintf("new_wx_bond_order failed,UID=%s,PlatformFee=%d,Bond=%d,FranchiseFee=%d\n",
			para.UID, para.PlatformFee, para.Bond, para.FranchiseFee)
		session.Forward(ret)
		return
	}

	if para.AgentLevel != constant.Agent_Level_Diamonds_A {

		e1 := ErrorLog("代理人的佣金层级不合法,AgentLevel=%s\n", para.AgentLevel)
		ret.Ret = "2"
		ret.Msg = e1.Error()
		session.Forward(ret)
		return
	}

	addr := session.RemoteAddr()
	i := strings.Index(addr, ":")

	para.TerminalIp = addr[:i]

	id, err := ider.GenOrderID()
	if err != nil {
		g_log.Error(err.Error())
		ret.Ret = "1"
		ret.Msg = err.Error()
		session.Forward(ret)
		return
	}
	para.OrderID = id

	Info("Agent OrderID=%s\n", id)

	para.ServiceTimeStamp = time.Now().Unix() + C_TIMEAREA*3600

	ch, e := wx_pub_pay(para)
	if e != nil {
		Error(e.Error())
		ret.Ret = "4"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}
	para.Charge = ch

	if err := common.SubmitBondOrder(para); err != nil {
		ret.Ret = "1"
		ret.Msg = err.Error()
		ret.Order = nil
		session.Forward(ret)
		return
	}
	ret.Ret = "0"
	ret.Order = para
	session.Forward(ret)
}

// func pay_wx_bond_success_cb(session *JsNet.StSession) {
// 	body := session.Body()

// 	paycb := &common.WxST_PayCb{}
// 	e := xml.Unmarshal(body, paycb)

// 	xml := ""
// 	if e != nil {
// 		xml = `<xml>
//   				<return_code><![CDATA[FAIL]]></return_code>
//   				<return_msg><![CDATA[` + e.Error() + `]]></return_msg>
// 			   </xml>`

// 	} else {
// 		if _, err := common.OrderPayBond(paycb.Out_trade_no, paycb); err != nil {
// 			Error(err.Error())
// 		}

// 		xml = `<xml>
//   				<return_code><![CDATA[SUCCESS]]></return_code>
//   				<return_msg><![CDATA[OK]]></return_msg>
// 			</xml>`
// 	}

// 	session.DirectWrite(xml)
// }

func refund_wx_bond_order(session *JsNet.StSession) {
	type refund_order_para struct {
		OrderID string
	}

	type refund_order_ret struct {
		Ret string
		Msg string
	}

	ret := &refund_order_ret{}
	para := &refund_order_para{}

	e := session.GetPara(para)
	if e != nil {
		Error(e.Error())
		ret.Ret = "1"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}

	order, e := common.QueryOrder(para.OrderID)
	if e != nil {
		Error(e.Error())
		ret.Ret = "2"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}
	/////保证金退款
	cb, e := wx_bond_refund(order)
	if e != nil {
		Error(e.Error())
		ret.Ret = "3"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}
	if cb["result_code"] == "FAIL" {
		ErrorLog("微信退款失败,return_msg:%s,err_code_des:%s\n", cb["return_msg"], cb["err_code_des"])
		ret.Ret = "3"
		ret.Msg = fmt.Sprintf("微信退款失败\n")
		session.Forward(ret)
		return
	}
	if _, err := common.WithdrawBond(order.OrderID, cb); err != nil {
		ret.Ret = "1"
		ret.Msg = err.Error()
		session.Forward(ret)
		return
	}

	ret.Ret = "0"
	ret.Msg = "success"
	session.Forward(ret)
}

func offline_agent(session *JsNet.StSession) {
	type refund_order_para struct {
		UID       string //uid
		JodNumber string //操作员工号
		JobName   string //操作员名字
		Reason    string //原因
	}

	type refund_order_ret struct {
		Ret string
		Msg string
	}

	ret := &refund_order_ret{}
	para := &refund_order_para{}

	e := session.GetPara(para)
	if e != nil {
		Error(e.Error())
		ret.Ret = "1"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}

	if para.UID == "" || para.JodNumber == "" {
		ret.Ret = "1"
		ret.Msg = fmt.Sprintf("offline_agent failed,UID=%s,JodNumber=%s\n", para.UID, para.JodNumber)
		session.Forward(ret)
		return
	}

	user, e := common.GetUserInfo(para.UID)
	if e != nil {
		Error(e.Error())
		ret.Ret = "2"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}
	if user.Agent == nil || user.Agent.BondOrderID == "" {
		ret.Ret = "2"
		ret.Msg = fmt.Sprintf("offline_agent failed,用户的代理信息不全，UID=%s，Agent=%v\n", para.UID, user.Agent)
		session.Forward(ret)
		return
	}

	order, e := common.QueryOrder(user.Agent.BondOrderID)
	if e != nil {
		Error(e.Error())
		ret.Ret = "2"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}
	/////refund  bond
	cb, e := wx_bond_refund(order)
	if e != nil {
		Error(e.Error())
		ret.Ret = "3"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}

	if err := common.AgentOffline(order.OrderID, para.JodNumber, para.JobName, para.Reason, cb); err != nil {
		ret.Ret = "1"
		ret.Msg = err.Error()
		session.Forward(ret)
		return
	}

	ret.Ret = "0"
	ret.Msg = "success"
	session.Forward(ret)
}

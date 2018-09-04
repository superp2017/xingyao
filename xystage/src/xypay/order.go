package main

import (
	_ "JsLib/JsConfig"
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"JsLib/JsOrder"
	"common"
	"encoding/xml"
	"ider"

	"fmt"
	"strings"
	"sync"
	"time"
)

func order_init() {

	JsDispatcher.Http("/newwxpuborder", new_wxpub_order)
	JsDispatcher.Http("/paysuccesscb", pay_success_cb)
	JsDispatcher.Http("/refundorder", refund_order) //退款
	// JsDispatcher.Http("/refundorders", refund_orders)
	JsDispatcher.Http("/withdrawbalance", BalanceWithdraw)
}

var g_pending_ids []string
var g_mutex sync.Mutex

func init() {

}

func pay_success_cb(session *JsNet.StSession) {
	body := session.Body()

	paycb := &common.WxST_PayCb{}

	e := xml.Unmarshal(body, paycb)

	xml := ""
	if e != nil {
		xml = `<xml>
  				<return_code><![CDATA[FAIL]]></return_code>
  				<return_msg><![CDATA[` + e.Error() + `]]></return_msg>
			   </xml>`

	} else {

		//get out the order

		order, err := common.QueryOrder(paycb.Out_trade_no)
		if err == nil {
			//judge the order and go to the correct path
			if order.OrderType == 1 {
				Info("common.OrderUserPaid .............................")
				if _, err := common.OrderUserPaid(paycb.Out_trade_no, paycb); err != nil {
					Error(err.Error())
				}

			} else {
				Info("common.OrderPayBond.............................")
				if _, err := common.OrderPayBond(paycb.Out_trade_no, paycb); err != nil {
					Error(err.Error())
				}

			}

		} else {
			Error("common.QueryOrder failed,err:\n", err.Error())
		}

		xml = `<xml>
  				<return_code><![CDATA[SUCCESS]]></return_code>
  				<return_msg><![CDATA[OK]]></return_msg>
			</xml>`
	}

	session.DirectWrite(xml)
}

func new_wxpub_order(session *JsNet.StSession) {
	Info("Enter order new")
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

	addr := session.RemoteAddr()
	i := strings.Index(addr, ":")

	para.TerminalIp = addr[:i]

	id, err := ider.GenOrderID()
	Info("OrderID=%s\n", id)
	if err != nil {
		g_log.Error(err.Error())
		ret.Ret = "1"
		ret.Msg = err.Error()
		session.Forward(ret)
		return
	}
	para.OrderID = id

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
	Info("Charge=%v\n", ch)

	if err := common.SubmitOrder(para); err != nil {
		ret.Ret = "1"
		ret.Msg = err.Error()
		ret.Order = nil
		session.Forward(ret)
		return
	}
	Info("Submit sucess\n")
	ret.Ret = "0"
	ret.Order = para
	session.Forward(ret)
}

func refund_order(session *JsNet.StSession) {
	type refund_order_para struct {
		OrderID string
		Msg     string
		JobNum  string
		JobName string
	}

	type refund_order_ret struct {
		Ret    string
		Msg    string
		Entity common.ST_Order
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

	cb, e := wx_pub_refund(order)
	if e != nil {
		Error(e.Error())
		ret.Ret = "3"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}

	Info("result_code=%v\n", cb)

	if cb["result_code"] == "FAIL" {
		ErrorLog("微信退款失败,return_msg:%s,err_code_des:%s\n", cb["return_msg"], cb["err_code_des"])
		ret.Ret = "3"
		ret.Msg = "微信退款失败\n"
		session.Forward(ret)
		return
	}

	//ErrorLog("Refund Money FeedBack=%v", cb)

	if err := common.OrderRefund(order, cb, para.Msg, para.JobNum, para.JobName); err != nil {
		ret.Ret = "1"
		ret.Msg = err.Error()
		session.Forward(ret)
		return
	}

	ret.Ret = "0"
	ret.Msg = "success"
	ret.Entity = *order
	session.Forward(ret)
}

func BalanceWithdraw(session *JsNet.StSession) {
	type withdraw_ret struct {
		Ret string
		Msg string
	}
	type st_para struct {
		UID   string
		Money int
		Des   string
	}
	st := &st_para{}
	ret := withdraw_ret{}
	if err := session.GetPara(st); err != nil {
		ret.Ret = "1"
		ret.Msg = err.Error()
		session.Forward(ret)
		return
	}
	if st.UID == "" || st.Money <= 0 {
		ret.Ret = "1"
		ret.Msg = fmt.Sprintf("param error,UID =%s,Money=%d", st.UID, st.Money)
		session.Forward(ret)
		return
	}
	user, err := common.GetUserInfo(st.UID)
	if err != nil {
		ret.Ret = "1"
		ret.Msg = err.Error()
		session.Forward(ret)
		return
	}

	transfer := &JsOrder.ST_Transfer{
		UserName:   user.Name,
		OpenId:     user.OpenId_web,
		UserHeader: user.HeadImageURL,
		TimeStamp:  time.Now().Unix(),
		Desc:       st.Des,
		Amount:     st.Money,
	}
	if _, err := direct_transfer(transfer); err != nil {
		ret.Ret = "1"
		ret.Msg = err.Error()
		session.Forward(ret)
		return
	}
	if _, err := common.WithdrawBalance(st.UID, st.Money); err != nil {
		ret.Ret = "1"
		ret.Msg = err.Error()
		session.Forward(ret)
		return
	}
	ret.Ret = "0"
	ret.Msg = "success"
	session.Forward(ret)
	return
}

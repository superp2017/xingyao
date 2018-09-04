package main

import (
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"common"
	"constant"
	"ider"
	"strconv"
	. "util"
)

//产品详情
//产品详情
type ST_RequestOrderPar struct {
	UID         string
	OrderType   string
	RequestPage int
}

type ST_OrderNet struct {
	OrderInfo common.ST_Order
}

func Order_init() {

	JsDispatcher.Http("/getorders", GetOrders)
	JsDispatcher.Http("/getordersn", GetOrdersN)
	JsDispatcher.Http("/getagentorders", GetAgentOrders)

	JsDispatcher.Http("/getorder", GetNetOrder)
	JsDispatcher.Http("/getlsorder", GetNetLsOrder)              //所有医院列表
	JsDispatcher.Http("/gettotalpagenumorder", GettotalNumOrder) //所有医院列表                    //所有医院列表

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

func GettotalNumOrder(session *JsNet.StSession) {

	type ST_PageNumTotal struct {
		PageNum   int
		OrderType string
		UID       string
	}

	st := &ST_RequestOrderPar{}
	lsOrderID := []string{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if st.OrderType == constant.Status_Order_AgentEntity || st.OrderType == constant.Status_Order_AgentWeChat {
		lsOrderID = GetAgentOrderList(st)

	} else {
		lsOrderID = GetOrderList(st)

	}

	totalNum := &ST_PageNumTotal{}
	totalNum.OrderType = st.OrderType
	totalNum.PageNum = getCeilNum(len(lsOrderID), constant.ItemAccountPerPage_Order)
	// totalNum.HmTotalPageNum = make(map[string]int)
	// totalNum.HmTotalPageNum[st.RequireType] = getCeilNum(len(lsOrderID), constant.ItemAccountPerPage_Order)
	totalNum.UID = st.UID
	ForwardEx(session, "0", totalNum, st.OrderType)
}

func GetOrders(session *JsNet.StSession) {

	lsOrderSend := []*common.ST_Order{}
	//Get the request info
	st := &ST_RequestOrderPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Info("ST=%v\n", st)

	listID := GetDedicateListIDOrder(st, constant.ItemAccountPerPage_Order)
	Info("List ID=%v\n", listID)
	// if err != nil {
	// 	ForwardEx(session, "1", nil, err.Error())
	// 	return
	// }

	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}
	lsOrderSend = common.QueryMoreOrders(listID)

	// superOrder.HmOrder[st.OrderType] = lsOrderSend
	ForwardEx(session, "0", lsOrderSend, st.OrderType)
}

func GetOrdersN(session *JsNet.StSession) {

	lsOrderSend := []*common.ST_Order{}
	//Get the request info
	st := &ST_RequestOrderPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Info("ST=%v\n", st)

	listID := GetDedicateListIDOrder(st, constant.ItemAccountPerPage_Order)
	Info("List ID=%v\n", listID)
	// if err != nil {
	// 	ForwardEx(session, "1", nil, err.Error())
	// 	return
	// }

	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}
	lsOrderSend = common.QueryMoreOrders(listID)

	lsOrderSuper := []*ST_OrderNet{}
	for _, v := range lsOrderSend {
		superOrder := &ST_OrderNet{}
		superOrder.OrderInfo = *v
		lsOrderSuper = append(lsOrderSuper, superOrder)
	}

	// superOrder.HmOrder[st.OrderType] = lsOrderSend
	ForwardEx(session, "0", lsOrderSend, st.OrderType)
}

func GetDedicateListIDOrder(st *ST_RequestOrderPar, itemPerPage int) []string {

	listPageID := []string{}
	listID := []string{}

	if st.OrderType == constant.Status_Order_AgentEntity || st.OrderType == constant.Status_Order_AgentWeChat {
		listID = GetAgentOrderList(st)

	} else {
		listID = GetOrderList(st)

	}
	ErrorLog("len(listID)=%d,RequestPage=%d,itemPerPage=%d", len(listID), st.RequestPage, itemPerPage)
	listStartDex := (st.RequestPage - 1) * itemPerPage
	if len(listID) > 0 && len(listID) > listStartDex {
		if len(listID)-listStartDex <= itemPerPage {
			listPageID = listID[listStartDex:]
			//Info("List ID1=%v\n", listPageID)
			return listPageID
		} else {
			Info("List ID2=%v\n", listPageID)
			listPageID = listID[listStartDex : listStartDex+itemPerPage]
		}
	}
	//Info("List ID=%v\n", listPageID)
	return listPageID
}

func GetAgentOrderList(st *ST_RequestOrderPar) []string {
	Info("The Request Info=%v\n", st)

	// GetGlobalOrder() (*ST_OrderStatusCache, error) {
	lsID := []string{}

	user, err := common.GetUserInfo(st.UID)
	if err != nil {
		Error(err.Error())
		return lsID

	}

	if user.Agent == nil {
		Error("User with UID=%s is not an Agent\n", st.UID)
	}

	if user.Agent.Orders == nil {
		Info("There is no order for the user uid=%s\n", st.UID)
		return lsID

	}
	for _, v := range user.Agent.Orders {

		lsID = append(lsID, v.OrderID)
	}

	Info("+++++++++++++++++++++++++++++++++lsID=%s\n", lsID)
	return lsID

}
func GetOrderList(st *ST_RequestOrderPar) []string {
	Info("The Request Info=%v\n", st)

	res := GetGlobalOrder()
	if res == nil {
		return []string{}
	}
	//	Info("-----------------------------res =%v\n", res)

	if st.OrderType == constant.Status_Order_PenddingAppointment {
		Info(constant.Status_Order_PenddingAppointment)
		return res.Check

	} else if st.OrderType == constant.Status_Order_PenddingVerify {
		Info(constant.Status_Order_PenddingVerify)
		return res.Verify
	} else if st.OrderType == constant.Status_Order_PendingStatements {
		Info(constant.Status_Order_PendingStatements)
		return res.Settlement
	} else if st.OrderType == constant.Status_Order_Succeed {
		Info(constant.Status_Order_Succeed)
		return res.Complete

	} else if st.OrderType == constant.Status_Order_Cancle {
		Info(constant.Status_Order_Cancle)
		return res.Cancle

	} else {
		return []string{}
	}
	Info("+++++++++++++++++++++++++++++++++\n")
	return []string{}

}

func GetNetOrder(session *JsNet.StSession) {
	//get the request
	//hosST_HospitalNet := &ST_HospitalNet{}
	type st_query struct {
		OrderID string //医院id
	}
	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.OrderID == "" {
		ForwardEx(session, "1", nil, "order id为空,QueryHospitalInfo()，查询失败!")
		return
	}
	orderST_OrderNet, err := getNetOrderL(st.OrderID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", orderST_OrderNet)
}

func GetNetLsOrder(session *JsNet.StSession) {
	type st_queryOrder struct {
		OrderID []string
	}

	lsNetOrder := []*ST_OrderNet{}

	st := &st_queryOrder{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	for _, v := range st.OrderID {
		orderST_OrderNet, err := getNetOrderL(v)
		if err == nil {
			lsNetOrder = append(lsNetOrder, orderST_OrderNet)
		}
	}
	Forward(session, "0", lsNetOrder)
}

func getNetOrderL(OrderID string) (*ST_OrderNet, error) {
	OrderNet := &ST_OrderNet{}
	//get the common
	orderInfo, err := common.QueryOrder(OrderID)

	if err != nil {
		return nil, err

	}
	OrderNet.OrderInfo = *orderInfo
	return OrderNet, nil
}

func GenerateOrder(session *JsNet.StSession) {
	order := &common.ST_Order{}
	for i := 1; i < 20; i++ {
		order.OrderID = "ord-" + strconv.Itoa(i)
		order.HosName = "Hospital" + strconv.Itoa(i)
		order.AgentInfo.Name = "AgentName" + strconv.Itoa(i)
		DirectWrite(constant.Hash_Order, order.OrderID, order)
	}

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
	ordernet := &ST_OrderNet{}
	ordernet.OrderInfo = *order

	Forward(session, "0", ordernet)
}

func GetAgentOrders(session *JsNet.StSession) {

	lsOrderSend := []*common.ST_Order{}
	//Get the request info
	st := &ST_RequestOrderPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Info("ST=%v\n", st)

	listID := GetDedicateListIDOrder(st, constant.ItemAccountPerPage_Order)
	Info("List ID=%v\n", listID)
	// if err != nil {
	// 	ForwardEx(session, "1", nil, err.Error())
	// 	return
	// }

	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}
	lsOrderSend = common.QueryMoreOrders(listID)

	// superOrder.HmOrder[st.OrderType] = lsOrderSend
	ForwardEx(session, "0", lsOrderSend, st.OrderType)
}

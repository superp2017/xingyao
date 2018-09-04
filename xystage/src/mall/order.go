package main

import (
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"

	"common"
	"constant"
	"ider"
	"strconv"

	"sort"
	. "util"
)

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
	JsDispatcher.Http("/getagentorders", GetAgentOrders)
	JsDispatcher.Http("/getordersn", GetOrdersN)
	JsDispatcher.Http("/getorder", GetNetOrder)
	JsDispatcher.Http("/getlsorder", GetNetLsOrder)              //所有医院列表
	JsDispatcher.Http("/gettotalpagenumorder", GettotalNumOrder) //所有医院列表
	// JsDispatcher.Http("/neworder", NewOrder)                     //所有医院列表
	JsDispatcher.Http("/orderinvalid", OrderInvalid) //所有医院列表
	JsDispatcher.Http("/sharearticle", shareArticleCB)

}

func GettotalNumOrder(session *JsNet.StSession) {

	type ST_PageNumTotal struct {
		PageNum   int
		OrderType string
		UID       string
	}

	st := &ST_RequestOrderPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	lsOrderID := GetOrderList(st)
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

	listID := GetDedicateListIDOrder(st, constant.ItemAccountPerPage_Order)

	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}
	lsOrderSend = common.QueryMoreOrders(listID)

	// superOrder.HmOrder[st.OrderType] = lsOrderSend
	ForwardEx(session, "0", lsOrderSend, st.OrderType)
}

func GetAgentOrders(session *JsNet.StSession) {

	lsOrderSend := []*common.ST_Order{}
	//Get the request info
	st := &ST_RequestOrderPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	listID := GetDedicateListIDOrder(st, constant.ItemAccountPerPage_Order)

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
	listID := GetDedicateListIDOrder(st, constant.ItemAccountPerPage_Order)

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
	ForwardEx(session, "0", lsOrderSuper, st.OrderType)
}

func GetDedicateListIDOrder(st *ST_RequestOrderPar, itemPerPage int) []string {

	listPageID := []string{}
	listID := GetOrderList(st)
	listStartDex := (st.RequestPage - 1) * itemPerPage
	if len(listID)-listStartDex <= itemPerPage {
		listPageID = listID
		return listPageID
	} else {
		listPageID = listID[:listStartDex+itemPerPage]
	}
	return listPageID
}

func GetOrderList(st *ST_RequestOrderPar) []string {
	lsID := []string{}

	// for i := 0; i < 8; i++ {
	// 	orderID := "order-" + strconv.Itoa(10001+i)
	// 	lsID = append(lsID, orderID)
	// }

	// Info("+++++++++++++++++++++++++++++++++Request ID=%v\n", st)

	// return lsID
	//lsID, _ := common.GetCityItemDoc(cityName, bodypart)

	//lsID, _ = getOrders(st.City, st.BodyPart, st.OrderItem, st.MinPrice, st.MaxPrice)

	// getOrders(city, st.BodyPart, st.Item, st.MinPrice, st.MaxPrice)

	//get the user

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

	//get the user out

	//get the order id

	//return the orderID
	sort.Sort(sort.Reverse(sort.StringSlice(lsID)))
	return lsID
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

func NewOrderL(session *JsNet.StSession) {
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

func OrderInvalid(session *JsNet.StSession) {
	type st_require struct {
		OrderID string
	}
	ordernet := &ST_OrderNet{}

	para := &st_require{}
	e := session.GetPara(para)
	if e != nil {
		Error(e.Error())

		ForwardEx(session, "1", nil, e.Error())
		return
	}

	order, e := common.OrderInvalid(para.OrderID)
	if e != nil {
		Error(e.Error())

		ForwardEx(session, "1", nil, e.Error())
		return
	}
	ordernet.OrderInfo = *order
	Forward(session, "0", ordernet)
}

func shareArticleCB(session *JsNet.StSession) {

	type st_query struct {
		UID       string //医院id
		ArticleID string
	}

	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	// AgentSharedArticle(UID, ArticleID string) (e error, userFB *ST_User) {

	err, user := common.AgentSharedArticle(st.UID, st.ArticleID)
	if err != nil {
		ForwardEx(session, "1", nil, "Fail")
		return
	}

	ForwardEx(session, "0", user, "Sucess")

}

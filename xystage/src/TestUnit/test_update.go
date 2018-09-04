package main

import (
	"JsLib/JsDispatcher"
	"JsLib/JsNet"
	"cache/cacheIO"

	. "util"
)

func init_update() {
	JsDispatcher.Http("/GetHosChange", GetHosChange)           ///获取医院的更新
	JsDispatcher.Http("/GetDocChange", GetDocChange)           ///获取医生的更新
	JsDispatcher.Http("/GetProChange", GetProChange)           ///获取产品的更新
	JsDispatcher.Http("/GetOrderChange", GetOrderChange)       ///获取单个对象订单的更新
	JsDispatcher.Http("/GetAgentChange", GetAgentChange)       ///获取代理的更新
	JsDispatcher.Http("/GetAllOrderChange", GetAllOrderChange) ///获取所有的订单的更新
}

func GetHosChange(session *JsNet.StSession) {
	Forward(session, "0", cacheIO.GetHosChange())
}

func GetDocChange(session *JsNet.StSession) {

	Forward(session, "0", cacheIO.GetDocChange())
}

func GetProChange(session *JsNet.StSession) {

	Forward(session, "0", cacheIO.GetProChange())
}

func GetAgentChange(session *JsNet.StSession) {

	Forward(session, "0", cacheIO.GetAgentChange())
}

func GetOrderChange(session *JsNet.StSession) {
	type st_object struct {
		KeyID string
		Reset bool
	}
	st := &st_object{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.KeyID == "" {
		ForwardEx(session, "1", nil, "KeyID is empty\n")
		return
	}

	Forward(session, "0", cacheIO.GetOrderChange(st.KeyID, st.Reset))
}
func GetAllOrderChange(session *JsNet.StSession) {
	Forward(session, "0", cacheIO.GetAllOrderChange())
}

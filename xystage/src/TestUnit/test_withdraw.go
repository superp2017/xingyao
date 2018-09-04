package main

import (
	// "JsLib/JsConfig"
	"JsLib/JsDispatcher"
	// . "JsLib/JsLogger"
	"JsLib/JsNet"
	// . "cache/cacheLib"
	"common"
	// "constant"
	// "encoding/json"
	// "ider"
	// "strconv"
	. "util"
)

func init_withdraw() {
	JsDispatcher.Http("/NewWithDraw", NewWithDraw)
}

func NewWithDraw(session *JsNet.StSession) {
	type st_Get struct {
		UID   string
		Money int
	}

	st := &st_Get{}
	session.GetPara(st)

	user, err := common.GetUserInfo(st.UID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	common.NewWithDraw(user, st.Money, 0, 0)

	Forward(session, "0", 1)
}

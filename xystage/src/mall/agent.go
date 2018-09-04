package main

import (
	"JsLib/JsDispatcher"

	"JsLib/JsNet"
	"common"

	// "strings"
	// "time"
	. "util"
)

func Agent_init() {
	JsDispatcher.Http("/withdrawbalance", WithDrawMoney)
}

func WithDrawMoney(session *JsNet.StSession) {
	type ST_DrawMoneyPar struct {
		UID   string
		Money int
	}

	st := &ST_DrawMoneyPar{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	_, err := common.WithdrawBalance(st.UID, st.Money)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return

	}

	ForwardEx(session, "0", st.Money, "sucess")
}

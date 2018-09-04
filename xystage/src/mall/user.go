package main

import (
	"JsLib/JsDispatcher"
	"JsLib/JsNet"

	// . "JsLib/JsLogger"
	// "JsLib/JsNet"
	// . "cache/cacheIO"
	. "util"

	"common"
)

func User_init() {

	JsDispatcher.Http("/getuserid", GetUserID)

}

func GetUserInfo(session *JsNet.StSession) {

	common.QueryUserInfo(session)
}

func GetUserID(session *JsNet.StSession) {
	type st_query struct {
		OpenID string //openID
	}

	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if st.OpenID == "" {
		ForwardEx(session, "1", nil, "OPen ID=null")
		return
	}
	UID, err := common.GetUIDFromOpenID(st.OpenID)
	if err != nil {
		ForwardEx(session, "1", nil, "get uid fail")
		return
	}
	Forward(session, "0", UID)
}

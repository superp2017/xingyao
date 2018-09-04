package main

import (
	"JsLib/JsDispatcher"
	"JsLib/JsNet"
	"common"
	"constant"
	"strconv"
	. "util"
)

func init_user() {
	JsDispatcher.Http("/UserHisToContinue", UserHisToContinue)
	JsDispatcher.Http("/AllUserHisToContinue", AllUserHisToContinue)
}

func moveOrderHisToContinue(UID string) error {
	data := common.ST_User{}
	if er := Update(constant.Hash_User, UID, &data, func() {

		for _, v := range data.HisOrders {
			exist := false
			for _, v1 := range data.Orders {
				if v1 == v {
					exist = true
					break
				}
			}
			if !exist {
				data.Orders = append(data.Orders, v)
			}
		}

	}); er != nil {
		return er
	}
	return nil
}

func AllUserHisToContinue(session *JsNet.StSession) {

	for i := 2; i < 56; i++ {
		moveOrderHisToContinue("user-" + strconv.Itoa(10000+i))
	}

	Forward(session, "0", nil)
}

// user-10055

func UserHisToContinue(session *JsNet.StSession) {
	type st_get struct {
		UID string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" {
		ForwardEx(session, "1", nil, "UID is Empty\n")
		return
	}
	if er := moveOrderHisToContinue(st.UID); er != nil {
		ForwardEx(session, "1", nil, er.Error())
		return
	}
	Forward(session, "0", nil)
}

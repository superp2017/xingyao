package main

import (
	"JsLib/JsConfig"
	"JsLib/JsDispatcher"
	// . "JsLib/JsLogger"
	"JsLib/JsNet"
	. "cache/cacheLib"
	"common"
	"constant"
	// "encoding/json"
	"ider"
	"strconv"
	. "util"
)

func init_agent() {

	JsDispatcher.Http("/applyagent", common.ApplyAgent)               //申请成为小B
	JsDispatcher.Http("/reapplyagent", common.ReApplyAgent)           //修改后重新为代理人
	JsDispatcher.Http("/modifyagentinfo", common.ModifyAgentBaseInfo) //修改代理信息
	JsDispatcher.Http("/reviewagent", common.RevieweAgent)            //审核代理
	JsDispatcher.Http("/offlineagent", common.OfflineAgent)           //下线代理
	// JsDispatcher.Http("/onlineagent", common.OnlineAgent)             //上线代理
	JsDispatcher.Http("/withdrawagentbond", WithdrawBond) //撤销代理保证金
	JsDispatcher.Http("/submitbond", OrderSubmitBond)     //提交代理保证金
	JsDispatcher.Http("/orderpaybond", OrderPayBond)      //支付代理保证金

	JsDispatcher.Http("/invitetoemployee", common.InviteToEmployee) //邀请成为店员
	JsDispatcher.Http("/userbindagent", common.UserBindAgent)       //用户绑定代理
	JsDispatcher.Http("/withdrawbalance", WithdrawBalance)          //代理提现余额

	JsDispatcher.Http("/getcitygent", common.GetCityAgent)                //获取城市代理信息
	JsDispatcher.Http("/getcityemployee", common.GetCityEmployee)         //获取城市店员信息
	JsDispatcher.Http("/getglobalemployee", common.GetGlobalEmployeeInfo) //获取全局的店员信息

	JsDispatcher.Http("/ChangeAgentStatus", ChangeAgentStatus)
	JsDispatcher.Http("/ChangeAllGlobalStatus", ChangeAllGlobalStatus)

	JsDispatcher.Http("/FixAgentRatio", FixAgentRatio)

	JsDispatcher.Http("/DownAgentLeave", DownAgentLeave)
	JsDispatcher.Http("/RemoveAllAgent", RemoveAllAgent)
	JsDispatcher.Http("/ChagngeAgentArticle", ChagngeAgentArticle)
	///

	//

	//

	//ErrorLog("CFG:%v", JsConfig.CFG)
}

func ChangeAgentStatus(session *JsNet.StSession) {
	type st_get struct {
		Status string
		UID    string
	}
	st := &st_get{}
	session.GetPara(st)
	user := &common.ST_User{}
	if err := Update(constant.Hash_User, st.UID, user, func() {
		if user.Agent != nil {
			// if user.Agent.Current.OpreatStatus == "" {
			user.Agent.Current.OpreatStatus = st.Status
			// }
		}
	}); err != nil {

	}
	Forward(session, "0", nil)
}

func ChagngeAgentArticle(session *JsNet.StSession) {
	type st_get struct {
		Index int
	}
	st := &st_get{}
	session.GetPara(st)
	// for i := 100001; i < st.Index; i++ {
	// 	UID := "user-" + strconv.Itoa(i)
	// 	user1 := &common.ST_User{}
	// 	user2 := &common.ST_User2{}
	// 	if err := ShareLock(constant.Hash_User, UID, user2); err != nil {
	// 		Info("111111111")
	// 		continue
	// 	}

	// 	b, err := json.Marshal(user2)
	// 	if err != nil {
	// 		Info("2222222222222")
	// 		continue
	// 	}

	// 	if err = json.Unmarshal(b, user1); err != nil {
	// 		Info("3333333333333333333")
	// 		continue
	// 	}
	// 	if user1.Agent != nil && user2.Agent != nil {
	// 		user1.Agent.Article = user2.Agent.ShareArticle
	// 		DirectWrite(constant.Hash_User, UID, user1)
	// 	}
	// }
	Forward(session, "0", nil)
}

func RemoveAllAgent(session *JsNet.StSession) {
	type st_get struct {
		Index int
	}
	st := &st_get{}
	session.GetPara(st)

	data := &ST_FullCache{}
	if err := WriteLock(constant.Hash_Global, constant.KEY_Global_Agent, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data = &ST_FullCache{}
	if err := WriteBack(constant.Hash_Global, constant.KEY_Global_Agent, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	data_E := &ST_FullCache{}
	if err := WriteLock(constant.Hash_Global, constant.KEY_Global_Employee, data_E); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data_E = &ST_FullCache{}
	if err := WriteBack(constant.Hash_Global, constant.KEY_Global_Employee, data_E); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	cityAgent := &TypListCache{}
	if err := WriteLock(constant.Hash_Global, constant.KEY_City_Agent, cityAgent); err == nil {
		cityAgent = &TypListCache{}
		WriteBack(constant.Hash_Global, constant.KEY_City_Agent, cityAgent)
	}

	cityemployee := &TypListCache{}
	if err := WriteLock(constant.Hash_Global, constant.KEY_City_Employee, cityemployee); err == nil {
		cityemployee = &TypListCache{}
		WriteBack(constant.Hash_Global, constant.KEY_City_Employee, cityemployee)
	}

	for i := 100001; i < st.Index; i++ {
		UID := "user-" + strconv.Itoa(i)
		user := &common.ST_User{}
		if err := Update(constant.Hash_User, UID, user, func() {
			user.Agent = nil
			user.AgentInfo = nil
			user.SupAgentInfo = nil
		}); err != nil {
			continue
		}
	}

	Forward(session, "0", nil)
}

func DownAgentLeave(session *JsNet.StSession) {
	type st_get struct {
		UID string
	}
	st := &st_get{}
	session.GetPara(st)
	data := &ST_FullCache{}
	if err := WriteLock(constant.Hash_Global, constant.KEY_Global_Agent, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.RemoveFromFull(st.UID)
	if err := WriteBack(constant.Hash_Global, constant.KEY_Global_Agent, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	data_E := &ST_FullCache{}
	if err := WriteLock(constant.Hash_Global, constant.KEY_Global_Employee, data_E); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data_E = &ST_FullCache{}
	if err := WriteBack(constant.Hash_Global, constant.KEY_Global_Employee, data_E); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	user := &common.ST_User{}
	if err := Update(constant.Hash_User, st.UID, user, func() {
		user.Agent = nil
		user.AgentInfo = nil
		user.SupAgentInfo = nil

	}); err != nil {
		ForwardEx(session, "1", nil, err.Error())
	}

	Forward(session, "0", nil)
}

func ChangeAllGlobalStatus(session *JsNet.StSession) {
	type st_get struct {
		Status string
	}
	st := &st_get{}
	session.GetPara(st)

	data := &ST_FullCache{}
	if err := ShareLock(constant.Hash_Global, constant.KEY_Global_Agent, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	for _, v := range data.Ids {
		user := &common.ST_User{}
		if err := Update(constant.Hash_User, v, user, func() {
			if user.Agent != nil {
				user.Agent.Current.OpreatStatus = st.Status
			}
		}); err != nil {
			continue
		}
	}
	Forward(session, "0", nil)
}

func FixAgentRatio(session *JsNet.StSession) {
	data := &ST_FullCache{}
	if err := ShareLock(constant.Hash_Global, constant.KEY_Global_Agent, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	for _, v := range data.Ids {
		user := &common.ST_User{}
		if err := Update(constant.Hash_User, v, user, func() {
			if user.Agent != nil && user.Agent.Current.OpreatStatus == constant.Agent_Online {

				if user.Agent.AgentLevel == constant.Agent_Level_Diamonds_A {
					user.Agent.FranchiseFee = JsConfig.CFG.Agent.Diamonds_A.FranchiseFee
					user.Agent.AgentRatio.CommonRate = JsConfig.CFG.Agent.Diamonds_A.CommonRate
					user.Agent.AgentRatio.SpecialRate = JsConfig.CFG.Agent.Diamonds_A.SpecialRate
					user.Agent.AgentRatio.CustomRate = JsConfig.CFG.Agent.Diamonds_A.CustomRate
					user.Agent.AgentRatio.RushRate = JsConfig.CFG.Agent.Diamonds_A.RushRate
				}
				if user.Agent.AgentLevel == constant.Agent_Level_TryUse {
					user.Agent.FranchiseFee = JsConfig.CFG.Agent.TryUse.FranchiseFee
					user.Agent.AgentRatio.CommonRate = JsConfig.CFG.Agent.TryUse.CommonRate
					user.Agent.AgentRatio.SpecialRate = JsConfig.CFG.Agent.TryUse.SpecialRate
					user.Agent.AgentRatio.CustomRate = JsConfig.CFG.Agent.TryUse.CustomRate
					user.Agent.AgentRatio.RushRate = JsConfig.CFG.Agent.TryUse.RushRate
				}
			}
		}); err != nil {
			continue
		}
	}
	Forward(session, "0", nil)
}

//提现代理保证金
func WithdrawBond(session *JsNet.StSession) {
	type st_get struct {
		UID string //
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	user, err := common.WithdrawBond(st.UID, make(map[string]string))
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", user)
}

///提交保证金
func OrderSubmitBond(session *JsNet.StSession) {
	type st_Get struct {
		UID         string
		PlatformFee int
		Bond        int
		AgentLevel  string
	}
	st := &st_Get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" || st.PlatformFee < 0 || st.Bond < 0 || st.AgentLevel == "" {
		ForwardEx(session, "1", nil, "OrderSubmitBond failed,UID=%s,PlatformFee=%d,Bond=%d,AgentLevel=%s\n",
			st.UID, st.PlatformFee, st.Bond, st.AgentLevel)
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
	order.Bond = st.Bond
	order.PlatformFee = st.PlatformFee
	order.AgentLevel = st.AgentLevel

	if err := common.SubmitBondOrder(order); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", order)
}

func OrderPayBond(session *JsNet.StSession) {
	type st_get struct {
		OrderID string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.OrderID == "" {
		ForwardEx(session, "1", nil, "OrderPayBond failed,UID=%s\n", st.OrderID)
		return
	}

	user, err := common.OrderPayBond(st.OrderID, nil)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", user)
}

func WithdrawBalance(session *JsNet.StSession) {
	type st_get struct {
		UID   string
		Money int
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.Money <= 0 || st.UID == "" {
		ForwardEx(session, "1", nil, "WithdrawBalance failed,UID=%S,Money=%d\n", st.UID, st.Money)
		return
	}
	data, err := common.WithdrawBalance(st.UID, st.Money)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

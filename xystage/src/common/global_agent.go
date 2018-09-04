package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	. "cache/cacheLib"
	"constant"
	. "util"
)

type ST_AgentStatusID struct {
	Apply          []string //申请中
	WaitPay        []string //待支付
	NoPass         []string //审核不通过
	Online         []string //在线的
	OfflineOnForce []string //解约的
	OfflineSelf    []string //解约的
}

//新增加一个全局代理
func NewGlobalAgent(City, UID string) error {
	Info("NewGlobalAgent....\n")
	if City == "" || UID == "" {
		return ErrorLog("NewGlobalAgent failed,param is empty,City=%s,UID=%s\n", City, UID)
	}
	data := &ST_FullCache{}
	err := WriteLock(constant.Hash_Global, constant.KEY_Global_Agent, data)
	data.AddToFull(UID)
	if err != nil {
		return DirectWrite(constant.Hash_Global, constant.KEY_Global_Agent, data)
	} else {
		if e := WriteBack(constant.Hash_Global, constant.KEY_Global_Agent, data); e != nil {
			return e
		}
		go addCityAgent(City, UID)
	}
	return nil
}

//新增加一个城市代理
func addCityAgent(City, agentID string) error {
	Info("addCityAgent....\n")
	cityAgent := &TypListCache{}
	if err := WriteLock(constant.Hash_Global, constant.KEY_City_Agent, cityAgent); err != nil {
		cityAgent.NewTypeList(City, agentID)
		return DirectWrite(constant.Hash_Global, constant.KEY_City_Agent, cityAgent)
	}
	cityAgent.NewTypeList(City, agentID)
	return WriteBack(constant.Hash_Global, constant.KEY_City_Agent, cityAgent)
}

//获取城市代理人
func getCityAgent(City string) ([]*ST_User, error) {
	cityAgent := &TypListCache{}
	if err := ShareLock(constant.Hash_Global, constant.KEY_City_Agent, cityAgent); err != nil {
		return []*ST_User{}, err
	}
	if list, ok := cityAgent.Data[City]; ok {
		return GetMoreUserInfo(list), nil
	}
	return []*ST_User{}, ErrorLog("城市:%s没有城市代理人\n", City)
}

///获取所有的代理状态
func GetGlobalAgent() (*ST_FullCache, error) {
	data := &ST_FullCache{}
	if err := ShareLock(constant.Hash_Global, constant.KEY_Global_Agent, data); err != nil {
		return data, ErrorLog(err.Error())
	}
	return data, nil
}

//获取城市代理人
func GetCityAgent(session *JsNet.StSession) {
	type st_get struct {
		City string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.City == "" {
		ForwardEx(session, "1", nil, "GetCityAgent failed,City is empty\n")
		return
	}
	data, err := getCityAgent(st.City)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func GetGlobalAgentID() (st *ST_AgentStatusID) {
	data, err := GetGlobalAgent()
	if err != nil {
		return nil
	}
	info := &ST_AgentStatusID{}

	UserList := GetMoreUserInfo(data.Ids)
	for _, v := range UserList {
		if v.Agent != nil {
			if v.Agent.Current.OpreatStatus == constant.Agent_Apply ||
				v.Agent.Current.OpreatStatus == constant.Agent_ReApply {
				info.Apply = append(info.Apply, v.UID)
			} else if v.Agent.Current.OpreatStatus == constant.Agent_PassReviewe {
				info.WaitPay = append(info.WaitPay, v.UID)
			} else if v.Agent.Current.OpreatStatus == constant.Agent_NoPassReviewe {
				info.NoPass = append(info.NoPass, v.UID)
			} else if v.Agent.Current.OpreatStatus == constant.Agent_Offline_force {
				info.OfflineOnForce = append(info.OfflineOnForce, v.UID)
			} else if v.Agent.Current.OpreatStatus == constant.Agent_Offline_self {
				info.OfflineSelf = append(info.OfflineSelf, v.UID)
			} else if v.Agent.Current.OpreatStatus == constant.Agent_Online {
				info.Online = append(info.Online, v.UID)
			} else {
				Error("v.Agent.Current.OpreatStatus=%s\n", v.Agent.Current.OpreatStatus)
				continue
			}
		}
	}
	return info
}

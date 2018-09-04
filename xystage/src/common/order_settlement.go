package common

import (
	. "JsLib/JsLogger"
	"constant"
)

///计算各种价格
func orderSettlement(data *ST_Order) error {
	orderHosSettlement(data)
	return orderAgentSettlement(data)
}

func testLog(data *ST_Order, info string) {
	ErrorLog("XingYaoPrice=%d\n", data.XingYaoPrice)
	ErrorLog("ProDeposit=%d\n", data.ProDeposit)
	ErrorLog("CouponPrice=%d\n", data.CouponPrice)
	ErrorLog("RedPrice=%d\n", data.RedPrice)
	ErrorLog("HosSettlementRatio=%d\n", data.HosSettlementRatio)
	ErrorLog("AppendPrice=%d\n", data.AppendPrice)
	ErrorLog("HosPayPrice=%d\n", data.HosPayPrice)
	ErrorLog("HosPayRealPrice=%d\n", data.HosPayRealPrice)
	ErrorLog("TotalPrice=%d\n", data.TotalPrice)
	ErrorLog("HosSettlementPrice=%d\n", data.HosSettlementPrice)
	ErrorLog("FullPayReturnPrice=%d\n", data.FullPayReturnPrice)
	ErrorLog("RealSettlementPrice=%d\n", data.RealSettlementPrice)
	ErrorLog(info + "\n")
}

func orderHosSettlement(data *ST_Order) {
	if data == nil {
		Error("orderHosSettlement failed,order is nil \n")
	}

	data.HosPayPrice = data.XingYaoPrice - data.ProDeposit - data.CouponPrice - data.RedPrice //到院支付价格

	data.HosPayRealPrice = data.HosPayPrice + data.AppendPrice //到院实际支付总额

	data.TotalPrice = data.XingYaoPrice + data.AppendPrice //消费总金额

	data.HosSettlementPrice = (data.TotalPrice-data.CouponPrice)*data.HosSettlementRatio/100 - data.ProDeposit //结算价格

	data.FullPayReturnPrice = data.TotalPrice - (data.TotalPrice-data.CouponPrice)*data.HosSettlementRatio/100 + data.RedPrice //全款支付返还医院金额

	data.RealSettlementPrice = data.HosSettlementPrice - data.RedPrice

}

func orderAgentSettlement(data *ST_Order) error {
	if data == nil {
		return ErrorLog("orderAgentSettlement failed,order  is nil \n")
	}
	if data.AgentInfo == nil {
		Error("orderAgentSettlement  ignore ,data.AgentInfo  is nil \n")
		return nil
	}

	if data.AgentInfo.UID != "" {
		agent, err := GetUserInfo(data.AgentInfo.UID)
		if err != nil {
			return err
		}
		if agent.Agent.Current.OpreatStatus != constant.Agent_Online {
			return ErrorLog("当前小B,不是在线状态,UID:%s,Status=%s \n", agent.UID, agent.Agent.Current.OpreatStatus)
		}
		data.AgentCommission = data.TotalPrice * data.AgentInfo.ProRatio * data.HosSettlementRatio / 10000 //结算比例
	}
	return nil
}

package main

import (
	"common"
	"constant"
	"sort"
)

type ST_GlobalOrder struct {
	Check      []string //待确认
	Verify     []string //待校验(已评论)
	Settlement []string //待结算
	Complete   []string //已完成
	Cancle     []string //已取消
}

func GetGlobalOrder() *ST_GlobalOrder {
	data := ST_GlobalOrder{}
	ids := common.GetGlobalOrderList()
	var list common.OrderList
	list = common.QueryMoreOrders(ids)
	if len(list) == 0 {
		return &data
	}
	sort.Sort(list)

	for _, v := range list {
		if v.Current.OpreatStatus == constant.Status_Order_PenddingAppointment {
			data.Check = append(data.Check, v.OrderID)
		}
		if v.Current.OpreatStatus == constant.Status_Order_PenddingVerify {
			data.Verify = append(data.Verify, v.OrderID)
		}
		if v.Current.OpreatStatus == constant.Status_Order_PendingStatements {
			data.Settlement = append(data.Settlement, v.OrderID)
		}
		if v.Current.OpreatUserStatus == constant.Status_Order_Succeed {
			data.Complete = append(data.Complete, v.OrderID)
		}
		if v.Current.OpreatStatus == constant.Status_Order_CancleBeforeVerfy ||
			v.Current.OpreatStatus == constant.Status_Order_CancleAfterVerfy {
			data.Cancle = append(data.Cancle, v.OrderID)
		}
	}

	return &data
}

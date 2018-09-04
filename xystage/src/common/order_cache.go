package common

import (
	. "JsLib/JsLogger"
	"constant"
	. "util"
)

type ST_HosOrder struct {
	Check   []*ST_Order
	Success []*ST_Order
}

//////获取医院端用的医院订单
func GetHosOrder(HosID string) (*ST_HosOrder, error) {
	data := []string{}
	if err := ShareLock(constant.Hash_HospitalOrder, HosID, &data); err != nil {
		return nil, err
	}
	order := QueryMoreOrders(data)
	list := ST_HosOrder{}

	for _, v := range order {
		if v.Current.OpreatStatus == constant.Status_Order_PendingConfirm {
			list.Check = append(list.Check, v)
		}
		if v.Current.OpreatStatus == constant.Status_Order_PendingStatements ||
			v.Current.OpreatStatus == constant.Status_Order_PenddingReconcile ||
			v.Current.OpreatStatus == constant.Status_OrderPenddingCollection ||
			v.Current.OpreatStatus == constant.Status_OrderSysConfirmCollection {
			list.Success = append(list.Success, v)
		}

	}
	return &list, nil
}

///往医院添加一条订单
func AppendHosOrder(HosID, OrderID string) error {
	data := &[]string{}
	err := WriteLock(constant.Hash_HospitalOrder, HosID, data)
	AppendUniqueString(data, OrderID)
	if err != nil {
		return DirectWrite(constant.Hash_HospitalOrder, HosID, data)
	}
	return WriteBack(constant.Hash_HospitalOrder, HosID, data)
}

///添加一个订单到全局
func addOrderToGlobal(OrderID string) error {
	data := []string{}
	if err := WriteLock(constant.Hash_Global, constant.KEY_GlobalOrderList, &data); err != nil {
		data = append(data, OrderID)
		return DirectWrite(constant.Hash_Global, constant.KEY_GlobalOrderList, &data)
	}
	exist := false
	for _, v := range data {
		if v == OrderID {
			exist = true
			break
		}
	}
	if !exist {
		data = append(data, OrderID)
	}
	return WriteBack(constant.Hash_Global, constant.KEY_GlobalOrderList, &data)
}

func GetGlobalOrderList() []string {
	data := []string{}
	if err := ShareLock(constant.Hash_Global, constant.KEY_GlobalOrderList, &data); err != nil {
		Error(err.Error())
	}
	return data
}

///追加一个订到到产品
func appendOrderToProduct(ProID, OrderID string) error {
	if ProID == "" || OrderID == "" {
		return ErrorLog("addpendOrderToProduct,ProID=%s,OrderID=%s\n", ProID, OrderID)
	}
	data := []string{}
	if err := WriteLock(constant.Has_ProOrderCache, ProID, &data); err != nil {
		data = append(data, OrderID)
		return DirectWrite(constant.Has_ProOrderCache, ProID, &data)
	}
	exist := false
	for _, v := range data {
		if v == OrderID {
			exist = true
			break
		}
	}
	if !exist {
		data = append(data, OrderID)
	}
	return WriteBack(constant.Has_ProOrderCache, ProID, &data)
}

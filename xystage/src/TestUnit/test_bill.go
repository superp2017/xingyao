package main

import (
	"JsLib/JsConfig"
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"common"
	"strconv"

	"time"
	. "util"
)

func init_bill() {
	JsDispatcher.Http("/startsettlement", StartSettlement) //开启结算
	/////////////////////////////////系统账单//////////////////////////////////////////
	JsDispatcher.Http("/querysystembill", common.QuerySystemBill) //查询系统账单
	JsDispatcher.Http("/gethosbillinfo", common.GetHosBillInfo)   //获取医院账单信息
	JsDispatcher.Http("/gethosallbills", common.GetHosAllBills)   //获取医院所有账单
	JsDispatcher.Http("/gethosmonthbill", common.GetHosMonthBill) //获取医院某个月的账单

}

////开始结算
func StartSettlement(session *JsNet.StSession) {

	checkSettlement()
	Forward(session, "0", nil)
}

////启动定时任务
func StartTimer() {
	go timeTask(testTask)
}

////定时任务
func timeTask(task func()) {
	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 1, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//////////定时任务//////////////////
		task()
	}
}

func testTask() {
	Error("当前时间:%s\n", CurTime())
}

///每个月的结算
func checkSettlement() {
	if strconv.Itoa(CurDay()) == JsConfig.CFG.BillDate.Date {
		hosList := common.GlobalHospitalList()
		if len(hosList) == 0 {
			Error("checkSettlement ,hoslist len is 0\n")
			return
		}
		bills := []*common.ST_HosBill{}
		for _, id := range hosList {
			if b, err := settlemntHospital(id); err == nil {
				if b != nil {
					bills = append(bills, b)
				}
			}
		}
		if err := common.UpdateHosMonthBill(bills); err != nil {
			Error("common.UpdateHosMonthBill failed,err:=%s\n", err.Error())
		}
	}
}

//单个医院结算
func settlemntHospital(HosID string) (*common.ST_HosBill, error) {
	ordercache, err := common.GetHosBillOrder(HosID)
	if err != nil {
		return nil, err
	}
	if ordercache == nil {
		return nil, nil
	}
	if ordercache != nil && len(ordercache.Refund) == 0 && len(ordercache.Settlement) == 0 {
		return nil, nil
	}
	///生成月账单
	bill, err := common.GenHosBill(HosID, ordercache)
	if err != nil {
		return nil, ErrorLog("settlemntHospital GenHosBill failed,HosID=%s\n", HosID)
	}
	return bill, nil
}

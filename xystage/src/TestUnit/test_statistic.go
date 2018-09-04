package main

import (
	"JsLib/JsDispatcher"
	"JsLib/JsNet"
	"common"
	"constant"
	"math/rand"
	"time"
	. "util"
)

type StHospitalStatistics struct {
	Date       string
	HosID      string
	VisitNum   int
	OrderNum   int
	TradeNum   int
	ProductNum int
	DoctorNum  int
}

type StTotalStatistics struct {
	Date        string
	VisitNum    int
	OrderNum    int
	TradeNum    int
	ProductNum  int
	HospitalNum int
	DoctorNum   int
}

func init_statistic() {
	JsDispatcher.Http("/RebuildHosStatistic", RebuildHosStatistic)
}

func RebuildHosStatistic(session *JsNet.StSession) {

	hosid := common.GlobalHospitalList()
	if len(hosid) == 0 {
		Forward(session, "0", nil)
	}

	glo := &StTotalStatistics{}
	for _, v := range hosid {
		data := []*StHospitalStatistics{}
		st := &StHospitalStatistics{}
		st.Date = time.Now().Format("2006-01-02")
		st.HosID = v
		st.VisitNum = rand.Intn(5000)

		order := []string{}
		if err := ShareLock(constant.Hash_HospitalOrder, v, &order); err == nil {
			list := common.QueryMoreOrders(order)
			for _, order := range list {
				st.TradeNum += order.XingYaoPrice
			}
		}

		st.OrderNum = len(order)
		doc := common.GetHosDoctorList(v)
		st.DoctorNum = len(doc)
		pro := common.GetHosProductList(v)
		st.ProductNum = len(pro)
		data = append(data, st)
		DirectWrite(constant.HOSPITAL_STATISTICS, v, &data)

		glo.Date = time.Now().Format("2006-01-02")
		glo.VisitNum += st.VisitNum
		glo.OrderNum += st.OrderNum
		glo.DoctorNum += st.DoctorNum
		glo.ProductNum += st.ProductNum
	}
	Set(constant.TOTAL_STATISTICS, glo)
	Forward(session, "0", nil)

}

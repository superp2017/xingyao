package main
//
//import (
//	. "JsLib/JsLogger"
//	"JsLib/JsNet"
//	"constant"
//	"sync"
//	"time"
//)
//
//type StTotalStatistics struct {
//	Date        string
//	VisitNum    int
//	OrderNum    int
//	TradeNum    int
//	ProductNum  int
//	HospitalNum int
//	DoctorNum   int
//}
//
//// type StHospitalTotal struct {
//// 	HosID      string
//// 	VisitNum   int
//// 	OrderNum   int
//// 	TradeNum   int
//// 	ProductNum int
//// 	DoctorNum  int
//// }
//
//type StHospitalStatistics struct {
//	Date       string
//	HosID      string
//	VisitNum   int
//	OrderNum   int
//	TradeNum   int
//	ProductNum int
//	DoctorNum  int
//}
//
//var g_total StTotalStatistics
//var g_map_sta map[string][]*StHospitalStatistics
//var g_date string
//var g_mutex sync.Mutex
//
//func init_statistics() {
//
//	//JsDispatcher.Http("/queryhospstatistics", QueryHospVisit)            //查询医院的各种统计信息
//	//JsDispatcher.Http("/addhospitalvisit", AddHospitalVisit)             //增加医院访问量
//	//JsDispatcher.Http("/addhosorderstatistic", AddHosOrderStatistic)     //增加医院的订单统计
//	//JsDispatcher.Http("/addhosproductstatistic", AddHosProductStatistic) //增加医院产品的统计
//	//JsDispatcher.Http("/adddoctorstatistic", AddDoctorStatistic)         //增加医生统计
//	//JsDispatcher.Http("/addhospitalstatistic", AddHospitalStatistic)     //增加医院统计
//
//	//e := Get(constant.TOTAL_STATISTICS, &g_total)
//	//if e != nil {
//	//	log.Fatalln(e.Error())
//	//}
//	//
//	//g_map_sta = make(map[string][]*StHospitalStatistics)
//	//
//	//hosid := common.GlobalHospitalList()
//	//if len(hosid) > 0 {
//	//	for _, v := range hosid {
//	//		data := []*StHospitalStatistics{}
//	//
//	//		if err := ShareLock(constant.HOSPITAL_STATISTICS, v, &data); err == nil {
//	//			g_map_sta[v] = data
//	//		}
//	//	}
//	//}
//}
//
//func up_date() {
//	g_date = time.Now().Format("2006-01-02")
//}
//
//func coolie() {
//	for {
//		up_date()
//		g_mutex.Lock()
//		Set(constant.TOTAL_STATISTICS, &g_total)
//		for k, v := range g_map_sta {
//			DirectWrite(constant.HOSPITAL_STATISTICS, k, &v)
//		}
//		g_mutex.Unlock()
//		time.Sleep(time.Minute * 5)
//	}
//}
//
////////////////////////////////////////////////////////////////////////////////////
//
//func AddHospitalVisit(session *JsNet.StSession) {
//	type st_get struct {
//		HosID string
//	}
//	st := &st_get{}
//	if err := session.GetPara(st); err != nil {
//		ForwardEx(session, "1", nil, err.Error())
//		return
//	}
//	if st.HosID == "" {
//		ForwardEx(session, "1", nil, "HosID is empty\n")
//		return
//	}
//
//	g_mutex.Lock()
//	defer g_mutex.Unlock()
//
//	g_total.VisitNum++
//
//	hospList, ok := g_map_sta[st.HosID]
//	if !ok {
//		hospList = make([]*StHospitalStatistics, 0)
//	}
//
//	hosp := &StHospitalStatistics{}
//	Inx := -1
//	for k, v := range hospList {
//		if g_date == v.Date {
//			hosp = v
//			Inx = k
//			break
//		}
//	}
//
//	hosp.VisitNum++
//	if Inx != -1 {
//		hospList[Inx] = hosp
//	} else {
//		hospList = append(hospList, hosp)
//	}
//
//	g_map_sta[st.HosID] = hospList
//
//	Forward(session, "0", nil)
//}
//
//func AddHosOrderStatistic(session *JsNet.StSession) {
//
//	type st_get struct {
//		HosID string
//		Fee   int
//	}
//	st := &st_get{}
//	if err := session.GetPara(st); err != nil {
//		ForwardEx(session, "1", nil, err.Error())
//		return
//	}
//	if st.HosID == "" || st.Fee < 0 {
//		ForwardEx(session, "1", nil, "HosID is empty or Fee <0\n")
//		return
//	}
//
//	g_mutex.Lock()
//	defer g_mutex.Unlock()
//
//	g_total.OrderNum++
//	g_total.TradeNum += st.Fee
//
//	hospList, ok := g_map_sta[st.HosID]
//	if !ok {
//		hospList = make([]*StHospitalStatistics, 0)
//	}
//
//	hosp := &StHospitalStatistics{}
//	Inx := -1
//	for k, v := range hospList {
//		if g_date == v.Date {
//			hosp = v
//			Inx = k
//			break
//		}
//	}
//
//	hosp.OrderNum++
//	hosp.OrderNum += st.Fee
//	if Inx != -1 {
//		hospList[Inx] = hosp
//	} else {
//		hospList = append(hospList, hosp)
//	}
//
//	g_map_sta[st.HosID] = hospList
//
//	Forward(session, "0", nil)
//}
//
//func AddHosProductStatistic(session *JsNet.StSession) {
//	type st_get struct {
//		HosID string
//	}
//	st := &st_get{}
//	if err := session.GetPara(st); err != nil {
//		ForwardEx(session, "1", nil, err.Error())
//		return
//	}
//	if st.HosID == "" {
//		ForwardEx(session, "1", nil, "HosID is empty \n")
//		return
//	}
//
//	g_mutex.Lock()
//	defer g_mutex.Unlock()
//
//	g_total.ProductNum++
//
//	hospList, ok := g_map_sta[st.HosID]
//	if !ok {
//		hospList = make([]*StHospitalStatistics, 0)
//	}
//
//	hosp := &StHospitalStatistics{}
//	Inx := -1
//	for k, v := range hospList {
//		if g_date == v.Date {
//			hosp = v
//			Inx = k
//			break
//		}
//	}
//
//	hosp.ProductNum++
//	if Inx != -1 {
//		hospList[Inx] = hosp
//	} else {
//		hospList = append(hospList, hosp)
//	}
//
//	g_map_sta[st.HosID] = hospList
//
//	Forward(session, "0", nil)
//}
//
//func AddDoctorStatistic(session *JsNet.StSession) {
//	type st_get struct {
//		HosID string
//	}
//	st := &st_get{}
//	if err := session.GetPara(st); err != nil {
//		ForwardEx(session, "1", nil, err.Error())
//		return
//	}
//	if st.HosID == "" {
//		ForwardEx(session, "1", nil, "HosID is empty \n")
//		return
//	}
//
//	g_mutex.Lock()
//	defer g_mutex.Unlock()
//
//	g_total.DoctorNum++
//
//	hospList, ok := g_map_sta[st.HosID]
//	if !ok {
//		hospList = make([]*StHospitalStatistics, 0)
//	}
//
//	hosp := &StHospitalStatistics{}
//	Inx := -1
//	for k, v := range hospList {
//		if g_date == v.Date {
//			hosp = v
//			Inx = k
//			break
//		}
//	}
//
//	hosp.DoctorNum++
//	if Inx != -1 {
//		hospList[Inx] = hosp
//	} else {
//		hospList = append(hospList, hosp)
//	}
//
//	g_map_sta[st.HosID] = hospList
//
//	Forward(session, "0", nil)
//}
//
//func AddHospitalStatistic(session *JsNet.StSession) {
//	g_mutex.Lock()
//	defer g_mutex.Unlock()
//
//	g_total.HospitalNum++
//
//}
//
////////////////////////////////////////////////////////////////////////////////////
//
//func QueryHospVisit(session *JsNet.StSession) {
//	type Para struct {
//		HosID string
//	}
//
//	type Ret struct {
//		Ret  string
//		Msg  string
//		Data []*StHospitalStatistics
//	}
//
//	para := &Para{}
//	ret := &Ret{}
//
//	e := session.GetPara(para)
//	if e != nil {
//		Error(e.Error())
//		ret.Ret = "1"
//		ret.Msg = e.Error()
//		session.Forward(ret)
//		return
//	}
//
//	g_mutex.Lock()
//	hospList, ok := g_map_sta[para.HosID]
//	g_mutex.Unlock()
//
//	if ok {
//		ret.Ret = "0"
//		ret.Msg = "success"
//		ret.Data = hospList
//		session.Forward(ret)
//		return
//	} else {
//		hList := make([]*StHospitalStatistics, 0)
//		e := ShareLock(constant.HOSPITAL_STATISTICS, para.HosID, &hList)
//		Info("StHospitalStatistics  =  %v\n", hList)
//		if e != nil {
//			Error(e.Error())
//			ret.Ret = "0"
//			ret.Msg = e.Error()
//			ret.Data = []*StHospitalStatistics{}
//			session.Forward(ret)
//			return
//		}
//
//		ret.Ret = "0"
//		ret.Msg = "success"
//		ret.Data = hList
//		session.Forward(ret)
//	}
//}

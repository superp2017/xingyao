package main

import (
	. "JsLib/JsLogger"
	. "cache/cacheLib"
	"common"
	"constant"
	// "log"
	// "sort"
	"strconv"
	"sync"
)

//全局的医院状态列表
var gl_HosONUM *ST_ONUMCache = &ST_ONUMCache{}
var gl_HosONUM_mutex sync.Mutex //全局医院锁

///全局的分城市的医院列表
var gl_CityHos *TypListCache = &TypListCache{}
var gl_CityHos_mutex sync.Mutex ///城市医院锁

///分城市分一级项目的医院列表
var gl_CityItemHos map[string]*TypListCache = make(map[string]*TypListCache)
var gl_CityItemHos_mutex sync.Mutex //分城市项目的医院锁

///////////////////////////////////////////////////////////////////////////////////
///////////////////////////全局医院获取函数////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////

//全局医院的句柄
func (cache *CacheHandle) CallGetGlobalHos(para *ST_CallPara, ret *ST_ONUMCache) error {
	Info("CallGetGlobalHos...\n")
	if gl_HosONUM == nil {
		return ErrorLog("gl_HosONUM is nil...\n")
	}
	data := &ST_ONUMCache{}
	if para.SortType == "" {
		data = gl_HosONUM
	} else {
		data = SortONUM(para.SortType, 0)
	}
	ret.New = data.New
	ret.Modify = data.Modify
	ret.UnPass = data.UnPass
	ret.OffLine = data.OffLine
	ret.OnLine = data.OnLine
	ret.Full = data.Full

	return nil
}

//获取全局的分城市医院
func (cache *CacheHandle) CallGetGlobalCityHos(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallGetGlobalCityHos ...\n")
	if para.City == "" {
		return ErrorLog("CallGetGlobalCityHos City is empty ...\n")
	}
	if data, ok := gl_CityHos.Data[para.City]; ok {
		ret.Ids = data
		return nil
	}
	return nil
}

//获取全局的分城市分项目医院
func (cache *CacheHandle) CallGetCityItemHos(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallGetCityItemHos ...\n")
	if para.City == "" || para.FirstItem == "" {
		return ErrorLog("CallGetCityItemHos City is empty ...\n")
	}
	if data, ok := gl_CityItemHos[para.City]; ok {
		if ls, ex := data.Data[para.FirstItem]; ex {
			ret.Ids = ls
			return nil
		}
		return nil
	}
	return nil
}

func (cache *CacheHandle) CallUpdateHosSale(para *ST_SortPara, ret *ST_CallRet) error {
	Info("CallUpdateHosSale ...\n")
	if para.ID == "" || para.Num < 1 {
		return ErrorLog("CallUpdateHosSale param empty,ID=%s,Num=%d\n", para.ID, para.Num)
	}
	addToSaleHos(para.ID, para.Num, true)
	return nil
}

func (cache *CacheHandle) CallUpdateHosComment(para *ST_SortPara, ret *ST_CallRet) error {
	Info("CallUpdateHosComment ...\n")
	if para.ID == "" || para.Num < 1 {
		return ErrorLog("CallUpdateHosComment param empty,ID=%s,Num=%d\n", para.ID, para.Num)
	}
	addToCommentHos(para.ID, para.Num, true)
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////全局医院的更新函数/////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//新增加一个医院
func (cache *CacheHandle) CallAddNewHospital(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallAddNewHospital ...\n")
	if para.HosID == "" {
		return ErrorLog("CallGetCityItemHos HosID is empty ...\n")
	}

	hos, err := common.GetHospitalInfo(para.HosID)
	if err != nil {
		return ErrorLog(err.Error())
	}
	load_hospital(hos)
	return nil
}

//更新所有医院的状态
func (cache *CacheHandle) CallUpdateHospitalStatus(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallUpdateHospital ...\n")
	update_hospital_status()
	Info(" Leave CallUpdateHospital ...\n")
	return nil
}

///更新所以医院的城市
func (cache *CacheHandle) CallUpdateHospitalCity(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallUpdateHospitalCity ...\n")
	update_hospital_city()
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////全局医院的内部函数/////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///加载某一医院
func load_hospital(hos *common.ST_Hospital) {
	if hos.Current.OpreatStatus != constant.OperatingStatus_Del {
		loadHosStatusList(hos)
		load_CityHos(hos)
		addToSaleHos(hos.HosID, hos.Reservepeople, false)
		addToCommentHos(hos.HosID, hos.Evaluatepeople, false)
		update_Search(hos.HosID, hos.HosName)
	}
}

//加载分城市的医院缓存
func load_CityHos(hos *common.ST_Hospital) {
	gl_CityHos_mutex.Lock()
	gl_CityHos.NewTypeList(hos.HosCity, hos.HosID)
	gl_CityHos.NewTypeList(constant.All, hos.HosID)
	gl_CityHos_mutex.Unlock()
}

//加载医院的状态缓存
func loadHosStatusList(hos *common.ST_Hospital) {
	Info("loadHosStatusList ....\n")
	gl_HosONUM_mutex.Lock()

	gl_HosONUM.Full = append(gl_HosONUM.Full, hos.HosID)
	if hos.Current.OpreatStatus == constant.OperatingStatus_new {
		gl_HosONUM.New = append(gl_HosONUM.New, hos.HosID)
	}
	if hos.Current.OpreatStatus == constant.OperatingStatus_online {
		gl_HosONUM.OnLine = append(gl_HosONUM.OnLine, hos.HosID)
	}
	if hos.Current.OpreatStatus == constant.OperatingStatus_Offline_self ||
		hos.Current.OpreatStatus == constant.OperatingStatus_Offline_onforce {
		gl_HosONUM.OffLine = append(gl_HosONUM.OffLine, hos.HosID)
	}
	if hos.Current.OpreatStatus == constant.OperatingStatus_modify {
		gl_HosONUM.Modify = append(gl_HosONUM.Modify, hos.HosID)
		if hos.LastOnlineDate != "" {
			gl_HosONUM.OnLine = append(gl_HosONUM.OnLine, hos.HosID)
		}
	}
	if hos.Current.OpreatStatus == constant.OperatingStatus_Reviewer_NotPass {
		gl_HosONUM.UnPass = append(gl_HosONUM.UnPass, hos.HosID)
		if hos.LastOnlineDate != "" {
			gl_HosONUM.OnLine = append(gl_HosONUM.OnLine, hos.HosID)
		}
	}

	gl_HosONUM_mutex.Unlock()
}

//清除所有的医院缓存
func clearHoscache() {
	gl_HosONUM_mutex.Lock()
	gl_HosONUM.Clear()
	gl_HosONUM_mutex.Unlock()

	gl_CityHos_mutex.Lock()
	gl_CityHos.Clear()
	gl_CityHos_mutex.Unlock()
}

//更新医院的状态缓存
func update_hospital_status() {

	gl_HosONUM_mutex.Lock()
	gl_HosONUM.Clear()
	gl_HosONUM_mutex.Unlock()

	rangeHospital(func(hos *common.ST_Hospital) {
		loadHosStatusList(hos)
		go checkHosStatus(hos, true)
	})
}

//更新分城市的医院缓存
func update_hospital_city() {

	gl_CityHos_mutex.Lock()
	gl_CityHos.Clear()
	gl_CityHos_mutex.Unlock()

	gl_CityPricePro_mutex.Lock()
	gl_CityPricePro = make(map[string]*TypListCache)
	gl_CityPricePro_mutex.Unlock()

	gl_CityItemHos_mutex.Lock()
	gl_CityItemHos = make(map[string]*TypListCache)
	gl_CityItemHos_mutex.Unlock()

	gl_CityItemDoc_mutex.Lock()
	gl_CityItemDoc = make(map[string]*TypListCache)
	gl_CityItemDoc_mutex.Unlock()

	gl_CityProduct_mutex.Lock()
	gl_CityProduct.Clear()
	gl_CityProduct_mutex.Unlock()

	gl_CityDoctor_mutex.Lock()
	gl_CityDoctor.Clear()
	gl_CityDoctor_mutex.Unlock()

	rangeHospital(func(hos *common.ST_Hospital) {
		if hos.Current.OpreatStatus != constant.OperatingStatus_Del {
			load_CityHos(hos)
			rangeProduct(hos.HosID, func(pro *common.ST_Product) {
				go common.ModifyProCity(pro.ProID, hos.HosCity)
				loadCityPricePro(strconv.Itoa(pro.XingYaoPrice), pro.ProID, hos.HosCity)
				loadProItemHosDoc(pro.Doctors, pro.HosID, pro.FirstItem, hos.HosCity)
				loadProCity(pro.ProID, hos.HosCity)
			})
			rangeDoctor(hos.HosID, func(doc *common.ST_Doctor) {
				go common.ModifyDocCity(doc.DocID, hos.HosCity)
				loadCityDoc(doc.DocID, hos.HosCity)
			})
		}
	})
}

//创建一个分城市分项目的医院缓存
func newCityItemHosCache(HosID, City, FirstItem string) {
	gl_CityItemHos_mutex.Lock()
	if _, ok := gl_CityItemHos[City]; !ok {
		gl_CityItemHos[City] = &TypListCache{}
	}
	gl_CityItemHos[City].NewTypeList(FirstItem, HosID)
	gl_CityItemHos_mutex.Unlock()
	return
}

//循环遍历所有的医院
func rangeHospital(F func(hos *common.ST_Hospital)) {
	hosid := common.GlobalHospitalList()
	if len(hosid) == 0 {
		Error("load hosid list failed\n\n")
		return
	}
	var hoslist common.HosList
	hoslist = common.QueryMoreHosInfo(hosid)

	if len(hoslist) == 0 {
		Warn("global hoslist len is 0\n")
		return
	}

	// sort.Sort(hoslist)

	for _, hos := range hoslist {
		F(hos)
	}
}

//检查医院的状态，如果发现下线的医院，去执行下线所有的对应的医生和产品
func checkHosStatus(hos *common.ST_Hospital, reload bool) {
	if hos.Current.OpreatStatus == constant.OperatingStatus_Offline_self ||
		hos.Current.OpreatStatus == constant.OperatingStatus_Offline_onforce ||
		hos.Current.OpreatStatus == constant.OperatingStatus_Del {
		if offlineAllDoctor(hos.HosID) && reload {
			update_doctor_status(hos.HosID)
		}
		if offlineAllProduct(hos.HosID) && reload {
			update_product_status(hos.HosID)
		}
	}
}

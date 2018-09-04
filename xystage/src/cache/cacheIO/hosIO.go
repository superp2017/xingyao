package cacheIO

import (
	. "JsLib/JsLogger"
	. "cache/cacheLib"
	//"constant"
)

/////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////医院获取接口////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////
//全局医院的句柄
func GetGlobalHosPtr(SortType string) (*ST_ONUMCache, error) {
	Info("GlabalGetHosPtr...\n")
	para := &ST_CallPara{SortType: SortType}
	ret := &ST_ONUMCache{}

	server := getCacheIO()
	if server == nil {
		return ret, ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallGetGlobalHos", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetGlobalHos error:%s", err.Error())
		return ret, err
	}
	return ret, nil
}

//全局分城市医院的句柄
func GetCityHos(City string) ([]string, error) {
	Info("GetCityHos...\n")
	if City == "" {
		return []string{}, ErrorLog("GetCityHos para is empty\n")
	}
	city := City
	if city == "全部" || city == "全国" {
		city = "All"
	}
	para := &ST_CallPara{City: city}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return nil, ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallGetGlobalCityHos", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetGlobalCityHos error:%s", err.Error())
		return ret.Ids, err
	}
	return ret.Ids, nil
}

//获取分城市，分项目的医院
func GetCityItemHos(City, FirstItem string) ([]string, error) {
	Info("GetCityItemHos...\n")
	if City == "" || FirstItem == "" {
		return []string{}, ErrorLog("GetCityItemHos para is empty\n")
	}
	city := City
	if city == "全部" || city == "全国" {
		city = "All"
	}

	para := &ST_CallPara{City: city, FirstItem: FirstItem}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return nil, ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallGetCityItemHos", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetCityItemHos error:%s", err.Error())
		return ret.Ids, err
	}
	return ret.Ids, nil
}

/////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////医院修改接口////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////

//添加一个医院到新增加的待审核缓存里面
func GlobalAddNewHospital(HosID string) error {
	Info("GlobalAddNewHospital...\n")
	if HosID == "" {
		return ErrorLog("GlobalAddNewHospital failed,HosID=%s\n", HosID)
	}

	para := &ST_CallPara{HosID: HosID}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallAddNewHospital", para, ret); err != nil {
		ErrorLog("CacheHandle.CallAddNewHospital error:%s", err.Error())
		return err
	}
	return nil
}

func GlobalUpdateHosStatus() error {
	para := &ST_CallPara{}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallUpdateHospitalStatus", para, ret); err != nil {
		ErrorLog("CacheHandle.CallUpdateHospitalStatus error:%s", err.Error())
		return err
	}
	return nil
}

func GlobalUpdateHosSale(HosID string, saleNum int) error {
	Info("GlobalUpdateHosSale...\n")
	if HosID == "" || saleNum < 1 {
		return ErrorLog("GlobalUpdateHosSale failed,HosID is empty,or SaleNum <1\n")
	}
	para := &ST_SortPara{ID: HosID, Num: saleNum}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallUpdateHosSale", para, ret); err != nil {
		return ErrorLog("CacheHandle.CallUpdateHosSale error:%s", err.Error())
	}
	return nil
}

func GlobalUpdateHosComment(HosID string, Num int) error {
	Info("GlobalUpdateHosComment...\n")
	if HosID == "" {
		return ErrorLog("GlobalUpdateHosComment failed,HosID is empty,or Num <1\n")
	}
	para := &ST_SortPara{ID: HosID, Num: Num}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallUpdateHosComment", para, ret); err != nil {
		ErrorLog("CacheHandle.CallUpdateHosComment error:%s", err.Error())
		return err
	}
	return nil
}

func GlobalUpdateHospitalCity() error {
	para := &ST_CallPara{}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallUpdateHospitalCity", para, ret); err != nil {
		ErrorLog("CacheHandle.CallUpdateHospitalCity error:%s", err.Error())
		return err
	}
	return nil
}

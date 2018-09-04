package cacheIO

import (
	. "JsLib/JsLogger"
	. "cache/cacheLib"
	//"constant"
)

////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////全局医生获取接口////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////
//全局医生的句柄
func GetGlobalDocPtr(SortType string) (*ST_ONUMCache, error) {
	Info("GlabalGetDocPtr...\n")
	para := &ST_CallPara{SortType: SortType}
	ret := &ST_ONUMCache{}
	server := getCacheIO()
	if server == nil {
		return ret, ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallGetGlobalDoc", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetGlobalDoc error:%s", err.Error())
		return ret, err
	}
	return ret, nil
}

//获取
func GetHosDocPtr(HosID string) (*ST_ONUMCache, error) {
	Info("GetHosDocPtr...\n")
	ret := &ST_ONUMCache{}

	if HosID == "" {
		return ret, ErrorLog("GetHosDocPtr failed,HosID=%s\n", HosID)
	}

	para := &ST_CallPara{HosID: HosID}

	server := getCacheIO()
	if server == nil {
		return ret, ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallGetHosDocPtr", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetHosDocPtr error:%s", err.Error())
		return ret, err
	}
	return ret, nil
}

//全局分城市医生
func GetCityDoc(City string) ([]string, error) {
	Info("GetCityDoc...\n")
	if City == "" {
		return []string{}, ErrorLog("GetCityDoc para is empty\n")
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

	if err := server.client.Call("CacheHandle.CallGetGlobalCityDoc", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetGlobalCityDoc error:%s", err.Error())
		return ret.Ids, err
	}
	return ret.Ids, nil
}

//获取分城市，分项目的医生
func GetCityItemDoc(City, FirstItem string) ([]string, error) {
	Info("GetCityItemDoc...\n")
	if City == "" || FirstItem == "" {
		return []string{}, ErrorLog("GetCityItemDoc para is empty\n")
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

	if err := server.client.Call("CacheHandle.CallGetCityItemDoc", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetCityItemDoc error:%s", err.Error())
		return ret.Ids, err
	}
	return ret.Ids, nil
}

//////////////////////////////////////////////////////////////////////////////////
////////////////////////////医生修改接口//////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////

func GolbalUpdateDoc(HosID string) error {
	Info("GolbalUpdateDoc...\n")
	if HosID == "" {
		return ErrorLog("GolbalUpdateDoc failed,HosID=%s\n", HosID)
	}

	para := &ST_CallPara{HosID: HosID}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallUpdateDoctor", para, ret); err != nil {
		ErrorLog("CacheHandle.CallUpdateDoctor error:%s", err.Error())
		return err
	}
	return nil
}

func GlobalUpdateDocSale(DocID string, saleNum int) error {
	Info("GlobalUpdateDocSale...\n")
	if DocID == "" || saleNum < 1 {
		return ErrorLog("GlobalUpdateDocSale failed,DocID is empty,or SaleNum <1\n")
	}
	para := &ST_SortPara{ID: DocID, Num: saleNum}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallUpdateDocSale", para, ret); err != nil {
		return ErrorLog("CacheHandle.CallUpdateDocSale error:%s", err.Error())
	}
	return nil
}

func GlobalUpdateDocComment(DocID string, Num int) error {
	Info("GlobalUpdateDocComment...\n")
	if DocID == "" {
		return ErrorLog("GlobalUpdateDocComment failed,DocID is empty,or Num <1\n")
	}
	para := &ST_SortPara{ID: DocID, Num: Num}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallUpdateDocComment", para, ret); err != nil {
		ErrorLog("CacheHandle.CallUpdateDocComment error:%s", err.Error())
		return err
	}
	return nil
}

//添加一个医生到新增加的待审核缓存里面
func GolbalAddNewDoctor(DocID, HosID, City, Name string) error {
	Info("GolbalAddNewDoctor...\n")
	if DocID == "" || HosID == "" || City == "" || Name == "" {
		return ErrorLog("GolbalAddNewDoctor failed,DocID=%s,HosID=%s,City=%s,Name=%s\n", DocID, HosID, City, Name)
	}

	para := &ST_CallPara{
		DocID:      DocID,
		HosID:      HosID,
		City:       City,
		ObjectName: Name,
	}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallAddNewDoctor", para, ret); err != nil {
		return ErrorLog("CacheHandle.CallAddNewDoctor error:%s", err.Error())
	}
	return nil
}

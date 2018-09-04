package cacheIO

import (
	. "JsLib/JsLogger"
	. "cache/cacheLib"
	//"constant"
)

////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////产品获取接口////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////

//全局产品的句柄
func GetGlobalProPtr(SortType string) (*ST_ONUMCache, error) {
	Info("GlabalGetProPtr...\n")
	para := &ST_CallPara{SortType: SortType}
	ret := &ST_ONUMCache{}

	server := getCacheIO()
	if server == nil {
		return ret, ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallGetGlobalPro", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetGlobalPro error:%s", err.Error())
		return ret, err
	}
	return ret, nil
}

//获取所有未通过审核的产品
func GetHosProPtr(HosID string) (*ST_ONUMCache, error) {
	Info("GetHosProPtr...\n")
	if HosID == "" {
		return nil, ErrorLog("GetHosProPtr failed,HosID=%s\n", HosID)
	}

	para := &ST_CallPara{HosID: HosID}
	ret := &ST_ONUMCache{}

	server := getCacheIO()
	if server == nil {
		return ret, ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallGetHosProPtr", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetHosProPtr error:%s", err.Error())
		return ret, err
	}
	return ret, nil
}

//获取一级菜单和二级菜单的映射
func GetFirstSecondItemMap() (map[string][]string, error) {
	Info("GetFirstSecondItemMap...\n")

	para := &ST_CallPara{}
	ret := &TypListCache{}

	server := getCacheIO()
	if server == nil {
		return ret.Data, ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallGetFirstMapSecondItem", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetFirstMapSecondItem error:%s", err.Error())
		return ret.Data, err
	}
	return ret.Data, nil
}

//获取全局分部位的产品列表
func GetFirstSecondItemPro(FirstItem string) ([]string, error) {
	Info("GetFirstSecondItemPro...\n")
	if FirstItem == "" {
		return []string{}, ErrorLog("GetFirstSecondItemPro faild,FirstItem=%s\n", FirstItem)
	}
	para := &ST_CallPara{FirstItem: FirstItem}
	ret := &ST_CallRet{}
	server := getCacheIO()
	if server == nil {
		return ret.Ids, ErrorLog("getCacheIO error \n")
	}
	if err := server.client.Call("CacheHandle.CallGetGlobalFirstItemPro", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetGlobalFirstItemPro error:%s", err.Error())
		return ret.Ids, err
	}
	return ret.Ids, nil
}

//获取全局分城市的产品列表
func GetGlobalCityPro(City string) ([]string, error) {
	Info("GetGlobalCityPro...\n")

	if City == "" {
		return []string{}, ErrorLog("GetGlobalCityPro faild,City=%s\n", City)
	}
	city := City
	if city == "全部" || city == "全国" {
		city = "All"
	}
	para := &ST_CallPara{City: city}
	ret := &ST_CallRet{}
	server := getCacheIO()
	if server == nil {
		return ret.Ids, ErrorLog("getCacheIO error \n")
	}
	if err := server.client.Call("CacheHandle.CallGetGlobalCityPro", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetGlobalCityPro error:%s", err.Error())
		return ret.Ids, err
	}
	return ret.Ids, nil
}

//获取分项目的产品列表
func GetItemProduct(FirstItem, SecondItem string) ([]string, error) {
	Info("GetItemProduct...\n")
	if FirstItem == "" {
		return []string{}, ErrorLog("GetItemProduct faild,FirstItem=%s,SecondItem=%s\n", FirstItem, SecondItem)
	}
	para := &ST_CallPara{FirstItem: FirstItem, SecondItem: SecondItem}
	ret := &ST_CallRet{}
	server := getCacheIO()
	if server == nil {
		return []string{}, ErrorLog("getCacheIO error \n")
	}
	if err := server.client.Call("CacheHandle.CallGetItemProduct", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetItemProduct error:%s", err.Error())
		return ret.Ids, err
	}
	return ret.Ids, nil
}

//获取分城市分价格的产品列表
func GetCityPriceProduct(City string) (*TypListCache, error) {
	Info("GetCityPriceProduct...\n")
	if City == "" {
		return nil, ErrorLog("GetCityPriceProduct faild,City=%s\n", City)
	}
	city := City
	if city == "全部" || city == "全国" {
		city = "All"
	}
	para := &ST_CallPara{City: city}
	ret := &TypListCache{}
	server := getCacheIO()
	if server == nil {
		return nil, ErrorLog("getCacheIO error \n")
	}
	if err := server.client.Call("CacheHandle.CallGetCityPriceProduct", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetCityPriceProduct error:%s", err.Error())
		return ret, err
	}
	return ret, nil
}

//获取所有分价格的产品列表
func GetPriceProduct() (*TypListCache, error) {
	Info("GetPriceProduct...\n")
	para := &ST_CallPara{}
	ret := &TypListCache{}
	server := getCacheIO()
	if server == nil {
		return nil, ErrorLog("getCacheIO error \n")
	}
	if err := server.client.Call("CacheHandle.CallGetPriceProduct", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetPriceProduct error:%s", err.Error())
		return ret, err
	}
	return ret, nil
}

///获取全局的产品的价格列表
func GetGlobalPrice() ([]string, error) {
	Info("GetGlobalPrice...\n")

	para := &ST_CallPara{}
	ret := &ST_CallRet{}
	server := getCacheIO()
	if server == nil {
		return ret.Ids, ErrorLog("getCacheIO error \n")
	}
	if err := server.client.Call("CacheHandle.CallGetGlobalPrice", para, ret); err != nil {
		ErrorLog("CacheHandle.CallGetGlobalPrice error:%s", err.Error())
		return ret.Ids, err
	}
	return ret.Ids, nil
}

////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////产品修改接口////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////

//添加一个医院产品到新增加的缓存里面
func GlobalAddNewProduct(ProID, HosID, DocID, FirstItem, City, Name string, SecondItems []string, ProPrice int) error {
	Info("GlobalAddNewProduct...\n")
	if ProID == "" || HosID == "" || DocID == "" ||
		FirstItem == "" || City == "" ||
		Name == "" || len(SecondItems) == 0 ||
		ProPrice <= 0 {
		return ErrorLog("GlobalAddNewProduct failed,param is empty\n")
	}

	para := &ST_CallPara{
		HosID:       HosID,
		DocID:       DocID,
		ProID:       ProID,
		City:        City,
		FirstItem:   FirstItem,
		SecondItems: SecondItems,
		ProPrice:    ProPrice,
		ObjectName:  Name,
	}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallAddNewProduct", para, ret); err != nil {
		ErrorLog("CacheHandle.CallAddNewProduct error:%s", err.Error())
		return err
	}
	return nil
}

func GlobalUpdateProStatus(HosID string) error {
	Info("GlobalUpdateProStatus...\n")
	if HosID == "" {
		return ErrorLog("GlobalUpdateProStatus failed,HosID is empty\n")
	}
	para := &ST_CallPara{HosID: HosID}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallUpdateProductStatus", para, ret); err != nil {
		ErrorLog("CacheHandle.CallUpdateProductStatus error:%s", err.Error())
		return err
	}
	return nil
}

func GlobalUpdateProPrice(City, ProID string, newPrice, oldPrice int) error {
	Info("GlobalUpdateProPrice...\n")
	if ProID == "" {
		return ErrorLog("GlobalUpdateProPrice failed,ProID is empty\n")
	}
	para := &ST_ProPricePara{HosCity: City, ProID: ProID, OldPrice: oldPrice, NewPrice: newPrice}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallUpdateProPrice", para, ret); err != nil {
		ErrorLog("CacheHandle.CallUpdateProPrice error:%s", err.Error())
		return err
	}
	return nil
}

func GlobalUpdateProSale(ProID string, saleNum int) error {
	Info("GlobalUpdateProSale...\n")
	if ProID == "" || saleNum < 1 {
		return ErrorLog("GlobalUpdateProSale failed,ProID is empty,or SaleNum <1\n")
	}
	para := &ST_SortPara{ID: ProID, Num: saleNum}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallUpdateProSale", para, ret); err != nil {
		return ErrorLog("CacheHandle.CallUpdateProSale error:%s", err.Error())
	}
	return nil
}

func GlobalUpdateProComment(ProID string, Num int) error {
	Info("GlobalUpdateProComment...\n")
	if ProID == "" {
		return ErrorLog("GlobalUpdateProComment failed,ProID is empty,or Num <1\n")
	}
	para := &ST_SortPara{ID: ProID, Num: Num}
	ret := &ST_CallRet{}

	server := getCacheIO()
	if server == nil {
		return ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.CallUpdateProComment", para, ret); err != nil {
		ErrorLog("CacheHandle.CallUpdateProComment error:%s", err.Error())
		return err
	}
	return nil
}

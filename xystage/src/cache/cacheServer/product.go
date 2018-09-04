package main

import (
	. "JsLib/JsLogger"
	. "cache/cacheLib"
	"common"
	"constant"
	"strconv"
	"sync"
)

var gl_HosProONUM map[string]*ST_ONUMCache = make(map[string]*ST_ONUMCache) //所有的产品状态列表
var gl_HosProONUM_mutex sync.Mutex

var gl_ProONUM *ST_ONUMCache = &ST_ONUMCache{} //全局的产品状态列表
var gl_ProONUM_mutex sync.Mutex

var gl_FirstItemProduct *TypListCache = &TypListCache{} //全局分一级项目的产品列表
var gl_FirstItemProduct_mutex sync.Mutex

var gl_FirstToSecond *TypListCache = &TypListCache{} //全局一级项目和二级项目的映射
var gl_FirstToSecond_mutex sync.Mutex

var gl_CityProduct *TypListCache = &TypListCache{} //全局分城市的产品列表
var gl_CityProduct_mutex sync.Mutex

var gl_PricePro *TypListCache = &TypListCache{} //分价格的产品列表
var gl_PricePro_mutex sync.Mutex

var gl_AllPrice *ST_FullCache = &ST_FullCache{} //存放所有产品的价格
var gl_AllPrice_mutex sync.Mutex

var gl_ItemProduct map[string]*TypListCache = make(map[string]*TypListCache) //所有分一级二级项目的产品列表
var gl_ItemProduct_mutex sync.Mutex

var gl_CityPricePro map[string]*TypListCache = make(map[string]*TypListCache) //分城市分价格的产品列表
var gl_CityPricePro_mutex sync.Mutex

///////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////全局产品的获取接口///////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////

////全局产品的句柄
func (cache *CacheHandle) CallGetGlobalPro(para *ST_CallPara, ret *ST_ONUMCache) error {
	Info("CallGetGlobalPro\n")
	if gl_ProONUM == nil {
		return ErrorLog("CallGetGlobalDoc failed,gl_ProONUM is nil\n")
	}
	data := &ST_ONUMCache{}
	if para.SortType == "" {
		data = gl_ProONUM
	} else {
		data = SortONUM(para.SortType, 2)
	}
	ret.New = data.New
	ret.Modify = data.Modify
	ret.UnPass = data.UnPass
	ret.OffLine = data.OffLine
	ret.OnLine = data.OnLine
	ret.Full = data.Full
	return nil
}

//获取单个医院的全部的产品列表
func (cache *CacheHandle) CallGetHosProPtr(para *ST_CallPara, ret *ST_ONUMCache) error {
	Info("CallGetHosProPtr...\n")

	if para == nil || para.HosID == "" {
		return ErrorLog("CallGetHosProPtr failed,HosID=%s\n", para.HosID)
	}

	data, err := getHosProONUM(para.HosID)
	if err != nil {
		return err
	}
	ret.New = data.New
	ret.OnLine = data.OnLine
	ret.OffLine = data.OffLine
	ret.Modify = data.Modify
	ret.UnPass = data.UnPass
	ret.Full = data.Full

	return nil
}

////全局一级菜单映射的产品
func (cache *CacheHandle) CallGetGlobalFirstItemPro(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallGetGlobalFirstItemPro ... \n")
	if para == nil || para.FirstItem == "" {
		return ErrorLog("CallGetGlobalFirstItemPro failed,FirstItem =%s\n", para.FirstItem)
	}
	if v, ok := gl_FirstItemProduct.Data[para.FirstItem]; ok {
		ret.Ids = v
		return nil
	}
	ret.Ids = []string{}
	return nil
}

////全局分城市的产品
func (cache *CacheHandle) CallGetGlobalCityPro(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallGetGlobalCityPro...\n")
	if para == nil || para.City == "" {
		return ErrorLog("CallGetGlobalCityPro failed,City =%s\n", para.City)
	}
	if v, ok := gl_CityProduct.Data[para.City]; ok {
		ret.Ids = v
		return nil
	}
	ret.Ids = []string{}
	return nil
}

////获取一级菜单和二级菜单的映射
func (cache *CacheHandle) CallGetFirstMapSecondItem(para *ST_CallPara, ret *TypListCache) error {
	Info("CallGetFirstMapSecondItem...\n")
	ret.Data = gl_FirstToSecond.Data
	return nil
}

//所有一级菜单和二级菜单映射的产品列表
func (cache *CacheHandle) CallGetItemProduct(para *ST_CallPara, ret *ST_CallRet) error {

	Info("CallGetItemProduct...\n")
	if para == nil || para.FirstItem == "" || para.SecondItem == "" {
		return ErrorLog("CallGetItemProduct failed,FirstItem=%s,SecondItem=%s\n", para.FirstItem, para.SecondItem)
	}

	for k, v := range gl_ItemProduct {
		ErrorLog("gl_ItemProduct:k=%v, v=%v\n", k, v)
	}

	if data, ok := gl_ItemProduct[para.FirstItem]; ok {
		o := false
		ret.Ids, o = data.Data[para.SecondItem]
		if !o {
			Warn("CallGetItemProduct not exist FirstItem=%s,SecondItem=%s\n", para.FirstItem, para.SecondItem)
		}
		return nil
	} else {
		Warn("CallGetItemProduct not exist FirstItem=%s\n", para.FirstItem)
	}
	return nil
}

//所有分城市分价格的产品列表
func (cache *CacheHandle) CallGetCityPriceProduct(para *ST_CallPara, ret *TypListCache) error {
	Info("CallGetCityPriceProduct...\n")
	if para == nil || para.City == "" {
		return ErrorLog("CallGetCityPriceProduct failed,BodyPart=%s\n", para.City)
	}
	////获取分城市的价格列表
	if data, ok := gl_CityPricePro[para.City]; ok {
		ret.Data = data.Data
	} else {
		Warn("CallGetCityPriceProduct not exist City=%s\n", para.City)
	}
	return nil
}

//所有分价格的产品列表
func (cache *CacheHandle) CallGetPriceProduct(para *ST_CallPara, ret *TypListCache) error {
	Info("CallGetPriceProduct...\n")
	if gl_PricePro == nil {
		return ErrorLog("CallGetiPriceProduct is nil \n")
	}
	ret.Data = gl_PricePro.Data

	return nil
}

//所有产品价格的列表
func (cache *CacheHandle) CallGetGlobalPrice(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallGetGlobalPrice...\n")
	if gl_AllPrice == nil {
		return ErrorLog("CallGetGlobalPrice is nil \n")
	}
	ret.Ids = gl_AllPrice.Ids
	return nil
}

func (cache *CacheHandle) CallUpdateProductStatus(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallUpdateProductStatus...\n")
	if para.HosID == "" {
		return ErrorLog("CallUpdateProductStatus param HosID is empty\n")
	}
	update_product_status(para.HosID)
	return nil
}

///更新某一产品的价格
func (cache *CacheHandle) CallUpdateProPrice(para *ST_ProPricePara, ret *ST_CallRet) error {
	Info("CallUpdateHospitalCity ...\n")
	if para.HosCity == "" || para.ProID == "" || para.OldPrice <= 0 || para.NewPrice <= 0 {
		return ErrorLog("CallUpdateProPrice param empty,HosCity=%s,ProID=%s,OldPrice=%d,NewPrice=%d\n",
			para.HosCity, para.ProID, para.OldPrice, para.NewPrice)
	}
	update_product_price(para.HosCity, para.ProID, strconv.Itoa(para.NewPrice), strconv.Itoa(para.OldPrice))
	addToPricePro(para.ProID, para.NewPrice, true)
	return nil
}

func (cache *CacheHandle) CallUpdateProSale(para *ST_SortPara, ret *ST_CallRet) error {
	Info("CallUpdateProSale ...\n")
	if para.ID == "" || para.Num < 1 {
		return ErrorLog("CallUpdateProSale param empty,ID=%s,Num=%d\n", para.ID, para.Num)
	}
	addToSalePro(para.ID, para.Num, true)
	return nil
}

func (cache *CacheHandle) CallUpdateProComment(para *ST_SortPara, ret *ST_CallRet) error {
	Info("CallUpdateProComment ...\n")
	if para.ID == "" || para.Num < 1 {
		return ErrorLog("CallUpdateProComment param empty,ID=%s,Num=%d\n", para.ID, para.Num)
	}
	addToCommentPro(para.ID, para.Num, true)
	return nil
}

///添加一个新的产品
func (cache *CacheHandle) CallAddNewProduct(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallAddNewProduct ...\n")
	if para.ProID == "" || para.HosID == "" || para.ObjectName == "" ||
		para.City == "" || para.FirstItem == "" ||
		len(para.SecondItems) == 0 || para.ProPrice <= 0 {
		return ErrorLog("CallAddNewProduct param ProID empty \n")
	}

	loadProStatus(constant.OperatingStatus_new, para.ProID, para.HosID)
	loadProItem(para.ProID, para.FirstItem, para.SecondItems)
	loadProPrice(strconv.Itoa(para.ProPrice), para.ProID, para.City)
	loadProItemHosDoc(para.DocID, para.HosID, para.FirstItem, para.City)
	loadProCity(para.ProID, para.City)
	loadGobalProStatus()
	update_Search(para.ProID, para.ObjectName)
	return nil
}

//加载有个医院的所有产品
func load_product(HosID, HosCity string) {
	if HosID == "" || HosCity == "" {
		Error("load_product failed,HosID=%s,HosCity=%s\n", HosID, HosCity)
		return
	}
	gl_HosProONUM_mutex.Lock()
	gl_HosProONUM[HosID] = &ST_ONUMCache{}
	gl_HosProONUM_mutex.Unlock()

	rangeProduct(HosID, func(pro *common.ST_Product) {
		if pro.Current.OpreatStatus != constant.OperatingStatus_Del {
			loadProStatus(pro.Current.OpreatStatus, pro.ProID, HosID)
			loadProItem(pro.ProID, pro.FirstItem, pro.SecondItems)
			loadProPrice(strconv.Itoa(pro.XingYaoPrice), pro.ProID, HosCity)
			loadProItemHosDoc(pro.Doctors, pro.HosID, pro.FirstItem, HosCity)
			loadProCity(pro.ProID, HosCity)
			update_Search(pro.ProID, pro.ProName)
			addToSalePro(pro.ProID, pro.AppointmentNums, false)
			addToCommentPro(pro.ProID, pro.Peoples, false)
			addToPricePro(pro.ProID, pro.XingYaoPrice, false)
		}
	})
	loadGobalProStatus()
	return
}

//循环遍历所有的的产品
func rangeProduct(HosID string, F func(pro *common.ST_Product)) {
	proids := common.GetHosProductList(HosID)
	if len(proids) == 0 {
		Warn("HosID:%s,proids len is 0\n", HosID)
		return
	}

	prolist := common.QueryMoreProducts(proids)
	if len(prolist) == 0 {
		Error("HosID:%s,prolist len is 0\n", HosID)
		return
	}
	for _, pro := range prolist {
		F(pro)
	}
}

///加载产品状态缓存
func loadProStatus(Status, ProID, HosID string) {

	gl_HosProONUM_mutex.Lock()
	gl_HosProONUM[HosID].Full = append(gl_HosProONUM[HosID].Full, ProID)
	if Status == constant.OperatingStatus_new {
		gl_HosProONUM[HosID].New = append(gl_HosProONUM[HosID].New, ProID)
	}
	if Status == constant.OperatingStatus_online {
		gl_HosProONUM[HosID].OnLine = append(gl_HosProONUM[HosID].OnLine, ProID)
	}
	if Status == constant.OperatingStatus_Offline_self ||
		Status == constant.OperatingStatus_Offline_onforce {
		gl_HosProONUM[HosID].OffLine = append(gl_HosProONUM[HosID].OffLine, ProID)
	}
	if Status == constant.OperatingStatus_modify {
		gl_HosProONUM[HosID].Modify = append(gl_HosProONUM[HosID].Modify, ProID)
	}
	if Status == constant.OperatingStatus_Reviewer_NotPass {
		gl_HosProONUM[HosID].UnPass = append(gl_HosProONUM[HosID].UnPass, ProID)
	}

	gl_HosProONUM_mutex.Unlock()
}

///加载全局所有产品的状态缓存
func loadGobalProStatus() {
	gl_ProONUM_mutex.Lock()
	gl_ProONUM.Clear()
	for _, v := range gl_HosProONUM {
		gl_ProONUM.AppendOther(v)
	}
	gl_ProONUM_mutex.Unlock()
}

///加载所有的跟价格的相关的产品缓存
func loadProPrice(Price, ProID, HosCity string) {
	loadPricePro(Price, ProID)
	loadAllProPrice(Price)
	loadCityPricePro(Price, ProID, HosCity)
}

////加载分价格的产品缓存
func loadPricePro(Price, ProID string) {
	gl_PricePro_mutex.Lock()
	gl_PricePro.NewTypeList(Price, ProID)
	gl_PricePro_mutex.Unlock()
}

///加载全局的产品价格列表缓存
func loadAllProPrice(Price string) {
	gl_AllPrice_mutex.Lock()
	gl_AllPrice.AddToFull(Price)
	gl_AllPrice_mutex.Unlock()
}

///加载分城市价格缓存
func loadCityPricePro(Price, ProID, HosCity string) {
	gl_CityPricePro_mutex.Lock()
	if _, ok := gl_CityPricePro[HosCity]; !ok {
		gl_CityPricePro[HosCity] = &TypListCache{}
	}

	if _, ok := gl_CityPricePro[constant.All]; !ok {
		gl_CityPricePro[constant.All] = &TypListCache{}
	}

	gl_CityPricePro[HosCity].NewTypeList(Price, ProID)
	gl_CityPricePro[constant.All].NewTypeList(Price, ProID)
	gl_CityPricePro_mutex.Unlock()
}

//加载一二级菜单的产品缓存
func loadProItem(ProID, FirstItem string, SecondItems []string) {
	if FirstItem != "" && ProID != "" {

		gl_FirstItemProduct_mutex.Lock()
		gl_FirstItemProduct.NewTypeList(FirstItem, ProID)
		gl_FirstItemProduct_mutex.Unlock()

		gl_ItemProduct_mutex.Lock()
		if _, ok := gl_ItemProduct[FirstItem]; !ok {
			gl_ItemProduct[FirstItem] = &TypListCache{}
		}
		gl_ItemProduct_mutex.Unlock()
	}
	for _, item := range SecondItems {

		gl_FirstToSecond_mutex.Lock()
		gl_FirstToSecond.NewTypeList(FirstItem, item)
		gl_FirstToSecond_mutex.Unlock()

		gl_ItemProduct_mutex.Lock()
		gl_ItemProduct[FirstItem].NewTypeList(item, ProID)
		gl_ItemProduct_mutex.Unlock()
	}
}

///加载所有跟项目相关的医院和医生的缓存
func loadProItemHosDoc(DocID, HosID, FirstItem, HosCity string) {
	newCityItemHosCache(HosID, HosCity, FirstItem)
	newCityItemHosCache(HosID, constant.All, FirstItem)
	newCityItemDocCache(DocID, HosCity, FirstItem)
	newCityItemDocCache(DocID, constant.All, FirstItem)
}

///加载分城市的产品缓存
func loadProCity(ProID, HosCity string) {
	gl_CityProduct_mutex.Lock()
	gl_CityProduct.NewTypeList(HosCity, ProID)
	gl_CityProduct.NewTypeList(constant.All, ProID)
	gl_CityProduct_mutex.Unlock()
}

///清楚所有的产品缓存
func clearProcache() {
	gl_ProONUM_mutex.Lock()
	gl_ProONUM.Clear()
	gl_ProONUM_mutex.Unlock()

	gl_FirstItemProduct_mutex.Lock()
	gl_FirstItemProduct.Clear()
	gl_FirstItemProduct_mutex.Unlock()

	gl_FirstToSecond_mutex.Lock()
	gl_FirstToSecond.Clear()
	gl_FirstToSecond_mutex.Unlock()

	gl_CityProduct_mutex.Lock()
	gl_CityProduct.Clear()
	gl_CityProduct_mutex.Unlock()

	gl_PricePro_mutex.Lock()
	gl_PricePro.Clear()
	gl_PricePro_mutex.Unlock()

	gl_AllPrice_mutex.Lock()
	gl_AllPrice.Clear()
	gl_AllPrice_mutex.Unlock()

	gl_HosProONUM_mutex.Lock()
	gl_HosProONUM = make(map[string]*ST_ONUMCache)
	gl_HosProONUM_mutex.Unlock()

	gl_ItemProduct_mutex.Lock()
	gl_ItemProduct = make(map[string]*TypListCache)
	gl_ItemProduct_mutex.Unlock()

	gl_CityPricePro_mutex.Lock()
	gl_CityPricePro = make(map[string]*TypListCache)
	gl_CityPricePro_mutex.Unlock()

	gl_CityItemHos_mutex.Lock()
	gl_CityItemHos = make(map[string]*TypListCache)
	gl_CityItemHos_mutex.Unlock()

	gl_CityItemDoc_mutex.Lock()
	gl_CityItemDoc = make(map[string]*TypListCache)
	gl_CityItemDoc_mutex.Unlock()
}

///更新某一个产品的价格
func update_product_price(HosCity, ProID, NewPrice, OldPrice string) {

	gl_PricePro_mutex.Lock()
	gl_PricePro.RemoveTypeList(OldPrice, ProID)
	gl_PricePro.NewTypeList(NewPrice, ProID)
	gl_PricePro_mutex.Unlock()

	gl_CityPricePro_mutex.Lock()
	if _, ok := gl_CityPricePro[HosCity]; ok {
		gl_CityPricePro[HosCity].RemoveTypeList(OldPrice, ProID)
	} else {
		gl_CityPricePro[HosCity] = &TypListCache{}
	}
	gl_CityPricePro[HosCity].NewTypeList(NewPrice, ProID)
	gl_CityPricePro_mutex.Unlock()

	loadAllProPrice(NewPrice)
}

////更新某一个医院的所有产品的状态
func update_product_status(HosID string) {

	gl_HosProONUM_mutex.Lock()
	gl_HosProONUM[HosID] = &ST_ONUMCache{}
	gl_HosProONUM_mutex.Unlock()

	rangeProduct(HosID, func(pro *common.ST_Product) {
		loadProStatus(pro.Current.OpreatStatus, pro.ProID, HosID)
	})
	loadGobalProStatus()
}

///从缓存中找出医院产品状态缓存
func getHosProONUM(HosID string) (*ST_ONUMCache, error) {
	if HosID == "" {
		return nil, ErrorLog("getHosProONUM failed,HosID=%s\n", HosID)
	}
	data, ok := gl_HosProONUM[HosID]
	if ok {
		return data, nil
	}
	return data, ErrorLog("getHosProONUM failed,gl_HosProONUM[%s]failed\n", HosID)
}

///下线所有对应医院的产品
func offlineAllProduct(HosID string) bool {
	change := false
	if HosID == "" {
		Info("offlineAllProduct failed,HosID=%s\n", HosID)
		return change
	}
	rangeProduct(HosID, func(pro *common.ST_Product) {
		if pro.Current.OpreatStatus != constant.OperatingStatus_Offline_self &&
			pro.Current.OpreatStatus != constant.OperatingStatus_Offline_onforce {
			common.OfflineProOnForce(HosID, pro.ProID)
			change = true
		}
	})
	return change
}

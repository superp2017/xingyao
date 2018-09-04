package main

import (
	. "JsLib/JsLogger"
	. "cache/cacheLib"
	"common"
	"constant"
	"sync"
)

//全局的医生状态列表
var gl_DocONUM *ST_ONUMCache = &ST_ONUMCache{}
var gl_DocONUM_mutex sync.Mutex

///所有的分城市的医生列表
var gl_CityDoctor *TypListCache = &TypListCache{}
var gl_CityDoctor_mutex sync.Mutex

var gl_HosDocONUM map[string]*ST_ONUMCache = make(map[string]*ST_ONUMCache) //所有的医生状态列表
var gl_HosDocONUM_mutex sync.Mutex

var gl_CityItemDoc map[string]*TypListCache = make(map[string]*TypListCache) //分城市分一级项目的医生列表
var gl_CityItemDoc_mutex sync.Mutex

//获取全局医生的句柄
func (cache *CacheHandle) CallGetGlobalDoc(para *ST_CallPara, ret *ST_ONUMCache) error {
	Info("CallGetGlobalDoc...\n")
	if gl_DocONUM == nil {
		return ErrorLog("CallGetGlobalDoc failed,gl_DocONUM is nil\n")
	}
	data := &ST_ONUMCache{}
	if para.SortType == "" {
		data = gl_DocONUM
	} else {
		data = SortONUM(para.SortType, 1)
	}
	ret.New = data.New
	ret.Modify = data.Modify
	ret.UnPass = data.UnPass
	ret.OffLine = data.OffLine
	ret.OnLine = data.OnLine
	ret.Full = data.Full
	return nil
}

//获取单个医院的全部的医生列表
func (cache *CacheHandle) CallGetHosDocPtr(para *ST_CallPara, ret *ST_ONUMCache) error {
	Info("CallGetHosDocPtr...\n")
	if para == nil || para.HosID == "" {
		return ErrorLog("CallGetHosDocPtr failed,HosID=%s\n", para.HosID)
	}
	data, err := getHosDocONUM(para.HosID)
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

//获取全局的分城市医生
func (cache *CacheHandle) CallGetGlobalCityDoc(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallGetGlobalCityDoc ...\n")
	if gl_CityDoctor == nil || para.City == "" {
		return ErrorLog("gl_CityDoctor is nil ...\n")
	}
	if data, ok := gl_CityDoctor.Data[para.City]; ok {
		ret.Ids = data
		return nil
	}
	return nil
}

//获取全局的分城市分项目医生
func (cache *CacheHandle) CallGetCityItemDoc(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallGetCityItemDoc ...\n")
	if para.City == "" || para.FirstItem == "" {
		return ErrorLog("CallGetCityItemDoc City is empty ...\n")
	}
	if data, ok := gl_CityItemDoc[para.City]; ok {
		if ls, ex := data.Data[para.FirstItem]; ex {
			ret.Ids = ls
			return nil
		}
		return nil
	}
	return nil
}

///更新某一个医院的所有的医生状态
func (cache *CacheHandle) CallUpdateDoctor(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallUpdateDoc ...\n")
	if para.HosID == "" {
		return ErrorLog("CallUpdateDoc failed,HosID=%s\n", para.HosID)
	}
	update_doctor_status(para.HosID)
	return nil
}

func (cache *CacheHandle) CallUpdateDocSale(para *ST_SortPara, ret *ST_CallRet) error {
	Info("CallUpdateDocSale ...\n")
	if para.ID == "" || para.Num < 1 {
		return ErrorLog("CallUpdateDocSale param empty,ID=%s,Num=%d\n", para.ID, para.Num)
	}
	addToSaleDoc(para.ID, para.Num, true)
	return nil
}

func (cache *CacheHandle) CallUpdateDocComment(para *ST_SortPara, ret *ST_CallRet) error {
	Info("CallUpdateDocComment ...\n")
	if para.ID == "" || para.Num < 1 {
		return ErrorLog("CallUpdateDocComment param empty,ID=%s,Num=%d\n", para.ID, para.Num)
	}
	addToCommentDoc(para.ID, para.Num, true)
	return nil
}

//新增加一个医生
func (cache *CacheHandle) CallAddNewDoctor(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallAddNewDoctor ....\n")
	if para.DocID == "" || para.HosID == "" || para.City == "" || para.ObjectName == "" {
		return ErrorLog("CallAddNewDoctor failed,DocID=%s,HosID=%s,HosCity=%s\n", para.DocID, para.HosID, para.City)
	}

	loadDocStatus(constant.OperatingStatus_new, para.DocID, para.HosID)
	loadGobalDocStatus()
	loadCityDoc(para.DocID, para.City)
	update_Search(para.DocID, para.ObjectName)
	return nil
}

///更新某一个医院的所有的医生状态
func update_doctor_status(HosID string) {

	gl_HosDocONUM_mutex.Lock()
	gl_HosDocONUM[HosID] = &ST_ONUMCache{}
	gl_HosDocONUM_mutex.Unlock()

	rangeDoctor(HosID, func(doc *common.ST_Doctor) {
		loadDocStatus(doc.Current.OpreatStatus, doc.DocID, HosID)
	})
	loadGobalDocStatus()
}

////加载医生
func load_doctor(HosID, HosCity string) {
	Info("load_doctor ....\n")
	if HosID == "" && HosCity == "" {
		Error("load_doctor failed,HosID=%s,HosCity=%s\n", HosID, HosCity)
		return
	}

	gl_HosDocONUM_mutex.Lock()
	gl_HosDocONUM[HosID] = &ST_ONUMCache{}
	gl_HosDocONUM_mutex.Unlock()

	rangeDoctor(HosID, func(doc *common.ST_Doctor) {
		if doc.Current.OpreatStatus != constant.OperatingStatus_Del {
			loadDocStatus(doc.Current.OpreatStatus, doc.DocID, HosID)
			loadCityDoc(doc.DocID, HosCity)
			update_Search(doc.DocID, doc.DocName)
			addToSaleDoc(doc.DocID, doc.Reservepeople, false)
			addToCommentDoc(doc.DocID, doc.Evaluatepeople, false)
		}
	})
	loadGobalDocStatus()
	return
}

//循环遍历一个医院下的所有医生
func rangeDoctor(HosID string, F func(doc *common.ST_Doctor)) {

	docids := common.GetHosDoctorList(HosID)
	if len(docids) == 0 {
		Warn("HosID:%s,docids len is 0\n", HosID)
		return
	}

	doclist := common.QueryMoreDocInfo(docids)
	if len(doclist) == 0 {
		Warn("HosID:%s,doclist len is 0\n", HosID)
		return
	}
	for _, doc := range doclist {
		F(doc)
	}
}

////加载某一医生的状态
func loadDocStatus(Status, DocID, HosID string) {
	Info("loadDocStatus ....\n")

	gl_HosDocONUM_mutex.Lock()
	gl_HosDocONUM[HosID].Full = append(gl_HosDocONUM[HosID].Full, DocID)
	if Status == constant.OperatingStatus_new {
		gl_HosDocONUM[HosID].New = append(gl_HosDocONUM[HosID].New, DocID)
	}
	if Status == constant.OperatingStatus_online {
		gl_HosDocONUM[HosID].OnLine = append(gl_HosDocONUM[HosID].OnLine, DocID)
	}
	if Status == constant.OperatingStatus_Offline_self ||
		Status == constant.OperatingStatus_Offline_onforce {
		gl_HosDocONUM[HosID].OffLine = append(gl_HosDocONUM[HosID].OffLine, DocID)
	}
	if Status == constant.OperatingStatus_modify {
		gl_HosDocONUM[HosID].Modify = append(gl_HosDocONUM[HosID].Modify, DocID)
	}
	if Status == constant.OperatingStatus_Reviewer_NotPass {
		gl_HosDocONUM[HosID].UnPass = append(gl_HosDocONUM[HosID].UnPass, DocID)
	}
	gl_HosDocONUM_mutex.Unlock()
}

//加载全局的所有医生的状态
func loadGobalDocStatus() {
	gl_DocONUM_mutex.Lock()
	gl_DocONUM.Clear()
	for _, v := range gl_HosDocONUM {
		gl_DocONUM.AppendOther(v)
	}
	gl_DocONUM_mutex.Unlock()
}

///加载分城市医生
func loadCityDoc(DocID, HosCity string) {
	Info("loadCityDoc ....\n")
	gl_CityDoctor_mutex.Lock()
	gl_CityDoctor.NewTypeList(HosCity, DocID)
	gl_CityDoctor.NewTypeList(constant.All, DocID)
	gl_CityDoctor_mutex.Unlock()
}

///清除所有医生相关的缓存
func clearDocCache() {

	gl_DocONUM_mutex.Lock()
	gl_DocONUM.Clear()
	gl_DocONUM_mutex.Unlock()

	gl_CityDoctor_mutex.Lock()
	gl_CityDoctor.Clear()
	gl_CityDoctor_mutex.Unlock()

	gl_HosDocONUM_mutex.Lock()
	gl_HosDocONUM = make(map[string]*ST_ONUMCache)
	gl_HosDocONUM_mutex.Unlock()
}

//创建一个分城市分项目的医生缓存
func newCityItemDocCache(DocID, City, FirstItem string) {
	gl_CityItemDoc_mutex.Lock()
	if _, ok := gl_CityItemDoc[City]; !ok {
		gl_CityItemDoc[City] = &TypListCache{}
	}
	gl_CityItemDoc[City].NewTypeList(FirstItem, DocID)
	gl_CityItemDoc_mutex.Unlock()
}

///从缓存中找出医生状态缓存
func getHosDocONUM(HosID string) (*ST_ONUMCache, error) {
	if HosID == "" {
		return nil, ErrorLog("getHosDocONUM failed,HosID=%s\n", HosID)
	}
	data := &ST_ONUMCache{}
	ok := false
	data, ok = gl_HosDocONUM[HosID]
	if ok {
		return data, nil
	}
	return nil, ErrorLog("getHosDocONUM gl_HosDocONUM[%s] failed\n", HosID)
}

///下线所有医生
func offlineAllDoctor(HosID string) bool {
	change := false
	if HosID == "" {
		Info("offlineAllDoctor failed,HosID=%s\n", HosID)
		return change
	}

	rangeDoctor(HosID, func(doc *common.ST_Doctor) {
		if doc.Current.OpreatStatus != constant.OperatingStatus_Offline_self &&
			doc.Current.OpreatStatus != constant.OperatingStatus_Offline_onforce {
			common.OfflineDocOnForce(HosID, doc.DocID)
			change = true
		}
	})
	return change
}

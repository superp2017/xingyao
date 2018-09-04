package main

import (
	"JsLib/JsConfig"
	. "JsLib/JsLogger"
	. "cache/cacheLib"
	"common"
	"constant"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type CacheHandle int

var handle net.Listener

///测试是否连接
func (cache *CacheHandle) CallIsConnect(para *ST_CallPara, ret *ST_CallRet) error {
	Info("CallIsConnect")
	return nil
}

///开启服务
func StartServer() {
	ch := new(CacheHandle)
	rpc.Register(ch)
	rpc.HandleHTTP()
	var e error
	handle, e = net.Listen("tcp", ":"+JsConfig.CFG.CacheServer.Listen)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(handle, nil)
}

///停止服务
func StopServer() {
	if nil != handle {
		handle.Close()
	}
}

////初始化加载所有的缓存
func load_cache() {
	rangeHospital(func(hos *common.ST_Hospital) {
		if hos.Current.OpreatStatus != constant.OperatingStatus_Del {
			load_hospital(hos)
			checkHosStatus(hos, false)
			load_doctor(hos.HosID, hos.HosCity)
			load_product(hos.HosID, hos.HosCity)
		}
	})

	sortlist()
	////默认排序
	cache_sort()
}

func sortlist() {
	SortSalePro()
	SortCommentPro()
	SortPricePro()
	SortSaleHos()
	SortCommentHos()
	SortSaleDoc()
	SortCommentDoc()
}

////////全部按销量排序
func cache_sort() {
	gl_HosONUM_mutex.Lock()
	gl_HosONUM = SortONUM(SORT_SALE, 0)
	Error("gl_HosONUM=%v\n", gl_HosONUM)
	gl_HosONUM_mutex.Unlock()
	gl_DocONUM_mutex.Lock()
	gl_DocONUM = SortONUM(SORT_SALE, 1)
	Error("gl_DocONUM=%v\n", gl_DocONUM)
	gl_DocONUM_mutex.Unlock()
	gl_ProONUM_mutex.Lock()
	gl_ProONUM = SortONUM(SORT_SALE, 2)
	Error("gl_ProONUM=%v\n", gl_ProONUM)
	gl_ProONUM_mutex.Unlock()
}

///清除所有的缓存
func clear_Cache() {
	clearHoscache()
	clearDocCache()
	clearProcache()
	clearSortcache()
}

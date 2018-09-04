package cacheIO

import (
	. "JsLib/JsLogger"
	. "cache/cacheLib"
)

func GetSearchHash() (*ST_SearchHash, error) {
	Info("getSearchHash...\n")
	para := &ST_CallPara{}
	ret := &ST_SearchHash{}

	server := getCacheIO()
	if server == nil {
		return ret, ErrorLog("getCacheIO error \n")
	}

	if err := server.client.Call("CacheHandle.GetGlobalSearch", para, ret); err != nil {
		return ret, ErrorLog("CacheHandle.GetGlobalSearch error:%s", err.Error())
	}
	return ret, nil
}

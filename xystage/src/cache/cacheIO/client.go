package cacheIO

import (
	"JsLib/JsConfig"
	. "JsLib/JsLogger"
	. "cache/cacheLib"
	"net/rpc"
)

type ST_CacheServer struct {
	source string
	client *rpc.Client
}

///全局的缓存句柄
var cacheServer *ST_CacheServer

func init() {
	cacheServer = &ST_CacheServer{}
	cacheServer.source = JsConfig.CFG.CacheClient.Ip + ":" + JsConfig.CFG.CacheClient.Port
	cacheServer.client = nil
}

///获取缓存IO
func getCacheIO() *ST_CacheServer {
	if cacheServer.client != nil {
		para := &ST_CallPara{}
		ret := &ST_CallRet{}
		if err := cacheServer.client.Call("CacheHandle.CallIsConnect", para, ret); err != nil {
			cacheServer.client.Close()
			ErrorLog("CacheHandle.CallIsConnect error:%s", err.Error())
		} else {
			return cacheServer
		}
	}

	var err error
	cacheServer.client, err = rpc.DialHTTP("tcp", cacheServer.source)
	if err != nil {
		ErrorLog("getCacheIO error:%s", err.Error())
		cacheServer.client = nil
		return nil
	}
	return cacheServer
}

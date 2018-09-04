package main

import (
	. "JsLib/JsLogger"
	. "cache/cacheLib"
	// "common"
	// "constant"
	// "log"
	// "strconv"
	"sync"
)

var gl_Search *ST_SearchHash = &ST_SearchHash{} ///全局医院医生产品id与名字的对应
var gl_Search_mutex sync.Mutex

func (cache *CacheHandle) GetGlobalSearch(para *ST_CallPara, ret *ST_SearchHash) error {
	Info("GetGlobalSearch ... \n")
	ret.Data = gl_Search.Data
	return nil
}

func init_global() {
	gl_Search_mutex.Lock()
	gl_Search.Data = make(map[string]string)
	gl_Search_mutex.Unlock()
}

////更新搜索
func update_Search(ID, Name string) {
	gl_Search_mutex.Lock()
	gl_Search.Data[Name] = ID
	gl_Search_mutex.Unlock()
}

package main

import (
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"cache/cacheIO"
	. "cache/cacheLib"
	"constant"
	. "util"
)

func init_cache() {

	JsDispatcher.Http("/UniqueGlobalPrice", UniqueGlobalPrice) //重新加载全局缓存

	JsDispatcher.Http("/CreatGloblHospital", CreatGloblHospital) //重新加载全局缓存
	JsDispatcher.Http("/CreatHosDoctor", CreatHosDoctor)         //重新加载全局缓存
	JsDispatcher.Http("/CreataHosProduct", CreataHosProduct)     //重新加载全局缓存

}

func ClearAllProOrder(session *JsNet.StSession) {
	ids := gl_hospital()

	for _, v := range ids {
		data, err := cacheIO.GetHosProPtr(v)
		if err != nil {
			Error(err.Error())
			continue
		}
		DirectWrite("Product_Cache", v, data.Full)
	}
	Forward(session, "0", nil)
}

//GlobalProPrice
func UniqueGlobalPrice(session *JsNet.StSession) {
	data := ST_FullCache{}
	if err := WriteLock(constant.Hash_Global, constant.KEY_GlobalProPrice, &data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	list := []string{}
	for _, v := range data.Ids {
		exist := false
		for _, v1 := range list {
			if v == v1 {
				exist = true
			}
		}
		if !exist {
			list = append(list, v)
		}
	}
	data.Ids = list

	WriteBack(constant.Hash_Global, constant.KEY_GlobalProPrice, &data)
	Forward(session, "0", data)
}

func CreatGloblHospital(session *JsNet.StSession) {
	ids := gl_hospital()

	DirectWrite("Hospital_Cache", constant.KEY_ALL_Hospital, &ids)
	Forward(session, "0", nil)
}

func CreatHosDoctor(session *JsNet.StSession) {
	ids := gl_hospital()

	for _, v := range ids {
		data, err := cacheIO.GetHosDocPtr(v)
		if err != nil {
			Error(err.Error())
			continue
		}
		DirectWrite("Doctor_Cache", v, data.Full)
	}
	Forward(session, "0", nil)
}

func CreataHosProduct(session *JsNet.StSession) {
	ids := gl_hospital()

	for _, v := range ids {
		data, err := cacheIO.GetHosProPtr(v)
		if err != nil {
			Error(err.Error())
			continue
		}
		DirectWrite("Product_Cache", v, data.Full)
	}
	Forward(session, "0", nil)
}

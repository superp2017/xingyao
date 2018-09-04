// Copyright 20 The Go Authors. All rights reserved.
// Author : tianfeng
// Create Date : 2016/11/29
//
package main

import (
	"encoding/json"

	"github.com/pingplusplus/pingpp-go/pingpp"
	"github.com/pingplusplus/pingpp-go/pingpp/charge"

	//	"math/rand"
	//	"strconv"

	"sync"
	//	"time"
	//
	"JsLib/JsConfig"
	"JsLib/JsLogger"

	"JsLib/JsOrder"
)

var g_once sync.Once
var g_log *JsLogger.ST_Logger
var g_success_url string

func init_pingpp() {
	g_success_url = JsConfig.CFG.Pay.PaySuccessUrl
	initOnce()
}

type StPayPara struct {
	OrderID string
	Channel string
}

type StPayRet struct {
	Ret    string
	Charge []byte
	Desc   string
}

func initOnce() {
	// Initial the edition
	// LogLevel 是 Go SDK 提供的 debug 开关
	pingpp.LogLevel = 2
	// 设置 API Key
	// pingpp.Key = "sk_live_D4ejnHPaLq58TqXHWLnLePy5"
	pingpp.Key = JsConfig.CFG.PingPP.Key
	//获取 SDK 版本
	//fmt.Println("Go SDK Version:", pingpp.Version())
	//设置错误信息语言，默认是中文
	pingpp.AcceptLanguage = "zh-CN"
	//Initial the logger
	var ok bool
	g_log, ok = JsLogger.GetComLogger()

	if !ok {
		g_log.Error("Log regisiter wrong")
	}
	g_log.Console(false)
}

func buildCharge(order *JsOrder.StOrder) ([]byte, error) {

	// var para StPayPara
	// session.GetPara(&para)

	// var ret StPayRet

	// order := JsOrder.GetOrder(para.OrderID)
	// if order == nil {
	// 	g_log.Error("order[%s] is nil", para.OrderID)
	// 	ret.Ret = "1"
	// 	ret.Desc = fmt.Sprintf("order[%s] is nil", para.OrderID)
	// 	session.Forward(ret)
	// 	return
	// }

	metadata := make(map[string]interface{})
	metadata["color"] = "red"

	extra := make(map[string]interface{})
	if order.Channel == "alipay_wap" {
		extra["success_url"] = g_success_url
	} else if order.Channel == "alipay" {

	} else if order.Channel == "wx" {

	} else if order.Channel == "wx_pub" {
		openid := order.OpenId
		extra["open_id"] = openid
	}

	params := &pingpp.ChargeParams{

		Order_no: order.OrderId,
		//App:       pingpp.App{Id: "app_HW1W9SHuvHuT0enn"},
		App:       pingpp.App{Id: JsConfig.CFG.PingPP.AppId},
		Amount:    (uint64)(order.Amount),
		Channel:   order.Channel,
		Currency:  "cny",
		Client_ip: order.TerminalIp,
		Subject:   order.Subject,
		Body:      order.Desc,
		Extra:     extra,
		Metadata:  metadata}

	//返回的第一个参数是 charge 对象，你需要将其转换成 json 给客户端，或者客户端接收后转换。
	ch, err := charge.New(params)

	if err != nil {
		return nil, err
	}
	ret, err := json.Marshal(&ch)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

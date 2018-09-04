package main

import (
	"JsGo/JsMobile"
	"JsLib/JsDispatcher"
	"JsLib/JsExit"
	//"JsLib/JsNet"
	// "log"
)

func exit() int {
	JsDispatcher.Close()
	return 0
}

func init() {

	init_pingpp()
	init_wx_pay()
	order_init()
	init_agent()
}

func main() {
	//JsNet.AppConf("./conf/app.conf")

	JsExit.RegisterExitCb(exit)
	JsMobile.AlidayuInit()
	JsDispatcher.Run()
}

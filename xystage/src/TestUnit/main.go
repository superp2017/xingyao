package main

import (
	"JsLib/JsDispatcher"
	"JsLib/JsExit"
	//"JsLib/JsNet"
)

func exit() int {
	JsDispatcher.Close()
	return 0
}

func main() {

	//JsNet.AppConf("conf/app.conf")

	JsExit.RegisterExitCb(exit)
	Init_hospital()
	init_onum()
	init_testOrder()
	init_bill()
	init_agent()
	init_user()
	init_cache()
	init_update()
	init_statistic()
	init_withdraw()
	// StartTimer()

	JsDispatcher.Run()
}

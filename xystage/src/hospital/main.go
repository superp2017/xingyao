package main

import (
	"JsGo/JsMobile"
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
	init_router()
	JsMobile.AlidayuInit()
	// JsDispatcher.StartSession()
	JsDispatcher.Run()
}

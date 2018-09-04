package main

import (
	"JsGo/JsMobile"
	"JsLib/JsDispatcher"
	"JsLib/JsExit"
	//"JsLib/JsNet"
	"article"
)

func exit() int {
	JsDispatcher.Close()
	return 0
}

func main() {
	//JsNet.AppConf("conf/app.conf")
	JsExit.RegisterExitCb(exit)
	home_init()
	hospital_init()
	product_init()
	Order_init()
	Init_Router()
	Business_init()
	article.InitalArticle()
	JsMobile.AlidayuInit()
	JsDispatcher.Run()
}

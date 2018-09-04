package main

import (
	"JsGo/JsMobile"
	"JsLib/JsDispatcher"
	"JsLib/JsExit"
	"article"
)

func exit() int {
	JsDispatcher.Close()
	return 0
}

func main() {

	JsExit.RegisterExitCb(exit)
	init_router()
	init_home()
	article.InitalArticle()
	JsMobile.AlidayuInit()
	hospital_init()
	doctor_init()
	product_init()
	Order_init()
	User_init()
	Agent_init()

	JsDispatcher.Run()

}

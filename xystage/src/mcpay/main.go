package main

import (
	"JsLib/JsDispatcher"
	"JsLib/JsExit"
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
	JsExit.RegisterExitCb(exit)
	JsDispatcher.Run()
}

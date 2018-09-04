package main

import (
	"JsLib/JsExit"
	"log"
)

func exit() int {
	log.Println("global cache server begin exit...")
	StopServer()
	log.Println("global cache server end exit...")
	return 0
}

func main() {
	JsExit.RegisterExitCb(exit)
	init_global()
	load_cache()  ////加载初始数据
	StartServer() ////启动RPC服务
}

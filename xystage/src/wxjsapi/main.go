package main

import (
	"log"

	. "JsLib/JsConfig"
	"JsLib/JsExit"

	"github.com/astaxie/beego"
)

// 程序退出函数
func clean() int {
	log.Printf("WxJsapi exit...")

	return 0
}

func init() {
	JsExit.RegisterExitCb(clean)
}

func main() {
	beego.Router("/wxjsapi", &ST_WeChatJsapiController{})

	beego.Run(":" + CFG.Http.Listen)
}

package main

import (
	"log"

	"JsLib/JsConfig"
	"JsLib/JsExit"
	. "JsLib/JsLogger"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	. "util"
)

type ResCenter int

var handle net.Listener

// id通道map
var g_IdChan map[string]chan string

// 停止id服务
var g_stop_coolie bool = false

// 初始化id服务
func initIdMap() {
	g_IdChan = make(map[string]chan string)
}

// id 服务生成线程
func idCoolie(key string, c chan string, v string) {
	L, _ := strconv.Atoi(v)
	for {
		c <- v
		L++
		v = strconv.Itoa(L)
		Set(key, v)
		if g_stop_coolie {
			break
		}
	}
}

func stopCoolie() {
	g_stop_coolie = true
	for _, c := range g_IdChan {
		<-c
	}
}

func checkPara(para string) string {
	id := ""
	if e := Get(para, &id); e != nil {
		Error("Get Id key[%s] error:%s\n", para, e.Error())
		return "0"
	}

	g_IdChan[para] = make(chan string)

	go idCoolie(para, g_IdChan[para], id)

	ret := <-g_IdChan[para]

	return ret
}

func GenId(para string) string {
	c, ok := g_IdChan[para]
	var ret string
	if ok {
		ret = <-c
	} else {
		ret = checkPara(para)
	}

	Info("ID = %s\n", ret)
	return ret
}

func initRpcServer() {
	initIdMap()
}

func (h *ResCenter) GetId(key *string, ret *string) error {
	*ret = GenId(*key)
	return nil
}

func startServer() {
	initRpcServer()

	hb := new(ResCenter)

	rpc.Register(hb)

	rpc.HandleHTTP()

	//handle, e := net.Listen("tcp", ":8521")
	handle, e := net.Listen("tcp", ":"+JsConfig.CFG.IDer.Port)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	http.Serve(handle, nil)
}

func stopServer() {
	if nil != handle {
		handle.Close()
	}
	stopCoolie()
}

func exit() int {

	stopServer()
	return 0
}

func main() {
	JsExit.RegisterExitCb(exit)

	startServer()

}

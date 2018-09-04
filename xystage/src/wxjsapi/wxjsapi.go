package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	_ "fmt"
	_ "log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"JsLib/JsConfig"

	. "JsLib/JsLogger"

	"github.com/astaxie/beego"
)

type ST_WeChat_AccessToken struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}

type ST_WeChat_Jsapi_Ticket struct {
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	Ticket     string `json:"ticket"`
	Expires_in int    `json:"expires_in"`
}

type ST_WeChatJsapiController struct {
	beego.Controller
}

type ST_Jsapi_Interface struct {
	AppId     string   `json:"appid"`
	Timestamp string   `json:"timestamp"`
	NonceStr  string   `json:"nonceStr"`
	Signature string   `json:"signature"`
	JsApiList []string `json:"jsApiList"`
}

type ST_JsApiRet struct {
	ST_Jsapi_Interface
	Token string
}

type ST_ParaUrl struct {
	Url string `json:"url"`
}

var g_wechat_token string = ""
var g_wechat_jsapi_ticket string = ""
var g_jsConfig ST_Jsapi_Interface
var g_lock sync.Mutex
var g_jString string = ""

func init() {

	accessPath := JsConfig.CFG.WxJsApi.WeChatAccessToken

	ticketPath := JsConfig.CFG.WxJsApi.WeChatJsapiTicket

	g_jsConfig.AppId = JsConfig.CFG.WxJsApi.WeChatAppId

	g_jsConfig.JsApiList = strings.Split(JsConfig.CFG.WxJsApi.WeChatJsapiList, ",")

	go wechat_token_coolie(accessPath, ticketPath)
}

func wechat_token_coolie(accessPath, ticketPath string) {

	for {

		response, e := http.Get(accessPath)
		defer response.Body.Close()

		if e != nil {
			b := make([]byte, 2048)
			response.Body.Read(b)

		} else {
			b := make([]byte, 2048)
			n, _ := response.Body.Read(b)

			var token ST_WeChat_AccessToken
			json.Unmarshal(b[:n], &token)
			g_wechat_token = token.Access_token

			Info("token=%s", g_wechat_token)

			ticket_path := (ticketPath + "?access_token=" + g_wechat_token + "&type=jsapi")

			update_jsapi_ticket(g_wechat_token, ticket_path)
		}

		//sleep one hour
		time.Sleep(time.Hour)
	}
}

func update_jsapi_ticket(token, ticketPath string) {
	response, e := http.Get(ticketPath)
	defer response.Body.Close()
	if e != nil {
		b := make([]byte, 2048)
		response.Body.Read(b)

	} else {
		b := make([]byte, 2048)
		n, _ := response.Body.Read(b)

		var ticket ST_WeChat_Jsapi_Ticket
		json.Unmarshal(b[:n], &ticket)
		g_wechat_jsapi_ticket = ticket.Ticket

	}
}

func buildSignature(url string) {

	g_jsConfig.NonceStr = JsConfig.CFG.WxJsApi.WeChatNoncestr

	now := time.Now().Nanosecond()
	timestamp := strconv.Itoa(now)
	g_jsConfig.Timestamp = timestamp

	JsConfig.CFG.WxJsApi.WeChatTimeStamp = timestamp

	str := "jsapi_ticket="
	str += g_wechat_jsapi_ticket
	str += "&noncestr="
	str += g_jsConfig.NonceStr
	str += "&timestamp="
	str += timestamp
	str += "&url="
	str += url

	//产生一个散列值得方式是 sha1.New()，sha1.Write(bytes)，然后 sha1.Sum([]byte{})。这里我们从一个新的散列开始。
	h := sha1.New()
	//写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
	h.Write([]byte(str))
	//这个用来得到最终的散列值的字符切片。Sum 的参数可以用来dui现有的字符切片追加额外的字节切片：一般不需要要。
	g_jsConfig.Signature = fmt.Sprintf("%x", string(h.Sum(nil)))

}

func (this *ST_WeChatJsapiController) doWxJsapi() {
	cb := this.GetString("callback")
	data := this.GetString("data")
	var req_url ST_ParaUrl
	json.Unmarshal([]byte(data), &req_url)

	buildSignature(req_url.Url)

	ret := &ST_JsApiRet{}
	ret.AppId = g_jsConfig.AppId
	ret.JsApiList = g_jsConfig.JsApiList
	ret.NonceStr = g_jsConfig.NonceStr
	ret.Signature = g_jsConfig.Signature
	ret.Timestamp = g_jsConfig.Timestamp
	ret.Token = g_wechat_token

	g_lock.Lock()
	b, _ := json.Marshal(&ret)
	g_jString = string(b)

	this.Ctx.WriteString(cb + "(" + g_jString + ")")
	g_lock.Unlock()
}

func (this *ST_WeChatJsapiController) Get() {
	this.doWxJsapi()
}

func (this *ST_WeChatJsapiController) Post() {
	this.doWxJsapi()
}

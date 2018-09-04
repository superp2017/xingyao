package main

import (
	"JsLib/JsConfig"
	"JsLib/JsDispatcher"
	_ "JsLib/JsLogger"
	"JsLib/JsNet"
	"common"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/chanxuehong/wechat/mp/user/oauth2"
)

func init_home() {

	JsDispatcher.Http("/home", home)
	JsDispatcher.Http("/agent", agent)
	JsDispatcher.Http("/employee", employee)
	JsDispatcher.Http("/mcauth", mc_auth)
	JsDispatcher.Http("/appauth", app_auth)

}

//开始地方
func home(session *JsNet.StSession) {

	openid := session.Get("openid")

	type home_ret struct {
		Ret string
		Msg string
	}

	ret := &home_ret{}

	if openid == "" {
		//openid = session.GetCookie("meicheng_openid")
		// log.Println("-- openid = " + openid)

		if openid == "" {
			JsNet.CheckWxAuth(session, authCb)
			return
		}
	}

	if openid != "" {
		//session.SetCookie("meicheng_openid", openid)
		//openid = "?openid=" + openid

		// rs := ""
		// if strings.Index(session.UrlPath, "?") == -1 {
		// 	rs = "?openid=" + openid
		// } else {
		// 	rs = "&openid=" + openid
		// }

		// log.Panicln(JsConfig.CFG.WxJsApi.WeChatRedirectHome + session.UrlPath + rs)
		//

		rs := strings.Split(session.UrlPath, "?")
		if len(rs) == 2 {
			log.Println(JsConfig.CFG.WxJsApi.WeChatRedirectHome + "?" + rs[1])

			JsNet.HttpRedict(session, JsConfig.CFG.WxJsApi.WeChatRedirectHome+"?"+rs[1])
		} else {
			ret.Ret = "2"
			ret.Msg = "openid == nil"
			session.Forward(ret)
		}

	} else {
		ret.Ret = "1"
		ret.Msg = "openid == nil"
		session.Forward(ret)
	}
}

func authCb(user *oauth2.UserInfo, session *JsNet.StSession) {
	common.NewXyUser(user, false)
}

//开始地方
func agent(session *JsNet.StSession) {

	openid := session.Get("openid")

	type home_ret struct {
		Ret string
		Msg string
	}

	ret := &home_ret{}

	if openid == "" {
		if openid == "" {
			JsNet.CheckWxAuth(session, authCb)
			return
		}
	}

	if openid != "" {
		UID, _ := common.GetUIDFromOpenID(openid)
		openid = "?openid=" + openid + "&UID=" + UID
		JsNet.HttpRedict(session, "http://yiqizhuan.heyluckystar.com"+openid)

	} else {
		ret.Ret = "1"
		ret.Msg = "openid == nil"
		session.Forward(ret)
	}
}

func employee(session *JsNet.StSession) {
	openid := session.Get("openid")

	type home_ret struct {
		Ret string
		Msg string
	}

	ret := &home_ret{}

	if openid == "" {
		openid = session.GetCookie("xingyao_openid")

		if openid == "" {
			JsNet.CheckWxAuth(session, authCb)
			return
		}
	}

	if openid != "" {

		session.SetCookie("xingyao_openid", openid)

		openid = "?openid=" + openid + "&employee=true"

		JsNet.HttpRedict(session, "http://agent.xingyaostar.com"+openid)

	} else {
		ret.Ret = "1"
		ret.Msg = "openid == nil"
		session.Forward(ret)
	}
}

func mc_auth(session *JsNet.StSession) {
	openid := session.Get("openid")

	if openid == "" {
		JsNet.CheckWxAuth(session, mc_authCb)
	} else {
		openid = "?openid=" + openid

		JsNet.HttpRedict(session, "http://vuemall.xingyaostar.com"+openid)
	}

}

func mc_authCb(user *oauth2.UserInfo, session *JsNet.StSession) {
	common.NewXyUser(user, false)
}

func app_auth(session *JsNet.StSession) {
	type Para struct {
		Code string
	}

	type Ret struct {
		Ret  string
		Msg  string
		User *common.ST_User
	}
	para := &Para{}
	ret := &Ret{}

	e := session.GetPara(para)
	if e != nil {
		ret.Ret = "1"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}

	weurl := `https://api.weixin.qq.com/sns/oauth2/access_token?appid=wx606c59dfc4792dfe&secret=6191c2d4f35c812b29a2b02729491da5&code=` + para.Code + `&grant_type=authorization_code`
	resp, e := http.Get(weurl)
	if e != nil {
		ret.Ret = "2"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}

	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		ret.Ret = "3"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}

	log.Printf("resp = %s\n", string(b))

	type AuthToken struct {
		Access_token  string `json:access_token`
		Expires_in    int    `json:expires_in`
		Refresh_token string `json:refresh_token`
		Openid        string `json:openid`
		Scope         string `json:scope`
		Unionid       string `json:unionid`
	}

	token := &AuthToken{}
	e = json.Unmarshal(b, token)
	if e != nil {
		ret.Ret = "4"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}

	wxurl := `https://api.weixin.qq.com/sns/userinfo?access_token=` + token.Access_token + `&openid=` + token.Openid
	resp, e = http.Get(wxurl)
	if e != nil {
		ret.Ret = "5"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}

	b, e = ioutil.ReadAll(resp.Body)
	if e != nil {
		ret.Ret = "6"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}

	log.Printf("response = %s\n", string(b))

	userinfo := &oauth2.UserInfo{}

	e = json.Unmarshal(b, userinfo)
	if e != nil {
		ret.Ret = "7"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}

	ret.User = common.NewXyUser(userinfo, true)
	ret.Ret = "0"
	ret.Msg = "success"
	session.Forward(ret)
}

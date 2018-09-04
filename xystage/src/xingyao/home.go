package main

import (
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"log"
	"util"
)

///具体账号信息
type XYAccountInfo struct {
	XYID         string //后台
	LoginAccount string //登录账号
	LoginCode    string //登录密码
	LoginName    string //登录者姓名
	LoginCell    string //登录者手机号码
	Author       string //权限或者角色
}

//所有与医院关联的账号信息
type ST_XYAccount struct {
	Admin    XYAccountInfo            //超级管理员
	Employee map[string]XYAccountInfo //所有员工账号
}

var g_admin ST_XYAccount

func home_init() {
	JsDispatcher.Http("/login", login)
	JsDispatcher.WhiteList("/login")

	err := util.Get(constant.KV_Admin, &g_admin)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func login(session *JsNet.StSession) {

	Info("*****************Enter Login\n")
	type Para struct {
		Account string
		Pwd     string
	}

	type Ret struct {
		Return
		Type string
	}

	ret := &Ret{}
	para := &Para{}

	err := session.GetPara(para)
	Info("---------------------login info=%v\n", para)
	if err != nil {
		Error(err.Error())
		ret.Ret = "1"
		ret.Msg = err.Error()
		session.Forward(ret)
		return
	}

	if g_admin.Admin.LoginAccount == para.Account {
		if g_admin.Admin.LoginCode == para.Pwd {
			ret.Ret = "0"
			ret.Msg = "success"
			ret.Type = "admin"
			session.Forward(ret)
			return
		}
	} else if g_admin.Employee != nil {
		employee, ok := g_admin.Employee[para.Account]
		if ok {
			if employee.LoginAccount == para.Account && employee.LoginCode == para.Pwd {
				ret.Ret = "0"
				ret.Msg = "success"
				ret.Type = employee.Author
				session.Forward(ret)
				return
			}
		}
	}

	ret.Ret = "2"
	ret.Msg = "wrong user or password"
	session.Forward(ret)
}

package ider

import (
	"JsLib/JsConfig"
	. "JsLib/JsLogger"
	"constant"
	"net/rpc"
)

//生成通用id
func GenID() (string, error) {
	ret := ""
	if err := getID(constant.Com_SEED, &ret); err != nil {
		return "", err
	}
	return ret, nil
}

///生成 医院id
func GenHosID() (string, error) {
	ret := ""
	if err := getID(constant.HosID_SEED, &ret); err != nil {
		return "", err
	}
	return "hos-" + ret, nil
}

///生成 医生id
func GenDocID() (string, error) {
	ret := ""
	if err := getID(constant.Doc_SEED, &ret); err != nil {
		return "", err
	}
	return "doc-" + ret, nil
}

///生成 产品id
func GenProID() (string, error) {
	ret := ""
	if err := getID(constant.Pro_SEED, &ret); err != nil {
		return "", err
	}
	return "cp" + ret, nil
}

///生成 医院账号
func GenHosAccountID(permission string) (string, error) {
	ret := ""
	if err := getID(constant.HosAccountID_SEED, &ret); err != nil {
		return "", err
	}
	return permission + ret, nil
}

///生成 article id
func GenArtID() (string, error) {
	ret := ""
	if err := getID(constant.Art_SEED, &ret); err != nil {
		return "", err
	}
	return "Art-" + ret, nil
}

//Generate SuperArticle ID

func GenSuperArtID() (string, error) {
	ret := ""
	if err := getID(constant.SuperArt_SEED, &ret); err != nil {
		return "", err
	}
	return "SuperArt-" + ret, nil
}

//获取用户id
func GenUserID() (string, error) {
	ret := ""
	if err := getID(constant.USER_SEED, &ret); err != nil {
		return "", err
	}
	return "user-" + ret, nil
}

//生成 账单id
func GenBillID() (string, error) {
	ret := ""
	if err := getID(constant.Bill_SEED, &ret); err != nil {
		return "", err
	}
	return "zd" + ret, nil
}

///生成 订单id
func GenOrderID() (string, error) {
	ret := ""
	if err := getID(constant.Order_SEED, &ret); err != nil {
		return "", err
	}
	return "xd" + ret, nil
}

//生产代理id
func GenAgentID() (string, error) {
	ret := ""
	if err := getID(constant.Agent_SEED, &ret); err != nil {
		return "", err
	}
	return "agent-" + ret, nil
}

func GenwWithDrwaID() (string, error) {
	ret := ""
	if err := getID("WithDraw_SEED", &ret); err != nil {
		return "", err
	}
	return "wd-" + ret, nil
}

func getHandle() (*rpc.Client, error) {

	client, err := rpc.DialHTTP("tcp", JsConfig.CFG.IDer.Ip+":"+JsConfig.CFG.IDer.Port)
	if err != nil {
		return nil, ErrorLog("Ider getHandle error:%s", err.Error())
	}
	return client, nil
}

func getID(para string, ret *string) error {

	client, err := getHandle()
	if err != nil {
		ErrorLog(err.Error())
		return err
	}
	defer client.Close()

	if err = client.Call("ResCenter.GetId", &para, ret); err != nil {
		return ErrorLog("ResCenter.GetId error:%s", err.Error())
	}
	return nil
}

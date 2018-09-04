package common

import (
	"JsGo/JsMobile"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"ider"
	. "util"
)

///具体账号信息
type AccountInfo struct {
	HosID        string //医院ID
	LoginAccount string //登录账号
	LoginCode    string //登录密码
	LoginName    string //登录者姓名
	LoginCell    string //登录者手机号码
	Author       string //权限或者角色
}

//所有与医院关联的账号信息
type ST_HosAccount struct {
	Admin    AccountInfo   //超级管理员
	Employee []AccountInfo //所有员工账号
}

//用于登录账号映射医院id
type ST_AccountMap struct {
	AccountInfo //具体账号信息
}

func HosLogin(session *JsNet.StSession) {
	type st_login struct {
		LoginAccount string //登录账号
		LoginCode    string //登录密码
	}
	st := &st_login{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.LoginAccount == "" || st.LoginCode == "" {
		ForwardEx(session, "1", nil, "医院端登录失败HosLogin(),LoginAccount=%s,LoginCode=%s\n",
			st.LoginAccount, st.LoginCode)
		return
	}
	data := &ST_AccountMap{}
	if err := ShareLock(constant.Hash_HosAccountMap, st.LoginAccount, data); err != nil {
		Error(err.Error())
		ForwardEx(session, "1", nil, "err:%s, 医院端登录失败ShareLock(),LoginAccount=%s,LoginCode=%s\n", err.Error(),
			st.LoginAccount, st.LoginCode)
		return
	}
	if data.LoginCode != st.LoginCode {
		ForwardEx(session, "1", nil, "医院端登录失败,密码错误,LoginAccount=%s,LoginCode=%s\n",
			st.LoginAccount, st.LoginCode)
		return
	}
	type st_ret struct {
		Account  *ST_AccountMap
		Hospital *ST_Hospital
	}
	ret := &st_ret{
		Account: data,
	}
	Info("登录成功,账号为=%s\n", st.LoginAccount)
	if data.HosID != "" {
		//	查询医院信息
		hos, e := GetHospitalInfo(data.HosID)
		if e != nil {
			ForwardEx(session, "1", ret, "获取医院信息失败,HosID=%s\n", data.HosID)
			return
		}
		ret.Hospital = hos
	}

	session.NewSession(st.LoginAccount)
	Forward(session, "0", ret)
}

func creatAdminAccount(HosID, LoginName, LoginCell, LoginCode string) (*AccountInfo, error) {

	LoginAccount, err := genAccount(constant.Author_admin)
	if err != nil {
		return nil, ErrorLog("账号注册失败genAccount\n")
	}
	info := &AccountInfo{
		HosID:        HosID,
		LoginAccount: LoginAccount,
		LoginName:    LoginName, LoginCell: LoginCell,
		LoginCode: LoginCode,
		Author:    constant.Author_admin,
	}
	if err := writeToAccountMap(info); err != nil {
		return nil, ErrorLog("注册账号失败writeToAccountMap\n")
	}
	if err := writeToCellMap(LoginCell, LoginAccount); err != nil {
		return nil, ErrorLog("注册账号失败writeToCellMap\n")
	}
	return info, nil
}

///医院密码重置
func HosResetLoginCode(session *JsNet.StSession) {
	type st_modify struct {
		LoginAccount string //登录账号
		LoginCode    string //登录密码
		CheckCode    string //短信验证码
		LoginCell    string //管理员手机号
	}
	st := &st_modify{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.LoginAccount == "" || st.LoginCode == "" || st.CheckCode == "" || st.LoginCell == "" {
		ForwardEx(session, "1", nil, "修改管理员账号失败,参数为空,st=%v\n", st)
		return
	}

	info := &AccountInfo{}
	if err := ShareLock(constant.Hash_HosAccountMap, st.LoginAccount, info); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if info.LoginCell != st.LoginCell {
		ForwardEx(session, "1", nil, "该手机号%s和登录的账号%s不匹配\n", st.LoginCell, st.LoginAccount)
		return
	}

	if ok := JsMobile.VerifySmsCode(st.LoginCell, st.CheckCode); !ok {
		ForwardEx(session, "2", nil, "修改管理员账号失败,验证码不正确\n")
		return
	}
	info.LoginCode = st.LoginCode

	if err := modifyAccountinfo(info, true); err != nil {
		ForwardEx(session, "3", nil, "修改管理员账号失败,modifyAccountinfo 失败,st=%v\n", st)
		return
	}
	Forward(session, "0", nil)
}

///修改管理员信息
func ModifyAdminInfoAccount(session *JsNet.StSession) {
	type st_modify struct {
		HosID        string //医院ID
		LoginAccount string //登录账号
		LoginCode    string //登录密码
		LoginName    string //登录者姓名
		LoginCell    string //登录者手机号码
		Author       string //用户权限
	}
	st := &st_modify{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" || st.LoginAccount == "" ||
		st.LoginName == "" || st.LoginCell == "" {
		ForwardEx(session, "1", nil, "修改管理员账号失败,参数为空,st=%v\n", st)
		return
	}

	data := &AccountInfo{
		HosID:        st.HosID,
		LoginAccount: st.LoginAccount,
		LoginCode:    st.LoginCode,
		LoginName:    st.LoginName,
		LoginCell:    st.LoginCell,
		Author:       st.Author,
	}
	if err := modifyAccountinfo(data, true); err != nil {
		ForwardEx(session, "3", nil, "修改管理员账号失败,modifyAccountinfo 失败,st=%v\n", st)
		return
	}
	Forward(session, "0", nil)
}

///创建医院员工账号
func AddEmployeeAccount(session *JsNet.StSession) {
	type st_employee struct {
		HosID     string
		City      string
		Author    string
		LoginCode string
	}
	st := &st_employee{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if st.HosID == "" || st.City == "" || st.Author == "" || st.LoginCode == "" {
		ForwardEx(session, "1", nil, "创建管理员的账号失败,HosID=%s,City=%s,Author=%s,LoginCode=%s\n",
			st.HosID, st.City, st.Author, st.LoginCode)
		return
	}
	number, err := genAccount(constant.Author_admin)
	if err != nil {
		ForwardEx(session, "1", nil, "创建账号失败,err:"+err.Error())
		return
	}
	info := &AccountInfo{
		HosID:        st.HosID,
		Author:       st.Author,
		LoginAccount: number,
		LoginCode:    st.LoginCode,
	}
	account := &ST_HosAccount{}
	if err := WriteLock(constant.Hash_HosAccount, st.HosID, account); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	account.Employee = append(account.Employee, *info)
	if err := WriteBack(constant.Hash_HosAccount, st.HosID, account); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", info)
}

///修改医院登录账号的密码
func ModifyEmployeeAccountInfo(session *JsNet.StSession) {
	st := &AccountInfo{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if err := modifyAccountinfo(st, false); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", nil)
}

///删除某一个账号
func DelHosAccount(session *JsNet.StSession) {
	type st_modify struct {
		HosID        string //医院id
		LoginAccount string //登录账号
	}
	st := &st_modify{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" || st.LoginAccount == "" {
		ForwardEx(session, "1", "删除账号失败,参数部分为空!,HosID=%s,LoginAccount=%s\n",
			st.HosID, st.LoginAccount)
		return
	}
	//删除医院账号
	delAccount(st.HosID, st.LoginAccount)
	Forward(session, "0", nil)
}

//写入到医院账户表
func writeToHosAccount(account *ST_HosAccount) error {
	if err := DirectWrite(constant.Hash_HosAccount, account.Admin.HosID, account); err != nil {
		Error(err.Error())
		return ErrorLog("创建管理员的账号写入失败DirectWrite(),LoginAccount=%s,HosID=%s\n",
			account.Admin.LoginAccount, account.Admin.HosID)
	}
	return nil
}

///获取医院的账户信息
func GetHosAccountInfo(HosID string) (*ST_HosAccount, error) {
	data := &ST_HosAccount{}
	if err := ShareLock(constant.Hash_HosAccount, HosID, data); err != nil {
		Error(err.Error())
		return data, err
	}
	return data, nil
}

//写入账号映射医院
func writeToAccountMap(info *AccountInfo) error {
	if err := DirectWrite(constant.Hash_HosAccountMap, info.LoginAccount, info); err != nil {
		Error(err.Error())
		return ErrorLog("账号映射医院写入失败DirectWrite(),LoginAccount=%s,HosID=%s\n",
			info.LoginAccount, info.HosID)
	}
	return nil
}

///手机映射账号信息
func writeToCellMap(LoginCell, LoginAccount string) error {
	if err := DirectWrite(constant.Hash_HosCellAccountMap, LoginCell, LoginAccount); err != nil {
		Error(err.Error())
		return ErrorLog("手机账号映射医院写入失败DirectWrite(),LoginCell=%s,LoginAccount=%s\n",
			LoginCell, LoginAccount)
	}
	return nil
}

///修改账号信息
func modifyAccountinfo(info *AccountInfo, isAdmin bool) error {
	////////////修改医院账户表///////////
	if err := modifyAccount(info, isAdmin); err != nil {
		Error(err.Error())
		return err
	}
	////////////修改医院账户映射表///////////
	if err := modifyAccountMap(info); err != nil {
		Error(err.Error())
		return err
	}
	return nil
}

//修改账号密码
func modifyAccount(info *AccountInfo, IsAdmin bool) error {
	data := &ST_HosAccount{}

	if err := WriteLock(constant.Hash_HosAccount, info.HosID, data); err != nil {
		return err
	}
	if IsAdmin {
		data.Admin = *info
		return WriteBack(constant.Hash_HosAccount, info.HosID, data)
	}
	for i, v := range data.Employee {
		if v.LoginAccount == info.LoginAccount {
			data.Employee[i] = *info
			break
		}
	}
	return WriteBack(constant.Hash_HosAccount, info.HosID, data)
}

//修改账户映射里面的账户密码
func modifyAccountMap(info *AccountInfo) error {
	data := ST_AccountMap{}
	if err := WriteLock(constant.Hash_HosAccountMap, info.LoginAccount, data); err != nil {
		return err
	}
	data.HosID = info.HosID
	data.LoginAccount = info.LoginAccount
	data.LoginCode = info.LoginCode
	data.LoginName = info.LoginName
	data.LoginCell = info.LoginCell
	data.Author = info.Author
	return WriteBack(constant.Hash_HosAccountMap, info.LoginAccount, data)
}

//删除医院账号
func delAccount(HosID, LoginAccount string) error {
	data := &ST_HosAccount{}
	if err := WriteLock(constant.Hash_HosAccount, HosID, data); err != nil {
		return err
	}
	index := 0
	for i, v := range data.Employee {
		if v.LoginAccount == LoginAccount {
			index = i
			break
		}
	}
	data.Employee = append(data.Employee[:index], data.Employee[index+1:]...)
	return WriteBack(constant.Hash_HosAccount, HosID, data)
}

//检查手机号是否已经注册为管理员
func checkcellMapExist(LoginCell string) (*AccountInfo, bool) {
	data := ""
	if err := ShareLock(constant.Hash_HosCellAccountMap, LoginCell, data); err != nil {
		return nil, false
	}
	info := &AccountInfo{}
	if err := ShareLock(constant.Hash_HosAccountMap, data, info); err != nil {
		return nil, true
	}
	return info, true
}

//生成账号规则
func genAccount(Author string) (string, error) {
	tag := "09"
	if Author == constant.Author_admin {
		tag = "00"
	} else if Author == constant.Author_doctor {
		tag = "01"
	} else if Author == constant.Author_finance {
		tag = "02"
	} else {
		tag = "09"
	}
	id, err := ider.GenHosAccountID(tag)
	if err != nil {
		return "", err
	}
	return id, nil
}

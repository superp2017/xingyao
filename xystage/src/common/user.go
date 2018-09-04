package common

import (
	"JsGo/JsMobile"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"ider"
	. "util"

	"github.com/chanxuehong/wechat/mp/user/oauth2"
)

type ST_UserBaseInfo struct {
	UID             string   //用户id
	Name            string   //用户姓名
	Cell            string   //电话
	City            string   //城市
	Province        string   //省份
	Country         string   //国家
	Address         string   //地址
	Pics            []string //形象照
	Sex             string   //性别
	Birthday        string   //生日
	Cosmetichistory []string //整容史
	Cosmeticing     []string //现在整容
	Weight          int      //体重
	Height          int      //身高
	Age             int      //年龄
	BWH             string   //三围
	HeadImageURL    string   //头像
	OpenId_web      string   //网页的openid
	OpenId_app      string   //App的openid
	UnionId         string   //UnionId 唯一标识
}

///用户的代理信息
type ST_UserAgentInfo struct {
	Agent     *ST_UserAgent       //用户的小B信息
	AgentInfo *ST_AgentSimpleInfo //上级代理信息
	BusSource string              //用户来源
}

type ST_UserRelation struct {
	HisPrograms []string //历史产品列表
	HisSearch   []string //搜索历史
	Orders      []string //用户订单列表
	HisOrders   []string //历史订单列表
}

type ST_User struct {
	ST_UserBaseInfo         //用户基本信息
	ST_UserAgentInfo        //用户的代理信息
	ST_UserRelation         //用户关联信息
	CreaDate         string //创建日期
}

func NewXyUser(user *oauth2.UserInfo, isApp bool) *ST_User {
	UID := ""
	var isE error = nil
	UID, isE = GetUIDFromUnionID(user.UnionId)
	if isE != nil || UID == "" {
		UID, isE = GetUIDFromOpenID(user.OpenId)
		if isE != nil || UID == "" {
			return NewUser(user, isApp)
		}
	}
	data := &ST_User{}
	if err := WriteLock(constant.Hash_User, data.UID, data); err != nil {
		Error(err.Error())
		return nil
	}
	if isApp {
		data.OpenId_app = user.OpenId
	} else {
		data.OpenId_web = user.OpenId
	}
	if err := WriteBack(constant.Hash_User, data.UID, data); err != nil {
		Error(isE.Error())
		return nil
	}
	if data != nil {
		if isApp {
			go openIDMapToUID(data.OpenId_app, UID)
		} else {
			go openIDMapToUID(data.OpenId_web, UID)
		}
		go UnionIDMapToUID(data.UnionId, UID)
	}
	return data
}

///新的用户
func NewUser(user *oauth2.UserInfo, isApp bool) *ST_User {
	st := &ST_User{CreaDate: CurTime()}
	if id, err := ider.GenUserID(); err == nil {
		st.UID = id
	} else {
		Error("生成用户id失败\n")
		return nil
	}
	st.Name = user.Nickname
	st.Province = user.Province
	st.Country = user.Country
	st.City = user.City
	if isApp {
		st.OpenId_app = user.OpenId
	} else {
		st.OpenId_web = user.OpenId
	}
	st.UnionId = user.UnionId
	st.HeadImageURL = user.HeadImageURL

	if user.Sex == 1 {
		st.Sex = "男"
	} else if user.Sex == 2 {
		st.Sex = "女"
	} else {
		st.Sex = "未知"
	}
	st.Agent = nil
	st.AgentInfo = nil
	if err := DirectWrite(constant.Hash_User, st.UID, st); err != nil {
		Error(err.Error())
		return nil
	}
	go newUserFav(st.UID)
	if isApp {
		go openIDMapToUID(st.OpenId_app, st.UID)
	} else {
		go openIDMapToUID(st.OpenId_web, st.UID)
	}
	go UnionIDMapToUID(st.UnionId, st.UID)
	go UserIncreaseOneDay()
	return st
}

///////将openID 映射 UID
func openIDMapToUID(openID, UID string) error {
	if err := DirectWrite(constant.Hash_OpenID_UID, openID, UID); err != nil {
		return err
	}
	return nil
}

///////unionID 映射 UID
func UnionIDMapToUID(unionID, UID string) error {
	if err := DirectWrite(constant.Hash_Union_UID, unionID, UID); err != nil {
		return err
	}
	return nil
}

///opendID 返回 UID
func GetUIDFromOpenID(openid string) (string, error) {
	UID := ""
	err := ShareLock(constant.Hash_OpenID_UID, openid, &UID)
	return UID, err
}

///opendID 返回 UID
func GetUIDFromUnionID(openid string) (string, error) {
	UID := ""
	err := ShareLock(constant.Hash_Union_UID, openid, &UID)
	return UID, err
}

///通过UnionId获取用户资料
func GetUserFormUnionID(session *JsNet.StSession) {
	type st_Get struct {
		UnionId string //openid
	}
	st := &st_Get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UnionId == "" {
		ForwardEx(session, "1", nil, "UnionId is empty\n")
		return
	}
	UID, err := GetUIDFromUnionID(st.UnionId)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	user, err := GetUserInfo(UID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if user != nil {
		go checkUserOrder(user)
	}
	Forward(session, "0", user)
}

///通过OpenId获取用户资料
func GetUserFormOpenID(session *JsNet.StSession) {
	type st_Get struct {
		OpenId string //openid
	}
	st := &st_Get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.OpenId == "" {
		ForwardEx(session, "1", nil, "OpenId is empty\n")
		return
	}
	UID, err := GetUIDFromOpenID(st.OpenId)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	user, err := GetUserInfo(UID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if user != nil {
		go checkUserOrder(user)
	}
	Forward(session, "0", user)
}

///修改用户
func ModifyUser(session *JsNet.StSession) {
	type st_modify struct {
		UID             string   //用户id
		Name            string   //用户姓名
		HeadImageURL    string   //头像
		Sex             string   //性别
		Birthday        string   //生日
		Age             int      //年龄
		City            string   //地址
		Cosmetichistory []string //整容史
		Cosmeticing     []string //现在整容
	}
	st := &st_modify{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" {
		ForwardEx(session, "1", nil, "ModifyUser failed,openid or uid is empty\n")
		return
	}
	data := &ST_User{}
	if err := WriteLock(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.Name = st.Name
	data.HeadImageURL = st.HeadImageURL
	data.Sex = st.Sex
	data.Birthday = st.Birthday
	data.Age = st.Age
	data.City = st.City
	data.Cosmetichistory = st.Cosmetichistory
	data.Cosmeticing = st.Cosmeticing
	if err := WriteBack(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

//绑定手机
func BindUserCell(session *JsNet.StSession) {
	change_usercell(session)
}

func change_usercell(session *JsNet.StSession) {
	type st_Cell struct {
		UID     string
		Cell    string
		MsgCode string
	}
	st := &st_Cell{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" || st.Cell == "" || st.MsgCode == "" {
		ForwardEx(session, "1", nil, "绑定手机号,参数不全,UID=%s,Cell=%s,MsgCode=%s\n", st.UID, st.Cell, st.MsgCode)
		return
	}
	///手机验证码
	if ok := JsMobile.VerifySmsCode(st.Cell, st.MsgCode); !ok {
		ForwardEx(session, "2", nil, "BindUserCell faild,验证码不正确\n")
		return
	}
	uid := ""
	ShareLock(constant.KEY_Cell_UID, st.Cell, &uid)
	if uid != "" {
		go removeUserCell(uid)
	}
	go CellMapToUID(st.UID, st.Cell)
	data := &ST_User{}
	if err := WriteLock(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	if data.Cell != st.Cell {
		if data.Cell != "" {
			go HDel(constant.KEY_Cell_UID, data.Cell)
		}
		data.Cell = st.Cell
	}
	if err := WriteBack(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	Forward(session, "0", nil)
}

//清空某一个用户的手机号
func removeUserCell(UID string) error {
	data := &ST_User{}
	err := WriteLock(constant.Hash_User, UID, data)
	if err == nil {
		data.Cell = ""
		return WriteBack(constant.Hash_User, UID, data)
	} else {
		return err
	}
	return nil
}

//换绑手机
func TochangeCell(session *JsNet.StSession) {
	change_usercell(session)
}

//订单无效的时候,将订单移除
func RemoveUserContinueOrder(UID, OrderID string) error {
	if UID == "" || OrderID == "" {
		return ErrorLog("RemoveUserContinueOrder failed,UID=%s,OrderID =%s\n", UID, OrderID)
	}
	data := &ST_User{}
	if err := WriteLock(constant.Hash_User, UID, data); err != nil {
		return err
	}
	for i, v := range data.Orders {
		if v == OrderID {
			data.Orders = append(data.Orders[:i], data.Orders[i+1:]...)
			break
		}
	}
	return WriteBack(constant.Hash_User, UID, data)
}

//将Cell 映射 UID
func CellMapToUID(UID, Cell string) error {
	if err := DirectWrite(constant.KEY_Cell_UID, Cell, UID); err != nil {
		return err
	}
	return nil
}

//查询用户信息
func QueryUserInfo(session *JsNet.StSession) {
	type st_get struct {
		UID string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data, err := GetUserInfo(st.UID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	go checkUserOrder(data)
	Forward(session, "0", data)

}

//后台查询用户信息
func GetUserInfo(UID string) (*ST_User, error) {
	if UID == "" {
		return nil, ErrorLog("UID为空,查询用户失败\n")
	}
	data := &ST_User{}
	if err := ShareLock(constant.Hash_User, UID, data); err != nil {
		return nil, ErrorLog("QueryUserInfo_Back,ShareLock failed" + err.Error())
	}
	return data, nil
}

///添加了一个订单ID
func NewUserOrder(UID, OrderID string) error {
	if UID == "" || OrderID == "" {
		return ErrorLog("NewUserOrder failed,UID=%s,OrderID=%s\n", UID, OrderID)
	}
	data := &ST_User{}
	if err := WriteLock(constant.Hash_User, UID, data); err != nil {
		return err
	}
	exist := false
	for _, v := range data.Orders {
		if v == OrderID {
			exist = true
			break
		}
	}
	if !exist {
		if data.Orders == nil || len(data.Orders) == 0 {
			data.Orders = []string{}
		}
		data.Orders = append(data.Orders, OrderID)
	}
	return WriteBack(constant.Hash_User, UID, data)

}

//获取用户订单列表
func GetUserAllOrderList(UID string) ([]string, error) {
	if UID == "" {
		return nil, ErrorLog("GetUserOrderList failed,UID=%s\n", UID)
	}
	data := &ST_User{}
	if err := ShareLock(constant.Hash_User, UID, data); err != nil {
		return []string{}, err
	}
	return data.Orders, nil
}

////检查无效的订单
func checkUserOrder(data *ST_User) {
	ids := []string{}
	for _, v := range data.Orders {
		or, err := QueryOrder(v)
		if err == nil {
			if or.Current.OpreatStatus == constant.Status_Order_PenddingPay &&
				CurStamp()-or.SubmitStamp >= 30*60 {
				ids = append(ids, v)
			}
		}
	}
	/////订单取消
	for _, v := range ids {
		OrderInvalid(v)
	}
}

///获取多个用户信息
func GetMoreUserInfo(uids []string) []*ST_User {
	data := []*ST_User{}
	for _, v := range uids {
		if user, err := GetUserInfo(v); err == nil {
			data = append(data, user)
		} else {
			ErrorLog(err.Error())
		}
	}
	return data
}

type ST_AgentInviteStatistic struct {
	UID       string
	UserName  string
	AgentID   string
	AgentCity string
	ApplyDate string
	CreatDate string
}

//代理邀请用户统计
func AgentInviteStatistics(user ST_AgentInviteStatistic) error {
	month := CurYearMonth()
	data := []ST_AgentInviteStatistic{}
	err := WriteLock(constant.Hash_UserFromAgent, month, &data)
	if err != nil {
		data = append(data, user)
		return DirectWrite(constant.Hash_UserFromAgent, month, &data)
	}
	exist := false
	for _, v := range data {
		if v.UID == user.UID {
			exist = true
			break
		}
	}
	if !exist {
		data = append(data, user)
	}
	return WriteBack(constant.Hash_UserFromAgent, month, &data)
}

///用户每天增加量
func UserIncreaseOneDay() error {
	data := make(map[string]int)
	err := WriteLock(constant.Hash_User, constant.KEY_UserIncrease, &data)
	if err != nil {
		data[CurDate()] = 1
		return DirectWrite(constant.Hash_User, constant.KEY_UserIncrease, &data)
	}
	if num, ok := data[CurDate()]; ok {
		data[CurDate()] = num + 1
	} else {
		data[CurDate()] = 1
	}
	return WriteBack(constant.Hash_User, constant.KEY_UserIncrease, &data)
}

//获取小B邀请统计
func GetAgentInviteStatistic(session *JsNet.StSession) {
	type st_Get struct {
		Date    string
		IsMonth bool
	}
	st := &st_Get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.Date == "" {
		ForwardEx(session, "1", nil, "GetUserStatistic faild,param is empty\n")
		return
	}
	yearMonth := st.Date
	if !st.IsMonth {
		time, err := GetTimeFormString(st.Date)
		if err != nil {
			ForwardEx(session, "1", nil, err.Error())
			return
		}
		yearMonth = GetYearMOnth(time)
	}

	data := []ST_AgentInviteStatistic{}
	if e := ShareLock(constant.Hash_UserFromAgent, yearMonth, &data); e != nil {
		ForwardEx(session, "1", nil, e.Error())
		return
	}
	if st.IsMonth {
		Forward(session, "0", data)
		return
	}
	retdata := []ST_AgentInviteStatistic{}
	for _, v := range data {
		if v.ApplyDate == st.Date {
			retdata = append(retdata, v)
		}
	}
	Forward(session, "0", retdata)
}

//获取用户每天的增加量
func GetUserIncrease(session *JsNet.StSession) {
	data := make(map[string]int)
	if err := ShareLock(constant.Hash_User, constant.KEY_UserIncrease, &data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)

}

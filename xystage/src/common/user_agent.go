package common

import (
	"JsGo/JsMobile"
	"JsLib/JsConfig"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	//"crypto"
	"strconv"
	"time"
	. "util"
)

//店员的信息
type ST_TeamInfo struct {
	UID              string               //uid
	Name             string               //姓名
	JoinDate         string               //加入时间
	Orders           []ST_CommissionOrder //带来的订单
	ToltleCommission int                  //带来的预计总提成
	TotalIncome      int                  //带来的实际收入
	HeadPic          string               //头像
}

type ST_CommissionOrder struct {
	OrderID       string //订单
	PreCommisson  int    //预计提成
	LastCommisson int    //最终提成
	PreDate       string //校验时间
	LastDate      string //评价时间
}

//代理的简短信息
type ST_AgentSimpleInfo struct {
	UID      string //代理ID
	Name     string //业务人员名称
	City     string //代理城市
	Cell     string //联系号码
	ProRatio int    //产品结算比例
}

///代理基本账户信息
type ST_AgentBaseInfo struct {
	IdentityNumber   string   // 身份证号
	IdentityPic      string   // 身份证照片
	BankNumber       string   // 银行卡号
	BankName         string   // 开户行
	WxNumber         string   // 微信号
	AgentQR          string   // 二维码
	ToltleCommission int      // 总提成
	TotalIncome      int      // 小B总收入
	Balance          int      // 小B钱包余额
	Withdraw         int      // 总提现金额
	WithDrawRecord   []string // 小B提现记录
}

//代理的详细信息
type ST_AgentDetail struct {
	AgentLevel string //代理级别
	AgentType  string //合伙人类别:实体店铺、微商
	////////////////////////////////////////////////////
	ShopName          string // （实体店铺）店铺名称
	ShopAddress       string // （实体店铺）店铺地址
	ShopBussiness     string // （实体店铺）主营业务
	ShopBL            string // （实体店铺）营业执照图片
	ShopMonthTurnover string // （实体店铺）月营业额
	////////////////////////////////////////////////////
	WxItem          string // （微商）主营项目
	WxMonthTurnover string // （微商）月营业额
	WxAddress       string // （微商）通讯地址
	////////////////////////////////////////////////////
	PlatformFee  int    // 平台使用费
	Agreement    string // 协议
	Bond         int    // 保证金
	FranchiseFee int    // 加盟费
}

//代理的状态
type ST_AgentStatus struct {
	OpreatStatus string //状态
	OpreatTime   string //操作时间
	OpreatReason string //操作原因
	OpreatJobNum string //操作人员工号
	OpreatName   string //操作人员姓名
}

type ST_AgentOpreat struct {
	ApplyDate    string           //申请时间
	ReApplyDate  string           //重新申请时间
	RevieweDate  string           //审核日期
	PayDate      string           //支付时间
	RePayDate    string           //重新支付时间
	OnlineDate   string           //上线时间
	OfflineDate  string           //下线时间
	OfflineStamp int64            //下线的时间戳
	WithdrawDate string           //撤回保证金时间
	OpreatRecord []ST_AgentStatus //操作记录
	Current      ST_AgentStatus   //当期操作记录
}

//代理的相关的列表
type ST_AgentRelation struct {
	Customer    []ST_TeamInfo        //邀请的客户
	Orders      []ST_CommissionOrder //相关的订单
	Article     []ST_AgentArticle    //关联的软文
	Bills       []string             //相关的账单
	BondOrderID string               //保证金支付订单id
}

type ST_UserAgent struct {
	ST_AgentSimpleInfo        //代理的简单信息
	ST_AgentBaseInfo          //代理的基本信息
	ST_AgentDetail            //代理的代理信息
	ST_AgentOpreat            //代理的操作信息
	ST_AgentRelation          //代理的相关列表
	CreatDate          string //创建日期
}

type ST_AgentArticle struct {
	ArticleID string
	ShareTime string
}

//申请成为代理
func ApplyAgent(session *JsNet.StSession) {
	st := &ST_UserAgent{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if err := checkAgent(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	st.ApplyDate = CurTime()
	st.CreatDate = CurTime()
	opreat := ST_AgentStatus{
		OpreatStatus: constant.Agent_Apply,
		OpreatTime:   CurTime(),
		OpreatReason: constant.Agent_Apply,
		OpreatJobNum: st.UID,
		OpreatName:   st.Name,
	}
	st.Current = opreat
	st.OpreatRecord = append(st.OpreatRecord, opreat)
	////////更新用户的代理信息/////////
	data := &ST_User{}
	if err := WriteLock(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_User, st.UID, data)
	if data.Agent != nil {
		if data.Agent.Current.OpreatStatus == constant.Agent_Online ||
			data.Agent.Current.OpreatStatus == constant.Agent_Apply ||
			data.Agent.Current.OpreatStatus == constant.Agent_ReApply {
			ForwardEx(session, "1", data, "该用户已经是小B，或者正在申请小B，UID=%s", data.UID)
			return
		} else {
			ErrorLog("不在线的小B，直接替换信息,UID=%s\n", data.UID)
		}
	}
	data.Agent = st
	go NewGlobalAgent(st.City, st.UID)
	Forward(session, "0", data)
}

//重新申请(需要审核)
func ReApplyAgent(session *JsNet.StSession) {
	st := &ST_UserAgent{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if err := checkAgent(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	/////////更新用户的代理信息////////////////
	data := &ST_User{}

	if err := WriteLock(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_User, st.UID, data)

	if data.Agent.Current.OpreatStatus != constant.Agent_Apply &&
		data.Agent.Current.OpreatStatus != constant.Agent_NoPassReviewe {
		ForwardEx(session, "1", nil, "ReApplyAgent, Update failed ,status =%s\n", data.Agent.Current.OpreatStatus)
		return
	}
	data.Agent = st
	data.Agent.ApplyDate = CurTime()
	data.Agent.ReApplyDate = CurTime()
	opreat := ST_AgentStatus{
		OpreatStatus: constant.Agent_Apply,
		OpreatTime:   CurTime(),
		OpreatReason: "修改后重新提交申请",
		OpreatJobNum: st.UID,
		OpreatName:   st.Name,
	}
	data.Agent.Current = opreat
	data.Agent.OpreatRecord = append(st.OpreatRecord, opreat)
	Forward(session, "0", data)
}

//修改小B产品结算比例
func ModifyAgentProRatio(session *JsNet.StSession) {
	type st_get struct {
		UID      string //UID
		ProRatio int    //产品结算比例
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" || st.ProRatio < 0 {
		ForwardEx(session, "1", nil, "ModifyAgentProRatio failed,UID=%s,ProRatio=%d\n", st.UID, st.ProRatio)
		return
	}
	data := &ST_User{}
	if err := WriteLock(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.Agent.ProRatio = st.ProRatio
	if err := WriteBack(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

////更新用户的代理的基本信息（不需要再次审核）
func ModifyAgentBaseInfo(session *JsNet.StSession) {
	type st_get struct {
		UID            string // uid
		Name           string // 名称
		Cell           string // 联系号码
		IdentityNumber string // 身份证号
		IdentityPic    string // 身份证照片
		BankNumber     string // 银行卡号
		BankName       string // 开户行
		WxNumber       string // 微信号
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" {
		ForwardEx(session, "1", nil, "ModifyAgentBaseInfo failed,UID is empty\n")
		return
	}
	/////////更新用户的代理信息////////////////
	data := &ST_User{}
	if err := WriteLock(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_User, st.UID, data)
	if data.Agent == nil {
		ForwardEx(session, "1", data, "当前用户还不是小B,UID=%s\n", st.UID)
		return
	}
	data.Agent.Name = st.Name
	data.Agent.Cell = st.Cell
	data.Agent.IdentityNumber = st.IdentityNumber
	data.Agent.IdentityPic = st.IdentityPic
	data.Agent.BankNumber = st.BankNumber
	data.Agent.BankName = st.BankName
	data.Agent.WxNumber = st.WxNumber
	Forward(session, "0", data)
}

//审核代理人
func RevieweAgent(session *JsNet.StSession) {
	type st_review struct {
		UID       string //用户id
		JobNumber string //审核人的工号
		JobName   string //审核人姓名
		Reason    string //审核备注
		IsPass    bool   //是否同意
		ProRatio  int    //产品结算比例
	}
	st := &st_review{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" || st.JobNumber == "" || st.ProRatio < 0 {
		ForwardEx(session, "1", nil, "RevieweAgent failed,param is empty,UID=%s,JobNumber=%s\n", st.UID, st.JobNumber)
		return
	}

	data := &ST_User{}
	if err := WriteLock(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_User, st.UID, data)
	if data.Agent == nil {
		ForwardEx(session, "1", nil, "RevieweAgent failed,User=%s is not apply Agent", st.UID)
		return
	}
	if data.Agent.Current.OpreatStatus != constant.Agent_ReApply && data.Agent.Current.OpreatStatus != constant.Agent_Apply {
		ForwardEx(session, "1", nil, "RevieweAgent  failed ,current status is   %s\n", data.Agent.Current.OpreatStatus)
		return
	}
	data.Agent.RevieweDate = CurTime()
	status := constant.Agent_NoPassReviewe
	if st.IsPass {
		data.AgentInfo = nil //上级代理
		status = constant.Agent_Online
		JsMobile.ComJsMobileVerify("喜妹儿", data.Agent.Cell, "SMS_126971460", "a", 3600, nil)
	}
	opreat := ST_AgentStatus{
		OpreatStatus: status,
		OpreatTime:   CurTime(),
		OpreatReason: st.Reason,
		OpreatJobNum: st.JobNumber,
		OpreatName:   st.JobName,
	}
	data.Agent.ProRatio = st.ProRatio //添加小B的结算比例
	data.Agent.Current = opreat
	data.Agent.OpreatRecord = append(data.Agent.OpreatRecord, opreat)

	Forward(session, "0", data)
}

//试用成为代理
func TryToUseAgent(session *JsNet.StSession) {
	type st_Get struct {
		UID string //uid
	}
	st := &st_Get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	data := &ST_User{}
	if err := WriteLock(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_User, st.UID, data)

	if data.Agent == nil {
		ForwardEx(session, "1", data, "TryToUseAgent failed,User=%s is not apply Agent", st.UID)
		return
	}
	if data.Agent.Current.OpreatStatus != constant.Agent_PassReviewe {
		ForwardEx(session, "1", data, "当前状态没法试用，status=%s\n", data.Agent.Current.OpreatStatus)
		return
	}
	data.AgentInfo = nil
	data.Agent.Bond = 0
	data.Agent.BondOrderID = ""
	data.Agent.PlatformFee = 0
	data.Agent.FranchiseFee = 0
	data.Agent.PayDate = ""
	data.Agent.OnlineDate = CurTime()
	opreat := ST_AgentStatus{
		OpreatStatus: constant.Agent_Online,
		OpreatTime:   CurTime(),
		OpreatReason: "试用上线",
		OpreatJobNum: data.UID,
		OpreatName:   data.Agent.Name,
	}
	data.Agent.Current = opreat
	data.Agent.OpreatRecord = append(data.Agent.OpreatRecord, opreat)
	fitAgentConfig(data)
	Forward(session, "0", data)
}

/////提交支付保证金的订单
func SubmitBondOrder(order *ST_Order) error {

	if order.OrderID == "" || order.UID == "" || order.PlatformFee < 0 || order.Bond < 0 { //|| cb == nil
		return ErrorLog("OrderPayBond failed,OrderID=%s,UID=%s,平台使用费=%d,保证金=%d\n",
			order.OrderID, order.UID, order.PlatformFee, order.Bond)
	}
	if order.AgentLevel != constant.Agent_Level_Diamonds_A {
		return ErrorLog("代理人的佣金层级不合法,AgentLevel=%s\n", order.AgentLevel)
	}

	user := &ST_User{}
	if err := WriteLock(constant.Hash_User, order.UID, user); err != nil {
		return err
	}
	defer WriteBack(constant.Hash_User, order.UID, user)

	if user.Agent == nil {
		return ErrorLog("SubmitBondOrder failed,User=%s is not apply Agent", order.UID)
	}
	if user.Agent.Current.OpreatStatus != constant.Agent_PassReviewe {
		if user.Agent.Current.OpreatStatus != constant.Agent_Offline_self {
			return ErrorLog("SubmitBondOrder failed,User=%s ,Agent status%s failed\n", order.UID, user.Agent.Current.OpreatStatus)
		}
	}
	user.Agent.AgentLevel = order.AgentLevel
	fitAgentConfig(user)

	user.Agent.PlatformFee = order.PlatformFee
	user.Agent.Bond = order.Bond
	opreat := ST_AgentStatus{
		OpreatStatus: user.Agent.Current.OpreatStatus,
		OpreatTime:   CurTime(),
		OpreatReason: constant.Agent_Submit,
		OpreatJobNum: order.UID,
		OpreatName:   user.Agent.Name,
	}
	user.AgentInfo = nil
	user.Agent.Current = opreat
	user.Agent.OpreatRecord = append(user.Agent.OpreatRecord, opreat)
	///////////////////更新代理费用订单id////////////////////
	user.Agent.BondOrderID = order.OrderID
	////////////////////订单添加到用户订单列表////////////////////
	user.Orders = append(user.Orders, order.OrderID)

	///////////////更新订单信息///////////////////////
	order.OrderSubmitDate = CurTime()
	order.SubmitStamp = CurStamp()
	order.OrderType = 2
	order.UserName = user.Name
	order.UserCell = user.Cell
	order.UserCity = user.City
	order.BusSource = user.BusSource
	order.AgentInfo = user.AgentInfo
	order.AgentCommission = 0 //上级代理的提成
	order.Amount = user.Agent.FranchiseFee + user.Agent.Bond
	order.RefundFee = user.Agent.Bond
	////////////////////////订单状态//////////////////////
	addOrderStatus(order,
		constant.Status_Order_PenddingPay, //待支付
		constant.Status_Order_PenddingPay, //待支付
		constant.Status_Order_PenddingPay, //待付款
		constant.Opreat_Order_Submit,      //用户提交订单
		"提交保证金订单",                         //用户提交订单
		constant.User_Side,                //用户
		order.UID,                         //用户id
		order.UserName)                    //用户姓名

	/////////////直接写///////////////////////
	return DirectWrite(constant.Hash_Order, order.OrderID, order)
}

//支付保证金订单
func OrderPayBond(OrderID string, cb *WxST_PayCb) (*ST_User, error) {
	if cb == nil {
		return nil, ErrorLog("OrderPayBond failed,WxST_PayCb is nil\n")
	}
	order := &ST_Order{}
	data := &ST_User{}
	if err := WriteLock(constant.Hash_Order, OrderID, order); err != nil {
		return data, err
	}
	defer WriteBack(constant.Hash_Order, OrderID, order)

	if order.Current.OpreatStatus != constant.Status_Order_PenddingPay {
		return data, ErrorLog("支付失败,当前状态不在未支付状态,OrderID=%s,Status=%s\n", OrderID, order.Current.OpreatStatus)
	}

	if err := WriteLock(constant.Hash_User, order.UID, data); err != nil {
		return data, err
	}
	defer WriteBack(constant.Hash_User, order.UID, data)

	///////////////////用户信息///////////////////////////////
	if data.Agent == nil {
		return data, ErrorLog("OrderPayBond failed,User=%s is not apply Agent", order.UID)
	}
	if (data.Agent.Current.OpreatStatus == constant.Agent_PassReviewe) ||
		(data.Agent.Current.OpreatStatus == constant.Agent_Offline_self &&
			((time.Now().Unix() - data.Agent.OfflineStamp) < 3600*24*365)) {
		data.Agent.PayDate = CurTime()
		data.Agent.OnlineDate = CurTime()
		opreat := ST_AgentStatus{
			OpreatStatus: constant.Agent_Online,
			OpreatTime:   CurTime(),
			OpreatReason: constant.Agent_Pay,
			OpreatJobNum: data.UID,
			OpreatName:   data.Agent.Name,
		}
		data.Agent.Current = opreat
		data.Agent.OpreatRecord = append(data.Agent.OpreatRecord, opreat)
		data.AgentInfo = nil
	} else {
		return data, ErrorLog("当前状态没法提交保证金，status=%s\n", data.Agent.Current.OpreatStatus)
	}

	///////////////////订单信息///////////////////////////////
	////支付信息/////
	order.WxPayCb = cb                            //微信支付回调
	order.PayMod = constant.PayMod_DepositPayment //支付方式，订金支付
	order.PayWay = constant.PayWay_Wechat         //微信支付
	order.PayDate = CurTime()                     //支付日期
	order.PayNumber = order.WxPayCb.Out_trade_no  //第三方支付的单号
	m, e := strconv.Atoi(cb.Cash_fee)
	if e == nil {
		order.RealPay = m ///支付价格
	}
	addOrderStatus(order,
		constant.Status_Order_Succeed, //完成
		constant.Status_Order_Succeed, //完成
		constant.Status_Order_Succeed, //完成
		constant.Opreat_Order_UserPay, //用户支付
		"支付保证金订单",                     //用户支付
		constant.User_Side,            //用户
		order.UID,                     //用户uid
		order.UserName)                //用户姓名

	return data, nil
}

///下线代理人(有保证金撤回)
func AgentOffline(OrderID, JodNumber, JobName, Reason string, refundCb map[string]string) error {
	if OrderID == "" || JodNumber == "" {
		return ErrorLog("AgentOffline failed,OrderID=%s,JobNumber=%s\n", OrderID, JodNumber)
	}
	order := &ST_Order{}
	if err := WriteLock(constant.Hash_Order, OrderID, order); err != nil {
		return err
	}
	defer WriteBack(constant.Hash_Order, OrderID, order)
	order.WxRefundCb = refundCb
	/////添加订单状态////////
	addOrderStatus(order,
		constant.Status_Order_AlreadyRefund,
		constant.Status_Order_AlreadyRefund,
		constant.Status_Order_AlreadyRefund,
		"代理被强制下钱", Reason,
		constant.Platform_Side, JodNumber,
		JobName)

	data := &ST_User{}
	if err := WriteLock(constant.Hash_User, order.UID, data); err != nil {
		return err
	}
	defer WriteBack(constant.Hash_User, order.UID, data)
	if data.Agent == nil {
		return ErrorLog("当前用户不是代理,%s", order.UID)
	}
	if data.Agent.Current.OpreatStatus != constant.Agent_Online {
		return ErrorLog("当前代理状态不能下线,UID=%s,Status=%s\n", order.UID, data.Agent.Current.OpreatStatus)
	}
	//////////////////添加状态///////////
	data.Agent.OfflineDate = CurTime()
	data.Agent.OfflineStamp = time.Now().Unix()
	opreat := ST_AgentStatus{
		OpreatStatus: constant.Agent_Offline_force,
		OpreatTime:   CurTime(),
		OpreatReason: Reason,
		OpreatJobNum: JodNumber,
		OpreatName:   JobName,
	}
	data.Agent.Current = opreat
	data.Agent.OpreatRecord = append(data.Agent.OpreatRecord, opreat)
	return nil
}

//后台强制上线代理（被下线的状态下）
func OnlineAgent(session *JsNet.StSession) {
	type st_online struct {
		UID string //用户Id
	}
	st := &st_online{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" {
		ForwardEx(session, "1", nil, "OnlineAgent failed,UID=%s,\n", st.UID)
		return
	}
	data := &ST_User{}
	if err := WriteLock(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, "OnlineAgent  WriteBack failed"+err.Error())
		return
	}
	defer WriteBack(constant.Hash_User, st.UID, data)
	if data.Agent == nil {
		ForwardEx(session, "1", data, "OnlineAgent failed,User=%s is not Agent", st.UID)
		return
	}
	if data.Agent.Current.OpreatStatus != constant.Agent_Offline_force && data.Agent.Current.OpreatStatus != constant.Agent_Offline_self {
		ForwardEx(session, "1", data, "当前代理状态不能手动上线,UID=%s,Status=%s\n", st.UID, data.Agent.Current.OpreatStatus)
		return
	}
	//////////////////添加状态///////////
	data.Agent.OnlineDate = CurTime()
	op := ST_AgentStatus{
		OpreatStatus: constant.Agent_Online,
		OpreatTime:   CurTime(),
		OpreatReason: "后台强制上线",
		OpreatJobNum: "Admin",
		OpreatName:   "Admin",
	}
	data.Agent.Current = op
	data.Agent.OpreatRecord = append(data.Agent.OpreatRecord, op)
	Forward(session, "0", data)
}

////下线代理人（后台强制）
func OfflineAgent(session *JsNet.StSession) {
	type st_online struct {
		UID       string //用户Id
		JodNumber string //审核人的工号
		JobName   string //审核人姓名
		Reason    string //审核备注
	}
	st := &st_online{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" || st.JodNumber == "" {
		ForwardEx(session, "1", nil, "OfflineAgent failed,UID=%s,JobNumber=%s\n", st.UID, st.JodNumber)
		return
	}
	data := &ST_User{}
	if err := WriteLock(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, "OfflineAgent  WriteBack failed"+err.Error())
		return
	}
	defer WriteBack(constant.Hash_User, st.UID, data)
	if data.Agent == nil {
		ForwardEx(session, "1", data, "OfflineAgent failed,User=%s is not Agent", st.UID)
		return
	}
	if data.Agent.Current.OpreatStatus != constant.Agent_Online {
		ForwardEx(session, "1", data, "当前代理状态不能下线,UID=%s,Status=%s\n", st.UID, data.Agent.Current.OpreatStatus)
		return
	}
	//////////////////添加状态///////////
	data.Agent.OfflineDate = CurTime()
	data.Agent.OfflineStamp = time.Now().Unix()
	op := ST_AgentStatus{
		OpreatStatus: constant.Agent_Offline_force,
		OpreatTime:   CurTime(),
		OpreatReason: st.Reason,
		OpreatJobNum: st.JodNumber,
		OpreatName:   st.JobName,
	}
	data.Agent.Current = op
	data.Agent.OpreatRecord = append(data.Agent.OpreatRecord, op)
	Forward(session, "0", data)
}

//撤回保证金
func WithdrawBond(OrderID string, refundCb map[string]string) (*ST_User, error) {
	if OrderID == "" {
		return nil, ErrorLog("WithdrawBond failed,OrderID is empty\n")
	}
	order := &ST_Order{}
	data := &ST_User{}
	if err := WriteLock(constant.Hash_Order, OrderID, order); err != nil {
		return data, err
	}
	defer WriteBack(constant.Hash_Order, OrderID, order)
	order.WxRefundCb = refundCb
	/////添加订单状态////////
	addOrderStatus(order,
		constant.Status_Order_AlreadyRefund,
		constant.Status_Order_AlreadyRefund,
		constant.Status_Order_AlreadyRefund,
		"", "", constant.User_Side,
		order.UID, order.UserName)

	if err := WriteLock(constant.Hash_User, order.UID, data); err != nil {
		return data, err
	}
	defer WriteBack(constant.Hash_User, order.UID, data)
	if data.Agent == nil {
		return data, ErrorLog("WithdrawBond failed,User=%s is not Agent", order.UID)
	}
	if data.Agent.Current.OpreatStatus != constant.Agent_Online {
		return data, ErrorLog("WithdrawBond failed,当前状态无法退款,Status=%s\n", data.Agent.Current.OpreatStatus)
	}
	//////////////////添加状态///////////
	data.Agent.Bond = 0
	data.Agent.WithdrawDate = CurTime()
	data.Agent.OfflineDate = CurTime()
	data.Agent.OfflineStamp = time.Now().Unix()
	opreat := ST_AgentStatus{
		OpreatStatus: constant.Agent_Offline_self,
		OpreatTime:   CurTime(),
		OpreatReason: constant.Agent_WithDraw,
		OpreatJobNum: data.UID,
		OpreatName:   data.Name,
	}
	data.Agent.Current = opreat
	data.Agent.OpreatRecord = append(data.Agent.OpreatRecord, opreat)
	return data, nil
}

//余额体现
func WithdrawBalance(UID string, money int) (*ST_User, error) {
	if UID == "" || money < 100 {
		return nil, ErrorLog("WithdrawBalance param failed,UID=%s,money=%d\n", UID, money)
	}
	user := &ST_User{}
	if err := WriteLock(constant.Hash_User, UID, user); err != nil {
		return user, err
	}
	defer WriteBack(constant.Hash_User, UID, user)
	if user.Agent == nil {
		return user, ErrorLog("当前用户不是小B,UID=%s\n", UID)
	}
	if user.Agent.Balance < money {
		return user, ErrorLog("当前用户余额不够提现,UID=%s,Balance=%d\n", UID, user.Agent.Balance)
	}
	withdraw := user.Agent.Withdraw + money
	Balance := user.Agent.TotalIncome - withdraw

	wd, e := NewWithDraw(user, money, Balance, 0)
	if e != nil {
		return user, e
	}
	user.Agent.Withdraw = withdraw
	user.Agent.Balance = Balance
	user.Agent.WithDrawRecord = append(user.Agent.WithDrawRecord, wd.WDID)
	return user, nil
}

//余额提现失败
func WithdrawBalanceFailed(UID string, money int) (*ST_User, error) {
	if UID == "" || money < 0 {
		return nil, ErrorLog("WithdrawBalanceFailed param failed,UID=%s,money=%d\n", UID, money)
	}
	user := &ST_User{}
	if err := WriteLock(constant.Hash_User, UID, user); err != nil {
		return user, err
	}
	defer WriteBack(constant.Hash_User, UID, user)
	if user.Agent == nil {
		return user, ErrorLog("当前用户不是小B或者店员,UID=%s\n", UID)
	}
	if user.Agent.Withdraw < money {
		return user, ErrorLog("user.Agent.Withdraw < money \n")
	}
	user.Agent.Withdraw -= money
	user.Agent.Balance = user.Agent.TotalIncome - user.Agent.Withdraw
	return user, nil
}

///新绑用户的代理
func UserBindAgent(session *JsNet.StSession) {
	type st_get struct {
		UID     string //用户
		AgentID string //代理id
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" || st.AgentID == "" {
		ForwardEx(session, "1", nil, "UserBindAgent failed ,param is empty,UID=%s,AgentID=%s\n", st.UID, st.AgentID)
		return
	}
	if st.AgentID == st.UID {
		ForwardEx(session, "1", nil, "自己不能邀请自己作为用户,UID=%s,AgentID=%s\n", st.UID, st.AgentID)
		return
	}
	Agent := &ST_User{}
	data := &ST_User{}

	if err := WriteLock(constant.Hash_User, st.AgentID, Agent); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_User, st.AgentID, Agent)
	if Agent.Agent == nil {
		ForwardEx(session, "1", nil, "user :%s is not agent \n", st.AgentID)
		return
	}
	if Agent.Agent.Current.OpreatStatus != constant.Agent_Online {
		ForwardEx(session, "1", nil, "当前小B,不是在线状态，不能邀请用户,UID:%s,Status=%s \n",
			Agent.UID, Agent.Agent.Current.OpreatStatus)
		return
	}
	if err := WriteLock(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_User, st.UID, data)

	if data.Agent != nil && data.Agent.Current.OpreatStatus == constant.Agent_Online {
		ForwardEx(session, "1", nil, "当前用户UID:%s已经是在线小B，不能接受邀请\n", data.UID)
		return
	}
	data.AgentInfo = &ST_AgentSimpleInfo{
		UID:      Agent.UID,  //用户ID
		Name:     Agent.Name, //姓名
		City:     Agent.City, //城市
		Cell:     Agent.Cell, //联系方式
		ProRatio: Agent.Agent.ProRatio,
	}
	if Agent.Agent.AgentType == constant.Agent_Type_Physical {
		data.BusSource = Agent.Agent.ShopBussiness
	} else {
		data.BusSource = Agent.Agent.WxItem
	}
	////////添加一个顾客///////////////////
	exist := false
	for _, v := range Agent.Agent.Customer {
		if v.UID == data.UID {
			exist = true
			break
		}
	}
	if !exist {
		Agent.Agent.Customer = append(Agent.Agent.Customer,
			ST_TeamInfo{
				UID:      data.UID,
				Name:     data.Name,
				HeadPic:  data.HeadImageURL,
				JoinDate: CurTime(),
			})
	}
	if data != nil && Agent != nil {
		go AgentInviteStatistics(ST_AgentInviteStatistic{
			UID:       data.UID,
			UserName:  data.Name,
			AgentID:   Agent.UID,
			AgentCity: Agent.Agent.City,
			ApplyDate: CurDate(),
			CreatDate: data.CreaDate,
		})
	}
	Forward(session, "0", Agent)
}

//更新用户代理信息
func updateUserAgent(UID string, agent *ST_UserAgent) (*ST_User, error) {
	if UID == "" {
		return nil, ErrorLog("updateUserAgent failed,UID is empty\n")
	}
	data := &ST_User{}
	if err := WriteLock(constant.Hash_User, UID, data); err != nil {
		return data, err
	}
	data.Agent = agent
	return data, WriteBack(constant.Hash_User, UID, data)
}

///添加代理的相关订单信息
func pushAgentOrder(order *ST_Order) error {
	if order.AgentInfo != nil && order.AgentInfo.UID != "" {
		agent := &ST_User{}
		if err := WriteLock(constant.Hash_User, order.AgentInfo.UID, agent); err != nil {
			return err
		}
		defer WriteBack(constant.Hash_User, order.AgentInfo.UID, agent)
		if agent.Agent == nil || agent.Agent.Current.OpreatStatus != constant.Agent_Online {
			return ErrorLog("当前小B,不是在线状态，不能推送,UID:%s,Status=%s \n",
				agent.UID, agent.Agent.Current.OpreatStatus)
		}
		updateAgentOrder(agent, order.OrderID, CurTime(), order.AgentCommission)
		updateCustomerOrder(agent, order.OrderID, CurTime(), order.UID, order.AgentCommission)
	} else {
		Info("%s没有上级代理\n", order.UID)
	}
	return nil
}

///追加一个相关的订单id
func updateAgentOrder(data *ST_User, OrderID, Date string, Commission int) {
	Info("updateAgentOrder...\n")
	if data.Agent == nil {
		Error("updateAgentOrder error,Agent is nil \n")
		return
	}
	exist := false
	com := 0
	in := 0
	for i, v := range data.Agent.Orders {
		if v.OrderID == OrderID {
			exist = true
			data.Agent.Orders[i].LastCommisson = Commission
			data.Agent.Orders[i].LastDate = Date
		}
		com += data.Agent.Orders[i].PreCommisson
		in += data.Agent.Orders[i].LastCommisson
	}
	if !exist {
		if data.Agent.Orders == nil || len(data.Agent.Orders) == 0 {
			data.Agent.Orders = []ST_CommissionOrder{}
		}
		data.Agent.Orders = append(data.Agent.Orders, ST_CommissionOrder{
			OrderID:       OrderID,
			PreCommisson:  Commission,
			LastCommisson: 0,
			PreDate:       Date,
		})
		com += Commission
	}
	data.Agent.ToltleCommission = com
	data.Agent.TotalIncome = in
	data.Agent.Balance = data.Agent.TotalIncome - data.Agent.Withdraw
}

//添加一个用户的相关订单
func updateCustomerOrder(data *ST_User, OrderID, Date, CusID string, Commission int) {
	Info("updateCustomerOrder...\n")
	if data.Agent == nil {
		Error("updateCustomerOrder error,Agent is nil \n")
		return
	}
	for i, v := range data.Agent.Customer {
		if v.UID == CusID {
			exist := false
			com := 0
			in := 0
			for j, v1 := range v.Orders {
				if v1.OrderID == OrderID {
					exist = true
					data.Agent.Customer[i].Orders[j].LastCommisson = Commission
					data.Agent.Customer[i].Orders[j].LastDate = Date
				}
				com += data.Agent.Customer[i].Orders[j].PreCommisson
				in += data.Agent.Customer[i].Orders[j].LastCommisson
			}
			if !exist {
				if data.Agent.Customer[i].Orders == nil || len(data.Agent.Customer[i].Orders) == 0 {
					data.Agent.Customer[i].Orders = []ST_CommissionOrder{}
				}
				data.Agent.Customer[i].Orders = append(data.Agent.Customer[i].Orders,
					ST_CommissionOrder{
						OrderID:       OrderID,
						PreCommisson:  Commission,
						LastCommisson: 0,
						PreDate:       Date,
					})
				com += Commission
			}
			data.Agent.Customer[i].ToltleCommission = com
			data.Agent.Customer[i].TotalIncome = in
			break
		}
	}
}

func AgentSharedArticle(UID, ArticleID string) (error, *ST_User) {
	if UID == "" || ArticleID == "" {
		return ErrorLog("参数不完整UID=%s,ArticleID=%s\n", UID, ArticleID), nil
	}
	user := &ST_User{}
	if err := WriteLock(constant.Hash_User, UID, user); err != nil {
		return err, user
	}
	defer WriteBack(constant.Hash_User, UID, user)

	if user.Agent == nil {
		return ErrorLog("当前用户:%s不是代理无法进行分享\n", UID), user
	}
	if user.Agent.Current.OpreatStatus != constant.Agent_Online {
		return ErrorLog("当前小B,不是在线状态,不能推送分享,UID:%s,Status=%s \n",
			UID, user.Agent.Current.OpreatStatus), user
	}
	exist := false
	for i, v := range user.Agent.Article {
		if v.ArticleID == ArticleID {
			user.Agent.Article[i].ShareTime = CurTime()
			exist = true
			break
		}
	}
	if !exist {
		if user.Agent.Article == nil || len(user.Agent.Article) == 0 {
			user.Agent.Article = []ST_AgentArticle{}
		}
		user.Agent.Article = append(user.Agent.Article, ST_AgentArticle{
			ArticleID: ArticleID,
			ShareTime: CurTime(),
		})
	}
	return nil, user
}

func checkAgent(agent *ST_UserAgent) error {
	if agent.UID == "" {
		return ErrorLog("代理人的UD不能为空\n")
	}
	if agent.Name == "" {
		return ErrorLog("代理人的UD或者姓名不能为空\n")
	}
	if agent.Cell == "" {
		return ErrorLog("代理人的手机号不能为空\n")
	}
	if agent.City == "" {
		return ErrorLog("代理人的城市不能为空\n")
	}
	//if agent.IdentityNumber == "" || agent.IdentityPic == "" {
	//	return ErrorLog("代理里人的身份证号码或者身份证照片不能为空\n")
	//}
	return nil
}

///检查小B是否下线超过1年，如果超过，就将小B资料删除
func checkUserAgent(user *ST_User) {
	//if user.Agent != nil && (user.Agent.Current.OpreatStatus == constant.Agent_Offline_force ||
	//	user.Agent.Current.OpreatStatus == constant.Agent_Offline_self) {
	//	if time.Now().Unix()-user.Agent.OfflineStamp > 3600*24*365 {
	//		updateUserAgent(user.UID, nil)
	//	}
	//}
}

func fitAgentConfig(data *ST_User) {
	if data.Agent.AgentLevel == constant.Agent_Level_Diamonds_A {
		data.Agent.FranchiseFee = JsConfig.CFG.Agent.Diamonds_A.FranchiseFee
	}
	if data.Agent.AgentLevel == constant.Agent_Level_TryUse {
		data.Agent.FranchiseFee = JsConfig.CFG.Agent.TryUse.FranchiseFee
	}
}

package common

import (
	"JsGo/JsMobile"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	// . "cache/cacheIO"

	"constant"
	"strconv"
	. "util"
)

////微信支付回调
type WxST_PayCb struct {
	AppId                string `xml:"appid"`
	Mch_id               string `xml:"mch_id"`
	Device_info          string `xml:"device_info"`
	Nonce_str            string `xml:"nonce_str"`
	Sign                 string `xml:"sign"`
	Sign_type            string `xml:"sign_type"`
	Result_code          string `xml:"result_code"`
	Err_code             string `xml:"err_code"`
	Err_code_des         string `xml:"err_code_des"`
	Openid               string `xml:"openid"`
	Is_subscribe         string `xml:"is_subscribe"`
	Trade_type           string `xml:"trade_type"`
	Bank_type            string `xml:"bank_type"`
	Total_fee            string `xml:"total_fee"`
	Settlement_total_fee string `xml:"settlement_total_fee"`
	Fee_type             string `xml:"fee_type"`
	Cash_fee             string `xml:"cash_fee"`
	Cash_fee_type        string `xml:"cash_fee_type"`
	Transaction_id       string `xml:"transaction_id"`
	Out_trade_no         string `xml:"out_trade_no"` //订单id
	Attach               string `xml:"attach"`
	Time_end             string `xml:"time_end"`
}

//订单价格信息
type ST_OrderPriceInfo struct {
	CouponPrice         int               //使用的优惠券的价格
	TotalPrice          int               //消费总金额
	XYCoin              int               //使用的星喜币
	RedPrice            int               //红包价格
	HosPayPrice         int               //到医院需要支付的价格
	AppendPrice         int               //新增加的项目金额
	HosPayRealPrice     int               //到院实际支付的价格
	FullPayReturnPrice  int               //全款支付返还医院金额
	HosSettlementRatio  int               //跟医院的产品结算比例
	HosSettlementPrice  int               //跟医院的产品结算价格
	RealSettlementPrice int               //实际计算的价格
	PayMod              string            //定金支付/全款付款
	PayWay              string            //支付途径:支付宝、微信、银联、平台余额（保留）
	PayAccount          int               //支付账号
	PayNumber           string            //第三方平台支付单号
	WxPayCb             *WxST_PayCb       //支付回调
	WxRefundCb          map[string]string //退款回调
	RefundMoney         int               //退款金额
	Charge              map[string]string //票据
}

//订单用户信息
type ST_OrderUserInfo struct {
	UID       string // 用户id
	UserName  string // 用户名字
	UserCell  string // 用户手机号
	UserCity  string // 用户城市
	VerifCode string // 用户的校验码
	BusSource string //

	AgentInfo       *ST_AgentSimpleInfo // 上级代理信息
	AgentCommission int                 // 上代理的提成
	AgentLevel      string              // 佣金层级
	PlatformFee     int                 // 平台使用费
	Bond            int                 // 保证金
	FranchiseFee    int                 // 加盟费
}

//订单的产品信息
type ST_OrderProInfo struct {
	ProID        string //产品id
	ProName      string //产品名称
	ProFirstItem string //产品的一级分类
	DocName      string //医生姓名
	DocID        string //医生Id
	ProType      string //产品类型:常规产品、特价产品、联合定制
	XingYaoPrice int    //产品价格
	ProDeposit   int    //产品定金
}

///分期
type ST_OrderInstalmentInfo struct {
	InstalmentID     string // 此次分期id
	InstalmentProID  string // 分期产品的id
	IsInstalment     string // 是否分期
	InstallmentCom   string // 分期公司
	InstallmentPrice string // 分期金额
	InstallmentNums  string // 分期期数
}

//保险
type ST_OrderInsurance struct {
	InsuranceProID string // 保险产品id
	InsuranceID    string // 此次保险id
	IsInsurance    string // 是否购买保险
	InsuranceType  string // 保险种类
	InsurancePrice int    // 保险保费
}

//订单的流程信息
type ST_OrderFlow struct {
	OpreatPart        string //操作方:用户、医院、平台
	OpreatStatus      string //系统订单状态
	OpreatUserStatus  string //用户的订单状态
	OpreatAgentStatus string //代理的订单状态
	OpreatTime        string //操作时间
	OpreatReason      string //操作原因
	OpreatAction      string //操作动作
	OpreatJobNum      string //操作人员工号
	OpreatName        string //操作人员姓名
}

//订单的状态信息
type OrderStatusInfo struct {
	SubmitStamp     int64          //下单的时间戳
	OrderSubmitDate string         //下单时间
	PayDate         string         //支付时间
	CheckDate       string         //校验时间
	CloseDate       string         //关闭时间
	RefundDate      string         //退款时间
	AppointmentDate string         //预约时间
	AppointmentDes  string         //预约说明
	Opreat          []ST_OrderFlow //状态变化流
	Current         ST_OrderFlow   //当前状态
}

///支付信息
type ST_OrderApp struct {
	TerminalIp       string // 支付主机IP
	LocalTimeStamp   int64  // 本地时间戳
	ServiceTimeStamp int64  // 服务端时间戳
	Amount           int    // 需要支付金额
	RealPay          int    // 实际支付价格
	Desc             string // 描述
	Nonce_str        string // 随机串
	Mch_id           string // 商家ID`
	AppId            string // 应用ID
	OpenId           string // opendid
	RefundId         string //退款ID
	RefundFee        int    //退款金额
}

//订单
type ST_Order struct {
	OrderID                string // 订单id
	OrderType              int    // 1-2-4
	HosID                  string // 医院id
	HosName                string // 医院名称
	ST_OrderApp                   //支付信息
	ST_OrderUserInfo              //用户信息
	ST_OrderProInfo               //订单的产品信息
	ST_OrderPriceInfo             //订单的价格信息
	ST_OrderInstalmentInfo        //产品的分期信息
	ST_OrderInsurance             //产品的保险信息
	OrderStatusInfo               //订单的状态信息
}

////用户提交订单，创建一个新的订单
func SubmitOrder(order *ST_Order) error {
	Info("Enter submint order")
	if order.UID == "" || order.OrderID == "" || order.ProID == "" {
		return ErrorLog("SubmitOrder param is empty,UID=%s,OrderID=%s,ProID=%s\n", order.UID, order.OrderID, order.ProID)
	}
	order.OrderSubmitDate = CurTime()
	order.SubmitStamp = CurStamp()
	////////////从数据库中获取用户信息///////////////////
	user, err := GetUserInfo(order.UID)
	if err != nil {
		return ErrorLog("创建订单失败,获取用户信息失败,UID=%s\n", order.UID)
	}
	order.OrderType = 1
	order.UserName = user.Name
	order.UserCell = user.Cell
	order.UserCity = user.City
	order.BusSource = user.BusSource
	order.AgentInfo = user.AgentInfo
	order.AgentCommission = 0 //上级代理的提成

	///////////////////////订单中的产品////////////////////////////////////
	pro, e := GetProductInfo(order.ProID)
	if e != nil {
		return ErrorLog("创建订单失败,获取产品信息失败,ProID=%s\n", order.ProID)
	}
	hos, e := GetHospitalInfo(pro.HosID)
	if e != nil {
		return ErrorLog("创建订单失败,获取医院信息失败,HosID=%s\n", pro.HosID)
	}
	order.HosID = pro.HosID               //医院id
	order.HosName = hos.HosName           //医院名称
	order.DocID = pro.Doctors             //医生id
	order.DocName = pro.DoctorName        //医生名称
	order.ProName = pro.ProName           //产品名字
	order.ProFirstItem = pro.FirstItem    //产品的一级菜单
	order.ProType = pro.ProType           //产品类型
	order.ProDeposit = pro.ProDeposit     //产品的定金
	order.XingYaoPrice = pro.XingYaoPrice //产品的星喜价格
	order.Amount = order.ProDeposit       ///支付价格等于产品定金

	//////////////////计算产品的结算比例///////////////////////////////
	order.HosSettlementRatio = pro.ProRatio //产品结算比例
	////////////////////////订单状态//////////////////////
	addOrderStatus(order,
		constant.Status_Order_PenddingPay, //待使用
		constant.Status_Order_PenddingPay, //代理待结算
		constant.Status_Order_PenddingPay, //待付款
		constant.Opreat_Order_Submit,      //用户提交订单
		constant.Opreat_Order_Submit,      //用户提交订单
		constant.User_Side,                //用户
		order.UID,                         //用户id
		order.UserName)                    //用户姓名
	/////////////直接写///////////////////////
	if err := DirectWrite(constant.Hash_Order, order.OrderID, order); err != nil {
		return ErrorLog("创建订单失败,DirectWrite():%s\n", err.Error())
	}
	//////////添加一个订单到用户///////////////
	go NewUserOrder(order.UID, order.OrderID)
	return nil
}

///用户已付款，状态：待付款->待预约
func OrderUserPaid(OrderID string, cb *WxST_PayCb) (*ST_Order, error) {
	Info("Enter OrderUserPaid .............................")

	if OrderID == "" {
		return nil, ErrorLog("OrderUserPaid faild,OrderID=%s\n", OrderID)
	}
	if cb == nil {
		return nil, ErrorLog("OrderUserPaid,WxST_PayCb is nil\n")
	}

	data := &ST_Order{}
	if err := WriteLock(constant.Hash_Order, OrderID, data); err != nil {
		return data, err
	}
	defer WriteBack(constant.Hash_Order, OrderID, data)
	if data.Current.OpreatStatus != constant.Status_Order_PenddingPay {
		return data, ErrorLog("支付失败,当前状态不在未支付状态,OrderID=%s,Status=%s\n", OrderID, data.Current.OpreatStatus)
	}
	///////////////////////结算信息////////////////////////////////
	orderSettlement(data)

	////支付信息/////
	data.WxPayCb = cb //微信支付回调
	m, e := strconv.Atoi(cb.Cash_fee)
	if e == nil {
		data.RealPay = m ///支付价格
	}
	data.RefundFee = data.RealPay                ///支付价格
	data.PayNumber = data.WxPayCb.Out_trade_no   //第三方支付的单号
	data.PayMod = constant.PayMod_DepositPayment //支付方式，订金支付
	data.PayWay = constant.PayWay_Wechat         //微信支付
	data.PayDate = CurTime()                     //支付日期
	/////////////订单状态///////////////
	addOrderStatus(data,
		constant.Status_Order_UserPenddingUse,     //待使用
		constant.Status_Order_PendingStatements,   //代理待结算
		constant.Status_Order_PenddingAppointment, //待预约
		constant.Opreat_Order_UserPay,             //用户支付
		constant.Opreat_Order_UserPay,             //用户支付
		constant.User_Side,                        //用户
		data.UID,                                  //用户uid
		data.UserName)                             //用户姓名
	if data != nil {
		///////订单统计信息//////////
		go order_statistics(data.HosID, data.DocID, data.ProID, data.OrderID, data.XingYaoPrice, data.ProDeposit, true)
		//////添加一个订单到全局///////////////
		go addOrderToGlobal(OrderID)

		city := "***"
		name := "***"
		cell := "***"
		if data.AgentInfo != nil {
			city = data.AgentInfo.City
			name = data.AgentInfo.Name
			cell = data.AgentInfo.Cell
		}
		go smsCustom(city, name, cell, data.UserName, data.UserCell, data.HosName, data.ProName, data.DocName,data.XingYaoPrice)
	}
	return data, nil
}

//客户下单后，短信通知客服
func smsCustom(City, Name, Cell, UserName, UserCell, HosName, ProName, DocName string,price int) {
	par := make(map[string]string)
	par["city"] = City
	par["name1"] = Name
	par["phone1"] = Cell
	par["name2"] = UserName
	par["phone2"] = UserCell
	par["hosp"] = HosName
	par["prod"] = ProName
	par["name3"] = DocName
	par["price"] = strconv.Itoa(price)+""

	 JsMobile.ComJsMobileVerify("喜妹儿", "13961150133", "SMS_130922367", "c", 60, par)
	 JsMobile.ComJsMobileVerify("喜妹儿", "18616215779", "SMS_130922367", "c", 60, par)
	 JsMobile.ComJsMobileVerify("喜妹儿", "13310030877", "SMS_130922367", "c", 60, par)
}

///后台预约，状态从待预约到待校验
func OrderAppointment(OrderID, OpreatJobNum, OpreatName, AppointmentDate, AppointmentDes string) (*ST_Order, error) {

	if OrderID == "" || OpreatJobNum == "" || OpreatName == "" || AppointmentDate == "" {
		return nil, ErrorLog("预约失败,参数不全,OrderID=%s,OpreatJobNum=%s,OpreatName=%s,AppointmentDate=%s",
			OrderID, OpreatJobNum, OpreatName, AppointmentDate)
	}
	data := &ST_Order{}
	if err := WriteLock(constant.Hash_Order, OrderID, data); err != nil {
		return data, err
	}
	defer WriteBack(constant.Hash_Order, OrderID, data)
	if data.Current.OpreatStatus != constant.Status_Order_PenddingAppointment &&
		data.Current.OpreatStatus != constant.Status_Order_PenddingVerify {
		return data, ErrorLog("预约失败,当前状态不在待预约状态状态,OrderID=%s,Status=%s\n", OrderID, data.Current.OpreatStatus)
	}
	////生成校验码
	var IsE error = nil
	data.VerifCode, IsE = NewVerifyCode(data.OrderID, data.HosID, data.UID)
	if IsE != nil {
		return data, IsE
	}
	data.AppointmentDate = AppointmentDate
	data.AppointmentDes = AppointmentDes
	/////////////订单状态///////////////
	addOrderStatus(data,
		constant.Status_Order_UserPenddingUse,   //待使用
		constant.Status_Order_PendingStatements, //代理待结算
		constant.Status_Order_PenddingVerify,    //待校验
		constant.Opreat_Order_SysAppoint,        //系统预约
		constant.Opreat_Order_SysAppoint,        //系统预约
		constant.Platform_Side,                  //用户
		OpreatJobNum,                            //用户uid
		OpreatName)                              //用户姓名

	if data != nil {
		go AppendHosOrder(data.HosID, OrderID)
	}
	return data, nil
}

///用户到医院验证，状态：待使用（待验证）->已使用（已验证）
func OrderHosVerify(HosID, VerifCode string) (*ST_Order, error) {

	if HosID == "" || VerifCode == "" {
		return nil, ErrorLog("验证失败,参数不全,HosID=%s,VerifCode=%s,OpreatJobNum=%s\n",
			HosID, VerifCode)
	}
	OrderID, e := CheckVerifyCode(VerifCode, HosID)
	if e != nil {
		return nil, ErrorLog("订单验证失败CheckVerifyCode(),VerifCode=%s\n", VerifCode)
	}
	data := &ST_Order{}

	if err := WriteLock(constant.Hash_Order, OrderID, data); err != nil {
		return data, err
	}
	defer WriteBack(constant.Hash_Order, OrderID, data)
	if data.Current.OpreatStatus != constant.Status_Order_PenddingVerify {
		return data, ErrorLog("校验失败,当前状态不在待校验状态,OrderID=%s,Status=%s\n", OrderID, data.Current.OpreatStatus)

	}
	if HosID != data.HosID || VerifCode != data.VerifCode {
		return data, ErrorLog("订单验证失败，订单的医院id或者校验码不正确！\n")
	}
	/////////////订单状态///////////////
	data.CheckDate = CurTime()
	addOrderStatus(data,
		constant.Status_Order_PenddingEvaluate, //待评价
		constant.Status_Order_PendingConfirm,   //待确认
		constant.Status_Order_PendingConfirm,   //待确认
		"医院校验",
		"医院校验",
		constant.Hospital_Side, //医院
		data.HosID,             //员工工号
		data.HosName)           //员工姓名
	//删除对应的校验码
	go DelVerifyCode(VerifCode)
	/////推送代理的订单////////
	if data != nil {
		go pushAgentOrder(data)
		go appendOrderToProduct(data.ProID, data.OrderID)
	}
	return data, nil
}

//医院确认手术，并且追加金额
func HosConfirmOperation(HosID, OrderID, Note string, price int) error {
	if HosID == "" || OrderID == "" || price < 0 {
		return ErrorLog("HosConfirmOperation error,HosID=%s,OrderID=%s,price=%d\n", HosID, OrderID, price)
	}
	data := &ST_Order{}
	if err := WriteLock(constant.Hash_Order, OrderID, data); err != nil {
		return err
	}
	if data.Current.OpreatStatus != constant.Status_Order_PendingConfirm && data.Current.OpreatStatus != constant.Status_Order_PendingStatements {
		return ErrorLog("HosConfirmOperation error,OrderID=%s 状态不对,status=%s\n", OrderID, data.Current.OpreatStatus)
	}

	data.AppendPrice = price
	///////////////////////重新结算////////////////////////////////
	orderSettlement(data)

	addOrderStatus(data,
		constant.Status_Order_PenddingEvaluate,  //待评价
		constant.Status_Order_PendingStatements, //待结算
		constant.Status_Order_PendingStatements, //待结算
		"确认手术，并追加金额",
		Note, constant.Hospital_Side, //医院
		HosID, "") //员工姓名
	return WriteBack(constant.Hash_Order, OrderID, data)
}

///限时的产品到期，状态：待付款->已失效（已取消）
func OrderInvalid(OrderID string) (*ST_Order, error) {
	Info("OrderInvalid ...OrderID=%s\n", OrderID)
	if OrderID == "" {
		return nil, ErrorLog("订单无效失败,OrderID=%s\n", OrderID)
	}
	data := &ST_Order{}
	if err := WriteLock(constant.Hash_Order, OrderID, data); err != nil {
		return data, err
	}
	defer WriteBack(constant.Hash_Order, OrderID, data)
	if data.Current.OpreatStatus != constant.Status_Order_PenddingPay {
		return data, ErrorLog("订单失效操作失败,状态不对,OrderID=%s,Status=%s\n", OrderID, data.Current.OpreatStatus)
	}
	data.CloseDate = CurTime()
	///////////订单状态///////////////
	addOrderStatus(data,
		constant.Status_Order_Invalid,
		constant.Status_Order_Invalid,
		constant.Status_Order_Invalid,                                //过期取消
		constant.Opreat_Order_Invalid, constant.Opreat_Order_Invalid, //过期取消
		constant.Platform_Side, "", "")

	//////////////更新订单状态////////////
	if data != nil {
		//////从用户正在进行的订单中移除
		go RemoveUserContinueOrder(data.UID, data.OrderID)
	}
	return data, nil
}

///退款
func OrderRefund(order *ST_Order, refundCb map[string]string, Msg, JobNum, JobName string) error {

	if err := WriteLock(constant.Hash_Order, order.OrderID, order); err != nil {
		return err
	}
	order.WxRefundCb = refundCb
	if err := WriteBack(constant.Hash_Order, order.OrderID, order); err != nil {
		return err
	}
	if _, err := OrderCancle(order.OrderID, Msg, JobNum, JobName); err != nil {
		return err
	}
	return nil
}

//订单取消
func OrderCancle(OrderID, Msg, JobNum, Name string) (*ST_Order, error) {
	if OrderID == "" || JobNum == "" {
		return nil, ErrorLog("OrderCancle failed,OrderID =%s,JobNum=%s\n", OrderID, JobNum)
	}
	list := []string{}
	data := &ST_Order{}

	if err := WriteLock(constant.Hash_Order, OrderID, data); err != nil {
		return data, err
	}
	list = append(list, data.UID)
	if data.Current.OpreatStatus == constant.Status_Order_PenddingPay {
		//////从用户正在进行的订单中移除
		if err := RemoveUserContinueOrder(data.UID, data.OrderID); err != nil {
			ErrorLog("RemoveUserContinueOrder,HosID=%s,DocID=%s\n", data.HosID, data.DocID)
		}
		addOrderStatus(data,
			constant.Status_Order_Cancle,
			constant.Status_Order_Cancle,
			constant.Status_Order_Cancle,
			constant.Opreat_Order_UserCancleBeforeVerify,
			Msg, constant.User_Side, JobNum, Name) //用户姓名
		data.RefundDate = ""
	} else if data.Current.OpreatStatus == constant.Status_Order_PenddingAppointment {
		addOrderStatus(data,
			constant.Status_Order_AlreadyRefund,
			constant.Status_Order_CancleBeforeAppointment,
			constant.Status_Order_CancleBeforeAppointment,
			constant.Opreat_Order_UserCancleBeforeVerify,
			Msg, constant.User_Side, JobNum, Name) //用户姓名
		list = append(list, "Admin")

	} else if data.Current.OpreatStatus == constant.Status_Order_PenddingVerify {
		addOrderStatus(data,
			constant.Status_Order_AlreadyRefund,
			constant.Status_Order_CancleBeforeVerfy,
			constant.Status_Order_CancleBeforeVerfy,
			constant.Opreat_Order_UserCancleBeforeVerify,
			Msg, constant.User_Side, JobNum, Name) //用户姓名
		////////////更新医院的订单/////////////////
		list = append(list, data.HosID)
		list = append(list, "Admin")
	} else {
		/////添加订单状态////////
		addOrderStatus(data,
			constant.Status_Order_AlreadyRefund,
			constant.Status_Order_CancleAfterVerfy,
			constant.Status_Order_CancleAfterVerfy,
			constant.Opreat_Order_UserCancleAfterVerify,
			Msg, constant.Platform_Side, JobNum, Name) //用户姓名
		////////////////代理订单的更新/////////////////
		if data.AgentInfo != nil && data.AgentInfo.UID != "" {
			list = append(list, data.AgentInfo.UID)
		}
	}
	data.CloseDate = CurTime()
	data.RefundDate = CurTime()

	if err := WriteBack(constant.Hash_Order, OrderID, data); err != nil {
		return data, err
	}
	//////////////更新订单状态////////////
	if data != nil {
		/////更新医院的退款信息//////
		go hospital_refundrate(data.HosID, data.XingYaoPrice)
		/////更新医生的订单信息//////
		go doctor_reserve_people(data.DocID, data.OrderID, false)
		/////删除校验码/////////
		go DelVerifyCode(data.VerifCode)
	}
	return data, nil
}

//系统生成订单,状态：待计算
func SysGenBill(OrderID string) (*ST_Order, error) {
	if OrderID == "" {
		return nil, ErrorLog("系统生成账单失败,OrderID=%s\n", OrderID)
	}
	data := &ST_Order{}
	var isE error = nil
	if err := WriteLock(constant.Hash_Order, OrderID, data); err != nil {
		return data, err
	}
	switch data.Current.OpreatStatus {
	case constant.Status_Order_PendingStatements:
		addOrderStatus(data,
			data.Current.OpreatUserStatus,           //保留原有的用户状态
			data.Current.OpreatAgentStatus,          //保留原来结算状态
			constant.Status_OrderPenddingCollection, //待收款
			constant.Status_OrderPenddingCollection, //待收款
			"后台生成账单",                                //原因
			constant.Platform_Side,                  //平台方
			"", "")
		break
	case constant.Status_Order_Cancle:
		addOrderStatus(data,
			data.Current.OpreatUserStatus,               //保留原状态
			data.Current.OpreatAgentStatus,              //保留原有的状态
			data.Current.OpreatStatus,                   //保留原状态
			constant.Opreat_Order_StatementsAfterCancle, //系统生成账单
			constant.Opreat_Order_StatementsAfterCancle, //原因
			constant.Platform_Side,                      //平台方
			"", "")
		break
	default:
		isE = ErrorLog("系统生成账单失败,当前状态不对,OrderID=%s,Status=%s\n", data.OrderID, data.Current.OpreatStatus)
		break
	}
	if err := WriteBack(constant.Hash_Order, OrderID, data); err != nil {
		return data, err

	}
	return data, isE
}

///系统确认收款，订单完成
func OrderSysConfirmCollection(OrderID, Msg, jobNum, Name string) error {
	if OrderID == "" || jobNum == "" {
		return ErrorLog("系统确认收款失败,参数不全,OrderID=%s,jobNum=%s\n", OrderID, jobNum)
	}
	data := &ST_Order{}
	if err := WriteLock(constant.Hash_Order, OrderID, data); err != nil {
		return err
	}
	defer WriteBack(constant.Hash_Order, OrderID, data)
	if data.Current.OpreatStatus != constant.Status_OrderPenddingCollection {
		return ErrorLog("系统确认收款失败,状态不对,OrderID=%s,Stauts=%s", OrderID, data.Current.OpreatStatus)
	}
	addOrderStatus(data,
		data.Current.OpreatUserStatus, //保留原有的用户状态
		constant.Status_Order_Succeed, //系统确认收款
		constant.Status_Order_Succeed, //系统确认收款
		constant.Status_Order_Succeed, //系统确认收款
		Msg, //原因
		constant.Platform_Side, //平台方
		jobNum,                 //星喜员工工号
		Name)                   //星喜员工姓名
	if data != nil {
		//////////////订单的代理结算//////////////////////
		go pushAgentOrder(data)
	}
	return nil
}

///更新订单的评论信息
func OrderComment(OrderID, UID, Name string) error {
	if OrderID == "" {
		return ErrorLog("updateOrderCommentInfo,参数不全,OrderID=%s,UID=%s\n", OrderID, UID)
	}
	data := &ST_Order{}
	if err := WriteLock(constant.Hash_Order, OrderID, data); err != nil {
		return err
	}
	defer WriteBack(constant.Hash_Order, OrderID, data)

	if data.Current.OpreatUserStatus != constant.Status_Order_PenddingEvaluate {
		return ErrorLog("OrderComment failed,状态不对,OrderID=%s,UID=%s,Stauts=%s",
			OrderID, UID, data.Current.OpreatUserStatus)
	}
	data.CloseDate = CurTime()
	addOrderStatus(data,
		constant.Status_Order_Succeed,         //订单完成
		constant.Status_OrderAlreadyStatement, //代理已结算
		data.Current.OpreatStatus,             //保持系统原有的状态
		constant.Opreat_Order_UserComment,     //用户评价
		constant.Opreat_Order_UserComment,     //用户评价
		constant.User_Side,                    //用户
		UID,                                   //用户id
		Name)                                  //用户姓名

	return nil
}

///添加订单的操作记录
func addOrderStatus(data *ST_Order, userStatus, agentstatus, sysStatus,
	action, reason, side, jobNum, Name string) {
	st := ST_OrderFlow{
		OpreatPart:        side,
		OpreatStatus:      sysStatus,
		OpreatAgentStatus: agentstatus,
		OpreatUserStatus:  userStatus,
		OpreatAction:      action,
		OpreatTime:        CurTime(),
		OpreatReason:      reason,
		OpreatJobNum:      jobNum,
		OpreatName:        Name,
	}
	data.Opreat = append(data.Opreat, st)
	data.Current = st
}

//查询订单
func QueryOrder(OrderID string) (*ST_Order, error) {
	if OrderID == "" {
		return nil, ErrorLog("查询订单失败,订单id为空\n")
	}
	order := &ST_Order{}
	err := ShareLock(constant.Hash_Order, OrderID, order)

	return order, err
}

///查询订单详情
func QueryOrderInfo(session *JsNet.StSession) {
	type st_query struct {
		OrderID string
	}
	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.OrderID == "" {
		ForwardEx(session, "1", nil, "param is empty,OrderID=%s\n", st.OrderID)
		return
	}
	order, err := QueryOrder(st.OrderID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", order)
}

//查询多个产品信息
func QueryMoreOrderInfo(session *JsNet.StSession) {
	type st_query struct {
		Orders []string //产品ids
	}
	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data := QueryMoreOrders(st.Orders)
	Forward(session, "0", data)
}

///查询多个订单信息
func QueryMoreOrders(ids []string) []*ST_Order {
	data := []*ST_Order{}
	if len(ids) == 0 {
		return data
	}
	for _, v := range ids {
		t, e := QueryOrder(v)
		if e != nil {
			continue
		}
		data = append(data, t)
	}
	return data
}

///验证码查看订单
func VerifyQueryOrder(HosID, VerifCode string) (*ST_Order, error) {

	if HosID == "" || VerifCode == "" {
		return nil, ErrorLog("验证失败,参数不全,HosID=%s,VerifCode=%s\n", HosID, VerifCode)
	}

	OrderID, e := CheckVerifyCode(VerifCode, HosID)
	if e != nil {
		return nil, e
	}
	order, err := QueryOrder(OrderID)
	if err != nil {
		return nil, err
	}
	if order.Current.OpreatStatus != constant.Status_Order_PenddingVerify {
		return nil, ErrorLog("该订单还未预约,不能通过校验码查询,OrderID=%s\n", order.OrderID)
	}
	return order, nil
}

///医院医生产品的订单统计信息
func order_statistics(HosID, DocID, ProID, OrderID string, ProPrcie, Deposit int, isAdd bool) {
	////更新医院的订单数量信息
	go hospital_reserve_people(HosID, ProPrcie)
	/////更新医生的订单信息
	go doctor_reserve_people(DocID, OrderID, isAdd)
	////更新产品的的预约销售信息
	go product_bespeak_people(ProID, OrderID, isAdd, Deposit, ProPrcie)
}

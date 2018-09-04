package constant

//系统数据库
const (
	C_XINGYAODB = "xingyaodb"
	C_CONSULTDB = "servicedb"
)

///ID 种子
const (
	KV_Admin          = "Admin"
	KV_MAXID          = "MaxId"
	Com_SEED          = "ComSEED"
	Doc_SEED          = "DocSEED"
	Pro_SEED          = "ProSEED"
	HosID_SEED        = "HosSEED"
	HosAccountID_SEED = "HosAccountSEED"
	Order_SEED        = "OrderSEED"
	Bill_SEED         = "BillSEED"
	Art_SEED          = "ArtSEED"
	USER_SEED         = "UserSEED"
	Agent_SEED        = "AgentSEED"
	SuperArt_SEED     = "SuperArtSEED"
)

//其他基础数据库
const (
	Hash_HosDiary        = "HosDiary"     //整容日记
	Hash_SysNotice       = "SysNotice"    //系统公告
	Hash_Consultation    = "Consultation" //咨询表
	Hash_FrontPageInfoLs = "FrontPageInfoLs"
	Hash_TopArticleInfo  = "TopArticleInfo"
	Hash_ArticleManager  = "ArticleManger"
	Hash_TopArticles     = "TopArticles"
	//Get the PingLun Article ID with the middle content
	Hash_ExampleShow = "ExampleShow"
	Hash_RemarkShow  = "RemarkShow"
)

///////全局相关的
const (
	All                 = "All"
	Hash_Global         = "Global"          //全局的列表
	KEY_GlobalOrderList = "GlobalOrderList" //全局订单列表
	KEY_GlobalBodyList  = "GlobalBodyList"  //存放所有身体部位的列表
	KEY_GlobalCityList  = "GlobalCityList"  //存放所有城市列表
	KEY_GlobalProPrice  = "GlobalProPrice"  //存放所有产品的价格
	KEY_GlobalCityMap   = "CityMap"         //全局城市映射

	KEY_Global_Agent    = "GlobalAgent"      //全局代理人Key
	KEY_City_Agent      = "CityAgent"        //分城市代理
	KEY_Global_Employee = "GlobalEmployee"   //全局店员
	KEY_City_Employee   = "CityEmployee"     //分城市店员
	KEY_Notice_Back     = "GlobalNoticeBack" //后台全局公告
)

///各种缓存表
const (
	Hash_HosConsultation  = "HosConsultation"  //医院咨询表
	Hash_HosConDateList   = "ConDateList"      //医院咨询表日期统计
	Hash_HosConUserList   = "ConUserList"      //医院咨询表，按用户
	Hash_HosVisitDateList = "HosVisitDateList" //医院访问量日期统计
	Hash_LatLot           = "LatLot"
)

/////用户相关的表和key
const (
	Hash_User           = "User"           //前端用户
	Hash_User_Fav       = "UserFavourite"  //用户的喜爱表
	Hash_OpenID_UID     = "OpenID_UID"     //openID映射UID
	Hash_Union_UID      = "UnionID_UID"    //Union映射UID
	KEY_Cell_UID        = "Cell_UID"       //cell-uid
	Hash_WithDraw       = "UserWithDraw"   //用户提现表
	KEY_Global_WithDraw = "GlobalWithDraw" //全局提现记录

	Hash_UserFromAgent = "UserFromAgent" //小B邀请用户统计
	KEY_UserIncrease   = "UserInCrease"  //用户每天的增加量
)

///医院相关
const (
	Hash_Hospital          = "Hospital"          //医院表
	Hash_Hos_Statistics    = "HosStaticsitc"     //医院统计表
	Hash_HosAccount        = "HosAccount"        //医院账号表
	Hash_HosAccountMap     = "HosAccountMap"     //医院账号和医院的映射表
	Hash_HosCellAccountMap = "HosCellAccountMap" //手机号和医院账号的映射表
	Hash_HospitalCache     = "Hospital_Cache"    //医院缓存表

	KEY_HosCityList   = "HosCityList"   //分城市的医院列表
	KEY_HosStatusList = "HosStatusList" //全局分状态的医院
	KEY_HosSimpleInfo = "HosSimpleInfo" //医院的id和名字映射
	//tianfeng use

	KEY_HosFullache = "HosFullCache"   //所有医院的列表
	KEY_HosORMCache = "HosStatusList"  //系统的医院线上线下审核缓存
	KEY_DocORMCache = "DocStatusList"  //系统的医生线上线下审核缓存
	KEY_ProORMCache = "ProdStatusList" //系统的产品线上线下审核缓存

	KEY_ALL_Hospital = "AllHospital" //所有的医院
)

//医生相关
const (
	Hash_Doctor         = "Doctor"        //医生表
	Hash_Doc_Statistics = "DocStaticsitc" //医生统计表
	Hash_DoctorCache    = "Doctor_Cache"  //医院缓存表
	KEY_DocCityList     = "DocCityList"   //分城市的医生列表

)

/////产品的相关
const (
	Hash_HosProduct        = "Product"          //产品表
	Hash_HosProComment     = "ProComment"       //产品评论
	Hash_HosProTransaction = "HosProTransactio" //产品的购买记录
	Hash_HosProVouchers    = "HosProVouchers"   //产品打折券
	Hash_ProductCache      = "Product_Cache"    //产品缓存表
	Hash_Pro_Statistics    = "ProStaticsitc"    //产品统计表

	KEY_ProBodyList  = "ProBodyList"  //分部位的产品列表
	KEY_ProCityList  = "ProCityList"  //分城市的产品列表
	KEY_ProPriceList = "ProPriceList" //分价格的产品列表
	KEY_BodyItemMap  = "BodyItemMap"  //身体部位和item的映射

)

//订单相关
const (
	Hash_Order            = "Order"             //订单表
	Hash_HistoryOrder     = "HistoryOrder"      //历史订单表
	Hash_Order_Verify     = "Order_Verify"      //订单的校验码
	Hash_UserOrderCache   = "UserOrderCache"    //用户订单缓存
	Hash_HosOrderDateList = "HosOrderDateList"  //医院订单日期统计
	Hash_HosOrderCache    = "HosOrderCache"     //医院的订单缓存
	Has_ProOrderCache     = "ProductOrderCache" //跟产品相关的订单缓存
	Hash_HospitalOrder    = "HospitalOrder"     //医院的订单列表

	KEY_OrderStatusList = "OrderStatusList" //订单状态缓存

)

const (
	Hash_HosBill  = "HospitalBill" //医院账单表
	Hash_UserBill = "UserBill"     //用户账单表
	Hash_SysBill  = "SystemBill"   //系统账单

	KEY_HosMonthBill = "HosMonthBill" //医院的各个月份的账单
)

const (
	TOTAL_STATISTICS    = "TOTAL-STATISTICS"
	HOSPITAL_STATISTICS = "HOSPITAL_STATISTICS"
)

////大牌名医相关
const (
	DocUptoExpert       = "DocUptoExpert"     //提升至大牌名医 的状态
	DocDownFromExpert   = "DocDownFromExpert" //从大牌名医移除 的状态
	KV_ExpertDoctor     = "ExpertDoctor"      //大牌名医的key
	KV_HomeExpertDoctor = "HomeExpertDoctor"  //首页大牌名医的key
)

///网络返回字段
const (
	CT_Ret    = "Ret"
	CT_Msg    = "Msg"
	CT_Entity = "Entity"
)

//医院端平台角色
const (
	Author_admin    = "Admin"    //管理员
	Author_doctor   = "Doctor"   //医生
	Author_finance  = "Finance"  //财务
	Author_operator = "Operator" //操作员
)

//状态缓存的状态
const (
	OperatingStatus_new              = "Status_New"            //"新建待审核"
	OperatingStatus_online           = "Status_Online"         //"线上运营"
	OperatingStatus_Offline_self     = "Status_OfflineSelf"    //"主动下线"
	OperatingStatus_Offline_onforce  = "Status_OfflineOnforce" //"强制下线"
	OperatingStatus_modify           = "Status_Modify"         //"修改待审核"
	OperatingStatus_Reviewer_NotPass = "Status_NoPass"         //"审核不通过"
	OperatingStatus_Del              = "Status_Del"            //"删除"

)

///公告在线状态
const (
	NoticeStatus_online  = "Online"  //在线
	NoticeStatus_offline = "Offline" //下线
)

//产品类型
const (
	Type_Product_common  = "Type_Product_common"  //标准产品
	Type_Product_special = "Type_Product_special" //特价产品
	Type_Product_custom  = "Type_Product_custom"  //联合定制
	Type_Product_rush    = "Type_Product_rush"    //抢购产品
)

////用来记录某一个动作操作的某一方
const (
	Platform_Side = "Platform_Side" //"平台方"
	Hospital_Side = "Hospital_Side" //"医院方"
	User_Side     = "User_Side"     //"用户"
)

//订单动作
const (
	Opreat_Order_Submit                      = "用户提交订单"
	Opreat_Order_Invalid                     = "支付超时,已取消订单"
	Opreat_Order_UserPay                     = "用户支付订单"
	Opreat_Order_UserComment                 = "用户评论"
	Opreat_Order_SysAppoint                  = "系统预约"
	Opreat_Order_HosVerify                   = "医院校验"
	Opreat_Order_HosReconcile                = "医院核实账单"
	Opreat_Order_SysAutoGenBill              = "系统自动生成账单"
	Opreat_Order_StatementsAfterCancle       = "退款后重新生成账单"
	Opreat_Order_UserCancleBeforeAppointment = "预约前用户退款"
	Opreat_Order_UserCancleBeforeVerify      = "到院校验前用户退款"
	Opreat_Order_UserCancleAfterVerify       = "到院校验后用户退款"
)

///订单状态
const (
	///
	Status_OrderAlreadyStatement = "Order_AlreadyStatement" ///已结算(小B)

	Status_Order_PenddingPay      = "Order_PenddingPay"      //"待支付"//用户、代理、系统
	Status_Order_UserPenddingUse  = "Order_UserPenddingUse"  //"待使用"//用户
	Status_Order_PenddingEvaluate = "Order_PenddingEvaluate" //"待评价"//用户
	Status_Order_Succeed          = "Order_Succeed"          //"已完成"//用户
	Status_Order_Invalid          = "Order_Invalid"          //"已失效"//用户
	Status_Order_Cancle           = "Order_Cancle"           //"已取消"//用户、医院、代理、系统
	Status_Order_AlreadyRefund    = "Order_AlreadyRefund"    //"已退款"//用户
	Status_Order_NeedDiary        = "Order_NeedDiary"        //"需要日记"//用户、医院、系统、代理

	Status_Order_PenddingAppointment     = "Order_PenddingAppointment"     //"待预约"//医院、代理、系统
	Status_Order_PenddingVerify          = "Order_PenddingVerify"          //"待校验"//医院、代理、系统
	Status_Order_PendingConfirm          = "Order_PendingConfirm"          //"待确认"//医院、系统、代理
	Status_Order_PendingStatements       = "Order_PendingStatements"       //"待结算"//医院、系统、代理
	Status_Order_PenddingReconcile       = "Order_PenddingReconcile"       //"待对账"//医院、系统、代理
	Status_OrderPenddingCollection       = "OrderPenddingCollection"       //"待收款"//医院、系统、代理
	Status_OrderSysConfirmCollection     = "OrderSysConfirmCollection"     //"已收款"//医院、系统、代理
	Status_Order_CancleBeforeAppointment = "Order_CancleBeforeAppointment" //用户预约前取消
	Status_Order_CancleBeforeVerfy       = "Order_CancleBeforeVerfy"       //用户校验前取消
	Status_Order_CancleAfterVerfy        = "Order_CancleAfterVerfy"        //用户校验后取消

	Status_Order_AgentEntity = "Order_AgentEntity"
	Status_Order_AgentWeChat = "Order_AgentWeChat"
)

/////支付方式
const (
	PayMod_DepositPayment = "DepositPayment" //"订金支付"
	PayMod_FullPayment    = "FullPayment"    //"全款支付"
	PayWay_Wechat         = "微信支付"
	PayWay_Ali            = "支付宝支付"
	PayWay_Bank           = "银联支付"
)

////////////
const (
	ItemAccountPerPage_Product  = 10
	ItemAccountPerPage_Doctor   = 10
	ItemAccountPerPage_Hospital = 10
	ItemAccountPerPage_Order    = 10
	ItemAccountPerPage_Bill     = 10
	ItemAccountPerPage_Business = 10
)

//代理人状态小B
const (
	Agent_Apply         = "Agent_Apply"         //申请成为代理
	Agent_ReApply       = "Agent_ReApply"       //重新申请成为代理
	Agent_PassReviewe   = "Agent_PassReviewe"   //资料审核通过
	Agent_NoPassReviewe = "Agent_NoPassReviewe" //资料审核不通过
	Agent_Online        = "Agent_Online"        //代理上线
	Agent_Offline_force = "Agent_Offline_force" //平台强制下线
	Agent_Offline_self  = "Agent_Offline_self"  //代理自己下线
)

//小B的代理级别
const (
	Agent_Level_Diamonds_A = "Level_Diamonds_A" //钻石A级
	// Agent_Level_Diamonds_B = "Level_Diamonds_B"    //钻石B级
	// Agent_Level_Gold_A     = "Level_Gold_A"        //黄金A级
	// Agent_Level_Gold_B     = "Level_Gold_B"        //黄金B级
	// Agent_Level_Silver     = "Level_Silver"        //白银级
	Agent_Level_TryUse = "Level_TryUse" //试用级别
	// Agent_Level_Employee   = "AgentLevel_Employee" //店员
)

//代理人操作
const (
	Agent_Submit   = "提交代理费用订单"
	Agent_Pay      = "代理费用支付"   //代理费用支付
	Agent_RePay    = "代理费用重新支付" //代理费用重新支付
	Agent_WithDraw = "申请代理费用退款" //申请代理费用退款
)

////代理商存在形式
const (
	Agent_Type_Physical = "Agent_Type_Physical" //实体店
	Agent_Type_Wechat   = "Agent_Type_Wechat"   //微商
)

const (
	ArticleNumPerPage_All = 10
	ArticleNumPerPage_His = 10
)

////文章相关的key
const (
	Hash_Article         = "Article"
	Hash_SuperArticle    = "SuperArticle"
	Key_ArticleTitleTips = "ArticleTitleTips"
	Key_ArticleHis       = "ArticleHis"
	Key_ArticleCurrent   = "ArticleCurrent"
	Key_ArticleAll       = "ArticleAll"
	Key_ArticleManager   = "ArticleManager"
	Key_ArticleSearch    = "ArticleSearch"
	Key_FrontPage        = "FrontPage"
)

////文章相关
const (
	ArticleType_AllArticle         = "AllArticle"
	ArticleType_MainPageHisArticle = "MainPageHisArticle"
	ArticleType_SubPageHisArticle  = "SubPageHisArticle"

	Article_toppageshow    = "toppageshow"
	Article_toppagesearch  = "toppagesearch"
	Article_topinfoxingyao = "topinfoxingyao"
	Article_smartdoctor    = "smartdoctor"
	Article_fourpartshow   = "fourpartshow"
	Article_editarticle    = "editarticle"
)

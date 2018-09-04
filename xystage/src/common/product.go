package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	. "cache/cacheIO"
	"constant"
	"encoding/json"
	"ider"
	. "util"
)

type ST_ProBaseInfo struct {
	ProName           string   //产品名称
	HosID             string   //医院id
	Hospital          string   //所属医院
	HosCity           string   //所在医院的城市
	Pics              []string //产品的宣传图片
	FirstItem         string   //一级项目
	SecondItems       []string //二级项目
	ThirdItems        []string //三级项目
	Doctors           string   //关联的医生
	DoctorName        string   //关联的医生的名字
	ProType           string   //产品类型:常规产品、特价产品、联合定制
	ServiceBodyPart   string   //服务部位
	SaleStartDate     string   //启售时间(特价产品)
	SaleEndDate       string   //售卖截止时间(特价产品)
	Inventoryquantity int      //库存数量
	ProCase           string   //产品案例
}

//价格信息
type ST_PriceInfo struct {
	MarketPrice    int      //市场价
	XingYaoPrice   int      //星喜价
	FullMoney      int      //全额
	ProDeposit     int      //订金
	IsFullPay      bool     //是否是全额付款：定金支付/全款付款
	IsCoupon       bool     //是否可以使用优惠券
	IncludeCost    []string //包含费用列表
	OutIncludeCost []string //不包含费用列表
	ProRatio       int      //产品结算比例
}

// 预约须知
type ST_AppointmentNotice struct {
	ForwardAppoint string   //提前预约时间
	ValidityTerm   string   //有效期
	UseTime        string   //使用时间
	ServicePromise []string //服务承诺:(无隐性消费、未使用随时退款）
}

//产品详情
type ST_ProductDetail struct {
	ST_AppointmentNotice          // 预约须知
	DetailPics           []string // 图片列表
	ProIntroduce         string   // 产品详细介绍
	ProDes               string   // 产品的一句话描述
	ProCategory          string   // 产品类别:玻尿酸、肉毒素、隆胸假体、超声刀设备、热玛吉、隆鼻假体
	ProBrand             string   // 产品品牌
	ProUnit              string   // 产品计量单位：毫升、单位、CC
	ProSpecifications    string   // 产品规格：毫升（0.5ml、0.75ml、1ml、1.5ml）,单位（50单位、100单位）,CC(200CC、250CC)
	OtherIngredients     string   // 其他成分：是（PRP、胶原蛋蛋白、vc、玻尿酸、肉毒素、其他）
	PreoperativeNotice   string   // 术前须知
	ApplicablePart       string   // 适用部位（单部位、多部位）

	TreatmentDuration string // 治疗时长(可选择填写)
	PainLevel         int    // 疼痛级别(可选择填写)
	RiskLevel         int    // 风险级别(可选择填写)
	TreatmentTimes    int    // 治疗次数(可选择填写)

	MaterialName    string // 材料名称(可选择填写)
	MaterialBrand   string // 材料品牌(可选择填写)
	MaterialOrigin  string // 材料源点(可选择填写)
	MaterialShape   string // 材料形状(可选择填写)
	MaterialTexture string // 材料质地(可选择填写)

	IsIncision      string // 是否有切口(可选择填写)
	IncisionPlace   string // 切口位置(可选择填写)
	InjectionPlace  string // 注射位置(可选择填写)
	InjectionDosage string // 注射用量(可选择填写)

	DeviceName   string // 设备名称(可选择填写)
	DeviceBrand  string // 设备品牌(可选择填写)
	DeviceOrigin string // 设备产地(可选择填写)

	AnesthesiaMethod string // 麻醉方式(可选择填写)
	ScarDescription  string // 疤痕描述(可选择填写)
	RecoveryProcess  string // 恢复过程(可选择填写)
	EffectDuration   string // 效果持续时间(可选择填写)
	ExperienceItems  string // 体检项目(可选择填写)

	ProjectAdvantages    string //项目优点
	ProjectDisadvantages string //项目缺点
	SideEffect           string //副作用及风险
	Convalescence        string //恢复期
	RecoveryReminder     string //恢复提醒
	NursingMethod        string //护理方法
	Considerations       string //注意事项
	CommonProblem        string //常见问题

}

type ST_ProStatics struct {
	AttentionNums    int // 关注数量
	Bespeakquota     int // 预约额度
	AppointmentNums  int // 预约数量
	CompositeScore   int // 综合评分
	StarGrade        int // 星级评分
	ServiceScore     int // 服务评分
	EnvironmentScore int // 环境评分
	EffectScore      int // 效果评分
	Peoples          int // 评论的人数
	SaveNum          int // 点赞人数
	Salesvolumes     int // 历史销售数量
	Salesquota       int // 历史销售额度
	VisitNum         int // 访问量
	ConsultNum          int      //咨询量
}

//产品结构
type ST_Product struct {
	ProID            string // 产品id
	ST_ProBaseInfo          // 基本信息
	ST_ProductDetail        // 产品详情
	ST_PriceInfo            // 价格信息
	ST_OpreatStatus         // 产品的操作信息
	ST_ProStatics           // 产品的统计信息
}

///创建一个产品
func ApplyNewProduct(session *JsNet.StSession) {
	st := &ST_Product{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "ApplyNewProduct,err:"+err.Error())
		return
	}
	if err := proCheckData(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	st.CreatDate = CurTime()
	if id, err := ider.GenProID(); err == nil {
		st.ProID = id
	} else {
		ForwardEx(session, "1", nil, "产品ID生成失败,HosID=%s,ProID=%s\n", st.HosID, st.ProID)
		return
	}
	hos, err := GetHospitalInfo(st.HosID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	st.CompositeScore = 500
	st.ServiceScore = 500
	st.EnvironmentScore = 500
	st.EffectScore = 500
	st.StarGrade = 5
	st.HosCity = hos.HosCity
	st.Hospital = hos.HosName
	appendProStatus(st,
		constant.Hospital_Side,
		"新增产品",
		constant.OperatingStatus_new,
		"无",
		"", "", "",
		hos.AdminAccount,
		hos.AdminName,
		hos.AdminCell)

	if err := DirectWrite(constant.Hash_HosProduct, st.ProID, st); err != nil {
		Error(err.Error())
		ForwardEx(session, "1", nil, "产品创建失败,HosID=%s,ProID=%s\n", st.HosID, st.ProID)
		return
	}
	go AddHosProduct(st.HosID, st.ProID)
	go AddProjectToDoc(st.Doctors, st.ProID)
	go ChangeHosProNum(st.HosID)
	go GlobalAddNewProduct(st.ProID, st.HosID, st.Doctors, st.FirstItem, st.HosCity, st.ProName, st.SecondItems, st.XingYaoPrice)
	Forward(session, "0", st)

}

//查询产品信息
func QueryProduct(session *JsNet.StSession) {
	type st_query struct {
		ProID string //产品id
	}
	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data, err := GetProductInfo(st.ProID)
	if err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	Forward(session, "0", data)
}

//查询多个产品信息
func QueryMoreProductInfo(session *JsNet.StSession) {
	type st_query struct {
		Pros []string //产品ids
	}
	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data := QueryMoreProducts(st.Pros)
	Forward(session, "0", data)
}

///查询多个订单信息
func QueryMoreProducts(ids []string) []*ST_Product {
	data := []*ST_Product{}
	if len(ids) == 0 {
		return data
	}
	for _, v := range ids {
		t, e := GetProductInfo(v)
		if e != nil {
			continue
		}
		data = append(data, t)
	}
	return data
}

//查产品信息
func GetProductInfo(ProID string) (*ST_Product, error) {
	if ProID == "" {
		return nil, ErrorLog("查询产品id失败,id为空!\n")
	}
	data := &ST_Product{}
	if err := ShareLock(constant.Hash_HosProduct, ProID, data); err != nil {
		return data, ErrorLog("查询产品id失败ShareLock(),ProID=%s!\n", ProID)
	}
	return data, nil
}

//审核产品
func RevieweNewProduct(session *JsNet.StSession) {
	type st_reviewe struct {
		HosID        string //医院id
		ProID        string //产品id
		ProDeposit   int    //订金
		OpreatReason string //操作原因
		OpreatJobNum string //操作人员工号
		OpreatName   string //操作人员姓名
		OpreatCell   string //操作人手机
		IsPass       bool   //审核是否通过:true:通过，false:不通过
		ProRatio     int    //产品结算比例
	}
	st := &st_reviewe{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" || st.ProID == "" || st.OpreatJobNum == "" || st.ProDeposit < 0 || st.ProRatio < 0 {
		ForwardEx(session, "1", nil,
			"产品审核失败,参数不全,HosID=%s,ProID=%s,OpreatJobNum=%s,HosCity=%s,ProDeposit=%d,ProRatio=%d\n",
			st.HosID, st.ProID, st.OpreatJobNum, st.ProDeposit, st.ProRatio)
		return
	}
	data := &ST_Product{}
	if err := WriteLock(constant.Hash_HosProduct, st.ProID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	defer WriteBack(constant.Hash_HosProduct, st.ProID, data)

	if data.Current.OpreatStatus == constant.OperatingStatus_online {
		Forward(session, "0", data)
		return
	}
	if data.Current.OpreatStatus != constant.OperatingStatus_new {
		ForwardEx(session, "1", data, "初次审核失败,当前产品状态为%s,不在待审核状态\n", data.Current.OpreatStatus)
		return
	}

	data.ProDeposit = st.ProDeposit
	//////////添加操作状态////////
	data.FirstOnlineDate = CurTime()
	data.LastOnlineDate = CurTime()
	data.LastOnlineName = st.OpreatName
	data.ProRatio = st.ProRatio //产品的结算比例
	Len := len(data.OpreatInfo)
	status := constant.OperatingStatus_Reviewer_NotPass
	action := "新增产品审核不通过"
	if st.IsPass {
		status = constant.OperatingStatus_online
		action = "新增产品审核通过"
	}
	appendProStatus(data,
		constant.Platform_Side,
		action,
		status,
		st.OpreatReason,
		st.OpreatJobNum,
		st.OpreatName,
		st.OpreatCell,
		data.OpreatInfo[Len-1].OpreatJobNum,
		data.OpreatInfo[Len-1].OpreatName,
		data.OpreatInfo[Len-1].OpreatCell)
	go GlobalUpdateProStatus(st.HosID)
	Forward(session, "0", data)
}

//审核修改后的产品
func RevieweModifyProduct(session *JsNet.StSession) {
	type st_reviewe struct {
		HosID        string //医院id
		ProID        string //产品id
		OpreatReason string //操作原因
		OpreatJobNum string //操作人员工号
		OpreatName   string //操作人员姓名
		OpreatCell   string //操作人手机
		IsPass       bool   //审核是否通过:true:通过，false:不通过
		ProDeposit   int
	}
	st := &st_reviewe{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if st.HosID == "" || st.ProID == "" || st.OpreatJobNum == "" {
		ForwardEx(session, "1", nil,
			"产品修改审核失败,参数不全,HosID=%s,ProID=%s,OpreatJobNum=%s,\n", st.HosID, st.ProID, st.OpreatJobNum)
		return
	}

	data := &ST_Product{}
	if err := WriteLock(constant.Hash_HosProduct, st.ProID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	defer WriteBack(constant.Hash_HosProduct, st.ProID, data)
	if data.Current.OpreatStatus == constant.OperatingStatus_online {
		Forward(session, "0", data)
		return
	}
	if data.Current.OpreatStatus != constant.OperatingStatus_modify {
		ForwardEx(session, "1", data, "修改审核失败,当前产品状态为%s,不在修改状态\n", data.Current.OpreatStatus)
		return
	}
	if st.ProDeposit > 0 {
		data.ProDeposit = st.ProDeposit
	}

	status := constant.OperatingStatus_Reviewer_NotPass
	action := "修改产品审核不通过"
	data.LastOnlineDate = CurTime()
	data.LastOnlineName = st.OpreatName
	Len := len(data.OpreatInfo)
	if st.IsPass {
		status = constant.OperatingStatus_online
		action = "修改产品审核通过"
	}
	appendProStatus(data,
		constant.Platform_Side,
		action,
		status,
		st.OpreatReason,
		st.OpreatJobNum,
		st.OpreatName,
		st.OpreatCell,
		data.OpreatInfo[Len-1].OpreatJobNum,
		data.OpreatInfo[Len-1].OpreatName,
		data.OpreatInfo[Len-1].OpreatCell)
	go GlobalUpdateProStatus(st.HosID)
	Forward(session, "0", data)
}

//修改产品的结算比例
func ModifyProductRatio(session *JsNet.StSession) {
	type st_info struct {
		ProID        string //产品id
		ProRatio     int    //产品结算比例
		OpreatJobNum string //操作人员工号
		OpreatName   string //操作人员姓名
		OpreatCell   string //操作人手机
	}

	st := &st_info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ProID == "" || st.ProRatio < 0 {
		ForwardEx(session, "1", nil, "修改产品信息失败,ProID=%s,ProRatio=%d\n", st.ProID, st.ProRatio)
		return
	}
	data := &ST_Product{}
	if err := WriteLock(constant.Hash_HosProduct, st.ProID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.ProRatio = st.ProRatio //修改产品结算比例
	data.LastOnlineDate = CurTime()
	data.LastOnlineName = st.OpreatName
	Len := len(data.OpreatInfo)
	appendProStatus(data,
		constant.Platform_Side,
		"修改产品结算比例",
		data.Current.OpreatStatus,
		"修改产品审核通过",
		st.OpreatJobNum,
		st.OpreatName,
		st.OpreatCell,
		data.OpreatInfo[Len-1].OpreatJobNum,
		data.OpreatInfo[Len-1].OpreatName,
		data.OpreatInfo[Len-1].OpreatCell)

	if err := WriteBack(constant.Hash_HosProduct, st.ProID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

//修改产品的信息
func ModifyProductInfo(session *JsNet.StSession) {
	type st_modify struct {
		ProID            string //产品id
		ST_ProBaseInfo          // 基本信息
		ST_ProductDetail        // 产品详情
		ST_PriceInfo            // 价格信息
		OpreatReason     string //操作原因
		OpreatJobNum     string //操作人员工号
		OpreatName       string //操作人员姓名
		OpreatCell       string //操作人手机
	}
	st := &st_modify{}
	data := &ST_Product{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "ModifyProductInfo1,err:"+err.Error())
		return
	}
	if err := session.GetPara(data); err != nil {
		ForwardEx(session, "1", nil, "ModifyProductInfo2,err:"+err.Error())
		return
	}
	if st.ProID == "" {
		ForwardEx(session, "1", nil, "修改产品信息失败,ProID为空\n")
		return
	}
	if err := proCheckData(data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	oldPrice := 0
	if err := WriteLock(constant.Hash_HosProduct, st.ProID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_HosProduct, st.ProID, data)
	oldPrice = data.XingYaoPrice
	b, err := json.Marshal(st)
	if err != nil {
		ForwardEx(session, "1", nil, "json.Marshal() error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(b, data); err != nil {
		ForwardEx(session, "1", nil, "json.Unmarshal() error: %s", err.Error())
		return
	}
	/////添加操作记录
	data.LastModifyDate = CurTime()
	appendProStatus(data,
		constant.Hospital_Side,
		"修改产品信息",
		constant.OperatingStatus_modify,
		st.OpreatReason,
		"", "", "",
		st.OpreatJobNum,
		st.OpreatName,
		st.OpreatCell)

	if st.XingYaoPrice != oldPrice {
		Warn("发现产品%s的将价格从%s更改为%d...\n", st.ProID, oldPrice, st.XingYaoPrice)
		go GlobalUpdateProPrice(st.HosCity, st.ProID, st.XingYaoPrice, oldPrice)
	}
	if data != nil {
		go GlobalUpdateProStatus(data.HosID)
	}
	Forward(session, "0", data)
}

//修改医生
func ModifyProCity(ProID, HosCity string) error {
	if ProID == "" || HosCity == "" {
		return ErrorLog("ModifyProCity failed,ProID=%s,HosCity=%s", ProID, HosCity)
	}
	data := &ST_Product{}
	if err := WriteLock(constant.Hash_HosProduct, ProID, data); err != nil {
		return err
	}
	data.HosCity = HosCity
	return WriteBack(constant.Hash_HosProduct, ProID, data)
}

//医院下线产品
func OfflineProductSelf(session *JsNet.StSession) {
	type st_offline struct {
		HosID        string //医院id
		ProID        string //产品id
		OpreatReason string //操作原因
		OpreatJobNum string //操作人员工号
		OpreatName   string //操作人员姓名
		OpreatCell   string //操作人联系方式
	}

	st := &st_offline{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ProID == "" || st.HosID == "" {
		ForwardEx(session, "1", nil, "医院主动下线产品失败,ProID=%s,HosID=%s\n",
			st.ProID, st.HosID)
		return
	}
	data := &ST_Product{}
	if err := WriteLock(constant.Hash_HosProduct, st.ProID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	defer WriteBack(constant.Hash_HosProduct, st.ProID, data)
	if data.Current.OpreatStatus == constant.OperatingStatus_Offline_self ||
		data.Current.OpreatStatus == constant.OperatingStatus_Offline_onforce {
		return
	}
	/////添加操作记录
	appendProStatus(data,
		constant.Hospital_Side,
		"医院主动下架产品",
		constant.OperatingStatus_Offline_self,
		st.OpreatReason,
		"", "", "",
		st.OpreatJobNum,
		st.OpreatName,
		st.OpreatCell)
	go GlobalUpdateProStatus(st.HosID)
	Forward(session, "0", data)
}

///内部强制下线
func OfflineProOnForce(HosID, ProID string) error {
	if HosID == "" || ProID == "" {
		return ErrorLog("OfflineProOnForce failed,HosID=%s,ProID=%s\n", HosID, ProID)
	}
	data := &ST_Product{}
	if err := WriteLock(constant.Hash_HosProduct, ProID, data); err != nil {
		return err
	}
	defer WriteBack(constant.Hash_HosProduct, ProID, data)
	if data.Current.OpreatStatus == constant.OperatingStatus_Offline_self ||
		data.Current.OpreatStatus == constant.OperatingStatus_Offline_onforce {
		return nil
	}
	appendProStatus(data,
		constant.Platform_Side,
		"医院下架,相关联的产品同步下架",
		constant.OperatingStatus_Offline_onforce,
		"", "", "",
		"", "", "", "")

	return nil
}

//平台强制下线产品
func OfflineProductOnforce(session *JsNet.StSession) {
	type st_offline struct {
		ProID        string //产品id
		HosID        string //医院id
		OpreatReason string //操作原因
		OpreatJobNum string //操作人员工号
		OpreatName   string //操作人员姓名
		OpreatCell   string //操作人手机
	}
	st := &st_offline{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ProID == "" || st.HosID == "" || st.OpreatJobNum == "" {
		ForwardEx(session, "1", nil, "产品强制下线失败,ProID=%s,HosID=%s,OpreatJobNum=%s\n",
			st.ProID, st.HosID, st.OpreatJobNum)
		return
	}
	data := &ST_Product{}
	if err := WriteLock(constant.Hash_HosProduct, st.ProID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_HosProduct, st.ProID, data)
	if data.Current.OpreatStatus == constant.OperatingStatus_Offline_self ||
		data.Current.OpreatStatus == constant.OperatingStatus_Offline_onforce {
		Forward(session, "0", data)
		return
	}
	/////添加操作记录
	appendProStatus(data,
		constant.Hospital_Side,
		"平台强制下架产品",
		constant.OperatingStatus_Offline_onforce,
		st.OpreatReason,
		st.OpreatJobNum,
		st.OpreatName,
		st.OpreatCell,
		"", "", "")
	go GlobalUpdateProStatus(st.HosID)
	Forward(session, "0", data)

}

//平台强制下线产品
func DelProduct(session *JsNet.StSession) {
	type info struct {
		ProID string //产品id
		HosID string //医院id
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ProID == "" || st.HosID == "" {
		ForwardEx(session, "1", nil, "产品删除失败,ProID=%s,HosID=%s\n", st.ProID, st.HosID)
		return
	}
	data := &ST_Product{}
	if err := WriteLock(constant.Hash_HosProduct, st.ProID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	/////添加操作记录
	appendProStatus(data,
		constant.Hospital_Side,
		"平台删除产品",
		constant.OperatingStatus_Del,
		"", "", "", "",
		"", "", "")

	if err := WriteBack(constant.Hash_HosProduct, st.ProID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	go GlobalUpdateProStatus(st.HosID)
	Forward(session, "0", data)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////

///获取医院下的所有产品
func GetHosProductList(HosID string) []string {
	proids := []string{}
	if HosID == "" {
		return proids
	}
	if err := ShareLock(constant.Hash_ProductCache, HosID, &proids); err != nil {
		Error(err.Error())
		return []string{}
	}
	return proids
}

//增加一个产品到医院
func AddHosProduct(HosID, ProID string) error {
	if HosID == "" || ProID == "" {
		return ErrorLog("AddHosProduct param failed,HosID =%S,DocID=%s\n", HosID, ProID)
	}
	proids := []string{}
	err := WriteLock(constant.Hash_ProductCache, HosID, &proids)
	if err != nil {
		proids = append(proids, ProID)
		return DirectWrite(constant.Hash_ProductCache, HosID, &proids)
	}
	exist := false
	for _, v := range proids {
		if v == ProID {
			exist = true
			break
		}
	}
	if !exist {
		proids = append(proids, ProID)
	}

	return WriteBack(constant.Hash_ProductCache, HosID, &proids)
}

///关注或者取消关注产品
func AttentionProduct(ProID string, IsAttention bool) error {
	if ProID == "" {
		return ErrorLog("AttentionProduct failed,ProID=%s,IsAttention=%v\n", ProID, IsAttention)
	}
	data := &ST_Product{}
	if err := WriteLock(constant.Hash_HosProduct, ProID, data); err != nil {
		return err
	}
	if IsAttention {
		data.AttentionNums++
	} else {
		data.AttentionNums--
	}
	if data.AttentionNums <= 0 {
		data.AttentionNums = 0
	}

	return WriteBack(constant.Hash_HosProduct, ProID, data)
}

////////////综合打分
func product_score(ProID string, Service int, Environmental int, Effect int) error {
	data := &ST_Product{}
	if ProID == "" || Service < 0 || Environmental < 0 || Effect < 0 {
		return ErrorLog("服务产品评分失败,ProID=%s,intService=%d, Environmental=%d, Effect int=%d\n", ProID,
			Service, Environmental, Effect)
	}

	if err := WriteLock(constant.Hash_HosProduct, ProID, data); err != nil {
		return err
	}
	ser := data.ServiceScore*data.Peoples + Service
	en := data.EnvironmentScore*data.Peoples + Environmental
	ef := data.EffectScore*data.Peoples + Effect
	pe := data.Peoples + 1

	data.EnvironmentScore = en / pe
	data.ServiceScore = ser / pe
	data.EffectScore = ef / pe
	data.Peoples = pe
	data.CompositeScore = (data.EnvironmentScore*4 + data.ServiceScore*4 + data.EffectScore*2) / 10
	data.StarGrade = data.CompositeScore / 100
	if data.StarGrade == 0 {
		data.StarGrade = 1
	}

	//////更新医生的环境分和服务分
	doctor_evaluate_people(data.Doctors, data.EnvironmentScore, data.ServiceScore, data.EnvironmentScore)
	hospital_evaluate_people(data.Doctors, data.EnvironmentScore, data.ServiceScore, data.EnvironmentScore)
	if data != nil {
		go GlobalUpdateProComment(ProID, data.Peoples)
	}
	return WriteBack(constant.Hash_HosProduct, ProID, data)
}

//产品点赞
func productSave(ProID string) (*ST_Product, error) {
	data := &ST_Product{}
	if ProID != "" {
		if err := WriteLock(constant.Hash_HosProduct, ProID, data); err != nil {
			return data, err
		}
		data.SaveNum++
		return data, WriteBack(constant.Hash_HosProduct, ProID, data)
	}
	return nil, ErrorLog("产品点赞失败,ProID=%s\n", ProID)
}

//更改产品的访问量
func ChangeProVisitNum(session *JsNet.StSession) {
	type st_offline struct {
		ProID string //产品id
	}
	st := &st_offline{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ProID == "" {
		ForwardEx(session, "1", nil, "增加产品访问量失败，ProID=%s\n", st.ProID)
		return
	}
	data := &ST_Product{}
	if err := WriteLock(constant.Hash_HosProduct, st.ProID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.VisitNum++
	if err := WriteBack(constant.Hash_HosProduct, st.ProID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}



//更改产品的访问量
func ChangeProConsultNum(session *JsNet.StSession) {
	type st_offline struct {
		ProID string //产品id
	}
	st := &st_offline{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ProID == "" {
		ForwardEx(session, "1", nil, "增加产品咨询量失败，ProID=%s\n", st.ProID)
		return
	}
	data := &ST_Product{}
	if err := WriteLock(constant.Hash_HosProduct, st.ProID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.ConsultNum++
	if err := WriteBack(constant.Hash_HosProduct, st.ProID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}





//预约人数
func product_bespeak_people(ProID, OrderID string, Bespeak bool, Bespeakquota, Salesquota int) (*ST_Product, error) {
	data := &ST_Product{}
	if ProID == "" && Bespeakquota < 0 || Salesquota < 0 {
		return nil, ErrorLog("预约人数失败,ProID=%s,Bespeakquota%d,Salesquota=%s\n", ProID, Bespeakquota, Salesquota)
	}

	if err := WriteLock(constant.Hash_HosProduct, ProID, data); err != nil {
		return data, err
	}
	if Bespeak {
		data.AppointmentNums++
		data.Bespeakquota += Bespeakquota
		data.Salesquota += Salesquota
		data.Salesvolumes++
	} else {
		data.AppointmentNums--
		data.Bespeakquota -= Bespeakquota
		if data.AppointmentNums < 0 {
			data.AppointmentNums = 0
		}
		if data.Bespeakquota < 0 {
			data.Bespeakquota = 0
		}
	}
	if data != nil {
		go GlobalUpdateProSale(ProID, data.Salesvolumes)
	}
	return data, WriteBack(constant.Hash_HosProduct, ProID, data)
}

//输入参数检查
func proCheckData(st *ST_Product) error {
	if st.ProName == "" {
		return ErrorLog("产品名称不能为空!\n")
	}
	if st.Doctors == "" {
		return ErrorLog("产品关联的医生不能为空!\n")
	}
	if st.DoctorName == "" {
		return ErrorLog("产品关联的医生名字不能为空!\n")
	}

	if st.ProType != constant.Type_Product_custom &&
		st.ProType != constant.Type_Product_special &&
		st.ProType != constant.Type_Product_common &&
		st.ProType != constant.Type_Product_rush {
		return ErrorLog("产品的类型不符合要求,ProType=%s\n", st.ProType)
	}

	if st.XingYaoPrice <= 0 {
		return ErrorLog("产品星喜价格不能为空!\n")
	}
	if len(st.SecondItems) == 0 || len(st.ThirdItems) == 0 {
		return ErrorLog("产品包含的项目不能为空,FirstItem=%s,SecondItems=%v,ThirdItems=%v\n",
			st.FirstItem, st.SecondItems, st.ThirdItems)
	}
	//fmt.Printf("订金：%d, 星喜价：%d\n", st.ProDeposit, st.XingYaoPrice)
	if st.ProDeposit > st.XingYaoPrice {
		return ErrorLog("产品订金不能大于星喜价!\n")
	}
	return nil
}

//往医院里增加一条操作记录
func appendProStatus(doc *ST_Product, OpreatPart,
	OpreatAction, OpreatStatus, OpreatReason,
	OpreatJobNum, OpreatName, OpreatCell,
	ApplyJobNum, ApplyName, ApplyCell string) {
	info := ST_Opreat{
		OpreatPart:   OpreatPart,
		OpreatAction: OpreatAction,
		OpreatStatus: OpreatStatus,
		OpreatReason: OpreatReason,
		OpreatJobNum: OpreatJobNum,
		OpreatName:   OpreatName,
		OpreatTime:   CurTime(),
		ApplyJobNum:  ApplyJobNum,
		ApplyName:    ApplyName,
		ApplyCell:    ApplyCell,
	}
	doc.Current = info
	doc.OpreatInfo = append(doc.OpreatInfo, info)
}

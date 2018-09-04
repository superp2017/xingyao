package common

import (
	"JsGo/JsMobile"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	. "cache/cacheIO"
	"constant"
	"encoding/json"
	"ider"
	. "util"
)

type ST_PicMap struct {
	Pic string //图片路径
	Des string //图片描述
}

type ST_HostBaseInfo struct {
	HosName        string      //医院名称
	HosCity        string      //医院所在城市
	AdminAccount   string      //管理员账号
	AdminName      string      //管理员姓名
	AdminCell      string      //管理员号码
	HosAddr        string      //详细地址
	HosLogo        string      //医院头像
	HosPic         []ST_PicMap //医院的宣传图片
	HosIntroduce   string      //医院详细介绍
	HosServe       string      //医院服务特色
	LRName         string      //法人代表名字LR:LegalRepresentative 法人代表
	LRNumber       string      //法人代表号码
	ContactsName   string      //医院联系人名字
	ContactsNumber string      //医院联系人号码
	HosCapitalType string      //医院资本类型,公立医院/民营医院
	ServiceType    string      //服务类型:医疗美容、综合、空腔专科、眼科专科、皮肤科专科
	HosType        string      //医院类型：医院、门诊部、诊所、其他
	ServiceLevel   string      //业务等级：一级、二级、三级、四级、其他
	ServiceRange   []string    //擅长的项目列表
	HotPrograms    []string    //热门项目
	RegisterDate   string      //医院注册日期
	Longitude      float64     //经度
	Latitude       float64     //纬度
}

type ST_HosQualification struct {
	HosBL          string //营业执照号码BL:business license
	BusinessPeriod string //营业执照有效周期,例如:2010年12月01日~2020年11月10日
	HosBLBegin     string //营业执照开始
	HosBLEnd       string //营业执照结束
	HosBLPic       string //营业执照照片

	MedicalPLNumber string //医疗机构执业许可证号码;    PL:practice license 执业许可证
	MedicalPLPeriod string //医疗机构执业许可证有效周期; PL:practice license 执业许可证
	MedicalPLBegin  string //医疗机构执业许可开始时间
	MedicalPLEnd    string //医疗机构职业虚空结束时间
	MedicalPLPics   string //医疗机构营业许可证照片

	MedicalAdvPLNumber string //医疗机构广告宣传执业许可证号码
	MedicalAdvPLPeriod string //医疗机构广告宣传执业许可证有效周期
	MedicalAdvPLBegin  string //医疗机构广告开始时间
	MedicalAdvPLEnd    string //医疗机构广告结束时间
	MedicalAdvPLPics   string //医疗机构广告宣传执业许可证照片
}

type ST_HosStatistics struct {
	Orderquantity  int //订单数量
	Orderquota     int //订单总额度
	Reservepeople  int //预约的人数
	ProductNum     int //产品数量
	DoctorNum      int //医生数量
	Evaluatepeople int //评价人数
	VisitNum       int //访问量
	ConsultNum     int //咨询人数
	/////////////////
	Complaintquantity   int //投诉数量
	Refundratequantity  int //退款总数
	RefundQuata         int //退款额度
	Refundrate          int //退款率
	Activitydegree      int //活跃度
	CompositeScore      int //综合评分
	StarGrade           int //评分星级
	HosCommentNum       int //医院评价数量
	HosAttentionNum     int //医院关注数量
	EnvironmentScore    int //环境分
	ServiceAttitude     int //服务态度
	PostoperativeEffect int //术后效果
}

type ST_Hospital struct {
	HosID               string //医院id
	ST_HostBaseInfo            //医院的基本信息
	ST_HosQualification        //医疗机构的资质信息
	ST_OpreatStatus            //状态信息
	ST_HosStatistics           //统计信息
}

func RegisterHospital(session *JsNet.StSession) {
	type st_get struct {
		Hos       *ST_Hospital
		CheckCode string
		LoginCode string
	}

	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	ErrorLog("RegisterHospital = %v", st)
	if st.Hos == nil {
		ForwardEx(session, "1", nil, "RegisterHospital Hos is nil\n")
		return
	}
	if st.CheckCode == "" || st.LoginCode == "" {
		ForwardEx(session, "1", nil, "验证码或者密码不能为空,CheckCode=%s,LoginCode=%s\n", st.CheckCode, st.LoginCode)
		return
	}
	///医院信息检查
	if err := hosCheckData(st.Hos); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	//短信校验码检查
	if ok := JsMobile.VerifySmsCode(st.Hos.AdminCell, st.CheckCode); !ok {
		ForwardEx(session, "2", nil, "注册失败,验证码不正确\n")
		return
	}

	////检查手机号是否注册医院
	data, ok := checkcellMapExist(st.Hos.AdminCell)
	if ok && data != nil && data.HosID != "" {
		hos, err := GetHospitalInfo(data.HosID)
		if err == nil && hos != nil {
			ForwardEx(session, "1", hos, "手机号%s已经创建账号,并且已经注册医院\n", st.Hos.AdminCell)
			return
		}
	}
	/////////////医院id生成/////////////////////////
	if id, err := ider.GenHosID(); err == nil {
		st.Hos.HosID = id
	} else {
		ForwardEx(session, "1", nil, "医院id生成失败\n")
		return
	}

	///创建管理员账号
	accountinfo, err := creatAdminAccount(st.Hos.HosID, st.Hos.AdminName, st.Hos.AdminCell, st.LoginCode)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	st.Hos.AdminAccount = accountinfo.LoginAccount
	//////创建医院//////
	if err := creatHospital(st.Hos); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	hosaccount := &ST_HosAccount{
		Admin:    *accountinfo,
		Employee: []AccountInfo{},
	}
	//////创建医院的账户信息/////
	if err := writeToHosAccount(hosaccount); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	go GlobalAddNewHospital(st.Hos.HosID)
	Forward(session, "0", accountinfo)
}

////创建医院
func creatHospital(st *ST_Hospital) error {
	st.CreatDate = CurTime()
	/////添加操作记录
	appendHosStatus(st,
		constant.Hospital_Side, "新创建医院",
		constant.OperatingStatus_new, "无",
		"", "", "",
		st.AdminAccount,
		st.AdminName,
		st.AdminCell)
	st.CompositeScore = 500
	st.EnvironmentScore = 500
	st.ServiceAttitude = 500
	st.PostoperativeEffect = 500
	st.StarGrade = 5
	if err := DirectWrite(constant.Hash_Hospital, st.HosID, st); err != nil {
		return ErrorLog("医院创建失败,err:%s\n", err.Error())
	}
	return AddGlobalHospotal(st.HosID)
}

//修改医院的基本信息和医院相关信息
func ModifyHospitalInfo(session *JsNet.StSession) {
	type st_modify struct {
		HosID               string //医院id
		ST_HostBaseInfo            //医院的基本信息
		ST_HosQualification        //医疗机构的资质信息
		OpreatReason        string //操作原因
		OpreatJobNum        string //操作人员工号
		OpreatName          string //操作人员姓名
		OpreatCell          string //操作人手机号
	}
	st := &st_modify{}
	data := &ST_Hospital{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if err := session.GetPara(data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" {
		ForwardEx(session, "1", nil, "修改医院基本信息失败,HosID为空,医院ID=%s\n", st.HosID)
		return
	}
	if err := hosCheckData(data); err != nil {
		Error(err.Error())
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if err := WriteLock(constant.Hash_Hospital, st.HosID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	oldCity := data.HosCity

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
	appendHosStatus(data,
		constant.Hospital_Side,
		"修改医院资料",
		constant.OperatingStatus_modify,
		st.OpreatReason,
		"", "", "",
		st.OpreatJobNum,
		st.OpreatName,
		st.OpreatCell)

	if err != WriteBack(constant.Hash_Hospital, st.HosID, data) {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if oldCity != st.HosCity {
		Warn("发现医院:(%s)的将城市从%s更改为%s...\n", st.HosID, oldCity, st.HosCity)
		go GlobalUpdateHospitalCity()
	}
	go GlobalUpdateHosStatus()
	Forward(session, "0", data)
}

//查询医院信息
func QueryHospitalInfo(session *JsNet.StSession) {
	type st_query struct {
		HosID string //医院id
	}
	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" {
		ForwardEx(session, "1", nil, "医院id为空,QueryHospitalInfo()，查询失败!")
		return
	}
	data, err := GetHospitalInfo(st.HosID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

//查询多个医院信息
func QueryMoreHospital(session *JsNet.StSession) {
	type st_query struct {
		HosIDs []string //医院id
	}
	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data := QueryMoreHosInfo(st.HosIDs)
	Forward(session, "0", data)
}

///查询多个医院信息
func QueryMoreHosInfo(ids []string) []*ST_Hospital {
	data := []*ST_Hospital{}
	if len(ids) == 0 {
		return data
	}
	for _, v := range ids {
		t, err := GetHospitalInfo(v)
		if err != nil {
			continue
		}
		data = append(data, t)
	}
	return data
}

//审核医院
func RevieweNewHospital(session *JsNet.StSession) {
	type st_reviewe struct {
		HosID        string //医院id
		OpreatReason string //操作原因
		OpreatJobNum string //操作人员工号
		OpreatName   string //操作人员姓名
		OpreatCell   string //操作人联系方式
		IsPass       bool   //审核是否通过:true:通过，false:不通过
		//	Ratio        ST_SettlementRatioInfo //结算信息
	}
	st := &st_reviewe{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" || st.OpreatJobNum == "" {
		ForwardEx(session, "1", nil, "初次医院审核,操作参数不全,HosID=%s,OpreatJobNum=%s\n",
			st.HosID, st.OpreatJobNum)
		return
	}
	data := &ST_Hospital{}
	if err := WriteLock(constant.Hash_Hospital, st.HosID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	defer WriteBack(constant.Hash_Hospital, st.HosID, data)
	if data.Current.OpreatStatus == constant.OperatingStatus_online {
		Forward(session, "0", data)
		return
	}
	if data.Current.OpreatStatus != constant.OperatingStatus_new {
		ForwardEx(session, "1", data, "初次审核失败,当前医院状态为%s,不在待审核状态\n", data.Current.OpreatStatus)
		return
	}
	/////添加操作记录
	data.FirstOnlineDate = CurTime()
	data.LastOnlineDate = CurTime()
	data.LastOnlineName = st.OpreatName
	Len := len(data.OpreatInfo)
	status := constant.OperatingStatus_Reviewer_NotPass
	action := "初次医院系统审核不通过"
	if st.IsPass {
		status = constant.OperatingStatus_online
		action = "初次医院系统审核通过"
	}
	appendHosStatus(data,
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

	///医院的id和姓名对应表
	go AddHosNameList(data.HosID, data.HosName)
	go GlobalUpdateHosStatus()
	Forward(session, "0", data)
}

//审核医院
func RevieweModifyHospital(session *JsNet.StSession) {
	type st_reviewe struct {
		HosID        string //医院id
		OpreatReason string //操作原因
		OpreatJobNum string //操作人员工号
		OpreatName   string //操作人员姓名
		OpreatCell   string //操作人手机号
		IsPass       bool   //审核是否通过:true:通过，false:不通过
	}
	st := &st_reviewe{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" || st.OpreatJobNum == "" {
		ForwardEx(session, "1", nil, "医院审核,参数不全,HosID=%s,OpreatJobNum=%s\n",
			st.HosID, st.OpreatJobNum)
		return
	}
	data := &ST_Hospital{}
	if err := WriteLock(constant.Hash_Hospital, st.HosID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	defer WriteBack(constant.Hash_Hospital, st.HosID, data)

	if data.Current.OpreatStatus == constant.OperatingStatus_online {
		Forward(session, "0", data)
		return
	}
	if data.Current.OpreatStatus != constant.OperatingStatus_modify {
		ForwardEx(session, "1", data, "修改审核失败,当前医院状态为%s,不在修改状态\n", data.Current.OpreatStatus)
		return
	}
	/////添加操作记录
	data.LastOnlineDate = CurTime()
	data.LastOnlineName = st.OpreatName
	Len := len(data.OpreatInfo)
	status := constant.OperatingStatus_Reviewer_NotPass
	action := "医院修改后台审核不通过"

	if st.IsPass {
		status = constant.OperatingStatus_online
		action = "医院修改后台审核通过"
	}
	appendHosStatus(data,
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

	go GlobalUpdateHosStatus()
	Forward(session, "0", data)
}

//医院强制下线
func HosOfflineOnForce(session *JsNet.StSession) {
	type st_offline struct {
		HosID        string //医院id
		OpreatReason string //操作原因
		OpreatJobNum string //操作人员工号
		OpreatName   string //操作人员姓名
		OpreatCell   string //操作人手机号
	}

	st := &st_offline{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" || st.OpreatJobNum == "" {
		ForwardEx(session, "1", nil, "医院强制下线失败,HosID=%s,OpreatJobNum=%s\n",
			st.HosID, st.OpreatJobNum)
		return
	}
	data := &ST_Hospital{}

	if err := WriteLock(constant.Hash_Hospital, st.HosID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	defer WriteBack(constant.Hash_Hospital, st.HosID, data)
	if data.Current.OpreatStatus == constant.OperatingStatus_Offline_self ||
		data.Current.OpreatStatus == constant.OperatingStatus_Offline_onforce {
		Forward(session, "0", data)
		return
	}

	/////添加操作记录
	appendHosStatus(data,
		constant.Platform_Side,
		"医院被系统强制下线",
		constant.OperatingStatus_Offline_onforce,
		st.OpreatReason,
		st.OpreatJobNum,
		st.OpreatName,
		st.OpreatCell,
		"", "", "")

	go GlobalUpdateHosStatus()
	Forward(session, "0", data)
}

//医院删除
func DelHospital(session *JsNet.StSession) {
	type info struct {
		HosID string //医院id
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" {
		ForwardEx(session, "1", nil, "医院删除失败,HosID=%s\n", st.HosID)
		return
	}
	data := &ST_Hospital{}

	if err := WriteLock(constant.Hash_Hospital, st.HosID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}

	/////添加操作记录
	appendHosStatus(data,
		constant.Platform_Side,
		"医院被删除",
		constant.OperatingStatus_Del, "", "", "", "",
		"", "", "")

	if err := WriteBack(constant.Hash_Hospital, st.HosID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}

	go GlobalUpdateHosStatus()
	Forward(session, "0", data)
}

//更改医院访问量
func ChangeHosVisitNum(session *JsNet.StSession) {
	type info struct {
		HosID string //医院id
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" {
		ForwardEx(session, "1", nil, "更改医院访问量失败,HosID=%s\n", st.HosID)
		return
	}
	data := &ST_Hospital{}
	if err := WriteLock(constant.Hash_Hospital, st.HosID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	data.VisitNum++
	if err := WriteBack(constant.Hash_Hospital, st.HosID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	Forward(session, "0", data)
}

//增加医院访问量
func ChangeHosConsultNum(session *JsNet.StSession) {
	type info struct {
		HosID string //医院id
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" {
		ForwardEx(session, "1", nil, "更改医院咨询量失败,HosID=%s\n", st.HosID)
		return
	}
	data := &ST_Hospital{}
	if err := WriteLock(constant.Hash_Hospital, st.HosID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	data.ConsultNum++
	if err := WriteBack(constant.Hash_Hospital, st.HosID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	Forward(session, "0", data)
}

//增加医院医生个数
func ChangeHosDocNum(HosID string) error {
	if HosID == "" {
		return ErrorLog("更改医院医生数量失败,HosID=%s\n", HosID)
	}
	data := &ST_Hospital{}
	if err := WriteLock(constant.Hash_Hospital, HosID, data); err != nil {
		return err
	}
	data.DoctorNum++
	if err := WriteBack(constant.Hash_Hospital, HosID, data); err != nil {
		return err
	}
	return nil
}

//增加医院产品个数
func ChangeHosProNum(HosID string) error {
	if HosID == "" {
		return ErrorLog("更改医院产品数量失败,HosID=%s\n", HosID)
	}
	data := &ST_Hospital{}
	if err := WriteLock(constant.Hash_Hospital, HosID, data); err != nil {
		return err
	}
	data.ProductNum++
	if err := WriteBack(constant.Hash_Hospital, HosID, data); err != nil {
		return err
	}
	return nil
}

//医院关注
func hospital_follow_people(HosID string, Follow bool) error {
	data := &ST_Hospital{}
	if HosID == "" {
		return ErrorLog("医院关注人数失败,HosID is empty\n")
	}
	if err := WriteLock(constant.Hash_Hospital, HosID, data); err != nil {
		return err
	}
	if Follow {
		data.HosAttentionNum++
	} else {
		if data.HosAttentionNum > 0 {
			data.HosAttentionNum--
		}
		if data.Reservepeople < 0 {
			data.HosAttentionNum = 0
		}
	}
	return WriteBack(constant.Hash_Hospital, HosID, data)
}

//医院下单
func hospital_reserve_people(HosID string, Orderquota int) error {
	data := &ST_Hospital{}
	if HosID != "" && Orderquota > 0 {
		if err := WriteLock(constant.Hash_Hospital, HosID, data); err != nil {
			return err
		}
		data.Reservepeople++
		data.Orderquantity++
		data.Orderquota += Orderquota
		if err := WriteBack(constant.Hash_Hospital, HosID, data); err != nil {
			return err
		}
		if data != nil {
			go GlobalUpdateHosSale(HosID, data.Reservepeople)
		}
	}
	return nil
}

///内部获取医院信息
func GetHospitalInfo(HosID string) (*ST_Hospital, error) {
	data := &ST_Hospital{}
	if err := ShareLock(constant.Hash_Hospital, HosID, data); err != nil {
		Error(err.Error())
		return nil, ErrorLog("获取医院信息失败ShareLock(),HosID=%s\n", HosID)
	}
	return data, nil
}

//获取医院的id和名字的映射表
func GetHosSimpleInfo(session *JsNet.StSession) {
	type st_HosSimpleInfo struct {
		HosID   string
		HosName string
	}
	v := st_HosSimpleInfo{}
	v.HosID = "0000"
	v.HosName = "全部医院"

	data := []st_HosSimpleInfo{}
	if err := ShareLock(constant.Hash_HospitalCache, constant.KEY_HosSimpleInfo, &data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data = append(data, v)
	Forward(session, "0", data)
}

///医院的id和姓名对应表
func AddHosNameList(HosID, HosName string) error {
	type st_HosSimpleInfo struct {
		HosID   string
		HosName string
	}
	st := []st_HosSimpleInfo{}
	if err := WriteLock(constant.Hash_HospitalCache, constant.KEY_HosSimpleInfo, &st); err != nil {
		return err
	}

	exist := false
	for _, v := range st {
		if v.HosID == HosID {
			exist = true
			break
		}
	}
	if !exist {
		st = append(st, st_HosSimpleInfo{HosID: HosID, HosName: HosName})
	}
	return WriteBack(constant.Hash_HospitalCache, constant.KEY_HosSimpleInfo, &st)
}

func GlobalHospitalList() []string {
	hosid := []string{}
	if err := ShareLock(constant.Hash_HospitalCache, constant.KEY_ALL_Hospital, &hosid); err != nil {
		Error(err.Error())
		return []string{}
	}
	return hosid
}

func AddGlobalHospotal(HosID string) error {
	hosid := []string{}
	err := WriteLock(constant.Hash_HospitalCache, constant.KEY_ALL_Hospital, &hosid)
	if err != nil {
		hosid = append(hosid, HosID)
		return DirectWrite(constant.Hash_HospitalCache, constant.KEY_ALL_Hospital, &hosid)
	}
	exist := false
	for _, v := range hosid {
		if v == HosID {
			exist = true
			break
		}
	}
	if !exist {
		hosid = append(hosid, HosID)
	}
	return WriteBack(constant.Hash_HospitalCache, constant.KEY_ALL_Hospital, &hosid)
}

//往医院里增加一条操作记录
func appendHosStatus(hos *ST_Hospital, OpreatPart, OpreatAction, OpreatStatus,
	OpreatReason, OpreatJobNum, OpreatName, OpreatCell,
	ApplyJobNum, ApplyName, ApplyCell string) {
	info := ST_Opreat{
		OpreatPart:   OpreatPart,
		OpreatAction: OpreatAction,
		OpreatStatus: OpreatStatus,
		OpreatReason: OpreatReason,
		OpreatJobNum: OpreatJobNum,
		OpreatName:   OpreatName,
		ApplyJobNum:  ApplyJobNum,
		ApplyName:    ApplyName,
		ApplyCell:    ApplyCell,
		OpreatTime:   CurTime(),
	}
	hos.Current = info
	hos.OpreatInfo = append(hos.OpreatInfo, info)
}

//输入检查
func hosCheckData(st *ST_Hospital) error {

	if st.HosName == "" {
		return ErrorLog("医院的名字不能为空!")
	}
	if st.HosCity == "" || st.HosAddr == "" {
		return ErrorLog("医院的城市或者详细地址不能为空!")
	}
	if st.HosBL == "" {
		return ErrorLog("医院的营业执照号码不能为空!")
	}
	if st.HosBLBegin == "" || st.HosBLEnd == "" {
		return ErrorLog("医院的营业执照号码开始时间或结束时间不能为空!")
	}
	if st.MedicalPLBegin == "" || st.MedicalPLEnd == "" {
		return ErrorLog("医疗机构执业许可开始时间或结束时间不能为空!")
	}
	if st.MedicalAdvPLBegin == "" || st.MedicalAdvPLEnd == "" {
		return ErrorLog("医疗机构广告宣传执业开始时间或结束时间不能为空!")
	}
	if st.LRName == "" || st.LRNumber == "" {
		return ErrorLog("医院法人代表的姓名或者号码不能为空!")
	}
	if st.ContactsName == "" || st.ContactsNumber == "" {
		return ErrorLog("医院联系人姓名或者号码不能为空!")
	}
	if st.MedicalPLNumber == "" {
		return ErrorLog("医疗机构执业许可证编号不能为空!")
	}
	if st.HosCapitalType == "" {
		return ErrorLog("医院资本类型不能为空!")
	}
	if st.ServiceType == "" {
		return ErrorLog("医院服务类型不能为空!")
	}
	if st.ServiceLevel == "" {
		return ErrorLog("医院服务等级不能为空!")
	}
	if st.HosType == "" {
		return ErrorLog("医院类型不能为空!")
	}
	if len(st.ServiceRange) == 0 {
		return ErrorLog("医院服务范围不能为空!")
	}
	if st.AdminCell == "" {
		return ErrorLog("医院管理员手机号不能为空!")
	}
	if st.AdminName == "" {
		return ErrorLog("医院管理员姓名不能为空!")
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////

//医院退单
func hospital_refundrate(HosID string, RefundQua int) error {
	data := &ST_Hospital{}
	if HosID == "" || RefundQua < 0 {
		return ErrorLog("hospital_refundrate修改退款率失败，HosID=%s,RefundQua=%d\n", HosID, RefundQua)
	}

	if err := WriteLock(constant.Hash_Hospital, HosID, data); err != nil {
		return err
	}
	data.RefundQuata += RefundQua
	data.Refundratequantity++
	data.Refundrate = (data.Refundratequantity * 100) / data.Orderquantity
	data.Reservepeople--
	if data.Reservepeople < 0 {
		data.Reservepeople = 0
	}
	return WriteBack(constant.Hash_Hospital, HosID, data)
}

//医院评价
func hospital_evaluate_people(HosID string, EnvironmentScore, ServiceAttitude, PostoperativeEffect int) error {
	data := &ST_Hospital{}
	if HosID != "" {

		if err := WriteLock(constant.Hash_Hospital, HosID, data); err != nil {
			return err
		}
		en := data.Evaluatepeople*data.EnvironmentScore + EnvironmentScore
		ser := data.Evaluatepeople*data.ServiceAttitude + ServiceAttitude
		eff := data.Evaluatepeople*data.PostoperativeEffect + PostoperativeEffect
		data.HosCommentNum++

		data.EnvironmentScore = en / data.Evaluatepeople
		data.ServiceAttitude = ser / data.Evaluatepeople
		data.PostoperativeEffect = eff / data.Evaluatepeople
		data.CompositeScore = (data.EnvironmentScore*4 + data.ServiceAttitude*4 + data.PostoperativeEffect*2) / 10

		data.StarGrade = data.CompositeScore / 100
		if data.StarGrade == 0 {
			data.StarGrade = 1
		}
		if err := WriteBack(constant.Hash_Hospital, HosID, data); err != nil {
			return err
		}
		if data != nil {
			go GlobalUpdateHosComment(HosID, data.Evaluatepeople)
		}
	}
	return nil
}

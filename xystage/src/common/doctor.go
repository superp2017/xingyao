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

//医生基本信息
type ST_DocBaseInfo struct {
	DocName             string   //医生姓名
	DocSex              string   //医生性别
	DocMobile           string   //医生手机
	DocHead             string   //医生头像
	DocEmail            string   //医生邮箱
	DocNationality      string   //医生国籍
	DocPhotos           []string //医生的形象照
	DocAge              string   //医生年龄
	DocIntroduce        string   //医生详细介绍
	WorkingTime         string   //医生从医时间
	Programs            []string //医生擅长的项目
	HisPrograms         []string //医生关联的项目
	HotPrograms         []string //热门项目
	IsExpertDoctor      bool     //是否是大牌名医
	ExpertDocPic        string   //大牌名医形象照
	ExpertDocDes        string   //大牌名医介绍
	ExpertDocBackground string   //大牌名医背景
	Hospital            string   //医生所在医院
	HosCity             string   //所在医院的城市
	DocPost             string   //医生所在医院的职务
}

//医生的资质信息
type ST_DocQualification struct {
	DocDuty       string //医生职务：院长、主任、医生
	DocEdu        string //医生学历Education：专科、本科、硕士、博士、博士后
	AcademicGrade string //医生学术等级:助教、讲师、副教授、教授
	DocTitle      string //医生职称Title:主任医生、副主任医生、主治医师、住院医师
	DocTitlePic   string //医生职称证书
	DocQuaLevel   string //医生职业资格等级qualification level：助理执业医师、执业医师、助理执业中医师、执业中医师
	DocQua        string //医生资i质 Qualifications
	DocQC         string //医生执业资格证书Qualification certificate
	DocCMC        string //医疗美容职业主诊资格证书Certificate of medical cosmetology
}

type ST_DocService struct {
	ServiceAttitude     int      //服务态度
	PostoperativeEffect int      //术后效果
	EnvironmentScore    int      //环境分
	GeneralEvaluation   int      //总体评价
	StarGrade           int      //评分星级
	Evaluatepeople      int      //评价人数
	ActivityLevel       int      //活跃程度
	Reservepeople       int      //预定人数
	DocAttentionNum     int      //医生关注人数
	Orders              []string //关联的订单
	VisitNum            int      //访问量
	ConsultNum          int      //咨询量
}

//医生
type ST_Doctor struct {
	DocID               string //医生ID
	HosID               string //医院ID
	ST_DocBaseInfo             //医生基本信息
	ST_DocQualification        //医生资质信息
	ST_OpreatStatus            //医生的操作信息
	ST_DocService              //医生的服务评价
}

//创建一个新的医生
func NewDoc(session *JsNet.StSession) {
	st := &ST_Doctor{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error()+",data=%v", st)
		return
	}
	if err := docCheckData(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	st.CreatDate = CurTime()
	if id, err := ider.GenDocID(); err == nil {
		st.DocID = id
	} else {
		ForwardEx(session, "1", nil, "医生ID生成失败!\n")
		return
	}
	hos, err := GetHospitalInfo(st.HosID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	st.HosCity = hos.HosCity
	st.GeneralEvaluation = 500
	st.ServiceAttitude = 500
	st.PostoperativeEffect = 500
	st.EnvironmentScore = 500
	st.StarGrade = 5
	/////添加操作记录
	appendDocStatus(st,
		constant.Hospital_Side,
		"新增医生",
		constant.OperatingStatus_new,
		"无",
		"", "", "",
		hos.AdminAccount,
		hos.AdminName,
		hos.AdminCell)
	if err := DirectWrite(constant.Hash_Doctor, st.DocID, st); err != nil {
		ForwardEx(session, "1", nil, "医生创建失败!\n"+err.Error())
		return
	}
	go AddHosDoctor(st.HosID, st.DocID)
	go GolbalAddNewDoctor(st.DocID, st.HosID, st.HosCity, st.DocName)
	go ChangeHosDocNum(st.HosID)
	Forward(session, "0", st)
}

//医生审核
func RevieweNewDoctor(session *JsNet.StSession) {
	type st_reviewe struct {
		HosID        string //医院id
		DocID        string //医生id
		OpreatReason string //操作原因
		OpreatJobNum string //操作人员工号
		OpreatName   string //操作人员姓名
		OpreatCell   string //操作人联系方式
		IsPass       bool   //审核是否通过:true:通过，false:不通过
	}
	st := &st_reviewe{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" || st.DocID == "" || st.OpreatJobNum == "" {
		ForwardEx(session, "1", nil,
			"医生审核失败,参数不全,HosID=%s,DocID=%s,OpreatJobNum=%s,HosCity=%s\n",
			st.HosID, st.DocID, st.OpreatJobNum)
		return
	}
	data := &ST_Doctor{}

	if err := WriteLock(constant.Hash_Doctor, st.DocID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_Doctor, st.DocID, data)

	if data.Current.OpreatStatus == constant.OperatingStatus_online {
		Forward(session, "0", data)
		return
	}
	if data.Current.OpreatStatus != constant.OperatingStatus_new {
		ForwardEx(session, "1", data, "初次审核失败,当前医生状态为%s,不在待审核状态\n", data.Current.OpreatStatus)
		return
	}

	/////添加操作记录
	Len := len(data.OpreatInfo)
	data.FirstOnlineDate = CurTime()
	data.LastOnlineDate = CurTime()
	data.LastOnlineName = st.OpreatName
	status := constant.OperatingStatus_Reviewer_NotPass
	action := "新增医生审核不通过"
	if st.IsPass {
		status = constant.OperatingStatus_online
		action = "新增医生审核通过"
	}
	appendDocStatus(data,
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

	Forward(session, "0", data)
	go GolbalUpdateDoc(st.HosID)
}

//医生审核
func RevieweModifyDoctor(session *JsNet.StSession) {
	type st_reviewe struct {
		HosID        string //医院id
		DocID        string //医生id
		OpreatReason string //操作原因
		OpreatJobNum string //操作人员工号
		OpreatName   string //操作人员姓名
		OpreatCell   string //操作人联系方式
		IsPass       bool   //审核是否通过:true:通过，false:不通过
	}
	st := &st_reviewe{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.HosID == "" || st.DocID == "" || st.OpreatJobNum == "" {
		ForwardEx(session, "1", nil,
			"医生审核失败,参数不全,HosID=%s,DocID=%s,OpreatJobNum=%s,HosCity=%s\n",
			st.HosID, st.DocID, st.OpreatJobNum)
		return
	}
	data := &ST_Doctor{}

	if err := WriteLock(constant.Hash_Doctor, st.DocID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_Doctor, st.DocID, data)

	if data.Current.OpreatStatus == constant.OperatingStatus_online {
		Forward(session, "0", data)
		return
	}
	if data.Current.OpreatStatus != constant.OperatingStatus_modify {
		ForwardEx(session, "1", nil, "修改审核失败,当前医生状态为%s,不在修改状态\n", data.Current.OpreatStatus)
		return
	}
	data.LastOnlineDate = CurTime()
	data.LastOnlineName = st.OpreatName
	Len := len(data.OpreatInfo)
	status := constant.OperatingStatus_Reviewer_NotPass
	action := constant.OperatingStatus_Reviewer_NotPass
	if st.IsPass {
		status = constant.OperatingStatus_online
		action = constant.OperatingStatus_online
	}
	appendDocStatus(data,
		constant.Platform_Side,
		action, status,
		st.OpreatReason,
		st.OpreatJobNum,
		st.OpreatName,
		st.OpreatCell,
		data.OpreatInfo[Len-1].OpreatJobNum,
		data.OpreatInfo[Len-1].OpreatName,
		data.OpreatInfo[Len-1].OpreatCell)

	Forward(session, "0", data)
	go GolbalUpdateDoc(st.HosID)

}

//查询医生
func QueryDocInfo(session *JsNet.StSession) {
	type st_query struct {
		DocID string //医生id
	}
	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.DocID == "" {
		ForwardEx(session, "1", "获取医生信息失败,DocID 为空,DocID=%s\n", st.DocID)
		return
	}
	data, err := QueryDoctor(st.DocID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

//查询多个医生信息
func QueryMoreDoctors(session *JsNet.StSession) {
	type st_query struct {
		DocIDs []string //医生ids
	}
	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", QueryMoreDocInfo(st.DocIDs))
}

///查询多个医生信息
func QueryMoreDocInfo(ids []string) []*ST_Doctor {
	data := []*ST_Doctor{}
	if len(ids) == 0 {
		return data
	}
	for _, v := range ids {
		t, e := QueryDoctor(v)
		if e != nil {
			continue
		}
		data = append(data, t)
	}
	return data
}

//修改医生信息
func ModifyDocInfo(session *JsNet.StSession) {
	type st_baseinfo struct {
		DocID               string //医生工号
		ST_DocBaseInfo             //医生基本信息
		ST_DocQualification        //医生资质信息
		OpreatReason        string //操作原因
		OpreatJobNum        string //操作人员工号
		OpreatName          string //操作人员姓名
		OpreatCell          string //操作人员手机号
	}
	st := &st_baseinfo{}
	data := &ST_Doctor{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if err := session.GetPara(data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.DocID == "" {
		ForwardEx(session, "1", "修改获取医生信息失败,DocID 为空,DocID=%s\n", st.DocID)
		return
	}

	if err := docCheckData(data); err != nil {
		Error(err.Error())
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if err := WriteLock(constant.Hash_Doctor, st.DocID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	defer WriteBack(constant.Hash_Doctor, st.DocID, data)

	d, err := json.Marshal(st)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if err1 := json.Unmarshal(d, data); err1 != nil {
		ForwardEx(session, "1", nil, err1.Error())
		return
	}

	data.LastModifyDate = CurTime()
	appendDocStatus(data,
		constant.Hospital_Side,
		"修改医生资料",
		constant.OperatingStatus_modify,
		st.OpreatReason,
		"", "", "",
		st.OpreatJobNum,
		st.OpreatName,
		st.OpreatCell)

	if data != nil {
		go GolbalUpdateDoc(data.HosID)
	}
	Forward(session, "0", data)
}

//修改医生
func ModifyDocCity(DocID, HosCity string) error {
	if DocID == "" || HosCity == "" {
		return ErrorLog("ModifyDocCity failed,DocID=%s,HosCity=%s", DocID, HosCity)
	}
	data := &ST_Doctor{}
	if err := WriteLock(constant.Hash_Doctor, DocID, data); err != nil {
		return err
	}
	defer WriteBack(constant.Hash_Doctor, DocID, data)
	data.HosCity = HosCity

	go removeHomeExpertDoctor(DocID, HosCity)
	go removeExpertDoctor(DocID, HosCity)
	return nil
}

//添加一个关联一个产品到医生
func AddProjectToDoc(DocID, ProID string) error {
	if DocID == "" || ProID == "" {
		return ErrorLog("AddProjectToDoc faild,param is empty,DocID=%s,ProID=%s\n", DocID, ProID)
	}
	data := &ST_Doctor{}
	if err := WriteLock(constant.Hash_Doctor, DocID, data); err != nil {
		return err
	}
	exist := false
	for _, v := range data.HisPrograms {
		if v == ProID {
			exist = true
			break
		}
	}
	if !exist {
		data.HisPrograms = append(data.HisPrograms, ProID)
	}
	exist = false
	for _, v := range data.HotPrograms {
		if v == ProID {
			exist = true
			break
		}
	}
	if !exist {
		data.HotPrograms = append(data.HotPrograms, ProID)
	}
	return WriteBack(constant.Hash_Doctor, DocID, data)
}

///提升至大牌名医
func UpToExpertDoc(session *JsNet.StSession) {
	type st_UP struct {
		DocID        string //医生id
		ExpertDocPic string //大牌名医形象照
		ExpertDocDes string //大牌名医介绍
	}
	st := &st_UP{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.DocID == "" || st.ExpertDocDes == "" || st.ExpertDocPic == "" {
		ForwardEx(session, "1", nil, "更改大排名医失败,DocID=%s\n", st.DocID)
		return
	}
	data := &ST_Doctor{}
	if err := WriteLock(constant.Hash_Doctor, st.DocID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_Doctor, st.DocID, data)
	//////// 更改医生的大牌名医状态///
	if err := addExpertDoctor(data.HosCity, st.DocID); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.IsExpertDoctor = true
	data.ExpertDocDes = st.ExpertDocDes
	data.ExpertDocPic = st.ExpertDocPic
	Forward(session, "0", data)
}

///下降大牌名医
func DownExpertDoc(session *JsNet.StSession) {
	type st_offline struct {
		DocID string //医生id
	}
	st := &st_offline{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.DocID == "" {
		ForwardEx(session, "1", nil, "更改大排名医失败,DocID=%s\n", st.DocID)
		return
	}
	data := &ST_Doctor{}
	if err := WriteLock(constant.Hash_Doctor, st.DocID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_Doctor, st.DocID, data)

	if err := removeExpertDoctor(st.DocID, data.HosCity); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if err := removeHomeExpertDoctor(data.HosCity, st.DocID); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.IsExpertDoctor = false
	Forward(session, "0", data)
}

///更改是否是大排名医
func changeDocExpertPhoto(DocID, ExpertDocBackground string) error {
	if DocID == "" || ExpertDocBackground == "" {
		return ErrorLog("ChangeDocPhoto faild,DocID=%s,ExpertDocBackground=%s\n", DocID, ExpertDocBackground)
	}
	data := &ST_Doctor{}

	if err := WriteLock(constant.Hash_Doctor, DocID, data); err != nil {
		return err
	}
	data.ExpertDocBackground = ExpertDocBackground
	return WriteBack(constant.Hash_Doctor, DocID, data)
}

//医院将医生下线
func DocOfflineSelf(session *JsNet.StSession) {
	type st_offline struct {
		HosID        string //医院id
		DocID        string //医生id
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
	if st.DocID == "" || st.HosID == "" {
		ForwardEx(session, "1", nil, "医院主动下线医生失败,DocID=%s,HosID=%s\n",
			st.DocID, st.HosID)
		return
	}
	data := &ST_Doctor{}

	if err := WriteLock(constant.Hash_Doctor, st.DocID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_Doctor, st.DocID, data)

	if data.Current.OpreatStatus == constant.OperatingStatus_Offline_self ||
		data.Current.OpreatStatus == constant.OperatingStatus_Offline_onforce {
		Forward(session, "0", data)
		return
	}
	appendDocStatus(data,
		constant.Hospital_Side,
		"医院将医生下架",
		constant.OperatingStatus_Offline_self,
		st.OpreatReason,
		"", "", "",
		st.OpreatJobNum,
		st.OpreatName,
		st.OpreatCell)

	if data != nil {
		go removeHomeExpertDoctor(data.DocID, data.HosCity)
		go removeExpertDoctor(data.DocID, data.HosCity)
	}
	go GolbalUpdateDoc(st.HosID)
	Forward(session, "0", data)
}

///内部强制下线
func OfflineDocOnForce(HosID, DocID string) error {
	if HosID == "" || DocID == "" {
		return ErrorLog("OfflineDocOnForce failed,HosID=%s,DocID=%s\n", HosID, DocID)
	}
	data := &ST_Doctor{}

	if err := WriteLock(constant.Hash_Doctor, DocID, data); err != nil {
		return err
	}
	defer WriteBack(constant.Hash_Doctor, DocID, data)

	if data.Current.OpreatStatus == constant.OperatingStatus_Offline_self ||
		data.Current.OpreatStatus == constant.OperatingStatus_Offline_onforce {
		return nil
	}
	appendDocStatus(data,
		constant.Platform_Side,
		"医院下架,相关联的医生同步下架",
		constant.OperatingStatus_Offline_onforce,
		"", "", "", "", "", "", "")
	if data != nil {
		go removeHomeExpertDoctor(data.DocID, data.HosCity)
		go removeExpertDoctor(data.DocID, data.HosCity)
	}
	return nil
}

//平台强制下线医生
func DocOfflineOnForce(session *JsNet.StSession) {
	type st_offline struct {
		DocID        string //医生id
		HosID        string //医院id
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
	if st.DocID == "" || st.HosID == "" || st.OpreatJobNum == "" {
		ForwardEx(session, "1", nil, "医生强制下线失败,DocID=%s,HosID=%s,OpreatJobNum=%s\n",
			st.DocID, st.HosID, st.OpreatJobNum)
		return
	}
	data := &ST_Doctor{}

	if err := WriteLock(constant.Hash_Doctor, st.DocID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	defer WriteBack(constant.Hash_Doctor, st.DocID, data)

	if data.Current.OpreatStatus == constant.OperatingStatus_Offline_self ||
		data.Current.OpreatStatus == constant.OperatingStatus_Offline_onforce {
		Forward(session, "0", data)
		return
	}

	appendDocStatus(data,
		constant.Hospital_Side,
		"平台将医生下架",
		constant.OperatingStatus_Offline_onforce,
		st.OpreatReason,
		st.OpreatJobNum,
		st.OpreatName,
		st.OpreatCell,
		"", "", "")

	if data != nil {
		go removeHomeExpertDoctor(data.DocID, data.HosCity)
		go removeExpertDoctor(data.DocID, data.HosCity)
	}
	go GolbalUpdateDoc(st.HosID)
	Forward(session, "0", data)
}

//平台删除医生
func DelDoctor(session *JsNet.StSession) {
	type info struct {
		DocID string //医生id
		HosID string //医院id
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.DocID == "" || st.HosID == "" {
		ForwardEx(session, "1", nil, "删除医生失败,DocID=%s,HosID=%s\n", st.DocID, st.HosID)
		return
	}
	data := &ST_Doctor{}
	if err := WriteLock(constant.Hash_Doctor, st.DocID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	appendDocStatus(data,
		constant.Hospital_Side, "平台将医生删除",
		constant.OperatingStatus_Del,
		"", "", "", "",
		"", "", "")

	if err := WriteBack(constant.Hash_Doctor, st.DocID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	if data != nil {
		go removeHomeExpertDoctor(data.DocID, data.HosCity)
		go removeExpertDoctor(data.DocID, data.HosCity)
	}
	go GolbalUpdateDoc(st.HosID)
	Forward(session, "0", data)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////
//查询医生信息
func QueryDoctor(DocID string) (*ST_Doctor, error) {
	data := &ST_Doctor{}
	if err := ShareLock(constant.Hash_Doctor, DocID, data); err != nil {
		Error(err.Error())
		return nil, ErrorLog("获取医生信息失败ShareLock(),DocID=%s\n", DocID)
	}
	return data, nil
}

//获取某一医院所有的医生
func GetHosDoctorList(HosID string) []string {
	docids := []string{}
	if HosID == "" {
		return docids
	}
	if err := ShareLock(constant.Hash_DoctorCache, HosID, &docids); err != nil {
		Error(err.Error())
		return []string{}
	}
	return docids
}

//添加一个医生到医院
func AddHosDoctor(HosID, DocID string) error {
	if HosID == "" || DocID == "" {
		return ErrorLog("AddHosDoctor param failed,HosID =%S,DocID=%s\n", HosID, DocID)
	}
	docids := []string{}
	err := WriteLock(constant.Hash_DoctorCache, HosID, &docids)
	if err != nil {
		docids = append(docids, DocID)
		return DirectWrite(constant.Hash_DoctorCache, HosID, &docids)
	}
	exist := false
	for _, v := range docids {
		if v == DocID {
			exist = true
			break
		}
	}
	if !exist {
		docids = append(docids, DocID)
	}
	return WriteBack(constant.Hash_DoctorCache, HosID, &docids)
}

//往医院里增加一条操作记录
func appendDocStatus(doc *ST_Doctor, OpreatPart,
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

//更改医生的访问量
func ChangeDocVisitNum(session *JsNet.StSession) {
	type info struct {
		DocID string //医生id
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.DocID == "" {
		ForwardEx(session, "1", nil, "医生访问量修改失败,DocID=%s\n", st.DocID)
		return
	}
	data := &ST_Doctor{}

	if err := WriteLock(constant.Hash_Doctor, st.DocID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.VisitNum++
	if err := WriteBack(constant.Hash_Doctor, st.DocID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}


//更改医生的访问量
func ChangeDocConsultNum(session *JsNet.StSession) {
	type info struct {
		DocID string //医生id
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.DocID == "" {
		ForwardEx(session, "1", nil, "医生咨询量修改失败,DocID=%s\n", st.DocID)
		return
	}
	data := &ST_Doctor{}

	if err := WriteLock(constant.Hash_Doctor, st.DocID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.ConsultNum++
	if err := WriteBack(constant.Hash_Doctor, st.DocID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}






//预定人数
func doctor_reserve_people(DocID, OrderID string, Reserve bool) error {
	data := &ST_Doctor{}

	if DocID != "" {
		if err := WriteLock(constant.Hash_Doctor, DocID, data); err != nil {
			return err
		}
		if Reserve {
			data.Reservepeople++
			///关联一个订单
			data.Orders = append(data.Orders, OrderID)
		} else {
			data.Reservepeople--
			if data.Reservepeople < 0 {
				data.Reservepeople = 0
			}

		}
		if data != nil {
			go GlobalUpdateDocSale(DocID, data.Reservepeople)
		}
		return WriteBack(constant.Hash_Doctor, DocID, data)
	}
	return nil
}

//评价人数
func doctor_evaluate_people(DocID string, EnvironmentScore, ServiceAttitude, PostoperativeEffect int) error {
	data := &ST_Doctor{}
	if DocID != "" {
		if err := WriteLock(constant.Hash_Doctor, DocID, data); err != nil {
			return err
		}
		en := data.Evaluatepeople*data.EnvironmentScore + EnvironmentScore
		ser := data.Evaluatepeople*data.ServiceAttitude + ServiceAttitude
		eff := data.Evaluatepeople*data.PostoperativeEffect + PostoperativeEffect
		data.Evaluatepeople++
		data.EnvironmentScore = en / data.Evaluatepeople
		data.ServiceAttitude = ser / data.Evaluatepeople
		data.PostoperativeEffect = eff / data.Evaluatepeople
		data.GeneralEvaluation = (data.EnvironmentScore*4 + data.ServiceAttitude*4 + data.PostoperativeEffect*2) / 10
		data.StarGrade = data.GeneralEvaluation / 100
		if data.StarGrade == 0 {
			data.StarGrade = 1
		}

		if data != nil {
			go GlobalUpdateDocComment(DocID, data.Reservepeople)
		}
		return WriteBack(constant.Hash_Doctor, DocID, data)
	}
	return nil
}

//医生关注人数
func doctor_follow_people(DocID string, Follow bool) error {
	data := &ST_Doctor{}
	if DocID == "" {
		return ErrorLog("doctor_follow_people failed,DocID is empty\n")
	}
	if err := WriteLock(constant.Hash_Doctor, DocID, data); err != nil {
		return err
	}
	if Follow {
		data.DocAttentionNum++
	} else {
		if data.DocAttentionNum > 0 {
			data.DocAttentionNum--
		}
		if data.DocAttentionNum < 0 {
			data.DocAttentionNum = 0
		}
	}
	return WriteBack(constant.Hash_Doctor, DocID, data)
}

//输入检查
func docCheckData(doc *ST_Doctor) error {

	if doc.HosID == "" {
		return ErrorLog("医院id不能为空!\n")
	}

	if doc.DocName == "" {
		return ErrorLog("医生名字不能为空!\n")
	}

	//if doc.WorkingTime == "" {
	//	return ErrorLog("医生工作年限不能为空!\n")
	//}

	if doc.Hospital == "" {
		return ErrorLog("医生所在医院不能为空!\n")
	}
	if len(doc.Programs) == 0 {
		return ErrorLog("医生擅长的项目不能为空!\n")
	}

	//if doc.DocQuaLevel == "" {
	//	return ErrorLog("医生职业资格等级不能为空!\n")
	//}

	// if doc.AcademicGrade == "" {
	// 	return ErrorLog("医生学术等级不能为空!\n")
	// }
	// if doc.DocNationality == "" {
	// 	return ErrorLog("医生国籍不能为空!\n")
	// }

	// if doc.DocTitle == "" {
	// 	return ErrorLog("医生职称不能为空!\n")
	// }

	// if doc.DocCMC == "" {
	// 	return ErrorLog("医疗美容职业主诊资格证书不能为空!\n")
	// }
	// if doc.DocSex == "" {
	// 	return ErrorLog("医生性别不能为空!\n")
	// }
	// if doc.DocQC == "" {
	// 	return ErrorLog("医生执业资格证书不能为空!\n")

	return nil
}

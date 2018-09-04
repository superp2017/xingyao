package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	. "util"
)

type SM_Collection struct {
	Uid       string   //用户ID
	UserName  string   //用户
	UserHead  string   //头像
	CreatTime string   //创建时间
	Hos       []string //收藏医院
	Doc       []string //收藏医生
	Por       []string //收藏产品
}

func newUserFav(UID string) error {
	data := &SM_Collection{Uid: UID}
	return DirectWrite(constant.Hash_User_Fav, UID, data)
}

///查询用户的喜好
func QueryUserFav(UID string) (*SM_Collection, error) {
	if UID == "" {
		return nil, ErrorLog("QueryUserFav failed,UID is empty\n")
	}
	data := &SM_Collection{}
	err := ShareLock(constant.Hash_User_Fav, UID, data)
	return data, err
}

//用户的喜好
func UseFav(session *JsNet.StSession) {
	type st_get struct {
		UID string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" {
		ForwardEx(session, "1", nil, "用户的喜好失败,UID=%s\n", st.UID)
		return
	}
	data := &SM_Collection{}
	if err := ShareLock(constant.Hash_User_Fav, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	list, err := QueryUserFav(data.Uid)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", list)
}

//关注医院
func FollowHospitial(session *JsNet.StSession) {
	type st_get struct {
		HosID    string
		UID      string
		IsAttent bool
	}
	st := &st_get{}

	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" || st.HosID == "" {
		ForwardEx(session, "1", nil, "关注医院失败,参数不完整,Uid=%s,HosID=%s",
			st.UID, st.HosID)
		return
	}
	data := &SM_Collection{}

	if err := WriteLock(constant.Hash_User_Fav, st.UID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	if data.CreatTime == "" {
		data.CreatTime = CurTime()
	}
	exist := false
	index := -1
	for i, v := range data.Hos {
		if v == st.HosID {
			exist = true
			index = i
			break
		}
	}
	if st.IsAttent {
		if !exist {
			data.Hos = append(data.Hos, st.HosID)
		}
	} else {
		if exist && index != -1 {
			data.Hos = append(data.Hos[:index], data.Hos[index+1:]...)
		}
	}
	if err := WriteBack(constant.Hash_User_Fav, st.UID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	/////添加一个关注量到医院
	go hospital_follow_people(st.HosID, st.IsAttent)
	Forward(session, "0", data)
}

//关注医生
func FollowDocto(session *JsNet.StSession) {
	type st_get struct {
		DocID    string
		UID      string
		IsAttent bool
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" || st.DocID == "" {
		ForwardEx(session, "1", nil, "关注医生失败,参数不完整,UID=%s,DocID=%s\n",
			st.UID, st.DocID)
		return
	}
	data := &SM_Collection{}
	if err := WriteLock(constant.Hash_User_Fav, st.UID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	if data.CreatTime == "" {
		data.CreatTime = CurTime()
	}
	exist := false
	index := -1
	for i, v := range data.Doc {
		if v == st.DocID {
			exist = true
			index = i
			break
		}
	}
	if st.IsAttent {
		if !exist {
			//切片增加
			data.Doc = append(data.Doc, st.DocID)
		}
	} else {
		if exist && index != -1 {
			data.Doc = append(data.Doc[:index], data.Doc[index+1:]...)
		}
	}
	if err := WriteBack(constant.Hash_User_Fav, st.UID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	/////添加一个关注量到医生
	go doctor_follow_people(st.DocID, st.IsAttent)
	Forward(session, "0", data)
}

//收藏产品
func CollectionPorduct(session *JsNet.StSession) {
	type st_get struct {
		ProID    string
		UID      string
		IsAttent bool
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" || st.ProID == "" {
		ForwardEx(session, "1", nil, "收藏产品失败,参数不完整,UID=%s,PorID=%s\n",
			st.UID, st.ProID)
		return
	}
	data := &SM_Collection{}
	if err := WriteLock(constant.Hash_User_Fav, st.UID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	if data.CreatTime == "" {
		data.CreatTime = CurTime()
	}
	exist := false
	index := -1
	for i, v := range data.Por {
		if v == st.ProID {
			exist = true
			index = i
			break
		}
	}
	if st.IsAttent {
		if !exist {
			//切片增加
			data.Por = append(data.Por, st.ProID)
		}
	} else {
		if exist && index != -1 {
			data.Por = append(data.Por[:index], data.Por[index+1:]...)
		}
	}
	if err := WriteBack(constant.Hash_User_Fav, st.UID, data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	////统计产品的关注量
	go AttentionProduct(st.ProID, st.IsAttent)
	Forward(session, "0", data)
}

//关注医院查看
func FollowHosView(session *JsNet.StSession) {
	type st_get struct {
		UID string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" {
		ForwardEx(session, "1", nil, "医院查看,UID=%s\n", st.UID)
		return
	}
	data := &SM_Collection{}
	if err := ShareLock(constant.Hash_User_Fav, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", QueryMoreHosInfo(data.Hos))
}

//关注医生查看
func FollowDocView(session *JsNet.StSession) {
	type st_get struct {
		UID string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" {
		ForwardEx(session, "1", nil, "医生查看,UID=%s\n", st.UID)
		return
	}
	data := &SM_Collection{}
	if err := ShareLock(constant.Hash_User_Fav, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", QueryMoreDocInfo(data.Doc))
}

//收藏产品查看
func CollectionProView(session *JsNet.StSession) {
	type st_get struct {
		UID string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" {
		ForwardEx(session, "1", nil, "产品查看,UID=%s\n", st.UID)
		return
	}
	data := &SM_Collection{}
	if err := ShareLock(constant.Hash_User_Fav, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", QueryMoreProducts(data.Por))
}

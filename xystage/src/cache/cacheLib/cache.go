package cacheLib

import (
	. "JsLib/JsLogger"
)

const (
	SORT_SALE    = "SORT_SALE"    //按销量排序
	SORT_COMMENT = "SORT_COMMENT" //按评论数排序
	SORT_PRICE   = "SORT_PRICE"   //按价格排序
)

type ST_CallPara struct {
	HosID       string   //医院Id
	DocID       string   //医生id
	OrderID     string   //订单id
	ProID       string   //产品id
	City        string   //城市
	FirstItem   string   //一级项目
	SecondItem  string   //二级菜单
	SecondItems []string //二级项目
	ProPrice    int      //产品价格
	ObjectName  string   //对象名称(医院、医生、产品)
	SortType    string   //排序方式
}

type ST_SortPara struct {
	ID  string
	Num int
}

type ST_CallRet struct {
	Ids []string //列表
}

type ST_OrderPara struct {
	HosID      string ///医院id
	OrderID    string //订单id
	SysStatus  string //当前系统订单状态
	UserStatus string //当前用户订单状态
}

type ST_ProPricePara struct {
	HosCity  string //
	ProID    string
	OldPrice int
	NewPrice int
}

type ST_ChangeItem struct {
	ID        string //各种ID
	TimeStamp int64  //时间戳
}

type ST_ChangeParam struct {
	HosID   string
	DocID   string
	ProID   string
	OrderID string
	UID     string
	Keys    []string
	IsRset  bool //是否清空
}

type ST_ChangeRet struct {
	List []ST_ChangeItem
}

type ST_ChangeMapRet struct {
	List map[string][]ST_ChangeItem
}

type ST_SearchHash struct {
	Data map[string]string
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

////本文件,主要记录线上、线下的ID列表
type ST_ONUMCache struct {
	OnLine  []string //在线运行列表
	OffLine []string //下架的列表
	New     []string //新创建的,待审核列表
	UnPass  []string //审核不通过列表
	Modify  []string //修改的列表,待重新审核
	Full    []string //全部列表
}

func (this *ST_ONUMCache) AppendOther(Other *ST_ONUMCache) {
	this.New = append(this.New, Other.New...)
	this.OnLine = append(this.OnLine, Other.OnLine...)
	this.OffLine = append(this.OffLine, Other.OffLine...)
	this.Modify = append(this.Modify, Other.Modify...)
	this.UnPass = append(this.UnPass, Other.UnPass...)
	this.Full = append(this.Full, Other.Full...)
}

func (this *ST_ONUMCache) Clear() {

	this.New = []string{}
	this.Modify = []string{}
	this.OnLine = []string{}
	this.OffLine = []string{}
	this.UnPass = []string{}
	this.Full = []string{}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////
//完整的緩存隊列
type ST_FullCache struct {
	Ids []string //緩存ids
}

//添加到full数据库
func (this *ST_FullCache) AddToFull(ID string) {
	exist := false
	L := len(this.Ids)
	for i := L - 1; i >= 0; i-- {
		if this.Ids[i] == ID {
			exist = true
			break
		}
	}
	if !exist {
		this.Ids = append(this.Ids, ID)
	}
}

func (this *ST_FullCache) RemoveFromFull(ID string) {
	for i := len(this.Ids) - 1; i >= 0; i-- {
		if this.Ids[i] == ID {
			this.Ids = append(this.Ids[:i], this.Ids[i+1:]...)
			break
		}
	}
}

func (this *ST_FullCache) Clear() {
	this.Ids = []string{}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////
type TypListCache struct {
	Data map[string][]string
}

func (this *TypListCache) NewTypeList(Type, ID string) (bool, error) {
	if Type == "" || ID == "" {
		return false, ErrorLog("NewTypeList failed,parame empty,Type=%s,ID=%s\n", Type, ID)
	}
	exist := false
	if this.Data == nil {
		this.Data = make(map[string][]string)
	}
	if list, ok := this.Data[Type]; ok {
		for _, v := range list {
			if v == ID {
				exist = true
				break
			}
		}
		if !exist {
			this.Data[Type] = append(this.Data[Type], ID)
		}
	} else {
		list := []string{ID}
		this.Data[Type] = list
	}

	return exist, nil
}

func (this *TypListCache) RemoveTypeList(Type, ID string) (bool, error) {
	if Type == "" || ID == "" {
		return false, ErrorLog("RemoveTypeList failed,parame empty,Type=%s,ID=%s\n", Type, ID)
	}
	exist := false
	if this.Data == nil {
		this.Data = make(map[string][]string)
		return false, nil
	}
	if list, ok := this.Data[Type]; ok {
		for i, v := range list {
			if v == ID {
				exist = true
				list = append(list[:i], list[i+1:]...)
				this.Data[Type] = list
				break
			}
		}
	} else {
		return false, nil
	}
	return exist, nil
}

func (this *TypListCache) Clear() {
	this.Data = make(map[string][]string)
	return
}

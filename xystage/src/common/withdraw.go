package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	. "cache/cacheLib"
	"constant"
	"ider"
	. "util"
)

type ST_WithDrawRecord struct {
	WDID        string //提现ID
	UID         string //uid
	Name        string //姓名
	Money       int    //金额
	BankNumber  string // 银行卡号
	BankName    string // 开户行
	Tax         int    //税额
	Balance     int    //提现后余额
	Status      int    //状态 0:申请提现  1:已经打款 2:提现有疑问
	Des         string //提现描述
	ConfirmDate string //确认或者打回时间
	Date        string //提现时间
}

//新建一条提现记录
func NewWithDraw(user *ST_User, Money, Balance, status int) (*ST_WithDrawRecord, error) {
	ErrorLog("000000000000000000000\n")
	if user == nil || user.Agent == nil {
		return nil, ErrorLog("NewWithDraw failed,Agent is empty\n")
	}
	if Money < 100 {
		return nil, ErrorLog("NewWithDraw failed,Money < 100\n")
	}
	id, err := ider.GenwWithDrwaID()
	if err != nil {
		return nil, err
	}
	st := &ST_WithDrawRecord{
		WDID:       id,
		UID:        user.UID,
		Name:       user.Name,
		Money:      Money,
		Status:     status,
		Balance:    Balance,
		BankName:   user.Agent.BankName,
		BankNumber: user.Agent.BankNumber,
		Date:       CurTime(),
	}
	e := DirectWrite(constant.Hash_WithDraw, id, st)
	if e == nil {
		go AppendGlobalWithDraw(id)
	}
	return st, e
}

//获取提现信息
func GetWithDrawInfo(WDID string) (*ST_WithDrawRecord, error) {
	data := &ST_WithDrawRecord{}
	err := ShareLock(constant.Hash_WithDraw, WDID, data)
	return data, err
}

func QueryMoreRecord(session *JsNet.StSession) {
	type st_info struct {
		WDIDs []string
	}
	st := &st_info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data := []*ST_WithDrawRecord{}
	if len(st.WDIDs) == 0 {
		ForwardEx(session, "0", data, "QueryWithDrawInfo failed,WDID is empty\n")
		return
	}

	for _, v := range st.WDIDs {
		dd, err := GetWithDrawInfo(v)
		if err != nil {
			continue
		}
		data = append(data, dd)
	}

	Forward(session, "0", data)
}

//获取提现信息
func QueryWithDrawInfo(session *JsNet.StSession) {
	type st_info struct {
		WDID string //提现id
	}
	st := &st_info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.WDID == "" {
		ForwardEx(session, "1", nil, "QueryWithDrawInfo failed,WDID is empty\n")
		return
	}
	data, err := GetWithDrawInfo(st.WDID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Forward(session, "0", data)
}

//确认提现打款
func confirmWithDraw(WDID, Des string, IsPass bool) error {
	data := &ST_WithDrawRecord{}
	if err := WriteLock(constant.Hash_WithDraw, WDID, data); err != nil {
		return err
	}
	defer WriteBack(constant.Hash_WithDraw, WDID, data)
	if !IsPass {
		_, IsE := WithdrawBalanceFailed(data.UID, data.Money)
		if IsE != nil {
			return IsE
		}
		data.Status = 2
	} else {
		data.Status = 1
	}
	data.Des = Des
	data.ConfirmDate = CurTime()
	return nil
}

//确认提现打款
func ComfirmWithDraw(session *JsNet.StSession) {
	type st_get struct {
		WDID   string //提现ID
		IsPass bool   //是否同意
		Des    string //描述
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.WDID == "" {
		ForwardEx(session, "1", nil, "ComfirmWithDrew param is empty\n")
		return
	}
	if err := confirmWithDraw(st.WDID, st.Des, st.IsPass); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", nil)
}

//添加一条全局的提现记录
func AppendGlobalWithDraw(WDID string) error {
	data := &TypListCache{}
	err := WriteLock(constant.Hash_WithDraw, constant.KEY_Global_WithDraw, data)
	data.NewTypeList(CurYearMonth(), WDID)
	if err != nil {
		return DirectWrite(constant.Hash_WithDraw, constant.KEY_Global_WithDraw, data)
	}
	return WriteBack(constant.Hash_WithDraw, constant.KEY_Global_WithDraw, data)
}

//获取全局的提现记录ID
func GetGlobalWithDrawID(Month string) ([]string, []string, []string, error) {
	data := &TypListCache{}
	if err := ShareLock(constant.Hash_WithDraw, constant.KEY_Global_WithDraw, data); err != nil {
		return []string{}, []string{}, []string{}, err
	}
	Apply := []string{}
	Success := []string{}
	UnPass := []string{}

	if list, ok := data.Data[Month]; ok {
		for _, v := range list {
			dd := &ST_WithDrawRecord{}
			if err := ShareLock(constant.Hash_WithDraw, v, dd); err != nil {
				continue
			}
			if dd.Status == 0 {
				Apply = append(Apply, v)
			}
			if dd.Status == 1 {
				Success = append(Success, v)
			}
			if dd.Status == 2 {
				UnPass = append(UnPass, v)
			}
		}
	}
	return Apply, Success, UnPass, nil
}

func GetWithDrawMonth(session *JsNet.StSession) {
	data := &TypListCache{}
	if err := ShareLock(constant.Hash_WithDraw, constant.KEY_Global_WithDraw, data); err != nil {
		ForwardEx(session, "1", []string{}, err.Error())
		return
	}
	list := []string{}
	if data != nil && data.Data != nil {
		for k, _ := range data.Data {
			list = append(list, k)
		}
	}
	if len(list) > 0 {
		list = append(list, "All")
	}
	Forward(session, "0", list)
}

//获取全局的提现记录信息
func GetGlobalWithDrawRecord(session *JsNet.StSession) {
	type st_Get struct {
		Month string
	}
	para := &st_Get{}
	if para.Month == "" {
		para.Month = "All"
	}
	data := &TypListCache{}
	if err := ShareLock(constant.Hash_WithDraw, constant.KEY_Global_WithDraw, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	type st_data struct {
		Apply   []*ST_WithDrawRecord
		Success []*ST_WithDrawRecord
		UnPass  []*ST_WithDrawRecord
	}
	st := &st_data{}
	list := []string{}
	if para.Month == "全部" || para.Month == "All" {
		if data != nil && data.Data != nil {
			for _, v := range data.Data {
				list = append(list, v...)
			}
		}
	} else {
		list, _ = data.Data[para.Month]
	}

	for _, v := range list {
		dd := &ST_WithDrawRecord{}
		if err := ShareLock(constant.Hash_WithDraw, v, dd); err != nil {
			continue
		}
		if dd.Status == 0 {
			st.Apply = append(st.Apply, dd)
		}
		if dd.Status == 1 {
			st.Success = append(st.Success, dd)
		}
		if dd.Status == 2 {
			st.UnPass = append(st.UnPass, dd)
		}
	}
	Forward(session, "0", st)
}

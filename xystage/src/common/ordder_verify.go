package common

import (
	. "JsLib/JsLogger"
	"constant"
	"math/rand"
	"time"
	. "util"
)

///订单校验映射表
type ST_OrdderVerify struct {
	HosID      string //医院id
	OrderID    string //订单id
	UID        string //用户id
	Status     string //状态
	VerifyDate string //校验时间
	CreatDate  string //创建时间
}

///创建一个校验码
func NewVerifyCode(OrderID, HosID, UID string) (string, error) {
	st := &ST_OrdderVerify{
		HosID:     HosID,
		OrderID:   OrderID,
		UID:       UID,
		Status:    "0",
		CreatDate: CurTime(),
	}
	if st.OrderID == "" {
		return "", ErrorLog("生成订单校验码失败,OrderID=%s\n", OrderID)
	}
	code := genVerifCode(6)
	if err := DirectWrite(constant.Hash_Order_Verify, code, st); err != nil {
		return "", ErrorLog("生成订单校验码失败DirectWrite(),OrderID=%s\n", OrderID)
	}
	return code, nil
}

//校验码验证,并返回订单id
func CheckVerifyCode(code, HosID string) (string, error) {
	if code == "" || HosID == "" {
		return "", ErrorLog("检验码验证失败,VerifyCode为空\n")
	}
	st := &ST_OrdderVerify{}
	if err := WriteLock(constant.Hash_Order_Verify, code, st); err != nil {
		return "", err
	}
	defer WriteBack(constant.Hash_Order_Verify, code, st)
	st.Status = "1"
	st.VerifyDate = CurTime()
	if st.HosID != HosID {
		return "", ErrorLog("检验码验证失败,医院id不匹配,code=%s,HosID=%s\n", code, HosID)
	}
	return st.OrderID, nil
}

///删除校验码映射表
func DelVerifyCode(code string) error {
	return HDel(constant.Hash_Order_Verify, code)
}

//生成随机字符串
func genVerifCode(length int64) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	var i int64 = 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for ; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

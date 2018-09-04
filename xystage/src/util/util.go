package util

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"errors"
	"fmt"
	"time"
)

//错误返回,并且带有错误日志输出
func ForwardEx(session *JsNet.StSession, ret interface{}, entry interface{}, info string, a ...interface{}) {
	res := make(map[string]interface{})
	res[constant.CT_Ret] = ret
	res[constant.CT_Entity] = entry
	if a == nil {
		Error(info)
		res[constant.CT_Msg] = fmt.Sprintf(info)
	} else {
		Error(info, a...)
		res[constant.CT_Msg] = fmt.Sprintf(info, a...)
	}
	session.Forward(res)
	return
}

//填充返回值
func Forward(session *JsNet.StSession, ret interface{}, entry interface{}) {
	res := make(map[string]interface{})
	res[constant.CT_Ret] = ret
	res[constant.CT_Entity] = entry
	res[constant.CT_Msg] = "Success"
	session.Forward(res)
	return
}

func Err(info string, a ...interface{}) error {
	i := fmt.Sprintf(info, a...)
	return errors.New(i)
}

//返回当前时间：例如 2017-02-17 16:33
func CurTime() string {
	return time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
}

//返回当前日期：例如：2017-02-17
func CurDate() string {
	return time.Unix(time.Now().Unix(), 0).Format("2006-01-02")
}

///当前时间的时间戳
func CurStamp() int64 {
	return time.Now().Unix()
}

//返回当前几号
func CurDay() int {
	return time.Now().Day()
}

//当前年份
func CurYear() string {
	return time.Unix(time.Now().Unix(), 0).Format("2006")
}

//当前的月
func CurMonth() string {
	return time.Unix(time.Now().Unix(), 0).Format("01")
}

//当前的月
func CurYearMonth() string {
	return time.Unix(time.Now().Unix(), 0).Format("2006-01")
}

func GetYearMOnth(time time.Time) string {
	return time.Format("2006-01")
}

func GetTimeFormString(date string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", date, time.Local)
}

//获取相差时间
func GetHourDiffer(start_time, end_time string) int64 {
	var hour int64
	t1, err1 := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	if err1 != nil {
		return hour
	}
	t2, err2 := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	if err2 == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		hour = diff / 3600
		return hour
	} else {
		return hour
	}
}

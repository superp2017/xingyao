package util

import (
	"JsLib/JsDBCache"
	. "JsLib/JsLogger"
	"constant"
	"encoding/json"
	"errors"
)

func Get(k string, v interface{}) error {

	dtio := JsDBCache.GenDtio(constant.C_XINGYAODB)
	if dtio == nil {
		return errors.New("GenDtio error")
	}

	b, e := dtio.Dtio_Get(k)
	if e != nil {
		Error(e.Error())
		return e
	}

	e = json.Unmarshal(b, v)
	if e != nil {
		Error(e.Error())
		return e
	}

	return nil
}

func Set(k string, v interface{}) error {
	dtio := JsDBCache.GenDtio(constant.C_XINGYAODB)
	if dtio == nil {
		return errors.New("GenDtio error")
	}

	b, e := json.Marshal(v)
	if e != nil {
		Error(e.Error())
		return e
	}

	e = dtio.Dtio_Set(k, b)
	if e != nil {
		Error(e.Error())
		return e
	}
	return nil
}

func ShareLock(t string, k string, v interface{}) error {
	dtio := JsDBCache.GenDtio(constant.C_XINGYAODB)
	if dtio == nil {
		return errors.New("GenDtio error")
	}

	_, e := dtio.Dtio_ShareLock(t, k, v)
	if e != nil {
		Error(e.Error())
		return e
	}

	return nil
}

func WriteLock(t string, k string, v interface{}) error {
	dtio := JsDBCache.GenDtio(constant.C_XINGYAODB)
	if dtio == nil {
		return errors.New("GenDtio error")
	}

	e := dtio.Dtio_WriteLock(t, k, v)
	if e != nil {
		Error(e.Error())
		return e
	}

	return nil
}

func WriteBack(t string, k string, v interface{}) error {
	dtio := JsDBCache.GenDtio(constant.C_XINGYAODB)
	if dtio == nil {
		return errors.New("GenDtio error")
	}

	e := dtio.Dtio_WriteBack(t, k, v)
	if e != nil {
		Error(e.Error())
		return e
	}
	return nil
}

//
//func Modify(t string, k string, v interface{}, m interface{}) error {
//	dtio := JsDBCache.GenDtio(constant.C_XINGYAODB)
//	if dtio == nil {
//		return errors.New("GenDtio error")
//	}
//	if e := dtio.Dtio_WriteLock(t, k, v); e != nil {
//		Error(e.Error())
//		return e
//	}
//	data, err := json.Marshal(m)
//	if err != nil {
//		Error("json.Marshal(m) error: %s", err.Error())
//		return err
//	}
//
//	if err = json.Unmarshal(data, &v); err != nil {
//		Error("json.Unmarshal(data, &v) error: %s", err.Error())
//		return err
//	}
//	if err := dtio.Dtio_WriteBack(t, k, v); err != nil {
//		Error(err.Error())
//		return err
//	}
//	return nil
//}
//
//func ModifyEx(t string, k string, v interface{}, m interface{}, f func()) error {
//	dtio := JsDBCache.GenDtio(constant.C_XINGYAODB)
//	if dtio == nil {
//		return errors.New("GenDtio error")
//	}
//	if e := dtio.Dtio_WriteLock(t, k, v); e != nil {
//		Error(e.Error())
//		return e
//	}
//	data, err := json.Marshal(m)
//	if err != nil {
//		Error("json.Marshal(m) error: %s", err.Error())
//		return err
//	}
//
//	if err = json.Unmarshal(data, &v); err != nil {
//		Error("json.Unmarshal(data, &v) error: %s", err.Error())
//		return err
//	}
//	f()
//	if err := dtio.Dtio_WriteBack(t, k, v); err != nil {
//		Error(err.Error())
//		return err
//	}
//	return nil
//}
//
//func Update(t string, k string, v interface{}, m func()) error {
//	dtio := JsDBCache.GenDtio(constant.C_XINGYAODB)
//	if dtio == nil {
//		return errors.New("GenDtio error")
//	}
//
//	if e := dtio.Dtio_Update(t, k, v, m); e != nil {
//		Error(e.Error())
//		return e
//	}
//	return nil
//}
//
//func UpdateEx(t string, k string, v interface{}, m func(error)) error {
//	dtio := JsDBCache.GenDtio(constant.C_XINGYAODB)
//	if dtio == nil {
//		return errors.New("GenDtio error")
//	}
//
//	if e := dtio.Dtio_UpdateEx(t, k, v, m); e != nil {
//		Error(e.Error())
//		return e
//	}
//	return nil
//}

func DirectWrite(t string, k string, v interface{}) error {
	dtio := JsDBCache.GenDtio(constant.C_XINGYAODB)
	if dtio == nil {
		return errors.New("GenDtio error")
	}

	if e := dtio.Dtio_UnsafeWriteEx(t, k, v); e != nil {
		Error(e.Error())
		return e
	}
	return nil
}

func HDel(table, id string) error {
	dtio := JsDBCache.GenDtio(constant.C_XINGYAODB)
	if dtio == nil {
		return errors.New("GenDtio error")
	}
	if e := dtio.Dtio_HDel(table, id); e != nil {
		Error(e.Error())
		return e
	}
	return nil
}

func Del(table string) error {
	dtio := JsDBCache.GenDtio(constant.C_XINGYAODB)
	if dtio == nil {
		return errors.New("GenDtio error")
	}
	if er := dtio.Dtio_Del(table); er != nil {
		Error(er.Error())
		return er
	}
	return nil
}

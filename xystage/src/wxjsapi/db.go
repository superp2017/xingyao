package main

import (
	. "JsLib/JsConfig"
	"JsLib/JsDBCache"
	. "JsLib/JsLogger"

	"encoding/json"
	"errors"
)

var C_XINGYAODB string = CFG.Http.DbName

func Get(k string, v interface{}) error {

	dtio := JsDBCache.GenDtio(C_XINGYAODB)
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
	dtio := JsDBCache.GenDtio(C_XINGYAODB)
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
	dtio := JsDBCache.GenDtio(C_XINGYAODB)
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
	dtio := JsDBCache.GenDtio(C_XINGYAODB)
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
	dtio := JsDBCache.GenDtio(C_XINGYAODB)
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

func Update(t string, k string, v interface{}, m func()) error {
	dtio := JsDBCache.GenDtio(C_XINGYAODB)
	if dtio == nil {
		return errors.New("GenDtio error")
	}

	e := dtio.Dtio_Update(t, k, v, m)
	if e != nil {
		Error(e.Error())
		return e
	}
	return nil
}

func UpdateEx(t string, k string, v interface{}, m func(error)) error {
	dtio := JsDBCache.GenDtio(C_XINGYAODB)
	if dtio == nil {
		return errors.New("GenDtio error")
	}

	e := dtio.Dtio_UpdateEx(t, k, v, m)
	if e != nil {
		Error(e.Error())
		return e
	}
	return nil
}

func DirectWrite(t string, k string, v interface{}) error {
	dtio := JsDBCache.GenDtio(C_XINGYAODB)
	if dtio == nil {
		return errors.New("GenDtio error")
	}

	e := dtio.Dtio_UnsafeWriteEx(t, k, v)
	if e != nil {
		Error(e.Error())
		return e
	}
	return nil
}

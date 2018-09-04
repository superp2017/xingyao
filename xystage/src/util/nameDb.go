package util

import (
	"JsLib/JsDBCache"
	. "JsLib/JsLogger"
	"encoding/json"
	"errors"
)

type NameDB struct {
	Name string
}

func NewDb(name string) *NameDB {
	return &NameDB{Name: name}
}

func (d *NameDB) Get(k string, v interface{}) error {

	dtio := JsDBCache.GenDtio(d.Name)
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

func (d *NameDB) Set(k string, v interface{}) error {
	dtio := JsDBCache.GenDtio(d.Name)
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

func (d *NameDB) ShareLock(t string, k string, v interface{}) error {
	dtio := JsDBCache.GenDtio(d.Name)
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

func (d *NameDB) WriteLock(t string, k string, v interface{}) error {
	dtio := JsDBCache.GenDtio(d.Name)
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

func (d *NameDB) WriteBack(t string, k string, v interface{}) error {
	dtio := JsDBCache.GenDtio(d.Name)
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

func (d *NameDB) Modify(t string, k string, v interface{}, m interface{}) error {
	dtio := JsDBCache.GenDtio(d.Name)
	if dtio == nil {
		return errors.New("GenDtio error")
	}
	if e := dtio.Dtio_WriteLock(t, k, v); e != nil {
		Error(e.Error())
		return e
	}
	data, err := json.Marshal(m)
	if err != nil {
		Error("json.Marshal(m) error: %s", err.Error())
		return err
	}

	if err = json.Unmarshal(data, &v); err != nil {
		Error("json.Unmarshal(data, &v) error: %s", err.Error())
		return err
	}
	if err := dtio.Dtio_WriteBack(t, k, v); err != nil {
		Error(err.Error())
		return err
	}
	return nil
}

func (d *NameDB) ModifyEx(t string, k string, v interface{}, m interface{}, f func()) error {
	dtio := JsDBCache.GenDtio(d.Name)
	if dtio == nil {
		return errors.New("GenDtio error")
	}
	if e := dtio.Dtio_WriteLock(t, k, v); e != nil {
		Error(e.Error())
		return e
	}
	data, err := json.Marshal(m)
	if err != nil {
		Error("json.Marshal(m) error: %s", err.Error())
		return err
	}

	if err = json.Unmarshal(data, &v); err != nil {
		Error("json.Unmarshal(data, &v) error: %s", err.Error())
		return err
	}
	f()
	if err := dtio.Dtio_WriteBack(t, k, v); err != nil {
		Error(err.Error())
		return err
	}
	return nil
}

func (d *NameDB) Update(t string, k string, v interface{}, m func()) error {
	dtio := JsDBCache.GenDtio(d.Name)
	if dtio == nil {
		return errors.New("GenDtio error")
	}

	if e := dtio.Dtio_Update(t, k, v, m); e != nil {
		Error(e.Error())
		return e
	}
	return nil
}

func (d *NameDB) UpdateEx(t string, k string, v interface{}, m func(error)) error {
	dtio := JsDBCache.GenDtio(d.Name)
	if dtio == nil {
		return errors.New("GenDtio error")
	}

	if e := dtio.Dtio_UpdateEx(t, k, v, m); e != nil {
		Error(e.Error())
		return e
	}
	return nil
}

func (d *NameDB) DirectWrite(t string, k string, v interface{}) error {
	dtio := JsDBCache.GenDtio(d.Name)
	if dtio == nil {
		return errors.New("GenDtio error")
	}

	if e := dtio.Dtio_UnsafeWriteEx(t, k, v); e != nil {
		Error(e.Error())
		return e
	}
	return nil
}

func (d *NameDB) HDel(table, id string) error {
	dtio := JsDBCache.GenDtio(d.Name)
	if dtio == nil {
		return errors.New("GenDtio error")
	}
	if e := dtio.Dtio_HDel(table, id); e != nil {
		Error(e.Error())
		return e
	}
	return nil
}

func (d *NameDB) Del(table string) error {
	dtio := JsDBCache.GenDtio(d.Name)
	if dtio == nil {
		return errors.New("GenDtio error")
	}
	if er := dtio.Dtio_Del(table); er != nil {
		Error(er.Error())
		return er
	}
	return nil
}

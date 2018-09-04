package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	. "cache/cacheLib"
	"constant"
	. "util"
)

type ST_HomeShow struct {
	FirstDoc  string
	SecondDoc string
	ThirdDoc  string
	FourthDoc string
}

type ST_HomeExpertDoc struct {
	Data map[string]*ST_HomeShow
}

///获取大牌名医
func GetExpertDoctor(session *JsNet.StSession) {
	type st_get struct {
		City string //城市
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	city := st.City
	if city == "全部" || city == "全国" {
		city = constant.All
	}
	data, err := getExpertDoc(city)
	if err == nil {
		Forward(session, "0", QueryMoreDocInfo(data))
		return
	}
	ForwardEx(session, "1", nil, err.Error())
}

///获取 首页大牌名医
func GetHomeExpertDoctor(session *JsNet.StSession) {
	type st_get struct {
		Pos   int
		DocID string
		City  string
	}
	st := &st_get{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.City == "" {
		ForwardEx(session, "1", nil, "GetHomeExpertDoctor failed..City=%s\n", st.City)
		return
	}
	data, err := getHomeExpertDoc()
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	city := st.City
	if city == "全部" || city == "全国" {
		city = constant.All
	}

	if d, ok := data.Data[city]; ok {
		ls := []*ST_Doctor{}
		if d != nil {
			if doc1, err := QueryDoctor(d.FirstDoc); err == nil {
				ls = append(ls, doc1)
			} else {
				ls = append(ls, &ST_Doctor{})
			}
			if doc2, err := QueryDoctor(d.SecondDoc); err == nil {
				ls = append(ls, doc2)
			} else {
				ls = append(ls, &ST_Doctor{})
			}
			if doc3, err := QueryDoctor(d.ThirdDoc); err == nil {
				ls = append(ls, doc3)
			} else {
				ls = append(ls, &ST_Doctor{})
			}
			if doc4, err := QueryDoctor(d.FourthDoc); err == nil {
				ls = append(ls, doc4)
			} else {
				ls = append(ls, &ST_Doctor{})
			}
		}
		Forward(session, "0", ls)
		return
	}
	ForwardEx(session, "0", nil, "该城市不存在首页大牌名医,City=%s", st.City)
}

///设置首页大牌名医
func SetHomeExpertDoctor(session *JsNet.StSession) {
	type st_get struct {
		Pos   int
		DocID string
		City  string
		Pic   string //大牌名医的背景照
	}
	st := &st_get{}

	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.City == "" || st.DocID == "" || st.Pic == "" || st.Pos < 0 {
		ForwardEx(session, "1", nil, "SetHomeExpertDoctor failed..City=%s,DocID=%s,Pos=%d,Pic=%s\n",
			st.City, st.DocID, st.Pos, st.Pic)
		return
	}
	city := st.City
	if city == "全部" || city == "全国" {
		city = constant.All
	}

	if err := setHomeExpertDoc(st.DocID, city, st.Pic, st.Pos); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", nil)
}

///获取首页的大牌名医
func getHomeExpertDoc() (*ST_HomeExpertDoc, error) {
	st := &ST_HomeExpertDoc{}
	if err := ShareLock(constant.Hash_Doctor, constant.KV_HomeExpertDoctor, st); err != nil {
		return nil, err
	}
	return st, nil
}

func setHomeExpertDoc(DocID, City, Pic string, Pos int) error {
	st := &ST_HomeExpertDoc{}
	err := WriteLock(constant.Hash_Doctor, constant.KV_HomeExpertDoctor, st)
	defer WriteBack(constant.Hash_Doctor, constant.KV_HomeExpertDoctor, st)

	show := &ST_HomeShow{}
	if Pos == 0 {
		show.FirstDoc = DocID
	}
	if Pos == 1 {
		show.SecondDoc = DocID
	}
	if Pos == 2 {
		show.ThirdDoc = DocID
	}
	if Pos == 3 {
		show.FourthDoc = DocID
	}
	if err != nil {
		st.Data = make(map[string]*ST_HomeShow)
		st.Data[City] = show
		if err1 := DirectWrite(constant.Hash_Doctor, constant.KV_HomeExpertDoctor, st); err1 != nil {
			return err1
		}
	} else {
		st.Data[City] = show
	}
	return changeDocExpertPhoto(DocID, Pic)
}

func removeHomeExpertDoctor(City, DocID string) error {
	st := &ST_HomeExpertDoc{}
	if err := WriteLock(constant.Hash_Doctor, constant.KV_HomeExpertDoctor, st); err != nil {
		return err
	}
	if v, ok := st.Data[City]; ok {
		if v.FirstDoc == DocID {
			v.FirstDoc = ""
		}
		if v.SecondDoc == DocID {
			v.SecondDoc = ""
		}
		if v.ThirdDoc == DocID {
			v.ThirdDoc = ""
		}
		if v.FourthDoc == DocID {
			v.FourthDoc = ""
		}
		st.Data[City] = v
	}
	return WriteBack(constant.Hash_Doctor, constant.KV_HomeExpertDoctor, st)

}

////添加一个大牌名医
func addExpertDoctor(City, DocID string) error {
	data := &TypListCache{}
	if err := WriteLock(constant.Hash_Doctor, constant.KV_HomeExpertDoctor, data); err != nil {
		return err
	}
	defer WriteBack(constant.Hash_Doctor, constant.KV_HomeExpertDoctor, data)
	if _, e := data.NewTypeList(City, DocID); e != nil {
		return e
	}
	if _, e := data.NewTypeList(constant.All, DocID); e != nil {
		return e
	}
	return nil
}

////从大牌名医列表中移除
func removeExpertDoctor(DocID, City string) error {
	data := &TypListCache{}
	if err := WriteLock(constant.Hash_Doctor, constant.KV_HomeExpertDoctor, data); err != nil {
		return err
	}
	defer WriteBack(constant.Hash_Doctor, constant.KV_HomeExpertDoctor, data)
	if _, e := data.RemoveTypeList(City, DocID); e != nil {
		return e
	}
	if _, e := data.RemoveTypeList(constant.All, DocID); e != nil {
		return e
	}
	return nil
}

///获取大排名医
func getExpertDoc(City string) ([]string, error) {
	if City == "" {
		return []string{}, ErrorLog("getExpertDocTor failed,parame is empty,City=%s\n", City)
	}
	data := &TypListCache{}
	if err := ShareLock(constant.Hash_Doctor, constant.KV_ExpertDoctor, data); err != nil {
		return []string{}, err
	}
	key := City
	if key == "全部" || key == "全国" {
		key = constant.All
	}
	if ls, ok := data.Data[key]; ok {
		return ls, nil
	}
	return []string{}, nil
}

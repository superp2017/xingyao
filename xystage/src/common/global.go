package common

import (
	"JsLib/JsNet"
	"constant"
	. "util"
)

//全局城市map
type ST_CityMap struct {
	Initial string
	City    []string
}

///获取全局城市映射表
func GetCityMap(session *JsNet.StSession) {

	data := &[]ST_CityMap{}
	if err := Get(constant.KEY_GlobalCityMap, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Forward(session, "0", data)
}

///设置城市map
func SetCityMap(session *JsNet.StSession) {
	type st_set struct {
		Initial string
		City    string
	}
	st := &st_set{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.City == "" || st.Initial == "" {
		ForwardEx(session, "1", nil, "SetCityMap failed ,param is empty,City=%s,Initial=%s\n", st.City, st.Initial)
		return
	}
	data := []ST_CityMap{}
	if err := Get(constant.KEY_GlobalCityMap, &data); err != nil {
		ForwardEx(session, "1", nil, "err:=%s,getCitMap failed ,param is empty,City=%s,Initial=%s\n", err.Error(), st.City, st.Initial)
		return
	}

	for i, v := range data {
		if v.Initial == st.Initial {
			ex := false
			for _, v1 := range v.City {
				if v1 == st.City {
					ex = true
					break
				}
			}
			if !ex {
				data[i].City = append(data[i].City, st.City)
			}
		}
	}

	if err := Set(constant.KEY_GlobalCityMap, &data); err != nil {
		ForwardEx(session, "1", nil, "err:=%s,Set failed ,param is empty,City=%s,Initial=%s\n", err.Error(), st.City, st.Initial)
		return
	}

	Forward(session, "0", data)
}

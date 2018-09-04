package main

import (
	. "JsLib/JsLogger"
	"cache/cacheIO"
	"constant"
	"strconv"
)

//获取产品
func getProducts(City, FirstItem, SecondItem string, min, max int) ([]string, error) {

	city := City
	if city == "全部" || city == "全国" {
		city = constant.All
	}

	////////////跟金钱无关/////////////////////
	if min <= 0 && max <= 0 {

		return getIems(city, FirstItem, SecondItem)
	}

	////////////跟金钱相关/////////////////////
	if min >= 0 && max >= min {

		data, er := getPrice(min, max)
		if City == "" && FirstItem == "" && SecondItem == "" {
			return data, er
		}

		if City != "" && FirstItem == "" && SecondItem == "" {
			return getCityPrice(City, min, max)
		}

		//////////////联合查询///////////////////
		res := []string{}
		ls, ls_err := getIems(city, FirstItem, SecondItem)
		if ls_err != nil {
			return []string{}, ls_err
		}

		for _, v1 := range data {
			for _, v2 := range ls {
				if v1 == v2 {
					res = append(res, v1)
				}
			}
		}
		return res, nil
	}
	return []string{}, nil
}

///根据城市、身体部位、二级菜单获取产品
func getIems(City, FirstItem, SecondItem string) ([]string, error) {

	if City == "" && FirstItem == "" && SecondItem == "" {

		return []string{}, nil
	}

	if City != "" && FirstItem == "" && SecondItem == "" {

		return cacheIO.GetGlobalCityPro(City)
	}

	if City == "" && FirstItem != "" && SecondItem == "" {

		return cacheIO.GetFirstSecondItemPro(FirstItem)
	}

	if City != "" && FirstItem != "" && SecondItem == "" {

		data1, er1 := cacheIO.GetGlobalCityPro(City)
		if er1 != nil {

			return []string{}, er1
		}
		data2, er2 := cacheIO.GetFirstSecondItemPro(FirstItem)
		if er2 != nil {

			return []string{}, er2
		}
		res := []string{}
		for _, v1 := range data1 {
			for _, v2 := range data2 {
				if v1 == v2 {
					res = append(res, v1)
				}
			}
		}

		return res, nil
	}

	if City == "" && FirstItem != "" && SecondItem != "" {

		return cacheIO.GetItemProduct(FirstItem, SecondItem)
	}

	if City != "" || FirstItem != "" && SecondItem != "" {
		d1, e1 := cacheIO.GetItemProduct(FirstItem, SecondItem)
		if e1 != nil {

			return d1, e1
		}
		d2, e2 := cacheIO.GetGlobalCityPro(City)
		if e2 != nil {
			return d2, e2
		}
		res := []string{}
		for _, v1 := range d1 {
			for _, v2 := range d2 {
				if v1 == v2 {
					res = append(res, v1)
				}
			}
		}

		return res, nil
	}

	return []string{}, nil
}

///根据价格区间获取产品id
func getPrice(min, max int) ([]string, error) {
	//////获取分价格的缓存
	priceMap, err := cacheIO.GetPriceProduct()
	if err != nil {
		return []string{}, err
	}
	list, err := getPriceRange(min, max)
	if err != nil {
		return []string{}, err
	}
	data := []string{}
	for _, v := range list {
		if l, ok := priceMap.Data[v]; ok {
			data = append(data, l...)
		}
	}
	return data, nil
}

///根据城市、价格区间获取产品id
func getCityPrice(City string, min, max int) ([]string, error) {
	priceMap, err := cacheIO.GetCityPriceProduct(City)
	if err != nil {
		return []string{}, err
	}
	list, err := getPriceRange(min, max)
	if err != nil {
		return []string{}, err
	}
	data := []string{}
	for _, v := range list {
		if l, ok := priceMap.Data[v]; ok {
			data = append(data, l...)
		}
	}
	return data, nil
}

//根据价格大小值获取系统的产品价格列表
func getPriceRange(min, max int) ([]string, error) {
	list := []string{}
	if min == max {
		list = append(list, strconv.Itoa(min))
	} else {
		////获取全局的价格列表
		pricelist, err := cacheIO.GetGlobalPrice()
		if err != nil {
			return []string{}, err
		}
		for _, v := range pricelist {
			t, err := strconv.Atoi(v)
			if err == nil {
				if t >= min && t <= max {
					list = append(list, v)
				}
			} else {
				ErrorLog("err:%s,price=%s\n", err.Error(), v)
			}
		}
	}
	return list, nil
}

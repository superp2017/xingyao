package common

import (
	. "JsLib/JsLogger"
	. "cache/cacheIO"
	. "cache/cacheLib"
	"constant"
	"errors"
	"sort"
)

//产品详情
type ST_RequestPar struct {
	ProductType string //1:new;2:modified;3 verifypass;4 verifypass;5 unsale
	//1: for xingyao product searching
	//2: for customer side of bodypart
	//3: for customer side of collect baby
	//4: for customer side of footmark
	ProductTypeSub int
	UserID         string
	RequestPage    int
	SortType       string
}

type ST_TotalNum struct {
	TotalPage_new        int
	TotalPage_modified   int
	TotalPage_verifypass int
	TotalPage_verifyfail int
	TotalPage_unsale     int
}

func GetDedicateListID(dbkey string, st *ST_RequestPar, itemPerPage int) (lsID []string, e error) {

	listID := []string{}
	listPageID := []string{}
	var resList *ST_ONUMCache
	resList = &ST_ONUMCache{}
	var isE error = nil
	//Get the  list id according to the type
	if st.RequestPage < 1 {
		return nil, errors.New("The request pagenumber should bigger than 1\n")
	}

	switch dbkey {
	case constant.Hash_HospitalCache:
		resList, isE = GetGlobalHosPtr(st.SortType)
		break

	case constant.Hash_DoctorCache:
		Info("Get the Doc ID..........\n")
		resList, isE = GetGlobalDocPtr(st.SortType)
		break

	case constant.Hash_ProductCache:
		Info("Get the product ID..........\n")
		resList, isE = GetGlobalProPtr(st.SortType)
		break
	default:
		return nil, errors.New("There is no such an key")
	}
	if isE != nil {
		return listID, isE
	}

	switch st.ProductType {
	case constant.OperatingStatus_modify:
		listID = resList.Modify
		break

	case constant.OperatingStatus_new:
		listID = resList.New
		break

	case constant.OperatingStatus_online:
		listID = resList.OnLine
		break

	case constant.OperatingStatus_Offline_self:
		listID = resList.OffLine
		break

	case constant.OperatingStatus_Reviewer_NotPass:
		listID = resList.UnPass
		break

	default:
		return nil, errors.New("The Type is not correct")
	}

	if dbkey == constant.Hash_DoctorCache || dbkey == constant.Hash_ProductCache {
		sort.Sort(sort.Reverse(sort.StringSlice(listID)))
	}

	listStartDex := (st.RequestPage - 1) * itemPerPage
	if len(listID)-listStartDex <= itemPerPage {
		listPageID = listID[listStartDex:]
		return listPageID, nil
	} else {
		listPageID = listID[listStartDex : listStartDex+itemPerPage]
	}
	return listPageID, nil
}

func GetTotalNum(sortType, dbkey string, totalNum *ST_TotalNum, numberOnePage int) (e error) {

	var resList *ST_ONUMCache
	var isE error = nil
	resList = &ST_ONUMCache{}

	switch dbkey {
	case constant.KEY_HosORMCache:
		resList, isE = GetGlobalHosPtr(sortType)
		break

	case constant.KEY_DocORMCache:

		resList, isE = GetGlobalDocPtr(sortType)
		break

	case constant.KEY_ProORMCache:
		resList, isE = GetGlobalProPtr(sortType)
		break
	default:
		return ErrorLog("There is no such an key")
	}

	if isE != nil {
		return isE
	}
	if resList == nil {
		return ErrorLog("get global ptr failed\n")
	}
	totalNum.TotalPage_modified = getCeilNum(len(resList.Modify), numberOnePage)
	totalNum.TotalPage_new = getCeilNum(len(resList.New), numberOnePage)
	totalNum.TotalPage_unsale = getCeilNum(len(resList.OffLine), numberOnePage)
	totalNum.TotalPage_verifyfail = getCeilNum(len(resList.UnPass), numberOnePage)
	totalNum.TotalPage_verifypass = getCeilNum(len(resList.OnLine), numberOnePage)
	return nil
}

func getCeilNum(a int, b int) int {
	c := 0.1
	c = float64(a) / float64(b)
	d := int(c)
	e := float64(d)

	if c > e {
		return d + 1
	} else {
		return d
	}

}

func GetProductDedicateID(st *ST_RequestPar, itemPerPage int) (lsID []string, e error) {
	listID := []string{}
	listPageID := []string{}
	if st.ProductTypeSub == 2 {
		listID = GetProductBodyPart(st.ProductType)

	} else if st.ProductTypeSub == 3 {
		listID = GetCustomerCollectProduct(st.UserID)
	} else if st.ProductTypeSub == 4 {
		listID = GetCustomerFootMarkProduct(st.UserID)
	}

	//Get the list entity according to the id list
	listStartDex := (st.RequestPage - 1) * itemPerPage
	if len(listID)-listStartDex <= itemPerPage {
		listPageID = listID[listStartDex:]
		return listPageID, nil
	} else {
		Info("List ID2=%v\n", listPageID)
		listPageID = listID[listStartDex : listStartDex+itemPerPage]
	}
	return listPageID, nil
}

func GetProductDedictateTotalNum(st *ST_RequestPar, numberOnePage int) int {
	listID := []string{}
	if st.ProductTypeSub == 2 {
		listID = GetProductBodyPart(st.ProductType)

	} else if st.ProductTypeSub == 3 {
		listID = GetCustomerCollectProduct(st.UserID)
	} else if st.ProductTypeSub == 4 {
		listID = GetCustomerFootMarkProduct(st.UserID)
	}
	totalNum := getCeilNum(len(listID), numberOnePage)
	return totalNum
}

//temp for get the List ID
func GetProductBodyPart(bodyPart string) (lsID []string) {
	ids := []string{}
	ids = append(ids, "wefwef")
	return ids
}

//temp for get the collect List ID
func GetCustomerCollectProduct(userID string) (lsID []string) {
	ids := []string{}
	ids = append(ids, "wefwef")
	return ids
}

//temp for get the footmark List ID
func GetCustomerFootMarkProduct(userID string) (lsID []string) {
	ids := []string{}
	ids = append(ids, "wefwef")
	return ids
}

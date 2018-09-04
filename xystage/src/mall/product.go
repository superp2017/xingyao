package main

import (
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	// "errors"

	"common"
	"constant"
	// "strconv"
	"cache/cacheIO"
	//"sort"
	. "util"
)

//产品详情
type ST_RequestProPar struct {
	City        string
	RequestArg  string
	RequestPage int
	BodyPart    string
	ProductItem string
	MinPrice    int
	MaxPrice    int
	SortType    string
}

type ST_ProductNet struct {
	ProductInfo     common.ST_Product
	LsProductCommon []common.ST_ProComment
}

func product_init() {

	JsDispatcher.Http("/getproducts", GetCityproduct)
	JsDispatcher.Http("/getproductsn", GetCityproductn)
	JsDispatcher.Http("/getproduct", GetNetProduct)
	JsDispatcher.Http("/getlsproduct", GetNetLsProduct)              //所有医院列表
	JsDispatcher.Http("/gettotalpagenumproduct", GettotalNumproduct) //所有医院列表

	JsDispatcher.Http("/getseconditem", GetSecondItem) //所有医院列表

}

func GettotalNumproduct(session *JsNet.StSession) {

	type ST_PageNumTotal struct {
		PageNum    int
		RequestArg string
		City       string
	}

	st := &ST_RequestProPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	lsproductID := GetproductList(st)
	totalNum := &ST_PageNumTotal{}
	totalNum.RequestArg = st.RequestArg
	totalNum.PageNum = getCeilNum(len(lsproductID), constant.ItemAccountPerPage_Product)
	// totalNum.HmTotalPageNum = make(map[string]int)
	// totalNum.HmTotalPageNum[st.RequireType] = getCeilNum(len(lsproductID), constant.ItemAccountPerPage_product)
	totalNum.City = st.City
	ForwardEx(session, "0", totalNum, st.RequestArg)
}

func GetCityproduct(session *JsNet.StSession) {

	lsproductSend := []*common.ST_Product{}

	//Get the request info
	st := &ST_RequestProPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	listID := GetDedicateListIDproduct(st, constant.ItemAccountPerPage_Product)

	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}
	lsproductSend = common.QueryMoreProducts(listID)
	// superproduct.Hmproduct[st.productType] = lsproductSend
	ForwardEx(session, "0", lsproductSend, st.RequestArg)
}

func GetCityproductn(session *JsNet.StSession) {

	//Get the request info
	st := &ST_RequestProPar{}
	if err := session.GetPara(st); err != nil {
		// Info("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^err!%s\v",err)

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Info("ST=%v\n", st)

	listID := GetDedicateListIDproduct(st, constant.ItemAccountPerPage_Product)
	// Info("********************************************************List ID=%v\n", listID)

	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}

	lsNetProduct := []*ST_ProductNet{}

	for _, v := range listID {
		proST_ProductNet, err := getNetProductL(v)
		if err == nil {
			// if proST_ProductNet.ProductInfo.Current.OpreatStatus == constant.OperatingStatus_online {

			lsNetProduct = append(lsNetProduct, proST_ProductNet)
			// }
		}
	}

	// superproduct.Hmproduct[st.productType] = lsproductSend
	ForwardEx(session, "0", lsNetProduct, st.RequestArg)
}

func GetDedicateListIDproduct(st *ST_RequestProPar, itemPerPage int) []string {

	//  Info("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~Get Dedicatelist ID Product with itemPage=%s\n",itemPerPage)
	listPageID := []string{}
	listID := GetproductList(st)

	listStartDex := (st.RequestPage - 1) * itemPerPage
	if len(listID)-listStartDex <= itemPerPage {
		listPageID = listID
		return listPageID
	} else {
		listPageID = listID[:listStartDex+itemPerPage]
	}
	//OutPutDocInfo("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@Get result\n")
	// Info("^^^^^^^^^^^^^^^^^^^^^^^^^^^List ID=%v\n", listPageID)
	return listPageID
}

func GetproductList(st *ST_RequestProPar) []string {
	// lsID := []string{}

	// for i := 0; i < 8; i++ {
	// 	hosID := "pro-" + strconv.Itoa(1001+i)
	// 	lsID = append(lsID, hosID)
	// }
	// return lsID

	// 	City        string
	// RequestArg  string
	// RequestPage int
	// BodyPart    string
	// ProductItem string
	// MinPrice    int
	// MaxPrice    int
	lsIDBackPro := []string{}
	lsPro, err := cacheIO.GetGlobalProPtr(st.SortType)
	if err != nil || lsPro == nil {
		return []string{}
	}

	lsProOnline := lsPro.OnLine

	//lsID, _ := common.GetCityItemDoc(cityName, bodypart)
	lsID, err := getProducts(st.City, st.BodyPart, st.ProductItem, st.MinPrice, st.MaxPrice)
	if err != nil {
		Error(err.Error())
		return lsIDBackPro
	}
	for _, v := range lsProOnline {
		for _, t := range lsID {
			if v == t {
				lsIDBackPro = append(lsIDBackPro, t)
			}
		}
	}

	IDB := make([]string, len(lsIDBackPro))

	if st.SortType == "SORT_PRICE" {
		for k, v := range lsIDBackPro {
			Info("%d;", k)
			IDB[len(IDB)-k-1] = v
		}

	} else {
		IDB = lsIDBackPro
	}

	//sort.Sort(sort.Reverse(sort.StringSlice(lsIDBackPro)))
	// getProducts(city, st.BodyPart, st.Item, st.MinPrice, st.MaxPrice)
	//Info("+++++++++++++++++++++++++++++++++lsOnLine=%s\n", lsProOnline)
	//Info("+++++++++++++++++++++++++++++++++lsID=%s\n", lsID)

	//return lsIDBackPro
	return IDB

}

func GetNetProduct(session *JsNet.StSession) {
	//get the request
	//hosST_HospitalNet := &ST_HospitalNet{}
	type st_query struct {
		ProID string //医院id
	}
	st := &st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ProID == "" {
		ForwardEx(session, "1", nil, "医院id为空,QueryHospitalInfo()，查询失败!")
		return
	}
	proST_ProductNet, err := getNetProductL(st.ProID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", proST_ProductNet)
}

func GetNetLsProduct(session *JsNet.StSession) {
	type st_queryproduct struct {
		ProID []string
	}

	lsNetProduct := []*ST_ProductNet{}

	st := &st_queryproduct{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	for _, v := range st.ProID {
		proST_ProductNet, err := getNetProductL(v)
		if err == nil {
			lsNetProduct = append(lsNetProduct, proST_ProductNet)
		}
	}
	Forward(session, "0", lsNetProduct)
}

func getNetProductL(productID string) (*ST_ProductNet, error) {
	productNet := &ST_ProductNet{}
	productNet.LsProductCommon = []common.ST_ProComment{}

	//get the product
	product, err := common.GetProductInfo(productID)
	// if product.Current.OpreatStatus != constant.OperatingStatus_online {
	// 	return nil, errors.New("The product status is not online")

	// }
	if err != nil {
		return nil, err
	}
	productNet.ProductInfo = *product

	//get the common
	lsCommon, err := common.Query_product_comments(productID)
	//Info("------------------------Product ID=%s has the common=%v\n", productID, lsCommon)
	if err == nil && len(lsCommon) > 0 {
		for _, v := range lsCommon {
			if v != nil {
				productNet.LsProductCommon = append(productNet.LsProductCommon, *v)
			}
		}
	}

	return productNet, nil
}

func GetSecondItem(session *JsNet.StSession) {
	res, err := cacheIO.GetFirstSecondItemMap()
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Forward(session, "0", res)
}

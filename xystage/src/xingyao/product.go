package main

import (
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"cache/cacheLib"
	"common"
	"constant"
	"strconv"
	. "util"
)

func product_init() {
	JsDispatcher.Http("/GetDedicateProduct", common.GetDedicateProduct)
	JsDispatcher.Http("/GenerateValues", GenerateValues)
	JsDispatcher.Http("/getpagenum", GettotalNum)
	JsDispatcher.Http("/reviewnewproduct", ReviewNewProduct)
	JsDispatcher.Http("/reviewmodifiedproduct", ReviewModifiedProduct)
	JsDispatcher.Http("/forceproductoffline", ForceProductOffLine)
	JsDispatcher.Http("/getproducts", common.GetDedicateProduct)
}

//Get all the new product (Test)--tianfeng
func GenerateValues(session *JsNet.StSession) {
	//Info("**************Enter GetProList\n")

	globalCashProduct := cacheLib.ST_ONUMCache{}

	for i := 1; i < 20; i++ {
		t := &common.ST_Product{}
		t.ProName = "Face Beautiful" + strconv.Itoa(i)
		t.ProType = "Face" + strconv.Itoa(i)
		t.Hospital = "Hospital" + strconv.Itoa(i)
		t.ProID = strconv.Itoa(i)
		t.Current.OpreatStatus = constant.OperatingStatus_new

		DirectWrite(constant.Hash_HosProduct, t.ProID, t)
		globalCashProduct.New = append(globalCashProduct.New, t.ProID)

		t.Current.OpreatStatus = constant.OperatingStatus_new

		t.Current.OpreatStatus = constant.OperatingStatus_new
		globalCashProduct.New = append(globalCashProduct.New, t.ProID)

		DirectWrite(constant.Hash_HosProduct, t.ProID, t)
	}

	for i := 20; i < 40; i++ {
		t := &common.ST_Product{}

		t.ProName = "Face Beautiful" + strconv.Itoa(i)
		t.ProType = "Face" + strconv.Itoa(i)
		t.Hospital = "Hospital" + strconv.Itoa(i)
		t.ProID = strconv.Itoa(i)

		DirectWrite(constant.Hash_HosProduct, t.ProID, t)
		globalCashProduct.Modify = append(globalCashProduct.Modify, t.ProID)

		t.Current.OpreatStatus = constant.OperatingStatus_modify

		globalCashProduct.Modify = append(globalCashProduct.Modify, t.ProID)
		DirectWrite(constant.Hash_HosProduct, t.ProID, t)

	}

	for i := 40; i < 60; i++ {
		t := &common.ST_Product{}

		t.Current.OpreatStatus = constant.OperatingStatus_online

		t.ProName = "Face Beautiful" + strconv.Itoa(i)
		t.ProType = "Face" + strconv.Itoa(i)
		t.Hospital = "Hospital" + strconv.Itoa(i)
		t.ProID = strconv.Itoa(i)
		t.Current.OpreatStatus = constant.OperatingStatus_online

		DirectWrite(constant.Hash_HosProduct, t.ProID, t)
		globalCashProduct.OnLine = append(globalCashProduct.OnLine, t.ProID)

		t.Current.OpreatStatus = constant.OperatingStatus_online

		t.Current.OpreatStatus = constant.OperatingStatus_online
		globalCashProduct.OnLine = append(globalCashProduct.OnLine, t.ProID)
		DirectWrite(constant.Hash_HosProduct, t.ProID, t)

	}

	for i := 60; i < 80; i++ {
		t := &common.ST_Product{}

		t.ProName = "Face Beautiful" + strconv.Itoa(i)
		t.ProType = "Face" + strconv.Itoa(i)
		t.Hospital = "Hospital" + strconv.Itoa(i)
		t.ProID = strconv.Itoa(i)
		t.Current.OpreatStatus = constant.OperatingStatus_Offline_self

		DirectWrite(constant.Hash_HosProduct, t.ProID, t)
		globalCashProduct.OffLine = append(globalCashProduct.OffLine, t.ProID)

		t.Current.OpreatStatus = constant.OperatingStatus_Offline_self
		globalCashProduct.OffLine = append(globalCashProduct.OffLine, t.ProID)
		DirectWrite(constant.Hash_HosProduct, t.ProID, t)

	}

	for i := 80; i < 100; i++ {
		t := &common.ST_Product{}

		t.ProName = "Face Beautiful" + strconv.Itoa(i)
		t.ProType = "Face" + strconv.Itoa(i)
		t.Hospital = "Hospital" + strconv.Itoa(i)
		t.ProID = strconv.Itoa(i)
		t.Current.OpreatStatus = constant.OperatingStatus_Reviewer_NotPass

		DirectWrite(constant.Hash_HosProduct, t.ProID, t)
		globalCashProduct.UnPass = append(globalCashProduct.UnPass, t.ProID)

		t.Current.OpreatStatus = constant.OperatingStatus_Reviewer_NotPass
		globalCashProduct.UnPass = append(globalCashProduct.UnPass, t.ProID)
		DirectWrite(constant.Hash_HosProduct, t.ProID, t)

	}

	DirectWrite(constant.Hash_Global, constant.Hash_ProductCache, globalCashProduct)

	ForwardEx(session, "0", nil, "")
}

func GettotalNum(session *JsNet.StSession) {
	//Get the request info
	type st_query struct {
		RequireType string
		SortType    string
	}

	st := &st_query{}
	totalNum := &common.ST_TotalNum{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if st.RequireType == "" {
		ForwardEx(session, "1", nil, "GettotalNum failed,RequireType %s,SortType%s\n", st.RequireType, st.SortType)
		return
	}

	err := common.GetTotalNum(st.SortType, st.RequireType, totalNum, constant.ItemAccountPerPage_Product)
	Info("Totalnum=%v\n", totalNum)

	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	ForwardEx(session, "0", totalNum, "sucess")
}

func ReviewNewProduct(session *JsNet.StSession) {
	common.RevieweNewProduct(session)
}

func ReviewModifiedProduct(session *JsNet.StSession) {
	common.RevieweModifyProduct(session)
}

func ForceProductOffLine(session *JsNet.StSession) {
	common.OfflineProductOnforce(session)
}

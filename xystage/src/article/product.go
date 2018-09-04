package article

import (
	"JsLib/JsNet"

	"common"
	"constant"
	. "util"
)

func GetDedicateProduct(session *JsNet.StSession) {
	st := &common.ST_RequestPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	listID, err := common.GetProductDedicateID(st, constant.ItemAccountPerPage_Product)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if len(listID) < 1 {
		ForwardEx(session, "1", nil, "the length is 0")
		return
	}

	lsProd := common.QueryMoreProducts(listID)

	ForwardEx(session, "0", lsProd, st.ProductType)

}

func GetDedicateProductTotalNum(session *JsNet.StSession) {
	type TotalNum struct {
		TotalNum       int
		RequestType    string
		RequestSubType int
	}
	tot := TotalNum{}

	st := &common.ST_RequestPar{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	totalNum := common.GetProductDedictateTotalNum(st, constant.ItemAccountPerPage_Product)
	tot.TotalNum = totalNum
	tot.RequestType = st.ProductType
	tot.RequestSubType = st.ProductTypeSub
	ForwardEx(session, "0", nil, tot.RequestType)

}

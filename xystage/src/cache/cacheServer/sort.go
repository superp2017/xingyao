package main

import (
	. "JsLib/JsLogger"
	. "cache/cacheLib"
	. "common"
	"sort"
	"sync"
)

var SaleHos SortList = SortList{}
var SaleHos_mutex sync.Mutex

var CommentHos SortList = SortList{}
var CommentHos_mutex sync.Mutex

var SaleDoc SortList = SortList{}
var SaleDoc_mutex sync.Mutex

var CommentDoc SortList = SortList{}
var CommentDoc_mutex sync.Mutex

var SalePro SortList = SortList{}
var SalePro_mutex sync.Mutex

var CommentPro SortList = SortList{}
var CommentPro_mutex sync.Mutex

var PricePro SortList = SortList{}
var PricePro_mutex sync.Mutex

func addToSalePro(ID string, Num int, check bool) {
	SalePro_mutex.Lock()
	defer SalePro_mutex.Unlock()
	SalePro.Append(ID, Num, check)
}

func addToCommentPro(ID string, Num int, check bool) {
	CommentPro_mutex.Lock()
	defer CommentPro_mutex.Unlock()

	CommentPro.Append(ID, Num, check)
}

func addToPricePro(ID string, Num int, check bool) {
	PricePro_mutex.Lock()
	defer PricePro_mutex.Unlock()
	PricePro.Append(ID, Num, check)
}

func addToSaleHos(ID string, Num int, check bool) {
	SaleHos_mutex.Lock()
	defer SaleHos_mutex.Unlock()

	SaleHos.Append(ID, Num, check)
}
func addToCommentHos(ID string, Num int, check bool) {
	CommentHos_mutex.Lock()
	defer CommentHos_mutex.Unlock()
	CommentHos.Append(ID, Num, check)
}

func addToSaleDoc(ID string, Num int, check bool) {
	SaleDoc_mutex.Lock()
	defer SaleDoc_mutex.Unlock()
	SaleDoc.Append(ID, Num, check)
}

func addToCommentDoc(ID string, Num int, check bool) {
	CommentDoc_mutex.Lock()
	defer CommentDoc_mutex.Unlock()
	CommentDoc.Append(ID, Num, check)
}

//产品销量排序
func SortSalePro() {
	SalePro_mutex.Lock()
	defer SalePro_mutex.Unlock()
	sort.Sort(SalePro)
}

//产品评论排序
func SortCommentPro() {
	CommentPro_mutex.Lock()
	defer CommentPro_mutex.Unlock()
	sort.Sort(CommentPro)
}

//产品价格排序
func SortPricePro() {
	PricePro_mutex.Lock()
	defer PricePro_mutex.Unlock()
	sort.Sort(PricePro)
}

///医院销量排序
func SortSaleHos() {
	SaleHos_mutex.Lock()
	defer SaleHos_mutex.Unlock()
	sort.Sort(SaleHos)
}

///医院评论排序
func SortCommentHos() {
	CommentHos_mutex.Lock()
	defer CommentHos_mutex.Unlock()
	sort.Sort(CommentHos)
}

///医生销量排序
func SortSaleDoc() {
	SaleDoc_mutex.Lock()
	defer SaleDoc_mutex.Unlock()
	sort.Sort(SaleDoc)
}

///医院评论排序
func SortCommentDoc() {
	CommentDoc_mutex.Lock()
	defer CommentDoc_mutex.Unlock()
	sort.Sort(CommentDoc)
}

///缓存排序
func SortONUM(Type string, obj int) *ST_ONUMCache {

	var ret ST_ONUMCache = ST_ONUMCache{}
	var list SortList = SortList{}
	var data *ST_ONUMCache = nil
	if obj == 0 {
		data = gl_HosONUM
		if Type == SORT_SALE {
			list = SaleHos
		} else if Type == SORT_COMMENT {
			list = CommentHos
		}
	} else if obj == 1 {
		data = gl_DocONUM
		if Type == SORT_SALE {
			list = SaleDoc
		} else if Type == SORT_COMMENT {
			list = CommentDoc
		}
	} else {
		data = gl_ProONUM
		if Type == SORT_SALE {
			list = SalePro
		} else if Type == SORT_COMMENT {
			list = CommentPro
		} else if Type == SORT_PRICE {
			list = PricePro
		}
	}
	if data == nil {
		return &ret
	}

	Error("list=%s\n", list)

	for _, v := range list {
		for _, i := range data.New {
			if v.ID == i {
				ret.New = append(ret.New, i)
			}
		}
		for _, i := range data.Modify {
			if v.ID == i {
				ret.Modify = append(ret.Modify, i)
			}
		}
		for _, i := range data.UnPass {
			if v.ID == i {
				ret.UnPass = append(ret.UnPass, i)
			}
		}
		for _, i := range data.OnLine {
			if v.ID == i {
				ret.OnLine = append(ret.OnLine, i)
			}
		}

		for _, i := range data.OffLine {
			if v.ID == i {
				ret.OffLine = append(ret.OffLine, i)
			}
		}
		for _, i := range data.Full {
			if v.ID == i {
				ret.Full = append(ret.Full, i)
			}
		}
	}
	Error("ret=%v\n", ret)

	return &ret
}

func clearSortcache() {
	SalePro_mutex.Lock()
	SalePro = SortList{}
	SalePro_mutex.Unlock()

	CommentPro_mutex.Lock()
	CommentPro = SortList{}
	CommentPro_mutex.Unlock()

	PricePro_mutex.Lock()
	PricePro = SortList{}
	PricePro_mutex.Unlock()

	SaleHos_mutex.Lock()
	SaleHos = SortList{}
	SaleHos_mutex.Unlock()

	CommentHos_mutex.Lock()
	CommentHos = SortList{}
	CommentHos_mutex.Unlock()

	SaleDoc_mutex.Lock()
	SaleDoc = SortList{}
	SaleDoc_mutex.Unlock()

	CommentDoc_mutex.Lock()
	CommentDoc = SortList{}
	CommentDoc_mutex.Unlock()

}

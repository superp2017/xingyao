package common

import (
	"sort"
	"time"
)

type OrderList []*ST_Order

func (list OrderList) Len() int {
	return len(list)
}

func (list OrderList) Less(i, j int) bool {

	t1, e1 := time.Parse("2006-01-02 15:04:05", list[i].Current.OpreatTime)
	t2, e2 := time.Parse("2006-01-02 15:04:05", list[j].Current.OpreatTime)
	if e1 == nil && e2 == nil {
		return t1.Unix() > t2.Unix()
	}
	return false
}

func (list OrderList) Swap(i, j int) {
	var temp *ST_Order = list[i]
	list[i] = list[j]
	list[j] = temp
}

//////////////////////////////////////////////////////////////////
type HosList []*ST_Hospital

func (list HosList) Len() int {
	return len(list)
}

func (list HosList) Less(i, j int) bool {

	t1, e1 := time.Parse("2006-01-02 15:04:05", list[i].Current.OpreatTime)
	t2, e2 := time.Parse("2006-01-02 15:04:05", list[j].Current.OpreatTime)
	if e1 == nil && e2 == nil {
		return t1.Unix() > t2.Unix()
	}
	return false
}

func (list HosList) Swap(i, j int) {
	var temp *ST_Hospital = list[i]
	list[i] = list[j]
	list[j] = temp
}

//////////////////////////////////////////////////////////////////
type DocList []*ST_Doctor

func (list DocList) Len() int {
	return len(list)
}

func (list DocList) Less(i, j int) bool {

	t1, e1 := time.Parse("2006-01-02 15:04:05", list[i].Current.OpreatTime)
	t2, e2 := time.Parse("2006-01-02 15:04:05", list[j].Current.OpreatTime)
	if e1 == nil && e2 == nil {
		return t1.Unix() > t2.Unix()
	}
	return false
}

func (list DocList) Swap(i, j int) {
	var temp *ST_Doctor = list[i]
	list[i] = list[j]
	list[j] = temp
}

//////////////////////////////////////////////////////////////////
type ProList []*ST_Product

func (list ProList) Len() int {
	return len(list)
}

func (list ProList) Less(i, j int) bool {

	t1, e1 := time.Parse("2006-01-02 15:04:05", list[i].Current.OpreatTime)
	t2, e2 := time.Parse("2006-01-02 15:04:05", list[j].Current.OpreatTime)
	if e1 == nil && e2 == nil {
		return t1.Unix() > t2.Unix()
	}
	return false
}

func (list ProList) Swap(i, j int) {
	var temp *ST_Product = list[i]
	list[i] = list[j]
	list[j] = temp
}

/////////////////////////////////////////////////////////////////////
type SortObject struct {
	Num int
	ID  string
}

type SortList []SortObject

func (list *SortList) Append(id string, num int, ischeck bool) {
	exist := false
	if ischeck {
		for i, v := range *list {
			if v.ID == id {
				exist = true
				(*list)[i].Num = num
			}
		}
	}
	if exist {
		sort.Sort(list)
	} else {
		*list = append(*list, SortObject{
			Num: num,
			ID:  id,
		})
	}
}

func (list SortList) Len() int {
	return len(list)
}

func (list SortList) Less(i, j int) bool {
	return list[i].Num > list[j].Num
}

func (list SortList) Swap(i, j int) {
	var temp SortObject = list[i]
	list[i] = list[j]
	list[j] = temp
}

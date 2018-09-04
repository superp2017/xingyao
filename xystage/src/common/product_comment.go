package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	. "util"
)

//评论列表
type ST_ProComment struct {
	ProID         string   //产品id
	UID           string   //用户id
	OrderID       string   //订单id
	UserName      string   //用户
	UserHead      string   //头像
	Comment       string   //评论内容
	ComPics       []string //图片列表
	Service       int      //服务产品评分
	Environmental int      //环境产品评分
	Effect        int      //效果产品评分
	IsAnonymous   string   //是否匿名
	ComTime       string   //评论时间
}

//新的产品评论
func AppendProductComment(session *JsNet.StSession) {
	st := &ST_ProComment{
		ComTime: CurTime(),
	}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ProID == "" || st.OrderID == "" || st.Comment == "" || st.Service < 0 || st.Environmental < 0 || st.Effect < 0 {
		ForwardEx(session, "1", nil, "评论失败,ProID=%s,Comment=%s,Service=%d,Environmental=%d,Effect=%d\n",
			st.ProID, st.Comment, st.Service, st.Environmental, st.Effect)
		return
	}
	data := []*ST_ProComment{}
	err := WriteLock(constant.Hash_HosProComment, st.ProID, &data)
	data = append(data, st)
	if err != nil {
		if e := DirectWrite(constant.Hash_HosProComment, st.ProID, &data); e != nil {
			ForwardEx(session, "1", data, e.Error())
			return
		}
	} else {
		if e := WriteBack(constant.Hash_HosProComment, st.ProID, &data); e != nil {
			ForwardEx(session, "1", data, e.Error())
			return
		}
	}

	////产品评分
	if err := product_score(st.ProID, st.Service, st.Environmental, st.Effect); err != nil {
		ErrorLog("product_score faild,err:%s\n", err.Error())
	}
	////更新订单的评论状态
	go OrderComment(st.OrderID, st.UID, st.UserName)
	Forward(session, "0", data)
}

//查询产品的评论
func QueryProductComments(session *JsNet.StSession) {
	type st_query struct {
		ProID string
	}
	st := st_query{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ProID == "" {
		ForwardEx(session, "1", nil, "查询评论失败,ProID为空\n")
		return
	}
	data, err := Query_product_comments(st.ProID)
	if err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	Forward(session, "0", data)
}

///查询产品的所有评论
func Query_product_comments(ProID string) ([]*ST_ProComment, error) {
	data := []*ST_ProComment{}
	if err := ShareLock(constant.Hash_HosProComment, ProID, &data); err != nil {
		return data, ErrorLog("查询失败,ShareLock(),ProID=%s\n", ProID)
	}
	return data, nil
}

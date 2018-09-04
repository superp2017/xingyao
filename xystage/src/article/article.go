package article

import (
	"JsLib/JsDispatcher"

	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"common"
	"constant"
	"encoding/json"
	"ider"
	"sort"
	"sync"
	"util"
	. "util"
)

//---------------------------------------------------
type ST_HinIn struct {
	Title      string
	TitleSub   string
	PictureURL string
}

//Store this to DB with xingyao manager managment
type ST_FirstPageInfoID struct {
	LsTopShowID       []string
	ExampleShowID     string
	RemarkShowID      string
	HinID             []string
	LsDoctorID        []string
	HmPreferArticleID map[string]([]string)
}

type ST_FirstPageInfo struct {
	LsTopShow   []ST_Article
	ExampleShow []ST_Article
	RemarkShow  []ST_Article
	Hin         []ST_Article
	LsDoctor    []common.ST_Doctor
}

//SuperArticelInfo
type DeskSearchArticlePush struct {
	ArticleSearch ST_Article
	KeyWords      string
	LsAllArticle  []ArticleTips
}

//---------------------------------------------------

type ST_Commons struct {
	Content    string
	Time       string
	AuthorID   string
	AuthorName string
}

type ST_XingYaoPic struct {
	PictureURL   string
	PictureNotes string
}

type ST_XingYaoLink struct {
	ProductURL string
	ProductDes string
}

//Store this to the DB including top article/Example/Remark three hash

type ST_Article struct {
	ArticleID    string
	Title        string
	TitleSub     string
	Instructions string
	AuthorNmae   string
	AuthorID     string
	ArticleTime  string
	Instruction  string
	Content      string
	CreatDate    string
	HotTips      string
	ShowPic      ST_XingYaoPic
	Commons      []ST_Commons
	LinkURL      []ST_XingYaoLink
}

type ArticleTips struct {
	ArticleID    string
	ArticleTitle string
}

type SuperArticleRet struct {
	SuperArticleID string
	PicURL         string
	MainTitle      string
	SubTitle       string
}

type DeskSearchArticle struct {
	ArticleSearch ST_Article
	KeyWords      string
}

type ProInfo struct {
	ProID   string //产品id
	ProName string //产品名称
}

type ST_SuperArticle struct {
	SuperArtileID string
	ArticleID     string
	ArticleTitle  string
	PicURL        string
	MainTitle     string
	SubTitle      string
	EntityTime    string
	EntityType    int //1--main page;2---four part select
	Position      int
	LsProduct     []string
	HmReader      map[string]int
	HmSupporter   map[string]int
	HmSharer      map[string]int
	HmShareTime   map[string]string

	Article ST_Article
}

type SuperArticleInfo struct {
	LsActiveTopShow   []ST_SuperArticle
	LsDeActiveTopShow []ST_SuperArticle
	LsActiveSubShow   []ST_SuperArticle
	LsDeActiveSubShow []ST_SuperArticle
	LsAgentArticle    []ST_SuperArticle
	LsAgentNotify     []ST_SuperArticle
	LsCustomerNotify  []ST_SuperArticle
	LsHospitalNotify  []ST_SuperArticle
	LsAllArticle      []ArticleTips
}

//SuperArticleInfoID ,for database
type ST_SuperArticleInfoID struct {
	LsActiveTopShowID   []string
	LsDeActiveTopShowID []string
	LsActiveSubShowID   []string
	LsDeActiveSubShowID []string
	LsAgentArticleID    []string
	LsAgentNotifyID     []string
	LsCustomerNotifyID  []string
	LsHospitalNotifyID  []string
	LsAllArticle        []ArticleTips
}

///////////////////////////////////////////////////
//
// added by mengzhaofeng 2017-03-08, to accelerate req
//
////////////////////////////////////////////////////
var global_superArticleInfo *SuperArticleInfo = nil
var global_LsActiveTopShow []*SuperArticleRet = nil
var global_LsActiveSubShow []*SuperArticleRet = nil
var super_mutex sync.Mutex

func initialSuperArticle(su *ST_SuperArticleInfoID) {
	if su.LsActiveTopShowID == nil {
		su.LsActiveTopShowID = []string{}
	}
	if su.LsDeActiveTopShowID == nil {
		su.LsDeActiveTopShowID = []string{}
	}

	if su.LsAllArticle == nil {
		su.LsAllArticle = []ArticleTips{}
	}

	if su.LsActiveSubShowID == nil {
		su.LsActiveSubShowID = make([]string, 4)
	}
	if su.LsAgentArticleID == nil {
		su.LsAgentArticleID = []string{}

	}
}
func InitalArticle() {
	//New
	JsDispatcher.Http("/newarticle", GenerateNewArticle)
	JsDispatcher.Http("/deletearticle", DeleteArticle) //删除文章
	JsDispatcher.Http("/updatearticle", UpdateArticle) //修改文章
	//Search
	JsDispatcher.Http("/getsearcharticle", GetSearchArticle)
	JsDispatcher.Http("/setsearcharticle", SetSearchArticle)
	//FrontPage
	JsDispatcher.Http("/removedeskarticle", DelSupArticle) //删除超级文章
	JsDispatcher.Http("/setdeskarticle", NewSuperArticle)
	JsDispatcher.Http("/modifysuparticle", ModifySupArticle) //修改超级文章
	JsDispatcher.Http("/delsuparticle", DelSupArticle)       //删除超级文章

	JsDispatcher.Http("/getactivetopshow", GetActiveTopShow)
	JsDispatcher.Http("/getactivesubshow", GetActiveSubShow)
	JsDispatcher.Http("/getdeskmanagecontent", GetFrontPageShowNet)
	JsDispatcher.Http("/gettotalpagenum", GetTotalPageNum)
	JsDispatcher.Http("/getarticles", GetDedicateArticle)
	JsDispatcher.Http("/getsuperarticles", GetDedicateArticle)
	JsDispatcher.Http("/getsuperarticlebyid", GetSuperArticleByID)
	JsDispatcher.Http("/getsuperarticlebyids", GetSuperArticleByIDS)

	JsDispatcher.Http("/addreading", NetAddReadingTimes)
	JsDispatcher.Http("/addsharearticle", NetAddShareTimes)
	JsDispatcher.Http("/addattention", NetAddAttentionTimes)
	JsDispatcher.Http("/recordlatlot", RecordLatLot)

	JsDispatcher.Http("/queryarticleinfo", queryArticleInfo) //查询单个文章内容

}

func GetActiveSubShow(session *JsNet.StSession) {
	if global_LsActiveSubShow == nil {
		ForwardEx(session, "1", "", "global_LsActiveSubShow is nil")
	} else {
		ForwardEx(session, "0", global_LsActiveSubShow, "success")
	}
}

func GetActiveTopShow(session *JsNet.StSession) {
	if global_LsActiveTopShow == nil {
		ForwardEx(session, "1", "", "global_LsActiveTopShow is nil")
	} else {
		ForwardEx(session, "0", global_LsActiveTopShow, "success")
	}
}

func GetSuperArticleByIDS(session *JsNet.StSession) {

	type Para struct {
		Ids []string
	}
	para := &Para{}

	if err := session.GetPara(&para); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	superArticle := make([]ST_SuperArticle, len(para.Ids))
	for k, v := range para.Ids {
		article, err := GetSuperArticleInfo(v)
		if err != nil {
			Error(err.Error())
			ForwardEx(session, "1", nil, err.Error())
			return
		}
		superArticle[k] = article
	}

	ForwardEx(session, "0", superArticle, "success")

}

func GetSuperArticleByID(session *JsNet.StSession) {
	var sid string
	if err := session.GetPara(&sid); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	superArticle, err := GetSuperArticleInfo(sid)
	if err == nil {
		ForwardEx(session, "0", superArticle, "success")
	} else {
		ForwardEx(session, "1", nil, err.Error())
	}

}

func GetFrontPageShowNet(session *JsNet.StSession) {
	superArticleInfo := GetFrontPageShow()
	ForwardEx(session, "0", superArticleInfo, "Sucess")
}

func NetAddReadingTimes(session *JsNet.StSession) {
	type RequestPar struct {
		UID            string
		SuperArticleID string
	}
	st := &RequestPar{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	_, err := addReadingTimes(st.SuperArticleID, st.UID)

	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	superArticleInfo, err := GetSuperArticleInfo(st.SuperArticleID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	ForwardEx(session, "0", &superArticleInfo, "Sucess")
}

func NetAddShareTimes(session *JsNet.StSession) {

	type RequestPar struct {
		UID            string
		SuperArticleID string
	}
	st := &RequestPar{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	_, err := addShareTimes(st.SuperArticleID, st.UID)

	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	superArticleInfo, err := GetSuperArticleInfo(st.SuperArticleID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	ForwardEx(session, "0", &superArticleInfo, "Sucess")
}

func NetAddAttentionTimes(session *JsNet.StSession) {
	type RequestPar struct {
		UID            string
		SuperArticleID string
		IsAttention    bool
	}
	st := &RequestPar{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	_, err := addAttentionTimes(st.SuperArticleID, st.UID, st.IsAttention)

	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	superArticleInfo, err := GetSuperArticleInfo(st.SuperArticleID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	ForwardEx(session, "0", &superArticleInfo, "Sucess")
}

func addReadingTimes(superArticleID string, UID string) (superArticleB *ST_SuperArticle, e error) {
	superArticle := &ST_SuperArticle{}
	errLockRead := util.WriteLock(constant.Hash_SuperArticle, superArticleID, superArticle)
	if errLockRead != nil {
		return nil, errLockRead
	}

	if superArticle.HmReader == nil {
		superArticle.HmReader = make(map[string]int)
	}
	_, ok := superArticle.HmReader[UID]
	if !ok {
		superArticle.HmReader[UID] = 0
	}
	superArticle.HmReader[UID] = superArticle.HmReader[UID] + 1
	util.WriteBack(constant.Hash_SuperArticle, superArticleID, superArticle)
	return superArticle, nil
}

func addShareTimes(superArticleID string, UID string) (superArticleB *ST_SuperArticle, e error) {

	superArticle := &ST_SuperArticle{}
	errLockRead := util.WriteLock(constant.Hash_SuperArticle, superArticleID, superArticle)
	if errLockRead != nil {
		return nil, errLockRead
	}

	if superArticle.HmSharer == nil {
		superArticle.HmSharer = make(map[string]int)
	}
	_, ok := superArticle.HmSharer[UID]
	if !ok {
		superArticle.HmSharer[UID] = 0
	}
	superArticle.HmSharer[UID] = superArticle.HmSharer[UID] + 1
	util.WriteBack(constant.Hash_SuperArticle, superArticleID, superArticle)
	return superArticle, nil
}

func addAttentionTimes(superArticleID string, UID string, isAdd bool) (superArticleB *ST_SuperArticle, e error) {

	superArticle := &ST_SuperArticle{}
	errLockRead := util.WriteLock(constant.Hash_SuperArticle, superArticleID, superArticle)
	if errLockRead != nil {
		return nil, errLockRead
	}

	if superArticle.HmSupporter == nil {
		superArticle.HmSupporter = make(map[string]int)
	}
	_, ok := superArticle.HmSupporter[UID]
	if !ok {
		superArticle.HmSupporter[UID] = 0
	}
	if isAdd {
		superArticle.HmSupporter[UID] = superArticle.HmSupporter[UID] + 1

	} else {
		if superArticle.HmSupporter[UID] > 0 {
			superArticle.HmSupporter[UID] = superArticle.HmSupporter[UID] - 1
		}
	}

	util.WriteBack(constant.Hash_SuperArticle, superArticleID, superArticle)
	return superArticle, nil

}

func GetFrontPageShowHospital(session *JsNet.StSession) {
	superArticleInfo := GetFrontPageShow()
	if superArticleInfo == nil {
		ForwardEx(session, "1", nil, "GetFrontPageShow error\n")
		return
	}
	// LsHospitalNotify  []ST_SuperArticle
	ForwardEx(session, "0", superArticleInfo.LsHospitalNotify, "Sucess")
}

func NewSuperArticle(session *JsNet.StSession) {

	Info("*****************************************Entering the new super article")
	st := &ST_SuperArticle{}
	article := &ST_Article{}

	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	err := ShareLock(constant.Hash_Article, st.ArticleID, article)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	id, _ := ider.GenSuperArtID()
	if id == "" {
		ForwardEx(session, "1", nil, "id is null")
		return
	}

	st.SuperArtileID = id

	st.EntityTime = CurTime()
	st.ArticleTitle = article.Title
	if st.MainTitle == "" {
		st.MainTitle = article.Title
	}
	if st.SubTitle == "" {
		st.SubTitle = article.Title
	}

	if err := DirectWrite(constant.Hash_SuperArticle, id, st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	//abstract the choised article,add to all article and tips
	deskManagerID := &ST_SuperArticleInfoID{}
	errLockRead := util.WriteLock(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)
	if errLockRead != nil {
		Error(errLockRead.Error())
		deskManagerID = &ST_SuperArticleInfoID{}
	}
	initialSuperArticle(deskManagerID)
	// LsAllArticle:= []ArticleTips
	//deskManagerID.LsAllArticle=append(deskManagerID.LsAllArticle,id)
	if st.EntityType == 1 {
		deskManagerID.LsActiveTopShowID = append(deskManagerID.LsActiveTopShowID, id)
	} else if st.EntityType == 2 {

		if st.Position > 0 && st.Position < (len(deskManagerID.LsActiveSubShowID)+1) {

			if deskManagerID.LsActiveSubShowID[st.Position-1] != "" {
				if deskManagerID.LsDeActiveSubShowID == nil {
					deskManagerID.LsDeActiveSubShowID = []string{}
				}
				deskManagerID.LsDeActiveSubShowID = append(deskManagerID.LsDeActiveSubShowID, deskManagerID.LsActiveSubShowID[st.Position-1])
			}
			deskManagerID.LsActiveSubShowID[st.Position-1] = id
		}

	} else if st.EntityType == 3 {
		deskManagerID.LsAgentArticleID = append(deskManagerID.LsAgentArticleID, id)

	} else if st.EntityType == 4 {
		deskManagerID.LsHospitalNotifyID = append(deskManagerID.LsHospitalNotifyID, id)
	} else if st.EntityType == 5 {
		deskManagerID.LsAgentNotifyID = append(deskManagerID.LsAgentNotifyID, id)
	} else if st.EntityType == 6 {
		deskManagerID.LsCustomerNotifyID = append(deskManagerID.LsCustomerNotifyID, id)
	}
	if errLockRead != nil {

		Error(errLockRead.Error())
		DirectWrite(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)

	} else {

		util.WriteBack(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)
	}
	//------------------------------------------
	frontPageShow := GetFrontPageShow()

	ForwardEx(session, "0", frontPageShow, "Sucess")
}

//修改SupArticle
func ModifySupArticle(session *JsNet.StSession) {
	st := &ST_SuperArticle{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data := &ST_SuperArticle{}
	if err := util.WriteLock(constant.Hash_SuperArticle, st.SuperArtileID, st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.ArticleID = st.ArticleID
	data.ArticleTitle = st.ArticleTitle
	data.PicURL = st.PicURL
	data.MainTitle = st.MainTitle
	data.SubTitle = st.SubTitle
	data.EntityType = st.EntityType
	data.Position = st.Position
	data.LsProduct = st.LsProduct
	data.HmReader = st.HmReader
	data.HmSupporter = st.HmSupporter
	data.HmSharer = st.HmSharer
	data.HmShareTime = st.HmShareTime
	if err := util.WriteLock(constant.Hash_SuperArticle, st.SuperArtileID, st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func initialSuperArticleInfo(superArticleInfo *SuperArticleInfo) {
	if superArticleInfo.LsActiveTopShow == nil {
		superArticleInfo.LsActiveTopShow = []ST_SuperArticle{}
	}

	if superArticleInfo.LsActiveSubShow == nil {
		superArticleInfo.LsActiveSubShow = make([]ST_SuperArticle, 4)
	}

	if superArticleInfo.LsDeActiveTopShow == nil {
		superArticleInfo.LsDeActiveTopShow = []ST_SuperArticle{}
	}

	if superArticleInfo.LsAgentArticle == nil {
		superArticleInfo.LsAgentArticle = []ST_SuperArticle{}
	}

	if superArticleInfo.LsHospitalNotify == nil {
		superArticleInfo.LsHospitalNotify = []ST_SuperArticle{}
	}

	if superArticleInfo.LsAgentNotify == nil {
		superArticleInfo.LsAgentNotify = []ST_SuperArticle{}
	}

	if superArticleInfo.LsCustomerNotify == nil {
		superArticleInfo.LsCustomerNotify = []ST_SuperArticle{}
	}

	if superArticleInfo.LsAllArticle == nil {
		superArticleInfo.LsAllArticle = []ArticleTips{}
	}

}

func GetFrontPageShow() *SuperArticleInfo {

	super_mutex.Lock()
	defer super_mutex.Unlock()

	superArticleInfo := &SuperArticleInfo{}
	global_superArticleInfo = superArticleInfo
	superArticleInfoID := &ST_SuperArticleInfoID{}
	initialSuperArticleInfo(superArticleInfo)
	// LsAllArticle []ArticleTips
	ShareLock(constant.Hash_ArticleManager, constant.Key_ArticleManager, superArticleInfoID)
	initialSuperArticle(superArticleInfoID)
	//pull out the main page history
	idnum := 0
	if len(superArticleInfoID.LsDeActiveTopShowID) < constant.ArticleNumPerPage_His {
		idnum = len(superArticleInfoID.LsDeActiveTopShowID)

	} else {
		idnum = constant.ArticleNumPerPage_All
	}
	allIDs := []string{}
	for _, v := range superArticleInfoID.LsDeActiveTopShowID[:idnum] {
		allIDs = append(allIDs, v)
	}
	//lsarticle := (allIDs)
	lssuperarticle_mainpage_his := GetLsSuperArticles(allIDs)
	superArticleInfo.LsDeActiveTopShow = lssuperarticle_mainpage_his
	//get the history sub page history article

	idnum = 0
	if len(superArticleInfoID.LsDeActiveSubShowID) < constant.ArticleNumPerPage_His {
		idnum = len(superArticleInfoID.LsDeActiveSubShowID)
	} else {
		idnum = constant.ArticleNumPerPage_All
	}
	allIDs = []string{}
	for _, v := range superArticleInfoID.LsDeActiveSubShowID[:idnum] {
		allIDs = append(allIDs, v)
	}
	lssuperarticle_subpage_his := GetLsSuperArticles(allIDs)
	superArticleInfo.LsDeActiveSubShow = lssuperarticle_subpage_his

	//get the history sub page history article

	allIDs = superArticleInfoID.LsActiveTopShowID
	lssuperarticle_mainpage_active := GetLsSuperArticles(allIDs)

	superArticleInfo.LsActiveTopShow = lssuperarticle_mainpage_active

	global_LsActiveTopShow = make([]*SuperArticleRet, len(lssuperarticle_mainpage_active))
	for k, v := range lssuperarticle_mainpage_active {
		art := &SuperArticleRet{}
		art.MainTitle = v.MainTitle
		art.SubTitle = v.SubTitle
		art.PicURL = v.PicURL
		art.SuperArticleID = v.SuperArtileID
		global_LsActiveTopShow[k] = art
	}
	//get the history sub page history article
	// allIDs = superArticleInfoID.LsActiveTopShowID
	// lssuperarticle_mainpage_active = GetLsSuperArticles(allIDs)

	// superArticleInfo.LsActiveTopShow = lssuperarticle_mainpage_active
	//get the four part show
	superArticleInfo.LsActiveSubShow = GetLsSuperArticles(superArticleInfoID.LsActiveSubShowID)

	global_LsActiveSubShow = make([]*SuperArticleRet, len(superArticleInfo.LsActiveSubShow))
	for k, v := range superArticleInfo.LsActiveSubShow {
		art := &SuperArticleRet{}
		art.MainTitle = v.MainTitle
		art.SubTitle = v.SubTitle
		art.PicURL = v.PicURL
		art.SuperArticleID = v.SuperArtileID
		global_LsActiveSubShow[k] = art
	}

	superArticleInfo.LsAllArticle = superArticleInfoID.LsAllArticle

	//get the agent article
	allIDs = superArticleInfoID.LsAgentArticleID
	superArticleInfo.LsAgentArticle = GetLsSuperArticles(superArticleInfoID.LsAgentArticleID)

	//get the hospital article
	idnum = 0
	if len(superArticleInfoID.LsHospitalNotifyID) < constant.ArticleNumPerPage_His {
		idnum = len(superArticleInfoID.LsHospitalNotifyID)
	} else {
		idnum = constant.ArticleNumPerPage_All
	}
	allIDs = []string{}
	for _, v := range superArticleInfoID.LsHospitalNotifyID[:idnum] {
		allIDs = append(allIDs, v)
	}
	lssuperarticle_hospitalnotify := GetLsSuperArticles(allIDs)
	superArticleInfo.LsHospitalNotify = lssuperarticle_hospitalnotify

	//get the agent article
	idnum = 0
	if len(superArticleInfoID.LsAgentNotifyID) < constant.ArticleNumPerPage_His {
		idnum = len(superArticleInfoID.LsAgentNotifyID)
	} else {
		idnum = constant.ArticleNumPerPage_All
	}
	allIDs = []string{}
	for _, v := range superArticleInfoID.LsAgentNotifyID[:idnum] {
		allIDs = append(allIDs, v)
	}
	lssuperarticle_agentnotify := GetLsSuperArticles(allIDs)
	superArticleInfo.LsAgentNotify = lssuperarticle_agentnotify

	//get the customer article
	idnum = 0
	if len(superArticleInfoID.LsCustomerNotifyID) < constant.ArticleNumPerPage_His {
		idnum = len(superArticleInfoID.LsCustomerNotifyID)
	} else {
		idnum = constant.ArticleNumPerPage_All
	}
	allIDs = []string{}
	for _, v := range superArticleInfoID.LsCustomerNotifyID[:idnum] {
		allIDs = append(allIDs, v)
	}
	lssuperarticle_customernotify := GetLsSuperArticles(allIDs)
	superArticleInfo.LsCustomerNotify = lssuperarticle_customernotify

	initialSuperArticleInfo(superArticleInfo)
	Info("Leave GetFrontPageShow-------------------------------\n")
	return superArticleInfo
}

func GetTotalPageNum(session *JsNet.StSession) {
	type articleTotalPageInfo struct {
		TotalPage_ArticleAll  int
		TotalPage_MainShowHis int
		TotalPage_SubShowHis  int
	}
	totalPageInfo := &articleTotalPageInfo{}
	deskManagerID := &ST_SuperArticleInfoID{}

	err := util.ShareLock(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	totalPageInfo.TotalPage_ArticleAll = getCeilNum(len(deskManagerID.LsAllArticle), constant.ArticleNumPerPage_All)
	totalPageInfo.TotalPage_MainShowHis = getCeilNum(len(deskManagerID.LsDeActiveTopShowID), constant.ArticleNumPerPage_His)
	totalPageInfo.TotalPage_SubShowHis = getCeilNum(len(deskManagerID.LsDeActiveSubShowID), constant.ArticleNumPerPage_His)
	ForwardEx(session, "0", totalPageInfo, "Sucess")
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

func GetSearchArticle(session *JsNet.StSession) {
	Info("Coming get search article\n")
	st := &DeskSearchArticle{}
	// deskManagerID := &ST_SuperArticleInfoID{}
	// err := util.ShareLock(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)

	pushArticle := &DeskSearchArticlePush{}
	err := util.ShareLock(constant.Hash_ArticleManager, constant.Key_ArticleSearch, st)
	if err != nil {
		Info("get wrong 1\n")

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	pushArticle.ArticleSearch = st.ArticleSearch
	pushArticle.KeyWords = st.KeyWords
	// pushArticle.LsAllArticle = deskManagerID.LsAllArticle
	Info("get 2\n")
	ForwardEx(session, "0", pushArticle, "Sucess")
	return
}

func SetSearchArticle(session *JsNet.StSession) {
	type SearchPar struct {
		ArticleID   string
		SearchValue string
		EntityTime  string
	}
	sp := &SearchPar{}
	article := &ST_Article{}

	if err := session.GetPara(sp); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	err := util.ShareLock(constant.Hash_Article, sp.ArticleID, article)
	if err := session.GetPara(sp); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	st := &DeskSearchArticle{}
	st.KeyWords = sp.SearchValue
	st.ArticleSearch = *article

	err = DirectWrite(constant.Hash_ArticleManager, constant.Key_ArticleSearch, st)
	if err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}

	ForwardEx(session, "0", st, err.Error())
	return

}

func RemoveActiveSuperArticles(session *JsNet.StSession) {
	st := &ST_SuperArticle{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	//abstract the choised article,add to all article and tips

	deskManagerID := &ST_SuperArticleInfoID{}
	errLockRead := util.WriteLock(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)
	if errLockRead != nil {
		deskManagerID = &ST_SuperArticleInfoID{}
	}
	initialSuperArticle(deskManagerID)

	if st.EntityType == 1 {
		for k, v := range deskManagerID.LsActiveTopShowID {
			if v == st.SuperArtileID {
				deskManagerID.LsActiveTopShowID = append(deskManagerID.LsActiveTopShowID[:k], deskManagerID.LsActiveTopShowID[k+1:]...)
				break
			}
		}
		deskManagerID.LsDeActiveTopShowID = append(deskManagerID.LsDeActiveTopShowID, st.SuperArtileID)
	} else if st.EntityType == 2 {
		deskManagerID.LsActiveSubShowID[st.Position-1] = ""
	} else if st.EntityType == 3 {
		for k, v := range deskManagerID.LsAgentArticleID {
			if v == st.SuperArtileID {
				deskManagerID.LsAgentArticleID = append(deskManagerID.LsAgentArticleID[:k], deskManagerID.LsAgentArticleID[k+1:]...)
				break
			}
		}
	}
	if errLockRead != nil {
		DirectWrite(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)

	} else {
		util.WriteBack(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)
	}
	//------------------------------------------
	frontPageShow := GetFrontPageShow()
	// frontPageShow = GetFrontPageShow()
	ForwardEx(session, "0", frontPageShow, "Sucess")
	return
}

// func GetDeskManageContent(session *JsNet.StSession) {
// 	deskManager := getDeskManageContentL()
// 	ForwardEx(session, "0", deskManager, "Sucess")
// 	return
// }

type ArticleRequestPar struct {
	ArticleType    string
	ArticleTypeSub int
	UserID         string
	RequestPage    int
}

func GetDedicateArticle(session *JsNet.StSession) {

	st := &ArticleRequestPar{}

	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}
	lsArticleID, _ := GetDedicateID(st)

	Info("Get the IDS=%v\n", lsArticleID)

	if st.ArticleType == constant.Article_editarticle {
		lsArticle := GetLsArticles(lsArticleID)
		ForwardEx(session, "0", lsArticle, "Sucess")
		return

	} else if st.ArticleType == constant.Article_toppageshow {
		lsSuperArticle := GetLsSuperArticles(lsArticleID)
		ForwardEx(session, "0", lsSuperArticle, "Sucess")
		return

	} else if st.ArticleType == constant.Article_fourpartshow {
		lsSuperArticle := GetLsSuperArticles(lsArticleID)
		ForwardEx(session, "0", lsSuperArticle, "Sucess")
		return
	}

	ForwardEx(session, "1", nil, "not found")
	return

}
func RecordLatLot(session *JsNet.StSession) {
	type LatLotRec struct {
		UID        string
		EntityTime string
		Lat        string
		Lot        string
	}
	st := &LatLotRec{}

	// 	st := &ST_Article{
	// 		CreatDate: CurTime(),
	// }

	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())

		return
	}

	st.EntityTime = CurTime()
	HDel(constant.Hash_LatLot, st.UID)
	if err := DirectWrite(constant.Hash_LatLot, st.UID, st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	ForwardEx(session, "0", nil, "sucess")
}

func GenerateNewArticle(session *JsNet.StSession) {

	st := &ST_Article{}

	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	id, _ := ider.GenArtID()
	if id == "" {
		ForwardEx(session, "1", nil, "id is none")
		return
	}
	st.ArticleID = id
	st.CreatDate = CurTime()
	if err := DirectWrite(constant.Hash_Article, id, st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	deskManagerID := &ST_SuperArticleInfoID{}
	errLockRead := util.WriteLock(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)
	if errLockRead != nil {
		deskManagerID = &ST_SuperArticleInfoID{}
	}
	if deskManagerID.LsAllArticle == nil {
		deskManagerID.LsAllArticle = []ArticleTips{}
	}
	articleTips := ArticleTips{}
	articleTips.ArticleID = st.ArticleID
	articleTips.ArticleTitle = st.Title
	deskManagerID.LsAllArticle = append(deskManagerID.LsAllArticle, articleTips)

	if errLockRead != nil {
		DirectWrite(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)
	} else {
		util.WriteBack(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)
	}
	idnum := 0
	if len(deskManagerID.LsAllArticle) < constant.ArticleNumPerPage_All {
		idnum = len(deskManagerID.LsAllArticle)
	} else {
		idnum = constant.ArticleNumPerPage_All
	}

	allIDs := []string{}
	for _, v := range deskManagerID.LsAllArticle[:idnum] {
		allIDs = append(allIDs, v.ArticleID)
	}
	lsarticle := GetLsArticles(allIDs)
	ForwardEx(session, "0", lsarticle, "Sucess")
}

func queryArticleInfo(session *JsNet.StSession) {
	type info struct {
		ArticleID string
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ArticleID == "" {
		ForwardEx(session, "1", nil, "ArticleID is empty\n")
	}
	data := &ST_Article{}
	if err := ShareLock(constant.Hash_Article, st.ArticleID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func DeleteArticle(session *JsNet.StSession) {
	type DelteArticle struct {
		ArticleID string
	}
	st := &DelteArticle{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ArticleID == "" {
		ForwardEx(session, "1", nil, "ArticleID is empty\n")
		return
	}
	if err := util.HDel(constant.Hash_Article, st.ArticleID); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	deskManagerID := &ST_SuperArticleInfoID{}
	if err := util.WriteLock(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID); err != nil {
		Forward(session, "0", nil)
		return
	}
	defer WriteBack(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)

	index := -1
	for k, v := range deskManagerID.LsAllArticle {
		if v.ArticleID == st.ArticleID {
			index = k
			break
		}
	}
	if index >= 0 && index < len(deskManagerID.LsAllArticle) {
		deskManagerID.LsAllArticle = append(deskManagerID.LsAllArticle[:index], deskManagerID.LsAllArticle[index+1:]...)
	}
	ForwardEx(session, "0", st, "Sucess")
}

//删除SupArticle
func DelSupArticle(session *JsNet.StSession) {
	type info struct {
		SuperArtileID string
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.SuperArtileID == "" {
		ForwardEx(session, "1", nil, "SuperArticleID is empty\n")
		return
	}
	if err := util.HDel(constant.Hash_SuperArticle, st.SuperArtileID); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	go func(id string) {
		data := &ST_SuperArticleInfoID{}
		if err := util.WriteLock(constant.Hash_ArticleManager, constant.Key_ArticleManager, data); err != nil {
			return
		}
		defer WriteBack(constant.Hash_ArticleManager, constant.Key_ArticleManager, data)

		index := -1
		for k, v := range data.LsActiveSubShowID {
			if v == id {
				index = k
				break
			}
		}
		if index >= 0 && index < len(data.LsActiveSubShowID) {
			data.LsActiveSubShowID = append(data.LsActiveSubShowID[:index], data.LsActiveSubShowID[index+1:]...)
		}
		index = -1
		for k, v := range data.LsActiveTopShowID {
			if v == id {
				index = k
				break
			}
		}
		if index >= 0 && index < len(data.LsActiveTopShowID) {
			data.LsActiveTopShowID = append(data.LsActiveTopShowID[:index], data.LsActiveTopShowID[index+1:]...)
		}
		index = -1
		for k, v := range data.LsAgentArticleID {
			if v == id {
				index = k
				break
			}
		}
		if index >= 0 && index < len(data.LsAgentArticleID) {
			data.LsAgentArticleID = append(data.LsAgentArticleID[:index], data.LsAgentArticleID[index+1:]...)
		}
		index = -1
		for k, v := range data.LsAgentNotifyID {
			if v == id {
				index = k
				break
			}
		}
		if index >= 0 && index < len(data.LsAgentNotifyID) {
			data.LsAgentNotifyID = append(data.LsAgentNotifyID[:index], data.LsAgentNotifyID[index+1:]...)
		}
		index = -1
		for k, v := range data.LsCustomerNotifyID {
			if v == id {
				index = k
				break
			}
		}
		if index >= 0 && index < len(data.LsCustomerNotifyID) {
			data.LsCustomerNotifyID = append(data.LsCustomerNotifyID[:index], data.LsCustomerNotifyID[index+1:]...)
		}

		index = -1
		for k, v := range data.LsHospitalNotifyID {
			if v == id {
				index = k
				break
			}
		}
		if index >= 0 && index < len(data.LsHospitalNotifyID) {
			data.LsHospitalNotifyID = append(data.LsHospitalNotifyID[:index], data.LsHospitalNotifyID[index+1:]...)
		}
	}(st.SuperArtileID)

	Forward(session, "0", nil)
}

func GetFrontPageDedicatedArticle(session *JsNet.StSession) {
	//Get the net parameters
	type ST_RequestParFrontPage struct {
		ArticleType string //1:new;2:modified;3 verifypass;4 verifypass;5 unsale
	}
	st := &ST_RequestParFrontPage{}
	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())
		return
	}
	//Get the FrontPageIDS
	frontPageIDs := &ST_FirstPageInfoID{}
	err := util.ShareLock(constant.Hash_FrontPageInfoLs, "constant.Key_FrontPage", frontPageIDs)
	if err != nil {
		return
	}
	//Get the Dedicate IDS
	lsArticleID := frontPageIDs.HmPreferArticleID[st.ArticleType]
	//Get the conresponding Article
	//lsProd := common.QueryMoreArticles(lsArticleID)
	//Return the conresponding Article
	ForwardEx(session, "0", lsArticleID, "sucess")
	return
}

func GetLsArticles(ids []string) []ST_Article {
	///查询多个订单信息
	data := []ST_Article{}
	if len(ids) == 0 {
		return data
	}
	for _, v := range ids {
		t, e := GetArticleInfo(v)
		if e != nil {
			continue
		}
		data = append(data, t)
	}
	return data
}

func GetArticleInfo(keyValue string) (st ST_Article, e error) {
	data := ST_Article{}
	if keyValue == "" {
		return data, ErrorLog("ID=nil!\n")
	}

	if err := ShareLock(constant.Hash_Article, keyValue, &data); err != nil {
		return data, ErrorLog("Search Front Article fail, ShareLock(),dbName=!;key=%s\n", keyValue)
	}
	return data, nil
}

func GetLsSuperArticles(ids []string) []ST_SuperArticle {

	sort.Sort(sort.Reverse(sort.StringSlice(ids)))

	///查询多个订单信息
	data := []ST_SuperArticle{}
	if len(ids) == 0 {
		return data
	}
	for _, v := range ids {
		t, e := GetSuperArticleInfo(v)
		if e != nil {
			continue
		}
		data = append(data, t)
	}
	return data
}

func GetSuperArticleInfo(keyValue string) (st ST_SuperArticle, e error) {
	data := ST_SuperArticle{}
	article := ST_Article{}
	if keyValue == "" {
		return data, ErrorLog("ID=nil!\n")
	}

	if err := ShareLock(constant.Hash_SuperArticle, keyValue, &data); err != nil {
		return data, ErrorLog("Search Front SuperArticle fail, ShareLock(),dbName=!;key=%s\n", keyValue)
	}

	ShareLock(constant.Hash_Article, data.ArticleID, &article)
	data.Article = article

	return data, nil
}

func GetDedicateID(st *ArticleRequestPar) (lsID []string, e error) {
	if st == nil {
		return []string{}, ErrorLog("GetDedicateID failed,ArticleRequestPar is nil\n")
	}

	itemPerPage := 10
	listID := []string{}
	listPageID := []string{}
	deskManagerID := &ST_SuperArticleInfoID{}
	errLockRead := util.ShareLock(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)
	if errLockRead != nil {
		deskManagerID = &ST_SuperArticleInfoID{}
	}

	initialSuperArticle(deskManagerID)

	if st.ArticleType == constant.Article_editarticle {
		itemPerPage = constant.ArticleNumPerPage_All
		for _, v := range deskManagerID.LsAllArticle {
			listID = append(listID, v.ArticleID)
		}
	} else if st.ArticleType == constant.Article_toppageshow {
		itemPerPage = constant.ArticleNumPerPage_His
		listID = deskManagerID.LsDeActiveTopShowID
	} else if st.ArticleType == constant.Article_fourpartshow {
		itemPerPage = constant.ArticleNumPerPage_His
		listID = deskManagerID.LsActiveSubShowID
	}

	//Get the list entity according to the id list
	listStartDex := (st.RequestPage - 1) * itemPerPage
	if listStartDex < 0 {
		listStartDex = 0
	}
	if listStartDex > len(listID) {
		listStartDex = len(listID)
	}

	if len(listID)-listStartDex <= itemPerPage {
		listPageID = listID[listStartDex:]
		return listPageID, nil
	} else {
		listPageID = listID[listStartDex : listStartDex+itemPerPage]
	}
	Info("List ID=%v\n", listPageID)
	return listPageID, nil
}

func UpdateArticle(session *JsNet.StSession) {
	// 	st.CreatDate = CurTime()
	Error("----------------------------------------------------------Update Article")
	st := &ST_Article{}

	if err := session.GetPara(st); err != nil {

		ForwardEx(session, "1", nil, err.Error())

		return
	}
	data := &ST_Article{}
	if err := WriteLock(constant.Hash_Article, st.ArticleID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer WriteBack(constant.Hash_Article, st.ArticleID, data)
	d, err := json.Marshal(data)
	if err != nil {
		Error("json.Marshal(m) error: %s", err.Error())
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if err = json.Unmarshal(d, st); err != nil {
		Error("json.Unmarshal(data, &v) error: %s", err.Error())
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	deskManagerID := &ST_SuperArticleInfoID{}
	errLockRead := util.WriteLock(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)
	if errLockRead != nil {
		deskManagerID = &ST_SuperArticleInfoID{}
	}
	if deskManagerID.LsAllArticle == nil {
		deskManagerID.LsAllArticle = []ArticleTips{}
	}
	articleTips := ArticleTips{}
	articleTips.ArticleID = st.ArticleID
	articleTips.ArticleTitle = st.Title

	for k, v := range deskManagerID.LsAllArticle {
		if v.ArticleID == st.ArticleID {
			deskManagerID.LsAllArticle[k] = articleTips
			break
		}

	}
	if errLockRead != nil {
		DirectWrite(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)

	} else {
		util.WriteBack(constant.Hash_ArticleManager, constant.Key_ArticleManager, deskManagerID)
	}
	//return the first page articles
	idnum := 0
	if len(deskManagerID.LsAllArticle) < constant.ArticleNumPerPage_All {
		idnum = len(deskManagerID.LsAllArticle)

	} else {
		idnum = constant.ArticleNumPerPage_All
	}

	allIDs := []string{}
	for _, v := range deskManagerID.LsAllArticle[:idnum] {
		allIDs = append(allIDs, v.ArticleID)
	}
	lsarticle := GetLsArticles(allIDs)
	ForwardEx(session, "0", lsarticle, "Sucess")
}

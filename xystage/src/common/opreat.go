package common

type ST_Opreat struct {
	OpreatPart   string //操作方:用户、医院、平台
	OpreatAction string //动作
	OpreatStatus string //状态
	OpreatTime   string //操作时间
	OpreatReason string //操作原因
	OpreatJobNum string //操作人员工号
	OpreatName   string //操作人员姓名
	OpreatCell   string //操作人联系方式
	ApplyJobNum  string //申请人员工号
	ApplyName    string //申请人员姓名
	ApplyCell    string //申请人联系方式
}

type ST_OpreatStatus struct {
	OpreatInfo      []ST_Opreat //历史操作记录和状态
	Current         ST_Opreat   //当前的操作信息
	FirstOnlineDate string      //初次审核通过时间
	OpreatReason    string      ///操作原因
	LastOnlineDate  string      //上次审核通过时间
	LastOnlineName  string      //上次审核人
	LastModifyDate  string      //上次修改时间
	CreatDate       string      //创建日期
}

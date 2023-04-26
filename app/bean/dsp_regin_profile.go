package bean

type DspRegionProfile struct {
	ProfileId string `json:"profile_id"`
	Region    string `json:"region"`
}

type ConsumerDetail struct {
	CtxId string
	// 时间戳
	CreateTime int64
	ReportType string
	ReportName string
	ProfileId  string
	// 0 初始化 1 成功 2 失败
	Status     int
	ReportDate string
	ErrDesc    string
	Batch      string
	CostTime   int64
	// 时间戳
	UpdateTime int64
}

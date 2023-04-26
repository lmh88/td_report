package bean

type Result struct {
	ModifileDate string
	FileName     string
	Exists       bool
	StartDate    string
	ProfileId    string
	Extrant      string // 标志sd类型的"T00020", "T00030", "remarketing"
}

// ChaData 计算服务器上面缺少的文件
type ChaData struct {
	ProfileId  string
	Date       string
	ReportType string
	ReportName string
}

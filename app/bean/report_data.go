package bean

// ReportData http接收请求的字段
type ReportData struct {
	ReportDataType int      `json:"reportDataType"`
	ReportName     []string `json:"reportName"`
	ReportType     string   `json:"reportType"`
	StartDate      string   `json:"startDate"`
	EndDate        string   `json:"endDate"`
	Profileids     []string `json:"profileids"`
	CallBackUrl    string   `json:"callBackUrl"`
	ProcessId      int      `json:"processId"`
	Batch          string   `json:"batch"`
}

// ReceiveMsgBody 消息体结构
type ReceiveMsgBody struct {
	MessageId   string `json:"MessageId"`
	MessageBody string `json:"MessageBody"`
}

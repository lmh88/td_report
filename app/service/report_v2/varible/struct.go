package varible

type Client struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type ClientRefreshToken struct {
	RefreshToken string `json:"refresh_token"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	ProfileId    string `json:"profile_id"`
}

type ProfileMsg struct {
	ReportType   string `json:"report_type"`
	ProfileId    string `json:"profile_id"`
	ProfileType    string `json:"profile_type"`
	Region       string `json:"region"`
	ReportDate   string `json:"report_date"`
	Timestamp    int64  `json:"timestamp"`
	RefreshToken string `json:"refresh_token"`
	BatchKey     string `json:"batch_key"`
	ClientTag    string `json:"clientTag"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type ReportIdMsg struct {
	ReportType    string         `json:"report_type"`
	ProfileId     string         `json:"profile_id"`
	ProfileType    string `json:"profile_type"`
	Region        string         `json:"region"`
	ReportDate    string         `json:"report_date"`
	Timestamp     int64          `json:"timestamp"`
	RefreshToken  string         `json:"refresh_token"`
	BatchKey      string         `json:"batch_key"`
	ClientTag     string         `json:"clientTag"`
	ClientId      string         `json:"clientId"`
	ClientSecret  string         `json:"clientSecret"`
	TraceId       string         `json:"trace_id"`
	ReportId      string         `json:"report_id"`
	ReportName    string         `json:"report_name"`
	ReportTactic  string         `json:"report_tactic"`   //sd专用
	RetryCount    map[string]int `json:"retry_count"`     //获取报表地址，尝试次数
	FirstTryCount int            `json:"first_try_count"` //第一次请求，尝试次数
}

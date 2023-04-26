package bean

type EventDetail struct {
	DataPath   string `json:"data_path"  binding:"required,min=1"`
	TaskType   string `json:"task_type"  binding:"required,min=1"`
	ReportType string `json:"report_type" binding:"required,min=1"`
	TaskStatus string `json:"task_status" binding:"required,min=1"`
	AdPlatform string `json:"ad_platform" binding:"required,min=1"`
	ReportDate string `json:"report_date" binding:"required,min=1"`
}

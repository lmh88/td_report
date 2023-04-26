package bean

// UploadS3Data s3上传
type UploadS3Data struct {
	Key  string
	Path string
}

// UploadS3DataCheck s3 上传统计
type UploadS3DataCheck struct {
	Key        string
	Path       string
	ProfileId  string
	ReportType string
	ReportName string
	StartDate  string
	EndDate    string
}

package v2

import (
	"net/http"
	"td_report/app/service/report_v3"
	"testing"
)

func TestGetReport(t *testing.T) {
	profile :=  getByProfile("2697320314504359")
	report_v3.CreateReport(profile)
}

func TestGetReportStatus(t *testing.T) {
	profile :=  getByProfile("2697320314504359")
	report_v3.GetReportStatus(profile, "356f1d84-59b3-45bb-bdb9-42983d63032f")


	http.Get("")
}

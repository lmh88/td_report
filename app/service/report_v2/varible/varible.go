package varible

import (
	"github.com/gogf/gf/frame/g"
	"td_report/vars"
)

const (
	// SpecialDay 特定天数，如果是2天的情况就需要过滤黑名单
	SpecialDay = 2
	// Retry429 429重试标记
	Retry429 = "retry_429"
	Seller = "seller"
	Vendor = "vendor"
)

const (
	ReportDefaultExchange = "report_default" //交换机名

	SpFastQueue     = "sp_profile_fast"
	SpMiddleQueue   = "sp_profile_middle"
	SpSlowQueue     = "sp_profile_slow"
	SpReportQueue   = "sp_report_ids"
	SpDelay10s      = "sp_delay_10s"
	SpDelay30s      = "sp_delay_30s"
	SpDelay60s      = "sp_delay_60s"
	SpDelay120s     = "sp_delay_120s"
	SpFailQueue     = "sp_fail"
	SpRetryQueue    = "sp_retry"
	SpRetryDelay20s = "sp_retry_delay_20s"

	SbFastQueue     = "sb_profile_fast"
	SbMiddleQueue   = "sb_profile_middle"
	SbSlowQueue     = "sb_profile_slow"
	SbReportQueue   = "sb_report_ids"
	SbDelay10s      = "sb_delay_10s"
	SbDelay30s      = "sb_delay_30s"
	SbDelay60s      = "sb_delay_60s"
	SbDelay120s     = "sb_delay_120s"
	SbFailQueue     = "sb_fail"
	SbRetryQueue    = "sb_retry"
	SbRetryDelay20s = "sb_retry_delay_20s"

	SdFastQueue     = "sd_profile_fast"
	SdMiddleQueue   = "sd_profile_middle"
	SdSlowQueue     = "sd_profile_slow"
	SdReportQueue   = "sd_report_ids"
	SdDelay10s      = "sd_delay_10s"
	SdDelay30s      = "sd_delay_30s"
	SdDelay60s      = "sd_delay_60s"
	SdDelay120s     = "sd_delay_120s"
	SdFailQueue     = "sd_fail"
	SdRetryQueue    = "sd_retry"
	SdRetryDelay20s = "sd_retry_delay_20s"
)

var QueueMap = map[string]map[string]string{
	vars.SP: {
		FastLevel:   SpFastQueue,
		MiddleLevel: SpMiddleQueue,
		SlowLevel:   SpSlowQueue,
	},
	vars.SB: {
		FastLevel:   SbFastQueue,
		MiddleLevel: SbMiddleQueue,
		SlowLevel:   SbSlowQueue,
	},
	vars.SD: {
		FastLevel:   SdFastQueue,
		MiddleLevel: SdMiddleQueue,
		SlowLevel:   SdSlowQueue,
	},
}

var ReportQueueMap = map[string]string{
	vars.SP: SpReportQueue,
	vars.SB: SbReportQueue,
	vars.SD: SdReportQueue,
}

const (
	FastLevel   = "fast"
	MiddleLevel = "middle"
	SlowLevel   = "slow"
)

var QueueLevelArray = [3]string{FastLevel, MiddleLevel, SlowLevel}

const (
	WaitTime10s  = "wt10s"
	WaitTime30s  = "wt30s"
	WaitTime60s  = "wt60s"
	WaitTime120s = "wt120s"
)

var WaitTimeArray = []string{WaitTime10s, WaitTime30s, WaitTime60s, WaitTime120s}

var RetryCount = map[string]int{
	WaitTime10s:  1,
	WaitTime30s:  3,
	WaitTime60s:  5,
	WaitTime120s: 20,
}

var DelayQueueMap = map[string]map[string]string{
	vars.SP: {
		WaitTime10s:  SpDelay10s,
		WaitTime30s:  SpDelay30s,
		WaitTime60s:  SpDelay60s,
		WaitTime120s: SpDelay120s,
	},
	vars.SB: {
		WaitTime10s:  SbDelay10s,
		WaitTime30s:  SbDelay30s,
		WaitTime60s:  SbDelay60s,
		WaitTime120s: SbDelay120s,
	},
	vars.SD: {
		WaitTime10s:  SdDelay10s,
		WaitTime30s:  SdDelay30s,
		WaitTime60s:  SdDelay60s,
		WaitTime120s: SdDelay120s,
	},
}

var DelayTimeMap = map[string]int{
	WaitTime10s:  10,
	WaitTime30s:  30,
	WaitTime60s:  60,
	WaitTime120s: 120,
}

var FailQueueMap = map[string]string{
	vars.SP: SpFailQueue,
	vars.SB: SbFailQueue,
	vars.SD: SdFailQueue,
}

var RetryQueueMap = map[string]string{
	vars.SP: SpRetryQueue,
	vars.SB: SbRetryQueue,
	vars.SD: SdRetryQueue,
}

var FirstRetryCountMap = map[string]int{
	vars.SP: 20,
	vars.SB: 20,
	vars.SD: 20,
}

var RetryDelayMap = map[string]string{
	vars.SP: SpRetryDelay20s,
	vars.SB: SbRetryDelay20s,
	vars.SD: SdRetryDelay20s,
}

var RetryDelayTimeMap = map[string]int{
	vars.SP: 20,
	vars.SB: 20,
	vars.SD: 20,
}

var OvertimeNoticeMap = map[string]int64{
	vars.SP: g.Cfg().GetInt64("overtime_notice.sp", 3600),
	vars.SB: g.Cfg().GetInt64("overtime_notice.sb", 3600),
}

var LimitRetryQueueMap = map[string]int{
	vars.SP: g.Cfg().GetInt("limit_queue.retry.sp", 300),
	vars.SB: g.Cfg().GetInt("limit_queue.retry.sb", 300),
	vars.SD: g.Cfg().GetInt("limit_queue.retry.sd", 300),
}

var LimitReportQueueMap = map[string]int{
	vars.SP: g.Cfg().GetInt("limit_queue.report.sp", 2400),
	vars.SB: g.Cfg().GetInt("limit_queue.report.sb", 2400),
	vars.SD: g.Cfg().GetInt("limit_queue.report.sd", 2400),
}

var ClientMap = map[string]Client{
	"c1": {
		ClientId: "amzn1.application-oa2-client.084663234c2143c3a3bf91fe34bbdf1e",
		ClientSecret: "29f982cfd2571585c52db5ba462f302716c87fd52507dcc355582b8821a80abc",
	},
	"c2": {
		ClientId: "amzn1.application-oa2-client.be4e30f5a0b14f488677728ec04c12e0",
		ClientSecret: "c8c0f06cd18861b9cd5b4cd228bf475ca08e7f6f7af4dad384e0772296a3dd24",
	},
}

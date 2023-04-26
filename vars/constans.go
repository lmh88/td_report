package vars

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
)

const (
	// SPL 定义campaignType 常量的长形式
	SPL  = "sponsoredProducts"
	SDL  = "sponsoredDisplay"
	SBL  = "sponsoredBrands"
	DSPL = "dsp"

	// SP 定义campaignType 常量的短形式
	SP       = "sp"
	SD       = "sd"
	SB       = "sb"
	SB_BRAND = "sb_brand"
	SB_ALL   = "sb_all"
	DSP      = "dsp"
	DspTemp  = "dsp_temp"
)

// 报表生产者天数和前缀
const (
	// ProductDay2 2天一个批次
	ProductDay2 = 2
	// ProductDay7 7天一个批次
	ProductDay7 = 7
	// ProductDay14 14天一个批次
	ProductDay14 = 14
	// ProductDay30 30天一个批次
	ProductDay30 = 30

	// 针对生产者有4种队列 2天进入的是快速队列 ， 7天的进入的是中塑队列 14天进入慢速队列 ,一个月的是后备队列
	FastQueue   = "fast"
	MiddleQueue = "middle"
	SlowQueue   = "slow"
	BackQueue   = "back"
)

// sb特殊报表类型
const (
	BrandMetricsWeekly  = "brandMetricsWeekly"
	BrandMetricsMonthly = "brandMetricsMonthly"
)

const (
	Error401txt   = "HTTP 401 Unauthorized"
	Timeouttxt    = "Client.Timeout exceeded while awaiting headers"
	Timeouttxt1   = "TLS handshake timeout"
	Timeouttxt2   = "i/o timeout"
	TimeFormatTpl = "2006-01-02"
	TimeLayout    = "20060102"
	TIMEFORMAT    = "2006-01-02 03:04:05"
	TIMEDSPFILE   = "2-Jan-06"
	Timezon       = "UTC"
)

const (
	SbPath      = "/sb"
	SdPath      = "/new_sd"
	SpPath      = "/new_sp"
	DspPath     = "/new_dsp"
	DspPathTemp = "/new_dsp_temp"
)

var RedisQueueList = []string{
	FastQueue,
	MiddleQueue,
	SlowQueue,
	BackQueue,
}

var PathMap = map[string]string{
	SD:  SdPath,
	SP:  SpPath,
	DSP: DspPath,
	SB:  SbPath,
}

var (
	SbRate  = g.Cfg().GetInt("limit.sb_rate")
	SpRate  = g.Cfg().GetInt("limit.sp_rate")
	DspRate = g.Cfg().GetInt("limit.dsp_rate")
	SdRate  = g.Cfg().GetInt("limit.sd_rate")
)

// LimitMap 评率限制
var LimitMap = map[string]int{
	SD:  SdRate,
	SP:  SpRate,
	DSP: DspRate,
	SB:  SbRate,
}

var MypathMap = map[string]string{
	SD:  g.Cfg().GetString("common.datapath") + SdPath,
	SP:  g.Cfg().GetString("common.datapath") + SpPath,
	SB:  g.Cfg().GetString("common.datapath") + SbPath,
	DSP: g.Cfg().GetString("common.datapath") + DspPath,
}

// Cache token直接缓存
var Cache = gcache.New()

// s3上传时间类型
const (
	TimeDay   = "day"
	TimeWeek  = "week"
	TimeMonth = "month"
)

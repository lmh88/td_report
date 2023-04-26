package vars

// ReportList 目前所有报表名称
var ReportList = map[string][]string{
	SB:       {"adGroup", "adGroupVideo", "campaigns", "campaignsVideo", "keywords", "keywordsVideo", "keywordsQuery", "keywordsQueryVideo", "targets", "targetsVideo", "ads"},
	SP:       {"adGroups", "asins", "campaigns", "campaignsPlacement", "keywords", "keywordsQuery", "productAds", "targetQuerys", "targets"},
	SD:       {"adGroups", "asins", "campaigns", "productAds", "targets"},
	DSP:      {"audience", "inventory", "order", "detail"},
	SB_ALL:   {"adGroup", "adGroupVideo", "campaigns", "campaignsVideo", "keywords", "keywordsVideo", "keywordsQuery", "keywordsQueryVideo", "targets", "targetsVideo", "brandMetricsWeekly", "brandMetricsMonthly"},
	SB_BRAND: {"brandMetricsWeekly", "brandMetricsMonthly"},
	DspTemp : {"order","detail"},
}

var ReportListFileExt = map[string]string{
	SP:       "gz",
	SD:       "gz",
	DSP:      "csv",
	SB:       "gz",
	SB_BRAND: "json",
}

// S3ReportMap 报表业务和老的s3报表名称映射关系
var S3ReportMap = map[string]string{
	"sb_adGroupVideo":        "ad_group_video",
	"sb_adGroup":             "ad_group",
	"sb_brandMetricsMonthly": "brand_metrics_month",
	"sb_brandMetricsWeekly":  "brand_metrics_week",
	"sb_campaignsVideo":      "campaign_video",
	"sb_campaigns":           "campaign",
	"sb_keywordsQueryVideo":  "keyword_query_video",
	"sb_keywordsQuery":       "keyword_query",
	"sb_keywordsVideo":       "keyword_video",
	"sb_keywords":            "keyword",
	"sb_targetsVideo":        "target_video",
	"sb_targets":             "target",
	"sb_ads":             "ad",

	"sd_adGroups":   "ad_group",
	"sd_asins":      "asin",
	"sd_campaigns":  "campaign",
	"sd_productAds": "product_ad",
	"sd_targets":    "target",

	"sp_adGroups":           "ad_group",
	"sp_asins":              "asin",
	"sp_campaigns":          "campaign",
	"sp_keywordsQuery":      "keyword_query",
	"sp_keywords":           "keyword",
	"sp_campaignsPlacement": "placement",
	"sp_productAds":         "product_ad",
	"sp_targetQuerys":       "target_query",
	"sp_targets":            "target",

	"dsp_audience":  "audience",
	"dsp_detail":    "detail",
	"dsp_inventory": "inventory",
	"dsp_order":     "order",
}

package vars

// AmazonRequestMsg 亚马逊请求结果
type AmazonRequestMsg struct {
	Code      int         `json:"code"`
	RequestId string      `json:"requestId"`
	Details   string      `json:"details,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type ResultMsg struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResultArrMsg struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

type ResultMapMsg struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type ResultMapsMsg struct {
	Code    string                   `json:"code"`
	Message string                   `json:"message"`
	Data    map[string][]interface{} `json:"data"`
}

// CampaignType 类型
var CampaignType = []string{SB, SP, DSP, SD}

type RegionType string

type BaseParams struct {
	Region    RegionType  `json:"region"`
	ProfileId int64       `json:"profileId"`
	TokenId   int         `json:"tokenId"`
	Data      interface{} `json:"data"`
}

type SpParams struct {
	BaseParams      BaseParams `json:"baseParams"`
	Query           map[string]interface{}
	PortfolioId     int64
	CampaignId      int64
	CampaignIds     map[string]interface{}
	AdGroupId       int64
	AdGroupIds      []interface{}
	KeywordId       int64
	Asin            string
	Asins           []interface{}
	AdId            int64
	TargetId        int64
	IncludeAncestor bool
	CategoryId      int64
}

type SdParams struct {
	BaseParams       BaseParams `json:"baseParams"`
	Query            map[string]interface{}
	CampaignId       int64
	AdGroupId        int64
	TargetId         int64
	Products         []interface{}
	TargetingClauses []interface{}
	AdId             int64
}

type SbParams struct {
	BaseParams      BaseParams `json:"baseParams"`
	PageUrl         string
	MediaId         string
	Query           map[string]interface{}
	PortfolioId     int64
	CampaignId      int64
	CampaignIds     map[string]interface{}
	AdGroupId       int64
	AdGroupIds      []interface{}
	KeywordId       int64
	Asin            string
	Asins           []string
	AdId            int64
	TargetId        int64
	IncludeAncestor bool
	CategoryId      int64
}

type OrderParams struct {
	BaseParams
	OrderId       string
	AdvertiserIds []string
	Status        string
}

type LineItemParams struct {
	BaseParams
	OrderIds   []string
	LineItemId string
	Status     string
}

type CreativeParams struct {
	BaseParams
	AdvertiserIds []string
}

type LineItemCreativeAssociation struct {
	BaseParams
	LineItemIds  []string
	AdvertiserId int64
}

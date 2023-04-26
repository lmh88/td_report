package fead

import "github.com/gogf/gf/frame/g"

type Subscription struct {
	SubscriptionId string `json:"subscriptionId"`
	DestinationArn string `json:"destinationArn"`
	DataSetId      string `json:"dataSetId"`
	Status         string `json:"status"`
	Notes          string `json:"notes"`
}


type SubscriptionList struct {
	Subscriptions []Subscription `json:"Subscriptions"`
	NextToken     string         `json:"nextToken"`
}

type SubscriptionId struct {
	SubscriptionId string `json:"subscriptionId"`
}

const (
	SpTraffic    = "sp-traffic"
	SpConversion = "sp-conversion"
	BudgetUsage  = "budget-usage"
)

var SqsMap = map[string]string{
	SpTraffic:    "arn:aws:sqs:us-east-1:207485092024:sp-traffic",
	SpConversion: "arn:aws:sqs:us-east-1:207485092024:sp-conversion",
	BudgetUsage:  "arn:aws:sqs:us-east-1:207485092024:budget-usage",
}


var NumMap = map[string]int{
	SpTraffic:    g.Cfg().GetInt("sqs.goroutine.sp-traffic", 5),
	SpConversion: g.Cfg().GetInt("sqs.goroutine.sp-conversion", 1),
	BudgetUsage:  g.Cfg().GetInt("sqs.goroutine.budget-usage", 1),
}

func getGoNum(name string) int {
	if res, ok := NumMap[name]; ok {
		return res
	}
	return 1
}

var KafkaMap = map[string]string{
	SpConversion: "amazon_streamapi_sp_ods_conversion",
	SpTraffic:    "amazon_streamapi_sp_ods_traffic",
	BudgetUsage:  "amazon_streamapi_sp_ods_budget_usage",
}

func getKafkaTopic(name string) string {
	if res, ok := KafkaMap[name]; ok {
		return res
	}
	return "amazon_streamapi_sp_ods_traffic_qa"
}


const StatusActive = "ACTIVE"     //可用
const StatusArchived = "ARCHIVED" //归档
const StatusProvisioning = "PROVISIONING" //配置中
const StatusPendingConfirmation  = "PENDING_CONFIRMATION" //等待确认
const StatusFailedConfirmation = "FAILED_CONFIRMATION" //确认失败
const StatusSuspended = "SUSPENDED" //暂停


func isMySqs(arn string) bool {
	for _, item := range SqsMap {
		if item == arn {
			return true
		}
	}
	return false
}

func isValidStatus(status string) bool {
	list := map[string]bool{
		StatusActive: true,
		StatusProvisioning: true,
		StatusPendingConfirmation: true,
	}
	_, ok := list[status]
	return ok
}

func middleStatus(status string) bool {
	list := map[string]bool{
		StatusProvisioning: true,
		StatusPendingConfirmation: true,
	}
	_, ok := list[status]
	return ok
}

func isInValidStatus(status string) bool {
	list := map[string]bool{
		StatusArchived: true,
		StatusFailedConfirmation: true,
		StatusSuspended: true,
	}
	_, ok := list[status]
	return ok
}


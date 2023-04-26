package bean

import "time"

type SubscriptNode struct {
	CreatedDate    time.Time `json:"createdDate"`
	DataSetId      string    `json:"dataSetId"`
	DestinationArn string    `json:"destinationArn"`
	Notes          string    `json:"notes"`
	Status         string    `json:"status"`
	SubscriptionId string    `json:"subscriptionId"`
	UpdatedDate    time.Time `json:"updatedDate"`
}

type ListAllStream struct {
	Subscriptions []SubscriptNode `json:"subscriptions"`
}

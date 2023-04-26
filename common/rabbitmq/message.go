package rabbitmq

type MessageBody struct {
	MessageId        string
	MessageBody      string
}

type TopicMessageRequest struct {
	Queuename string      `json:"queuename,omitempty"`
	Body  MessageBody `json:"message_body"`
}

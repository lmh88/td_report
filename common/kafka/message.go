package kafka

type MessageBody struct {
	MessageId        string
	MessageBody      string
}

type TopicMessageRequest struct {
	Topic string      `json:"topic,omitempty"`
	Body  MessageBody `json:"message_body"`
}

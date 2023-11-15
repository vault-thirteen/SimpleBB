package models

type Message struct {
	Id uint `json:"id"`

	// ID of a thread containing this message.
	ThreadId uint `json:"threadId"`

	Text string `json:"text"`

	// Message meta-data.
	EventData
}

func NewMessage() (msg *Message) {
	return &Message{
		EventData: EventData{
			Creator: &EventParameters{},
			Editor:  &OptionalEventParameters{},
		},
	}
}

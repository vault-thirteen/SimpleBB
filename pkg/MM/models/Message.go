package models

type Message struct {
	// Identifier of this message.
	Id uint `json:"id"`

	// Identifier of a thread containing this message.
	ThreadId uint `json:"threadId"`

	// Textual information of this message.
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

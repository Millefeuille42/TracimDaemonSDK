package TracimDaemonSDK

import "time"

type TLMEvent struct {
	EventId   int         `json:"event_id"`
	EventType string      `json:"event_type"`
	Read      interface{} `json:"read"`
	Created   time.Time   `json:"created"`
	Fields    interface{} `json:"fields"`
}

type Event struct {
	Data       []byte
	DataParsed TLMEvent
	Size       int
}

type EventHandler func(*TracimDaemonClient, *Event)

const (
	EventTypeGeneric = "custom_message"
)

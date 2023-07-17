package TracimDaemonSDK

import "time"

// TLMEvent is the struct that represents the data sent by Tracim (see tracim TLM documentation)
type TLMEvent struct {
	EventId   int         `json:"event_id"`
	EventType string      `json:"event_type"`
	Read      interface{} `json:"read"`
	Created   time.Time   `json:"created"`
	Fields    interface{} `json:"fields"`
}

// Event is a wrapper for the TLMEvent struct, containing the raw data, the parsed data and the size of the raw data
type Event struct {
	Data       []byte
	DataParsed TLMEvent
	Size       int
}

// EventHandler is the function definition for the event handlers
// it takes a TracimDaemonClient and an Event as parameters
type EventHandler func(*TracimDaemonClient, *Event)

const (
	// EventTypeGeneric is the event type for generic events (every message sent by Tracim)
	EventTypeGeneric = "custom_message"
)

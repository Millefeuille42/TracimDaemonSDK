package TracimDaemonSDK

import (
	"fmt"
	"time"
)

// TLMEvent is the struct that represents the data sent by Tracim (see tracim TLM documentation)
type TLMEvent struct {
	EventId   int         `json:"event_id"`
	EventType string      `json:"event_type"`
	Read      interface{} `json:"read"`
	Created   time.Time   `json:"created"`
	Fields    interface{} `json:"fields"`
}

// EventHandler is the function definition for the event handlers
// it takes a TracimDaemonClient and an DaemonEvent as parameters
type EventHandler func(*TracimDaemonClient, *DaemonEvent)

const (
	// EventTypeGeneric is the event type for generic events (every DaemonEvent)
	EventTypeGeneric = "custom_message"
)

func defaultPingHandler(c *TracimDaemonClient, e *DaemonEvent) {
	err := c.SendDaemonEvent(&DaemonEvent{
		Path: c.ClientSocketPath,
		Type: DaemonPong,
		Data: nil,
	}, e.Path)

	if err != nil {
		c.log(err)
		return
	}

	c.log(fmt.Sprintf("SOCKET: SEND: %s -> %s", DaemonPong, e.Path))
}

func defaultAccountInfoHandler(c *TracimDaemonClient, e *DaemonEvent) {
	if e.Path != c.MasterSocketPath {
		return
	}

	err := ParseDaemonData(e, &DaemonAccountInfoData{})
	if err != nil {
		c.log(err)
		return
	}

	c.UserID = e.Data.(*DaemonAccountInfoData).UserId

	c.log(fmt.Sprintf("SOCKET: RECV: %s -> %s", e.Type, e.Path))
}

package TracimDaemonSDK

import (
	"encoding/json"
	"net"
)

// DaemonSubscriptionEvent is the event sent by the client to the master when it wants to subscribe or unsubscribe
type DaemonSubscriptionEvent struct {
	// Path is the path to the client socket
	Path string `json:"path"`
	// Action is the action to perform (add or delete)
	Action string `json:"action"`
}

const (
	// DaemonSubscriptionActionDelete is the action to send to delete a client
	DaemonSubscriptionActionDelete = "client_delete"
	// DaemonSubscriptionActionAdd is the action to send to add a client
	DaemonSubscriptionActionAdd = "client_add"
)

func (c *TracimDaemonClient) sendDaemonSubscriptionEvent(eventType string) error {
	masterSocket, err := net.Dial("unix", c.MasterSocketPath)
	if err != nil {
		return err
	}
	defer masterSocket.Close()

	message := DaemonSubscriptionEvent{
		Path:   c.ClientSocketPath,
		Action: eventType,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = masterSocket.Write(data)
	return err
}

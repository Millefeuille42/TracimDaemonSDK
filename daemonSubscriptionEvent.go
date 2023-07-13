package TracimDaemonSDK

import (
	"encoding/json"
	"net"
)

type DaemonSubscriptionEvent struct {
	Path   string `json:"path"`
	Action string `json:"action"`
}

const (
	DaemonSubscriptionActionDelete = "client_delete"
	DaemonSubscriptionActionAdd    = "client_add"
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

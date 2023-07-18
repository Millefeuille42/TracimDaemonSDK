package TracimDaemonSDK

import (
	"encoding/json"
	"net"
)

// DaemonEvent is the event sent by the client to the master when it wants to subscribe or unsubscribe
type DaemonEvent struct {
	// Path is the path to the client socket
	Path string `json:"path"`
	// Action is the action to perform (add or delete)
	Action string `json:"action"`
	// Additional Data transited in the message
	Data interface{} `json:"data,omitempty"`
}

const (
	// DaemonSubscriptionActionDelete is the action to send to delete a client
	DaemonSubscriptionActionDelete = "daemon_client_delete"
	// DaemonSubscriptionActionAdd is the action to send to add a client
	DaemonSubscriptionActionAdd = "daemon_client_add"
	// DaemonAck is the action sent by the master to acknowledge the previous action
	DaemonAck = "daemon_ack"
	// DaemonPing is the action sent, expecting a DaemonPong response, used to test responsiveness of the other end
	DaemonPing = "daemon_ping"
	// DaemonPong is the action sent, in response of a DaemonPing, used to test responsiveness of the other end
	DaemonPong = "daemon_pong"
	// DaemonAccountInfo is the action sent to relay the info of the current logged-in user
	DaemonAccountInfo = "daemon_account_info"
	// DaemonTracimEvent is the action sent to relay tracim events
	DaemonTracimEvent = "daemon_tracim_event"
)

func (c *TracimDaemonClient) SendDaemonEvent(event *DaemonEvent, socket string) error {
	masterSocket, err := net.Dial("unix", socket)
	if err != nil {
		return err
	}
	defer masterSocket.Close()

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = masterSocket.Write(data)
	return err
}

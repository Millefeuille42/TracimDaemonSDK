package TracimDaemonSDK

import (
	"encoding/json"
	"errors"
	"net"
	"reflect"
)

// DaemonEvent is the event sent by the client to the master when it wants to subscribe or unsubscribe
type DaemonEvent struct {
	// Path is the path to the client socket
	Path string `json:"path"`
	// Type is the DaemonEvent type
	Type string `json:"type"`
	// Additional Data transited in the message
	Data interface{} `json:"data,omitempty"`
}

// SendDaemonEvent sends a daemon event to the provided socket
func (c *TracimDaemonClient) SendDaemonEvent(event *DaemonEvent, socket string) error {
	targetSocket, err := net.Dial("unix", socket)
	if err != nil {
		return err
	}
	defer targetSocket.Close()

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = targetSocket.Write(data)
	return err
}

/* --- Daemon Event Types --- */

/* Client to Daemon events */

// DaemonClientAdd is the event sent by a client to the daemon to register as a client
const DaemonClientAdd = "daemon_client_add"

// DaemonClientDelete is the event sent by a client to the daemon to unregister as a client
const DaemonClientDelete = "daemon_client_delete"

// DaemonGetClients is the event sent by a client to the daemon to get a list of the currently active clients
const DaemonGetClients = "daemon_get_clients"

// DaemonGetAccountInfo is the event sent by a client to the daemon to get info about the logged-in user
const DaemonGetAccountInfo = "daemon_get_account_info"

// DaemonDoRequest is the event sent by a client to the daemon to make a request to the tracim API
const DaemonDoRequest = "daemon_do_request"

/* Any to Any events */

// DaemonAck is the event sent by any to acknowledge the previous action
const DaemonAck = "daemon_ack"

// DaemonPing is the event sent by any, expecting a DaemonPong response, used to test responsiveness of the other end
const DaemonPing = "daemon_ping"

// DaemonPong is the event sent by any, in response of a DaemonPing, used to test responsiveness of the other end
const DaemonPong = "daemon_pong"

/* Daemon to Client events */

// DaemonRequestResult is the event sent by the daemon in response to the DaemonDoRequest event
const DaemonRequestResult = "daemon_request_result"

// DaemonAccountInfo is the event sent by the daemon in response to a DaemonClientAdd or DaemonGetAccountInfo event, it contains info about the logged-in user
const DaemonAccountInfo = "daemon_account_info"

// DaemonClients is the event sent by the daemon in response to a DaemonGetClients event, it contains info about the currently registered clients
const DaemonClients = "daemon_clients"

/* Daemon to all Clients events */

// DaemonTracimEvent is the event sent by the daemon to every client, to broadcast tracim events
const DaemonTracimEvent = "daemon_tracim_event"

// DaemonClientAdded is the event sent by the daemon to every client when a new client registers
const DaemonClientAdded = "daemon_client_added"

// DaemonClientDeleted is the event sent by the daemon to every client when a client is unregistered
const DaemonClientDeleted = "daemon_client_deleted"

/* --- DaemonEvent Data Types --- */

// ParseDaemonData converts DaemonEvent Data field into a reference struct pointer
// unless an error is returned, DaemonEvent.Data is safely cast-able as a *reference type
func ParseDaemonData(e *DaemonEvent, reference interface{}) error {
	refV := reflect.ValueOf(reference)
	if refV.Kind() != reflect.Pointer || refV.IsNil() {
		return errors.New("reference is not a pointer")
	}

	var msi map[string]interface{}
	if reflect.ValueOf(e.Data).IsNil() || reflect.TypeOf(e.Data) != reflect.TypeOf(msi) {
		return errors.New("data is nil or of an invalid type")
	}

	srcAsMap := e.Data.(map[string]interface{})
	structValue := reflect.ValueOf(reference).Elem()
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)
		if tagVal, ok := fieldType.Tag.Lookup("json"); ok {
			if val, ok := srcAsMap[tagVal]; ok {
				if reflect.TypeOf(val) == field.Type() {
					field.Set(reflect.ValueOf(val))
				}
			}
		}
	}

	e.Data = reference

	return nil
}

/* Generic types */

// DaemonClientData is an intermediary struct for data concerning a daemon client
type DaemonClientData struct {
	// Path to the socket of the client
	Path string `json:"path"`
	// Pid of the client
	Pid string `json:"pid"`
}

/* Client to Daemon types */

// DaemonClientAddData is the data sent for DaemonClientAdd events
type DaemonClientAddData DaemonClientData

// DaemonClientDeleteData is the data sent for DaemonClientDelete events
type DaemonClientDeleteData DaemonClientData

// DaemonDoRequestData is the data sent for DaemonDoRequest events
type DaemonDoRequestData struct {
	// Method of the request
	Method string
	// Endpoint of the request (appended to <protocol>://<tracim_host>:<port>/api)
	Endpoint string
	// Body of the request
	Body []byte
}

/* Any to Any types */

// DaemonAckData is the data sent for DaemonAck events
type DaemonAckData struct {
	// Type is the type of the DaemonEvent being acknowledged
	Type DaemonEvent
}

/* Daemon to Client types */

// DaemonRequestResultData is the data sent for DaemonRequestResult events
type DaemonRequestResultData struct {
	// Request contains the originating DaemonDoRequestData data,
	Request DaemonDoRequestData
	// StatusCode corresponds to the eponymous field in the http.Response struct
	StatusCode int
	// Status corresponds to the eponymous field in the http.Response struct
	Status string
	// Data is the []byte data found in the http.Response associated with the previous request
	Data []byte `json:"data"`
}

// DaemonAccountInfoData is the data sent for DaemonAccountInfo events
type DaemonAccountInfoData struct {
	// UserId is the tracim user ID used by the daemon
	UserId string `json:"user_id"`
}

// DaemonClientsData is the data sent for DaemonClients events
type DaemonClientsData []DaemonClientData

/* Daemon to all Clients types */

// DaemonTracimEventData is the data sent for DaemonTracimEvent events
type DaemonTracimEventData struct {
	TLMEvent
}

// DaemonClientAddedData is the data sent for DaemonClientAdded events
type DaemonClientAddedData DaemonClientData

// DaemonClientDeletedData is the data sent for DaemonClientDeleted events
type DaemonClientDeletedData DaemonClientData

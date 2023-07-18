# TracimDaemonSDK

An SDK for the [TracimDaemon](https://github.com/Millefeuille42/TracimDaemon) project

## Usage

Import the package

```go
import "github.com/Millefeuille42/TracimDaemonSDK"
```

Create a new TracimDaemon client

```go
client := TracimDaemonSDK.NewClient(TracimDaemonSDK.Config{
	MasterSocketPath: "path/to/master/socket",
	ClientSocketPath: "path/to/client/socket",
})
```

Set up handlers for signals (for proper shutdown)

```go
client.HandleCloseOnSig(os.Interrupt)
client.HandleCloseOnSig(os.Kill)
```

Create and listen to the plugin socket

```go
err := client.CreateClientSocket()
if err != nil {
	log.Fatal(err)
	return
}
defer client.ClientSocket.Close()
```

Set up various handlers

```go
client.RegisterHandler(TracimDaemonSDK.EventTypeGeneric, genericHandler)
```

With `genericHandler` being a function with the following signature:

```go
func genericHandler(c *TracimDaemonSDK.TracimDaemonClient, e *TracimDaemonSDK.Event) {
    log.Printf("%s RECV: %s\n", c.Config.ClientSocketPath, e.DataParsed.EventType)
}
```

Register the plugin to the master daemon

```go
err = client.RegisterToMaster()
if err != nil {
    log.Fatal(err)
    return
}
```

Start the client

```go
client.ListenToEvents()
```

The "minimal" client code is as follows:

```go
package main

import (
	"github.com/Millefeuille42/TracimDaemonSDK"
	"log"
	"os"
)

func genericHandler(c *TracimDaemonSDK.TracimDaemonClient, e *TracimDaemonSDK.Event) {
	log.Printf("%s RECV: %s\n", c.Config.ClientSocketPath, e.DataParsed.EventType)
}

func main() {
	client := TracimDaemonSDK.NewClient(TracimDaemonSDK.Config{
		MasterSocketPath: os.Getenv("TRACIM_MINICLIENT_MASTER_SOCKET_PATH"),
		ClientSocketPath: os.Getenv("TRACIM_MINICLIENT_CLIENT_SOCKET_PATH"),
	})

	client.HandleCloseOnSig(os.Interrupt)
	client.HandleCloseOnSig(os.Kill)
	err := client.CreateClientSocket()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.ClientSocket.Close()

	client.RegisterHandler(TracimDaemonSDK.EventTypeGeneric, genericHandler)
	err = client.RegisterToMaster()
	if err != nil {
		log.Fatal(err)
		return
	}

	client.ListenToEvents()
}
```

## Definitions

### TLMEvent

TLMEvent is the struct that represents the data sent by Tracim (see tracim TLM documentation)

```go
type TLMEvent struct {
	EventId   int         `json:"event_id"`
	EventType string      `json:"event_type"`
	Read      interface{} `json:"read"`
	Created   time.Time   `json:"created"`
	Fields    interface{} `json:"fields"`
}
```

### DaemonEvent

DaemonEvent is the event format used to communicate between apps

```go
type DaemonEvent struct {
	Path   string `json:"path"`
	Action string `json:"action"`
	Data interface{} `json:"data,omitempty"`
}
```

- The `Data` field can contain additional information of any format
- The `Path` field is the path to the plugin socket (as defined in the config)
- The `Action` field is any of the following:

```go
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
````

For now, only the `DaemonAccountInfo` and `DaemonTracimEvent` should contain additional data.

### EventHandler

EventHandler is the function definition for the event handlers.
It takes a `TracimDaemonClient` and a `DaemonEvent` as parameters

```go
type EventHandler func(*TracimDaemonClient, *DaemonEvent)
```

By default, handlers for `DaemonAccountInfo` and `DaemonPing` are already defined, it is possible to override them.

### Event types

Event types are defined by tracim. It is also possible to set handlers for every `DaemonEvent` type.
There also is events defined by the SDK, for convenience.

```go
// EventTypeGeneric is the event type for generic events (every DaemonEvent)
EventTypeGeneric = "custom_message"
```

## Protocol (for developers of another language)

### Registering / unregistering a plugin

When registering / unregistering a plugin a basic message must be sent on the master socket.

The message must be a JSON object with the following structure:

```json
{
    "action": "client_add",
    "path": "/path/to/plugin/socket",
    "data": {}
}
```

The `action` field can be one of the previously demonstrated types.

To register a plugin, the plugin must send the message to the master socket, with the `action` field set to `client_add`.
To unregister a plugin, the plugin must send the message to the master socket, with the `action` field set to `client_delete`.

### Receiving events

Once registered, the plugin will receive events from the master socket.

The events are JSON objects stored in the `data` field of a `DaemonEvent` with the following structure:

```json
{
    "event_id": 1,
    "event_type": "custom_message",
    "read": false,
    "created": "2020-01-01T00:00:00Z",
    "fields": {}
}
```

With the `fields` field being a JSON object containing the event data.


### Ack and Keep-Alive

#### Ack

The master daemon will send a `DaemonAck` upon receiving any events not expecting a response, otherwise the 
expected response is sent. As for now, a `DaemonPong` for a `DaemonPing` and a `DaemonAccountInfo` for a `DaemonSubscriptionActionAdd`

The master daemon expects no `DaemonAck` on its messages.

#### Keep-Alive

The master daemon will periodically (once every minute) send a `DaemonPing` event, clients have a minute to respond with `DaemonPong`,
If not, it will unregister un-responding clients at the next ping.

It is possible to test the master daemon responsiveness by sending it `DaemonPing` events. It will respond with a `DaemonPong` as soon as possible.

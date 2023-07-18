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

### Event

Event is a wrapper for the TLMEvent struct, containing the raw data, the parsed data and the size of the raw data.

It is what's sent from the master to the handlers.
```go
type Event struct {
    Data       []byte
    DataParsed TLMEvent
    Size       int
}
```

### EventHandler

EventHandler is the function definition for the event handlers.
It takes a TracimDaemonClient and an Event as parameters

```go
type EventHandler func(*TracimDaemonClient, *Event)
```

### Event types

Event types are defined by tracim. However, the SDK defines some constants for convenience.

```go
// EventTypeGeneric is the event type for generic events (every message sent by Tracim)
EventTypeGeneric = "custom_message"
```

### DaemonSubscriptionEvent

DaemonSubscriptionEvent is the event sent by the client to the master when it wants to subscribe or unsubscribe

```go
type DaemonSubscriptionEvent struct {
Path   string `json:"path"`
Action string `json:"action"`
}
```

The `Path` field is the path to the plugin socket (as defined in the config)
The `Action` field is any of the following:

```go
// DaemonSubscriptionActionDelete is the action to send to delete a client
DaemonSubscriptionActionDelete = "client_delete"
// DaemonSubscriptionActionAdd is the action to send to add a client
DaemonSubscriptionActionAdd    = "client_add"
````

## Protocol (for developers of another language)

### Registering / unregistering a plugin

When registering / unregistering a plugin a basic message must be sent on the master socket.

The message must be a JSON object with the following structure:

```json
{
    "action": "client_add",
    "path": "/path/to/plugin/socket"
}
```

The `action` field can be either `client_add` or `client_delete`

To register a plugin, the plugin must send the message to the master socket, with the `action` field set to `client_add`.
To unregister a plugin, the plugin must send the message to the master socket, with the `action` field set to `client_delete`.

### Receiving events

Once registered, the plugin will receive events from the master socket.

The events are JSON objects with the following structure:

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

### Note

For the moment there is no keep_alive mechanism. The plugin MUST unregister itself when it is stopped,
to avoid multiple messages on the same socket.

For now, no acknowledgement is sent by the master daemon when a plugin is registered / unregistered.
It is planned to communicate user information upon registration, but it is not implemented yet.

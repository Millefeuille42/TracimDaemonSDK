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
	MasterSocketPath: "path/to/daemon/socket",
	ClientSocketPath: "path/to/client/socket",
})
```

Set up handlers for signals (for proper shutdown)

```go
client.HandleCloseOnSig(os.Interrupt)
client.HandleCloseOnSig(os.Kill)
```

Create and listen to the client socket

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

Register the client to the daemon daemon

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
	Type string `json:"type"`
	Data interface{} `json:"data,omitempty"`
}
```

- The `Data` field can contain additional information of any format
- The `Path` field is the path to the client socket (as defined in the config)
- The `Type` field is any of the Daemon* constants defined in `daemonEvents.go`

A `Type` is expected to contain additional data if there is a `<eventType>Data` struct defined in `daemonEvents.go`.

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

### Communication

Client and daemons communicate with the `DaemonEvent` format, i.e. JSON data following this format:

```json
{
    "type": "event_type",
    "path": "/path/to/client/socket",
    "data": {}
}
```

(See above `DaemonEvent` section for details about each field)

### Registering / unregistering a client

When registering / unregistering a client a `DaemonEvent` must be sent on the daemon socket.

```json
{
    "type": "client_add",
    "path": "/path/to/client/socket",
    "data": {
      "path": "/path/to/client/socket",
      "pid": 999
    }
}
```

To register a client, the client must send the message to the daemon socket, with the `type` field set to `client_add`.
To unregister a client, the client must send the message to the daemon socket, with the `type` field set to `client_delete`.

In both, additional info, defined as follows, is required:

```json
{
  "path": "/path/to/client/socket",
  "pid": 999
}
```

With `pid` being the PID of the client process.

### Receiving events

Once registered, the client will receive `DaemonEvent`s from the daemon.

(See above `DaemonEvent` section for details about types and data)

### Ack and Keep-Alive

#### Ack

The daemon will send a `DaemonAck` upon receiving any managed events not expecting a response, otherwise the 
expected response is sent. As for now, a `DaemonPong` for a `DaemonPing` and a `DaemonAccountInfo` for a `DaemonClientAdd`

The daemon expects no `DaemonAck` on its messages.

#### Keep-Alive[

The daemon will periodically (once every minute) send a `DaemonPing` event, clients have a minute to respond with `DaemonPong`,
If not, it will unregister un-responding clients at the next ping.

It is possible to test the daemon responsiveness by sending it `DaemonPing` events. It will respond with a `DaemonPong` as soon as possible.
]()

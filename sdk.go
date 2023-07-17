package TracimDaemonSDK

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
)

// Config is the configuration of the TracimDaemonClient
type Config struct {
	// MasterSocketPath is the path to the master socket
	MasterSocketPath string
	// ClientSocketPath is the path to the client socket
	ClientSocketPath string
}

// TracimDaemonClient is a client for the TracimDaemon
type TracimDaemonClient struct {
	Config
	// ClientSocket is the listener for the client socket
	ClientSocket net.Listener
	// EventHandlers is the map of event handlers
	EventHandlers map[string]EventHandler
}

func (c *TracimDaemonClient) callHandler(eventType string, eventData *Event) {
	if _, ok := c.EventHandlers[eventType]; ok {
		c.EventHandlers[eventType](c, eventData)
	}
}

// RegisterToMaster registers the client to the master
func (c *TracimDaemonClient) RegisterToMaster() error {
	return c.sendDaemonSubscriptionEvent(DaemonSubscriptionActionAdd)
}

// UnregisterFromMaster unregisters the client from the master
func (c *TracimDaemonClient) UnregisterFromMaster() error {
	return c.sendDaemonSubscriptionEvent(DaemonSubscriptionActionDelete)
}

// CreateClientSocket creates the client socket and attaches a listener to it
func (c *TracimDaemonClient) CreateClientSocket() error {
	var err error
	c.ClientSocket, err = net.Listen("unix", c.ClientSocketPath)
	return err
}

// ListenToEvents listens to events sent by the daemon on the client socket
func (c *TracimDaemonClient) ListenToEvents() {
	for {
		conn, err := c.ClientSocket.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			buf := make([]byte, 4096)

			n, err := conn.Read(buf)
			if err != nil {
				log.Print(err)
				return
			}
			eventData := Event{
				Data: buf[:n],
				Size: n,
			}
			err = json.Unmarshal(eventData.Data, &eventData.DataParsed)
			if err != nil {
				log.Print(err)
				return
			}

			c.callHandler(EventTypeGeneric, &eventData)
			c.callHandler(eventData.DataParsed.EventType, &eventData)
		}(conn)
	}
}

// RegisterHandler registers an event handler for a specific event type
func (c *TracimDaemonClient) RegisterHandler(eventType string, handler EventHandler) {
	c.EventHandlers[eventType] = handler
}

// HandleCloseOnSig handles the closing of the client on a specific signal
func (c *TracimDaemonClient) HandleCloseOnSig(sig os.Signal) {
	cc := make(chan os.Signal, 1)
	signal.Notify(cc, sig)
	go func() {
		<-cc
		err := c.UnregisterFromMaster()
		if err != nil {
			log.Print(err)
		}
		_ = os.Remove(c.ClientSocketPath)
		os.Exit(1)
	}()
}

// NewClient creates a new TracimDaemonClient
func NewClient(conf Config) (client *TracimDaemonClient) {
	client = &TracimDaemonClient{
		Config:        conf,
		EventHandlers: make(map[string]EventHandler),
	}

	return client
}

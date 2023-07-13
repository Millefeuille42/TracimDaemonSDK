package TracimDaemonSDK

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
)

type Config struct {
	MasterSocketPath string
	ClientSocketPath string
}

type TracimDaemonClient struct {
	Config
	ClientSocket  net.Listener
	EventHandlers map[string]EventHandler
}

func (c *TracimDaemonClient) callHandler(eventType string, eventData *Event) {
	if _, ok := c.EventHandlers[eventType]; ok {
		c.EventHandlers[eventType](c, eventData)
	}
}

func (c *TracimDaemonClient) RegisterToMaster() error {
	return c.sendDaemonSubscriptionEvent(DaemonSubscriptionActionAdd)
}

func (c *TracimDaemonClient) UnregisterFromMaster() error {
	return c.sendDaemonSubscriptionEvent(DaemonSubscriptionActionDelete)
}

func (c *TracimDaemonClient) CreateClientSocket() error {
	var err error
	c.ClientSocket, err = net.Listen("unix", c.ClientSocketPath)
	return err
}

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

func (c *TracimDaemonClient) RegisterHandler(eventType string, handler EventHandler) {
	c.EventHandlers[eventType] = handler
}

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

func NewClient(conf Config) (client *TracimDaemonClient) {
	client = &TracimDaemonClient{
		Config: conf,
	}

	return client
}

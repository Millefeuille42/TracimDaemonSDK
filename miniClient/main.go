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
		MasterSocketPath: "/tmp/tracim_master.sock",
		ClientSocketPath: "/tmp/tracim_mini_client.sock",
	})

	client.HandleCloseOnSig(os.Interrupt)
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

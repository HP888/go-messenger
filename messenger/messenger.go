package messenger

import (
	"fmt"
	"gomessenger/config"
	"net"
)

const PREFIX = "[GoMessenger] "

func Listen(config config.ServerConfig) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		panic(err)
	}

	acceptConnections(listener.(*net.TCPListener))
}

func acceptConnections(listener *net.TCPListener) {
	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			PrintErr(err)
			continue
		}

		handler := NewHandler(connection)
		Println("Client has connected to the server! [" + handler.Address + "]")
		go handler.Handle()
	}
}

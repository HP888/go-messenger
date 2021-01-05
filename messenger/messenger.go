package messenger

import (
	"fmt"
	"gomessenger/config"
	"log"
	"net"
)

const PREFIX = "[GoMessenger] "

func Listen(config config.ServerConfig) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		panic(err)
	}

	acceptConnection(listener.(*net.TCPListener))
}

func acceptConnection(listener *net.TCPListener) {
	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
			continue
		}

		handler := NewHandler(connection)
		Println("Client has connected to the server! [" + handler.Address + "]")
		go handler.Handle()
	}
}

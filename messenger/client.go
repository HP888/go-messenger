package messenger

import (
	"container/list"
	"fmt"
	"net"
	"sync"
)

type Client struct {
	PacketQueue chan Packet
	Connection  *net.TCPConn
	Element     *list.Element
	Address     string
	Mutex       sync.Mutex
}

var handlers = list.New()

func NewHandler(connection *net.TCPConn) *Client {
	client := new(Client)
	client.Address = connection.RemoteAddr().(*net.TCPAddr).IP.String()
	client.PacketQueue = make(chan Packet)
	client.Connection = connection

	client.Element = handlers.PushBack(client)

	return client
}

func (client *Client) Handle() {
	go client.ReadPackets()
	err := client.WritePackets()

	if client.Connection != nil {
		_ = client.Connection.Close()
		client.PacketQueue = nil
		client.Connection = nil
	}

	Println("Client has disconnected from the server! [" + client.Address + "]")
	if err != nil {
		PrintErr(err)
	}

	handlers.Remove(client.Element)
}

func (client *Client) ReadPackets() {
	for {
		reader := client.Connection
		bytes, err := ReadBytes(reader, MaxPacketLength)

		if err != nil {
			client.PacketQueue <- Packet{Data: []byte(""), Error: err}
			break
		}

		client.PacketQueue <- Packet{Data: bytes, Error: nil}
	}
}

func (client *Client) WritePackets() error {
	var err error
	packetChannel := client.PacketQueue

	for {
		packet := <-packetChannel
		if packet.Error != nil {
			err = fmt.Errorf("error in handler's WritePackets(): %w", packet.Error)
			break
		}

		bytes := packet.Data

		for element := handlers.Front(); element != nil; element = element.Next() {
			targetHandler := element.Value.(*Client)

			targetHandler.Mutex.Lock()

			err = WriteBytes(targetHandler.Connection, bytes)
			if err != nil {
				break
			}

			targetHandler.Mutex.Unlock()
		}
	}

	close(packetChannel)
	return err
}

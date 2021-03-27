package messenger

import (
	"container/list"
	"fmt"
	"net"
	"sync"
)

type Handler struct {
	PacketQueue chan Packet
	Connection  *net.TCPConn
	Element     *list.Element
	Address     string
	Mutex       sync.Mutex
}

var handlers = list.New()

func NewHandler(connection *net.TCPConn) *Handler {
	handler := new(Handler)
	handler.Address = connection.RemoteAddr().(*net.TCPAddr).IP.String()
	handler.PacketQueue = make(chan Packet)
	handler.Connection = connection

	handler.Element = handlers.PushBack(handler)

	return handler
}

func (handler *Handler) Handle() {
	go handler.ReadPackets()
	err := handler.WritePackets()

	if handler.Connection != nil {
		_ = handler.Connection.Close()
		handler.PacketQueue = nil
		handler.Connection = nil
	}

	Println("Client has disconnected from the server! [" + handler.Address + "]")
	if err != nil {
		PrintErr(err)
	}

	handlers.Remove(handler.Element)
}

func (handler *Handler) ReadPackets() {
	for {
		reader := handler.Connection
		bytes, err := ReadBytes(reader, MaxPacketLength)

		if err != nil {
			handler.PacketQueue <- Packet{Data: []byte(""), Error: err}
			break
		}

		handler.PacketQueue <- Packet{Data: bytes, Error: nil}
	}
}

func (handler *Handler) WritePackets() error {
	var err error
	packetChannel := handler.PacketQueue

	for {
		packet := <-packetChannel
		if packet.Error != nil {
			err = fmt.Errorf("error in handler's WritePackets(): %w", packet.Error)
			break
		}

		bytes := packet.Data

		for element := handlers.Front(); element != nil; element = element.Next() {
			targetHandler := element.Value.(*Handler)

			targetHandler.Mutex.Lock()
			_ = WriteBytes(targetHandler.Connection, bytes)
			targetHandler.Mutex.Unlock()
		}
	}

	close(packetChannel)
	return err
}

package messenger

import (
	"container/list"
	"fmt"
	"net"
)

type Handler struct {
	PacketQueue chan []byte
	Connection  *net.TCPConn
	Element     *list.Element
	Address     string
}

var handlers = list.New()

func NewHandler(connection *net.TCPConn) *Handler {
	handler := new(Handler)
	handler.Element = handlers.PushBack(handler)

	handler.Address = connection.RemoteAddr().(*net.TCPAddr).IP.String()
	handler.PacketQueue = make(chan []byte)
	handler.Connection = connection

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
		bytes, err := ReadBytes(reader, 32767)

		if err != nil {
			handler.PacketQueue <- nil
			break
		}

		handler.PacketQueue <- bytes
		// Println("Readed: " + string(bytes))
	}
}

func (handler *Handler) WritePackets() error {
	var err error
	packetChannel := handler.PacketQueue

	for {
		packet := <-packetChannel
		if packet == nil {
			err = fmt.Errorf("error in handler's WritePackets(): Connection closed")
			break
		}

		for element := handlers.Front(); element != nil; element = element.Next() {
			targetHandler := element.Value.(*Handler)
			_ = WriteBytes(targetHandler.Connection, packet)
		}
	}

	close(packetChannel)
	return err
}

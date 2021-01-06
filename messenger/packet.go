package messenger

type Packet struct {
	Data  []byte
	Error error
}

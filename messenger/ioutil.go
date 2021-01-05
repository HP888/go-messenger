package messenger

import (
	"encoding/binary"
	"fmt"
	"io"
)

func ReadInt(reader io.Reader) (n int32, err error) {
	err = binary.Read(reader, binary.BigEndian, &n)
	return
}

func ReadBytes(reader io.Reader, maxLength int32) ([]byte, error) {
	length, err := ReadInt(reader)
	if err != nil {
		return nil, err
	}

	if length > maxLength {
		return nil, fmt.Errorf("byte array longer than maximum: %d > %d", length, maxLength)
	}

	arr := make([]byte, length)
	_, err = io.ReadFull(reader, arr)

	if err != nil {
		return nil, err
	}

	return arr, err
}

func WriteInt(writer io.Writer, number int32) error {
	return binary.Write(writer, binary.BigEndian, number)
}

func WriteBytes(writer io.Writer, array []byte) error {
	err := WriteInt(writer, int32(len(array)))
	if err != nil {
		return err
	}

	_, err = writer.Write(array)
	return err
}

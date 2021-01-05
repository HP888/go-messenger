package messenger

import "fmt"

func Println(line string) {
	fmt.Println(PREFIX + line)
}

func PrintErr(err error) {
	fmt.Println(err)
}

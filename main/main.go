package main

import (
	"gomessenger/config"
	"gomessenger/messenger"
)

func main() {
	config := config.Load()
	messenger.Listen(config.Server)
}

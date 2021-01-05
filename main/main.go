package main

import (
	"gomessenger/config"
	"gomessenger/messenger"
)

func main() {
	configuration := config.Load()
	messenger.Listen(configuration.Server)
}

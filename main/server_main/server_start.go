package main

import (
	server "IM_system/server"
)

func main() {
	server := server.NewServer("127.0.0.1", 8888)
	server.Start()
}

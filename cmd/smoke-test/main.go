package main

import (
	"io"
	"net"
	"protohackers"
)

func main() {
	config := protohackers.ParseConfig()
	server := protohackers.NewServer(config)
	server.StartTCP(handle)
}

func handle(conn net.Conn) {
	io.Copy(conn, conn)
}

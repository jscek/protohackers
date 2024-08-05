package main

import (
	"io"
	"net"
	"protohackers"
)

func main() {
	protohackers.StartTCPServer(handle)
}

func handle(conn net.Conn) {
	io.Copy(conn, conn)
}

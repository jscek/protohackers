package protohackers

import (
	"flag"
	"fmt"
	"log"
	"net"
)

type UDPConnHandler func(net.PacketConn)

func StartUDPServer(handler UDPConnHandler) {
	port := flag.Int("port", 8080, "port number")
	flag.Parse()

	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("udp server started on :%d\n", *port)
	defer conn.Close()

	handler(conn)
}

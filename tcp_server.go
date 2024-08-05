package protohackers

import (
	"flag"
	"fmt"
	"log"
	"net"
)

type ConnHandler func(net.Conn)

func StartTCPServer(handler ConnHandler) {
	port := flag.Int("port", 8080, "port number")
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tcp server started on :%d\n", *port)
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			log.Println("client connected")

			defer func() {
				log.Println("client disconnected")
				conn.Close()
			}()

			handler(conn)
		}()
	}
}

package protohackers

import (
	"flag"
	"fmt"
	"log"
	"net"
)

type Config struct {
	port int
}

func ParseConfig() *Config {
	config := Config{}
	flag.IntVar(&config.port, "port", 8080, "port number")
	flag.Parse()

	return &config
}

type Server struct {
	config *Config
}

func NewServer(config *Config) Server {
	return Server{config: config}
}

type ConnHandler func(net.Conn)

type UDPConnHandler func(net.PacketConn)

func (server *Server) StartTCP(handler ConnHandler) {
	port := server.config.port

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tcp server started on :%d\n", port)
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

func (server *Server) StartUDP(handler UDPConnHandler) {
	port := server.config.port

	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("udp server started on :%d\n", port)
	defer conn.Close()

	handler(conn)
}

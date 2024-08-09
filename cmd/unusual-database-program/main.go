package main

import (
	"fmt"
	"log"
	"net"
	"protohackers"
	"strings"
)

var (
	db = map[string]string{"version": "KVS 1.0"}
)

func main() {
	config := protohackers.ParseConfig()
	server := protohackers.NewServer(config)
	server.StartUDP(handle)
}

func handle(conn net.PacketConn) {
	for {
		buf := make([]byte, 1000)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		args := strings.SplitN(string(buf[:n]), "=", 2)
		key := args[0]

		if len(args) > 1 {
			if key == "version" {
				continue
			}

			value := args[1]
			db[key] = value
		} else {
			value, exists := db[key]

			if exists {
				_, err := conn.WriteTo([]byte(fmt.Sprintf("%s=%s", key, value)), addr)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

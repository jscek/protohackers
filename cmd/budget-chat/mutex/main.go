package main

import (
	"bufio"
	"fmt"
	"net"
	"protohackers"
	"regexp"
	"sync"
)

type Client struct {
	address string
	name    string
	conn    net.Conn
}

var (
	clients = make(map[string]Client)
	mu      sync.Mutex
)

func main() {
	protohackers.StartTCPServer(handler)
}

func handler(conn net.Conn) {
	conn.Write([]byte("welcome! what's your name?\n"))

	scanner := bufio.NewScanner(conn)

	var client Client

	if scanner.Scan() {
		name := scanner.Text()
		if !isLegalName(name) {
			conn.Write([]byte("sorry, this is not the best name\n"))
			return
		}

		client = Client{name: name, conn: conn, address: conn.RemoteAddr().String()}

		mu.Lock()
		clients[client.address] = client
		mu.Unlock()

		broadcast(&client, fmt.Sprintf("* %s has entered the room\n", client.name))

		defer func() {
			broadcast(&client, fmt.Sprintf("* %s has left the room\n", client.name))
			mu.Lock()
			delete(clients, client.address)
			mu.Unlock()
		}()

		mu.Lock()
		others := []string{}
		for a, c := range clients {
			if a != client.address {
				others = append(others, c.name)
			}
		}
		mu.Unlock()

		conn.Write([]byte(fmt.Sprintf("* The room contains: %v\n", others)))
	}

	for scanner.Scan() {
		message := scanner.Text()
		broadcast(&client, fmt.Sprintf("[%s] %s\n", client.name, message))

		fmt.Printf("[%s] %s\n", client.name, message)
	}
}

var isLegalName = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

func broadcast(client *Client, message string) {
	mu.Lock()
	defer mu.Unlock()
	for a, c := range clients {
		if a != client.address {
			c.conn.Write([]byte(message))
		}
	}
}

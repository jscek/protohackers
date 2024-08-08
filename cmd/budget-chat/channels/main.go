package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
)

type Channels struct {
	join     chan Message
	messages chan Message
	left     chan Message
}

type Message struct {
	from    net.Conn
	content string
}

type User struct {
	name string
	conn net.Conn
}

func (user *User) Send(message string) error {
	_, err := user.conn.Write([]byte(message))
	return err
}

type Room struct {
	users map[net.Conn]User
}

func NewRoom() Room {
	return Room{users: make(map[net.Conn]User)}
}

func (room *Room) AddUser(user User) {
	room.users[user.conn] = user

	var others []string
	for _, u := range room.users {
		if u.conn != user.conn {
			others = append(others, u.name)
		}
	}

	user.Send(fmt.Sprintf("* The room contains: %v\n", others))
}

func (room *Room) Send(from User, message string) error {
	for _, user := range room.users {
		if user.conn != from.conn {
			err := user.Send(fmt.Sprintf("[%s] %s\n", from.name, message))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (room *Room) Notify(message string) {
	for _, user := range room.users {
		user.Send(message)
	}
}

func main() {
	port := flag.Int("port", 8080, "port number")
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tcp server started on :%d\n", *port)
	defer ln.Close()

	channels := Channels{join: make(chan Message), left: make(chan Message), messages: make(chan Message)}
	go handleMessages(channels)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			log.Println("client connected", conn)

			defer func() {
				log.Println("client disconnected", conn)
				conn.Close()
			}()

			handler(conn, channels)
		}()
	}
}

func handler(conn net.Conn, channels Channels) {
	defer func() {
		channels.left <- Message{from: conn, content: ""}
	}()

	conn.Write([]byte("welcome! what's your name?\n"))
	scanner := bufio.NewScanner(conn)

	if scanner.Scan() {
		message := scanner.Text()
		channels.join <- Message{from: conn, content: message}
	}

	for scanner.Scan() {
		message := scanner.Text()
		channels.messages <- Message{from: conn, content: message}
	}
}

func handleMessages(channels Channels) {
	room := NewRoom()

	for {
		select {
		case msg := <-channels.join:
			var isLegalName = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
			name := msg.content
			user := User{name: name, conn: msg.from}

			if !isLegalName(name) {
				user.Send("That's not a very good name\n")
				user.conn.Close()
			}

			room.Notify(fmt.Sprintf("* %s has joined\n", user.name))
			room.AddUser(user)

			fmt.Println("users", room.users)

		case msg := <-channels.messages:
			user := room.users[msg.from]
			room.Send(user, msg.content)

		case msg := <-channels.left:
			user, exists := room.users[msg.from]
			if exists {
				room.Notify(fmt.Sprintf("* %s has left\n", user.name))
				delete(room.users, msg.from)
			}
		}
	}
}

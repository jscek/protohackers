package main

import (
	"bufio"
	"log"
	"net"
	"protohackers"
	"regexp"
	"sync"
)

var boguscoinAddress = regexp.MustCompile(`(^|\s*)(7[a-zA-Z0-9]{25,34})(\s|$)`)

func main() {
	config := protohackers.ParseConfig()
	server := protohackers.NewServer(config)

	server.StartTCP(handle)
}

func handle(downstreamConn net.Conn) {
	upstreamConn, err := net.Dial("tcp", "chat.protohackers.com:16963")
	if err != nil {
		log.Fatal(err)
	}

	close := sync.Once{}

	go pipe(downstreamConn, upstreamConn, &close)
	pipe(upstreamConn, downstreamConn, &close)
}

func pipe(src net.Conn, target net.Conn, close *sync.Once) {
	defer close.Do(func() {
		src.Close()
		target.Close()
	})

	reader := bufio.NewReader(src)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		msg = boguscoinAddress.ReplaceAllString(msg, `${1}7YWHMfk9JZe0LM0g1ZauHuiSxhI${3}`)

		_, err = target.Write([]byte(msg))
		if err != nil {
			log.Fatal(err)
		}
	}
}

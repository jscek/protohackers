package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"protohackers"
)

type Message struct {
	t      byte
	first  int32
	second int32
}

func (msg *Message) UnmarshallBinary(data []byte) error {
	if len(data) != 9 {
		return errors.New("invalid data length")
	}

	if data[0] != 'I' && data[0] != 'Q' {
		return errors.New("invalid type")
	}

	msg.t = data[0]
	msg.first = int32(binary.BigEndian.Uint32(data[1:5]))
	msg.second = int32(binary.BigEndian.Uint32(data[5:]))

	return nil
}

func main() {
	protohackers.StartTCPServer(handle)
}

func handle(conn net.Conn) {

	reader := bufio.NewReader(conn)
	buf := make([]byte, 9)
	prices := make(map[int32]int32)

	for {
		n, _ := io.ReadFull(reader, buf)
		if n == 0 {
			break
		}

		var msg Message
		err := msg.UnmarshallBinary(buf)
		if err != nil {
			log.Println(err)
			break
		}

		if msg.t == 'I' {
			insertPrice(prices, msg.first, msg.second)
		} else if msg.t == 'Q' {
			mean := queryMeanPrice(prices, msg.first, msg.second)
			res := make([]byte, 4)
			binary.BigEndian.PutUint32(res, uint32(mean))
			conn.Write(res)
		}
	}
}

func insertPrice(prices map[int32]int32, ts int32, price int32) {
	prices[ts] = price
}

func queryMeanPrice(prices map[int32]int32, from int32, to int32) int32 {
	var sum, count int64

	for ts, price := range prices {
		if ts >= from && ts <= to {
			sum += int64(price)
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return int32(sum / count)
}

package main

import (
	"bufio"
	"encoding/json"
	"log"
	"math/big"
	"net"
	"protohackers"
)

type Request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func main() {
	config := protohackers.ParseConfig()
	server := protohackers.NewServer(config)
	server.StartTCP(handle)
}

func handle(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Bytes()

		log.Println("request: ", string(message))

		var req Request
		err := json.Unmarshal([]byte(message), &req)
		if err != nil || req.Method != "isPrime" || req.Number == nil {
			conn.Write([]byte("boom\n"))
			break
		}

		res := Response{Method: "isPrime", Prime: isPrime(int(*req.Number))}
		data, _ := json.Marshal(res)

		_, err = conn.Write(append(data, '\n'))
		if err != nil {
			log.Println("error: ", err)
			break
		} else {
			log.Println("response: ", string(data))
		}
	}
}

func isPrime(n int) bool {
	return big.NewInt(int64(n)).ProbablyPrime(20)
}

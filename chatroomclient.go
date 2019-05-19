package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func MessageSend(conn net.Conn) {
	var input string
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)

		if strings.ToUpper(input) == "EXIT" {
			conn.Close()
			break
		}

		_, err := conn.Write([]byte(input))
		if err != nil {
			fmt.Println("client connect failed:", err.Error())
			break
		}
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	go MessageSend(conn)
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("over")
			os.Exit(0)
		}
		fmt.Println("receive server send message content:", string(buf))
	}

	fmt.Println("client is end")
}

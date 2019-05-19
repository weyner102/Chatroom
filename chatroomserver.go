package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	LOG_DIRECTORY = "./test.log"
)

var onlineConns = make(map[string]net.Conn)
var messageQueue = make(chan string, 1000)
var quitChan = make(chan bool)
var logger *log.Logger

func ProcessInfo(conn net.Conn) {
	buf := make([]byte, 1024)
	defer func(conn net.Conn) {
		addr := conn.RemoteAddr().String()
		delete(onlineConns, addr)
		conn.Close()
	}(conn)

	for {
		//numOfBytes  成功的读取的字节数
		//numOfBytes, err := conn.Read(buf)
		numOfBytes, err := conn.Read(buf)
		if err != nil {
			break
		}

		if numOfBytes != 0 {
			message := string(buf[:numOfBytes])
			messageQueue <- message
		}

	}
}

func doProcessMessage(message string) {
	contents := strings.Split(message, "#")
	if len(contents) > 1 {
		addr := contents[0]
		sendMessage := strings.Join(contents[1:], "#")
		addr = strings.Trim(addr, " ")

		if conn, ok := onlineConns[addr]; ok {
			_, err := conn.Write([]byte(sendMessage))
			if err != nil {
				fmt.Println("onlineContents send failer")
			}
		}
	} else {
		contents := strings.Split(message, "*")
		if strings.ToUpper(contents[1]) == "LIST" {
			var ips string = ""
			for i := range onlineConns {
				ips = ips + "|" + i
			}

			if conn, ok := onlineConns[contents[0]]; ok {
				_, err := conn.Write([]byte(ips))
				if err != nil {
					fmt.Println("onlineContents send failer")
				}
			}
		}
	}
}

func ConsumeMessage() {
	for {
		select {
		case message := <-messageQueue:
			//对消息进行解析
			doProcessMessage(message)
		case <-quitChan:
			break
		}
	}
}

func main() {
	logFile, err := os.OpenFile(LOG_DIRECTORY, os.O_RDWR|os.O_CREATE, 0)
	if err != nil {
		fmt.Println("logFile create failure!")
		os.Exit(-1)
	}
	defer logFile.Close()

	logger = log.New(logFile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)

	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	fmt.Println("server is waiting...")
	logger.Println("server is start")

	go ConsumeMessage()

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println(conn.RemoteAddr().String())
		//将conn存储到映射表中
		onlineConns[conn.RemoteAddr().String()] = conn

		go ProcessInfo(conn)
	}

}

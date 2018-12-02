package main

import (
	"bytes"
	"io"
	"time"
	"fmt"
	"strings"
	"net"
)

const (
	SERVER_NETWORK = "tcp"
	SERVER_ADDRESS = "127.0.0.1:8001"
	// 数据边界
	DELIMITER = '\t'
)

func main() {

}

// 服务端
func serverGo() {
	// 首先创建一个监听器
	var listener net.Listener
	listener, err := net.Listen(SERVER_NETWORK, SERVER_ADDRESS)
	if err != nil {
		printServerLog("Listen Error: %s", err)
		return
	}
	defer listener.Close()
	printServerLog("Got listener for the server.(local address: %s)", listener.Addr())

	// 等待用户连接请求
	for {
		// 该方法会一直阻塞, 直到新连接到来
		conn, err := listener.Accept()
		if err != nil {
			printServerLog("Accept Error: %s", err)
		}
		printServerLog("Established a connection with a client application. (remote address: %s)", conn.RemoteAddr())

		// 单独处理每一个连接
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	for {
		// 设置读取超时时间
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		strReq, err := read(conn)
		if err != nil {
			if err == io.EOF {
				printServerLog("The connection is closed by another side.")
			} else {
				printServerLog("Read Error: %s", err)
			}
			break
		}
		printServerLog("Received request: %s.", strReq)

		// 没写完

	}
}

func read(conn net.Conn) (string, error) {
	readBytes := make([]byte, 1)
	var buffer bytes.Buffer

	for {
		_, err := conn.Read(readBytes)
		if err != nil {
			return "", err
		}
		readByte := readBytes[0]
		if readByte == DELIMITER {
			break
		}
		buffer.WriteByte(readByte)
	}
	return buffer.String(), nil
}

func printServerLog(format string, args ...interface{}){
	printLog("Server", 0, format, args...)
}

func printClientLog(sn int, format string, args ...interface{}){
	printLog("Client", sn, format, args...)
}

func printLog(role string, sn int, format string, args ...interface{}){
	// format 指的是格式,后面应该有换行符
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Printf("%s[%d]: %s", role, sn, fmt.Sprintf(format, args...))
}
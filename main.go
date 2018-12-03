package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	SERVER_NETWORK = "tcp"
	SERVER_ADDRESS = "127.0.0.1:8001"
	// 数据边界
	DELIMITER = '\t'
)

var wg sync.WaitGroup

func main() {
	wg.Add(2)
	go serverGo()
	time.Sleep(500 * time.Millisecond)
	go clientGo(1)
	wg.Wait()
}

// 客户端
func clientGo(id int) {
	defer wg.Done()
	// 建立连接
	conn, err := net.DialTimeout(SERVER_NETWORK, SERVER_ADDRESS, 2*time.Second)
	if err != nil {
		printClientLog(id, "Dial Error: %s.", err)
		return
	}
	defer conn.Close()
	printClientLog(id, "Connected to Server. (remote address: %s, local address: %s)", conn.RemoteAddr(), conn.LocalAddr())
	time.Sleep(200 * time.Millisecond)

	// 发送请求
	requestNumber := 5
	conn.SetDeadline(time.Now().Add(5 * time.Millisecond))
	for i := 0; i < requestNumber; i++ {
		req := rand.Int31()
		n, err := write(conn, fmt.Sprintf("%d", req))
		if err != nil {
			printClientLog(id, "Write Error: %s.", err)
			continue
		}
		printClientLog(id, "Sent request (written %d bytes): %d.", n, req)
	}

	// 接收请求
	for j := 0; j < requestNumber; j++ {
		strResp, err := read(conn)
		if err != nil {
			if err == io.EOF {
				printClientLog(id, "The connection is closed by another side.")
			} else {
				printClientLog(id, "Read Error: %s.", err)
			}
			break
		}
		printClientLog(id, "Received response: %s.", strResp)
	}
}

// 服务端
func serverGo() {
	defer wg.Done()
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
	defer conn.Close()
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

		// 处理请求
		intReq, err := strToInt32(strReq)
		if err != nil {
			// 将错误发送出去
			n, err := write(conn, err.Error())
			printServerLog("Sent error message (written %d bytes): %s.", n, err)
			continue
		}
		floatResq := cbrt(intReq)
		respMsg := fmt.Sprintf("The cube root of %d is %f.", intReq, floatResq)
		// 将结果发送回客户端
		n, err := write(conn, respMsg)
		if err != nil {
			printServerLog("Write Error: %s", err)
		}
		printServerLog("Sent response (written %d bytes): %s", n, respMsg)
	}
}

// 求立方根
func cbrt(param int32) float64 {
	return math.Cbrt(float64(param))
}

// 将字符串转换成 Int32
func strToInt32(str string) (int32, error) {
	num, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("\"%s\" is not integer", str)
	}
	if num > math.MaxInt32 || num < math.MinInt32 {
		return 0, fmt.Errorf("%d is not 32-bit integer", num)
	}
	return int32(num), nil
}

// 向 connection 写入内容
func write(conn net.Conn, content string) (int, error) {
	var buffer bytes.Buffer
	buffer.WriteString(content)
	buffer.WriteByte(DELIMITER)
	return conn.Write(buffer.Bytes())
}

// 从 connection 中读取内容
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

// 服务端日志输出
func printServerLog(format string, args ...interface{}) {
	printLog("Server", 0, format, args...)
}

// 客户端日志输出
func printClientLog(sn int, format string, args ...interface{}) {
	printLog("Client", sn, format, args...)
}

func printLog(role string, sn int, format string, args ...interface{}) {
	// format 指的是格式,后面应该有换行符
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Printf("%s[%d]: %s", role, sn, fmt.Sprintf(format, args...))
}

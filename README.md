# socket_by_go

使用 go 语言的并发机制来学习 socket 编程。

## 概述

一个小的项目实例，该实例包含服务端程序和客户端程序，以网络的 TCP 协议作为通信基础。

服务端程序的功能：接收客户端程序请求，计算请求数据的立方根，并把对结果的描述返回给客户端程序。

客户端程序的功能：向服务端程序发送若干个整数的请求数据，接收服务端程序返回的响应数据并记录他们。

## socket 流程

![socket.png](https://i.loli.net/2018/12/03/5c04fa4b860b8.png)

## 一些细节

    func Listen(net, laddr string) (Listener, error)

使用 net.Listen 函数获取监听器，接受两个 string 类型的参数。
第一个参数含义是以何种协议监听给定的地址。
第二个参数 laddr( Local Address 的简写)表示当前程序在网络中的标识，格式是： `host:port`
该函数会返回一个 net.Listener 类型的监听器

    listener, err := net.Listen("tcp", "127.0.0.1:8001")
    conn, err := listener.Accept()

当监听器调用 Accept 方法时，流程被阻塞，直到某个客户端程序与当前程序建立 TCP 连接，然后返回当前 TCP 连接的 net.Conn 类型值。


    func Dial(network, address string) (Conn, error)

**客户端**程序调用 Dial 函数用于向指定的网络地址发送连接建立申请。参数含义与服务端 Listen 函数基本一致(不是完全一致)

可以使用以下函数设置超时时间：

    func DialTimeout(network, address string, timeout time.Duration) (Conn, error)
  
time.Duration 是 int64 类型的别名，单位是纳秒。可以使用 time 标准代码库中预先声明的与常用时间单位对应的 time.Duration 类型的常量。例如： time.Nanosecond 代表一纳秒，值为 1； time.Microsecond 代表 1 微秒，其值为 1000 * Nanosecond。

在 net.Conn 类型中包含了八个方法。

- Read 方法

      Read(b []byte) (n int, err error)

    第一个参数相当于用来存放从连接上接收到的数据的**容器**。
    在一般情况下， Read 方法只有把参数值填满之后才返回。因此，我们可以这么做：

        b := make([]byte, 10)
        n, err := conn.Read(b)
        content := string(b[:n])

    此外，如果读取数据时发现 TCP 连接已经被另外一端关闭了，会返回一个 error 类型的值，该值与 io.EOF 值相等。

- Write 方法

      Write(b []byte) (n int, err error)

    第一个参数和 Read 方法一样，是一个容器，不过它是向当前 connection 写入数据。

- Close 方法

    该方法会关闭当前连接。

- LocalAddr 和 RemoteAddr 方法

    都返回一个 net.Addr 类型的结果，分别代表参与当前通信的客户端地址和服务端地址。

    net.Addr 类型包含两个方法

        conn.LocalAddr().NetWork()     ----->    返回所使用的协议名称
        conn.LocalAddr().String()      ----->    返回对应的网络地址

    - SetDeadline、SetReadDeadline、SetWriteDeadline 方法

    三个方法都接受一个 time.Time 类型的值作为参数，并返回一个 error 类型的值作为结果。

    SetDeadline() 设置当前连接的 I/O 操作的超时时间(包括当不限于读和写)。
    SetReadDeadline() 和 SetWriteDeadline() 分别设计读和写操作的超时时间。


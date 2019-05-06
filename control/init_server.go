package control

import (
	"net"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"bufio"
	"io"
	"assets_server/depend/server"
	"assets_server/handler"
)

func InitServer() {

	// Listen
	serverAddr := "127.0.0.1:8004"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", serverAddr)
	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		fmt.Println("net.Listen err = ", err)
		return
	}
	defer listener.Close()

	// Register & Path
	RegisterServer()
	for route := range server.ServerHandler.ServerList {
		fmt.Println("support route : ", route)
	}
	fmt.Println("tcp server success listen ", serverAddr)

	// 监听退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		quitMsg := <-quit
		listener.Close()

		quitLog := fmt.Sprintf("server stop by signal %s", quitMsg)
		fmt.Println(quitLog)
		os.Exit(1)
	}()

	//接受多个用户
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("listener Accept err=", err)
			return
		}

		go HandleConn(conn)
	}
}

//处理用户请求
func HandleConn(conn *net.TCPConn) {

	//获取客户端的网络地址信息
	addr := conn.RemoteAddr().String()
	fmt.Println(addr, " conncet successful")

	r := bufio.NewReader(conn)
	for {
		//读取用户数据
		//sc := bufio.NewScanner(Conn.)
		buf, err := r.ReadString('\n')

		if err != nil && err != io.EOF {
			fmt.Println("err = ", err)
			return
		}

		addr := fmt.Sprintf("[%s]", addr)
		if "exit" == string(buf) {
			fmt.Println("connect ", addr, " exit")
			conn.Close()
			return
		}

		// 分发请求
		//writeMsg := fmt.Sprintf("客户端: %s, 发来消息: %s", addr, string(buf))
		//write, err := comm.Utf8ToGbk(buf)
		//fmt.Println(write, err)
		// 发送到分发器
		NewDispatcher.requestChan <- &requestInfo{
			conn: conn,
			msg:  []byte(buf),
		}
	}
}

// 注册服务
func RegisterServer() {

	// 测试方法
	server.ServerHandler.RegisterService("say_hello", handler.SayHelloWorld)
	server.ServerHandler.RegisterService("is_me", handler.IsMe)

	// 正式方法
}

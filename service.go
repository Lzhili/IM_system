package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User //在线用户的map
	mapLock   sync.RWMutex     //锁
	Message   chan string      //消息广播的channel
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// Server的启动服务器的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	//close listen socket
	defer listener.Close()

	//启动监听Message的goroutine
	go this.ListenMessage()

	//不断循环等待新的客户端链接，每次建立新的链接就创建一个handler的goroutine
	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}
		//do handler
		go this.Handler(conn)
	}
}

// 监听service的Message广播消息channel的goroutine，一旦有消息就广播给全部在线的User
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		//将msg发送给全部在线user，user监听到再发送给客户端
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "] " + user.Name + ": " + msg
	this.Message <- sendMsg
}

// handler
func (this *Server) Handler(conn net.Conn) {
	//当前链接的业务
	//fmt.Println("Link established successfully")

	//创建一个新用户
	user := NewUser(conn, this)

	//用户上线
	user.Online()

	//监听用户是否活跃的channel
	isLive := make(chan bool)

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline() //用户下线
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn.Read err:", err)
				return
			}
			//提取用户的消息（去除"\n"）
			//msg := string(buf[:n-1])
			msg := string(buf[:n])

			//用户针对msg消息进行处理
			user.DoMessage(msg)

			//用户任意消息，代表当前用户是活跃的
			isLive <- true
		}
	}()

	//当前handler阻塞
	for {
		select {
		case <-isLive:
			//当前用户是活跃的，应该重置定时器
			//不做任何事情，为了激活select，更新下面的定时器
		case <-time.After(time.Second * 300):
			//已经超时，将当前的User强制关闭
			user.sendMsg("You got kicked")
			//释放资源
			close(user.C)
			//关闭连接
			conn.Close()
			//退出当前的handler
			return
		}
	}
}

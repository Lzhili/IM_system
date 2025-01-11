package main

import (
	client "IM_system/client"
	"flag"
	"fmt"
)

// 全局变量
var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "server ip address (default 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "server port (default 8888)")
}

func main() {
	//命令行解析
	flag.Parse()

	client := client.NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>> Link to server failed.....")
	}

	//一旦有新的client就单独开启一个goroutine去处理server的回执消息
	go client.DealResponse()

	fmt.Println(">>>>> Successfully linked to server......")
	fmt.Println(">>>>> serverIp:", serverIp, " and serverPort:", serverPort)

	//启动客户端的业务
	client.Run()
}

package client

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	//1.创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	//2.链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn
	//3.返回对象
	return client
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.Public chat mode")
	fmt.Println("2.Private chat mode")
	fmt.Println("3.Update username")
	fmt.Println("0.Exit")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>> Please enter a valid number <<<<")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 { //不断循环直到输入0退出
		for client.menu() != true {
		} //不断循环直到合法输入0-3

		//根据不同模型处理不同的业务
		switch client.flag {
		case 1:
			//公聊模式
			fmt.Println(">>>>> Public chat mode")
			break
		case 2:
			//私聊模式
			fmt.Println(">>>>> Private chat mode")
			break
		case 3:
			//更新用户名
			fmt.Println(">>>>> Update username")
			break
		}
	}
}

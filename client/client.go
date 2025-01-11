package client

import (
	"fmt"
	"io"
	"net"
	"os"
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

// 处理server回应的消息，直接显示到标准输出即可
func (client *Client) DealResponse() {
	//一旦client.conn有消息，直接copy到stdout标准输出，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
	//上面这行代码等价于下面这段
	//for {
	//	buf := make([]byte, xxx)
	//	client.conn.Read(buf)
	//	fmt.Println(string(buf))
	//}

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

// 公聊模式
func (client *Client) PublicChat() {
	//提示用户输入消息
	var chatMsg string
	fmt.Println(">>>> Please enter the chat message, or input \"exit\" to exit")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" { //一直循环直到输入exit退出
		//消息不为空则发送给服务器
		if len(chatMsg) != 0 {
			_, err := client.conn.Write([]byte(chatMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}
		chatMsg = "" //记得清空
		fmt.Println(">>>> Please enter the chat message, or input \"exit\" to exit")
		fmt.Scanln(&chatMsg)
	}
}

// 更新用户名
func (client *Client) UpdateName() bool {
	fmt.Println(">>>> Please enter new name: ")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
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
			client.PublicChat()
			break
		case 2:
			//私聊模式
			fmt.Println(">>>>> Private chat mode")
			break
		case 3:
			//更新用户名
			fmt.Println(">>>>> Update username")
			client.UpdateName()
			break
		}
	}
}

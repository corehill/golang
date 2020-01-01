package process

import (
	"../../common/message"
	"../utils"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func ShowMenu()  {
	fmt.Println("login success")
	fmt.Println("1 online user list")
	fmt.Println("2 send msg")
	fmt.Println("3 meg list")
	fmt.Println("4 quit")
	fmt.Println("1-4")
	var content string
	var key int
	fmt.Scanln(&key)
	//总会使用到smsProcess实例，因此外面将其定义在swith外部
	smsProcess := &SmsProcess{}
	switch key {
	case 1:
		outputOnlineUser()
		//fmt.Println()
	case 2:
		fmt.Println("please tying the group sms")
		fmt.Scanln(&content)
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println()
	case 4:
		fmt.Println()
		os.Exit(0)
	default:
		fmt.Println("the num is wrong,try again")


	}
}


//和服务器端保持通信
func serverProcessMes(conn net.Conn)  {
	//创建一个transfer实例，不停的读取服务器发送的信息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("client read the msg from server")
		mes,err := tf.ReadPkg()
		if err != nil {
			fmt.Println("read server err", err)
			return
		}
		//如果读取到消息，又是下一步处理逻辑
		//fmt.Println(mes)
		switch mes.Type {
		case message.NotifyUserStatuesMesType://有人上线
		//处理
			//1.取出notfiymes
			var notifyUserStatuesMes message.NotifyUserStatuesMes
			json.Unmarshal([]byte(mes.Date), &notifyUserStatuesMes)
			//2.把该用户的信息 状态保存到客户map中
			updateUserStatus(&notifyUserStatuesMes)
		case message.SmsMesType://有人群发消息
			outputGroupMes(&mes)
		default:
			fmt.Println("server return unknown mes")
		}
	}
}
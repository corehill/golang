package main

import (
	"fmt"
	"io"
	"net"
	"../../common/message"
	"../process"
	"../utils"
)

//先创建一个processor的结构体
type Processor struct {
	Conn net.Conn
	//这里调用总控 创建一个实例

}

//编写一个serverproceessmes函数
//功能： 根据客户端发送消息的种类不同，决定调用哪个函数来处理
func (this *Processor)serverProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		//处理登录
		//创建一个UserProcess实例
		up := &process.UserProcess{
			Conn : this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		//处理注册
		up := &process.UserProcess{
			Conn : this.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		//创建一个SmsProcess实例完成转发群聊消息的任务
		smsProcess := process.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	default:
		fmt.Println("the message not exists")
	}
	return
}

func (this *Processor) process2() (err error) {
	for {
		//这里我们将读取数据包， 直接封装成一个函数readPkg()， 返回Message, err
		//创建一个transfer 实例 完成读包任务
		tf := &utils.Transfer{
			Conn : this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF{
				fmt.Println("client shut down")
				return err
			} else{
				fmt.Println("read conn fail", err)
				return err
			}
		}
		fmt.Println("mes=",mes)

		err = this.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}
}

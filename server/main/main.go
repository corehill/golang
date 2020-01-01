package main

import (
	"../model"
	"fmt"
	"net"
)

//func readPkg(conn net.Conn) (mes message.Message, err error) {
//	buf := make([]byte, 8096)
//	fmt.Println("waiting for read")
//	//conn.read 在conn没有被关系的情况下，才会阻塞
//	//如果客户端关闭了conn，就不会阻塞
//	_, err = conn.Read(buf[:4])
//	if err != nil {
//		//fmt.Println("conn read err", err)
//		return
//	}
//
//	//根据读到的长度 转成一个uint32类型
//	var pkgLen uint32
//	pkgLen = binary.BigEndian.Uint32(buf[:4])
//
//	n, err := conn.Read(buf[:pkgLen])
//	if n != int(pkgLen) || err != nil {
//		//fmt.Println("conn read fail err", err)
//		return
//	}
//
//	//把pkgLen 反序列化
//	err = json.Unmarshal(buf[:pkgLen], &mes)
//	if err != nil {
//		fmt.Println("json unmarshal faril err", err)
//		return
//	}
//
//	return
//}
//
//func writePkg(conn net.Conn, data []byte) (err error) {
//	//发送一个长度给对方
//	var pkgLen uint32
//	pkgLen = uint32(len(data))
//	var buf [4]byte
//	binary.BigEndian.PutUint32(buf[:4], pkgLen)
//	//发送长度
//	n, err := conn.Write(buf[:4])
//	if n != 4 || err != nil {
//		fmt.Println("conn fail", err)
//		return
//	}
//
//	//发送消息本身
//	n, err = conn.Write(data)
//	if n != int(pkgLen) || err != nil{
//		fmt.Println("send data err", err)
//		return
//	}
//	return
//}

////编写一个函数serverProcessLogin函数，专门处理登录请求
//func serverProcessLogin(conn net.Conn, mes *message.Message) (err error) {
//	var loginMes message.LoginMes
//	err = json.Unmarshal([]byte(mes.Date), &loginMes)
//	if err != nil {
//		fmt.Println("server unmarshal err=", err)
//		return
//	}
//	//先声明一个res.Mes
//	var resMes message.Message
//	resMes.Type = message.LoginResMesType
//
//	//再申明一个loginRes
//	var logResMes message.LoginResMes
//
//	if loginMes.UsrId == 100 && loginMes.UsrPwd == "123456" {
//		logResMes.Code = 200
//		fmt.Println("login success")
//	}else {
//		logResMes.Code = 500 //500表示用户不存在
//		logResMes.Error = "this client not exiest, register first"
//	}
//
//	//将resresmsg序列化并赋值给resmes.data
//	data, err := json.Marshal(logResMes)
//	if err != nil {
//		fmt.Println("logresmes marshal err=", err)
//		return
//	}
//
//	resMes.Date = string(data)
//
//	//将整个resMes序列化，并发送
//	data, err = json.Marshal(resMes)
//	if err != nil {
//		fmt.Println("logresmes marshal err=", err)
//		return
//	}
//	err = writePkg(conn, data)
//	if err != nil {
//		fmt.Println("write to client", err)
//	}
//	return
//}
//
//
//
////编写一个serverproceessmes函数
////功能： 根据客户端发送消息的种类不同，决定调用哪个函数来处理
//func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {
//	switch mes.Type {
//	case message.LoginMesType:
//		err = serverProcessLogin(conn, mes)
//	case message.RegisterMesType:
//		fmt.Println()
//	default:
//		fmt.Println("the message not exists")
//	}
//	return
//}

func init()  {
	//当服务器启动前， 我们就去初始化我们的redis的连接池
	initPool("localhost:6379",16,0,300)
	initUserDao()
}

//处理和客户端的通讯
func processMain(conn net.Conn)  {
	//读客户端发送的信息
	defer conn.Close()
	//循环读取
	processor := &Processor{Conn:conn}
	err := processor.process2()
	if err != nil {
		fmt.Println("client and server err", err)
		return
	}
}

//这里我们编写一个函数，完成对UserDao的初始化任务
func initUserDao()  {
	//这里的pool本身就是一个全局的变量
	//这里需要注意一个初始化顺序问题
	//initPool， 在initUserDao之前
	model.MyUserDao = model.NewUserDao(pool)
}

func main()  {
	//提示信息
	fmt.Println("server 8889")
	listen,err := net.Listen("tcp","0.0.0.0:8889")
	if err != nil {
		fmt.Println("listen err=", err)
		return
	}
	defer listen.Close()
	for {
		fmt.Println("wait the client")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen accept err=", err)
		}
		//一旦链接成功，则启动一个协程和客户端保持通讯
		go processMain(conn)
	}

}

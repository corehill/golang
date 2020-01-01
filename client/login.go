package main

import (
	"../common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	_"time"
)

func login(usrId int, usrPwd string) (err error) {
	//下一个就要开始定协议
	//fmt.Printf("usrId=%d usrpwd", usrId, usrPwd)
	//return nil
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net dial err=", err)
		return
	}
	defer conn.Close()
	//准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType

	//创建一个loginmes 结构体
	var loginMes message.LoginMes
	loginMes.UsrId = usrId
	loginMes.UsrPwd = usrPwd

	//将loginmes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json marshal err=", err)
		return
	}

	//把  data赋值给mes.data
	mes.Date = string(data)

	//将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json marshal sec err=", err)
		return
	}

	//此时，data就是我们要发送的数据
	//先获取到dat的长度 转成一个表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:4], pkgLen)
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn fail", err)
		return
	}

	//fmt.Println("client send pkglen success", len(data), string(data))
	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("send data fail", err)
		return
	}

	//这里还需要处理服务器端返回的消息
	mes, err = readPkg(conn)
	fmt.Println("mes", mes)
	if err != nil {
		fmt.Println("read err", err)
	}
	var LoginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Date), &LoginResMes)
	if LoginResMes.Code == 200 {
		fmt.Println("login success")
	} else if LoginResMes.Code == 500 {
		fmt.Println(LoginResMes.Error)
	} else {
		fmt.Println("read unmarshal err", err)
	}
	return
}

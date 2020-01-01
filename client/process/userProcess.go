package process

import (
	"../../common/message"
	"../utils"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {

}

//给关联一个用户登录的方法
func (this *UserProcess) Login(usrId int, usrPwd string) (err error) {
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
	tf := &utils.Transfer{
		Conn : conn,
	}
	mes, err = tf.ReadPkg()
	fmt.Println("mes", mes)
	if err != nil {
		fmt.Println("read err", err)
	}
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Date), &loginResMes)
	if loginResMes.Code == 200 {
		//初始化CureUser
		CurUser.Conn = conn
		CurUser.UserId = usrId
		CurUser.UserStatus = message.UserOnline
		//可以显示当前用户列表
		fmt.Println("all user list:")
		for _, id := range loginResMes.UsersId{
			if id == usrId{
				continue
			}
			fmt.Println("user list:",id)
			//完成客户端的onlineUsers完成舒适化
			user := &message.User{
				UserId:     id,
				UserStatus: message.UserOnline,
			}
			onlineUsers[id] = user
		}
		fmt.Println()
		//这里我们还需要在客户端启动一个协程
		//该协程保持和服务器端的通信，如果服务器有数据推送给客户端
		//则接受并显示在客户端的终端
		go serverProcessMes(conn)
		//fmt.Println("login success")
		//循环显示登陆成功的菜单
		for {
			ShowMenu()
		}
	}  else {
		fmt.Println(loginResMes.Error)
	}
	return
}

func (this *UserProcess) Register(usrId int, usrPwd string, usrName string) (err error) {
	//链接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net dial err=", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//准备通过conn发送消息给服务
	var mes message.Message
	mes.Type = message.RegisterMesType
	//创建一个loginMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserPwd = usrPwd
	registerMes.User.UserId = usrId
	registerMes.User.UserName = usrName
	//将registerRes序列化

	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json marshal err=", err)
		return
	}
	mes.Date = string(data)
	//将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json marshal sec err=", err)
		return
	}
	//创建一个transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}

	//发送data给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("register writepkg err=", err)
		return
	}

	mes, err = tf.ReadPkg() //mes就是mesgisterResMes
	if err != nil {
		fmt.Println("register readpkg err=", err)
		return
	}
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Date), &registerResMes)
	if registerResMes.Code == 200 {
		//这里我们还需要在客户端启动一个协程
		//该协程保持和服务器端的通信，如果服务器有数据推送给客户端
		//则接受并显示在客户端的终端
		fmt.Println("register success")
		os.Exit(0)
	}  else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return
}
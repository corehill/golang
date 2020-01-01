package process

import (
	"../../common/message"
	"../model"
	"../utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	// 字段
	Conn net.Conn
	// 增加一个字段，表示conn是哪一个用户
	UserId int
}

//这里我们编写通知所有在线用户的方法
//userId 要通知其他在线用户 我上线
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {
	//遍历onlineUsers， 然后依次发送NotifyUserStatuesMes
	for id, up := range userMgr.onlineUsers{
		//过滤掉自己
		if id == userId {
			continue
		}
		//开始通知：单独写一个方法
		up.NotifyMeOnline(userId)
	}
}

func (this *UserProcess) NotifyMeOnline(userId int) {
	//组装NotifyUserStatuesMes
	var mes message.Message
	mes.Type = message.NotifyUserStatuesMesType

	var notifyUserStatuesMes message.NotifyUserStatuesMes
	notifyUserStatuesMes.UserId = userId
	notifyUserStatuesMes.Status = message.UserOffline

	//将notifyUserStatuesMes序列化
	data, err := json.Marshal(notifyUserStatuesMes)
	if err != nil {
		fmt.Println("notify err", err)
		return
	}

	mes.Date = string(data)
	//再次序列化 准备发送
	data,err = json.Marshal(mes)
	if err != nil {
		fmt.Println("notify2 err", err)
		return
	}

	//发送
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("notify3 err", err)
		return
	}
}


func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	//从mes中取出mes.Data 并直接反序列化成RegisterMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Date), &registerMes)
	if err != nil {
		fmt.Println("server unmarshal err=", err)
		return
	}

	var resMes message.Message
	resMes.Type = message.RegisterResMesType

	//我们需要到redis数据库去完成注册
	//再申明一个loginRes
	var registerResMes message.RegisterResMes

	//使用model.MyUserDao 到redis去验证
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err  == model.ERROR_USER_EXISTS{
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误"
		}
	} else {
		registerResMes.Code = 200
	}
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("logresmes marshal err=", err)
		return
	}
	resMes.Date = string(data)

	//将整个resMes序列化，并发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("logresmes marshal err=", err)
		return
	}

	//发送data 我们将封装到writePKG
	tf := &utils.Transfer {
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("write to client", err)
	}
	return

}

//编写一个函数serverProcessLogin函数，专门处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Date), &loginMes)
	if err != nil {
		fmt.Println("server unmarshal err=", err)
		return
	}
	//先声明一个res.Mes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	//再申明一个loginRes
	var logResMes message.LoginResMes

	//我们需要到redis数据库去完成验证
	//使用model.MyUserDao 到redis去验证
	user, err := model.MyUserDao.Login(loginMes.UsrId,loginMes.UsrPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			logResMes.Code = 500
			logResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD{
			logResMes.Code = 403
			logResMes.Error = err.Error()
		} else {
			logResMes.Code = 505
			logResMes.Error = "unknow error"
		}
	} else {
		logResMes.Code = 200
		//这里， 因为用户登录成功，我们就把改登录成功的用户放入到userMgr中
		//将登录成功的用户userId 赋给this
		this.UserId = loginMes.UsrId
		userMgr.AddOnlineUser(this)
		this.NotifyOthersOnlineUser(loginMes.UsrId)
		//通知其他在线用户， 我上线了
		fmt.Println("login success",user)
		//将当前在线用户的id放入到logResMes.UsersId
		//遍历onlineUsers
		for id, _ := range userMgr.onlineUsers{
			logResMes.UsersId = append(logResMes.UsersId, id)
			fmt.Println("online client:",id)
		}
	}


	//if loginMes.UsrId == 100 && loginMes.UsrPwd == "123456" {
	//	logResMes.Code = 200
	//	fmt.Println("login success")
	//}else {
	//	logResMes.Code = 500 //500表示用户不存在
	//	logResMes.Error = "this client not exiest, register first"
	//}

	//将resresmsg序列化并赋值给resmes.data
	data, err := json.Marshal(logResMes)
	if err != nil {
		fmt.Println("logresmes marshal err=", err)
		return
	}

	resMes.Date = string(data)

	//将整个resMes序列化，并发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("logresmes marshal err=", err)
		return
	}

	//发送data 我们将封装到writePKG
	tf := &utils.Transfer {
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("write to client", err)
	}
	return
}





package message

const (
	LoginMesType = "LoginMes"
	LoginResMesType = "LoginResMes"
	RegisterMesType = "RegisterMes"
	RegisterResMesType = "RegisterResMes"
	NotifyUserStatuesMesType = "NotifyUserStatuesMes"
	SmsMesType = "SmsMes"
)

//定义几个用户状态的常量
const(
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"`//消息类型
	Date string `json:"date"`//消息内容类型
}


//定义两个消息。。需要再增加
type LoginMes struct {
	UsrId int `json:"usrId"`//用户id
	UsrPwd string `json:"usrPwd"`//用户密码
}


type LoginResMes struct {
	Code int `json:"code"`//返回状态码 500 表示用户未注册 200 表示登录成功
	UsersId []int `json:"users"`
	Error string `json:"error"`//返回错误信息
}

type RegisterMes struct {
	User User `json:"user"`//类型就是User结构体
}

type RegisterResMes struct {
	Code int `json:"code"`//返回状态码 400 表示用户已存在 200 表示注册成功
	Error string `json:"error"`//返回错误信息
}

//为了配合服务器端推送用户状态变化的消息
type NotifyUserStatuesMes struct {
	UserId int `json:"userId"` //用户id
	Status int `json:"status"` //用户状态
}

//增加一个SmsMes //发送的消息
type SmsMes struct {
	Content string `json:"content"`
	User //匿名结构体
}

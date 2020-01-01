package process

import (
	"../../common/message"
	"../model"
	"fmt"
)

//客户端要维护的map

var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
var CurUser model.CurUser //我们在登陆成功后。完成对CurUser的初始化

//在客户端显示当前在线的用户
func outputOnlineUser() {
	//遍历onlineuser
	fmt.Println("online user list now")
	for id, user := range onlineUsers{
		//如果不显示自己
		fmt.Println("user id",id,"user", user)
	}
}

//编写一个方法，处理返回的信息
func updateUserStatus(notifyUserStatuesMes *message.NotifyUserStatuesMes)  {
	//适当优化
	user, ok := onlineUsers[notifyUserStatuesMes.UserId]
	if !ok { //原来没有
		user = &message.User{
			UserId: notifyUserStatuesMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatuesMes.Status
	onlineUsers[notifyUserStatuesMes.UserId] = user
	outputOnlineUser()
}
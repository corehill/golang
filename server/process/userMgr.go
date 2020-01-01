package process

import "fmt"

var (
	userMgr *UserMgr
)

//因为UserMgr实例在服务器端有且只有一个
type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

//完成对userMgr的初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//完成对OnlineUsers添加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}

//删除
func (this *UserMgr) DelOnlineUser(up *UserProcess) {
	delete(this.onlineUsers, up.UserId)
}

//查询，返回当前所有在线用户
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}

//根据ID返回对应的值
func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	//如何从map中取出一个值，带检测的方式
	up, ok := this.onlineUsers[userId]
	if !ok {//说明你要查找的该用户不在线
		err = fmt.Errorf("client %d not online", userId)
		return
	}
	return
}
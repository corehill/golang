package model

import (
	"../../common/message"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

//在服务器启动后，就初始化一个userDao实例
//把它做成全局的变量，在需要和redis操作时，就直接使用即可
var (
	MyUserDao *UserDao
)

//定义一个UserDao结构体
//完成对User结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

//使用工厂模式，创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao){
	userDao = &UserDao{
		pool,
	}
	return
}


//思考一下 在UserDao应该提供哪些方法给我们
//根据用户id 返回一个User实例+err
func (this *UserDao) getUserBtId(conn redis.Conn, id int) (user *User, err error){
	//通过给定的ID去REDIS里面查询用户
	res, err := redis.String(conn.Do("hget", "users", id))
	if err != nil {
		//错误
		if err == redis.ErrNil {//表示在users哈希中， 没有找到对应id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}

	user = &User{}
	//这里我们需要把res反序列化成user实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json unmarshal err", err)
		return
	}
	return
}

//完成登录校验
//login 完成对用户的验证
//如果用户id和pwd都正确，则返回一个user实例
//如果用户的id和pwd有错误，则返回对应的错误信息
func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	//先从连接池中取出连接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserBtId(conn, userId)
	if err != nil {
		return
	}
	//此时用户信息已经获取到
	if user.UserPwd != userPwd{
		err = ERROR_USER_PWD
		return
	}
	return
}


func (this *UserDao) Register(user *message.User) (err error) {
	//先从连接池中取出连接
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.getUserBtId(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}
	//这时， 说明id在redis还没有，则可以完成注册
	data, err := json.Marshal(user)  //序列化user
	if err != nil {
		return
	}
	//入库
	fmt.Println(string(data))
	_, err = conn.Do("hset", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误 err=", err)
		return
	}
	return
}
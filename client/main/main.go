package main

import "fmt"
//import "../client/login"
import "../process"

var usrId int
var usrPwd string
var usrName string

func main() {
	//接收用户选择
	var key int
	loop := true
	for {
		fmt.Println("group talk")
		fmt.Println("1 login")
		fmt.Println("2 register")
		fmt.Println("3 exit")
		fmt.Println("choose 1-3")

		fmt.Scanln(&key)

		switch key {
		case 1:
			fmt.Println("chat room")
			fmt.Println("typing id")
			fmt.Scanln(&usrId)
			fmt.Println("typing pwd")
			fmt.Scanln(&usrPwd)
			//创建一个userprocess的实例
			up :=  &process.UserProcess{}
			up.Login(usrId,usrPwd)

		case 2:
			fmt.Println("register")
			fmt.Println("typing id")
			fmt.Scanln(&usrId)
			fmt.Println("typing pwd")
			fmt.Scanln(&usrPwd)
			fmt.Println("typing nickname")
			fmt.Scanln(&usrName)
			//2 调用UserProcess 完成注册的请求
			up :=  &process.UserProcess{}
			up.Register(usrId,usrPwd,usrName)
		case 3:
			fmt.Println()
			loop =false
		default:
			fmt.Println("wrong tying")
		}
		if loop == false{
			break
		}
	}
	//根据用户的输入，显示新的提示信息
	//if key == 1 {
	//	//用户需登录
	//
	//	//login(usrId, usrPwd)
	//	//因为使用了新的程序结构，我们创建
	//	//if err != nil {
	//	//	fmt.Println("login failure")
	//	//} else {
	//	//	fmt.Println("login sucess")
	//	//}
	//	//先把登录函数。写到另外一个文件
	//}else if key==2{
	//	fmt.Println("register")
	//}
}

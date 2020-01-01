package process

import (
	"encoding/json"
	"fmt"
	"../../common/message"
	"../utils"
)

type SmsProcess struct {

}

//发送群聊的消息

func(this *SmsProcess) SendGroupMes(content string) (err error) {
	//创建一个mes
	var mes message.Message
	mes.Type = message.SmsMesType

	//创建一个SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	//将smsMes进行序列化
	data,err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("Sms Marshal err", err)
		return
	}

	mes.Date = string(data)
	//将mes进行序列化，准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("Sms Marshal2 err", err)
		return
	}

	//将序列化后的mes发送
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("send Sms err", err)
		return
	}
	return
}
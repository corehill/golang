package process

import (
	"../../common/message"
	"../utils"
	"encoding/json"
	"fmt"
	"net"
)


type SmsProcess struct {
	//暂时不需要字段
}

func (this *SmsProcess) SendGroupMes(mes *message.Message) {
	//遍历服务器的onlineusers
	//将消息转发出去

	//取出mes内容 SmsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Date),&smsMes)
	if err != nil {
		fmt.Println("json unmarshal err", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json marshal err", err)
		return
	}

	for id, up := range userMgr.onlineUsers{
		//这里我们还需要过滤掉自己
		if id == smsMes.UserId{
			continue
		}
		this.SendMesToEachOnlineUser(data, up.Conn)
	}
}


func (this *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("transfer message fail")
	}
}
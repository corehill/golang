package process

import (
	"../../common/message"
	"encoding/json"
	"fmt"
)

func outputGroupMes(mes *message.Message)  {//这个地方的mes一定SmsMes
	//显示即可
	//1. 反序列化
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Date), &smsMes)
	if err != nil {
		fmt.Println("json Unmarshal err",err)
		return
	}

	//显示
	fmt.Println(smsMes.Content)
}

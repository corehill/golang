package model

import (
	"net"
	"../../common/message"
)

//因为在科幻段，我们很多i地方会使用到curUser，我们将其作为一个全局变量
type CurUser struct {
	Conn net.Conn
	message.User
}
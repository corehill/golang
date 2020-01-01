package utils

import (
	"../../common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type Transfer struct {
	Conn net.Conn
	Buf [8096]byte //传输时  使用缓冲
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	//buf := make([]byte, 8096)
	fmt.Println("waiting for read")
	//conn.read 在conn没有被关系的情况下，才会阻塞
	//如果客户端关闭了conn，就不会阻塞
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		//fmt.Println("conn read err", err)
		return
	}

	//根据读到的长度 转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4])

	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//fmt.Println("conn read fail err", err)
		return
	}

	//把pkgLen 反序列化
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json unmarshal faril err", err)
		return
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	//发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[:4], pkgLen)
	//发送长度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn fail", err)
		return
	}

	//发送消息本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("send data err", err)
		return
	}
	return
}

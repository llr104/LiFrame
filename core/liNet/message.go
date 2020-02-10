package liNet

import (
	"LiFrame/utils"
)

type Message struct {
	NameLen uint32
	BodyLen uint32 //消息的长度
	Name    []byte //消息的ID
	Body    []byte //消息的内容
}


//创建一个Message消息包
func NewMsgPackage(id string, data []byte) *Message {
	m := Message{}
	m.SetMsgName(id)
	m.SetBody(data)
	return &m
}

//获取消息数据段长度
func (msg *Message) GetBodyLen() uint32 {
	return msg.BodyLen
}

func (msg *Message) GetNameLen() uint32 {
	return msg.NameLen
}

//获取消息ID
func (msg *Message) GetMsgName() string {
	var name string
	utils.DecodeObject(msg.Name,&name)
	return name
}

func (msg *Message) GetMsgNameByte() []byte{
	return msg.Name
}

//获取消息内容
func (msg *Message) GetBody() []byte {
	return msg.Body
}

//设置消息数据段长度
func (msg *Message) SetBodyLen(len uint32) {
	msg.BodyLen = len
}

func (msg *Message) SetNameLen(len uint32) {
	msg.NameLen = len
}

func (msg *Message) SetMsgNameByte(data []byte) {
	msg.NameLen = uint32(len(data))
	msg.Name = data
}

//设置消息ID
func (msg *Message) SetMsgName(name string) {
	msg.Name,_ = utils.EncodeObject(name)
	msg.NameLen = uint32(len(msg.Name))
}

//设置消息内容
func (msg *Message) SetBody(data []byte) {
	msg.BodyLen = uint32(len(data))
	msg.Body = data
}


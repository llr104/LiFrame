package liFace

/*
	将请求的一个消息封装到message中，定义抽象层接口
 */
type IMessage interface {
	GetBodyLen() uint32 //获取消息数据段长度
	GetNameLen() uint32

	GetMsgNameByte() []byte
	GetMsgName() string //获取消息ID
	GetBody() []byte    //获取消息内容
	GetSeq() uint32
	SetMsgNameByte([]byte) 	//设置消息ID
	SetMsgName(string) 		//设置消息ID
	SetNameLen(uint32)

	SetBody([]byte)    //设置消息内容
	SetBodyLen(uint32) //设置消息数据段长度

}

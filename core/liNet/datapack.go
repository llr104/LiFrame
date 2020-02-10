package liNet

import (
	"LiFrame/core/liFace"
	"LiFrame/utils"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/thinkoner/openssl"
)


//|---4 bytes---|---4 bytes---|-------dataLen-------|
//|-------------------------------------------------|
//|---dataLen---|---nameLen---|---name---|---body---|
//|-------------------------------------------------|

var DataPackKey = []byte("msgprotokey12345")
//封包拆包类实例，暂时不需要成员
type DataPack struct {}

//封包拆包实例初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包头长度方法
func(dp *DataPack) GetHeadLen() uint32 {
	//NameLen uint32(4字节) +  BodyLen uint32(4字节)
	return 8
}
//封包方法(压缩数据)
func(dp *DataPack) Pack(msg liFace.IMessage)([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//加密
	var dataLen uint32
	body := msg.GetBody()
	if dst, err := openssl.AesECBEncrypt(body, DataPackKey, openssl.PKCS7_PADDING); err != nil {
		return nil, err
	}else{
		body = dst
		dataLen = uint32(len(dst)) + msg.GetNameLen()
	}

	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, dataLen); err != nil {
		return nil, err
	}

	//写nameLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetNameLen()); err != nil {
		return nil, err
	}

	//写msgName
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgNameByte()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, body); err != nil {
		return nil ,err
	}

	return dataBuff.Bytes(), nil
}
//拆包方法(解压数据)
func(dp *DataPack) Unpack(binaryData []byte)(liFace.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	var dataLen uint32 = 0
	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &dataLen); err != nil {
		return nil, err
	}

	//读nameLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.NameLen); err != nil {
		return nil, err
	}

	//bodyLen
	msg.BodyLen = dataLen-msg.NameLen
	//判断dataLen的长度是否超出我们允许的最大包长度
	if utils.GlobalObject.MaxPacketSize > 0 && dataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg Data received")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}

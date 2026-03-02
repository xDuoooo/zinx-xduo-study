package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx-xduo-study/src/utils"
	"zinx-xduo-study/src/ziface"
)

//封包，拆包的具体实现

type DataPack struct {
}

func (dp *DataPack) GetHeadLen() uint32 {
	//Data length uint32(4字节) + ID uint32(4字节)
	return 8
}

// Pack |data len|msgID|data|
func (dp *DataPack) Pack(message ziface.IMessage) ([]byte, error) {
	//创建一个存放byte字节的缓冲
	buffer := bytes.NewBuffer([]byte{})
	//将dataLen写进dataBuff中
	if err := binary.Write(buffer, binary.LittleEndian, message.GetMsgLen()); err != nil {
		return nil, err
	}
	//将MsgId写进databuff中
	if err := binary.Write(buffer, binary.LittleEndian, message.GetMsgId()); err != nil {
		return nil, err
	}
	//将data数据写进databuff中
	if err := binary.Write(buffer, binary.LittleEndian, message.GetData()); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// UnPack 将包的Head信息读出来，之后再根据head信息里的data的长度再进行一次读
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	msg := &Message{}
	//读datalen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("this package is bigger than maxPackageSize")
	}

	//判断datalen是否已经超出了最大包大小
	return msg, nil

}

// NewDataPack 拆包封包初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

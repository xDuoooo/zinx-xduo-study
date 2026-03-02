package ziface

//定义一个对message解决TCP粘包问题的封包拆包模块 TLV 格式
//封包、拆包模块
//面向TCP连接中的数据流，用于解决TCP粘包问题

type IDataPack interface {
	GetHeadLen() uint32
	Pack(message IMessage) ([]byte, error)
	UnPack([]byte) (IMessage, error)
}

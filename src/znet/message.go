package znet

type Message struct {
	Id      uint32 //消息ID
	DateLen uint   //消息长度
	Data    []byte //消息内容
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}
func (m *Message) GetMsgLen() uint {
	return m.DateLen
}
func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}
func (m *Message) SetMsgLen(len uint) {
	m.DateLen = len
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}

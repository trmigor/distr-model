package messages

// Message type represents a message between processes.
type Message struct {
	SendTime     int64
	DeliveryTime int64
	From         int32
	To           int32
	Ptr          int
	Body         []byte
}

// NewMessage creates new instance of Message type by sender number, receiver number and body.
func NewMessage(from int32, to int32, body []byte) *Message {
	return &Message{
		From: from,
		To:   to,
		Body: body,
	}
}

// append adds data to body of a message.
func (msg *Message) append(a *MessageArg) {
	msg.Body = append(msg.Body, a.Body...)
}

// NewMessageByArgs creates new instance of Message type by message arguments.
func NewMessageByArgs(args ...*MessageArg) *Message {
	msg := NewMessage(-1, -1, make([]byte, 0))
	for _, n := range args {
		msg.append(n)
	}
	return msg
}

// GetInt32 extracts the earliest message argument that is not yet extracted if it is of type int32.
// If it is not or there is nothing to extract, panics.
func (msg *Message) GetInt32() int32 {
	if msg.Ptr+4 < len(msg.Body) && msg.Body[msg.Ptr] == Int32Type {
		var res uint32
		msg.Ptr++
		for i := 0; i <= 24; i += 8 {
			res |= uint32(msg.Body[msg.Ptr]&0xFF) << i
			msg.Ptr++
		}
		return int32(res)
	}
	panic("Expected int32")
}

// GetInt64 extracts the earliest message argument that is not yet extracted if it is of type int64.
// If it is not or there is nothing to extract, panics.
func (msg *Message) GetInt64() int64 {
	if msg.Ptr+8 < len(msg.Body) && msg.Body[msg.Ptr] == Int64Type {
		var res uint64
		msg.Ptr++
		for i := 0; i < 8; i++ {
			res = (res << 8) | uint64(msg.Body[msg.Ptr+7-i]&0xFF)
		}
		msg.Ptr += 8
		return int64(res)
	}
	panic("Expected int64")
}

// GetString extracts the earliest message argument
// that is not yet extracted if it is of type []byte.
// If it is not or there is nothing to extract, panics.
func (msg *Message) GetString() []byte {
	if msg.Ptr < len(msg.Body) && msg.Body[msg.Ptr] == StringType {
		res := make([]byte, 0)
		msg.Ptr++
		for msg.Ptr < len(msg.Body) && msg.Body[msg.Ptr] != 0 {
			res = append(res, msg.Body[msg.Ptr])
			msg.Ptr++
		}
		return res
	}
	panic("Expected string")
}

// GetData extracts the earliest message argument that is not yet extracted of any type.
// If there is no data, returns nil. If data is broken, panics.
func (msg *Message) GetData() interface{} {
	if msg.Ptr >= len(msg.Body) {
		return nil
	}
	switch msg.Body[msg.Ptr] {
	case Int32Type:
		return msg.GetInt32()
	case Int64Type:
		return msg.GetInt64()
	case StringType:
		return msg.GetString()
	}
	panic("Not expected type")
}

// Greater compares two messages by delivery time.
func Greater(first *Message, second *Message) bool {
	return first.DeliveryTime > second.DeliveryTime
}

package messages

const (
	// Int32Type marks message arguments of type int32.
	Int32Type byte = iota + 'A'
	// Int64Type marks message arguments of type int64.
	Int64Type
	// StringType marks message arguments of type []byte.
	StringType
)

// MessageArg type represents argument of a message.
type MessageArg struct {
	Body []byte
}

// NewMessageArg creates new instance of a MessageArg type according to input.
func NewMessageArg(q interface{}) *MessageArg {
	body := make([]byte, 0)
	switch v := q.(type) {
	case int32:
		body = append(body, Int32Type)
		for i := 0; i <= 24; i += 8 {
			body = append(body, byte((v>>i)&0xFF))
		}
	case int64:
		body = append(body, Int64Type)
		for i := 0; i < 8; i++ {
			body = append(body, byte(v&0xFF))
			v >>= 8
		}
	case []byte:
		body = append(body, StringType)
		body = append(body, v...)
		body = append(body, 0)
	default:
		panic("Not expected type")
	}
	return &MessageArg{body}
}

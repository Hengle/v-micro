package gotcp

import (
	"encoding/binary"
	"fmt"

	"github.com/fananchong/v-micro/transport"
)

/*
对 transport.Message 做序列化与反序列化
格式：
    |<--------------------------------------- Header ------------------------------------------------------->|<------- Body ------->|
    +-----[2]---------+----[2]-----+---[key1 len]---+-----[2]------+---[value1 len]---+---[more key-value]---+------[body len]------+
    | key-value count |  key1 len  |      key1      |  value1 len  |      value1      |         ...          |         body         |
    +-----------------+------------+----------------+--------------+------------------+----------------------+----------------------+

*/

func marshal(msg *transport.Message) ([]byte, error) {
	totalSize := size(msg)
	data := make([]byte, totalSize)
	offset := 0
	// header
	headerSize := uint16(len(msg.Header))
	if ok := marshalUint16(data, &offset, headerSize); !ok {
		return nil, fmt.Errorf("marshal header size error")
	}
	for k, v := range msg.Header {
		if ok := marshalString2(data, &offset, k); !ok {
			return nil, fmt.Errorf("marshal header key error. key:%s, value:%s", k, v)
		}
		if ok := marshalString2(data, &offset, v); !ok {
			return nil, fmt.Errorf("marshal header value error. key:%s, value:%s", k, v)
		}
	}
	// body
	if ok := marshalSliceByte(data, &offset, msg.Body); !ok {
		return nil, fmt.Errorf("marshal body error")
	}
	return data, nil
}

func unmarshal(data []byte) (*transport.Message, error) {
	maxLen := len(data)
	offset := 0
	var headerSize uint16
	var ok bool
	msg := &transport.Message{
		Header: make(map[string]string),
	}
	// header
	if headerSize, ok = unmarshalUint16(data, &offset, maxLen); !ok {
		return nil, fmt.Errorf("unmarshal header size error")
	}
	for i := 0; i < int(headerSize); i++ {
		var k string
		var v string
		var ok1 bool
		var ok2 bool
		if k, ok1 = unmarshalString2(data, &offset, maxLen); ok1 {
			v, ok2 = unmarshalString2(data, &offset, maxLen)
		}
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("unmarshal header data error")
		}
		msg.Header[k] = v
	}
	// body
	msg.Body = data[offset:]
	return msg, nil
}

func size(msg *transport.Message) (s int) {
	s = 2
	for k, v := range msg.Header {
		s += 2 + len(k) + 2 + len(v)
	}
	s += len(msg.Body)
	return
}

func unmarshalUint16(buf []byte, offset *int, maxLen int) (uint16, bool) {
	needLen := 2
	if *offset+needLen > maxLen {
		return 0, false
	}
	ret := binary.LittleEndian.Uint16(buf[*offset : *offset+needLen])
	*offset = *offset + needLen
	return ret, true
}

func unmarshalString(buf []byte, offset *int, needLen, maxLen int) (string, bool) {
	if *offset+needLen > maxLen {
		return "", false
	}
	s := buf[*offset : *offset+needLen]
	ret := string(s[:clen(s)])
	*offset = *offset + needLen
	return ret, true
}

func unmarshalString2(buf []byte, offset *int, maxLen int) (string, bool) {
	needLen, ok := unmarshalUint16(buf, offset, maxLen)
	if !ok {
		return "", false
	}
	return unmarshalString(buf, offset, int(needLen), maxLen)
}

func unmarshalSliceByte(buf []byte, offset *int, needLen, maxLen int) ([]byte, bool) {
	if *offset+needLen > maxLen {
		return []byte{}, false
	}
	s := buf[*offset : *offset+needLen]
	*offset = *offset + needLen
	return s, true
}

func marshalUint16(buf []byte, offset *int, value uint16) bool {
	needLen := 2
	if *offset+needLen > len(buf) {
		return false
	}
	binary.LittleEndian.PutUint16(buf[*offset:], value)
	*offset = *offset + needLen
	return true
}

func marshalString(buf []byte, offset *int, value string, needLen int) bool {
	realLen := len(value)
	if *offset+needLen > len(buf) || realLen > needLen {
		return false
	}
	copy(buf[*offset:*offset+realLen], []byte(value))
	copy(buf[*offset+realLen:*offset+needLen], zeroSlice[:])
	*offset = *offset + needLen
	return true
}

func marshalString2(buf []byte, offset *int, value string) bool {
	needLen := len(value)
	if ok := marshalUint16(buf, offset, uint16(needLen)); !ok {
		return false
	}
	return marshalString(buf, offset, value, needLen)
}

func marshalSliceByte(buf []byte, offset *int, value []byte) bool {
	needLen := len(value)
	if *offset+needLen > len(buf) {
		return false
	}
	copy(buf[*offset:*offset+needLen], value)
	*offset = *offset + needLen
	return true
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

var zeroSlice [40960]byte

package gmq

import (
	"encoding/binary"
)

const (
	MAX_UINT16 = int(^uint16(0))
)

func encodeLength(payloadLen uint64) []byte {
	if payloadLen == 0 {
		return []byte{0}
	}
	result := []byte{}
	for l := payloadLen; l > 0; {
		b := byte(l % 128)
		l = l / 128
		if l > 0 {
			b = b | 128
		}
		result = append(result, b)
	}
	return result
}

func extractLength(b []byte) (uint64, []byte) {
	var mut uint64 = 1
	var result uint64
	idx := 0
	for idx = 0; idx < len(b); idx++ {
		b := b[idx]
		result += uint64(b&127) * mut
		mut *= 128
		if b&128 == 0 {
			break
		}
	}
	return result, b[(idx + 1):]
}

func encodeString(s string) []byte {
	result := []byte(s)

	if l := len(result); l >= MAX_UINT16 {
		result = result[:MAX_UINT16-1]
	}
	return append(encodeUint16(uint16(len(result))), result...)
}

func extractString(b []byte) (string, []byte) {
	l, buf := extractUint16(b)
	if int(l) > len(buf) {
		l = uint16(len(buf))
	}
	return string(buf[:l]), buf[l:]
}

func encodeBinary(b []byte) []byte {
	l := len(b)
	if l >= MAX_UINT16 {
		l = MAX_UINT16 - 1
	}
	return append(encodeUint16(uint16(l)), b[:l]...)
}

func extractBinary(b []byte) ([]byte, []byte) {
	l, buf := extractUint16(b)
	if int(l) > len(buf) {
		l = uint16(len(buf))
	}
	return buf[:l], buf[l:]
}

func encodeUint32(i uint32) []byte {
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, i)
	return result
}

func extractUint32(b []byte) (uint32, []byte) {
	if len(b) < 4 {
		return 0, []byte{}
	}
	return binary.BigEndian.Uint32(b[:4]), b[4:]
}

func encodeUint16(i uint16) []byte {
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, i)
	return result
}

func extractUint16(b []byte) (uint16, []byte) {
	if len(b) < 2 {
		return 0, []byte{}
	}
	return binary.BigEndian.Uint16(b[:2]), b[2:]
}

func encodeBool(b bool) []byte {
	var result byte
	if b {
		result = 1
	}
	return []byte{result}
}

func extractBool(b []byte) (bool, []byte) {
	if len(b) == 0 {
		return false, []byte{}
	}
	return b[0] > 0, b[1:]
}

func extractByte(b []byte) (byte, []byte) {
	if len(b) == 0 {
		return 0, []byte{}
	}
	return b[0], b[1:]
}

func encodeKeyValue(v KeyValue) []byte {
	return append(encodeString(v.Key), encodeString(v.Value)...)
}

func extractKeyValue(b []byte) (KeyValue, []byte) {
	var result KeyValue
	result.Key, b = extractString(b)
	result.Value, b = extractString(b)
	return result, b
}

package gmq

import (
	"errors"
	"fmt"
)

type PacketType byte

const (
	_ PacketType = iota << 4
	CONNECT
	CONNACK
	PUBLISH
	PUBACK
	PUBREC
	PUBREL
	PUBCOMP
	SUBSCRIBE
	SUBACK
	UNSUBSCRIBE
	UNSUBACK
	PINGREQ
	PINGRESP
	DISCONNECT
	AUTH
)

type ReasonCode byte

const (
	REASON_SUCCESS                       ReasonCode = 0
	REASON_UNSPECIFIED_ERROR             ReasonCode = 128
	REASON_MALFORMED_PACKET              ReasonCode = 129
	REASON_PROTOCOL_ERROR                ReasonCode = 130
	REASON_IMPLEMENTATION_SPECIFIC_ERROR ReasonCode = 131
	REASON_UNSUPPORTED_PROTOCOL_VERSION  ReasonCode = 132
	REASON_CLIENT_IDENTIFIER_NOT_VALID   ReasonCode = 133
	REASON_BAD_USERNAME_OR_PASSWORD      ReasonCode = 134
	REASON_NOT_AUTHRIZED                 ReasonCode = 135
	REASON_SERVER_UNAVAILABLE            ReasonCode = 136
	REASON_SERVER_BUSY                   ReasonCode = 137
	REASON_BANNED                        ReasonCode = 138
	REASON_BAD_AUTHENTICATION_METHOD     ReasonCode = 140
	REASON_TOPIC_NAME_INVALID            ReasonCode = 144
	REASON_PACKET_TOO_LARGE              ReasonCode = 149
	REASON_QUOTA_EXCEEDED                ReasonCode = 151
	REASON_PAYLOAD_FORMAT_UNVALID        ReasonCode = 153
	REASON_RETAIN_NOT_SUPPORTED          ReasonCode = 154
	REASON_QOS_NOT_SUPPORTED             ReasonCode = 155
	REASON_USE_ANOTHER_SERVER            ReasonCode = 156
	REASON_SERVER_MOVED                  ReasonCode = 157
	REASON_CONNECTION_RATE_EXCEEDED      ReasonCode = 159
)

type QOSLevel byte

const (
	QOS_AT_MOST_ONCE QOSLevel = iota
	QOS_AT_LEAST_ONCE
	QOS_EXACTLY_ONCE
)

// ===== connect =====

type Connect struct {
	Name    string
	Version byte
	// flags
	WillRetain bool
	QOS        QOSLevel
	CleanStart bool

	KeepAlive uint16
	// properties
	Properties []Property
	// payload
	ClientId       string
	WillProperties []Property
	WillTopic      string
	WillPayload    []byte
	Username       string
	Password       []byte
}

func ConnectFromBytes(b []byte) (Connect, error) {
	var c Connect

	if len(b) == 0 {
		return c, errors.New("connect: buffer is empty") // TODO error with reason-code
	}
	if ptype := PacketType(b[0] & 0xf0); ptype != CONNECT {
		return c, fmt.Errorf("connect: control type is not CONNECT (%v)", ptype) // TODO error with reason-code
	}

	length, buf := extractLength(b[1:])
	if lb := uint64(len(buf)); lb != length {
		return c, fmt.Errorf("connect: give length %d is not equal to package size %d", length, lb) // TODO error with reason-code
	}

	c.Name, buf = extractString(buf)
	c.Version, buf = extractByte(buf)
	var flags byte
	flags, buf = extractByte(buf)
	useUsername := flags&128 == 128
	usePassword := flags&64 == 64
	c.WillRetain = flags&32 == 32
	c.QOS = QOSLevel(flags & 24)
	useWill := flags&4 == 4
	c.CleanStart = flags&2 == 2

	c.KeepAlive, buf = extractUint16(buf)
	c.Properties, buf = extractProperties(buf)
	c.ClientId, buf = extractString(buf)
	if useWill {
		c.WillProperties, buf = extractProperties(buf)
		c.WillTopic, buf = extractString(buf)
		c.WillPayload, buf = extractBinary(buf)
	}
	if useUsername {
		c.Username, buf = extractString(buf)
	}
	if usePassword {
		c.Password, buf = extractBinary(buf)
	}

	if rest := len(buf); rest != 0 {
		return c, fmt.Errorf("connect: too many bytes %d", rest) // TODO error with reason-code
	}

	// TODO validate Connect: 3.1.4 CONNECT Actions / https://docs.oasis-open.org/mqtt/mqtt/v5.0/os/mqtt-v5.0-os.html#_Toc3901033

	return c, nil
}

func (c Connect) UseUsername() bool {
	return c.Username != ""
}

func (c Connect) UsePassword() bool {
	return len(c.Password) > 0
}

func (c Connect) UseWill() bool {
	return len(c.WillProperties) > 0 && c.WillTopic != "" && len(c.WillPayload) > 0
}

func (c Connect) version() []byte {
	return append(encodeString(c.Name), c.Version)
}

func (c Connect) flags() byte {
	result := byte(c.QOS) << 3
	if c.UseUsername() {
		result |= 128
	}
	if c.UsePassword() {
		result |= 64
	}
	if c.WillRetain {
		result |= 32
	}
	if c.UseWill() {
		result |= 4
	}
	if c.CleanStart {
		result |= 2
	}
	return result
}

func (c Connect) payload() []byte {
	result := make([]byte, 0)
	if c.ClientId != "" {
		result = append(result, encodeString(c.ClientId)...)
	}
	if len(c.WillProperties) > 0 {
		result = append(result, propertyBytes(c.WillProperties...)...)
	}
	if c.WillTopic != "" {
		result = append(result, encodeString(c.WillTopic)...)
	}
	if len(c.WillPayload) > 0 {
		result = append(result, encodeBinary(c.WillPayload)...)
	}
	if c.Username != "" {
		result = append(result, encodeString(c.Username)...)
	}
	if len(c.Password) > 0 {
		result = append(result, encodeBinary(c.Password)...)
	}
	return result
}

func (c Connect) Bytes() []byte {
	header := []byte{byte(CONNECT)}

	var payload []byte
	payload = append(payload, c.version()...)
	payload = append(payload, c.flags())
	payload = append(payload, encodeUint16(c.KeepAlive)...)
	payload = append(payload, propertyBytes(c.Properties...)...)
	payload = append(payload, c.payload()...)

	header = append(header, encodeLength(uint64(len(payload)))...)
	return append(header, payload...)
}

// ===== connack =====

type ConnAck struct {
	SessionPresent bool
	ReasonCode     ReasonCode
	Properties     []Property
}

func ConnAckFromBytes(b []byte) (ConnAck, error) {
	var ca ConnAck

	if len(b) == 0 {
		return ca, errors.New("connack: buffer is empty") // TODO error with reason-code
	}
	if ptype := PacketType(b[0] & 0xf0); ptype != CONNACK {
		return ca, fmt.Errorf("connack: control type is not CONNACK (%v)", ptype) // TODO error with reason-code
	}

	length, buf := extractLength(b[1:])
	if lb := uint64(len(buf)); lb != length {
		return ca, fmt.Errorf("connack: give length %d is not equal to package size %d", length, lb) // TODO error with reason-code
	}

	ca.SessionPresent, buf = extractBool(buf)
	var reasonCode byte
	reasonCode, buf = extractByte(buf)
	ca.ReasonCode = ReasonCode(reasonCode)
	ca.Properties, buf = extractProperties(buf)

	if rest := len(buf); rest != 0 {
		return ca, fmt.Errorf("connack: too many bytes %d", rest) // TODO error with reason-code
	}

	return ca, nil
}

func (ca ConnAck) Bytes() []byte {
	header := []byte{byte(CONNACK)}

	var payload []byte
	payload = append(payload, encodeBool(ca.SessionPresent)...)
	payload = append(payload, byte(ca.ReasonCode))
	payload = append(payload, propertyBytes(ca.Properties...)...)

	header = append(header, encodeLength(uint64(len(payload)))...)
	return append(header, payload...)
}

// ===== disconnect =====

type Disconnect struct {
	ReasonCode ReasonCode
	Properties []Property
}

func DisconnectFromBytes(b []byte) (Disconnect, error) {
	var d Disconnect

	if len(b) == 0 {
		return d, errors.New("disconnect: buffer is empty") // TODO error with reason-code
	}
	if ptype := PacketType(b[0] & 0xf0); ptype != DISCONNECT {
		return d, fmt.Errorf("disconnect: control type is not DISCONNECT (%v)", ptype) // TODO error with reason-code
	}

	length, buf := extractLength(b[1:])
	if lb := uint64(len(buf)); lb != length {
		return d, fmt.Errorf("disconnect: given length %d is not equal to package size %d", length, lb) // TODO error with reason-code
	}

	var reasonCode byte
	reasonCode, buf = extractByte(buf)
	d.ReasonCode = ReasonCode(reasonCode)
	d.Properties, buf = extractProperties(buf)

	if rest := len(buf); rest != 0 {
		return d, fmt.Errorf("disconnect: too many bytes %d", rest) // TODO error with reason-code
	}

	return d, nil
}

func (d Disconnect) Bytes() []byte {
	header := []byte{byte(DISCONNECT)}

	var payload []byte
	payload = append(payload, byte(d.ReasonCode))
	payload = append(payload, propertyBytes(d.Properties...)...)

	header = append(header, encodeLength(uint64(len(payload)))...)
	return append(header, payload...)
}

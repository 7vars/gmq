package gmq

type PropertyType byte

const (
	PROP_PAYLOAD_FORMAT                  PropertyType = 1
	PROP_MESSAGE_EXPIRY_INTERVAL         PropertyType = 2
	PROP_CONTENT_TYPE                    PropertyType = 3
	PROP_RESPONSE_TOPIC                  PropertyType = 8
	PROP_CORRELATION_DATA                PropertyType = 9
	PROP_SESSION_EXPIRY_INTERVAL         PropertyType = 17
	PROP_ASSIGNED_CLIENT_ID              PropertyType = 18
	PROP_SERVER_KEEP_ALIVE               PropertyType = 19
	PROP_AUTHENTICATION_METHOD           PropertyType = 21
	PROP_AUTHENTICATION_DATA             PropertyType = 22
	PROP_REQUEST_PROBLEM_INFO            PropertyType = 23
	PROP_DELAY_INTERVAL                  PropertyType = 24
	PROP_REQUEST_RESPONSE_INFO           PropertyType = 25
	PROP_RESPONSE_INFO                   PropertyType = 26
	PROP_SERVER_REFERENCE                PropertyType = 28
	PROP_REASON_STRING                   PropertyType = 31
	PROP_RECEIVE_MAX                     PropertyType = 33
	PROP_ALIAS_MAX                       PropertyType = 34
	PROP_MAX_QOS                         PropertyType = 36
	PROP_RETAIN_AVAILABLE                PropertyType = 37
	PROP_USER_PROPERTY                   PropertyType = 38
	PROP_MAX_PACKET_SIZE                 PropertyType = 39
	PROP_WILDCARD_SUBSCRIPTION_AVAILABLE PropertyType = 40
	PROP_SUBSCRIPTION_IDS_AVAILABLE      PropertyType = 41
	PROP_SHARED_SUBSCTION_AVAILABLE      PropertyType = 42
)

type KeyValue struct {
	Key   string
	Value string
}

type Property struct {
	Type  PropertyType
	Value any
}

func (p Property) Bytes() []byte {
	result := []byte{byte(p.Type)}
	switch val := p.Value.(type) {
	case bool:
		result = append(result, encodeBool(val)...)
	case byte:
		result = append(result, val)
	case uint16:
		result = append(result, encodeUint16(val)...)
	case uint32:
		result = append(result, encodeUint32(val)...)
	case []byte:
		result = append(result, encodeBinary(val)...)
	case string:
		result = append(result, encodeString(val)...)
	case KeyValue:
		result = append(result, encodeKeyValue(val)...)
	default:
		return []byte{}
	}
	return result
}

func propertyBytes(properties ...Property) []byte {
	result := make([]byte, 0)
	for _, prop := range properties {
		result = append(result, prop.Bytes()...)
	}
	l := uint64(len(result))
	return append(encodeLength(l), result...)
}

func extractProperties(b []byte) ([]Property, []byte) {
	length, b := extractLength(b)
	buf := b[:length]
	result := []Property{}

	for len(buf) != 0 {
		var prop Property
		var btype byte
		btype, buf = extractByte(buf)
		switch ptype := PropertyType(btype); ptype {
		case PROP_PAYLOAD_FORMAT:
			prop.Type = ptype
			prop.Value, buf = extractBool(buf)
		case PROP_MESSAGE_EXPIRY_INTERVAL:
			prop.Type = ptype
			prop.Value, buf = extractUint32(buf)
		case PROP_CONTENT_TYPE:
			prop.Type = ptype
			prop.Value, buf = extractString(buf)
		case PROP_RESPONSE_TOPIC:
			prop.Type = ptype
			prop.Value, buf = extractString(buf)
		case PROP_CORRELATION_DATA:
			prop.Type = ptype
			prop.Value, buf = extractBinary(buf)
		case PROP_SESSION_EXPIRY_INTERVAL:
			prop.Type = ptype
			prop.Value, buf = extractUint32(buf)
		case PROP_ASSIGNED_CLIENT_ID:
			prop.Type = ptype
			prop.Value, buf = extractString(buf)
		case PROP_SERVER_KEEP_ALIVE:
			prop.Type = ptype
			prop.Value, buf = extractUint16(buf)
		case PROP_AUTHENTICATION_METHOD:
			prop.Type = ptype
			prop.Value, buf = extractString(buf)
		case PROP_AUTHENTICATION_DATA:
			prop.Type = ptype
			prop.Value, buf = extractBinary(buf)
		case PROP_REQUEST_PROBLEM_INFO:
			prop.Type = ptype
			prop.Value, buf = extractBool(buf)
		case PROP_DELAY_INTERVAL:
			prop.Type = ptype
			prop.Value, buf = extractUint32(buf)
		case PROP_REQUEST_RESPONSE_INFO:
			prop.Type = ptype
			prop.Value, buf = extractBool(buf)
		case PROP_RESPONSE_INFO:
			prop.Type = ptype
			prop.Value, buf = extractString(buf)
		case PROP_SERVER_REFERENCE:
			prop.Type = ptype
			prop.Value, buf = extractString(buf)
		case PROP_REASON_STRING:
			prop.Type = ptype
			prop.Value, buf = extractString(buf)
		case PROP_RECEIVE_MAX:
			prop.Type = ptype
			prop.Value, buf = extractUint16(buf)
		case PROP_ALIAS_MAX:
			prop.Type = ptype
			prop.Value, buf = extractUint16(buf)
		case PROP_MAX_QOS:
			prop.Type = ptype
			prop.Value, buf = extractBool(buf)
		case PROP_RETAIN_AVAILABLE:
			prop.Type = ptype
			prop.Value, buf = extractBool(buf)
		case PROP_USER_PROPERTY:
			prop.Type = ptype
			prop.Value, buf = extractKeyValue(buf)
		case PROP_MAX_PACKET_SIZE:
			prop.Type = ptype
			prop.Value, buf = extractUint32(buf)
		case PROP_WILDCARD_SUBSCRIPTION_AVAILABLE:
			prop.Type = ptype
			prop.Value, buf = extractBool(buf)
		case PROP_SUBSCRIPTION_IDS_AVAILABLE:
			prop.Type = ptype
			prop.Value, buf = extractBool(buf)
		case PROP_SHARED_SUBSCTION_AVAILABLE:
			prop.Type = ptype
			prop.Value, buf = extractBool(buf)
		default:
			return []Property{}, []byte{}
		}
		result = append(result, prop)
	}

	return result, b[length:]
}

func SessionExpiryIntervalProperty(interval uint32) Property {
	return Property{
		Type:  PROP_SESSION_EXPIRY_INTERVAL,
		Value: interval,
	}
}

func ReceiveMaxProperty(max uint16) Property {
	return Property{
		Type:  PROP_RECEIVE_MAX,
		Value: max,
	}
}

func MaxPacketSizeProperty(size uint32) Property {
	return Property{
		Type:  PROP_MAX_PACKET_SIZE,
		Value: size,
	}
}

func AliasMaxProperty(max uint16) Property {
	return Property{
		Type:  PROP_ALIAS_MAX,
		Value: max,
	}
}

func RequestResponseInfoProperty(b bool) Property {
	return Property{
		Type:  PROP_REQUEST_RESPONSE_INFO,
		Value: b,
	}
}

func RequestProblemInfoProperty(b bool) Property {
	return Property{
		Type:  PROP_REQUEST_PROBLEM_INFO,
		Value: b,
	}
}

func UserPropertyProperty(key, value string) Property {
	return Property{
		Type: PROP_USER_PROPERTY,
		Value: KeyValue{
			Key:   key,
			Value: value,
		},
	}
}

func AuthenticationMethodProperty(method string) Property {
	return Property{
		Type:  PROP_AUTHENTICATION_METHOD,
		Value: method,
	}
}

func AuthenticationDataProperty(data []byte) Property {
	return Property{
		Type:  PROP_AUTHENTICATION_DATA,
		Value: data,
	}
}

func DelayIntervalProperty(interval uint32) Property {
	return Property{
		Type:  PROP_DELAY_INTERVAL,
		Value: interval,
	}
}

func PayloadFormatProperty(b bool) Property {
	return Property{
		Type:  PROP_PAYLOAD_FORMAT,
		Value: b,
	}
}

func MessageExpiryIntervalProperty(interval uint32) Property {
	return Property{
		Type:  PROP_MESSAGE_EXPIRY_INTERVAL,
		Value: interval,
	}
}

func ContentTypeProperty(contentType string) Property {
	return Property{
		Type:  PROP_CONTENT_TYPE,
		Value: contentType,
	}
}

func ResponseTopicProperty(topic string) Property {
	return Property{
		Type:  PROP_RESPONSE_TOPIC,
		Value: topic,
	}
}

func CorrelationDataProperty(b []byte) Property {
	return Property{
		Type:  PROP_CORRELATION_DATA,
		Value: b,
	}
}

func MaxQOSProperty(b bool) Property {
	return Property{
		Type:  PROP_MAX_QOS,
		Value: b,
	}
}

func RetainAvailableProperty(b bool) Property {
	return Property{
		Type:  PROP_MAX_QOS,
		Value: b,
	}
}

func AssignedClientIdProperty(clientId string) Property {
	return Property{
		Type:  PROP_ASSIGNED_CLIENT_ID,
		Value: clientId,
	}
}

func ReasonStringProperty(reason string) Property {
	return Property{
		Type:  PROP_ASSIGNED_CLIENT_ID,
		Value: reason,
	}
}

func WildcardSubscriptionAvailableProperty(b bool) Property {
	return Property{
		Type:  PROP_WILDCARD_SUBSCRIPTION_AVAILABLE,
		Value: b,
	}
}

func SubscriptionIdsAvailableProperty(b bool) Property {
	return Property{
		Type:  PROP_SUBSCRIPTION_IDS_AVAILABLE,
		Value: b,
	}
}

func ShardSubscriptionAvailableProperty(b bool) Property {
	return Property{
		Type:  PROP_SHARED_SUBSCTION_AVAILABLE,
		Value: b,
	}
}

func ServerKeepAliveProperty(v uint16) Property {
	return Property{
		Type:  PROP_SERVER_KEEP_ALIVE,
		Value: v,
	}
}

func ResponseInfoProperty(info string) Property {
	return Property{
		Type:  PROP_RESPONSE_INFO,
		Value: info,
	}
}

func ServerReferenceProperty(ref string) Property {
	return Property{
		Type:  PROP_SERVER_REFERENCE,
		Value: ref,
	}
}

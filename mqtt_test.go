package gmq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	connect := Connect{
		Name:      "MQTT",
		Version:   5,
		KeepAlive: 30,
		Properties: []Property{
			SessionExpiryIntervalProperty(0),
		},
		ClientId: "mqttx_5b75ccbb",
		Username: "test",
		Password: []byte("test123"),
	}

	buf := connect.Bytes()
	assert.True(t, len(buf) > 0)

	con, err := ConnectFromBytes(buf)
	assert.NoError(t, err)
	assert.Equal(t, connect, con)
}

func TestConnAck(t *testing.T) {
	connack := ConnAck{
		SessionPresent: false,
		ReasonCode:     REASON_SUCCESS,
		Properties: []Property{
			AliasMaxProperty(10),
			ReceiveMaxProperty(10),
		},
	}

	buf := connack.Bytes()
	assert.True(t, len(buf) > 0)

	cak, err := ConnAckFromBytes(buf)
	assert.NoError(t, err)
	assert.Equal(t, connack, cak)
}

func TestDisconnect(t *testing.T) {
	disconnect := Disconnect{
		ReasonCode:     REASON_SUCCESS,
		Properties: []Property{},
	}

	buf := disconnect.Bytes()
	assert.True(t, len(buf) > 0)

	dis, err := DisconnectFromBytes(buf)
	assert.NoError(t, err)
	assert.Equal(t, disconnect, dis)
}
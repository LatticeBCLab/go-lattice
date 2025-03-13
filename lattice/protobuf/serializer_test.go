package protobuf

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializer(t *testing.T) {
	proto := `syntax = "proto3";
	package contract.v1;
	import "google/protobuf/wrappers.proto";
	
	message Person {
		string name = 1;
		google.protobuf.BoolValue man = 2;
	}
	`

	fd := MakeFileDescriptor(strings.NewReader(proto))
	assert.NotNil(t, fd)

	t.Run("Marshal message", func(t *testing.T) {
		json := `{"name": "Jack", "man": null}`
		bs, err := MarshallMessage(fd, json)
		assert.Nil(t, err)
		t.Log(string(bs))
	})

	t.Run("Unmarshal message", func(t *testing.T) {
		json := `{"name": "Jack", "man": null}`
		bs, err := MarshallMessage(fd, json)
		assert.Nil(t, err)
		str, err := UnmarshallMessage(fd, bs)
		assert.Nil(t, err)
		t.Log(str)
	})
}

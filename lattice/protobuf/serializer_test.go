package protobuf

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializer(t *testing.T) {
	proto := `syntax = "proto3";
	import "google/protobuf/wrappers.proto";
	
	message Student {
		google.protobuf.StringValue id = 1;
		google.protobuf.StringValue name = 2;
		google.protobuf.Int32Value age = 3;
		google.protobuf.DoubleValue balance = 4;
	}
	`

	fd := MakeFileDescriptor(strings.NewReader(proto))
	assert.NotNil(t, fd)

	t.Run("Repeated register wrapper proto", func(t *testing.T) {
		fd = MakeFileDescriptor(strings.NewReader(proto))
		assert.NotNil(t, fd)
	})

	t.Run("Marshal message", func(t *testing.T) {
		json := `{"name": "Jack", "age": null}`
		bs, err := MarshallMessage(fd, json)
		assert.Nil(t, err)
		t.Log(string(bs))
	})

	t.Run("Unmarshal message", func(t *testing.T) {
		json := `{"name": "Jack", "age": null}`
		bs, err := MarshallMessage(fd, json)
		assert.Nil(t, err)
		str, err := UnmarshallMessage(fd, bs)
		assert.Nil(t, err)
		t.Log(str)
	})

}

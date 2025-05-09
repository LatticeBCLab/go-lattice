package protobuf

import (
	"google.golang.org/protobuf/reflect/protoregistry"
	"io"
	"strings"
	"sync"

	"github.com/bufbuild/protocompile/parser"
	"github.com/bufbuild/protocompile/reporter"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

var registerWrapperOnce sync.Once

// MakeFileDescriptor 生成proto的文件描述
//
// Parameters:
//   - reader io.Reader
//
// Returns:
//   - pref.FileDescriptor
func MakeFileDescriptor(reader io.Reader) pref.FileDescriptor {
	errHandler := reporter.NewHandler(nil)
	ast, err := parser.Parse("example.proto", reader, errHandler)
	if err != nil {
		panic(err)
	}

	result, err := parser.ResultFromAST(ast, true, errHandler)
	if err != nil {
		panic(err)
	}

	fdp := result.FileDescriptorProto()

	resolver := protoregistry.GlobalFiles
	registerWrapperOnce.Do(func() {
		if err = resolver.RegisterFile(MakeWrapperProtoFileDescriptor()); err != nil {
			panic(err)
		}
	})

	// get FileDescriptor
	fd, err := protodesc.NewFile(fdp, resolver)
	if err != nil {
		panic(err)
	}
	return fd
}

func MakeWrapperProtoFileDescriptor() pref.FileDescriptor {
	reader := strings.NewReader(wrapperProto)
	errHandler := reporter.NewHandler(nil)
	ast, err := parser.Parse("google/protobuf/wrappers.proto", reader, errHandler)
	if err != nil {
		panic(err)
	}

	result, err := parser.ResultFromAST(ast, true, errHandler)
	if err != nil {
		panic(err)
	}

	fdp := result.FileDescriptorProto()

	// get FileDescriptor
	fd, err := protodesc.NewFile(fdp, nil)
	if err != nil {
		panic(err)
	}
	return fd
}

// MarshallMessage 序列化
//
// Parameters:
//   - fd pref.FileDescriptor
//   - json string
//
// Returns:
//   - []byte
//   - error
func MarshallMessage(fd pref.FileDescriptor, json string) ([]byte, error) {
	messageDescriptor := fd.Messages().Get(0)
	// messageDescriptor := fd.Messages().ByName(pref.Name(name))
	message := dynamicpb.NewMessage(messageDescriptor)

	if err := protojson.Unmarshal([]byte(json), message); err != nil {
		return nil, err
	}

	bytes, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

// UnmarshallMessage 反序列化
//
// Parameters:
//   - fd pref.FileDescriptor
//   - data []byte
//
// Returns:
//   - string
//   - error
func UnmarshallMessage(fd pref.FileDescriptor, data []byte) (string, error) {
	messageDescriptor := fd.Messages().Get(0)
	// messageDescriptor := fd.Messages().ByName(pref.Name(name))
	message := dynamicpb.NewMessage(messageDescriptor)

	err := proto.Unmarshal(data, message)
	if err != nil {
		return "", err
	}

	options := protojson.MarshalOptions{
		EmitUnpopulated: true,
	}
	jsonBytes, err := options.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

const wrapperProto = `syntax = "proto3";

package google.protobuf;

option cc_enable_arenas = true;
option go_package = "google.golang.org/protobuf/types/known/wrapperspb";
option java_package = "com.google.protobuf";
option java_outer_classname = "WrappersProto";
option java_multiple_files = true;
option objc_class_prefix = "GPB";
option csharp_namespace = "Google.Protobuf.WellKnownTypes";

message DoubleValue {
    // The double value.
    double value = 1;
}

// Wrapper message for float.
//
// The JSON representation for FloatValue is JSON number.
message FloatValue {
    // The float value.
    float value = 1;
}

// Wrapper message for int64.
//
// The JSON representation for Int64Value is JSON string.
message Int64Value {
    // The int64 value.
    int64 value = 1;
}

// Wrapper message for uint64.
//
// The JSON representation for UInt64Value is JSON string.
message UInt64Value {
    // The uint64 value.
    uint64 value = 1;
}

// Wrapper message for int32.
//
// The JSON representation for Int32Value is JSON number.
message Int32Value {
    // The int32 value.
    int32 value = 1;
}

// Wrapper message for uint32.
//
// The JSON representation for UInt32Value is JSON number.
message UInt32Value {
    // The uint32 value.
    uint32 value = 1;
}

// Wrapper message for bool.
//
// The JSON representation for BoolValue is JSON true and false.
message BoolValue {
    // The bool value.
    bool value = 1;
}

// Wrapper message for string.
//
// The JSON representation for StringValue is JSON string.
message StringValue {
    // The string value.
    string value = 1;
}

// Wrapper message for bytes.
//
// The JSON representation for BytesValue is JSON string.
message BytesValue {
    // The bytes value.
    bytes value = 1;
}
`

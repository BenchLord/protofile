package protofile

import (
	"fmt"
	"testing"
)

func TestProtofile(t *testing.T) {
	f, err := New("test.proto")
	if err != nil {
		t.Fatal(err)
	}
	s := fmt.Sprintf("%#v", f.GetServices())
	if s != "[]protofile.Service{protofile.Service{name:\"Test\", methods:[]protofile.Method{protofile.Method{name:\"getAll\", req:protofile.Message{name:\"Empty\", fields:map[string]protofile.Field{}}, isReqStream:false, res:protofile.Message{name:\"Thing\", fields:map[string]protofile.Field{\"awards\":protofile.Field{name:\"awards\", fieldType:\"string\", repeated:true, number:3}, \"id\":protofile.Field{name:\"id\", fieldType:\"Id\", repeated:false, number:1}, \"name\":protofile.Field{name:\"name\", fieldType:\"string\", repeated:false, number:2}}}, isResStream:true}, protofile.Method{name:\"getOne\", req:protofile.Message{name:\"Id\", fields:map[string]protofile.Field{\"value\":protofile.Field{name:\"value\", fieldType:\"int64\", repeated:false, number:1}}}, isReqStream:false, res:protofile.Message{name:\"Thing\", fields:map[string]protofile.Field{\"awards\":protofile.Field{name:\"awards\", fieldType:\"string\", repeated:true, number:3}, \"id\":protofile.Field{name:\"id\", fieldType:\"Id\", repeated:false, number:1}, \"name\":protofile.Field{name:\"name\", fieldType:\"string\", repeated:false, number:2}}}, isResStream:false}}}}" {
		t.Error("Did not get expected output")
	}
}

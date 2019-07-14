package protofile

import (
	"testing"
)

func TestProtofile(t *testing.T) {
	f, err := New("test.proto")
	if err != nil {
		t.Fatal(err)
	}
	service := f.GetServices()[0]

	if service.GetName() != "Test" {
		t.Fatalf("Failed to get service name.\nWanted: Test\nGot:%s", service.GetName())
	}

	message := f.GetMessages()
	if _, ok := message["Empty"]; !ok {
		t.Fatal("Failed to find Empty message")
	}
	if _, ok := message["Thing"]; !ok {
		t.Fatal("Failed to find Thing message")
	}
	if _, ok := message["Id"]; !ok {
		t.Fatal("Failed to find Id message")
	}
}

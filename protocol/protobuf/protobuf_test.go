package protobuf

import (
	"fmt"
	"strings"
	"testing"
)

var _testProto = `
syntax = "proto3";

import "google/protobuf/any.proto";
import "google/protobuf/duration.proto";

// option java_package = "com.iqiyi.itv.proto";
// option go_package = "pb";

enum Gender {
    MALE = 0;
    FEMALE = 1;
}

enum PetType {
	DOG = 0;
	CAT = 1;
	SNAKE = 2;
}
message Person {
    string name = 1;
    uint32 age = 2;
    Gender gender = 3;
    repeated string hobbies = 4; 
    repeated google.protobuf.Any desc = 6;
    google.protobuf.Duration online_time = 5;
	repeated Pet pets = 7;
    reserved 10 to 16;
    reserved "nick_name";
}

message Pet {
	PetType type = 1;
	string name = 2;
}

message Integer {
    int32 value = 1;
}

message PersonGetRequest {
	uint32 id = 1;
}

service PersonService {
  rpc GetByID (PersonGetRequest) returns (Person);
}

`

var _proto *ProtobufProtocol

func init() {
	_proto, _ = NewProtobufProtocol()
}

func TestParse(t *testing.T) {
	err := _proto.Parse("test", strings.NewReader(_testProto))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(_proto.parser)
}

func TestMarshal(t *testing.T) {
	err := _proto.Parse("test", strings.NewReader(_testProto))
	if err != nil {
		t.Fatal(err)
	}

	b, err := _proto.Marshal("Pet", map[string]interface{}{
		"type": 1,
		"name": "xiaoqiang",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("message: [% X]", b)
}

package str

import (
	"reflect"
	"testing"
)

var str = "hello,你好！"
var bytes = []byte(str)

func Test_Bytes(t *testing.T) {
	bs := Bytes(str)
	if !reflect.DeepEqual(bs, bytes) {
		t.Fail()
	}
}

func Test_String(t *testing.T) {
	s := String(bytes)
	if s != str {
		t.Fail()
	}
}

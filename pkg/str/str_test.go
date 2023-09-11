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

func Test_Connect(t *testing.T) {
	s := []string{"hello", ",", "你好", "！", " "}
	str, err := Connect(s...)
	if err != nil {
		t.Fatal(err)
	}
	if str != "hello,你好！ " {
		t.Fail()
	}
}

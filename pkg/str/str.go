package str

import (
	"reflect"
	"strings"
	"unsafe"
)

func Bytes(str string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))

	var bytes []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	sh.Data = stringHeader.Data
	sh.Len = stringHeader.Len
	sh.Cap = stringHeader.Len

	return bytes
}

func Connect(s ...string) (string, error) {
	var res strings.Builder

	for _, str := range s {
		_, _ = res.WriteString(str)
	}

	return res.String(), nil
}

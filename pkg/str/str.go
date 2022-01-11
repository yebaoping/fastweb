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

func String(bytes []byte) string {
	var sb strings.Builder
	_, err := sb.Write(bytes)
	if err != nil {
		return string(bytes)
	}

	return sb.String()
}

func Connect(s ...string) (string, error) {
	var res strings.Builder

	for _, str := range s {
		_, err := res.WriteString(str)
		if err != nil {
			return "", err
		}
	}

	return res.String(), nil
}

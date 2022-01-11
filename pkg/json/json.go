package json

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func MarshalToString(data interface{}) (string, error) {
	return json.MarshalToString(data)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func UnmarshalFromString(data string, v interface{}) error {
	return json.UnmarshalFromString(data, v)
}

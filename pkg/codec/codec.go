package codec

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
)

var (
	DefaultMarshaler = NewJsoniterMarshaler()
)

type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type jsonMarshaler struct{}

func NewJsonMarshaler() Marshaler {
	return &jsonMarshaler{}
}

func (j jsonMarshaler) Marshal(v interface{}) ([]byte, error) {
	bytes, err := json.Marshal(v)
	return bytes, err
}

func (j jsonMarshaler) Unmarshal(bytes []byte, v interface{}) error {
	return json.Unmarshal(bytes, v)
}

type jsoniterMarshaler struct {
	json jsoniter.API
}

func NewJsoniterMarshaler() Marshaler {
	j := jsoniter.ConfigCompatibleWithStandardLibrary
	return &jsoniterMarshaler{j}
}

func (j jsoniterMarshaler) Marshal(v interface{}) ([]byte, error) {
	return j.json.Marshal(v)
}

func (j jsoniterMarshaler) Unmarshal(bytes []byte, v interface{}) error {
	return j.json.Unmarshal(bytes, v)
}

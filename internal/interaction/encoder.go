package interaction

import (
	"bytes"
	"encoding/gob"

	"github.com/openimsdk/tools/errs"
)

type Encoder interface {
	Encode(data interface{}) ([]byte, error)
	Decode(encodeData []byte, decodeData interface{}) error
}

type GobEncoder struct {
}

func NewGobEncoder() *GobEncoder {
	return &GobEncoder{}
}
func (g *GobEncoder) Encode(data interface{}) ([]byte, error) {
	buff := bytes.Buffer{}
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
func (g *GobEncoder) Decode(encodeData []byte, decodeData interface{}) error {
	buff := bytes.NewBuffer(encodeData)
	dec := gob.NewDecoder(buff)
	err := dec.Decode(decodeData)
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

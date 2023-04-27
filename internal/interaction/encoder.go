// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package interaction

import (
	"bytes"
	"encoding/gob"
	"open_im_sdk/pkg/utils"
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
		return utils.Wrap(err, "")
	}
	return nil
}

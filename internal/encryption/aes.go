package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var ErrLength = errors.New("key length error,too short")
var ErrPaddingNotMatch = errors.New("padding not match")

type AESEncryptor struct {
}

func NewAESEncryptor() Encryptor {
	return &AESEncryptor{}
}

func (a *AESEncryptor) pKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padText...)
}

func (a *AESEncryptor) pKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unPadding := int(origData[length-1])
	if unPadding > length {
		return nil, ErrPaddingNotMatch
	} else {
		return origData[:(length - unPadding)], nil
	}
}

func (a *AESEncryptor) Encryption(rawData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(key) < blockSize {
		return nil, ErrLength
	}
	rawData = a.pKCS5Padding(rawData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(rawData))
	blockMode.CryptBlocks(encrypted, rawData)
	return encrypted, nil
}

func (a *AESEncryptor) Decryption(encryptedData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(key) < blockSize {
		return nil, ErrLength
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) //初始向量的长度必须等于块block的长度16字节
	origData := make([]byte, len(encryptedData))
	blockMode.CryptBlocks(origData, encryptedData)
	return a.pKCS5UnPadding(origData)
}

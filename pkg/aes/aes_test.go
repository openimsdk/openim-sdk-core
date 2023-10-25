package aes

import (
	"fmt"
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	data := "1                2        [111]"
	encrypt, err := EncryptByAes([]byte(data), []byte("78a07ea4e875ea75"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	aes, err := DecryptByAes(string(encrypt), []byte("78a07ea4e875ea75"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(data)
	fmt.Println(encrypt)
	fmt.Println(string(aes))
}

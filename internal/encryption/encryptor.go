package encryption

import (
	"context"
	"github.com/OpenIMSDK/protocol/sdkws"
)

type Encryptor interface {
	Encryption(rawData, key []byte) ([]byte, error)
	Decryption(encryptedData, key []byte) ([]byte, error)
}

type MessageEncryptor interface {
	Encryption(ctx context.Context, message *sdkws.MsgData, conversationID string) error
	Decryption(ctx context.Context, messageList []*sdkws.MsgData, conversationID string) error
}

package encryption

import (
	"context"
	"errors"
	"github.com/OpenIMSDK/protocol/encryption"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/sdkerrs"
	"open_im_sdk/pkg/utils"
	"sync"
	"time"
)

type ConversationEncryptor struct {
	LoginUserID string
	m           sync.Map
	a           Encryptor
}

func NewConversationEncryptor(loginUserID string) MessageEncryptor {
	return &ConversationEncryptor{LoginUserID: loginUserID}
}

func (c *ConversationEncryptor) Encryption(ctx context.Context, message *sdkws.MsgData, conversationID string) error {
	key, err := c.GetMaxVersionKey(ctx, conversationID)
	if err != nil {
		return err
	}
	encryptedData, err := c.a.Encryption(message.Content, []byte(key.Key))
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "encryption success", "message", message, "key", key)
	message.KeyVersion = key.Version
	message.Content = encryptedData
	return nil
}

func (c *ConversationEncryptor) Decryption(ctx context.Context, messageList []*sdkws.MsgData, conversationID string) error {
	for _, message := range messageList {
		log.ZDebug(ctx, "decryption", "message", message, "conversation_id", conversationID)
		if message.KeyVersion != 0 {
			if message.RecvID != c.LoginUserID && message.SendID != c.LoginUserID {
				log.ZWarn(ctx, "maybe message come from app manager", errors.New("manager message"), "message", message)
				continue
			}
			key, err := c.GetKeyByMessageVersion(ctx, conversationID, message.KeyVersion)
			if err != nil {
				log.ZWarn(ctx, "get key by message version failed", err, "message", message, "conversation_id", conversationID)
				continue
			}
			decryptedData, err := c.a.Decryption(message.Content, []byte(key.Key))
			if err != nil {
				log.ZWarn(ctx, "decryption failed", err, "message", message, "conversation_id", conversationID)
				continue
			}
			log.ZDebug(ctx, "decryption success", "message", message, "key", key)
			message.Content = decryptedData
		}
	}
	return nil
}
func (c *ConversationEncryptor) GetKeyByMessageVersion(ctx context.Context, conversationID string, version int32) (*encryption.VersionKey, error) {
	key, ok := c.m.Load(c.genConversationIDVersionKey(conversationID, version))
	if ok {
		return key.(*encryption.VersionKey), nil
	} else {
		versionKeyList, _, err := c.getEncryptionKeyFromSvr(ctx, conversationID, version)
		if err != nil {
			return nil, err
		}
		if len(versionKeyList) == 0 {
			return nil, sdkerrs.ErrMsgEncryptionKeyNotFound
		}
		return versionKeyList[0], nil
	}
}

func (c *ConversationEncryptor) GetMaxVersionKey(ctx context.Context, conversationID string) (*encryption.VersionKey, error) {
	key, ok := c.m.Load(c.genConversationIDMaxVersionKey(conversationID))
	if ok {
		return key.(*encryption.VersionKey), nil
	} else {
		_, maxKeyVersion, err := c.getEncryptionKeyFromSvr(ctx, conversationID, 0)
		if err != nil {
			return nil, err
		}
		return maxKeyVersion, nil
	}
}
func (c *ConversationEncryptor) getEncryptionKeyFromSvr(ctx context.Context, conversationID string, keyVersion int32) ([]*encryption.VersionKey, *encryption.VersionKey, error) {
	apiReq := encryption.GetEncryptionKeyReq{ConversationID: conversationID, KeyVersion: keyVersion}
	var apiResp *encryption.GetEncryptionKeyResp
	var err error
	for i := 0; i < 10; i++ {
		apiResp, err = util.CallApi[encryption.GetEncryptionKeyResp](ctx, constant.GetEncryptionKeyRouter, &apiReq)
		if err != nil {
			log.ZError(ctx, "getEncryptionKeyFromSvr failed", err, "conversation_id", conversationID, "key_version", keyVersion)
			time.Sleep(1 * time.Second)
			continue
		}
	}
	if err != nil {
		return nil, nil, err
	}
	var tempVersion encryption.VersionKey
	for _, v := range apiResp.VersionKeyList {
		if v.Version > tempVersion.Version {
			tempVersion = *v
		}
		c.m.Store(c.genConversationIDVersionKey(conversationID, v.Version), v)
	}
	c.m.Store(c.genConversationIDMaxVersionKey(conversationID), &tempVersion)
	return apiResp.VersionKeyList, &tempVersion, nil
}
func (c *ConversationEncryptor) genConversationIDVersionKey(conversationID string, version int32) string {
	return conversationID + "v" + utils.Int32ToString(version)
}
func (c *ConversationEncryptor) genConversationIDMaxVersionKey(conversationID string) string {
	return conversationID + "v" + "max"
}

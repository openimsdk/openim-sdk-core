package aes_key

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
)

type AesKey struct {
	db db_interface.DataBase
}

func NewAesKey(db db_interface.DataBase) *AesKey {
	return &AesKey{db: db}
}

type GetKeyReq struct {
	sessionType                   int32
	groupID, friendUserID, userId string
}
type GetKeyResp struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Data    Data   `json:"data"`
}
type Data struct {
	Key string `json:"key"`
}

func (k *AesKey) GetKey(sessionType int32, groupID, friendUserID, userId string) (string, error) {
	return "78a07ea4e875ea75", nil
}

//func (k *AesKey) GetKey(sessionType int32, groupID, friendUserID, userId string) (string, error) {
//	switch sessionType {
//	case constant.GroupChatType:
//		if groupID == "" {
//			return "", errors.New("Parameter error ")
//		}
//	case constant.SingleChatType:
//		if friendUserID == "" || userId == "" {
//			return "", errors.New("Parameter error ")
//		}
//	default:
//		return "", errors.New("sessionType error ")
//	}
//	resp := GetKeyResp{}
//	err := util.ApiPost(k.ctx, "", GetKeyReq{sessionType, groupID, friendUserID, userId}, &resp)
//	if err != nil {
//		return "", err
//	}
//	return resp.Data.Key, nil
//}

package indexdb

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/indexdb/temp_struct"
)

type FriendRequest struct {
	loginUserID string
}

func NewFriendRequest(loginUserID string) *FriendRequest {
	return &FriendRequest{loginUserID: loginUserID}
}

func (i FriendRequest) InsertFriendRequest(friendRequest *model_struct.LocalFriendRequest) error {
	_, err := Exec(utils.StructToJsonString(friendRequest))
	return err
}

func (i FriendRequest) DeleteFriendRequestBothUserID(fromUserID, toUserID string) error {
	_, err := Exec(fromUserID, toUserID)
	return err
}

func (i FriendRequest) UpdateFriendRequest(friendRequest *model_struct.LocalFriendRequest) error {
	tempLocalFriendRequest := temp_struct.LocalFriendRequest{
		FromUserID:    friendRequest.FromUserID,
		FromNickname:  friendRequest.FromNickname,
		FromFaceURL:   friendRequest.FromFaceURL,
		FromGender:    friendRequest.FromGender,
		ToUserID:      friendRequest.ToUserID,
		ToNickname:    friendRequest.ToNickname,
		ToFaceURL:     friendRequest.ToFaceURL,
		ToGender:      friendRequest.ToGender,
		HandleResult:  friendRequest.HandleResult,
		ReqMsg:        friendRequest.ReqMsg,
		CreateTime:    friendRequest.CreateTime,
		HandlerUserID: friendRequest.HandlerUserID,
		HandleMsg:     friendRequest.HandleMsg,
		HandleTime:    friendRequest.HandleTime,
		Ex:            friendRequest.Ex,
		AttachedInfo:  friendRequest.AttachedInfo,
	}
	_, err := Exec(utils.StructToJsonString(tempLocalFriendRequest))
	return err
}

func (i FriendRequest) GetRecvFriendApplication() (result []*model_struct.LocalFriendRequest, err error) {
	gList, err := Exec(i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalFriendRequest
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i FriendRequest) GetSendFriendApplication() (result []*model_struct.LocalFriendRequest, err error) {
	gList, err := Exec(i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalFriendRequest
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i FriendRequest) GetFriendApplicationByBothID(fromUserID, toUserID string) (*model_struct.LocalFriendRequest, error) {
	c, err := Exec(fromUserID, toUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := c.(string); ok {
			result := model_struct.LocalFriendRequest{}
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return &result, err
		} else {
			return nil, ErrType
		}
	}
}

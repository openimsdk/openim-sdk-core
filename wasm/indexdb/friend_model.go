package indexdb

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/indexdb/temp_struct"
)

type Friend struct {
	loginUserID string
}

func NewFriend(loginUserID string) *Friend {
	return &Friend{loginUserID: loginUserID}
}

func (i Friend) InsertFriend(friend *model_struct.LocalFriend) error {
	_, err := Exec(utils.StructToJsonString(friend))
	return err
}

func (i Friend) DeleteFriendDB(friendUserID string) error {
	_, err := Exec(friendUserID, i.loginUserID)
	return err
}

func (i Friend) UpdateFriend(friend *model_struct.LocalFriend) error {
	tempLocalFriend := temp_struct.LocalFriend{
		OwnerUserID:    friend.OwnerUserID,
		FriendUserID:   friend.FriendUserID,
		Remark:         friend.Remark,
		CreateTime:     friend.CreateTime,
		AddSource:      friend.AddSource,
		OperatorUserID: friend.OperatorUserID,
		Nickname:       friend.Nickname,
		FaceURL:        friend.FaceURL,
		Gender:         friend.Gender,
		PhoneNumber:    friend.PhoneNumber,
		Birth:          friend.Birth,
		Email:          friend.Email,
		Ex:             friend.Ex,
		AttachedInfo:   friend.AttachedInfo,
	}
	_, err := Exec(utils.StructToJsonString(tempLocalFriend))
	return err
}

func (i Friend) GetAllFriendList() (result []*model_struct.LocalFriend, err error) {
	gList, err := Exec(i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalFriend
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

func (i Friend) SearchFriendList(keyword string, isSearchUserID, isSearchNickname, isSearchRemark bool) (result []*model_struct.LocalFriend, err error) {
	gList, err := Exec(keyword, isSearchUserID, isSearchNickname, isSearchRemark)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalFriend
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

func (i Friend) GetFriendInfoByFriendUserID(FriendUserID string) (*model_struct.LocalFriend, error) {
	c, err := Exec(FriendUserID, i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := c.(string); ok {
			result := model_struct.LocalFriend{}
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

func (i Friend) GetFriendInfoList(friendUserIDList []string) (result []*model_struct.LocalFriend, err error) {
	gList, err := Exec(utils.StructToJsonString(friendUserIDList))
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalFriend
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

package indexdb

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/indexdb/temp_struct"
)

type Black struct {
	loginUserID string
}

func NewBlack(loginUserID string) *Black {
	return &Black{loginUserID: loginUserID}
}

func (i Black) GetBlackListDB() (result []*model_struct.LocalBlack, err error) {
	gList, err := Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalBlack
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

func (i Black) GetBlackListUserID() (result []string, err error) {
	gList, err := Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i Black) GetBlackInfoByBlockUserID(blockUserID string) (result *model_struct.LocalBlack, err error) {
	gList, err := Exec(blockUserID, i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp model_struct.LocalBlack
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return &temp, err
		} else {
			return nil, ErrType
		}
	}
}

func (i Black) GetBlackInfoList(blockUserIDList []string) (result []*model_struct.LocalBlack, err error) {
	gList, err := Exec(utils.StructToJsonString(blockUserIDList))
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalBlack
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

func (i Black) InsertBlack(black *model_struct.LocalBlack) error {
	_, err := Exec(utils.StructToJsonString(black))
	return err
}

func (i Black) UpdateBlack(black *model_struct.LocalBlack) error {
	tempLocalBlack := temp_struct.LocalBlack{
		Nickname:       black.Nickname,
		FaceURL:        black.FaceURL,
		Gender:         black.Gender,
		CreateTime:     black.CreateTime,
		AddSource:      black.AddSource,
		OperatorUserID: black.OperatorUserID,
		Ex:             black.Ex,
		AttachedInfo:   black.AttachedInfo,
	}
	_, err := Exec(black.OwnerUserID, black.BlockUserID, utils.StructToJsonString(tempLocalBlack))
	return err
}

func (i Black) DeleteBlack(blockUserID string) error {
	_, err := Exec(blockUserID, i.loginUserID)
	return err
}

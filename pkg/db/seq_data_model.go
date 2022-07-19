package db

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetMinSeq(ID string) (uint32, error) {
	var seqData model_struct.LocalSeq
	return seqData.MinSeq, utils.Wrap(d.conn.First(&seqData).Error, "GetMinSeq failed")
}

func (d *DataBase) SetMinSeq(ID string, minSeq uint32) error {
	seqData := model_struct.LocalSeq{ID: ID, MinSeq: minSeq}
	t := d.conn.Updates(&seqData)
	if t.RowsAffected == 0 {
		return utils.Wrap(d.conn.Create(seqData).Error, "Updates failed")
	} else {
		return utils.Wrap(t.Error, "SetMinSeq failed")
	}
}

func (d *DataBase) GetUserMinSeq() (uint32, error) {
	return 0, nil
}

func (d *DataBase) GetGroupMinSeq(groupID string) (uint32, error) {
	return 0, nil
}

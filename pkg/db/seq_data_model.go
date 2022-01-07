package db

import (
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetMinSeq() (int32, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var seqData LocalSeqData
	return seqData.Seq, utils.Wrap(d.conn.First(&seqData).Error, "GetMinSeq failed")
}

func (d *DataBase) SetMinSeq(seq int32) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	seqData := LocalSeqData{UserID: d.loginUserID, Seq: seq}
	t := d.conn.Updates(&seqData)
	if t.RowsAffected == 0 {
		return utils.Wrap(d.conn.Create(seqData).Error, "_setMinSeq failed")
	} else {
		return utils.Wrap(t.Error, "UpdateLoginUser failed")
	}
}

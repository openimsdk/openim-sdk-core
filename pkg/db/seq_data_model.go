package db

import (
	"open_im_sdk/pkg/utils"
)

func (u *DataBase) _getMinSeq() (int32, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var seqData LocalSeqData
	return seqData.Seq, utils.Wrap(u.conn.First(&seqData).Error, "_getMinSeq failed")
}

func (u *DataBase) _setMinSeq(seq int32) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()

	seqData := LocalSeqData{UserID: u.loginUserID, Seq: seq}
	t := u.conn.Updates(&seqData)
	if t.RowsAffected == 0 {
		return utils.Wrap(u.conn.Create(seqData).Error, "_setMinSeq failed")
	} else {
		return utils.Wrap(t.Error, "_updateLoginUser failed")
	}
}

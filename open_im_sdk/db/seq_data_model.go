package db

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/utils"
)

func (u *open_im_sdk.UserRelated) _getMinSeq() (int32, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var seqData open_im_sdk.LocalSeqData
	return seqData.Seq, utils.wrap(u.imdb.First(&seqData).Error, "_getMinSeq failed")
}

func (u *open_im_sdk.UserRelated) _setMinSeq(seq int32) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()

	seqData := open_im_sdk.LocalSeqData{UserID: u.loginUserID, Seq: seq}
	t := u.imdb.Updates(&seqData)
	if t.RowsAffected == 0 {
		return utils.wrap(u.imdb.Create(seqData).Error, "_setMinSeq failed")
	} else {
		return utils.wrap(t.Error, "_updateLoginUser failed")
	}
}

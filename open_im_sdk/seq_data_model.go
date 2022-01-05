package open_im_sdk

func (u *UserRelated) _getMinSeq() (int32, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var seqData LocalSeqData
	return seqData.Seq, wrap(u.imdb.First(&seqData).Error, "_getMinSeq failed")
}

func (u *UserRelated) _setMinSeq(seq int32) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()

	seqData := LocalSeqData{UserID: u.loginUserID, Seq: seq}
	t := u.imdb.Updates(&seqData)
	if t.RowsAffected == 0 {
		return wrap(u.imdb.Create(seqData).Error, "_setMinSeq failed")
	} else {
		return wrap(t.Error, "_updateLoginUser failed")
	}
}

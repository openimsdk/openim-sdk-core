package open_im_sdk

import "time"

func InsertIntoTheFriendToFriendInfo(friend *Friend) error {
	if friend.CreateTime.Unix() < 0 {
		friend.CreateTime = time.Now()
	}
	return Wrap(imdb.Create(friend).Error, "insert failed")
}

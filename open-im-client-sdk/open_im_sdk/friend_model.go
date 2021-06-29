package open_im_sdk

import (
	"fmt"
)

func insertIntoTheFriendToFriendInfo(uid, name, comment, icon string, gender int32, mobile, birth, email, ex string) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("insert into friend_info(uid,name,comment,icon,gender,mobile,birth,email,ex) values (?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		log(err.Error())
		return err
	}
	_, err = stmt.Exec(uid, name, comment, icon, gender, mobile, birth, email, ex)
	if err != nil {
		log(err.Error())
		return err
	}
	return nil
}

func delTheFriendFromFriendInfo(uid string) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from friend_info where uid=?")
	defer stmt.Close()
	if err != nil {
		log(err.Error())
		return err
	}
	_, err = stmt.Exec(uid)
	if err != nil {
		log(err.Error())
		return err
	}
	return nil
}
func updateTheFriendInfo(uid, name, comment, icon string, gender int32, mobile, birth, email, ex string) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("replace into friend_info(uid,name,comment,icon,gender,mobile,birth,email,ex) values (?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		log(err.Error())
		return err
	}
	_, err = stmt.Exec(uid, name, comment, icon, gender, mobile, birth, email, ex)
	if err != nil {
		log(err.Error())
		return err
	}
	return nil
}

func updateFriendInfo(uid, name, icon string, gender int32, mobile, birth, email, ex string) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update friend_info set `name` = ?, icon = ?, gender = ?, mobile = ?, birth = ?, email = ?, ex = ? where uid = ?")
	defer stmt.Close()
	if err != nil {
		log(err.Error())
		return err
	}
	_, err = stmt.Exec(name, icon, gender, mobile, birth, email, ex, uid)
	if err != nil {
		log(err.Error())
		return err
	}
	return nil
}

func insertIntoTheUserToBlackList(info userInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("insert into black_list(uid,name,icon,gender,mobile,birth,email,ex) values (?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		log(err.Error())
		return err
	}
	_, err = stmt.Exec(info.Uid, info.Name, info.Icon, info.Gender, info.Mobile, info.Birth, info.Email, info.Ex)
	if err != nil {
		log(err.Error())
		return err
	}
	return nil
}

func updateBlackList(info userInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("replace into black_list(uid,name,icon,gender,mobile,birth,email,ex) values (?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		log(err.Error())
		return err
	}
	_, err = stmt.Exec(info.Uid, info.Name, info.Icon, info.Gender, info.Mobile, info.Birth, info.Email, info.Ex)
	if err != nil {
		fmt.Println(err)
		log(err.Error())
		return err
	}
	return nil
}

func delTheUserFromBlackList(uid string) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from black_list where uid=?")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		log(err.Error())
		return err
	}
	_, err = stmt.Exec(uid)
	if err != nil {
		fmt.Println(err)
		log(err.Error())
		return err
	}
	return nil
}

func insertIntoTheUserToApplicationList(appUserInfo applyUserInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("insert into friend_request(uid,name,icon,gender,mobile,birth,email,ex,flag,req_message,create_time) values (?,?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		sdkLog("Prepare failed ", err.Error())
		return err
	}
	_, err = stmt.Exec(appUserInfo.Uid, appUserInfo.Name, appUserInfo.Icon, appUserInfo.Gender, appUserInfo.Mobile, appUserInfo.Birth, appUserInfo.Email, appUserInfo.Ex, appUserInfo.Flag, appUserInfo.ReqMessage, appUserInfo.ApplyTime)
	if err != nil {
		sdkLog("Exec failed, ", err.Error())
		return err
	}
	return nil
}

func delTheUserFromApplicationList(uid string) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from friend_request where uid=?")
	defer stmt.Close()
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(uid)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}

func updateApplicationList(info applyUserInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("replace into friend_request(uid,name,icon,gender,mobile,birth,email,ex,flag,req_message,create_time) values (?,?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec(info.Uid, info.Name, info.Icon, info.Gender, info.Mobile, info.Birth, info.Email, info.Ex, info.Flag, info.ReqMessage, info.ApplyTime)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func getFriendInfoByFriendUid(friendUid string) (*friendInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from friend_info  where uid=? ", friendUid)
	if err != nil {
		sdkLog("query failed, ", err.Error())
		return nil, err
	}

	var (
		uid           string
		name          string
		icon          string
		gender        int32
		mobile        string
		birth         string
		email         string
		ex            string
		comment       string
		isInBlackList int32
	)
	for stmt.Next() {
		err = stmt.Scan(&uid, &name, &comment, &icon, &gender, &mobile, &birth, &email, &ex)
		if err != nil {
			sdkLog("scan failed, ", err.Error())
			continue
		}
	}
	blackUser, _ := getBlackUsInfoByUid(uid)
	if blackUser.Uid != "" {
		isInBlackList = 1
	}
	return &friendInfo{uid, name, icon, gender, mobile, birth, email, ex, comment, isInBlackList}, nil
}

func getLocalFriendList() ([]friendInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from friend_info")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	friends := make([]friendInfo, 0)
	for stmt.Next() {
		var (
			uid           string
			name          string
			icon          string
			gender        int32
			mobile        string
			birth         string
			email         string
			ex            string
			comment       string
			isInBlackList int32
		)
		err = stmt.Scan(&uid, &name, &comment, &icon, &gender, &mobile, &birth, &email, &ex)
		if err != nil {
			sdkLog("scan failed, ", err.Error())
			continue
		}
		//check friend is in blacklist
		blackUser, _ := getBlackUsInfoByUid(uid)
		if blackUser.Uid != "" {
			isInBlackList = 1
		}
		friends = append(friends, friendInfo{uid, name, icon, gender, mobile, birth, email, ex, comment, isInBlackList})
	}
	return friends, nil
}

func getLocalFriendApplication() ([]applyUserInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from friend_request order by create_time desc")
	defer stmt.Close()
	if err != nil {
		println(err.Error())
		return nil, err
	}
	applyUsersInfo := make([]applyUserInfo, 0)
	for stmt.Next() {
		var (
			uid        string
			name       string
			icon       string
			gender     int32
			mobile     string
			birth      string
			email      string
			ex         string
			reqMessage string
			applyTime  string
			flag       int32
		)
		err = stmt.Scan(&uid, &name, &icon, &gender, &mobile, &birth, &email, &ex, &flag, &reqMessage, &applyTime)
		if err != nil {
			sdkLog("sqlite scan failed", err.Error())
			continue
		}
		applyUsersInfo = append(applyUsersInfo, applyUserInfo{uid, name, icon, gender, mobile, birth, email, ex, reqMessage, applyTime, flag})
	}
	return applyUsersInfo, nil
}

func getLocalBlackList() ([]userInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from black_list")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	usersInfo := make([]userInfo, 0)
	for stmt.Next() {
		var (
			uid    string
			name   string
			icon   string
			gender int32
			mobile string
			birth  string
			email  string
			ex     string
		)
		err = stmt.Scan(&uid, &name, &icon, &gender, &mobile, &birth, &email, &ex)
		if err != nil {
			sdkLog("sqlite scan failed", err.Error())
			continue
		}
		usersInfo = append(usersInfo, userInfo{uid, name, icon, gender, mobile, birth, email, ex})
	}
	return usersInfo, nil
}

func getBlackUsInfoByUid(blackUid string) (*userInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from black_list where uid=?", blackUid)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	var (
		uid    string
		name   string
		icon   string
		gender int32
		mobile string
		birth  string
		email  string
		ex     string
	)
	for stmt.Next() {
		err = stmt.Scan(&uid, &name, &icon, &gender, &mobile, &birth, &email, &ex)
		if err != nil {
			sdkLog("scan failed, ", err.Error())
			continue
		}
	}

	return &userInfo{uid, name, icon, gender, mobile, birth, email, ex}, nil
}

/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/6/8 9:39).
 */
package open_im_sdk

import (
	"encoding/json"
	"fmt"
)

type ConversationStruct struct {
	ConversationID    string `json:"conversationID"`
	ConversationType  int    `json:"conversationType"`
	UserID            string `json:"userID"`
	GroupID           string `json:"groupID"`
	ShowName          string `json:"showName"`
	FaceURL           string `json:"faceUrl"`
	RecvMsgOpt        int    `json:"recvMsgOpt"`
	UnreadCount       int    `json:"unreadCount"`
	LatestMsg         string `json:"latestMsg"`
	LatestMsgSendTime int64  `json:"latestMsgSendTime"`
	DraftText         string `json:"draftText"`
	DraftTimestamp    int64  `json:"draftTimestamp"`
	IsPinned          int    `json:"isPinned"`
}

func getAllConversationListModel() (err error, list []*ConversationStruct) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	rows, err := initDB.Query("SELECT * FROM conversation order by  case when is_pinned=1 then 0 else 1 end,latest_msg_send_time DESC")
	for rows.Next() {
		c := new(ConversationStruct)
		err = rows.Scan(&c.ConversationID, &c.ConversationType, &c.UserID, &c.GroupID, &c.ShowName,
			&c.FaceURL, &c.RecvMsgOpt, &c.UnreadCount, &c.LatestMsg, &c.LatestMsgSendTime, &c.DraftText, &c.DraftTimestamp, &c.IsPinned)
		if err != nil {
			sdkLog("getAllConversationListModel ,err:", err.Error())
			continue
		} else {
			list = append(list, c)
		}
	}
	return nil, list
}

func insertConversationModel(c *ConversationStruct) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("INSERT INTO conversation(conversation_id, conversation_type, " +
		"user_id,group_id,show_name,face_url,recv_msg_opt,unread_count,latest_msg,latest_msg_send_time,draft_text,draft_timestamp,is_pinned) values(?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log(err.Error())
		return err
	}
	_, err = stmt.Exec(c.ConversationID, c.ConversationType, c.UserID, c.GroupID, c.ShowName, c.FaceURL, c.RecvMsgOpt, c.UnreadCount, c.LatestMsg, c.LatestMsgSendTime, c.DraftText, c.DraftTimestamp, c.IsPinned)

	if err != nil {
		log(err.Error())
		return err
	}
	return nil
}
func getConversationLatestMsgModel(conversationID string) (err error, latestMsg *MsgStruct) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	var s string
	msg := new(MsgStruct)
	rows, err := initDB.Query("SELECT latest_msg FROM conversation where  conversation_id=?", conversationID)
	if err != nil {
		log(err.Error())
		return err, nil
	}
	for rows.Next() {
		err = rows.Scan(&s)
		if err != nil {
			sdkLog("getConversationLatestMsgModel ,err:", err.Error())
			continue
		}
	}

	err = json.Unmarshal([]byte(s), msg)
	if err != nil {
		log(err.Error())
		return err, nil
	}
	return nil, msg
}
func setConversationLatestMsgModel(c *ConversationStruct, conversationID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update conversation set latest_msg=?,latest_msg_send_time=? where conversation_id=?")
	if err != nil {
		log(err.Error())
		return err
	}

	_, err = stmt.Exec(c.LatestMsg, c.LatestMsgSendTime, conversationID)
	if err != nil {
		log(err.Error())
		return err
	}
	return nil

}
func judgeConversationIfExists(conversationID string) bool {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	var count int
	rows, err := initDB.Query("select count(*) from conversation where  conversation_id=?", conversationID)
	if err != nil {
		fmt.Println("judge err")
		log(err.Error())
		return false
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			sdkLog("getConversationLatestMsgModel ,err:", err.Error())
			continue
		} else {
			if count == 1 {
				return true
			} else {
				return false
			}
		}
	}
	return false
}
func addConversationOrUpdateLatestMsg(c *ConversationStruct, conversationID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("INSERT INTO conversation(conversation_id, conversation_type, user_id,group_id,show_name,face_url,recv_msg_opt,unread_count,latest_msg,latest_msg_send_time,draft_text,draft_timestamp,is_pinned)" +
		" values(?,?,?,?,?,?,?,?,?,?,?,?,?)" +
		"ON CONFLICT(conversation_id) DO UPDATE SET latest_msg = ?")
	if err != nil {
		log(err.Error())
		return err
	}
	_, err = stmt.Exec(c.ConversationID, c.ConversationType, c.UserID, c.GroupID, c.ShowName, c.FaceURL, c.RecvMsgOpt, c.UnreadCount, c.LatestMsg, c.LatestMsgSendTime, c.DraftText, c.DraftTimestamp, c.IsPinned, c.LatestMsg)

	if err != nil {
		log(err.Error())
		return err
	}
	return nil

}
func getOneConversationModel(conversationID string) (err error, c ConversationStruct) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	rows, err := initDB.Query("SELECT * FROM conversation where  conversation_id=?", conversationID)
	if err != nil {
		fmt.Println("com herre")
		log(err.Error())
		return err, c
	}
	for rows.Next() {
		err = rows.Scan(&c.ConversationID, &c.ConversationType, &c.UserID, &c.GroupID, &c.ShowName,
			&c.FaceURL, &c.RecvMsgOpt, &c.UnreadCount, &c.LatestMsg, &c.LatestMsgSendTime, &c.DraftText, &c.DraftTimestamp, &c.IsPinned)
		if err != nil {
			sdkLog("getOneConversationModel ,err:", err.Error())
			continue
		}
	}
	return nil, c

}
func deleteConversationModel(conversationID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from conversation where conversation_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(conversationID)
	if err != nil {
		log(err.Error())
		return err
	}
	return nil
}
func setConversationDraftModel(conversationID, draftText string, DraftTimestamp int64) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update conversation set draft_text=?,draft_timestamp=?where conversation_id=?")
	if err != nil {
		log(err.Error())
		return err
	}

	_, err = stmt.Exec(draftText, DraftTimestamp, conversationID)
	if err != nil {
		log(err.Error())
		return err
	}
	return nil

}
func pinConversationModel(conversationID string, isPinned int) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update conversation set is_pinned=? where conversation_id=?")
	if err != nil {
		log(err.Error())
		return err
	}

	_, err = stmt.Exec(isPinned, conversationID)
	if err != nil {
		log(err.Error())
		return err
	}
	return nil

}
func setConversationUnreadCount(unreadCount int, conversationID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update conversation set unread_count=? where conversation_id=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(unreadCount, conversationID)
	if err != nil {
		return err
	}
	return nil

}
func incrConversationUnreadCount(conversationID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update conversation set unread_count = unread_count+1 where conversation_id=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(conversationID)
	if err != nil {
		return err
	}
	return nil

}
func getTotalUnreadMsgCountModel() (totalUnreadCount int32, err error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	rows, err := initDB.Query("SELECT IFNULL(SUM(unread_count), 0) FROM conversation")
	if err != nil {
		log(err.Error())
		return totalUnreadCount, err
	}
	for rows.Next() {
		err = rows.Scan(&totalUnreadCount)
		if err != nil {
			sdkLog("getTotalUnreadMsgCountModel ,err:", err.Error())
			continue
		}
	}
	return totalUnreadCount, err

}
func getMultipleConversationModel(conversationIDList []string) (err error, list []*ConversationStruct) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	rows, err := initDB.Query("SELECT * FROM conversation where conversation_id in (" + sqlStringHandle(conversationIDList) + ")")
	fmt.Println("SELECT * FROM conversation where conversation_id in (" + sqlStringHandle(conversationIDList) + ")")
	for rows.Next() {
		temp := new(ConversationStruct)
		err = rows.Scan(&temp.ConversationID, &temp.ConversationType, &temp.UserID, &temp.GroupID, &temp.ShowName,
			&temp.FaceURL, &temp.RecvMsgOpt, &temp.UnreadCount, &temp.LatestMsg, &temp.LatestMsgSendTime, &temp.DraftText, &temp.DraftTimestamp, &temp.IsPinned)
		if err != nil {
			sdkLog("getMultipleConversationModel err:", err.Error())
			return err, nil
		} else {
			list = append(list, temp)
		}
	}
	return nil, list
}
func sqlStringHandle(ss []string) (s string) {
	for i := 0; i < len(ss); i++ {
		s += "'" + ss[i] + "'"
		if i < len(ss)-1 {
			s += ","
		}
	}
	return s
}

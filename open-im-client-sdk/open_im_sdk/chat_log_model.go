/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/6/7 15:22).
 */
package open_im_sdk

import (
	"database/sql"
	"fmt"
)

type ChatLog struct {
	MsgId            string
	SendID           string
	IsRead           int32
	Seq              int64
	Status           int32
	SessionType      int32
	RecvID           string
	ContentType      int32
	MsgFrom          int32
	Content          string
	Remark           sql.NullString
	SenderPlatformID int32
	SendTime         int64
	CreateTime       int64
}

func insertSendMessageToChatLog(message *MsgStruct) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("INSERT INTO chat_log(msg_id, send_id, is_read," +
		" seq,status, session_type, recv_id, content_type, msg_from, content, remark,sender_platform_id, send_time,create_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log(err.Error())
		return err
	}
	_, err = stmt.Exec(message.ClientMsgID, message.SendID,
		getIsRead(message.IsRead), message.Seq, message.Status, message.SessionType, message.RecvID, message.ContentType,
		message.MsgFrom, message.Content, message.Remark, message.PlatformID, message.SendTime, message.CreateTime)
	if err != nil {
		log(err.Error())
		return err
	}
	return nil
}
func insertPushMessageToChatLog(message *MsgStruct) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("INSERT INTO chat_log(msg_id, send_id, is_read," +
		" seq,status, session_type, recv_id, content_type, msg_from, content, remark,sender_platform_id, send_time,create_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)" +
		"ON CONFLICT(msg_id) DO UPDATE SET seq = ?")
	if err != nil {
		sdkLog("Prepare failed, ", err.Error())
		return err
	}
	_, err = stmt.Exec(message.ServerMsgID, message.SendID,
		getIsRead(message.IsRead), message.Seq, message.Status, message.SessionType, message.RecvID, message.ContentType,
		message.MsgFrom, message.Content, message.Remark, message.PlatformID, message.SendTime, message.CreateTime, message.Seq)
	if err != nil {
		sdkLog("Exec failed, ", err.Error())
		return err
	}
	return nil
}
func updateMessageSeq(message *MsgStruct) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update chat_log set seq=? where msg_id=?")
	if err != nil {
		sdkLog("Prepare failed, ", err.Error())
		return err
	}
	_, err = stmt.Exec(message.Seq, message.ServerMsgID)
	if err != nil {
		sdkLog("Exec failed ", err.Error())
		return err
	}
	return nil
}
func judgeMessageIfExists(message *MsgStruct) bool {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	var count int
	rows, err := initDB.Query("select count(*) from chat_log where  msg_id=?", message.ServerMsgID)
	if err != nil {
		sdkLog("Query failed, ", err.Error())
		return false
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log(err.Error())

			return false
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
func getChatLog(msgID string) (chat *ChatLog, err error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	chat = &ChatLog{}
	// query
	rows, err := initDB.Query("SELECT * FROM chat_log where msg_id = " + msgID)
	if err != nil {
		log(err.Error())
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&chat.MsgId, &chat.SendID, &chat.IsRead,
			&chat.Seq, &chat.Status, &chat.SessionType, &chat.RecvID, &chat.ContentType,
			&chat.MsgFrom, &chat.Content, &chat.Remark, &chat.SenderPlatformID, &chat.SendTime, &chat.CreateTime)
		if err != nil {
			sdkLog("getChatLog,err:", err.Error())
			continue
		}
	}

	rows.Close()
	return chat, nil
}

//func GetHistoryMessage(recvID string,count int)  {
//	rows, err := defInitDB.Query("SELECT * FROM chat_log where seq = 0 And recv_id =" +recvID+"And send_id = "+LoginUid )
//	if err != nil {
//		return nil, err
//	}
//	rows.c
//	for rows.Next() {
//		t = c
//		err = rows.Scan(&t.ConversationID, &t.ConversationType, &t.UserID, &t.GroupID, &t.ShowName,
//			&t.FaceURL, &t.RecvMsgOpt, &t.UnreadCount, &t.LatestMsg, &t.DraftText, &t.DraftTimestamp, &t.IsPinned)
//		if err != nil {
//			return err, nil
//		} else {
//			list = append(list, &t)
//		}
//	}
//
//}
func setSingleMessageHasRead(sendID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update chat_log set is_read=? where send_id=?And is_read=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(HasRead, sendID, NotRead)
	if err != nil {
		return err
	}
	return nil

}
func setMessageStatus(msgID string, status int) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update chat_log set status=? where msg_id=?")
	if err != nil {
		sdkLog("prepare failed, err: ", err.Error())
		return err
	}
	_, err = stmt.Exec(status, msgID)
	if err != nil {
		sdkLog("exec failed, err: ", err.Error())
		return err
	}
	return nil
}
func setMessageHasReadByMsgID(msgID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update chat_log set is_read=? where msg_id=?And is_read=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(HasRead, msgID, NotRead)
	if err != nil {
		return err
	}
	return nil

}
func getHistoryMessage(sourceConversationID string, startTime int64, count int) (err error, list MsgFormats) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	rows, err := initDB.Query("select * from chat_log WHERE (send_id = ? OR recv_id =? )AND content_type<=?AND status!=?AND send_time<?  order by send_time DESC  LIMIT ? OFFSET 0 ",
		sourceConversationID, sourceConversationID, AcceptFriendApplicationTip, MsgStatusHasDeleted, startTime, count)
	for rows.Next() {
		temp := new(MsgStruct)
		err = rows.Scan(&temp.ServerMsgID, &temp.SendID, &temp.IsRead, &temp.Seq, &temp.Status, &temp.SessionType,
			&temp.RecvID, &temp.ContentType, &temp.MsgFrom, &temp.Content, &temp.Remark, &temp.PlatformID, &temp.SendTime, &temp.CreateTime)
		if err != nil {
			sdkLog("getHistoryMessage,err:", err.Error())
			continue
		} else {
			temp.ClientMsgID = temp.ServerMsgID
			err = msgHandleByContentType(temp)
			if err != nil {
				sdkLog("getHistoryMessage,err:", err.Error())
				continue
			}
			list = append(list, temp)
		}
	}
	return nil, list
}
func deleteMessageByConversationModel(sourceConversationID string, maxSeq int64) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from chat_log  where send_id=? or recv_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(sourceConversationID, sourceConversationID)
	if err != nil {
		log(err.Error())
		return err
	}
	return nil
}
func deleteMessageByMsgID(msgID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from chat_log  where msg_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(msgID)
	if err != nil {
		log(err.Error())
		return err
	}
	return nil
}
func updateMessageTimeAndMsgIDStatus(ClientMsgID, ServerMsgID string, sendTime int64, status int) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update chat_log set msg_id=?,send_time=?, status=? where msg_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(ServerMsgID, sendTime, status, ClientMsgID)
	if err != nil {
		return err
	}
	return nil
}
func getMultipleMessageModel(messageIDList []string) (err error, list []*MsgStruct) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	rows, err := initDB.Query("SELECT * FROM chat_log where msg_id in (" + sqlStringHandle(messageIDList) + ")")
	fmt.Println("SELECT * FROM conversation where conversation_id in (" + sqlStringHandle(messageIDList) + ")")
	defer rows.Close()
	for rows.Next() {
		temp := new(MsgStruct)
		err = rows.Scan(&temp.ServerMsgID, &temp.SendID, &temp.IsRead, &temp.Seq, &temp.Status, &temp.SessionType,
			&temp.RecvID, &temp.ContentType, &temp.MsgFrom, &temp.Content, &temp.Remark, &temp.PlatformID, &temp.SendTime, &temp.CreateTime)
		if err != nil {
			sdkLog("getMultipleMessageModel,err:", err.Error())
			continue
		} else {
			err = msgHandleByContentType(temp)
			if err != nil {
				sdkLog("getMultipleMessageModel,err:", err.Error())
				continue
			}
			list = append(list, temp)
		}
	}
	return nil, list
}
func getLocalMaxSeqModel() (seq int64, err error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	rows, err := initDB.Query("SELECT IFNULL(MAX(seq), 0) FROM chat_log")
	if err != nil {
		log(err.Error())
		return seq, err
	}
	for rows.Next() {
		err = rows.Scan(&seq)
		if err != nil {
			sdkLog("getTotalUnreadMsgCountModel ,err:", err.Error())
			continue
		}
	}
	return seq, err
}

func msgHandleByContentType(msg *MsgStruct) (err error) {
	switch msg.ContentType {
	case Text:
	case Picture:
		err = jsonStringToStruct(msg.Content, &msg.PictureElem)
	case Sound:
		err = jsonStringToStruct(msg.Content, &msg.SoundElem)
	case Video:
		err = jsonStringToStruct(msg.Content, &msg.VideoElem)
	case File:
		err = jsonStringToStruct(msg.Content, &msg.FileElem)
	}
	return err
}

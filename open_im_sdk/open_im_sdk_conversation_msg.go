package open_im_sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

func GetAllConversationList(callback Base) {
	err, list := getAllConversationListModel()
	if err != nil {
		callback.OnError(203, err.Error())
	} else {
		if list != nil {
			callback.OnSuccess(structToJsonString(list))
		} else {
			callback.OnSuccess(structToJsonString([]ConversationStruct{}))
		}
	}

}
func GetOneConversation(conversationID string, callback Base) {
	err, c := getOneConversationModel(conversationID)
	if err != nil {
		callback.OnError(203, err.Error())
	} else {
		callback.OnSuccess(structToJsonString(c))
	}

}
func GetMultipleConversation(conversationIDList string, callback Base) {
	var c []string
	err := json.Unmarshal([]byte(conversationIDList), &c)
	if err != nil {
		callback.OnError(200, err.Error())
		fmt.Println("json PARSING ERROR", 111)
		log("json ")
	}
	err, list := getMultipleConversationModel(c)
	if err != nil {
		callback.OnError(203, err.Error())
	} else {
		if list != nil {
			callback.OnSuccess(structToJsonString(list))
		} else {
			callback.OnSuccess(structToJsonString([]ConversationStruct{}))
		}
	}
}
func DeleteConversation(conversationID string, callback Base) {
	//Transaction operation required
	var sourceID string
	err, c := getOneConversationModel(conversationID)
	if err != nil {
		callback.OnError(201, err.Error())
		return
	}
	switch c.ConversationType {
	case SingleChatType:
		sourceID = c.UserID
	case GroupChatType:
		sourceID = c.GroupID
	}
	maxSeq, err := getLocalMaxSeqModel()
	if err != nil {
		callback.OnError(201, err.Error())
		return
	}
	err = deleteMessageByConversationModel(sourceID, maxSeq)
	if err != nil {
		callback.OnError(202, err.Error())
		return
	}
	err = deleteConversationModel(conversationID)
	if err != nil {
		callback.OnError(203, err.Error())
		return
	} else {
		callback.OnSuccess("")
		_ = triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: ConAndUnreadChange})
	}
}
func SetConversationDraft(conversationID, draftText string, callback Base) {
	err := setConversationDraftModel(conversationID, draftText, getCurrentTimestampBySecond())
	if err != nil {
		callback.OnError(203, err.Error())
	} else {
		callback.OnSuccess("")
		_ = triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: ConAndUnreadChange})
	}
}
func PinConversation(conversationID string, isPinned bool, callback Base) {
	var i int
	if isPinned {
		i = 1
	} else {
		i = 0
	}
	err := pinConversationModel(conversationID, i)
	if err != nil {
		callback.OnError(203, err.Error())
	} else {
		callback.OnSuccess("")
		_ = triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: ConAndUnreadChange})
	}

}
func GetTotalUnreadMsgCount(callback Base) {
	count, err := getTotalUnreadMsgCountModel()
	if err != nil {
		callback.OnError(203, err.Error())
	} else {
		callback.OnSuccess(int32ToString(count))
	}

}

type OnConversationListener interface {
	OnSyncServerStart()
	OnSyncServerFinish()
	OnSyncServerFailed()
	OnNewConversation(conversationList string)
	OnConversationChanged(conversationList string)
	OnTotalUnreadMessageCountChanged(totalUnreadCount int32)
}

func SetConversationListener(listener OnConversationListener) {
	ConListener.ConversationListener = listener
}

type OnAdvancedMsgListener interface {
	OnRecvNewMessage(message string)
	OnRecvC2CReadReceipt(msgReceiptList string)
	OnRecvMessageRevoked(msgId string)
}

func AddAdvancedMsgListener(listener OnAdvancedMsgListener) {
	if listener == nil {
		fmt.Println("AddAdvancedMsgListener listener is null")
		return
	}
	ConListener.MsgListenerList = append(ConListener.MsgListenerList, listener)
}

func RemoveAdvancedMsgListener(listener OnAdvancedMsgListener) {
}

func ForceSyncMsg() {

	if SdkInitManager.conn != nil {
		SdkInitManager.conn.Close()
		SdkInitManager.conn = nil
	}

	c2v := cmd2Value{Cmd: CmdForceSyncMsg}
	sendCmd(ConversationCh, c2v, 1)
}

type SendMsgCallBack interface {
	Base
	OnProgress(progress int)
}

func CreateTextMessage(text string) string {
	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, Text)
	s.Content = text
	return structToJsonString(s)
}
func CreateTextAtMessage(text, atUserList string) string {
	var users []string
	_ = json.Unmarshal([]byte(atUserList), &users)
	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, Text)
	s.Content = text
	s.ForceList = users
	return structToJsonString(s)
}

func CreateVideoMessageFromFullPath(videoFullPath string, videoType string, duration int64, snapshotFullPath string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := fileTmpPath(videoFullPath) //a->b
		s, err := copyFile(videoFullPath, dstFile)
		if err != nil {
			sdkLog("open file failed: ", err, videoFullPath)
		}
		sdkLog("videoFullPath dstFile", videoFullPath, dstFile, s)
		dstFile = fileTmpPath(snapshotFullPath) //a->b
		s, err = copyFile(snapshotFullPath, dstFile)
		if err != nil {
			sdkLog("open file failed: ", err, snapshotFullPath)
		}
		sdkLog("snapshotFullPath dstFile", snapshotFullPath, dstFile, s)
		wg.Done()
	}()

	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, Video)
	s.VideoElem.VideoPath = videoFullPath
	s.VideoElem.VideoType = videoType
	s.VideoElem.Duration = duration
	if snapshotFullPath == "" {
		s.VideoElem.SnapshotPath = ""
	} else {
		s.VideoElem.SnapshotPath = snapshotFullPath
	}
	fi, err := os.Stat(s.VideoElem.VideoPath)
	if err != nil {
		sdkLog(err.Error())
		return ""
	}
	s.VideoElem.VideoSize = fi.Size()
	if snapshotFullPath != "" {
		imageInfo, err := getImageInfo(s.VideoElem.SnapshotPath)
		if err != nil {
			sdkLog("CreateVideoMessage err:", err.Error())
			return ""
		}
		s.VideoElem.SnapshotHeight = imageInfo.Height
		s.VideoElem.SnapshotWidth = imageInfo.Width
		s.VideoElem.SnapshotSize = imageInfo.Size
	}
	wg.Wait()
	return structToJsonString(s)
}

func CreateImageMessageFromFullPath(imageFullPath string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := fileTmpPath(imageFullPath) //a->b
		_, err := copyFile(imageFullPath, dstFile)
		if err != nil {
			sdkLog("open file failed: ", err, imageFullPath)
		}
		wg.Done()
	}()

	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, Picture)
	s.PictureElem.SourcePath = imageFullPath
	sdkLog("ImageMessage  path:", s.PictureElem.SourcePath)
	imageInfo, err := getImageInfo(s.PictureElem.SourcePath)
	if err != nil {
		sdkLog("CreateImageMessage err:", err.Error())
		return ""
	}
	s.PictureElem.SourcePicture.Width = imageInfo.Width
	s.PictureElem.SourcePicture.Height = imageInfo.Height
	s.PictureElem.SourcePicture.Type = imageInfo.Type
	s.PictureElem.SourcePicture.Size = imageInfo.Size
	wg.Wait()
	return structToJsonString(s)
}

func CreateImageMessage(imagePath string) string {
	sdkLog("start1: ", time.Now())
	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, Picture)
	s.PictureElem.SourcePath = SvrConf.DbDir + imagePath
	sdkLog("ImageMessage  path:", s.PictureElem.SourcePath)
	sdkLog("end1", time.Now())

	sdkLog("start2 ", time.Now())
	imageInfo, err := getImageInfo(s.PictureElem.SourcePath)
	if err != nil {
		sdkLog("CreateImageMessage err:", err.Error())
		return ""
	}
	sdkLog("end2", time.Now())

	s.PictureElem.SourcePicture.Width = imageInfo.Width
	s.PictureElem.SourcePicture.Height = imageInfo.Height
	s.PictureElem.SourcePicture.Type = imageInfo.Type
	s.PictureElem.SourcePicture.Size = imageInfo.Size
	return structToJsonString(s)
}
func CreateSoundMessage(soundPath string, duration int64) string {
	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, Sound)
	s.SoundElem.SoundPath = SvrConf.DbDir + soundPath
	s.SoundElem.Duration = duration
	fi, err := os.Stat(s.SoundElem.SoundPath)
	if err != nil {
		sdkLog(err.Error())
		return ""
	}
	s.SoundElem.DataSize = fi.Size()
	return structToJsonString(s)
}

func CreateVideoMessage(videoPath string, videoType string, duration int64, snapshotPath string) string {
	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, Video)
	s.VideoElem.VideoPath = SvrConf.DbDir + videoPath
	s.VideoElem.VideoType = videoType
	s.VideoElem.Duration = duration
	if snapshotPath == "" {
		s.VideoElem.SnapshotPath = ""
	} else {
		s.VideoElem.SnapshotPath = SvrConf.DbDir + snapshotPath
	}
	fi, err := os.Stat(s.VideoElem.VideoPath)
	if err != nil {
		sdkLog(err.Error())
		return ""
	}
	s.VideoElem.VideoSize = fi.Size()
	if snapshotPath != "" {
		imageInfo, err := getImageInfo(s.VideoElem.SnapshotPath)
		if err != nil {
			sdkLog("CreateVideoMessage err:", err.Error())
			return ""
		}
		s.VideoElem.SnapshotHeight = imageInfo.Height
		s.VideoElem.SnapshotWidth = imageInfo.Width
		s.VideoElem.SnapshotSize = imageInfo.Size
	}
	return structToJsonString(s)
}
func CreateFileMessage(filePath string, fileName string) string {
	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, File)
	s.FileElem.FilePath = SvrConf.DbDir + filePath
	s.FileElem.FileName = fileName
	fi, err := os.Stat(s.FileElem.FilePath)
	if err != nil {
		sdkLog(err.Error())
		return ""
	}
	s.FileElem.FileSize = fi.Size()
	return structToJsonString(s)
}
func CreateMergerMessage(messageList, title, summaryList string) string {
	var messages []*MsgStruct
	var summaries []string
	s := MsgStruct{}
	_ = json.Unmarshal([]byte(messageList), messages)
	_ = json.Unmarshal([]byte(summaryList), &summaries)
	initBasicInfo(&s, UserMsgType, Merger)
	s.MergeElem.AbstractList = summaries
	s.MergeElem.Title = title
	s.MergeElem.MultiMessage = messages
	return structToJsonString(s)
}
func CreateForwardMessage(m string) string {
	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, 0)
	return structToJsonString(s)
}

func SendMessage(callback SendMsgCallBack, message, receiver, groupID string, onlineUserOnly bool) string {
	var conversationID string
	r := SendMsgRespFromServer{}
	a := paramsUserSendMsg{}
	s := MsgStruct{}
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(message), &s)
	if err != nil {
		callback.OnError(2038, err.Error())
		sdkLog("json unmarshal err:", err.Error())
		return ""
	}

	var delFile []string
	switch s.ContentType {
	case Text:
	case Picture:
		var sourcePath string
		if fileExist(s.PictureElem.SourcePath) {
			sourcePath = s.PictureElem.SourcePath
		} else {
			sourcePath = fileTmpPath(s.PictureElem.SourcePath)
			delFile = append(delFile, sourcePath)
		}
		sourceUrl, uuid, err := uploadImage(sourcePath, callback)
		if err != nil {
			callback.OnError(301, err.Error())
			fmt.Println("oss Picture upload err", 111)
			return ""
		} else {
			s.PictureElem.SourcePicture.Url = sourceUrl
			s.PictureElem.SourcePicture.UUID = uuid
			s.PictureElem.SnapshotPicture.Url = sourceUrl + "?imageView2/2/w/" + ZoomScale + "/h/" + ZoomScale
			s.PictureElem.SnapshotPicture.Width = int32(stringToInt(ZoomScale))
			s.PictureElem.SnapshotPicture.Height = int32(stringToInt(ZoomScale))
			s.Content = structToJsonString(s.PictureElem)
		}
	case Sound:
		soundURL, uuid, err := uploadSound(s.SoundElem.SoundPath, callback)
		if err != nil {
			callback.OnError(301, err.Error())
			fmt.Println("oss Sound upload err", 111)
			sdkLog("uploadSound err:", err.Error())
			return ""
		} else {
			s.SoundElem.SourceURL = soundURL
			s.SoundElem.UUID = uuid
			s.Content = structToJsonString(s.SoundElem)
		}
	case Video:
		var videoPath string
		var snapPath string
		if fileExist(s.PictureElem.SourcePath) {
			videoPath = s.VideoElem.VideoPath
			snapPath = s.VideoElem.SnapshotPath
		} else {
			videoPath = fileTmpPath(s.VideoElem.VideoPath)
			snapPath = fileTmpPath(s.VideoElem.SnapshotPath)
			delFile = append(delFile, videoPath)
			delFile = append(delFile, snapPath)
		}

		snapshotURL, snapshotUUID, videoURL, videoUUID, err := uploadVideo(videoPath, snapPath, callback)
		if err != nil {
			callback.OnError(301, err.Error())
			sdkLog("oss  Video upload err:", err.Error())
			return ""
		} else {
			s.VideoElem.VideoURL = videoURL
			s.VideoElem.SnapshotUUID = snapshotUUID
			s.VideoElem.SnapshotURL = snapshotURL
			s.VideoElem.VideoUUID = videoUUID
			s.Content = structToJsonString(s.VideoElem)
		}
	case File:
		fileURL, fileUUID, err := uploadFile(s.FileElem.FilePath, callback)
		if err != nil {
			callback.OnError(301, err.Error())
			sdkLog("oss  File upload err:", err.Error())
			return ""
		} else {
			s.FileElem.SourceURL = fileURL
			s.FileElem.UUID = fileUUID
			s.Content = structToJsonString(s.FileElem)
		}
	default:
		callback.OnError(2038, "Not currently supported ")
	}

	if receiver == "" {
		s.SessionType = GroupChatType
		s.RecvID = groupID
		conversationID = GetConversationIDBySessionType(groupID, GroupChatType)
	} else if groupID == "" {
		s.SessionType = SingleChatType
		s.RecvID = receiver
		conversationID = GetConversationIDBySessionType(receiver, SingleChatType)
	} else {
		callback.OnError(201, "args err")
		return ""
	}
	go func() {
		//Store messages to local database
		err = insertSendMessageToChatLog(&s)
		if err != nil {
			callback.OnError(202, err.Error())
			return
			fmt.Println("INSERTION ERROR", 22221)

		}
		c := ConversationStruct{
			ConversationID:    conversationID,
			ConversationType:  int(s.SessionType),
			UserID:            s.RecvID,
			RecvMsgOpt:        1,
			LatestMsg:         structToJsonString(s),
			LatestMsgSendTime: s.CreateTime,
		}
		_ = triggerCmdUpdateConversation(updateConNode{conversationID, AddConOrUpLatMsg,
			c})
		//Protocol conversion
		a.ReqIdentifier = 1003
		a.PlatformID = s.PlatformID
		a.SendID = s.SendID
		a.OperationID = operationIDGenerator()
		a.Data.SessionType = s.SessionType
		a.Data.MsgFrom = s.MsgFrom
		a.MsgIncr = 1
		a.Data.ForceList = []string{}
		a.Data.ContentType = s.ContentType
		a.Data.RecvID = s.RecvID
		a.Data.ForceList = s.ForceList
		a.Data.Content = s.Content
		a.Data.ClientMsgID = s.ClientMsgID
		if onlineUserOnly {
			a.Data.Options["history"] = 0
			a.Data.Options["persistent"] = 0
		} else {
			a.Data.Options = m
		}
		a.Data.OffLineInfo = m
		bMsg, err := post2Api(sendMsgRouter, a, token)
		if err != nil {
			callback.OnError(http.StatusInternalServerError, err.Error())
			_ = updateMessageTimeAndMsgIDStatus(s.ClientMsgID, s.ClientMsgID, s.CreateTime, MsgStatusSendFailed)

		} else if err = json.Unmarshal(bMsg, &r); err != nil {
			callback.OnError(200, err.Error()+"  "+string(bMsg))
			_ = updateMessageTimeAndMsgIDStatus(s.ClientMsgID, s.ClientMsgID, s.CreateTime, MsgStatusSendFailed)
		} else {
			if r.ErrCode != 0 {
				callback.OnError(r.ErrCode, r.ErrMsg)
				_ = updateMessageTimeAndMsgIDStatus(s.ClientMsgID, s.ClientMsgID, s.CreateTime, MsgStatusSendFailed)
			} else {
				callback.OnSuccess("")
				callback.OnProgress(100)

				for _, v := range delFile {
					err := os.Remove(v)
					if err != nil {
						sdkLog(err.Error())
					}
				}
				_ = updateMessageTimeAndMsgIDStatus(r.Data.ClientMsgID, r.Data.ServerMsgID, r.Data.SendTime, MsgStatusSendSuccess)
				s.ServerMsgID = r.Data.ServerMsgID
				s.SendTime = r.Data.SendTime
				c.LatestMsg = structToJsonString(s)
				c.LatestMsgSendTime = s.SendTime
				_ = triggerCmdUpdateConversation(updateConNode{conversationID, AddConOrUpLatMsg,
					c})
			}
		}
	}()
	return s.ClientMsgID
}

func GetHistoryMessageList(callback Base, getMessageOptions string) {
	var sourceID string
	var conversationID string
	var startTime int64
	p := PullMsgReq{}
	err := json.Unmarshal([]byte(getMessageOptions), &p)
	if err != nil {
		callback.OnError(200, err.Error())
		return
	}
	if p.UserID == "" {
		sourceID = p.GroupID
		conversationID = GetConversationIDBySessionType(sourceID, GroupChatType)
	} else {
		sourceID = p.UserID
		conversationID = GetConversationIDBySessionType(sourceID, SingleChatType)

	}
	if p.StartMsg == nil {
		err, m := getConversationLatestMsgModel(conversationID)
		if err != nil {
			fmt.Println("get history err :", err)
			return
		}
		startTime = m.SendTime + TimeOffset

	} else {
		startTime = p.StartMsg.SendTime
	}
	fmt.Println("sourceID:", sourceID, "startTime:", startTime, "count:", p.Count)
	err, list := getHistoryMessage(sourceID, startTime, p.Count)
	sort.Sort(list)
	if err != nil {
		callback.OnError(203, err.Error())
	} else {
		if list != nil {
			callback.OnSuccess(structToJsonString(list))
		} else {
			callback.OnSuccess(structToJsonString([]MsgStruct{}))
		}
	}

}
func RevokeMessage(callback Base, message string) {
	s := MsgStruct{}
	err := json.Unmarshal([]byte(message), &s)
	if err != nil {
		callback.OnError(200, err.Error())
	}
	//Send message internally
	internalSendMsg(callback, s, s.RevokeMessage.RecvID, s.RevokeMessage.GroupID, false)
	callback.OnSuccess("")
}

func MarkC2CMessageAsRead(callback Base, receiver string, msgList string) {
	go func() {
		var msgIdList []string
		err := json.Unmarshal(([]byte)(msgList), &msgIdList)
		if err != nil {
			sdkLog("unmarshal failed, ", msgIdList, err.Error())
			callback.OnError(ErrCodeConversation, err.Error())
			return
		}
		s := MsgStruct{}
		initBasicInfo(&s, UserMsgType, C2CMessageAsRead)
		s.Content = msgList
		fmt.Println("xxxxxxxxxxxxxx", receiver, msgIdList)
		internalSendMsg(callback, s, receiver, "", false)
	}()
}

func MarkSingleMessageHasRead(callback Base, userID string) {
	conversationID := GetConversationIDBySessionType(userID, SingleChatType)
	if err := setSingleMessageHasRead(userID); err != nil {
		callback.OnError(201, err.Error())
	} else {
		callback.OnSuccess("")
		triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: UnreadCountSetZero})
	}

}

func DeleteMessageFromLocalStorage(callback Base, message string) {
	var conversation ConversationStruct
	var conversationID string
	var sourceID string
	s := MsgStruct{}
	err := json.Unmarshal([]byte(message), &s)
	if err != nil {
		callback.OnError(200, err.Error())
		return
	}
	maxSeq, err := getLocalMaxSeqModel()
	if err != nil {
		callback.OnError(201, err.Error())
		return
	}
	//如果删除的消息是最大的seq，则标记为已经删除的
	if maxSeq == s.Seq {
		err = setMessageStatus(s.ServerMsgID, MsgStatusHasDeleted)
		if err != nil {
			callback.OnError(202, err.Error())
			return
		}
	} else {
		if err = deleteMessageByMsgID(s.ServerMsgID); err != nil {
			callback.OnError(203, err.Error())
			return
		}
	}
	callback.OnSuccess("")
	if s.SessionType == GroupChatType {
		conversationID = GetConversationIDBySessionType(s.RecvID, GroupChatType)
		sourceID = s.RecvID
	} else if s.SessionType == SingleChatType {
		if s.SendID != LoginUid {
			conversationID = GetConversationIDBySessionType(s.SendID, SingleChatType)
			sourceID = s.SendID
		} else {
			conversationID = GetConversationIDBySessionType(s.RecvID, SingleChatType)
			sourceID = s.RecvID

		}
	}
	_, latestMsg := getConversationLatestMsgModel(conversationID)

	if s.ServerMsgID == latestMsg.ServerMsgID { //If the deleted message is the latest message of the conversation, update the latest message of the conversation
		err, list := getHistoryMessage(sourceID, s.SendTime+TimeOffset, 1)
		if err != nil {
			sdkLog("DeleteMessageFromLocalStorage database err:", err.Error())
		}
		if list == nil {
			conversation.LatestMsg = ""
			conversation.LatestMsgSendTime = getCurrentTimestampBySecond()
		} else {
			conversation.LatestMsg = structToJsonString(list[0])
			conversation.LatestMsgSendTime = list[0].SendTime

		}
		err = triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: AddConOrUpLatMsg, Args: conversation})
		if err != nil {
			sdkLog("DeleteMessageFromLocalStorage triggerCmdUpdateConversation err:", err.Error())
		}

	}

}
func InsertSingleMessageToLocalStorage(callback Base, message, userID, sender string) string {
	s := MsgStruct{}
	err := json.Unmarshal([]byte(message), &s)
	if err != nil {
		callback.OnError(200, err.Error())
		return ""
	}
	s.SendID = sender
	s.RecvID = userID
	//Generate client message primary key
	s.ClientMsgID = getMsgID(s.SendID)
	s.SendTime = getCurrentTimestampBySecond()
	if err = insertSendMessageToChatLog(&s); err != nil {
		callback.OnError(201, err.Error())
	} else {
		callback.OnSuccess("")
	}

	return s.ClientMsgID
}
func FindMessages(callback Base, messageIDList string) {
	var c []string
	err := json.Unmarshal([]byte(messageIDList), &c)
	if err != nil {
		callback.OnError(200, err.Error())
		fmt.Println("json PARSING ERROR", 111)
		log("json ")
	}
	err, list := getMultipleMessageModel(c)
	if err != nil {
		callback.OnError(203, err.Error())
	} else {
		if list != nil {
			callback.OnSuccess(structToJsonString(list))
		} else {
			callback.OnSuccess(structToJsonString([]MsgStruct{}))
		}
	}

}

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
	go func() {
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
	}()
}
func GetOneConversation(sourceID string, sessionType int, callback Base) {
	go func() {
		conversationID := GetConversationIDBySessionType(sourceID, sessionType)
		err, c := getOneConversationModel(conversationID)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			//
			if c.ConversationID == "" {

				c.ConversationID = conversationID
				c.ConversationType = sessionType
				c.RecvMsgOpt = 1

				switch sessionType {
				case SingleChatType:
					c.UserID = sourceID
					faceUrl, name, err := getUserNameAndFaceUrlByUid(sourceID)
					if err != nil {
						callback.OnError(301, err.Error())
						sdkLog("getUserNameAndFaceUrlByUid err:", err)
						return
					}
					c.ShowName = name
					c.FaceURL = faceUrl
				case GroupChatType:
					c.GroupID = sourceID
					faceUrl, name, err := getGroupNameAndFaceUrlByUid(sourceID)
					if err != nil {
						callback.OnError(301, err.Error())
						sdkLog("getGroupNameAndFaceUrlByUid err:", err)
					}
					c.ShowName = name
					c.FaceURL = faceUrl

				}
				err = addConversationOrUpdateLatestMsg(&c, conversationID)
				if err != nil {
					callback.OnError(301, err.Error())
					return
				}
				callback.OnSuccess(structToJsonString(c))

			} else {
				callback.OnSuccess(structToJsonString(c))
			}
		}
	}()
}

func GetMultipleConversation(conversationIDList string, callback Base) {
	go func() {
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
	}()
}
func DeleteConversation(conversationID string, callback Base) {
	go func() {
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
		//标记删除与此会话相关的消息
		err = setMessageStatusBySourceID(sourceID, MsgStatusHasDeleted, c.ConversationType)
		if err != nil {
			callback.OnError(202, err.Error())
			return
		}
		//重置该会话信息，空会话
		err = ResetConversation(conversationID)
		if err != nil {
			callback.OnError(203, err.Error())
			return
		} else {
			callback.OnSuccess("")
			_ = triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: ConAndUnreadChange})
		}
	}()
}
func SetConversationDraft(conversationID, draftText string, callback Base) {
	err := setConversationDraftModel(conversationID, draftText, getCurrentTimestampByNano())
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
	SdkInitManager.syncSeq2Msg()
}

func ForceSyncJoinedGroup() {
	groupManager.syncJoinedGroupInfo()
}

func ForceSyncJoinedGroupMember() {

	groupManager.syncJoinedGroupMember()
}

func ForceSyncGroupRequest() {
	groupManager.syncGroupRequest()
}

func ForceSyncApplyGroupRequest() {
	groupManager.syncApplyGroupRequest()
}

type SendMsgCallBack interface {
	Base
	OnProgress(progress int)
}

func CreateTextMessage(text string) string {
	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, Text)
	s.Content = text
	s.AtElem.AtUserList = []string{}
	return structToJsonString(s)
}
func CreateTextAtMessage(text, atUserList string) string {
	var users []string
	_ = json.Unmarshal([]byte(atUserList), &users)
	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, AtText)
	s.ForceList = users
	s.AtElem.Text = text
	s.AtElem.AtUserList = users
	s.Content = structToJsonString(s.AtElem)
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
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.VideoElem)
	return structToJsonString(s)
}

func CreateImageMessageFromFullPath(imageFullPath string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := fileTmpPath(imageFullPath) //a->b
		_, err := copyFile(imageFullPath, dstFile)
		sdkLog("copy file, ", imageFullPath, dstFile)
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
		sdkLog("getImageInfo err:", err.Error())
		return ""
	}
	s.PictureElem.SourcePicture.Width = imageInfo.Width
	s.PictureElem.SourcePicture.Height = imageInfo.Height
	s.PictureElem.SourcePicture.Type = imageInfo.Type
	s.PictureElem.SourcePicture.Size = imageInfo.Size
	wg.Wait()
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.PictureElem)
	return structToJsonString(s)
}

func CreateSoundMessageFromFullPath(soundPath string, duration int64) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := fileTmpPath(soundPath) //a->b
		_, err := copyFile(soundPath, dstFile)
		sdkLog("copy file, ", soundPath, dstFile)
		if err != nil {
			sdkLog("open file failed: ", err, soundPath)
		}
		wg.Done()
	}()
	sdkLog("init base info ")
	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, Sound)
	s.SoundElem.SoundPath = soundPath
	s.SoundElem.Duration = duration
	fi, err := os.Stat(s.SoundElem.SoundPath)
	if err != nil {
		sdkLog(err.Error(), s.SoundElem.SoundPath)
		return ""
	}
	s.SoundElem.DataSize = fi.Size()
	wg.Wait()
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.SoundElem)
	sdkLog("to string")
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
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.PictureElem)
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
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.SoundElem)
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
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.VideoElem)
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
	s.AtElem.AtUserList = []string{}
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
	go func() {
		c := ConversationStruct{
			ConversationType:  int(s.SessionType),
			RecvMsgOpt:        1,
			LatestMsgSendTime: s.CreateTime,
		}
		if receiver == "" && groupID == "" {
			callback.OnError(201, "args err")
			return
		} else if receiver == "" {
			s.SessionType = GroupChatType
			s.RecvID = groupID
			s.GroupID = groupID
			conversationID = GetConversationIDBySessionType(groupID, GroupChatType)
			c.GroupID = groupID
			faceUrl, name, err := getGroupNameAndFaceUrlByUid(groupID)
			if err != nil {
				sdkLog("getGroupNameAndFaceUrlByUid err:", err)
				callback.OnError(202, err.Error())
				return
			}
			c.ShowName = name
			c.FaceURL = faceUrl
			groupMemberList, err := getLocalGroupMemberListByGroupID(groupID)
			if err != nil {
				sdkLog("getLocalGroupMemberListByGroupID err:", err)
				callback.OnError(202, err.Error())
				return
			}
			isExistInGroup := func(target string, groupMemberList []groupMemberFullInfo) bool {

				for _, element := range groupMemberList {

					if target == element.UserId {
						return true
					}
				}
				return false

			}(s.SendID, groupMemberList)
			if !isExistInGroup {
				sdkLog("SendGroupMessage err:", "not exist in this group")
				callback.OnError(208, "not exist in this group")
				return
			}

		} else {
			s.SessionType = SingleChatType
			s.RecvID = receiver
			conversationID = GetConversationIDBySessionType(receiver, SingleChatType)
			c.UserID = receiver
			faceUrl, name, err := getUserNameAndFaceUrlByUid(receiver)
			if err != nil {
				sdkLog("getUserNameAndFaceUrlByUid err:", err)
				callback.OnError(301, err.Error())
				return
			}
			c.FaceURL = faceUrl
			c.ShowName = name
		}
		userInfo, err := getLoginUserInfoFromLocal()
		if err != nil {
			sdkLog("getLoginUserInfoFromLocal err:", err)
			return
		}
		s.SenderFaceURL = userInfo.Icon
		s.SenderNickName = userInfo.Name
		c.ConversationID = conversationID
		c.LatestMsg = structToJsonString(s)
		err = insertMessageToLocalOrUpdateContent(&s)
		if err != nil {
			sdkLog("insertMessageToLocalOrUpdateContent err:", err)
			callback.OnError(202, err.Error())
			return
		}
		_ = triggerCmdUpdateConversation(updateConNode{conversationID, AddConOrUpLatMsg,
			c})
		_ = triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
		var delFile []string
		switch s.ContentType {
		case Text:
		case AtText:
		case Picture:
			var sourcePath string
			if fileExist(s.PictureElem.SourcePath) {
				sourcePath = s.PictureElem.SourcePath
				delFile = append(delFile, fileTmpPath(s.PictureElem.SourcePath))
			} else {
				sourcePath = fileTmpPath(s.PictureElem.SourcePath)
				delFile = append(delFile, sourcePath)
			}
			sdkLog("file: ", sourcePath, delFile)
			sourceUrl, uuid, err := uploadImage(sourcePath, callback)
			if err != nil {
				sdkLog("oss Picture upload err", err.Error())
				callback.OnError(301, err.Error())
				sendMessageFailedHandle(&s, &c, conversationID)
				return
			} else {
				s.PictureElem.SourcePicture.Url = sourceUrl
				s.PictureElem.SourcePicture.UUID = uuid
				s.PictureElem.SnapshotPicture.Url = sourceUrl + "?imageView2/2/w/" + ZoomScale + "/h/" + ZoomScale
				s.PictureElem.SnapshotPicture.Width = int32(stringToInt(ZoomScale))
				s.PictureElem.SnapshotPicture.Height = int32(stringToInt(ZoomScale))
				s.Content = structToJsonString(s.PictureElem)
			}
		case Sound:
			var sourcePath string
			if fileExist(s.SoundElem.SoundPath) {
				sourcePath = s.SoundElem.SoundPath
				delFile = append(delFile, fileTmpPath(s.SoundElem.SoundPath))
			} else {
				sourcePath = fileTmpPath(s.SoundElem.SoundPath)
				delFile = append(delFile, sourcePath)
			}
			sdkLog("file: ", sourcePath, delFile)
			soundURL, uuid, err := uploadSound(sourcePath, callback)
			if err != nil {
				callback.OnError(301, err.Error())
				fmt.Println("oss Sound upload err", 111)
				sdkLog("uploadSound err:", err.Error())
				sendMessageFailedHandle(&s, &c, conversationID)
				return
			} else {
				s.SoundElem.SourceURL = soundURL
				s.SoundElem.UUID = uuid
				s.Content = structToJsonString(s.SoundElem)
			}
		case Video:
			var videoPath string
			var snapPath string
			if fileExist(s.VideoElem.VideoPath) {
				videoPath = s.VideoElem.VideoPath
				snapPath = s.VideoElem.SnapshotPath
				delFile = append(delFile, fileTmpPath(s.VideoElem.VideoPath))
				delFile = append(delFile, fileTmpPath(s.VideoElem.SnapshotPath))
			} else {
				videoPath = fileTmpPath(s.VideoElem.VideoPath)
				snapPath = fileTmpPath(s.VideoElem.SnapshotPath)
				delFile = append(delFile, videoPath)
				delFile = append(delFile, snapPath)
			}
			sdkLog("file: ", videoPath, snapPath, delFile)
			snapshotURL, snapshotUUID, videoURL, videoUUID, err := uploadVideo(videoPath, snapPath, callback)
			if err != nil {
				callback.OnError(301, err.Error())
				sdkLog("oss  Video upload err:", err.Error())
				sendMessageFailedHandle(&s, &c, conversationID)
				return
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
				sendMessageFailedHandle(&s, &c, conversationID)
				return

			} else {
				s.FileElem.SourceURL = fileURL
				s.FileElem.UUID = fileUUID
				s.Content = structToJsonString(s.FileElem)
			}
		default:
			callback.OnError(2038, "Not currently supported ")
			sdkLog("Not currently supported ", s.ContentType)
			return
		}
		//Store messages to local database
		err = insertMessageToLocalOrUpdateContent(&s)
		if err != nil {
			callback.OnError(202, err.Error())
			return
		}
		//Protocol conversion
		a.ReqIdentifier = 1003
		a.PlatformID = s.PlatformID
		a.SendID = s.SendID
		a.SenderFaceURL = s.SenderFaceURL
		a.SenderNickName = s.SenderNickName
		a.OperationID = operationIDGenerator()
		a.Data.SessionType = s.SessionType
		a.Data.MsgFrom = s.MsgFrom
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
			sendMessageFailedHandle(&s, &c, conversationID)
		} else if err = json.Unmarshal(bMsg, &r); err != nil {
			callback.OnError(200, err.Error()+"  "+string(bMsg))
			sendMessageFailedHandle(&s, &c, conversationID)
		} else {
			if r.ErrCode != 0 {
				callback.OnError(r.ErrCode, r.ErrMsg)
				sendMessageFailedHandle(&s, &c, conversationID)
			} else {
				callback.OnSuccess("")
				callback.OnProgress(100)

				for _, v := range delFile {
					err := os.Remove(v)
					if err != nil {
						sdkLog("remove failed,", err.Error(), v)
					}
					sdkLog("remove file: ", v)
				}
				_ = updateMessageTimeAndMsgIDStatus(r.Data.ClientMsgID, r.Data.SendTime, MsgStatusSendSuccess)
				s.ServerMsgID = r.Data.ServerMsgID
				s.SendTime = r.Data.SendTime
				s.Status = MsgStatusSendSuccess
				c.LatestMsg = structToJsonString(s)
				c.LatestMsgSendTime = s.SendTime
				_ = triggerCmdUpdateConversation(updateConNode{conversationID, AddConOrUpLatMsg,
					c})
				_ = triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
			}
		}
	}()
	return s.ClientMsgID
}

func GetHistoryMessageList(callback Base, getMessageOptions string) {
	go func() {
		fmt.Println("GetHistoryMessageList", getMessageOptions)
		var sourceID string
		var conversationID string
		var startTime int64
		var latestMsg MsgStruct
		var sessionType int
		p := PullMsgReq{}
		err := json.Unmarshal([]byte(getMessageOptions), &p)
		if err != nil {
			callback.OnError(200, err.Error())
			return
		}
		if p.UserID == "" {
			sourceID = p.GroupID
			conversationID = GetConversationIDBySessionType(sourceID, GroupChatType)
			sessionType = GroupChatType
		} else {
			sourceID = p.UserID
			conversationID = GetConversationIDBySessionType(sourceID, SingleChatType)
			sessionType = SingleChatType
		}
		if p.StartMsg == nil {
			err, m := getConversationLatestMsgModel(conversationID)
			if err != nil {
				callback.OnError(200, err.Error())
				return
			}
			if m == "" {
				startTime = 0
			} else {
				err := json.Unmarshal([]byte(m), &latestMsg)
				if err != nil {
					fmt.Println("get history err :", err)
					callback.OnError(200, err.Error())
					return
				}
				startTime = latestMsg.SendTime + TimeOffset
			}

		} else {
			startTime = p.StartMsg.SendTime
		}
		fmt.Println("sourceID:", sourceID, "startTime:", startTime, "count:", p.Count)
		err, list := getHistoryMessage(sourceID, startTime, p.Count, sessionType)
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
	}()
}
func RevokeMessage(callback Base, message string) {
	go func() {
		var receiver, groupID string
		c := MsgStruct{}
		err := json.Unmarshal([]byte(message), &c)
		if err != nil {
			callback.OnError(200, err.Error())
			return
		}
		s, err := getOneMessage(c.ClientMsgID)
		if err != nil {
			callback.OnError(201, "getOneMessage err")
			return
		}
		if s.Status != MsgStatusSendSuccess {
			callback.OnError(201, "only send success message can be revoked")
			return
		}
		sdkLog("test data", s)
		//Send message internally
		switch s.SessionType {
		case SingleChatType:
			receiver = s.RecvID
		case GroupChatType:
			groupID = s.RecvID
		default:
			callback.OnError(200, "args err")
		}
		s.Content = s.ClientMsgID
		s.ClientMsgID = getMsgID(s.SendID)
		s.ContentType = Revoke
		err = autoSendMsg(s, receiver, groupID, false, true, false)
		if err != nil {
			sdkLog("autoSendMsg revokeMessage err:", err.Error())
			callback.OnError(300, err.Error())

		} else {
			err = setMessageStatus(s.Content, MsgStatusRevoked)
			if err != nil {
				sdkLog("setLocalMessageStatus revokeMessage err:", err.Error())
				callback.OnError(300, err.Error())
			} else {
				callback.OnSuccess("")
			}
		}
	}()
}
func TypingStatusUpdate(receiver, msgTip string) {
	go func() {
		s := MsgStruct{}
		initBasicInfo(&s, UserMsgType, Typing)
		s.Content = msgTip
		err := autoSendMsg(&s, receiver, "", true, false, false)
		if err != nil {
			sdkLog("TypingStatusUpdate err:", err)
		} else {
			sdkLog("TypingStatusUpdate success!!!")
		}
	}()
}

func MarkC2CMessageAsRead(callback Base, receiver string, msgList string) {
	go func() {
		s := MsgStruct{}
		initBasicInfo(&s, UserMsgType, HasReadReceipt)
		s.Content = msgList
		sdkLog("MarkC2CMessageAsRead: send Message")
		err := autoSendMsg(&s, receiver, "", false, false, false)
		if err != nil {
			sdkLog("MarkC2CMessageAsRead  err:", err.Error())
			callback.OnError(300, err.Error())
		} else {
			callback.OnSuccess("")
		}
	}()
}

func MarkSingleMessageHasRead(callback Base, userID string) {
	//go func() {
	//	conversationID := GetConversationIDBySessionType(userID, SingleChatType)
	//	if err := setSingleMessageHasRead(userID); err != nil {
	//		callback.OnError(201, err.Error())
	//	} else {
	//		callback.OnSuccess("")
	//		triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: UnreadCountSetZero})
	//		n := NotificationContent{
	//			IsDisplay:   0,
	//			DefaultTips: userID,
	//			Detail:      userID,
	//		}
	//		msg := createTextSystemMessage(n, HasReadReceipt)
	//		autoSendMsg(msg, userID, "", false, false, false)
	//	}
	//}()

	go func() {
		conversationID := GetConversationIDBySessionType(userID, SingleChatType)
		if err := setSingleMessageHasRead(userID); err != nil {
			callback.OnError(201, err.Error())
		} else {
			callback.OnSuccess("")
			triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: UnreadCountSetZero})
		}
	}()
}
func MarkGroupMessageHasRead(callback Base, groupID string) {
	go func() {
		conversationID := GetConversationIDBySessionType(groupID, GroupChatType)
		if err := setGroupMessageHasRead(groupID); err != nil {
			callback.OnError(201, err.Error())
		} else {
			callback.OnSuccess("")
			triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: UnreadCountSetZero})
		}
	}()
}

func DeleteMessageFromLocalStorage(callback Base, message string) {
	go func() {
		var conversation ConversationStruct
		var latestMsg MsgStruct
		var conversationID string
		var sourceID string
		s := MsgStruct{}
		err := json.Unmarshal([]byte(message), &s)
		if err != nil {
			callback.OnError(200, err.Error())
			return
		}
		err = setMessageStatus(s.ClientMsgID, MsgStatusHasDeleted)
		if err != nil {
			callback.OnError(202, err.Error())
			return
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
		_, m := getConversationLatestMsgModel(conversationID)
		if m != "" {
			err := json.Unmarshal([]byte(m), &latestMsg)
			if err != nil {
				fmt.Println("DeleteMessage err :", err)
				callback.OnError(200, err.Error())
				return
			}
		} else {
			sdkLog("err ,conversation has been deleted")
		}

		if s.ClientMsgID == latestMsg.ClientMsgID { //If the deleted message is the latest message of the conversation, update the latest message of the conversation
			err, list := getHistoryMessage(sourceID, s.SendTime+TimeOffset, 1, int(s.SessionType))
			if err != nil {
				sdkLog("DeleteMessageFromLocalStorage database err:", err.Error())
			}
			conversation.ConversationID = conversationID
			if list == nil {
				conversation.LatestMsg = ""
				conversation.LatestMsgSendTime = getCurrentTimestampByNano()
			} else {
				conversation.LatestMsg = structToJsonString(list[0])
				conversation.LatestMsgSendTime = list[0].SendTime
			}
			err = triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: AddConOrUpLatMsg, Args: conversation})
			if err != nil {
				sdkLog("DeleteMessageFromLocalStorage triggerCmdUpdateConversation err:", err.Error())
			}
			_ = triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})

		}
	}()
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
	s.SendTime = getCurrentTimestampByNano()
	if err = insertMessageToLocalOrUpdateContent(&s); err != nil {
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

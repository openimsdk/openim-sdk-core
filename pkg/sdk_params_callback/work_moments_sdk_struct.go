package sdk_params_callback

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
)

type WorkMomentsNotificationCallback string

const ClearWorkMomentsMessageCallback = constant.SuccessCallbackDefault

type GetWorkMomentsUnReadCountCallback db.LocalWorkMomentsNotificationUnreadCount

type Comment struct {
	UserID        string `json:"userID"`
	UserName      string `json:"userName"`
	FaceURL       string `json:"faceURL"`
	ReplyUserID   string `json:"replyUserID"`
	ReplyUserName string `json:"replyUserName"`
	ContentID     string `json:"contentID"`
	Content       string `json:"content"`
	CreateTime    int32  `json:"createTime"`
}

type WorkMomentNotificationMsg struct {
	NotificationMsgType int32  `json:"notificationMsgType"`
	ReplyUserName       string `json:"replyUserName"`
	ReplyUserID         string `json:"replyUserID"`
	Content             string `json:"content"`
	ContentID           string `json:"contentID"`
	WorkMomentID        string `json:"workMomentID"`
	UserID              string `json:"userID"`
	UserName            string `json:"userName"`
	FaceURL             string `json:"faceURL"`
	WorkMomentContent   string `json:"workMomentContent"`
	CreateTime          int32  `json:"createTime"`
}

type GetWorkMomentsNotificationCallback []WorkMomentNotificationMsg

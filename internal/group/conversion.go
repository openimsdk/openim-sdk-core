package group

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/pkg/db/model_struct"
)

func ServerGroupToLocalGroup(info *sdkws.GroupInfo) *model_struct.LocalGroup {
	return &model_struct.LocalGroup{
		GroupID:                info.GroupID,
		GroupName:              info.GroupName,
		Notification:           info.Notification,
		Introduction:           info.Introduction,
		FaceURL:                info.FaceURL,
		CreateTime:             info.CreateTime,
		Status:                 info.Status,
		CreatorUserID:          info.CreatorUserID,
		GroupType:              info.GroupType,
		OwnerUserID:            info.OwnerUserID,
		MemberCount:            int32(info.MemberCount),
		Ex:                     info.Ex,
		NeedVerification:       info.NeedVerification,
		LookMemberInfo:         info.LookMemberInfo,
		ApplyMemberFriend:      info.ApplyMemberFriend,
		NotificationUpdateTime: info.NotificationUpdateTime,
		NotificationUserID:     info.NotificationUserID,
		//AttachedInfo:           info.AttachedInfo, // TODO
	}
}

func ServerGroupRequestToLocalGroupRequest(info *sdkws.GroupRequest) *model_struct.LocalGroupRequest {
	return &model_struct.LocalGroupRequest{
		GroupID:       info.GroupInfo.GroupID,
		GroupName:     info.GroupInfo.GroupName,
		Notification:  info.GroupInfo.Notification,
		Introduction:  info.GroupInfo.Introduction,
		GroupFaceURL:  info.GroupInfo.FaceURL,
		CreateTime:    info.GroupInfo.CreateTime,
		Status:        info.GroupInfo.Status,
		CreatorUserID: info.GroupInfo.CreatorUserID,
		GroupType:     info.GroupInfo.GroupType,
		OwnerUserID:   info.GroupInfo.OwnerUserID,
		MemberCount:   int32(info.GroupInfo.MemberCount),
		UserID:        info.UserInfo.UserID,
		Nickname:      info.UserInfo.Nickname,
		UserFaceURL:   info.UserInfo.FaceURL,
		Gender:        info.UserInfo.Gender,
		HandleResult:  info.HandleResult,
		ReqMsg:        info.ReqMsg,
		HandledMsg:    info.HandleMsg,
		ReqTime:       info.ReqTime,
		HandleUserID:  info.HandleUserID,
		HandledTime:   info.HandleTime,
		Ex:            info.Ex,
		//AttachedInfo:  info.AttachedInfo,
		JoinSource:    info.JoinSource,
		InviterUserID: info.InviterUserID,
	}
}

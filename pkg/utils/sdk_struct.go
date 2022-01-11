package utils

import (
	"github.com/gorilla/websocket"
	"open_im_sdk/internal/controller/init"
	"open_im_sdk/pkg/server_api_params"
)

//fixme---------------UiParam-->client---------------

type paramsUiAddFriend struct {
	UID        string `json:"uid" validate:"required"`
	ReqMessage string `json:"reqMessage"`
}
type ui2AcceptFriend struct {
	UID string `json:"uid" validate:"required"`
}
type uid2Comment struct {
	Uid     string `json:"uid" validate:"required"`
	Comment string `json:"comment"`
}

type delUid struct {
	Uid string `json:"uid" validate:"required"`
}
type ui2UpdateUserInfo struct {
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Gender int32  `json:"gender"`
	Mobile string `json:"mobile"`
	Birth  string `json:"birth"`
	Email  string `json:"email"`
	Ex     string `json:"ex"`
}
type ui2ClientCommonReq struct {
	UidList []string `json:"uidList" validate:"required"`
}

//fixme ----------msg to ui--------------------
type acceptOrRefuseFriendMsgCallback struct {
	Uid        string `json:"uid"`
	ResUltCode string `json:"resUltCode"`
	ResultInfo string `json:"resultInfo"`
}

//fixme---------------user struct--------------
type getUserInfoResp struct {
	Data    []userInfo `json:"data"`
	ErrCode int32      `json:"errCode"`
	ErrMsg  string     `json:"errMsg"`
}
type userInfo struct {
	Uid    string `json:"uid"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Gender int32  `json:"gender"`
	Mobile string `json:"mobile"`
	Birth  string `json:"birth"`
	Email  string `json:"email"`
	Ex     string `json:"ex"`
}

type paramsGetUserInfo struct {
	UidList     []string `json:"uidList"`
	OperationID string   `json:"operationID"`
}

type paramsUpdateUserInfo struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Gender      int32  `json:"gender"`
	Mobile      string `json:"mobile"`
	Birth       string `json:"birth"`
	Email       string `json:"email"`
	Ex          string `json:"ex"`
	OperationID string `json:"operationID"`
}

//fixme------------friend struct---------------
type paramsAddFriend struct {
	OperationID string `json:"operationID" binding:"required"`
	UID         string `json:"uid" binding:"required"`
	ReqMessage  string `json:"reqMessage"`
}

type paramsDeleteFriend struct {
	Uid         string `json:"uid"`
	OperationID string `json:"operationID"`
}
type paramsCommonReq struct {
	OperationID string `json:"operationID"`
}
type paramsAddFriendResponse struct {
	Uid         string `json:"uid"`
	OperationID string `json:"operationID"`
	Flag        int    `json:"flag"`
}
type paramsSetFriendInfo struct {
	Uid         string `json:"uid"`
	OperationID string `json:"operationID"`
	Comment     string `json:"comment"`
}
type paramsAddBlackList struct {
	OperationID string `json:"operationID" binding:"required"`
	UID         string `json:"uid" binding:"required"`
}

type paramsRemoveBlackList struct {
	OperationID string `json:"operationID" binding:"required"`
	UID         string `json:"uid" binding:"required"`
}

type commonResp struct {
	ErrCode int32  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}
type getFriendListResp struct {
	Data    []friendInfo `json:"data"`
	ErrCode int          `json:"errCode"`
	ErrMsg  string       `json:"errMsg"`
}
type getFriendResp struct {
	Data    friendInfo `json:"data"`
	ErrCode int        `json:"errCode"`
	ErrMsg  string     `json:"errMsg"`
}

type getBlackListResp struct {
	Data    []userInfo `json:"data"`
	ErrCode int32      `json:"errCode"`
	ErrMsg  string     `json:"errMsg"`
}
type getFriendApplyListResp struct {
	Data    []applyUserInfo `json:"data"`
	ErrCode int             `json:"errCode"`
	ErrMsg  string          `json:"errMsg"`
}
type friendInfo struct {
	UID           string `json:"uid"`
	Name          string `json:"name"`
	Icon          string `json:"icon"`
	Gender        int32  `json:"gender"`
	Mobile        string `json:"mobile"`
	Birth         string `json:"birth"`
	Email         string `json:"email"`
	Ex            string `json:"ex"`
	Comment       string `json:"comment"`
	IsInBlackList int32  `json:"isInBlackList"`
}

type applyUserInfo struct {
	Uid        string `json:"uid"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	Gender     int32  `json:"gender"`
	Mobile     string `json:"mobile"`
	Birth      string `json:"birth"`
	Email      string `json:"email"`
	Ex         string `json:"ex"`
	ReqMessage string `json:"reqMessage"`
	ApplyTime  string `json:"applyTime"`
	Flag       int32  `json:"flag"`
}
type Uid2Flag struct {
	Uid  string `json:"uid"`
	Flag int32  `json:"flag"`
}

//fixme--------------message struct-----------------

type MessageReceipt struct {
	UserID      string   `json:"uid"`
	MsgIdList   []string `json:"msgIDList"`
	ReadTime    int64    `json:"readTime"`
	MsgFrom     int32    `json:"msgFrom"`
	ContentType int32    `json:"contentType"`
	SessionType int32    `json:"sessionType"`
}

//type MsgData struct {
//	SendID           string
//	RecvID           string
//	SessionType      int32
//	MsgFrom          int32
//	ContentType      int32
//	ServerMsgID      string
//	Content          string
//	SendTime         int64
//	Seq              int64
//	SenderPlatformID int32
//	SenderNickName   string
//	SenderFaceURL    string
//	ClientMsgID      string
//}

type WsMsgData struct {
	PlatformID  int32                  `mapstructure:"platformID" validate:"required"`
	SessionType int32                  `mapstructure:"sessionType" validate:"required"`
	MsgFrom     int32                  `mapstructure:"msgFrom" validate:"required"`
	ContentType int32                  `mapstructure:"contentType" validate:"required"`
	RecvID      string                 `mapstructure:"recvID" validate:"required"`
	ForceList   []string               `mapstructure:"forceList" validate:"required"`
	Content     string                 `mapstructure:"content" validate:"required"`
	Options     map[string]interface{} `mapstructure:"options" validate:"required"`
	ClientMsgID string                 `mapstructure:"clientMsgID" validate:"required"`
	OfflineInfo map[string]interface{} `mapstructure:"offlineInfo" validate:"required"`
	Ext         map[string]interface{} `mapstructure:"ext"`
}
type WsSubMsg struct {
	SendTime    int64  `json:"sendTime"`
	ServerMsgID string `json:"serverMsgID"`
	ClientMsgID string `json:"clientMsgID"`
}

//type GeneralWsResp struct {
//	ReqIdentifier int    `json:"reqIdentifier"`
//	ErrCode       int    `json:"errCode"`
//	ErrMsg        string `json:"errMsg"`
//	MsgIncr       string `json:"msgIncr"`
//	OperationID   string `json:"operationID"`
//	Data          []byte `json:"data"`
//}

//type GeneralWsReq struct {
//	ReqIdentifier int32  `json:"reqIdentifier"`
//	Token         string `json:"token"`
//	SendID        string `json:"sendID"`
//	OperationID   string `json:"operationID"`
//	MsgIncr       string `json:"msgIncr"`
//	Data          []byte `json:"data"`
//}

type PullUserMsgResp struct {
	ErrCode       int                       `json:"errCode"`
	ErrMsg        string                    `json:"errMsg"`
	ReqIdentifier int                       `json:"reqIdentifier"`
	MsgIncr       int                       `json:"msgIncr"`
	Data          paramsPullUserMsgDataResp `json:"data"`
}
type paramsPullUserMsgDataResp struct {
	Group  []*server_api_params.GatherFormat `json:"group"`
	MaxSeq int64                             `json:"maxSeq"`
	MinSeq int64                             `json:"minSeq"`
	Single []*server_api_params.GatherFormat `json:"single"`
}

type ArrMsg struct {
	SingleData []server_api_params.MsgData
	GroupData  []server_api_params.MsgData
}

type IMConfig struct {
	Platform int32  `json:"platform"`
	ApiAddr  string `json:"api_addr"`
	WsAddr   string `json:"ws_addr"`
	DbDir    string `json:"db_dir"`
	logLevel int32  `json:"log_level"`
}

type IMManager struct {
	conn *websocket.Conn
	cb   init.IMSDKListener

	LoginState int
}

type paramsPullUserMsgDataReq struct {
	SeqBegin int64 `json:"seqBegin"`
	SeqEnd   int64 `json:"seqEnd"`
}

type paramsPullUserMsg struct {
	ReqIdentifier int                      `json:"reqIdentifier"`
	OperationID   string                   `json:"operationID"`
	SendID        string                   `json:"sendID"`
	Data          paramsPullUserMsgDataReq `json:"data"`
}

type paramsPullUserMsgBySeq struct {
	ReqIdentifier int     `json:"reqIdentifier"`
	OperationID   string  `json:"operationID"`
	SendID        string  `json:"sendID"`
	SeqList       []int64 `json:"seqList" binding:"required"`
}

type paramsPullUserGroupMsgDataResp struct {
	paramsPullUserSingleMsgDataResp
}
type paramsPullUserSingleList struct {
	SendID           string `json:"sendID"`
	RecvID           string `json:"recvID"`
	SendTime         int64  `json:"sendTime"`
	ContentType      int32  `json:"contentType"`
	MsgFrom          int32  `json:"msgFrom"`
	Content          string `json:"content"`
	Seq              int64  `json:"seq"`
	ServerMsgID      string `json:"serverMsgID"`
	SenderPlatformID int32  `json:"senderPlatformID"`
	SenderNickName   string `json:"senderNickName"`
	SenderFaceURL    string `json:"senderFaceUrl"`
	ClientMsgID      string `json:"clientMsgID"`
}

type paramsPullUserSingleMsgDataResp struct {
	ID   string                     `json:"id"`
	List []paramsPullUserSingleList `json:"list"`
}

//type paramsPullUserMsgDataResp struct {
//	Group  []*server_api_params.GatherFormat `json:"group"`
//	MaxSeq int64                             `json:"maxSeq"`
//	MinSeq int64                             `json:"minSeq"`
//	Single []*server_api_params.GatherFormat `json:"single"`
//}
//
//type PullUserMsgResp struct {
//	ErrCode       int                       `json:"errCode"`
//	ErrMsg        string                    `json:"errMsg"`
//	ReqIdentifier int                       `json:"reqIdentifier"`
//	MsgIncr       int                       `json:"msgIncr"`
//	Data          paramsPullUserMsgDataResp `json:"data"`
//}

type paramsNewestSeqReq struct {
	ReqIdentifier int    `json:"reqIdentifier"`
	OperationID   string `json:"operationID"`
	SendID        string `json:"sendID"`
	MsgIncr       int    `json:"msgIncr"`
}

type paramsNewestSeqDataResp struct {
	Seq    int64 `json:"seq"`
	MinSeq int64 `json:"minSeq"`
}
type paramsNewestSeqResp struct {
	ErrCode       int                     `json:"errCode"`
	ErrMsg        string                  `json:"errMsg"`
	ReqIdentifier int                     `json:"reqIdentifier"`
	MsgIncr       int                     `json:"msgIncr"`
	Data          paramsNewestSeqDataResp `json:"data"`
}

type paramsTencentOssCredentialReq struct {
	OperationID string `json:"operationID"`
}

type paramsTencentOssCredentialResp struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Bucket  string `json:"bucket"`
	Region  string `json:"region"`
	Data    struct {
		ExpiredTime int64
		Expiration  string
		StartTime   int64
		RequestId   string
		Credentials struct {
			TmpSecretId  string
			TmpSecretKey string
			Token        string
		}
	} `json:"data"`
}

////////////////////////// message/////////////////////////

type WsSendMsgResp struct {
	ServerMsgID string `json:"serverMsgID"`
	ClientMsgID string `json:"clientMsgID"`
	SendTime    int64  `json:"sendTime"`
}

type PullMsgReq struct {
	UserID   string     `json:"userID"`
	GroupID  string     `json:"groupID"`
	StartMsg *MsgStruct `json:"startMsg"`
	Count    int        `json:"count"`
}
type SendMsgRespFromServer struct {
	ErrCode       int    `json:"errCode"`
	ErrMsg        string `json:"errMsg"`
	ReqIdentifier int    `json:"reqIdentifier"`
	Data          struct {
		ServerMsgID string `json:"serverMsgID"`
		ClientMsgID string `json:"clientMsgID"`
		SendTime    int64  `json:"sendTime"`
	}
}
type paramsUserSendMsg struct {
	SenderPlatformID int32  `json:"senderPlatformID" binding:"required"`
	SendID           string `json:"sendID" binding:"required"`
	SenderNickName   string `json:"senderNickName"`
	SenderFaceURL    string `json:"senderFaceUrl"`
	OperationID      string `json:"operationID" binding:"required"`
	Data             struct {
		SessionType int32                             `json:"sessionType" binding:"required"`
		MsgFrom     int32                             `json:"msgFrom" binding:"required"`
		ContentType int32                             `json:"contentType" binding:"required"`
		RecvID      string                            `json:"recvID" `
		GroupID     string                            `json:"groupID" `
		ForceList   []string                          `json:"forceList"`
		Content     []byte                            `json:"content" binding:"required"`
		Options     map[string]bool                   `json:"options" `
		ClientMsgID string                            `json:"clientMsgID" binding:"required"`
		CreateTime  int64                             `json:"createTime" binding:"required"`
		OffLineInfo server_api_params.OfflinePushInfo `json:"offlineInfo" `
	}
}
type ImageInfo struct {
	Width  int32  `json:"x"`
	Height int32  `json:"y"`
	Type   string `json:"type"`
	Size   int64  `json:"size"`
}
type PictureBaseInfo struct {
	UUID   string `json:"uuid"`
	Type   string `json:"type"`
	Size   int64  `json:"size"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
	Url    string `json:"url"`
}
type SoundBaseInfo struct {
	UUID      string `json:"uuid"`
	SoundPath string `json:"soundPath"`
	SourceURL string `json:"sourceUrl"`
	DataSize  int64  `json:"dataSize"`
	Duration  int64  `json:"duration"`
}
type VideoBaseInfo struct {
	VideoPath      string `json:"videoPath"`
	VideoUUID      string `json:"videoUUID"`
	VideoURL       string `json:"videoUrl"`
	VideoType      string `json:"videoType"`
	VideoSize      int64  `json:"videoSize"`
	Duration       int64  `json:"duration"`
	SnapshotPath   string `json:"snapshotPath"`
	SnapshotUUID   string `json:"snapshotUUID"`
	SnapshotSize   int64  `json:"snapshotSize"`
	SnapshotURL    string `json:"snapshotUrl"`
	SnapshotWidth  int32  `json:"snapshotWidth"`
	SnapshotHeight int32  `json:"snapshotHeight"`
}
type FileBaseInfo struct {
	FilePath  string `json:"filePath"`
	UUID      string `json:"uuid"`
	SourceURL string `json:"sourceUrl"`
	FileName  string `json:"fileName"`
	FileSize  int64  `json:"fileSize"`
}
type MsgStruct struct {
	ClientMsgID      string                            `json:"clientMsgID"`
	ServerMsgID      string                            `json:"serverMsgID"`
	CreateTime       int64                             `json:"createTime"`
	SendTime         int64                             `json:"sendTime"`
	SessionType      int32                             `json:"sessionType"`
	SendID           string                            `json:"sendID"`
	RecvID           string                            `json:"recvID"`
	MsgFrom          int32                             `json:"msgFrom"`
	ContentType      int32                             `json:"contentType"`
	SenderPlatformID int32                             `json:"platformID"`
	ForceList        []string                          `json:"forceList"`
	SenderNickname   string                            `json:"senderNickname"`
	SenderFaceURL    string                            `json:"senderFaceUrl"`
	GroupID          string                            `json:"groupID"`
	Content          string                            `json:"content"`
	Seq              int64                             `json:"seq"`
	IsRead           bool                              `json:"isRead"`
	Status           int32                             `json:"status"`
	Remark           string                            `json:"remark"`
	OfflinePush      server_api_params.OfflinePushInfo `json:"offlinePush"`
	PictureElem      struct {
		SourcePath      string          `json:"sourcePath"`
		SourcePicture   PictureBaseInfo `json:"sourcePicture"`
		BigPicture      PictureBaseInfo `json:"bigPicture"`
		SnapshotPicture PictureBaseInfo `json:"snapshotPicture"`
	} `json:"pictureElem"`
	SoundElem struct {
		UUID      string `json:"uuid"`
		SoundPath string `json:"soundPath"`
		SourceURL string `json:"sourceUrl"`
		DataSize  int64  `json:"dataSize"`
		Duration  int64  `json:"duration"`
	} `json:"soundElem"`
	VideoElem struct {
		VideoPath      string `json:"videoPath"`
		VideoUUID      string `json:"videoUUID"`
		VideoURL       string `json:"videoUrl"`
		VideoType      string `json:"videoType"`
		VideoSize      int64  `json:"videoSize"`
		Duration       int64  `json:"duration"`
		SnapshotPath   string `json:"snapshotPath"`
		SnapshotUUID   string `json:"snapshotUUID"`
		SnapshotSize   int64  `json:"snapshotSize"`
		SnapshotURL    string `json:"snapshotUrl"`
		SnapshotWidth  int32  `json:"snapshotWidth"`
		SnapshotHeight int32  `json:"snapshotHeight"`
	} `json:"videoElem"`
	FileElem struct {
		FilePath  string `json:"filePath"`
		UUID      string `json:"uuid"`
		SourceURL string `json:"sourceUrl"`
		FileName  string `json:"fileName"`
		FileSize  int64  `json:"fileSize"`
	} `json:"fileElem"`
	MergeElem struct {
		Title        string       `json:"title"`
		AbstractList []string     `json:"abstractList"`
		MultiMessage []*MsgStruct `json:"multiMessage"`
	} `json:"mergeElem"`
	AtElem struct {
		Text       string   `json:"text"`
		AtUserList []string `json:"atUserList"`
		IsAtSelf   bool     `json:"isAtSelf"`
	} `json:"atElem"`
	LocationElem struct {
		Description string  `json:"description"`
		Longitude   float64 `json:"longitude"`
		Latitude    float64 `json:"latitude"`
	} `json:"locationElem"`
	CustomElem struct {
		Data        string `json:"data"`
		Description string `json:"description"`
		Extension   string `json:"extension"`
	} `json:"customElem"`
	QuoteElem struct {
		Text         string     `json:"text"`
		QuoteMessage *MsgStruct `json:"quoteMessage"`
	} `json:"quoteElem"`
	//RevokeMessage struct {
	//	ServerMsgID    string `json:"serverMsgID"`
	//	SendID         string `json:"sendID"`
	//	SenderNickname string `json:"senderNickname"`
	//	RecvID         string `json:"recvID"`
	//	GroupID        string `json:"groupID"`
	//	ContentType    int32  `json:"contentType"`
	//	SendTime       int64  `json:"sendTime"`
	//}
}

////////////////////////// group/////////////////////////

type changeGroupInfo struct {
	data       groupInfo `json:"data"`
	changeType int32     `json:"changeType"`
}

type groupInfo struct {
	GroupId      string `json:"groupID"`
	GroupName    string `json:"groupName"`
	Notification string `json:"notification"`
	Introduction string `json:"introduction"`
	FaceUrl      string `json:"faceUrl"`
	Ex           string `json:"ex"`
	OwnerId      string `json:"ownerId"`
	CreateTime   uint64 `json:"createTime"`
	MemberCount  uint32 `json:"memberCount"`
}

type groupMemberFullInfo struct {
	GroupId  string `json:"groupID"`
	UserId   string `json:"userId"`
	Role     int    `json:"role"`
	JoinTime uint64 `json:"joinTime"`
	NickName string `json:"nickName"`
	FaceUrl  string `json:"faceUrl"`
}

type groupApplication struct {
	GroupId          string `json:"groupID"`
	FromUser         string `json:"fromUserID"`
	FromUserNickName string `json:"fromUserNickName"`
	FromUserFaceUrl  string `json:"fromUserFaceUrl"`
	ToUser           string `json:"toUserID"`
	AddTime          int    `json:"addTime"`
	RequestMsg       string `json:"requestMsg"`
	HandledMsg       string `json:"handledMsg"`
	Type             int    `json:"type"`
	HandleStatus     int    `json:"handleStatus"`
	HandleResult     int    `json:"handleResult"`
}

type createGroupMemberInfo struct {
	Uid     string `json:"uid"`
	SetRole int32  `json:"setRole"`
}

type idResult struct {
	UId    string `json:"uid"`
	Result int32  `json:"result"`
}
type groupMemberInfoResult struct {
	Nextseq int32                 `json:"nextSeq"`
	Data    []groupMemberFullInfo `json:"data"`
	commonResp
}

type getGroupMemberListResult struct {
	NextSeq int32                 `json:"nextSeq"`
	Data    []groupMemberFullInfo `json:"data"`
}

type groupMemberOperationResult struct {
	commonResp
	createGroupMemberInfo
}

type groupApplicationResult struct {
	UnReadCount          int                `json:"count"`
	GroupApplicationList []GroupReqListInfo `json:"user"`
}

type updateGroupNode struct {
	groupId string
	Action  int //
	Args    interface{}
}

type createGroupArgs struct {
	uIdCreator     string
	initMemberList []groupMemberFullInfo
}

type joinGroupArgs struct {
	applyUser groupMemberFullInfo
	reason    string
}

type quiteGroupArgs struct {
	quiteUser groupMemberFullInfo
}

type setGroupInfoArgs struct {
	group groupInfo
}

type kickGroupAgrs struct {
	kickedList []groupMemberFullInfo
	op         groupMemberFullInfo
}

type transferGroupArgs struct {
	oldOwner groupMemberFullInfo
	newOwner groupMemberFullInfo
}

type inviteUserToGroupArgs struct {
	op      groupMemberFullInfo
	invited []groupMemberFullInfo
	reason  string
}

type applyGroupProcessedArgs struct {
	op        groupMemberFullInfo
	applyList []applyProcessed
}

type applyProcessed struct {
	member groupMemberFullInfo
	reason string
}

//fixme------------------group--------------------

type createGroupReq struct {
	MemberList   []createGroupMemberInfo `json:"memberList"`
	GroupName    string                  `json:"groupName"`
	Introduction string                  `json:"introduction"`
	Notification string                  `json:"notification"`
	FaceUrl      string                  `json:"faceUrl"`
	OperationID  string                  `json:"operationID"`
	Ex           string                  `json:"ex"`
}

type createGroupResp struct {
	ErrCode int       `json:"errCode"`
	ErrMsg  string    `json:"errMsg"`
	Data    groupInfo `json:"data"`
}

type setGroupInfoReq struct {
	GroupId      string `json:"groupID"`
	GroupName    string `json:"groupName"`
	Notification string `json:"notification"`
	Introduction string `json:"introduction"`
	FaceUrl      string `json:"faceUrl"`
	OperationID  string `json:"operationID"`
}

type joinGroupReq struct {
	GroupID     string `json:"groupID"`
	Message     string `json:"message"`
	OperationID string `json:"operationID"`
}

type quitGroupReq struct {
	GroupID     string `json:"groupID"`
	OperationID string `json:"operationID"`
}

type getGroupsInfoReq struct {
	GroupIDList []string `json:"groupIDList"`
	OperationID string   `json:"operationID"`
}

type getGroupsInfoResp struct {
	ErrCode int         `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	Data    []groupInfo `json:"data"`
}

type getGroupMemberListReq struct {
	GroupID     string `json:"groupID"`
	Filter      int32  `json:"filter"`
	NextSeq     int32  `json:"nextSeq"`
	OperationID string `json:"operationID"`
}

type getGroupAllMemberReq struct {
	GroupID     string `json:"groupID"`
	OperationID string `json:"operationID"`
}

type getGroupMembersInfoReq struct {
	GroupID     string   `json:"groupID"`
	MemberList  []string `json:"memberList"`
	OperationID string   `json:"operationID"`
}
type getGroupMembersInfoResp struct {
	commonResp
	Data []groupMemberFullInfo `json:"data"`
}

type inviteUserToGroupReq struct {
	GroupID     string   `json:"groupID"`
	UidList     []string `json:"uidList"`
	Reason      string   `json:"reason"`
	OperationID string   `json:"operationID"`
}

type inviteUserToGroupResp struct {
	commonResp
	Data []idResult `json:"data"`
}

type getJoinedGroupListReq struct {
	paramsCommonReq
}

type getJoinedGroupListResp struct {
	commonResp
	Data []groupInfo `json:"data"`
}

type kickGroupMemberApiReq struct {
	GroupID     string                `json:"groupID"`
	UidListInfo []groupMemberFullInfo `json:"uidListInfo"`
	Reason      string                `json:"reason"`
	OperationID string                `json:"operationID"`
}

type kickGroupMemberApiResp struct {
	commonResp
	Data []idResult `json:"data"`
}
type transferGroupReq struct {
	GroupID     string `json:"groupID"`
	Uid         string `json:"uid"`
	OperationID string `json:"operationID"`
}
type getGroupApplicationListReq struct {
	OperationID string `json:"operationID"`
}

type getGroupApplicationListResp struct {
	commonResp
	Data groupApplicationResult `json:"data"`
}

type accessOrRefuseGroupApplicationReq struct {
	OperationID      string `json:"operationID"`
	GroupId          string `json:"groupID"`
	FromUser         string `json:"fromUserID"`
	FromUserNickName string `json:"fromUserNickName"`
	FromUserFaceUrl  string `json:"fromUserFaceUrl"`
	ToUserNickname   string `json:"toUserNickName"`
	ToUserFaceUrl    string `json:"toUserFaceURL"`
	ToUser           string `json:"toUserID"`
	AddTime          int64  `json:"addTime"`
	RequestMsg       string `json:"requestMsg"`
	HandledMsg       string `json:"handledMsg"`
	Type             int32  `json:"type"`
	HandleStatus     int32  `json:"handleStatus"`
	HandleResult     int32  `json:"handleResult"`
}

type SoundElem struct {
	UUID      string `json:"uuid"`
	SoundPath string `json:"soundPath"`
	SourceURL string `json:"sourceUrl"`
	DataSize  int64  `json:"dataSize"`
	Duration  int64  `json:"duration"`
}

/*
OperationID          string   `protobuf:"bytes,2,opt,name=operationID,proto3" json:"operationID,omitempty"`
GroupID              string   `protobuf:"bytes,3,opt,name=groupID,proto3" json:"groupID,omitempty"`
Reason               string   `protobuf:"bytes,4,opt,name=reason,proto3" json:"reason,omitempty"`
UidList              []string `protobuf:"bytes,5,rep,name=uidList,proto3" json:"uidList,omitempty"`

*/

/*
type InviteUserToGroupReq struct {
	Op      string   `json:"op"`
	GroupID string   `json:"groupID"`
	Reason  string   `json:"reason"`
	UidList []string `json:"uidList"`
}*/

/*
type KickGroupMemberReq struct {
	Op      string   `json:"op"`
	GroupID string   `json:"groupID"`
	Reason  string   `json:"reason"`
	UidList []string `json:"uidList"`
}*/

//type TransferGroupOwnerReq struct {
//	GroupID     string
//	OldOwner    string
//	NewOwner    string
//	OperationID string
//}
type GroupApplicationResponseReq struct {
	OperationID      string
	OwnerID          string
	GroupID          string
	FromUserID       string
	FromUserNickName string
	FromUserFaceUrl  string
	ToUserID         string
	ToUserNickName   string
	ToUserFaceUrl    string
	AddTime          int64
	RequestMsg       string
	HandledMsg       string
	Type             int32
	HandleStatus     int32
	HandleResult     int32
}
type AgreeOrRejectGroupMember struct {
	GroupId  string `json:"groupID"`
	UserId   string `json:"userId"`
	Role     int    `json:"role"`
	JoinTime uint64 `json:"joinTime"`
	NickName string `json:"nickName"`
	FaceUrl  string `json:"faceUrl"`
	Reason   string `json:"reason"`
}

type GroupReqListInfo struct {
	ID               string `json:"id"`
	GroupID          string `json:"groupID"`
	FromUserID       string `json:"fromUserID"`
	ToUserID         string `json:"toUserID"`
	Flag             int32  `json:"flag"`
	RequestMsg       string `json:"reqMsg"`
	HandledMsg       string `json:"handledMsg"`
	AddTime          int64  `json:"createTime"`
	FromUserNickname string `json:"fromUserNickName"`
	ToUserNickname   string `json:"toUserNickName"`
	FromUserFaceUrl  string `json:"fromUserFaceURL"`
	ToUserFaceUrl    string `json:"toUserFaceURL"`
	HandledUser      string `json:"handledUser"`
	Type             int32  `json:"type"`
	HandleStatus     int32  `json:"handleStatus"`
	HandleResult     int32  `json:"handleResult"`
	//IsRead           int32  `json:"is_read"`
}

type NotificationContent struct {
	IsDisplay   int32  `json:"isDisplay"`
	DefaultTips string `json:"defaultTips"`
	Detail      string `json:"detail"`
}

type GroupApplicationInfo struct {
	Info         accessOrRefuseGroupApplicationReq `json:"info"`
	HandUserID   string                            `json:"handUserID"`
	HandUserName string                            `json:"handUserName"`
	HandUserIcon string                            `json:"handUserIcon"`
}

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}
type paramsSetReceiveMessageOpt struct {
	OperationID        string   `json:"operationID" binding:"required"`
	Option             int32    `json:"option" binding:"required"`
	ConversationIdList []string `json:"conversationIdList" binding:"required"`
}
type paramGetReceiveMessageOpt struct {
	ConversationIdList []string `json:"conversationIdList" binding:"required"`
	OperationID        string   `json:"operationID" binding:"required"`
}

//m:=make(map[string]interface)
//map["operationID"] = "3434"
//map["option"] = 1
//{
//	"operationID":"3434",
//	"option":1
//}
type paramGetAllConversationMessageOpt struct {
	OperationID string `json:"operationID" binding:"required"`
}

type optResult struct {
	ConversationId string `json:"conversationId" binding:"required"`
	Result         int32  `json:"result" binding:"required"`
}

type getReceiveMessageOptResp struct {
	ErrCode int32       `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	Data    []optResult `json:"data"`
}

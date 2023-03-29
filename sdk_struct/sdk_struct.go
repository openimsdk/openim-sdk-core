package sdk_struct

import "open_im_sdk/pkg/server_api_params"

////////////////////////// message/////////////////////////

type MessageReceipt struct {
	GroupID     string   `json:"groupID"`
	UserID      string   `json:"userID"`
	MsgIDList   []string `json:"msgIDList"`
	ReadTime    int64    `json:"readTime"`
	MsgFrom     int32    `json:"msgFrom"`
	ContentType int32    `json:"contentType"`
	SessionType int32    `json:"sessionType"`
}
type MessageRevoked struct {
	RevokerID                   string `json:"revokerID"`
	RevokerRole                 int32  `json:"revokerRole"`
	ClientMsgID                 string `json:"clientMsgID"`
	RevokerNickname             string `json:"revokerNickname"`
	RevokeTime                  int64  `json:"revokeTime"`
	SourceMessageSendTime       int64  `json:"sourceMessageSendTime"`
	SourceMessageSendID         string `json:"sourceMessageSendID"`
	SourceMessageSenderNickname string `json:"sourceMessageSenderNickname"`
	SessionType                 int32  `json:"sessionType"`
	Seq                         uint32 `json:"seq"`
	Ex                          string `json:"ex"`
}
type MessageReaction struct {
	ClientMsgID  string `json:"clientMsgID"`
	ReactionType int    `json:"reactionType"`
	Counter      int32  `json:"counter,omitempty"`
	UserID       string `json:"userID"`
	GroupID      string `json:"groupID"`
	SessionType  int32  `json:"sessionType"`
	Info         string `json:"info,omitempty"`
}
type ImageInfo struct {
	Width  int32  `json:"x"`
	Height int32  `json:"y"`
	Type   string `json:"type,omitempty"`
	Size   int64  `json:"size"`
}
type PictureBaseInfo struct {
	UUID   string `json:"uuid,omitempty"`
	Type   string `json:"type,omitempty"`
	Size   int64  `json:"size"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
	Url    string `json:"url,omitempty"`
}
type SoundBaseInfo struct {
	UUID      string `json:"uuid,omitempty"`
	SoundPath string `json:"soundPath,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
	DataSize  int64  `json:"dataSize"`
	Duration  int64  `json:"duration"`
}
type VideoBaseInfo struct {
	VideoPath      string `json:"videoPath,omitempty"`
	VideoUUID      string `json:"videoUUID,omitempty"`
	VideoURL       string `json:"videoUrl,omitempty"`
	VideoType      string `json:"videoType,omitempty"`
	VideoSize      int64  `json:"videoSize"`
	Duration       int64  `json:"duration"`
	SnapshotPath   string `json:"snapshotPath,omitempty"`
	SnapshotUUID   string `json:"snapshotUUID,omitempty"`
	SnapshotSize   int64  `json:"snapshotSize"`
	SnapshotURL    string `json:"snapshotUrl,omitempty"`
	SnapshotWidth  int32  `json:"snapshotWidth"`
	SnapshotHeight int32  `json:"snapshotHeight"`
}
type FileBaseInfo struct {
	FilePath  string `json:"filePath,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
	FileName  string `json:"fileName,omitempty"`
	FileSize  int64  `json:"fileSize"`
}

type MsgStruct struct {
	ClientMsgID          string                            `json:"clientMsgID,omitempty"`
	ServerMsgID          string                            `json:"serverMsgID,omitempty"`
	CreateTime           int64                             `json:"createTime"`
	SendTime             int64                             `json:"sendTime"`
	SessionType          int32                             `json:"sessionType"`
	SendID               string                            `json:"sendID,omitempty"`
	RecvID               string                            `json:"recvID,omitempty"`
	MsgFrom              int32                             `json:"msgFrom"`
	ContentType          int32                             `json:"contentType"`
	SenderPlatformID     int32                             `json:"platformID"`
	SenderNickname       string                            `json:"senderNickname,omitempty"`
	SenderFaceURL        string                            `json:"senderFaceUrl,omitempty"`
	GroupID              string                            `json:"groupID,omitempty"`
	Content              string                            `json:"content,omitempty"`
	Seq                  uint32                            `json:"seq"`
	IsRead               bool                              `json:"isRead"`
	Status               int32                             `json:"status"`
	IsReact              bool                              `json:"isReact,omitempty"`
	IsExternalExtensions bool                              `json:"isExternalExtensions,omitempty"`
	OfflinePush          server_api_params.OfflinePushInfo `json:"offlinePush,omitempty"`
	AttachedInfo         string                            `json:"attachedInfo,omitempty"`
	Ex                   string                            `json:"ex,omitempty"`
	PictureElem          struct {
		SourcePath      string          `json:"sourcePath,omitempty"`
		SourcePicture   PictureBaseInfo `json:"sourcePicture,omitempty"`
		BigPicture      PictureBaseInfo `json:"bigPicture,omitempty"`
		SnapshotPicture PictureBaseInfo `json:"snapshotPicture,omitempty"`
	} `json:"pictureElem,omitempty"`
	SoundElem struct {
		UUID      string `json:"uuid,omitempty"`
		SoundPath string `json:"soundPath,omitempty"`
		SourceURL string `json:"sourceUrl,omitempty"`
		DataSize  int64  `json:"dataSize"`
		Duration  int64  `json:"duration"`
	} `json:"soundElem,omitempty"`
	VideoElem struct {
		VideoPath      string `json:"videoPath,omitempty"`
		VideoUUID      string `json:"videoUUID,omitempty"`
		VideoURL       string `json:"videoUrl,omitempty"`
		VideoType      string `json:"videoType,omitempty"`
		VideoSize      int64  `json:"videoSize"`
		Duration       int64  `json:"duration"`
		SnapshotPath   string `json:"snapshotPath,omitempty"`
		SnapshotUUID   string `json:"snapshotUUID,omitempty"`
		SnapshotSize   int64  `json:"snapshotSize"`
		SnapshotURL    string `json:"snapshotUrl,omitempty"`
		SnapshotWidth  int32  `json:"snapshotWidth"`
		SnapshotHeight int32  `json:"snapshotHeight"`
	} `json:"videoElem,omitempty"`
	FileElem struct {
		FilePath  string `json:"filePath,omitempty"`
		UUID      string `json:"uuid,omitempty"`
		SourceURL string `json:"sourceUrl,omitempty"`
		FileName  string `json:"fileName,omitempty"`
		FileSize  int64  `json:"fileSize"`
	} `json:"fileElem,omitempty"`
	MergeElem struct {
		Title             string           `json:"title,omitempty"`
		AbstractList      []string         `json:"abstractList,omitempty"`
		MultiMessage      []*MsgStruct     `json:"multiMessage,omitempty"`
		MessageEntityList []*MessageEntity `json:"messageEntityList,omitempty"`
	} `json:"mergeElem,omitempty"`
	AtElem struct {
		Text         string     `json:"text,omitempty"`
		AtUserList   []string   `json:"atUserList,omitempty"`
		AtUsersInfo  []*AtInfo  `json:"atUsersInfo,omitempty"`
		QuoteMessage *MsgStruct `json:"quoteMessage,omitempty"`
		IsAtSelf     bool       `json:"isAtSelf"`
	} `json:"atElem,omitempty"`
	FaceElem struct {
		Index int    `json:"index"`
		Data  string `json:"data,omitempty"`
	} `json:"faceElem,omitempty"`
	LocationElem struct {
		Description string  `json:"description,omitempty"`
		Longitude   float64 `json:"longitude"`
		Latitude    float64 `json:"latitude"`
	} `json:"locationElem,omitempty"`
	CustomElem struct {
		Data        string `json:"data,omitempty"`
		Description string `json:"description,omitempty"`
		Extension   string `json:"extension,omitempty"`
	} `json:"customElem,omitempty"`
	QuoteElem struct {
		Text              string           `json:"text,omitempty"`
		QuoteMessage      *MsgStruct       `json:"quoteMessage,omitempty"`
		MessageEntityList []*MessageEntity `json:"messageEntityList,omitempty"`
	} `json:"quoteElem,omitempty"`
	NotificationElem struct {
		Detail      string `json:"detail,omitempty"`
		DefaultTips string `json:"defaultTips,omitempty"`
	} `json:"notificationElem,omitempty"`
	MessageEntityElem struct {
		Text              string           `json:"text,omitempty"`
		MessageEntityList []*MessageEntity `json:"messageEntityList,omitempty"`
	} `json:"messageEntityElem,omitempty"`
	AttachedInfoElem AttachedInfoElem `json:"attachedInfoElem,omitempty"`
}

type AtInfo struct {
	AtUserID      string `json:"atUserID,omitempty"`
	GroupNickname string `json:"groupNickname,omitempty"`
}
type AttachedInfoElem struct {
	GroupHasReadInfo          GroupHasReadInfo `json:"groupHasReadInfo,omitempty"`
	IsPrivateChat             bool             `json:"isPrivateChat"`
	BurnDuration              int32            `json:"burnDuration"`
	HasReadTime               int64            `json:"hasReadTime"`
	NotSenderNotificationPush bool             `json:"notSenderNotificationPush"`
	MessageEntityList         []*MessageEntity `json:"messageEntityList,omitempty"`
	IsEncryption              bool             `json:"isEncryption"`
	InEncryptStatus           bool             `json:"inEncryptStatus"`
	//MessageReactionElem       []*ReactionElem  `json:"messageReactionElem,omitempty"`
}

//type ReactionElem struct {
//	Counter          int32               `json:"counter,omitempty"`
//	Type             int                 `json:"type,omitempty"`
//	UserReactionList []*UserReactionElem `json:"userReactionList,omitempty"`
//	CanRepeat        bool                `json:"canRepeat,omitempty"`
//	Info             string              `json:"info,omitempty"`
//}
//type UserReactionElem struct {
//	UserID  string `json:"userID,omitempty"`
//	Counter int32  `json:"counter,omitempty"`
//	Info    string `json:"info,omitempty"`
//}

type MessageEntity struct {
	Type   string `json:"type,omitempty"`
	Offset int32  `json:"offset"`
	Length int32  `json:"length"`
	Url    string `json:"url,omitempty"`
	Info   string `json:"info,omitempty"`
}
type GroupHasReadInfo struct {
	HasReadUserIDList []string `json:"hasReadUserIDList,omitempty"`
	HasReadCount      int32    `json:"hasReadCount"`
	GroupMemberCount  int32    `json:"groupMemberCount"`
}
type NewMsgList []*MsgStruct

// Implement the sort.Interface interface to get the number of elements method
func (n NewMsgList) Len() int {
	return len(n)
}

//Implement the sort.Interface interface comparison element method
func (n NewMsgList) Less(i, j int) bool {
	return n[i].SendTime < n[j].SendTime
}

//Implement the sort.Interface interface exchange element method
func (n NewMsgList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

type IMConfig struct {
	Platform             int32  `json:"platform"`
	ApiAddr              string `json:"api_addr"`
	WsAddr               string `json:"ws_addr"`
	DataDir              string `json:"data_dir"`
	LogLevel             uint32 `json:"log_level"`
	ObjectStorage        string `json:"object_storage"` //"cos"(default)  "oss"
	EncryptionKey        string `json:"encryption_key"`
	IsCompression        bool   `json:"is_compression"`
	IsExternalExtensions bool   `json:"is_external_extensions"`
}

var SvrConf IMConfig

type CmdNewMsgComeToConversation struct {
	MsgList       []*server_api_params.MsgData
	OperationID   string
	SyncFlag      int
	MaxSeqOnSvr   uint32
	MaxSeqOnLocal uint32
	CurrentMaxSeq uint32
	PullMsgOrder  int
}

type CmdPushMsgToMsgSync struct {
	Msg         *server_api_params.MsgData
	OperationID string
}

type CmdMaxSeqToMsgSync struct {
	MaxSeqOnSvr            uint32
	OperationID            string
	MinSeqOnSvr            uint32
	GroupID2MinMaxSeqOnSvr map[string]*server_api_params.MaxAndMinSeq
}

type CmdJoinedSuperGroup struct {
	OperationID string
}

type OANotificationElem struct {
	NotificationName    string `mapstructure:"notificationName" validate:"required"`
	NotificationFaceURL string `mapstructure:"notificationFaceURL" validate:"required"`
	NotificationType    int32  `mapstructure:"notificationType" validate:"required"`
	Text                string `mapstructure:"text" validate:"required"`
	Url                 string `mapstructure:"url"`
	MixType             int32  `mapstructure:"mixType"`
	Image               struct {
		SourceUrl   string `mapstructure:"sourceURL"`
		SnapshotUrl string `mapstructure:"snapshotURL"`
	} `mapstructure:"image"`
	Video struct {
		SourceUrl   string `mapstructure:"sourceURL"`
		SnapshotUrl string `mapstructure:"snapshotURL"`
		Duration    int64  `mapstructure:"duration"`
	} `mapstructure:"video"`
	File struct {
		SourceUrl string `mapstructure:"sourceURL"`
		FileName  string `mapstructure:"fileName"`
		FileSize  int64  `mapstructure:"fileSize"`
	} `mapstructure:"file"`
	Ex string `mapstructure:"ex"`
}
type MsgDeleteNotificationElem struct {
	GroupID     string   `json:"groupID"`
	IsAllDelete bool     `json:"isAllDelete"`
	SeqList     []string `json:"seqList"`
}

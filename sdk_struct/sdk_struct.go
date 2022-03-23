package sdk_struct

import "open_im_sdk/pkg/server_api_params"

////////////////////////// message/////////////////////////

type MessageReceipt struct {
	GroupID     string   `json:"groupID"`
	UserID      string   `json:"userID"`
	MsgIdList   []string `json:"msgIDList"`
	ReadTime    int64    `json:"readTime"`
	MsgFrom     int32    `json:"msgFrom"`
	ContentType int32    `json:"contentType"`
	SessionType int32    `json:"sessionType"`
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
	SenderNickname   string                            `json:"senderNickname"`
	SenderFaceURL    string                            `json:"senderFaceUrl"`
	GroupID          string                            `json:"groupID"`
	Content          string                            `json:"content"`
	Seq              uint32                            `json:"seq"`
	IsRead           bool                              `json:"isRead"`
	Status           int32                             `json:"status"`
	OfflinePush      server_api_params.OfflinePushInfo `json:"offlinePush"`
	AttachedInfo     string                            `json:"attachedInfo"`
	Ex               string                            `json:"ex"`
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
	FaceElem struct {
		Index int    `json:"index"`
		Data  string `json:"data"`
	} `json:"faceElem"`
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
	NotificationElem struct {
		Detail      string `json:"detail"`
		DefaultTips string `json:"defaultTips"`
	} `json:"notificationElem"`
	AttachedInfoElem AttachedInfoElem `json:"attachedInfoElem"`
}
type AttachedInfoElem struct {
	GroupHasReadInfo GroupHasReadInfo `json:"groupHasReadInfo"`
}
type GroupHasReadInfo struct {
	HasReadUserIDList []string `json:"hasReadUserIDList"`
	HasReadCount      int32    `json:"hasReadCount"`
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
	Platform      int32  `json:"platform"`
	ApiAddr       string `json:"api_addr"`
	WsAddr        string `json:"ws_addr"`
	DataDir       string `json:"data_dir"`
	LogLevel      uint32 `json:"log_level"`
	ObjectStorage string `json:"object_storage"` //"cos"(default)  "oss"
}

var SvrConf IMConfig

type CmdNewMsgComeToConversation struct {
	MsgList     []*server_api_params.MsgData
	OperationID string
}

type CmdPushMsgToMsgSync struct {
	Msg         *server_api_params.MsgData
	OperationID string
}

type CmdMaxSeqToMsgSync struct {
	MaxSeqOnSvr uint32
	OperationID string
}

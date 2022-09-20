package temp_struct

type LocalChatLog struct {
	ServerMsgID      string ` json:"serverMsgID,omitempty"`
	SendID           string ` json:"sendID,omitempty"`
	RecvID           string ` json:"recvID,omitempty"`
	SenderPlatformID int32  ` json:"senderPlatformID,omitempty"`
	SenderNickname   string ` json:"senderNickname,omitempty"`
	SenderFaceURL    string ` json:"senderFaceURL,omitempty"`
	SessionType      int32  ` json:"sessionType,omitempty"`
	MsgFrom          int32  ` json:"msgFrom,omitempty"`
	ContentType      int32  ` json:"contentType,omitempty"`
	Content          string ` json:"content,omitempty"`
	IsRead           bool   ` json:"isRead,omitempty"`
	Status           int32  ` json:"status,omitempty"`
	Seq              uint32 ` json:"seq,omitempty"`
	SendTime         int64  ` json:"sendTime,omitempty"`
	CreateTime       int64  ` json:"createTime,omitempty"`
	AttachedInfo     string ` json:"attachedInfo,omitempty"`
	Ex               string ` json:"ex,omitempty"`
}

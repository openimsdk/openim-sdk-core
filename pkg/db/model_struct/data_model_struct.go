package model_struct

import "time"

//
//message FriendInfo{
//string OwnerUserID = 1;
//string Remark = 2;
//int64 CreateTime = 3;
//UserInfo FriendUser = 4;
//int32 AddSource = 5;
//string OperatorUserID = 6;
//string Ex = 7;
//}
//open_im_sdk.FriendInfo(FriendUser) != imdb.Friend(FriendUserID)
//	table = ` CREATE TABLE IF NOT EXISTS friends(
//     owner_user_id CHAR (64) NOT NULL,
//     friend_user_id CHAR (64) NOT NULL ,
//     name varchar(64) DEFAULT NULL ,
//	 face_url varchar(100) DEFAULT NULL ,
//     remark varchar(255) DEFAULT NULL,
//     gender int DEFAULT NULL ,
//   	 phone_number varchar(32) DEFAULT NULL ,
//	 birth INTEGER DEFAULT NULL ,
//	 email varchar(64) DEFAULT NULL ,
//	 create_time INTEGER DEFAULT NULL ,
//	 add_source int DEFAULT NULL ,
//	 operator_user_id CHAR(64) DEFAULT NULL,
//  	 ex varchar(1024) DEFAULT NULL,
//  	 PRIMARY KEY (owner_user_id,friend_user_id)
// 	)`

type LocalFriend struct {
	OwnerUserID    string `gorm:"column:owner_user_id;primary_key;type:varchar(64)" json:"ownerUserID"`
	FriendUserID   string `gorm:"column:friend_user_id;primary_key;type:varchar(64)" json:"userID"`
	Remark         string `gorm:"column:remark;type:varchar(255)" json:"remark"`
	CreateTime     uint32 `gorm:"column:create_time" json:"createTime"`
	AddSource      int32  `gorm:"column:add_source" json:"addSource"`
	OperatorUserID string `gorm:"column:operator_user_id;type:varchar(64)" json:"operatorUserID"`
	Nickname       string `gorm:"column:name;type:varchar;type:varchar(255)" json:"nickname"`
	FaceURL        string `gorm:"column:face_url;type:varchar;type:varchar(255)" json:"faceURL"`
	Gender         int32  `gorm:"column:gender" json:"gender"`
	PhoneNumber    string `gorm:"column:phone_number;type:varchar(32)" json:"phoneNumber"`
	Birth          uint32 `gorm:"column:birth" json:"birth"`
	Email          string `gorm:"column:email;type:varchar(64)" json:"email"`
	Ex             string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo   string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
}

//message FriendRequest{
//string  FromUserID = 1;
//string ToUserID = 2;
//int32 HandleResult = 3;
//string ReqMsg = 4;
//int64 CreateTime = 5;
//string HandlerUserID = 6;
//string HandleMsg = 7;
//int64 HandleTime = 8;
//string Ex = 9;
//}
//open_im_sdk.FriendRequest == imdb.FriendRequest
type LocalFriendRequest struct {
	FromUserID   string `gorm:"column:from_user_id;primary_key;type:varchar(64)" json:"fromUserID"`
	FromNickname string `gorm:"column:from_nickname;type:varchar;type:varchar(255)" json:"fromNickname"`
	FromFaceURL  string `gorm:"column:from_face_url;type:varchar;type:varchar(255)" json:"fromFaceURL"`
	FromGender   int32  `gorm:"column:from_gender" json:"fromGender"`

	ToUserID   string `gorm:"column:to_user_id;primary_key;type:varchar(64)" json:"toUserID"`
	ToNickname string `gorm:"column:to_nickname;type:varchar;type:varchar(255)" json:"toNickname"`
	ToFaceURL  string `gorm:"column:to_face_url;type:varchar;type:varchar(255)" json:"toFaceURL"`
	ToGender   int32  `gorm:"column:to_gender" json:"toGender"`

	HandleResult  int32  `gorm:"column:handle_result" json:"handleResult"`
	ReqMsg        string `gorm:"column:req_msg;type:varchar(255)" json:"reqMsg"`
	CreateTime    uint32 `gorm:"column:create_time" json:"createTime"`
	HandlerUserID string `gorm:"column:handler_user_id;type:varchar(64)" json:"handlerUserID"`
	HandleMsg     string `gorm:"column:handle_msg;type:varchar(255)" json:"handleMsg"`
	HandleTime    uint32 `gorm:"column:handle_time" json:"handleTime"`
	Ex            string `gorm:"column:ex;type:varchar(1024)" json:"ex"`

	AttachedInfo string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
}

//message GroupInfo{
//  string GroupID = 1;
//  string GroupName = 2;
//  string Notification = 3;
//  string Introduction = 4;
//  string FaceUrl = 5;
//  string OwnerUserID = 6;
//  uint32 MemberCount = 8;
//  int64 CreateTime = 7;
//  string Ex = 9;
//  int32 Status = 10;
//  string CreatorUserID = 11;
//  int32 GroupType = 12;
//}
//  open_im_sdk.GroupInfo (OwnerUserID ,  MemberCount )> imdb.Group
//    	group_id char(64) NOT NULL,
//		name varchar(64) DEFAULT NULL ,
//    	introduction varchar(255) DEFAULT NULL,
//    	notification varchar(255) DEFAULT NULL,
//    	face_url varchar(100) DEFAULT NULL,
//    	group_type int DEFAULT NULL,
//    	status int DEFAULT NULL,
//    	creator_user_id char(64) DEFAULT NULL,
//    	create_time INTEGER DEFAULT NULL,
//    	ex varchar(1024) DEFAULT NULL,
//    	PRIMARY KEY (group_id)
//	)`

type LocalGroup struct {
	//`json:"operationID" binding:"required"`
	//`protobuf:"bytes,1,opt,name=GroupID" json:"GroupID,omitempty"` `json:"operationID" binding:"required"`
	GroupID                string `gorm:"column:group_id;primary_key;type:varchar(64)" json:"groupID" binding:"required"`
	GroupName              string `gorm:"column:name;size:255" json:"groupName"`
	Notification           string `gorm:"column:notification;type:varchar(255)" json:"notification"`
	Introduction           string `gorm:"column:introduction;type:varchar(255)" json:"introduction"`
	FaceURL                string `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
	CreateTime             uint32 `gorm:"column:create_time" json:"createTime"`
	Status                 int32  `gorm:"column:status" json:"status"`
	CreatorUserID          string `gorm:"column:creator_user_id;type:varchar(64)" json:"creatorUserID"`
	GroupType              int32  `gorm:"column:group_type" json:"groupType"`
	OwnerUserID            string `gorm:"column:owner_user_id;type:varchar(64)" json:"ownerUserID"`
	MemberCount            int32  `gorm:"column:member_count" json:"memberCount"`
	Ex                     string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo           string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	NeedVerification       int32  `gorm:"column:need_verification"  json:"needVerification"`
	LookMemberInfo         int32  `gorm:"column:look_member_info" json:"lookMemberInfo"`
	ApplyMemberFriend      int32  `gorm:"column:apply_member_friend" json:"applyMemberFriend"`
	NotificationUpdateTime uint32 `gorm:"column:notification_update_time" json:"notificationUpdateTime"`
	NotificationUserID     string `gorm:"column:notification_user_id;size:64" json:"notificationUserID"`
}

//message GroupMemberFullInfo {
//string GroupID = 1 ;
//string UserID = 2 ;
//int32 roleLevel = 3;
//int64 JoinTime = 4;
//string NickName = 5;
//string FaceUrl = 6;
//int32 JoinSource = 8;
//string OperatorUserID = 9;
//string Ex = 10;
//int32 AppMangerLevel = 7; //if >0
//}  open_im_sdk.GroupMemberFullInfo(AppMangerLevel) > imdb.GroupMember
//  group_id char(64) NOT NULL,
//   user_id char(64) NOT NULL,
//   nickname varchar(64) DEFAULT NULL,
//   user_group_face_url varchar(64) DEFAULT NULL,
//   role_level int DEFAULT NULL,
//   join_time INTEGER DEFAULT NULL,
//   join_source int DEFAULT NULL,
//   operator_user_id char(64) NOT NULL,

type LocalGroupMember struct {
	GroupID        string `gorm:"column:group_id;primary_key;type:varchar(64)" json:"groupID"`
	UserID         string `gorm:"column:user_id;primary_key;type:varchar(64)" json:"userID"`
	Nickname       string `gorm:"column:nickname;type:varchar(255)" json:"nickname"`
	FaceURL        string `gorm:"column:user_group_face_url;type:varchar(255)" json:"faceURL"`
	RoleLevel      int32  `gorm:"column:role_level" json:"roleLevel"`
	JoinTime       uint32 `gorm:"column:join_time;index:index_join_time;" json:"joinTime"`
	JoinSource     int32  `gorm:"column:join_source" json:"joinSource"`
	InviterUserID  string `gorm:"column:inviter_user_id;size:64"  json:"inviterUserID"`
	MuteEndTime    uint32 `gorm:"column:mute_end_time;default:0" json:"muteEndTime"`
	OperatorUserID string `gorm:"column:operator_user_id;type:varchar(64)" json:"operatorUserID"`
	Ex             string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo   string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
}

//message GroupRequest{
//string UserID = 1;
//string GroupID = 2;
//string HandleResult = 3;
//string ReqMsg = 4;
//string  HandleMsg = 5;
//int64 ReqTime = 6;
//string HandleUserID = 7;
//int64 HandleTime = 8;
//string Ex = 9;
//}open_im_sdk.GroupRequest == imdb.GroupRequest
type LocalGroupRequest struct {
	GroupID       string `gorm:"column:group_id;primary_key;type:varchar(64)" json:"groupID"`
	GroupName     string `gorm:"column:group_name;size:255" json:"groupName"`
	Notification  string `gorm:"column:notification;type:varchar(255)" json:"notification"`
	Introduction  string `gorm:"column:introduction;type:varchar(255)" json:"introduction"`
	GroupFaceURL  string `gorm:"column:face_url;type:varchar(255)" json:"groupFaceURL"`
	CreateTime    uint32 `gorm:"column:create_time" json:"createTime"`
	Status        int32  `gorm:"column:status" json:"status"`
	CreatorUserID string `gorm:"column:creator_user_id;type:varchar(64)" json:"creatorUserID"`
	GroupType     int32  `gorm:"column:group_type" json:"groupType"`
	OwnerUserID   string `gorm:"column:owner_user_id;type:varchar(64)" json:"ownerUserID"`
	MemberCount   int32  `gorm:"column:member_count" json:"memberCount"`

	UserID      string `gorm:"column:user_id;primary_key;type:varchar(64)" json:"userID"`
	Nickname    string `gorm:"column:nickname;type:varchar(255)" json:"nickname"`
	UserFaceURL string `gorm:"column:user_face_url;type:varchar(255)" json:"userFaceURL"`
	Gender      int32  `gorm:"column:gender" json:"gender"`

	HandleResult  int32  `gorm:"column:handle_result" json:"handleResult"`
	ReqMsg        string `gorm:"column:req_msg;type:varchar(255)" json:"reqMsg"`
	HandledMsg    string `gorm:"column:handle_msg;type:varchar(255)" json:"handledMsg"`
	ReqTime       uint32 `gorm:"column:req_time" json:"reqTime"`
	HandleUserID  string `gorm:"column:handle_user_id;type:varchar(64)" json:"handleUserID"`
	HandledTime   uint32 `gorm:"column:handle_time" json:"handledTime"`
	Ex            string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo  string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	JoinSource    int32  `gorm:"column:join_source" json:"joinSource"`
	InviterUserID string `gorm:"column:inviter_user_id;size:64"  json:"inviterUserID"`
}

//string UserID = 1;
//string Nickname = 2;
//string FaceUrl = 3;
//int32 Gender = 4;
//string PhoneNumber = 5;
//string Birth = 6;
//string Email = 7;
//string Ex = 8;
//int64 CreateTime = 9;
//int32 AppMangerLevel = 10;
//open_im_sdk.User == imdb.User
type LocalUser struct {
	UserID           string    `gorm:"column:user_id;primary_key;type:varchar(64)" json:"userID"`
	Nickname         string    `gorm:"column:name;type:varchar(255)" json:"nickname"`
	FaceURL          string    `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
	Gender           int32     `gorm:"column:gender" json:"gender"`
	PhoneNumber      string    `gorm:"column:phone_number;type:varchar(32)" json:"phoneNumber"`
	Birth            uint32    `gorm:"column:birth" json:"birth"`
	Email            string    `gorm:"column:email;type:varchar(64)" json:"email"`
	CreateTime       uint32    `gorm:"column:create_time" json:"createTime"`
	AppMangerLevel   int32     `gorm:"column:app_manger_level" json:"-"`
	Ex               string    `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo     string    `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	GlobalRecvMsgOpt int32     `gorm:"column:global_recv_msg_opt" json:"globalRecvMsgOpt"`
	BirthTime        time.Time `gorm:"column:birth_time" json:"birthTime"`
}

//message BlackInfo{
//string OwnerUserID = 1;
//int64 CreateTime = 2;
//PublicUserInfo BlackUserInfo = 4;
//int32 AddSource = 5;
//string OperatorUserID = 6;
//string Ex = 7;
//}
// open_im_sdk.BlackInfo(BlackUserInfo) != imdb.Black (BlockUserID)
type LocalBlack struct {
	OwnerUserID    string `gorm:"column:owner_user_id;primary_key;type:varchar(64)" json:"ownerUserID"`
	BlockUserID    string `gorm:"column:block_user_id;primary_key;type:varchar(64)" json:"userID"`
	Nickname       string `gorm:"column:nickname;type:varchar(255)" json:"nickname"`
	FaceURL        string `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
	Gender         int32  `gorm:"column:gender" json:"gender"`
	CreateTime     uint32 `gorm:"column:create_time" json:"createTime"`
	AddSource      int32  `gorm:"column:add_source" json:"addSource"`
	OperatorUserID string `gorm:"column:operator_user_id;type:varchar(64)" json:"operatorUserID"`
	Ex             string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo   string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
}

type LocalSeqData struct {
	UserID string `gorm:"column:user_id;primary_key;type:varchar(64)"`
	Seq    uint32 `gorm:"column:seq"`
}

type LocalSeq struct {
	ID     string `gorm:"column:id;primary_key;type:varchar(64)"`
	MinSeq uint32 `gorm:"column:min_seq"`
}

//`create table if not exists  chat_log (
//      client_msg_id char(64)   NOT NULL,
//      server_msg_id char(64)   DEFAULT NULL,
//	  send_id char(64)   NOT NULL ,
//	  is_read int NOT NULL ,
//	  seq INTEGER DEFAULT NULL ,
//	  status int NOT NULL ,
//	  session_type int NOT NULL ,
//	  recv_id char(64)   NOT NULL ,
//	  content_type int NOT NULL ,
//      sender_face_url varchar(100) DEFAULT NULL,
//      sender_nick_name varchar(64) DEFAULT NULL,
//	  msg_from int NOT NULL ,
//	  content varchar(1000)   NOT NULL ,
//	  sender_platform_id int NOT NULL ,
//	  send_time INTEGER DEFAULT NULL ,
//	  create_time INTEGER  DEFAULT NULL,
//      ex varchar(1024) DEFAULT NULL,
//	  PRIMARY KEY (client_msg_id)
//	)`
type LocalChatLog struct {
	ClientMsgID          string `gorm:"column:client_msg_id;primary_key;type:char(64)" json:"clientMsgID"`
	ServerMsgID          string `gorm:"column:server_msg_id;type:char(64)" json:"serverMsgID"`
	SendID               string `gorm:"column:send_id;type:char(64)" json:"sendID"`
	RecvID               string `gorm:"column:recv_id;index:index_recv_id;type:char(64)" json:"recvID"`
	SenderPlatformID     int32  `gorm:"column:sender_platform_id" json:"senderPlatformID"`
	SenderNickname       string `gorm:"column:sender_nick_name;type:varchar(255)" json:"senderNickname"`
	SenderFaceURL        string `gorm:"column:sender_face_url;type:varchar(255)" json:"senderFaceURL"`
	SessionType          int32  `gorm:"column:session_type" json:"sessionType"`
	MsgFrom              int32  `gorm:"column:msg_from" json:"msgFrom"`
	ContentType          int32  `gorm:"column:content_type;index:content_type_alone" json:"contentType"`
	Content              string `gorm:"column:content;type:varchar(1000)" json:"content"`
	IsRead               bool   `gorm:"column:is_read" json:"isRead"`
	Status               int32  `gorm:"column:status" json:"status"`
	Seq                  uint32 `gorm:"column:seq;index:index_seq;default:0" json:"seq"`
	SendTime             int64  `gorm:"column:send_time;index:index_send_time;" json:"sendTime"`
	CreateTime           int64  `gorm:"column:create_time" json:"createTime"`
	AttachedInfo         string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex                   string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	IsReact              bool   `gorm:"column:is_react" json:"isReact"`
	IsExternalExtensions bool   `gorm:"column:is_external_extensions" json:"isExternalExtensions"`
	MsgFirstModifyTime   int64  `gorm:"column:msg_first_modify_time" json:"msgFirstModifyTime"`
}

type LocalErrChatLog struct {
	Seq              uint32 `gorm:"column:seq;primary_key" json:"seq"`
	ClientMsgID      string `gorm:"column:client_msg_id;type:char(64)" json:"clientMsgID"`
	ServerMsgID      string `gorm:"column:server_msg_id;type:char(64)" json:"serverMsgID"`
	SendID           string `gorm:"column:send_id;type:char(64)" json:"sendID"`
	RecvID           string `gorm:"column:recv_id;type:char(64)" json:"recvID"`
	SenderPlatformID int32  `gorm:"column:sender_platform_id" json:"senderPlatformID"`
	SenderNickname   string `gorm:"column:sender_nick_name;type:varchar(255)" json:"senderNickname"`
	SenderFaceURL    string `gorm:"column:sender_face_url;type:varchar(255)" json:"senderFaceURL"`
	SessionType      int32  `gorm:"column:session_type" json:"sessionType"`
	MsgFrom          int32  `gorm:"column:msg_from" json:"msgFrom"`
	ContentType      int32  `gorm:"column:content_type" json:"contentType"`
	Content          string `gorm:"column:content;type:varchar(1000)" json:"content"`
	IsRead           bool   `gorm:"column:is_read" json:"isRead"`
	Status           int32  `gorm:"column:status" json:"status"`
	SendTime         int64  `gorm:"column:send_time" json:"sendTime"`
	CreateTime       int64  `gorm:"column:create_time" json:"createTime"`
	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}
type TempCacheLocalChatLog struct {
	ClientMsgID      string `gorm:"column:client_msg_id;primary_key;type:char(64)" json:"clientMsgID"`
	ServerMsgID      string `gorm:"column:server_msg_id;type:char(64)" json:"serverMsgID"`
	SendID           string `gorm:"column:send_id;type:char(64)" json:"sendID"`
	RecvID           string `gorm:"column:recv_id;type:char(64)" json:"recvID"`
	SenderPlatformID int32  `gorm:"column:sender_platform_id" json:"senderPlatformID"`
	SenderNickname   string `gorm:"column:sender_nick_name;type:varchar(255)" json:"senderNickname"`
	SenderFaceURL    string `gorm:"column:sender_face_url;type:varchar(255)" json:"senderFaceURL"`
	SessionType      int32  `gorm:"column:session_type" json:"sessionType"`
	MsgFrom          int32  `gorm:"column:msg_from" json:"msgFrom"`
	ContentType      int32  `gorm:"column:content_type" json:"contentType"`
	Content          string `gorm:"column:content;type:varchar(1000)" json:"content"`
	IsRead           bool   `gorm:"column:is_read" json:"isRead"`
	Status           int32  `gorm:"column:status" json:"status"`
	Seq              uint32 `gorm:"column:seq;default:0" json:"seq"`
	SendTime         int64  `gorm:"column:send_time;" json:"sendTime"`
	CreateTime       int64  `gorm:"column:create_time" json:"createTime"`
	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

//`create table if not exists  conversation (
//  conversation_id char(128) NOT NULL,
//  conversation_type int(11) NOT NULL,
//  user_id varchar(128)  DEFAULT NULL,
//  group_id varchar(128)  DEFAULT NULL,
//  show_name varchar(128)  NOT NULL,
//  face_url varchar(128)  NOT NULL,
//  recv_msg_opt int(11) NOT NULL ,
//  unread_count int(11) NOT NULL ,
//  latest_msg varchar(255)  NOT NULL ,
//  latest_msg_send_time INTEGER(255)  NOT NULL ,
//  draft_text varchar(255)  DEFAULT NULL ,
//  draft_timestamp INTEGER(255)  DEFAULT NULL ,
//  is_pinned int(10) NOT NULL ,
//  PRIMARY KEY (conversation_id)
//)
type LocalConversation struct {
	ConversationID        string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	ConversationType      int32  `gorm:"column:conversation_type" json:"conversationType"`
	UserID                string `gorm:"column:user_id;type:char(64)" json:"userID"`
	GroupID               string `gorm:"column:group_id;type:char(128)" json:"groupID"`
	ShowName              string `gorm:"column:show_name;type:varchar(255)" json:"showName"`
	FaceURL               string `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
	RecvMsgOpt            int32  `gorm:"column:recv_msg_opt" json:"recvMsgOpt"`
	UnreadCount           int32  `gorm:"column:unread_count" json:"unreadCount"`
	GroupAtType           int32  `gorm:"column:group_at_type" json:"groupAtType"`
	LatestMsg             string `gorm:"column:latest_msg;type:varchar(1000)" json:"latestMsg"`
	LatestMsgSendTime     int64  `gorm:"column:latest_msg_send_time;index:index_latest_msg_send_time" json:"latestMsgSendTime"`
	DraftText             string `gorm:"column:draft_text" json:"draftText"`
	DraftTextTime         int64  `gorm:"column:draft_text_time" json:"draftTextTime"`
	IsPinned              bool   `gorm:"column:is_pinned" json:"isPinned"`
	IsPrivateChat         bool   `gorm:"column:is_private_chat" json:"isPrivateChat"`
	BurnDuration          int32  `gorm:"column:burn_duration;default:30" json:"burnDuration"`
	IsNotInGroup          bool   `gorm:"column:is_not_in_group" json:"isNotInGroup"`
	UpdateUnreadCountTime int64  `gorm:"column:update_unread_count_time" json:"updateUnreadCountTime"`
	AttachedInfo          string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex                    string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}
type LocalConversationUnreadMessage struct {
	ConversationID string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	ClientMsgID    string `gorm:"column:client_msg_id;primary_key;type:char(64)" json:"clientMsgID"`
	SendTime       int64  `gorm:"column:send_time" json:"sendTime"`
	Ex             string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

//message GroupRequest{
//string UserID = 1;
//string GroupID = 2;
//string HandleResult = 3;
//string ReqMsg = 4;
//string  HandleMsg = 5;
//int64 ReqTime = 6;
//string HandleUserID = 7;
//int64 HandleTime = 8;
//string Ex = 9;
//}open_im_sdk.GroupRequest == imdb.GroupRequest
type LocalAdminGroupRequest struct {
	LocalGroupRequest
}

type LocalDepartment struct {
	DepartmentID     string `gorm:"column:department_id;primary_key;size:64" json:"departmentID"`
	FaceURL          string `gorm:"column:face_url;size:255" json:"faceURL"`
	Name             string `gorm:"column:name;size:256" json:"name" binding:"required"`
	ParentID         string `gorm:"column:parent_id;size:64" json:"parentID" binding:"required"` // "0" or Real parent id
	Order            int32  `gorm:"column:order_department" json:"order" `                       // 1, 2, ...
	DepartmentType   int32  `gorm:"column:department_type" json:"departmentType"`                //1, 2...
	CreateTime       uint32 `gorm:"column:create_time" json:"createTime"`
	SubDepartmentNum uint32 `gorm:"column:sub_department_num" json:"subDepartmentNum"`
	MemberNum        uint32 `gorm:"column:member_num" json:"memberNum"`
	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
}

type LocalDepartmentMember struct {
	UserID      string `gorm:"column:user_id;primary_key;size:64" json:"userID"`
	Nickname    string `gorm:"column:nickname;size:256" json:"nickname"`
	EnglishName string `gorm:"column:english_name;size:256" json:"englishName"`
	FaceURL     string `gorm:"column:face_url;size:256" json:"faceURL"`
	Gender      int32  `gorm:"column:gender" json:"gender"` //1 ,2
	Mobile      string `gorm:"column:mobile;size:32" json:"mobile"`
	Telephone   string `gorm:"column:telephone;size:32" json:"telephone"`
	Birth       uint32 `gorm:"column:birth" json:"birth"`
	Email       string `gorm:"column:email;size:64" json:"email"`

	DepartmentID string `gorm:"column:department_id;primary_key;size:64" json:"departmentID"`
	Order        int32  `gorm:"column:order_member" json:"order"` //1,2
	Position     string `gorm:"column:position;size:256" json:"position"`
	Leader       int32  `gorm:"column:leader" json:"leader"` //-1, 1
	Status       int32  `gorm:"column:status" json:"status"` //-1, 1
	CreateTime   uint32 `gorm:"column:create_time" json:"createTime"`
	Ex           string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
}

type SearchDepartmentMemberResult struct {
	LocalDepartmentMember
	DepartmentName string `gorm:"column:name;size:256" json:"departmentName"`
}
type LocalChatLogReactionExtensions struct {
	ClientMsgID             string `gorm:"column:client_msg_id;primary_key;type:char(64)" json:"clientMsgID"`
	LocalReactionExtensions []byte `gorm:"column:local_reaction_extensions" json:"localReactionExtensions"`
}

package open_im_sdk

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

type Friend struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;type:varchar(64)"`
	FriendUserID   string    `gorm:"column:friend_user_id;primary_key;type:varchar;size:64"`
	Remark         string    `gorm:"column:remark;type:varchar;size:255"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;type:varchar;size:64"`
	Nickname       string    `gorm:"column:name;type:varchar;size:255"`
	FaceUrl        string    `gorm:"column:face_url;type:varchar;size:255"`
	Gender         int32     `gorm:"column:gender"`
	PhoneNumber    string    `gorm:"column:phone_number;type:varchar;size:32"`
	Birth          time.Time `gorm:"column:birth"`
	Email          string    `gorm:"column:email;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
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
type FriendRequest struct {
	FromUserID    string    `gorm:"column:from_user_id;primary_key;size:64"`
	ToUserID      string    `gorm:"column:to_user_id;primary_key;size:64"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg;size:255"`
	CreateTime    time.Time `gorm:"column:create_time"`
	HandlerUserID string    `gorm:"column:handler_user_id;size:64"`
	HandleMsg     string    `gorm:"column:handle_msg;size:255"`
	HandleTime    time.Time `gorm:"column:handle_time"`
	Ex            string    `gorm:"column:ex;size:1024"`
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

type Group struct {
	//`json:"operationID" binding:"required"`
	//`protobuf:"bytes,1,opt,name=GroupID" json:"GroupID,omitempty"` `json:"operationID" binding:"required"`
	GroupID       string    `gorm:"column:group_id;primary_key;size:64" json:"groupID" binding:"required"`
	GroupName     string    `gorm:"column:name;size:255" json:"groupName"`
	Notification  string    `gorm:"column:notification;size:255" json:"notification"`
	Introduction  string    `gorm:"column:introduction;size:255" json:"introduction"`
	FaceUrl       string    `gorm:"column:face_url;size:255" json:"faceUrl"`
	CreateTime    time.Time `gorm:"column:create_time"`
	Status        int32     `gorm:"column:status"`
	CreatorUserID string    `gorm:"column:creator_user_id;size:64"`
	GroupType     int32     `gorm:"column:group_type"`
	OwnerUserID   string    `gorm:"column:owner_user_id;size:64"`
	MemberCount   int32     `gorm:"column:member_count"`
	Ex            string    `gorm:"column:ex" json:"ex;size:1024"`
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

type GroupMember struct {
	GroupID        string    `gorm:"column:group_id;primary_key;size:64"`
	UserID         string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname       string    `gorm:"column:nickname;size:255"`
	FaceUrl        string    `gorm:"column:user_group_face_url;size:255"`
	RoleLevel      int32     `gorm:"column:role_level"`
	JoinTime       time.Time `gorm:"column:join_time"`
	JoinSource     int32     `gorm:"column:join_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
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
type GroupRequest struct {
	GroupID      string    `gorm:"column:group_id;primary_key;size:64"`
	UserID       string    `gorm:"column:user_id;primary_key;size:64"`
	HandleResult int32     `gorm:"column:handle_result"`
	ReqMsg       string    `gorm:"column:req_msg;size:255"`
	HandledMsg   string    `gorm:"column:handle_msg;size:255"`
	ReqTime      time.Time `gorm:"column:req_time"`
	HandleUserID string    `gorm:"column:handle_user_id;size:64"`
	HandledTime  time.Time `gorm:"column:handle_time"`
	Ex           string    `gorm:"column:ex;size:1024"`
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
type User struct {
	UserID         string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname       string    `gorm:"column:name;size:255"`
	FaceUrl        string    `gorm:"column:face_url;size:255"`
	Gender         int32     `gorm:"column:gender"`
	PhoneNumber    string    `gorm:"column:phone_number;size:32"`
	Birth          time.Time `gorm:"column:birth"`
	Email          string    `gorm:"column:email;size:64"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AppMangerLevel int32     `gorm:"column:app_manger_level"`
	Ex             string    `gorm:"column:ex;size:1024"`
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
type Black struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;size:64"`
	BlockUserID    string    `gorm:"column:block_user_id;primary_key;size:64"`
	Nickname       string    `gorm:"column:nick_name;size:255"`
	FaceUrl        string    `gorm:"column:face_url;size:255"`
	Gender         int32     `gorm:"column:gender"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
}
type LocalData struct {
	UserID string `gorm:"column:user_id;primary_key;size:64"`
	Seq    int32  `gorm:"column:seq;default: '1'"`
}

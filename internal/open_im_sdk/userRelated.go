package open_im_sdk

import (
	"open_im_sdk/internal/controller/conversation_msg"
	"open_im_sdk/internal/controller/friend"
	"open_im_sdk/internal/controller/group"
	"open_im_sdk/internal/controller/interaction"
	"open_im_sdk/internal/controller/ws"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"sync"
)

type UserRelated struct {
	ConversationCh chan conversation_msg.Cmd2Value //cmd：

	token       string
	loginUserID string

	utils.IMManager
	friend.Friend
	conversation_msg.ConversationListener
	group.Group

	imDB *db.DataBase
	//validate *validator.Validate

	wsRespAsyn *interaction.WsRespAsyn
	wsConn     *ws.WsConn
	stateMutex sync.Mutex
	//	mRWMutex   sync.RWMutex

	//Global minimum seq lock
	minSeqSvr        int64
	minSeqSvrRWMutex sync.RWMutex

	//Global cache seq map lock
	seqMsg      map[int32]*server_api_params.MsgData
	seqMsgMutex sync.RWMutex

	//	receiveMessageOpt sync.Map
	//Global message not disturb cache lock
	receiveMessageOpt      map[string]int32
	receiveMessageOptMutex sync.RWMutex
}

var userForSDK *UserRelated

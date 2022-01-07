package open_im_sdk

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"open_im_sdk/internal/controller/conversation_msg"
	"open_im_sdk/internal/controller/friend"
	"open_im_sdk/internal/controller/group"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"sync"
)

type UserRelated struct {
	ConversationCh chan utils.cmd2Value //cmd：

	token          string
	loginUserID    string
	wsNotification map[string]chan utils.GeneralWsResp
	wsMutex        sync.RWMutex
	utils.IMManager
	friend.Friend
	conversation_msg.ConversationListener
	group.groupListener

	db       *sql.DB
	imdb     *gorm.DB
	validate *validator.Validate

	mRWMutex   sync.RWMutex
	stateMutex sync.Mutex

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

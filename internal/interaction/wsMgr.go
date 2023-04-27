package interaction

import (
	"context"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"sync"
)

type LongConnMgr struct {
	w sync.Mutex
	// The long connection,can be set tcp or websocket.
	conn     LongConn
	listener open_im_sdk_callback.OnConnListener
	// Buffered channel of outbound messages.
	send               chan Message
	pushMsgAndMaxSeqCh chan common.Cmd2Value
	//
	closedErr  error
	ctx        *ConnContext
	isCompress bool
	connStatus int
	syncer     *WsRespAsyn
	encoder    Encoder
	compressor Compressor
}

func NewLongConnMgr(ctx context.Context, listener open_im_sdk_callback.OnConnListener, pushMsgAndMaxSeqCh chan common.Cmd2Value) *LongConnMgr {
	return &LongConnMgr{listener: listener, pushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh}
}

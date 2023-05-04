package interaction

import (
	"sync"
)

// 1.协程消息的收发模块一直运行，直到收到关闭信号（一般是用户退出登陆）
// 2.消息调用模块传递channle，收到消息后，调用channel，将消息传递给协程消息收发模块
// 3.消息通过ws发送后
// 4.

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	w sync.Mutex
	// The long connection,can be set tcp or websocket.
	conn LongConn

	// Buffered channel of outbound messages.
	send chan Message
	//
	closedErr  error
	ctx        *ConnContext
	isCompress bool
	connStatus int
	syncer     *WsRespAsyn
	encoder    Encoder
	compressor Compressor
}

package interaction

import (
	"net/http"
	"time"
)

type PongHandler func(string) error
type LongConn interface {
	//Close this connection
	Close() error
	// WriteMessage Write message to connection,messageType means data type,can be set binary(2) and text(1).
	WriteMessage(messageType int, message []byte) error
	// ReadMessage Read message from connection.
	ReadMessage() (int, []byte, error)
	//SetReadTimeout sets the read deadline on the underlying network connection,
	//after a read has timed out, will return an error.
	SetReadDeadline(timeout time.Duration) error
	//SetWriteTimeout sets to write deadline when send message,when read has timed out,will return error.
	SetWriteDeadline(timeout time.Duration) error
	// Dial Try to dial a connection,url must set auth args,header can control compress data
	Dial(urlStr string, requestHeader http.Header) (*http.Response, error)
	// IsNil Whether the connection of the current long connection is nil
	IsNil() bool
	// SetConnNil Set the connection of the current long connection to nil
	SetConnNil()
	// SetReadLimit sets the maximum size for a message read from the peer.bytes
	SetReadLimit(limit int64)
	SetPongHandler(handler PongHandler)
	// GenerateLongConn Check the connection of the current and when it was sent are the same
	GenerateLongConn(w http.ResponseWriter, r *http.Request) error
	// CheckSendConnDiffNow Check the connection of the current and when it was sent are the same
	CheckSendConnDiffNow() bool
	// LocalAddr returns the local network address.
	LocalAddr() string
}

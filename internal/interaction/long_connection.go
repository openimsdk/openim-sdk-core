package interaction

import "net/http"

type LongConn interface {
	//Close this connection
	Close() error
	//Write message to connection,messageType means data type,can be set binary(2) and text(1).
	WriteMessage(messageType int, message []byte) error
	//Read message from connection.
	ReadMessage() (int, []byte, error)
	//SetReadTimeout sets the read deadline on the underlying network connection,
	//after a read has timed out, will return an error.
	SetReadTimeout(timeout int) error
	//SetWriteTimeout sets to write deadline when send message,when read has timed out,will return error.
	SetWriteTimeout(timeout int) error
	//Try to dial a connection,url must set auth args,header can control compress data
	Dial(urlStr string, requestHeader http.Header) (*http.Response, error)
	//Whether the connection of the current long connection is nil
	IsNil() bool
	//Set the connection of the current long connection to nil
	SetConnNil()
	//Check the connection of the current and when it was sent are the same
	CheckSendConnDiffNow() bool
}

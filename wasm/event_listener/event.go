package event_listener

import (
	"open_im_sdk/pkg/utils"
	"syscall/js"
)

type EventData struct {
	Event       string      `json:"event"`
	ErrCode     int32       `json:"errCode"`
	ErrMsg      string      `json:"errMsg"`
	Data        interface{} `json:"data,omitempty"`
	OperationID string      `json:"operationID"`
	callback    *js.Value
}

func (e *EventData) GetOperationID() string {
	return e.OperationID
}

func NewEventData(callback *js.Value) *EventData {
	return &EventData{callback: callback}
}
func (e *EventData) SendMessage() {
	e.callback.Invoke(utils.StructToJsonString(e))
}
func (e *EventData) SetEvent(event string) CallbackWriter {
	e.Event = event
	return e
}

func (e *EventData) SetData(data interface{}) CallbackWriter {
	e.Data = data
	return e
}
func (e *EventData) SetErrCode(errCode int32) CallbackWriter {
	e.ErrCode = errCode
	return e
}
func (e *EventData) SetOperationID(operationID string) CallbackWriter {
	e.OperationID = operationID
	return e
}
func (e *EventData) SetErrMsg(errMsg string) CallbackWriter {
	e.ErrMsg = errMsg
	return e
}

package conversation_msg

import (
	"context"
	"encoding/json"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"open_im_sdk/internal/file"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/sdk_struct"
)

func NewFileCallback(ctx context.Context, progress func(progress int), msg *sdk_struct.MsgStruct, db db_interface.DataBase) file.PutFileCallback {
	if msg.AttachedInfoElem == nil {
		msg.AttachedInfoElem = &sdk_struct.AttachedInfoElem{}
	}
	if msg.AttachedInfoElem.Progress == nil {
		msg.AttachedInfoElem.Progress = &sdk_struct.UploadProgress{}
	}
	return &FileCallback{ctx: ctx, progress: progress, msg: msg, db: db}
}

type FileCallback struct {
	ctx      context.Context
	db       db_interface.DataBase
	msg      *sdk_struct.MsgStruct
	progress func(progress int)
}

func (c *FileCallback) Open(size int64) {}

func (c *FileCallback) HashProgress(current, total int64) {}

func (c *FileCallback) HashComplete(hash string, total int64) {}

func (c *FileCallback) PutStart(current, total int64) {}

func (c *FileCallback) PutProgress(save int64, current, total int64) {
	c.msg.AttachedInfoElem.Progress.Save = save
	c.msg.AttachedInfoElem.Progress.Current = current
	c.msg.AttachedInfoElem.Progress.Total = total
	data, err := json.Marshal(c.msg.AttachedInfoElem)
	if err != nil {
		panic(err)
	}
	if err := c.db.UpdateMessageByClientMsgID(c.ctx, c.msg.ClientMsgID, map[string]any{"attached_info": string(data)}); err != nil {
		log.ZError(c.ctx, "update PutProgress message attached info failed", err)
	}
	c.progress(int(float64(current) / float64(total)))
}

func (c *FileCallback) PutComplete(total int64, putType int) {
	c.msg.AttachedInfoElem.Progress = nil
	data, err := json.Marshal(c.msg.AttachedInfoElem)
	if err != nil {
		panic(err)
	}
	if err := c.db.UpdateMessageByClientMsgID(c.ctx, c.msg.ClientMsgID, map[string]any{"attached_info": string(data)}); err != nil {
		log.ZError(c.ctx, "update PutComplete message attached info failed", err)
	}
}

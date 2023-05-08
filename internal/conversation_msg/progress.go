package conversation_msg

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"open_im_sdk/internal/file"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/sdk_struct"
)

func NewFileCallback(ctx context.Context, progress func(progress int), msg *sdk_struct.MsgStruct, db db_interface.DataBase) file.PutFileCallback {
	if msg.AttachedInfoElem.Progress == nil {
		msg.AttachedInfoElem.Progress = &sdk_struct.UploadProgress{}
	}
	return &FileCallback{progress: progress, msg: msg, db: db}
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
	if err := c.db.UpdateMessageAttachedInfo(c.ctx, c.msg); err != nil {
		log.ZError(c.ctx, "update message attached info failed", err)
	}
	c.progress(int(float64(current) / float64(total)))
}

func (c *FileCallback) PutComplete(total int64, putType int) {}

package conversation_msg

import (
	"open_im_sdk/internal/file"
	"open_im_sdk/sdk_struct"
)

func NewFileCallback(progress func(progress int), p *sdk_struct.MsgStruct) file.PutFileCallback {
	if p.AttachedInfoElem.Progress == nil {
		p.AttachedInfoElem.Progress = &sdk_struct.UploadProgress{}
	}
	return &FileCallback{progress: progress, up: p.AttachedInfoElem.Progress}
}

type FileCallback struct {
	progress func(progress int)
	up       *sdk_struct.UploadProgress
}

func (c *FileCallback) Open(size int64) {}

func (c *FileCallback) HashProgress(current, total int64) {}

func (c *FileCallback) HashComplete(hash string, total int64) {}

func (c *FileCallback) PutStart(current, total int64) {}

func (c *FileCallback) PutProgress(save int64, current, total int64) {
	c.up.Save = save
	c.up.Total = total
	c.up.Current = current
	c.progress(int(float64(current) / float64(total)))
}

func (c *FileCallback) PutComplete(total int64, putType int) {}

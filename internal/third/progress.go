package third

import (
	"context"
)

type Progress interface {
	OnProgress(current int64, size int64)
}

type progressConvert struct {
	ctx context.Context
	p   Progress
}

func (p *progressConvert) Open(size int64) {
	p.p.OnProgress(0, size)
}

func (p *progressConvert) PartSize(partSize int64, num int) {}

func (p *progressConvert) HashPartProgress(index int, size int64, partHash string) {}

func (p *progressConvert) HashPartComplete(partsHash string, fileHash string) {}

func (p *progressConvert) UploadID(uploadID string) {}

func (p *progressConvert) UploadPartComplete(index int, partSize int64, partHash string) {}

func (p *progressConvert) UploadComplete(fileSize int64, streamSize int64, storageSize int64) {
	//log.ZDebug(p.ctx, "upload log progress", "fileSize", fileSize, "current", streamSize)
	p.p.OnProgress(streamSize, fileSize)
}

func (p *progressConvert) Complete(size int64, url string, typ int) {
	p.p.OnProgress(size, size)
}

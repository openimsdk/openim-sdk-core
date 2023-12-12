package third

type Progress interface {
	OnProgress(current int64, size int64)
}

type progressConvert struct {
	p Progress
}

func (p *progressConvert) Open(size int64) {}

func (p *progressConvert) PartSize(partSize int64, num int) {}

func (p *progressConvert) HashPartProgress(index int, size int64, partHash string) {}

func (p *progressConvert) HashPartComplete(partsHash string, fileHash string) {}

func (p *progressConvert) UploadID(uploadID string) {}

func (p *progressConvert) UploadPartComplete(index int, partSize int64, partHash string) {}

func (p *progressConvert) UploadComplete(fileSize int64, streamSize int64, storageSize int64) {
	if p != nil {
		p.p.OnProgress(streamSize, fileSize)
	}
}

func (p *progressConvert) Complete(size int64, url string, typ int) {}

package file

type PutFileCallback interface {
	Open(size int64)
	HashProgress(current, total int64)
	HashComplete(hash string, total int64)
	PutStart(current, total int64)
	PutProgress(save int64, current, total int64)
	PutComplete(total int64, putType int)
}

package third

import (
	"sync"

	"github.com/openimsdk/openim-sdk-core/v3/internal/third/file"
)

type Third struct {
	platform      int32
	loginUserID   string
	appFramework  string
	LogFilePath   string
	fileUploader  *file.File
	logUploadLock sync.Mutex
}

func (t *Third) SetPlatform(platform int32) {
	t.platform = platform
}

func (t *Third) SetLoginUserID(loginUserID string) {
	t.loginUserID = loginUserID
}

func (t *Third) SetAppFramework(appFramework string) {
	t.appFramework = appFramework
}

func (t *Third) SetLogFilePath(LogFilePath string) {
	t.LogFilePath = LogFilePath
}

func NewThird(fileUploader *file.File) *Third {
	return &Third{fileUploader: fileUploader}
}

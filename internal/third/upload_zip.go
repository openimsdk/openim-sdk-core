package third

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/internal/third/file"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/openim-sdk-core/v3/version"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/third"
)

func (c *Third) uploadZipAndReport(ctx context.Context, zipDir string, files []string, namePrefix, cause, ex, cancelID string, progress Progress) error {
	if len(files) == 0 {
		return fmt.Errorf("zip files empty")
	}
	if zipDir == "" {
		zipDir = filepath.Dir(files[0])
	}
	if cancelID != "" {
		var cancel context.CancelFunc
		ctx, cancel = c.fileUploader.WithUploadCancel(ctx, cancelID)
		defer cancel()
	}

	zippath := filepath.Join(zipDir, fmt.Sprintf("%s_%d.zip", time.Now().Format("2006-01-02_15-04-05_999"), rand.Uint32()))
	defer os.Remove(zippath)

	if err := zipFiles(ctx, zippath, files); err != nil {
		return err
	}

	reqUpload := &file.UploadFileReq{
		Filepath:    zippath,
		Name:        fmt.Sprintf("%s_%s_%s_%s_%s_%s", namePrefix, c.loginUserID, c.appFramework, constant.PlatformID2Name[int(c.platform)], version.Version, filepath.Base(zippath)),
		Cause:       cause,
		ContentType: "application/zip",
		CancelID:    cancelID,
	}
	var cb file.UploadFileCallback
	if progress != nil {
		cb = &progressConvert{ctx: ctx, p: progress}
	}
	resp, err := c.fileUploader.UploadFile(ctx, reqUpload, cb)
	if err != nil {
		return err
	}
	reqLog := &third.UploadLogsReq{
		Platform:     c.platform,
		AppFramework: c.appFramework,
		Version:      version.Version,
		FileURLs:     []*third.FileURL{{Filename: zippath, URL: resp.URL}},
		Ex:           ex,
	}
	return api.UploadLogs.Execute(ctx, reqLog)
}

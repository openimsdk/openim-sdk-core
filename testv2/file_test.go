package testv2

import (
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"open_im_sdk/internal/file"
	"open_im_sdk/open_im_sdk"
	"path/filepath"
	"testing"
	"time"
)

type FilePutCallback struct{}

func (c *FilePutCallback) Open(size int64) {
	fmt.Println("open put file", size)
}

func (c *FilePutCallback) HashProgress(current, total int64) {
	//fmt.Println("hash", current, total)
}

func (c *FilePutCallback) HashComplete(hash string, total int64) {
	fmt.Println("hash complete", hash, total)
}

func (c *FilePutCallback) PutStart(current, total int64) {
	fmt.Println("put start", current, total)
}

func (c *FilePutCallback) PutProgress(save int64, current, total int64) {
	fmt.Printf("put progress [%d/%d] put %f%% save %f%%\n", current, total, float64(current)/float64(total)*100, float64(save)/float64(total)*100)
}

func (c *FilePutCallback) PutComplete(total int64, putType int) {
	fmt.Println("put complete", total, putType)
}

func TestPut(t *testing.T) {
	ctx := mcontext.NewCtx("123456")

	putID := "asdkmjakslaaa"

	go func() {
		time.Sleep(time.Second * 3)
		fmt.Println("###################### CANCEL PUT ######################")
		open_im_sdk.UserForSDK.File().Cancel(ctx, putID)
	}()

	req := &file.PutArgs{
		PutID:    putID,
		Filepath: "C:\\Users\\Admin\\Desktop\\VMware-workstation-full-17.0.0-20800274.exe",
	}
	req.Name = filepath.Base(req.Filepath)
	str, err := open_im_sdk.UserForSDK.File().PutFile(ctx, req, &FilePutCallback{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("url", str)
}

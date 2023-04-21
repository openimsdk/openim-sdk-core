package testv2

import (
	"context"
	"fmt"
	"math/rand"
	"open_im_sdk/internal/file"
	"open_im_sdk/open_im_sdk"
	"strconv"
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
	fmt.Printf("put progress [%d/%d] %d\n", current, total, save)
}

func (c *FilePutCallback) PutComplete(total int64, putType int) {
	fmt.Println("put complete", total, putType)
}

func TestPut(t *testing.T) {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		time.Sleep(time.Second * 2)
		cancel()
		fmt.Println("-----------------------------------cancel")
	}()

	//req := &file.PutArgs{
	//	Filepath:    "C:\\Users\\Admin\\Desktop\\landscape.png",
	//	Name:        "landscape.png",
	//	ContentType: "image/png",
	//}
	req := &file.PutArgs{
		PutID:       strconv.FormatUint(rand.Uint64(), 10),
		Filepath:    "C:\\Users\\Admin\\Desktop\\VMware-workstation-full-17.0.0-20800274.exe",
		Name:        "VMware-workstation-full-17.0.0-20800274.exe",
		ContentType: "app/exe",
	}
	str, err := open_im_sdk.UserForSDK.File().putFilePath(ctx, req, &FilePutCallback{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("url", str)
}

func TestContinuePut(t *testing.T) {
	req := &file.PutArgs{
		PutID:       "127519572988560128231",
		Filepath:    "C:\\Users\\Admin\\Desktop\\VMware-workstation-full-17.0.0-20800274.exe",
		Name:        "VMware-workstation-full-17.0.0-20800274.exe",
		ContentType: "app/exe",
	}
	str, err := open_im_sdk.UserForSDK.File().PutFile(ctx, req, &FilePutCallback{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("url", str)
}

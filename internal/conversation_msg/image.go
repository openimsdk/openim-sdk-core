package conversation_msg

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/tools/errs"
)

func getImageInfo(filePath string) (*sdk_struct.ImageInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errs.WrapMsg(err, "image file  open err")
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, errs.WrapMsg(err, "image file  decode err")
	}
	size := img.Bounds().Max
	return &sdk_struct.ImageInfo{Width: int32(size.X), Height: int32(size.Y), Type: "image/" + format, Size: info.Size()}, nil
}

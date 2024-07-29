package initialization

import (
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/openim-sdk-core/v3/version"
	"github.com/openimsdk/tools/log"
)

const (
	rotateCount  uint = 0
	rotationTime uint = 24
)

func InitLog(cf sdk_struct.IMConfig) error {
	if err := log.InitFromConfig("open-im-sdk-core", "", int(cf.LogLevel), cf.IsLogStandardOutput, false, cf.LogFilePath, rotateCount, rotationTime, version.Version, true); err != nil {
		fmt.Println("log init failed ", err.Error())
		return err
	}
	return nil
}

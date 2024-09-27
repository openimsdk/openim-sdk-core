package initialization

import (
	"fmt"

	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/openim-sdk-core/v3/version"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/log"
)

const (
	rotateCount  uint = 1
	rotationTime uint = 24
)

func InitLog(cf sdk_struct.IMConfig) error {
	if err := log.InitLoggerFromConfig("open-im-sdk-core", "", cf.SystemType, constant.PlatformID2Name[int(cf.PlatformID)], int(cf.LogLevel), cf.IsLogStandardOutput, false, cf.LogFilePath, rotateCount, rotationTime, version.Version, true); err != nil {
		fmt.Println("log init failed ", err.Error())
		return err
	}
	return nil
}

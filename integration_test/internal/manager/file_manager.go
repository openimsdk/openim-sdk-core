package manager

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/decorator"
	"github.com/openimsdk/tools/errs"
	"os"
)

type TestFileManager struct {
	*MetaManager
}

func NewFileManager(m *MetaManager) *TestFileManager {
	return &TestFileManager{m}
}

func (m *TestFileManager) DeleteLocalDB(ctx context.Context) error {
	defer decorator.FuncLog(ctx)()

	conf := config.GetConf()
	err := os.RemoveAll(conf.DataDir)
	if err != nil {
		return errs.WrapMsg(err, "remove db failed")
	}
	err = os.MkdirAll(conf.DataDir, os.ModePerm)
	if err != nil {
		return errs.WrapMsg(err, "make db dir failed")
	}
	return nil
}

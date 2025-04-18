package cliconf

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/log"
)

const (
	ccConversationActiveNumKey = "CONVERSATION_ACTIVE_NUM"
)

const (
	ccConversationActiveNumDefault = 50
)

type ClientConfig struct {
	ConversationActiveNum int
	RawConfig             map[string]string
}

type currentClientConfig struct {
	wait   chan struct{}
	err    error
	config *ClientConfig
}

type clientConfig struct {
	userID string
	config atomic.Pointer[currentClientConfig]
}

func (c *clientConfig) parseServerUserConfig(configs map[string]string) (*ClientConfig, error) {
	var config ClientConfig
	config.ConversationActiveNum, _ = strconv.Atoi(configs[ccConversationActiveNumKey])
	if config.ConversationActiveNum <= 0 {
		config.ConversationActiveNum = ccConversationActiveNumDefault
	}
	config.RawConfig = configs
	return &config, nil
}

func (c *clientConfig) getServerUserConfig(ctx context.Context) (*ClientConfig, error) {
	configs, err := api.ExtractField(ctx, api.UserClientConfig.Invoke, &user.GetUserClientConfigReq{UserID: c.userID}, (*user.GetUserClientConfigResp).GetConfigs)
	if err != nil {
		return nil, err
	}
	return c.parseServerUserConfig(configs)
}

func (c *clientConfig) asyncGetConfig(ctx context.Context, curr *currentClientConfig) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	for {
		configs, err := c.getServerUserConfig(ctx)
		if c.config.Load() != curr {
			log.ZDebug(ctx, "get config end, but call clear config", "config", configs, "error", err)
			continue
		}
		if err != nil {
			curr.err = err
			close(curr.wait)
			log.ZWarn(ctx, "get config failed", err)
			return
		}
		curr.config = configs
		close(curr.wait)
		log.ZDebug(ctx, "get config success", "config", configs)
		return
	}
}

func (c *clientConfig) getCurrConfig(ctx context.Context) (*currentClientConfig, error) {
	for i := 0; ; i++ {
		if i > 10 {
			return nil, sdkerrs.ErrSdkInternal.WrapMsg("get client config timeout")
		}
		curr := c.config.Load()
		if curr != nil {
			return curr, nil
		}
		curr = &currentClientConfig{wait: make(chan struct{})}
		if !c.config.CompareAndSwap(nil, curr) {
			close(curr.wait)
			continue
		}
		go c.asyncGetConfig(context.WithoutCancel(ctx), curr)
		return curr, nil
	}
}

func (c *clientConfig) ClearConfig() {
	c.config.Swap(nil)
	log.ZDebug(context.Background(), "clear config")
}

func (c *clientConfig) GetConfig(ctx context.Context) (*ClientConfig, error) {
	curr, err := c.getCurrConfig(ctx)
	if err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, context.Cause(ctx)
	case <-curr.wait:
		if curr.err != nil {
			return nil, curr.err
		}
		return curr.config, nil
	}
}

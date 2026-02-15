// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// 本地启动 Go 版 SDK 的 cmd 入口，用于联调或本地验证。
//
// 与 Rust im_client 一致：使用内置默认密码获取 token，无参即可启动。
//
//	go run ./cmd/sdk
//
// 仅支持密码登录，可选 -phone / -password 覆盖默认。
package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/openimsdk/protocol/constant"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

// 与 Rust im_client 相同的默认密码，便于无参启动.
const defaultPassword = "284f3d09ea0695538e4ded1c1766d73a"

var (
	apiAddr    = flag.String("api", "http://127.0.0.1:10002", "IM API 地址")
	wsAddr     = flag.String("ws", "ws://127.0.0.1:10001", "WebSocket 地址")
	phone      = flag.String("phone", "17764228283", "手机号，与 Rust 一致")
	password   = flag.String("password", defaultPassword, "密码，与 Rust 内置默认一致")
	areaCode   = flag.String("area", "+86", "区号，与 Rust 一致")
	accountAPI = flag.String("account-api", "http://localhost:10008", "账号登录服务地址，与 Rust 一致")
	dataDir    = flag.String("data-dir", "./sdk_data", "本地数据目录（数据库与缓存）")
)

// POST /account/login 请求（与 Rust 一致：仅 areaCode + phoneNumber + password + platform）.
type accountLoginReq struct {
	AreaCode    string `json:"areaCode"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Platform    int32  `json:"platform"`
	VerifyCode  string `json:"verifyCode"`
}

type accountLoginResp struct {
	ErrCode int32             `json:"errCode"`
	ErrMsg  string            `json:"errMsg"`
	Data    *accountLoginData `json:"data"`
}

type accountLoginData struct {
	ImToken   string `json:"imToken"`
	ChatToken string `json:"chatToken"`
	UserID    string `json:"userID"`
}

// getTokenByAccount 调用 POST {accountAPI}/account/login（与 Rust 一致，仅传密码不传验证码）.
func getTokenByAccount(accountBase, phoneNum, pwd, area string, platform int32) (uid, imToken string, err error) {
	url := accountBase + "/account/login"
	reqBody := accountLoginReq{
		AreaCode:    area,
		PhoneNumber: phoneNum,
		// Password:    pwd,
		Platform:   platform,
		VerifyCode: "666666",
	}
	body, _ := json.Marshal(reqBody)
	if os.Getenv("SDK_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "[SDK] POST %s body: %s\n", url, string(body))
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", "", fmt.Errorf("account login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// 账号服务中间件要求 header 带 operationID
	b := make([]byte, 16)
	if _, rErr := rand.Read(b); rErr != nil {
		return "", "", fmt.Errorf("operationID: %w", rErr)
	}
	req.Header.Set("operationID", hex.EncodeToString(b))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("account login request: %w", err)
	}
	defer resp.Body.Close()
	rawBody, _ := io.ReadAll(resp.Body)
	var r accountLoginResp
	if err := json.Unmarshal(rawBody, &r); err != nil {
		return "", "", fmt.Errorf("account login decode: %w", err)
	}
	if r.ErrCode != 0 {
		if os.Getenv("SDK_DEBUG") != "" {
			fmt.Fprintf(os.Stderr, "[SDK] account login 响应 body: %s\n", string(rawBody))
		}

		return "", "", fmt.Errorf("account login errCode=%d errMsg=%s", r.ErrCode, r.ErrMsg)
	}
	if r.Data == nil {
		return "", "", fmt.Errorf("account login: data is nil")
	}

	return r.Data.UserID, r.Data.ImToken, nil
}

func main() {
	flag.Parse()

	// 登录与 SDK 配置共用同一平台标识，与 Rust 一致（Web=5）
	platformID := int32(constant.WebPlatformID)

	pwd := *password
	if pwd == "" {
		pwd = os.Getenv("SDK_PASSWORD")
	}
	uid, imToken, err := getTokenByAccount(*accountAPI, *phone, pwd, *areaCode, platformID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Login] 账号登录失败:", err)
		fmt.Fprintln(os.Stderr, "[Login] 可设置 SDK_DEBUG=1 查看请求体，或通过 -phone / -password 覆盖默认")
		os.Exit(1)
	}
	fmt.Println("[Login] 账号登录成功, userID:", uid)

	logPath := *dataDir + "/logs"
	// 日志级别: 0=Panic 1=Fatal 2=Error 3=Warn 4=Info 5=Debug，设为 Debug 便于观察
	const logLevelDebug = 4
	cfg := sdk_struct.IMConfig{
		SystemType:          "cmd",
		PlatformID:          platformID,
		ApiAddr:             *apiAddr,
		WsAddr:              *wsAddr,
		DataDir:             *dataDir,
		LogLevel:            logLevelDebug,
		IsLogStandardOutput: true,
		LogFilePath:         logPath,
	}
	configJSON, err := json.Marshal(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "marshal config:", err)
		os.Exit(1)
	}

	fmt.Println("[SDK] 数据目录:", *dataDir)
	fmt.Println("[SDK] 日志路径:", logPath)
	fmt.Println("[SDK] 日志级别: Debug (5)")

	connListener := &connListener{}
	if !open_im_sdk.InitSDK(connListener, "cmd_sdk", string(configJSON)) {
		fmt.Fprintln(os.Stderr, "InitSDK 失败")
		os.Exit(1)
	}
	fmt.Println("[SDK] InitSDK 成功")

	// 不设置会话/消息监听器，统一使用 SDK 默认空回调（em.go）

	loginCB := &loginCallback{}
	open_im_sdk.Login(loginCB, "cmd_login", uid, imToken)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("[SDK] 收到退出信号，退出")
}

type connListener struct{}

func (c *connListener) OnConnecting() {
	fmt.Println("[Conn] OnConnecting")
}

func (c *connListener) OnConnectSuccess() {
	fmt.Println("[Conn] OnConnectSuccess")
}

func (c *connListener) OnConnectFailed(errCode int32, errMsg string) {
	fmt.Printf("[Conn] OnConnectFailed code=%d msg=%s\n", errCode, errMsg)
}

func (c *connListener) OnKickedOffline() {
	fmt.Println("[Conn] OnKickedOffline")
}

func (c *connListener) OnUserTokenExpired() {
	fmt.Println("[Conn] OnUserTokenExpired")
}

func (c *connListener) OnUserTokenInvalid(errMsg string) {
	fmt.Println("[Conn] OnUserTokenInvalid:", errMsg)
}

type loginCallback struct{}

func (l *loginCallback) OnError(errCode int32, errMsg string) {
	fmt.Printf("[Login] OnError code=%d msg=%s\n", errCode, errMsg)
	os.Exit(1)
}

func (l *loginCallback) OnSuccess(data string) {
	fmt.Println("[Login] OnSuccess, userID:", open_im_sdk.GetLoginUserID())
}

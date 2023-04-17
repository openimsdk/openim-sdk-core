package chao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"io"
	"net/http"
	"open_im_sdk/open_im_sdk"
	"os"
	"reflect"
	"time"
	"unsafe"
)

var HOST = "192.168.44.128"
var APIADDR = "http://" + HOST + ":10002"
var WSADDR = "ws://" + HOST + ":10001"

func GetUserToken(ctx context.Context, userID string) (string, error) {
	jsonReqData, err := json.Marshal(map[string]any{
		"userID":   userID,
		"platform": 1,
		"secret":   "openIM123",
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, APIADDR+"/auth/user_token", bytes.NewReader(jsonReqData))
	if err != nil {
		return "", err
	}
	req.Header.Set("operationID", ctx.Value("operationID").(string))
	client := http.Client{Timeout: time.Second * 3}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	type Result struct {
		ErrCode int    `json:"errCode"`
		ErrMsg  string `json:"errMsg"`
		ErrDlt  string `json:"errDlt"`
		Data    struct {
			Token             string `json:"token"`
			ExpireTimeSeconds int    `json:"expireTimeSeconds"`
		} `json:"data"`
	}
	var result Result
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("errCode:%d, errMsg:%s, errDlt:%s", result.ErrCode, result.ErrMsg, result.ErrDlt)
	}
	return result.Data.Token, nil
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetResValue[T any](value T, err error) T {
	CheckErr(err)
	return value
}

func call[T any](ctx context.Context, fn any, args ...any) (T, error) {
	var t T
	operationID := ctx.Value("operationID").(string)
	var err error
	var data string
	var done = make(chan struct{})
	onErr := func(errCode int32, errMsg string) {
		err = fmt.Errorf("errCode:%d, errMsg:%s", errCode, errMsg)
		close(done)
	}
	onSucc := func(v string) {
		data = v
		close(done)
	}
	in := make([]reflect.Value, 0, len(args)+2)
	in = append(in, reflect.ValueOf(NewCbFn(onErr, onSucc)))
	in = append(in, reflect.ValueOf(operationID))
	for i := range args {
		in = append(in, reflect.ValueOf(args[i]))
	}
	reflect.ValueOf(fn).Call(in)
	<-done
	if err != nil {
		return t, err
	}
	if _, ok := any(t).(string); ok {
		return *(*T)(unsafe.Pointer(&data)), nil
	}
	if err := json.Unmarshal([]byte(data), &t); err != nil {
		return t, err
	}
	return t, nil
}

func CallRaw(ctx context.Context, fn any, args ...any) string {
	return Call[string](ctx, fn, args...)
}

func Call[T any](ctx context.Context, fn any, args ...any) T {
	return GetResValue(call[T](ctx, fn, args...))
}

func Main() {
	operationID := "op123"
	ctx := context.WithValue(context.Background(), "operationID", operationID)
	userID := "123456"
	token := GetResValue(GetUserToken(ctx, userID))
	fmt.Println("token:", token)

	open_im_sdk.InitSDK(&Listener{}, operationID, string(GetResValue(json.Marshal(GetConf()))))
	CallRaw(ctx, open_im_sdk.Login, userID, token)

	req := group.CreateGroupReq{
		InitMembers: nil,
		OwnerUserID: userID,
		GroupInfo: &sdkws.GroupInfo{
			GroupName: "testgroup",
		},
	}
	info := Call[*group.CreateGroupResp](ctx, open_im_sdk.CreateGroupV2, string(GetResValue(json.Marshal(&req)))).GroupInfo
	fmt.Println("-------------------------------------")
	fmt.Println(info.String())
	fmt.Println("-------------------------------------")
	os.Exit(0)
	//time.Sleep(time.Second * 100)

}

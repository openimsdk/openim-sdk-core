package chao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"reflect"
	"runtime"
	"runtime/debug"
	"time"
	"unsafe"
)

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

func PrintTest() {
	var color bool
	color = true
	pc, _, _, _ := runtime.Caller(0)
	name := path.Base(runtime.FuncForPC(pc).Name())
	r := recover()
	if r == nil {
		s := fmt.Sprintf("[DoTest] %s %s pass", time.Now().Format("2006-01-02 15:04:05"), name)
		if color {
			s = fmt.Sprintf("%c[1;42;32m%s%c[0m\n\n", 0x1B, s, 0x1B)
		}
		fmt.Println(s)
		return
	}
	s := fmt.Sprintf("[DoTest] %s %s fail %v", time.Now().Format("2006-01-02 15:04:05"), name, r)
	if color {
		s = fmt.Sprintf("%c[1;45;31m%s%c[0m\n\n", 0x1B, s, 0x1B)
	}
	fmt.Println(s)
	fmt.Println(string(debug.Stack()))
}

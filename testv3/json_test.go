package testv3

import (
	"encoding/json"
	"fmt"
	"testing"
)

// 定义JSON解码后的结构体类型
type AccountCheckRespSingleUserStatus struct {
	UserID        string `json:"userID"`
	AccountStatus string `json:"accountStatus"`
}

type AccountCheckRespData struct {
	Results []AccountCheckRespSingleUserStatus `json:"results"`
}

type AccountCheckResp struct {
	//ErrCode int                  `json:"errCode"`
	//ErrMsg  string               `json:"errMsg"`
	//ErrDlt  string               `json:"errDlt"`
	Data AccountCheckRespData `json:"data"`
	//Results []*AccountCheckRespSingleUserStatus `json:"results"`
}

func Test_unmarshal(t *testing.T) {
	// 从字符串解码json到byte数组
	r := []byte(`{"errCode":0,"errMsg":"","errDlt":"","data":{"results":[{"userID":"169.254.125.227_reliability_1687675248__0","accountStatus":"unregistered"}]}}`)

	// 解码JSON到结构体
	var resp AccountCheckResp
	err := json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 打印结果
	//fmt.Printf("errCode: %d\n", resp.ErrCode)
	//fmt.Printf("errMsg: %s\n", resp.ErrMsg)
	//fmt.Printf("errDlt: %s\n", resp.ErrDlt)
	//for _, r := range resp.Data.Results {
	//	fmt.Printf("userID: %s, accountStatus: %s\n", r.UserID, r.AccountStatus)
	//}
}

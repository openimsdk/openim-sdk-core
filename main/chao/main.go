package chao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/internal/util"
	"open_im_sdk/open_im_sdk"
)

func TestCreateGroupV2(ctx context.Context, userID string) {
	defer PrintTest()
	req := group.CreateGroupReq{
		InitMembers: nil,
		OwnerUserID: userID,
		GroupInfo: &sdkws.GroupInfo{
			GroupName: "test",
		},
	}
	info := Call[sdkws.GroupInfo](ctx, open_im_sdk.CreateGroupV2, string(GetResValue(json.Marshal(&req))))
	fmt.Println(info.String())
}

func Main() {

	util.BaseURL = APIADDR
	operationID := "op123"
	ctx := context.WithValue(context.Background(), "operationID", operationID)
	userID := "123456"
	token := GetResValue(GetUserToken(ctx, userID))
	fmt.Println("token:", token)
	util.Token = token
	open_im_sdk.InitSDK(&Listener{}, operationID, string(GetResValue(json.Marshal(GetConf()))))
	CallRaw(ctx, open_im_sdk.Login, userID, token)
	TestCreateGroupV2(ctx, userID)

}

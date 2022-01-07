package open_im_sdk

type GetDesignatedFriendsInfoParams []string
type GetDesignatedFriendsInfoCallback []LocalFriend

type AddFriendParams struct {
	ToUserID string `json:"toUserID" validate:"required"`
	ReqMsg   string `json:"reqMsg"`
}

const AddFriendCallback = SuccessCallbackDefault

type GetRecvFriendApplicationListParams struct{}
type GetRecvFriendApplicationListCallback []*LocalFriendRequest

type GetSendFriendApplicationListParams struct{}
type GetSendFriendApplicationListCallback []LocalFriendRequest

type ProcessFriendApplicationParams struct {
	ToUserID  string `json:"toUserID" validate:"required"`
	HandleMsg string `json:"handleMsg"`
}

//type ProcessFriendApplicationCallback struct{}
const ProcessFriendApplicationCallback = SuccessCallbackDefault

type CheckFriendParams []string
type CheckFriendCallback []UserIDResult

type DeleteFriendParams string

//type DeleteFriendCallback struct{}
const DeleteFriendCallback = SuccessCallbackDefault

type GetFriendListParams struct{}
type GetFriendListCallback []LocalFriend

type SetFriendRemarkParams struct {
	ToUserID string `json:"toUserID" validate:"required"`
	Remark   string `json:"remark" validate:"required"`
}

//type SetFriendRemarkCallback struct{}
const SetFriendRemarkCallback = SuccessCallbackDefault

type AddBlackParams string

//type AddBlackCallback struct{}
const AddBlackCallback = SuccessCallbackDefault

type GetBlackListParams struct{}

//type GetBlackListCallback []LocalBlack
const GetBlackListCallback = SuccessCallbackDefault

type DeleteBlackParams string

//type DeleteBlackCallback struct{}
const DeleteBlackCallback = SuccessCallbackDefault

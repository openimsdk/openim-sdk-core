package sdk_params_callback

import "open_im_sdk/pkg/server_api_params"

type InviteCallback  *server_api_params.SignalInviteReply

type InviteInGroupCallback server_api_params.SignalInviteInGroupReply


type CancelCallback struct {

}

type RejectCallback struct {

}


type AcceptCallback struct {

}

type HungUpCallback struct {

}

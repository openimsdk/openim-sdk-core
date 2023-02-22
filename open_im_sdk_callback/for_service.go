package open_im_sdk_callback

type OnListenerForService interface {
	//有人申请进群
	OnGroupApplicationAdded(groupApplication string)
	//进群申请被同意
	OnGroupApplicationAccepted(groupApplication string)
	//有人申请添加你为好友
	OnFriendApplicationAdded(friendApplication string)
	//好友申请被同意
	OnFriendApplicationAccepted(groupApplication string)
	//收到新消息
	OnRecvNewMessage(message string)
}

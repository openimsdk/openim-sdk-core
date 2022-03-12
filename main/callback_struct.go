package main

type GoFunVoid = func()
type GoFunInt = func(int)
type GoFunString = func(string)
type GoFunIntString = func(int, string)

type Base struct {
	onError   GoFunIntString
	onSuccess GoFunString
}
func (s *Base) OnError(errCode int32, errMsg string) {
	s.onError(int(errCode), errMsg)
}

func (s *Base) OnSuccess(data string) {
	s.onSuccess(data)
}

type SendMsgCallBack struct {
	Base
	onProgress GoFunInt
}

func (s *SendMsgCallBack) OnError(errCode int32, errMsg string) {
	s.onError(int(errCode), errMsg)
}

func (s *SendMsgCallBack) OnSuccess(data string) {
	s.onSuccess(data)
}

func (s *SendMsgCallBack) OnProgress(progress int) {
	s.onProgress(progress)
}

type OnConnListener struct {
	onConnecting       GoFunVoid
	onConnectSuccess   GoFunVoid
	onConnectFailed    GoFunIntString
	onKickedOffline    GoFunVoid
	onUserTokenExpired GoFunVoid
}

func (s *OnConnListener) OnConnecting() {
	s.onConnecting()
}

func (s *OnConnListener) OnConnectSuccess() {
	s.onConnectSuccess()
}

func (s *OnConnListener) OnConnectFailed(errCode int32, errMsg string) {
	s.onConnectFailed(int(errCode), errMsg)
}

func (s *OnConnListener) OnKickedOffline() {
	s.onKickedOffline()
}

func (s *OnConnListener) OnUserTokenExpired() {
	s.onUserTokenExpired()
}

type OnGroupListener struct {
	onJoinedGroupAdded         GoFunString
	onJoinedGroupDeleted       GoFunString
	onGroupMemberAdded         GoFunString
	onGroupMemberDeleted       GoFunString
	onGroupApplicationAdded    GoFunString
	onGroupApplicationDeleted  GoFunString
	onGroupInfoChanged         GoFunString
	onGroupMemberInfoChanged   GoFunString
	onGroupApplicationAccepted GoFunString
	onGroupApplicationRejected GoFunString
}

func (s *OnGroupListener) OnJoinedGroupAdded(groupInfo string) {
	s.onJoinedGroupAdded(groupInfo)
}

func (s *OnGroupListener) OnJoinedGroupDeleted(groupInfo string) {
	s.onJoinedGroupDeleted(groupInfo)
}
func (s *OnGroupListener) OnGroupMemberAdded(groupMemberInfo string) {
	s.onGroupMemberAdded(groupMemberInfo)
}
func (s *OnGroupListener) OnGroupMemberDeleted(groupMemberInfo string) {
	s.onGroupMemberDeleted(groupMemberInfo)
}
func (s *OnGroupListener) OnGroupApplicationAdded(groupApplication string) {
	s.onGroupApplicationAdded(groupApplication)
}
func (s *OnGroupListener) OnGroupApplicationDeleted(groupApplication string) {
	s.onGroupApplicationDeleted(groupApplication)
}
func (s *OnGroupListener) OnGroupInfoChanged(groupInfo string) {
	s.onGroupInfoChanged(groupInfo)
}
func (s *OnGroupListener) OnGroupMemberInfoChanged(groupMemberInfo string) {
	s.onGroupMemberInfoChanged(groupMemberInfo)
}
func (s *OnGroupListener) OnGroupApplicationAccepted(groupApplication string) {
	s.onGroupApplicationAccepted(groupApplication)
}
func (s *OnGroupListener) OnGroupApplicationRejected(groupApplication string) {
	s.onGroupApplicationRejected(groupApplication)
}

type OnFriendshipListener struct {
	onFriendApplicationAdded    GoFunString
	onFriendApplicationDeleted  GoFunString
	onFriendApplicationAccepted GoFunString
	onFriendApplicationRejected GoFunString
	onFriendAdded               GoFunString
	onFriendDeleted             GoFunString
	onFriendInfoChanged         GoFunString
	onBlackAdded                GoFunString
	onBlackDeleted              GoFunString
}

func (s *OnFriendshipListener) OnFriendApplicationAdded(friendApplication string) {
	s.onFriendApplicationAdded(friendApplication)
}
func (s *OnFriendshipListener) OnFriendApplicationDeleted(friendApplication string) {
	s.onFriendApplicationDeleted(friendApplication)
}
func (s *OnFriendshipListener) OnFriendApplicationAccepted(groupApplication string) {
	s.onFriendApplicationAccepted(groupApplication)
}
func (s *OnFriendshipListener) OnFriendApplicationRejected(friendApplication string) {
	s.onFriendApplicationRejected(friendApplication)
}
func (s *OnFriendshipListener) OnFriendAdded(friendInfo string) {
	s.onFriendAdded(friendInfo)
}
func (s *OnFriendshipListener) OnFriendDeleted(friendInfo string) {
	s.onFriendDeleted(friendInfo)
}
func (s *OnFriendshipListener) OnFriendInfoChanged(friendInfo string) {
	s.onFriendInfoChanged(friendInfo)
}
func (s *OnFriendshipListener) OnBlackAdded(blackInfo string) {
	s.onBlackAdded(blackInfo)
}
func (s *OnFriendshipListener) OnBlackDeleted(blackInfo string) {
	s.onBlackDeleted(blackInfo)
}

type OnConversationListener struct {
	onSyncServerStart                GoFunVoid
	onSyncServerFinish               GoFunVoid
	onSyncServerFailed               GoFunVoid
	onNewConversation                GoFunString
	onConversationChanged            GoFunString
	onTotalUnreadMessageCountChanged GoFunInt
}

func (s *OnConversationListener) OnSyncServerStart() {
	s.onSyncServerStart()
}
func (s *OnConversationListener) OnSyncServerFinish() {
	s.onSyncServerFinish()
}
func (s *OnConversationListener) OnSyncServerFailed() {
	s.onSyncServerFailed()
}
func (s *OnConversationListener) OnNewConversation(conversationList string) {
	s.onNewConversation(conversationList)
}
func (s *OnConversationListener) OnConversationChanged(conversationList string) {
	s.onConversationChanged(conversationList)
}
func (s *OnConversationListener) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	s.onTotalUnreadMessageCountChanged(int(totalUnreadCount))
}

type OnAdvancedMsgListener struct {
	onRecvNewMessage     GoFunString
	onRecvC2CReadReceipt GoFunString
	onRecvMessageRevoked GoFunString
}

func (s *OnAdvancedMsgListener) OnRecvNewMessage(message string) {
	s.onRecvNewMessage(message)
}
func (s *OnAdvancedMsgListener) OnRecvC2CReadReceipt(msgReceiptList string) {
	s.onRecvC2CReadReceipt(msgReceiptList)
}
func (s *OnAdvancedMsgListener) OnRecvMessageRevoked(msgId string) {
	s.onRecvMessageRevoked(msgId)
}

type OnUserListener struct {
	onSelfInfoUpdated GoFunString
}

func (s *OnUserListener) OnSelfInfoUpdated(userInfo string) {
	s.onSelfInfoUpdated(userInfo)
}

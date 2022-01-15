package sdk_params_callback

import "open_im_sdk/pkg/db"

////////////////////////////////friend////////////////////////////////////
type FriendApplicationAddedCallback db.LocalFriendRequest
type FriendApplicationAcceptCallback db.LocalFriendRequest
type FriendApplicationRejectCallback db.LocalFriendRequest
type FriendApplicationDeletedCallback db.LocalFriendRequest
type FriendAddedCallback db.LocalFriend
type FriendDeletedCallback db.LocalFriend
type BlackAddCallback db.LocalBlack
type BlackDeletedCallback db.LocalBlack
type FriendInfoChangedCallback db.LocalFriend


////////////////////////////////group////////////////////////////////////


OnJoinedGroupAdded
OnJoinedGroupDeleted
OnMemberAdded
OnMemberDeleted
OnReceiveJoinApplication
OnApplicationAccept
OnApplicationReject
OnGroupInfoChanged
OnMemberInfoChanged
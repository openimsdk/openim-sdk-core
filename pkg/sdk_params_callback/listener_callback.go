package sdk_params_callback

import "open_im_sdk/pkg/db"

////////////////////////////////friend////////////////////////////////////
type FriendApplicationAddedCallback db.LocalFriendRequest
type FriendApplicationAcceptCallback db.LocalFriendRequest
type FriendApplicationRejectCallback db.LocalFriendRequest
type FriendApplicationDeletedCallback db.LocalFriendRequest
type FriendAddedCallback db.LocalFriend
type FriendDeletedCallback db.LocalFriend
type FriendInfoChangedCallback db.LocalFriend
type BlackAddCallback db.LocalBlack
type BlackDeletedCallback db.LocalBlack

////////////////////////////////group////////////////////////////////////

type JoinedGroupAddedCallback db.LocalGroup
type JoinedGroupDeletedCallback db.LocalGroup
type GroupMemberAddedCallback db.LocalGroupMember
type GroupMemberDeletedCallback db.LocalGroupMember
type GroupApplicationAddedCallback db.LocalAdminGroupRequest
type GroupApplicationDeletedCallback db.LocalAdminGroupRequest
type GroupApplicationAcceptCallback db.LocalAdminGroupRequest
type GroupApplicationRejectCallback db.LocalAdminGroupRequest
type GroupInfoChangedCallback db.LocalGroup
type GroupMemberInfoChangedCallback db.LocalGroupMember

//////////////////////////////user////////////////////////////////////////
type SelfInfoUpdatedCallback db.LocalUser

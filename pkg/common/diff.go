package common

import (
	"github.com/jinzhu/copier"
	"open_im_sdk/pkg/db"
	log2 "open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"reflect"
)

type diff interface {
	Key() string
	Value() interface{}
}

func CompFields(a interface{}, b interface{}, fields ...string) bool {
	return false
	//	at := reflect.TypeOf(a)
	av := reflect.ValueOf(a)
	bt := reflect.TypeOf(b)
	bv := reflect.ValueOf(b)

	av = reflect.ValueOf(av.Interface())

	_fields := make([]string, 0)
	if len(fields) > 0 {
		_fields = fields
	} else {
		for i := 0; i < bv.NumField(); i++ {
			_fields = append(_fields, bt.Field(i).Name)
		}
	}

	if len(_fields) == 0 {
		return false
	}

	for i := 0; i < len(_fields); i++ {
		name := _fields[i]
		f := av.Elem().FieldByName(name)
		bValue := bv.FieldByName(name)

		if f.IsValid() && f.Kind() == bValue.Kind() {
			f.Set(bValue)
		} else {

		}
	}
	return false
}

func friendCopyToLocal(localFriend *db.LocalFriend, apiFriend *server_api_params.FriendInfo) {
	copier.Copy(localFriend, apiFriend)
	copier.Copy(localFriend, apiFriend.FriendUser)
	localFriend.FriendUserID = apiFriend.FriendUser.UserID
}

func friendRequestCopyToLocal(localFriendRequest *db.LocalFriendRequest, apiFriendRequest *server_api_params.FriendRequest) {
	copier.Copy(localFriendRequest, apiFriendRequest)
}

func blackCopyToLocal(localBlack *db.LocalBlack, apiBlack *server_api_params.PublicUserInfo, ownerUserID string) {
	copier.Copy(localBlack, apiBlack)
	localBlack.OwnerUserID = ownerUserID
	localBlack.BlockUserID = apiBlack.UserID
}

func TransferToLocalFriend(apiFriendList []*server_api_params.FriendInfo) []*db.LocalFriend {
	localFriendList := make([]*db.LocalFriend, 0)
	for _, v := range apiFriendList {
		var localFriend db.LocalFriend
		friendCopyToLocal(&localFriend, v)
		localFriendList = append(localFriendList, &localFriend)
	}
	return localFriendList
}

func checkListDiff(a []diff, b []diff) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]interface{})
	for _, v := range a {
		mapA[v.Key()] = v
	}
	mapB := make(map[string]interface{})
	for _, v := range b {
		mapB[v.Key()] = v
	}

	aInBNot = make([]int, 0)
	bInANot = make([]int, 0)
	sameA = make([]int, 0)
	sameB = make([]int, 0)

	//for a
	for i, v := range a {
		ia, ok := mapB[v.Key()]
		if !ok {
			//in a, but not in b
			aInBNot = append(aInBNot, i)
		} else {
			if v != ia {
				sameA = append(sameA, i)
			}
		}
	}
	//for b
	for i, v := range b {
		ib, ok := mapA[v.Key()]
		if !ok {
			bInANot = append(bInANot, i)
		} else {
			if ib != v {
				sameB = append(sameB, i)
			}
		}
	}
	return aInBNot, bInANot, sameA, sameB
}

func TransferToLocalGroupMember(apiData []*server_api_params.GroupMemberFullInfo) []*db.LocalGroupMember {
	local := make([]*db.LocalGroupMember, 0)
	for _, v := range apiData {
		var node db.LocalGroupMember
		log2.NewDebug("0", "local test api ", v)
		GroupMemberCopyToLocal(&node, v)
		log2.NewDebug("0", "local test local  ", node)
		local = append(local, &node)
	}
	log2.NewDebug("0", "local test local all ", local)
	return local
}

func GroupMemberCopyToLocal(dst *db.LocalGroupMember, src *server_api_params.GroupMemberFullInfo) {
	copier.Copy(dst, src)
}



func TransferToLocalGroupInfo(apiData []*server_api_params.GroupInfo) []*db.LocalGroup {
	local := make([]*db.LocalGroup, 0)
	for _, v := range apiData {
		var node db.LocalGroup
		log2.NewDebug("0", "local test api ", v)
		GroupInfoCopyToLocal(&node, v)
		log2.NewDebug("0", "local test local  ", node)
		local = append(local, &node)
	}
	log2.NewDebug("0", "local test local all ", local)
	return local
}

func GroupInfoCopyToLocal(dst *db.LocalGroup, src *server_api_params.GroupInfo) {
	copier.Copy(dst, src)
}



func TransferToLocalGroupRequest(apiData []*server_api_params.GroupRequest) []*db.LocalGroupRequest {
	local := make([]*db.LocalGroupRequest, 0)
	for _, v := range apiData {
		var node db.LocalGroupRequest
		log2.NewDebug("0", "local test api ", v)
		GroupRequestCopyToLocal(&node, v)
		log2.NewDebug("0", "local test local  ", node)
		local = append(local, &node)
	}
	log2.NewDebug("0", "local test local all ", local)
	return local
}

func GroupRequestCopyToLocal(dst *db.LocalGroupRequest, src *server_api_params.GroupRequest) {
	copier.Copy(dst, src)
}
//
//func TransferToLocalUserInfo(apiData []*server_api_params.UserInfo) []*db.LocalUser {
//	localData := make([]*db.LocalUser, 0)
//	for _, v := range apiData {
//		var localNode db.LocalUser
//		log2.NewDebug("0", "local test api ", v)
//		UserInfoCopyToLocal(&localNode, v)
//		log2.NewDebug("0", "local test local  ", localNode)
//		localData = append(localData, &localNode)
//	}
//	log2.NewDebug("0", "local test local all ", localData)
//	return localData
//}
//
//func UserInfoCopyToLocal(dst *db.LocalUser, src *server_api_params.UserInfo) {
//	copier.Copy(dst, src)
//}


func TransferToLocalUserInfo(apiData *server_api_params.UserInfo) *db.LocalUser{
	var localNode db.LocalUser
	copier.Copy(&localNode, apiData)
	return &localNode
}

func TransferToLocalFriendRequest(apiFriendList []*server_api_params.FriendRequest) []*db.LocalFriendRequest {
	localFriendList := make([]*db.LocalFriendRequest, 0)
	for _, v := range apiFriendList {
		var localFriendRequest db.LocalFriendRequest
		log2.NewDebug("0", "local test api ", v)
		friendRequestCopyToLocal(&localFriendRequest, v)
		log2.NewDebug("0", "local test local  ", localFriendRequest)
		localFriendList = append(localFriendList, &localFriendRequest)
	}
	log2.NewDebug("0", "local test local all ", localFriendList)
	return localFriendList
}

func TransferToLocalBlack(apiBlackList []*server_api_params.PublicUserInfo, ownerUserID string) []*db.LocalBlack {
	localBlackList := make([]*db.LocalBlack, 0)
	for _, v := range apiBlackList {
		var localBlack db.LocalBlack
		blackCopyToLocal(&localBlack, v, ownerUserID)
		localBlackList = append(localBlackList, &localBlack)
	}

	return localBlackList
}

func CheckFriendListDiff(a []*db.LocalFriend, b []*db.LocalFriend) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*db.LocalFriend)
	for _, v := range a {
		mapA[v.FriendUserID] = v
	}
	mapB := make(map[string]*db.LocalFriend)
	for _, v := range b {
		mapB[v.FriendUserID] = v
	}

	aInBNot = make([]int, 0)
	bInANot = make([]int, 0)
	sameA = make([]int, 0)
	sameB = make([]int, 0)

	//for a
	for i, v := range a {
		ia, ok := mapB[v.FriendUserID]
		if !ok {
			//in a, but not in b
			aInBNot = append(aInBNot, i)
		} else {
			if v != ia {
				// key of a and b is equal, but value different
				sameA = append(sameA, i)
			}
		}
	}
	//for b
	for i, v := range b {
		ib, ok := mapA[v.FriendUserID]
		if !ok {
			bInANot = append(bInANot, i)
		} else {
			if ib != v {
				sameB = append(sameB, i)
			}
		}
	}
	return aInBNot, bInANot, sameA, sameB
}

func CheckFriendRequestDiff(a []*db.LocalFriendRequest, b []*db.LocalFriendRequest) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*db.LocalFriendRequest)
	for _, v := range a {
		mapA[v.ToUserID] = v
	}
	mapB := make(map[string]*db.LocalFriendRequest)
	for _, v := range b {
		mapB[v.ToUserID] = v
	}

	aInBNot = make([]int, 0)
	bInANot = make([]int, 0)
	sameA = make([]int, 0)
	sameB = make([]int, 0)

	//for a
	for i, v := range a {
		ia, ok := mapB[v.ToUserID]
		if !ok {
			//in a, but not in b
			aInBNot = append(aInBNot, i)
		} else {
			if v != ia {
				// key of a and b is equal, but value different
				sameA = append(sameA, i)
			}
		}
	}
	//for b
	for i, v := range b {
		ib, ok := mapA[v.ToUserID]
		if !ok {
			bInANot = append(bInANot, i)
		} else {
			if ib != v {
				sameB = append(sameB, i)
			}
		}
	}
	return aInBNot, bInANot, sameA, sameB
}
func CheckBlackListDiff(a []*db.LocalBlack, b []*db.LocalBlack) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*db.LocalBlack)
	for _, v := range a {
		mapA[v.BlockUserID] = v
	}
	mapB := make(map[string]*db.LocalBlack)
	for _, v := range b {
		mapB[v.BlockUserID] = v
	}

	aInBNot = make([]int, 0)
	bInANot = make([]int, 0)
	sameA = make([]int, 0)
	sameB = make([]int, 0)

	//for a
	for i, v := range a {
		ia, ok := mapB[v.BlockUserID]
		if !ok {
			//in a, but not in b
			aInBNot = append(aInBNot, i)
		} else {
			if v == ia {
				// key of a and b is equal, but value different
				sameA = append(sameA, i)
			}
		}
	}
	//for b
	for i, v := range b {
		ib, ok := mapA[v.BlockUserID]
		if !ok {
			bInANot = append(bInANot, i)
		} else {
			if ib != v {
				sameB = append(sameB, i)
			}
		}
	}
	return aInBNot, bInANot, sameA, sameB
}



func CheckGroupInfoDiff(a []*db.LocalGroup, b []*db.LocalGroup) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*db.LocalGroup)
	for _, v := range a {
		mapA[v.GroupID] = v
	}
	mapB := make(map[string]*db.LocalGroup)
	for _, v := range b {
		mapB[v.GroupID] = v
	}

	aInBNot = make([]int, 0)
	bInANot = make([]int, 0)
	sameA = make([]int, 0)
	sameB = make([]int, 0)

	//for a
	for i, v := range a {
		ia, ok := mapB[v.GroupID]
		if !ok {
			//in a, but not in b
			aInBNot = append(aInBNot, i)
		} else {
			if v != ia {
				// key of a and b is equal, but value different
				sameA = append(sameA, i)
			}
		}
	}
	//for b
	for i, v := range b {
		ib, ok := mapA[v.GroupID]
		if !ok {
			bInANot = append(bInANot, i)
		} else {
			if ib != v {
				sameB = append(sameB, i)
			}
		}
	}
	return aInBNot, bInANot, sameA, sameB
}



func CheckGroupMemberDiff(a []*db.LocalGroupMember, b []*db.LocalGroupMember) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*db.LocalGroupMember)
	for _, v := range a {
		mapA[v.GroupID+v.UserID] = v
	}
	mapB := make(map[string]*db.LocalGroupMember)
	for _, v := range b {
		mapB[v.GroupID+v.UserID] = v
	}

	aInBNot = make([]int, 0)
	bInANot = make([]int, 0)
	sameA = make([]int, 0)
	sameB = make([]int, 0)

	//for a
	for i, v := range a {
		ia, ok := mapB[v.GroupID+v.UserID]
		if !ok {
			//in a, but not in b
			aInBNot = append(aInBNot, i)
		} else {
			if v != ia {
				// key of a and b is equal, but value different
				sameA = append(sameA, i)
			}
		}
	}
	//for b
	for i, v := range b {
		ib, ok := mapA[v.GroupID]
		if !ok {
			bInANot = append(bInANot, i)
		} else {
			if ib != v {
				sameB = append(sameB, i)
			}
		}
	}
	return aInBNot, bInANot, sameA, sameB
}




func CheckGroupRequestDiff(a []*db.LocalGroupRequest, b []*db.LocalGroupRequest) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*db.LocalGroupRequest)
	for _, v := range a {
		mapA[v.GroupID+v.UserID] = v
	}
	mapB := make(map[string]*db.LocalGroupRequest)
	for _, v := range b {
		mapB[v.GroupID+v.UserID] = v
	}

	aInBNot = make([]int, 0)
	bInANot = make([]int, 0)
	sameA = make([]int, 0)
	sameB = make([]int, 0)

	//for a
	for i, v := range a {
		ia, ok := mapB[v.GroupID+v.UserID]
		if !ok {
			//in a, but not in b
			aInBNot = append(aInBNot, i)
		} else {
			if v != ia {
				// key of a and b is equal, but value different
				sameA = append(sameA, i)
			}
		}
	}
	//for b
	for i, v := range b {
		ib, ok := mapA[v.GroupID+v.UserID]
		if !ok {
			bInANot = append(bInANot, i)
		} else {
			if ib != v {
				sameB = append(sameB, i)
			}
		}
	}
	return aInBNot, bInANot, sameA, sameB
}



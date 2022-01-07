package commom

import (
	"github.com/jinzhu/copier"
	"open_im_sdk/pkg/db"
	log2 "open_im_sdk/pkg/log"
	"reflect"
)

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

func friendCopyToLocal(localFriend *open_im_sdk.LocalFriend, apiFriend *open_im_sdk.FriendInfo) {
	copier.Copy(localFriend, apiFriend)
	copier.Copy(localFriend, apiFriend.FriendUser)
	localFriend.FriendUserID = apiFriend.FriendUser.UserID
}

func friendRequestCopyToLocal(localFriendRequest *open_im_sdk.LocalFriendRequest, apiFriendRequest *open_im_sdk.FriendRequest) {
	copier.Copy(localFriendRequest, apiFriendRequest)

}

func blackCopyToLocal(localBlack *open_im_sdk.LocalBlack, apiBlack *open_im_sdk.PublicUserInfo, ownerUserID string) {
	copier.Copy(localBlack, apiBlack)
	localBlack.OwnerUserID = ownerUserID
	localBlack.BlockUserID = apiBlack.UserID
}

func transferToLocalFriend(apiFriendList []*open_im_sdk.FriendInfo) []*open_im_sdk.LocalFriend {
	localFriendList := make([]*db.LocalFriend, 0)
	for _, v := range apiFriendList {
		var localFriend db.LocalFriend
		friendCopyToLocal(&localFriend, v)
		localFriendList = append(localFriendList, &localFriend)
	}
	return localFriendList
}

func checkListDiff(a []open_im_sdk.diff, b []open_im_sdk.diff) (aInBNot, bInANot, sameA, sameB []int) {
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

func transferToLocalFriendRequest(apiFriendList []*open_im_sdk.FriendRequest) []*open_im_sdk.LocalFriendRequest {
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

func transferToLocalBlack(apiBlackList []*open_im_sdk.PublicUserInfo, ownerUserID string) []*open_im_sdk.LocalBlack {
	localBlackList := make([]*db.LocalBlack, 0)
	for _, v := range apiBlackList {
		var localBlack db.LocalBlack
		blackCopyToLocal(&localBlack, v, ownerUserID)
		localBlackList = append(localBlackList, &localBlack)
	}

	return localBlackList
}

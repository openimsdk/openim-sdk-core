package common

import (
	"fmt"
	"open_im_sdk/pkg/db"

	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/copier"

	//log2 "open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
)

type diff interface {
	Key() string
	Value() interface{}
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
			if !cmp.Equal(v, ia) {
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
	//operationID := utils.OperationIDGenerator()
	for _, v := range apiData {
		var node db.LocalGroupMember
		//		log2.NewDebug(operationID, "local test api ", v)
		GroupMemberCopyToLocal(&node, v)
		//		log2.NewDebug(operationID, "local test local  ", node)
		local = append(local, &node)
	}
	//	log2.NewDebug(operationID, "local test local all ", local)
	return local
}

func GroupMemberCopyToLocal(dst *db.LocalGroupMember, src *server_api_params.GroupMemberFullInfo) {
	copier.Copy(dst, src)
}

func TransferToLocalGroupInfo(apiData []*server_api_params.GroupInfo) []*db.LocalGroup {
	local := make([]*db.LocalGroup, 0)
	//operationID := utils.OperationIDGenerator()
	for _, v := range apiData {
		var node db.LocalGroup
		//	log2.NewDebug(operationID, "local test api ", v)
		GroupInfoCopyToLocal(&node, v)
		//		log2.NewDebug(operationID, "local test local  ", node)
		local = append(local, &node)
	}
	//	log2.NewDebug(operationID, "local test local all ", local)
	return local
}

func GroupInfoCopyToLocal(dst *db.LocalGroup, src *server_api_params.GroupInfo) {
	copier.Copy(dst, src)
}

func TransferToLocalGroupRequest(apiData []*server_api_params.GroupRequest) []*db.LocalGroupRequest {
	local := make([]*db.LocalGroupRequest, 0)
	//operationID := utils.OperationIDGenerator()
	for _, v := range apiData {
		var node db.LocalGroupRequest
		//	log2.NewDebug(operationID, "local test api ", v)
		GroupRequestCopyToLocal(&node, v)
		//		log2.NewDebug(operationID, "local test local  ", node)
		local = append(local, &node)
	}
	//	log2.NewDebug(operationID, "local test local all ", local)
	return local
}

func GroupRequestCopyToLocal(dst *db.LocalGroupRequest, src *server_api_params.GroupRequest) {
	copier.Copy(dst, src)
	copier.Copy(dst, src.GroupInfo)
	copier.Copy(dst, src.UserInfo)
	dst.GroupFaceURL = src.GroupInfo.FaceURL
	dst.UserFaceURL = src.UserInfo.FaceURL
}

func AdminGroupRequestCopyToLocal(dst *db.LocalAdminGroupRequest, src *server_api_params.GroupRequest) {
	copier.Copy(dst, src)
	copier.Copy(dst, src.GroupInfo)
	copier.Copy(dst, src.UserInfo)
	dst.GroupFaceURL = src.GroupInfo.FaceURL
	dst.UserFaceURL = src.UserInfo.FaceURL
}

func SendGroupRequestCopyToLocal(dst *db.LocalGroupRequest, src *server_api_params.GroupRequest) {
	copier.Copy(dst, src)
	copier.Copy(dst, src.GroupInfo)
	copier.Copy(dst, src.UserInfo)
	dst.GroupFaceURL = src.GroupInfo.FaceURL
	dst.UserFaceURL = src.UserInfo.FaceURL
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

func TransferToLocalUserInfo(apiData *server_api_params.UserInfo) *db.LocalUser {
	var localNode db.LocalUser
	copier.Copy(&localNode, apiData)
	return &localNode
}

func TransferToLocalFriendRequest(apiFriendList []*server_api_params.FriendRequest) []*db.LocalFriendRequest {
	localFriendList := make([]*db.LocalFriendRequest, 0)
	//operationID := utils.OperationIDGenerator()
	for _, v := range apiFriendList {
		var localFriendRequest db.LocalFriendRequest
		//	log2.NewDebug(operationID, "local test api ", v)
		friendRequestCopyToLocal(&localFriendRequest, v)
		//	log2.NewDebug(operationID, "local test local  ", localFriendRequest)
		localFriendList = append(localFriendList, &localFriendRequest)
	}
	//	log2.NewDebug(operationID, "local test local all ", localFriendList)
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
			if !cmp.Equal(v, ia) {
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
			if !cmp.Equal(v, ib) {
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
		mapA[v.FromUserID+v.ToUserID] = v
	}
	mapB := make(map[string]*db.LocalFriendRequest)
	for _, v := range b {
		mapB[v.FromUserID+v.ToUserID] = v
	}

	aInBNot = make([]int, 0)
	bInANot = make([]int, 0)
	sameA = make([]int, 0)
	sameB = make([]int, 0)

	//for a
	for i, v := range a {
		ia, ok := mapB[v.FromUserID+v.ToUserID]
		if !ok {
			//in a, but not in b
			aInBNot = append(aInBNot, i)
		} else {
			if !cmp.Equal(v, ia) {
				// key of a and b is equal, but value different
				sameA = append(sameA, i)
			}
		}
	}
	//for b
	for i, v := range b {
		ib, ok := mapA[v.FromUserID+v.ToUserID]
		if !ok {
			bInANot = append(bInANot, i)
		} else {
			if !cmp.Equal(v, ib) {
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
			if !cmp.Equal(v, ia) {
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
			if !cmp.Equal(v, ib) {
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
		//fmt.Println("mapa   ", *v)
		mapA[v.GroupID] = v

	}
	mapB := make(map[string]*db.LocalGroup)
	for _, v := range b {
		//	fmt.Println("mapb   ", *v)
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
			if !cmp.Equal(v, ia) {
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
			if !cmp.Equal(v, ib) {
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
			//reflect.DeepEqual(a, b)
			//	reflect.DeepEqual(v, ia)
			//if !cmp.Equal(v, ia)
			if !cmp.Equal(v, ia) {
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
			if !cmp.Equal(v, ib) {
				sameB = append(sameB, i)
			}
		}
	}
	return aInBNot, bInANot, sameA, sameB
}

func CheckDepartmentMemberDiff(a []*db.LocalDepartmentMember, b []*db.LocalDepartmentMember) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*db.LocalDepartmentMember)
	for _, v := range a {
		mapA[v.DepartmentID+v.UserID] = v
	}
	mapB := make(map[string]*db.LocalDepartmentMember)
	for _, v := range b {
		mapB[v.DepartmentID+v.UserID] = v
	}

	aInBNot = make([]int, 0)
	bInANot = make([]int, 0)
	sameA = make([]int, 0)
	sameB = make([]int, 0)

	//for a
	for i, v := range a {
		ia, ok := mapB[v.DepartmentID+v.UserID]
		if !ok {
			//in a, but not in b
			aInBNot = append(aInBNot, i)
		} else {
			//reflect.DeepEqual(a, b)
			//	reflect.DeepEqual(v, ia)
			//if !cmp.Equal(v, ia)
			if !cmp.Equal(v, ia) {
				// key of a and b is equal, but value different
				sameA = append(sameA, i)
			}
		}
	}
	//for b
	for i, v := range b {
		ib, ok := mapA[v.DepartmentID+v.UserID]
		if !ok {
			bInANot = append(bInANot, i)
		} else {
			if !cmp.Equal(v, ib) {
				sameB = append(sameB, i)
			}
		}
	}
	return aInBNot, bInANot, sameA, sameB

}

func CheckDepartmentDiff(a []*db.LocalDepartment, b []*db.LocalDepartment) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*db.LocalDepartment)
	for _, v := range a {
		mapA[v.DepartmentID] = v
	}
	mapB := make(map[string]*db.LocalDepartment)
	for _, v := range b {
		mapB[v.DepartmentID] = v
	}

	aInBNot = make([]int, 0)
	bInANot = make([]int, 0)
	sameA = make([]int, 0)
	sameB = make([]int, 0)

	//for a
	for i, v := range a {
		ia, ok := mapB[v.DepartmentID]
		if !ok {
			//in a, but not in b
			aInBNot = append(aInBNot, i)
		} else {
			if !cmp.Equal(v, ia) {
				// key of a and b is equal, but value different
				sameA = append(sameA, i)
			}
		}
	}
	//for b
	for i, v := range b {
		ib, ok := mapA[v.DepartmentID]
		if !ok {
			bInANot = append(bInANot, i)
		} else {
			if !cmp.Equal(v, ib) {
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
			if !cmp.Equal(v, ia) {
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
			if !cmp.Equal(v, ib) {
				sameB = append(sameB, i)
			}
		}
	}
	return aInBNot, bInANot, sameA, sameB
}

func CheckAdminGroupRequestDiff(a []*db.LocalAdminGroupRequest, b []*db.LocalAdminGroupRequest) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*db.LocalAdminGroupRequest)
	for _, v := range a {
		mapA[v.GroupID+v.UserID] = v
		//	fmt.Println("mapA   ", v)
	}
	mapB := make(map[string]*db.LocalAdminGroupRequest)
	for _, v := range b {
		mapB[v.GroupID+v.UserID] = v
		//	fmt.Println("mapB   ", v)
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
			//fmt.Println("aInBNot", a[i], ia)
			aInBNot = append(aInBNot, i)
		} else {
			if !cmp.Equal(v, ia) {
				// key of a and b is equal, but value different
				//fmt.Println("sameA", a[i], ia)
				sameA = append(sameA, i)
			}
		}
	}
	//for b
	for i, v := range b {
		ib, ok := mapA[v.GroupID+v.UserID]
		if !ok {
			//fmt.Println("bInANot", b[i], ib)

			bInANot = append(bInANot, i)
		} else {
			if !cmp.Equal(v, ib) {
				//	fmt.Println("sameB", b[i], ib)
				sameB = append(sameB, i)
			}
		}
	}
	return aInBNot, bInANot, sameA, sameB
}

func CheckConversationListDiff(conversationsOnServer, conversationsOnLocal []*tempConversation) (aInBNot, bInANot, sameA, sameB []int) {
	mapA := make(map[string]*tempConversation)
	mapB := make(map[string]*tempConversation)
	for _, v := range conversationsOnServer {
		mapA[v.ConversationID] = v
	}
	for _, v := range conversationsOnLocal {
		mapB[v.ConversationID] = v
	}
	aInBNot = make([]int, 0)
	bInANot = make([]int, 0)
	sameA = make([]int, 0)
	sameB = make([]int, 0)
	for i, v := range conversationsOnServer {
		ia, ok := mapB[v.ConversationID]
		if !ok {
			//in a, but not in b
			//fmt.Println("aInBNot", conversationsOnServer[i], ia)
			aInBNot = append(aInBNot, i)
		} else {
			//fmt.Println("test result is v", v)
			//fmt.Println("test result is ia", ia)
			if !cmp.Equal(v, ia) {
				fmt.Println(v, ia)
				// key of a and b is equal, but value different
				//fmt.Println("sameA", conversationsOnServer[i], ia)
				sameA = append(sameA, i)
			}
		}
	}

	for i, v := range conversationsOnLocal {
		ib, ok := mapA[v.ConversationID]
		if !ok {
			//fmt.Println("bInANot", conversationsOnLocal[i], ib)
			bInANot = append(bInANot, i)
		} else {
			if !cmp.Equal(v, ib) {
				//	fmt.Println("sameB", conversationsOnLocal[i], ib)
				sameB = append(sameB, i)
			}
		}
	}

	return aInBNot, bInANot, sameA, sameB
}

//
//func CheckSendGroupRequestDiff(a []*db.LocalGroupRequest, b []*db.LocalGroupRequest) (aInBNot, bInANot, sameA, sameB []int) {
//	//to map, friendid_>friendinfo
//	mapA := make(map[string]*db.LocalGroupRequest)
//	for _, v := range a {
//		mapA[v.GroupID+v.UserID] = v
//		fmt.Println("mapA   ", v)
//	}
//	mapB := make(map[string]*db.LocalGroupRequest)
//	for _, v := range b {
//		mapB[v.GroupID+v.UserID] = v
//		fmt.Println("mapB   ", v)
//	}
//
//	aInBNot = make([]int, 0)
//	bInANot = make([]int, 0)
//	sameA = make([]int, 0)
//	sameB = make([]int, 0)
//
//	//for a
//	for i, v := range a {
//		ia, ok := mapB[v.GroupID+v.UserID]
//		if !ok {
//			//in a, but not in b
//			fmt.Println("aInBNot", a[i], ia)
//			aInBNot = append(aInBNot, i)
//		} else {
//			if !cmp.Equal(v, ia) {
//				// key of a and b is equal, but value different
//				fmt.Println("sameA", a[i], ia)
//				sameA = append(sameA, i)
//			}
//		}
//	}
//	//for b
//	for i, v := range b {
//		ib, ok := mapA[v.GroupID+v.UserID]
//		if !ok {
//			fmt.Println("bInANot", b[i], ib)
//
//			bInANot = append(bInANot, i)
//		} else {
//			if !cmp.Equal(v, ib) {
//				fmt.Println("sameB", b[i], ib)
//				sameB = append(sameB, i)
//			}
//		}
//	}
//	return aInBNot, bInANot, sameA, sameB
//}

func TransferToLocalAdminGroupRequest(apiData []*server_api_params.GroupRequest) []*db.LocalAdminGroupRequest {
	local := make([]*db.LocalAdminGroupRequest, 0)
	//operationID := utils.OperationIDGenerator()
	for _, v := range apiData {
		var node db.LocalAdminGroupRequest
		//	log2.NewDebug(operationID, "local test api ", v)
		AdminGroupRequestCopyToLocal(&node, v)
		//		log2.NewDebug(operationID, "local test local  ", node)
		local = append(local, &node)
	}
	//	log2.NewDebug(operationID, "local test local all ", local)
	return local
}

func TransferToLocalDepartmentMember(apiData []*server_api_params.UserDepartmentMember) []*db.LocalDepartmentMember {
	local := make([]*db.LocalDepartmentMember, 0)
	for _, v := range apiData {
		var node db.LocalDepartmentMember
		copier.Copy(&node, v.DepartmentMember)
		copier.Copy(&node, v.OrganizationUser)
		local = append(local, &node)
	}
	return local
}

func TransferToLocalDepartment(apiData []*server_api_params.Department) []*db.LocalDepartment {
	local := make([]*db.LocalDepartment, 0)
	for _, v := range apiData {
		var node db.LocalDepartment
		copier.Copy(&node, v)
		local = append(local, &node)
	}
	return local
}

func TransferToLocalSendGroupRequest(apiData []*server_api_params.GroupRequest) []*db.LocalGroupRequest {
	local := make([]*db.LocalGroupRequest, 0)
	for _, v := range apiData {
		var node db.LocalGroupRequest
		SendGroupRequestCopyToLocal(&node, v)
		local = append(local, &node)
	}
	return local
}

type tempConversation struct {
	RecvMsgOpt       int32
	ConversationID   string
	ConversationType int32
	UserID           string
	GroupID          string
	IsPrivateChat    bool
	IsPinned         bool
	GroupAtType      int32
	IsNotInGroup     bool
	AttachedInfo     string
	Ex               string
}

func ServerTransferToTempConversation(resp server_api_params.GetAllConversationsResp) []*tempConversation {
	var tempConversations []*tempConversation
	for _, serverConversation := range resp.Conversations {
		tempConversations = append(tempConversations, &tempConversation{
			RecvMsgOpt:       serverConversation.RecvMsgOpt,
			ConversationID:   serverConversation.ConversationID,
			ConversationType: serverConversation.ConversationType,
			UserID:           serverConversation.UserID,
			GroupID:          serverConversation.GroupID,
			IsPrivateChat:    serverConversation.IsPrivateChat,
			IsPinned:         serverConversation.IsPinned,
			GroupAtType:      serverConversation.GroupAtType,
			IsNotInGroup:     serverConversation.IsNotInGroup,
			AttachedInfo:     serverConversation.AttachedInfo,
			Ex:               serverConversation.Ex,
		})
	}
	return tempConversations
}

func LocalTransferToTempConversation(local []*db.LocalConversation) []*tempConversation {
	var tempConversations []*tempConversation
	for _, localConversation := range local {
		tempConversations = append(tempConversations, &tempConversation{
			RecvMsgOpt:       localConversation.RecvMsgOpt,
			ConversationID:   localConversation.ConversationID,
			ConversationType: localConversation.ConversationType,
			UserID:           localConversation.UserID,
			GroupID:          localConversation.GroupID,
			IsPrivateChat:    localConversation.IsPrivateChat,
			IsPinned:         localConversation.IsPinned,
			GroupAtType:      localConversation.GroupAtType,
			IsNotInGroup:     localConversation.IsNotInGroup,
			AttachedInfo:     localConversation.AttachedInfo,
			Ex:               localConversation.Ex,
		})
	}
	return tempConversations
}

func TransferToLocalConversation(resp server_api_params.GetAllConversationsResp) []*db.LocalConversation {
	var localConversations []*db.LocalConversation
	for _, serverConversation := range resp.Conversations {
		localConversations = append(localConversations, &db.LocalConversation{
			RecvMsgOpt:       serverConversation.RecvMsgOpt,
			ConversationID:   serverConversation.ConversationID,
			ConversationType: serverConversation.ConversationType,
			UserID:           serverConversation.UserID,
			GroupID:          serverConversation.GroupID,
			IsPrivateChat:    serverConversation.IsPrivateChat,
			IsPinned:         serverConversation.IsPinned,
			GroupAtType:      serverConversation.GroupAtType,
			IsNotInGroup:     serverConversation.IsNotInGroup,
			AttachedInfo:     serverConversation.AttachedInfo,
			Ex:               serverConversation.Ex,
		})
	}
	return localConversations
}

func TransferToServerConversation(local []*db.LocalConversation) server_api_params.GetAllConversationsResp {
	var serverConversations server_api_params.GetAllConversationsResp
	for _, localConversation := range local {
		serverConversations.Conversations = append(serverConversations.Conversations, server_api_params.Conversation{
			RecvMsgOpt:       localConversation.RecvMsgOpt,
			ConversationID:   localConversation.ConversationID,
			ConversationType: localConversation.ConversationType,
			UserID:           localConversation.UserID,
			GroupID:          localConversation.GroupID,
			IsPrivateChat:    localConversation.IsPrivateChat,
			IsPinned:         localConversation.IsPinned,
			AttachedInfo:     localConversation.AttachedInfo,
			Ex:               localConversation.Ex,
		})
	}
	return serverConversations
}

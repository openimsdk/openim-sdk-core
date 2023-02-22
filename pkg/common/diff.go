package common

import (
	"fmt"
	"open_im_sdk/pkg/db/model_struct"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/copier"

	//log2 "open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
)

type diff interface {
	Key() string
	Value() interface{}
}

func friendCopyToLocal(localFriend *model_struct.LocalFriend, apiFriend *server_api_params.FriendInfo) {
	copier.Copy(localFriend, apiFriend)
	copier.Copy(localFriend, apiFriend.FriendUser)
	localFriend.FriendUserID = apiFriend.FriendUser.UserID
}

func friendRequestCopyToLocal(localFriendRequest *model_struct.LocalFriendRequest, apiFriendRequest *server_api_params.FriendRequest) {
	copier.Copy(localFriendRequest, apiFriendRequest)
}

func blackCopyToLocal(localBlack *model_struct.LocalBlack, apiBlack *server_api_params.PublicUserInfo, ownerUserID string) {
	copier.Copy(localBlack, apiBlack)
	localBlack.OwnerUserID = ownerUserID
	localBlack.BlockUserID = apiBlack.UserID
}

func TransferToLocalFriend(apiFriendList []*server_api_params.FriendInfo) []*model_struct.LocalFriend {
	localFriendList := make([]*model_struct.LocalFriend, 0)
	for _, v := range apiFriendList {
		var localFriend model_struct.LocalFriend
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

func TransferToLocalGroupMember(apiData []*server_api_params.GroupMemberFullInfo) []*model_struct.LocalGroupMember {
	local := make([]*model_struct.LocalGroupMember, 0)
	//operationID := utils.OperationIDGenerator()
	for _, v := range apiData {
		var node model_struct.LocalGroupMember
		//		log2.NewDebug(operationID, "local test api ", v)
		GroupMemberCopyToLocal(&node, v)
		//		log2.NewDebug(operationID, "local test local  ", node)
		local = append(local, &node)
	}
	//	log2.NewDebug(operationID, "local test local all ", local)
	return local
}

func GroupMemberCopyToLocal(dst *model_struct.LocalGroupMember, src *server_api_params.GroupMemberFullInfo) {
	copier.Copy(dst, src)
}

func TransferToLocalGroupInfo(apiData []*server_api_params.GroupInfo) []*model_struct.LocalGroup {
	local := make([]*model_struct.LocalGroup, 0)
	//operationID := utils.OperationIDGenerator()
	for _, v := range apiData {
		var node model_struct.LocalGroup
		//	log2.NewDebug(operationID, "local test api ", v)
		GroupInfoCopyToLocal(&node, v)
		//		log2.NewDebug(operationID, "local test local  ", node)
		local = append(local, &node)
	}
	//	log2.NewDebug(operationID, "local test local all ", local)
	return local
}

func GroupInfoCopyToLocal(dst *model_struct.LocalGroup, src *server_api_params.GroupInfo) {
	copier.Copy(dst, src)
}

func TransferToLocalGroupRequest(apiData []*server_api_params.GroupRequest) []*model_struct.LocalGroupRequest {
	local := make([]*model_struct.LocalGroupRequest, 0)
	//operationID := utils.OperationIDGenerator()
	for _, v := range apiData {
		var node model_struct.LocalGroupRequest
		//	log2.NewDebug(operationID, "local test api ", v)
		GroupRequestCopyToLocal(&node, v)
		//		log2.NewDebug(operationID, "local test local  ", node)
		local = append(local, &node)
	}
	//	log2.NewDebug(operationID, "local test local all ", local)
	return local
}

func GroupRequestCopyToLocal(dst *model_struct.LocalGroupRequest, src *server_api_params.GroupRequest) {
	copier.Copy(dst, src)
	copier.Copy(dst, src.GroupInfo)
	copier.Copy(dst, src.UserInfo)
	dst.GroupFaceURL = src.GroupInfo.FaceURL
	dst.UserFaceURL = src.UserInfo.FaceURL
}

func AdminGroupRequestCopyToLocal(dst *model_struct.LocalAdminGroupRequest, src *server_api_params.GroupRequest) {
	copier.Copy(dst, src)
	copier.Copy(dst, src.GroupInfo)
	copier.Copy(dst, src.UserInfo)
	dst.GroupFaceURL = src.GroupInfo.FaceURL
	dst.UserFaceURL = src.UserInfo.FaceURL
}

func SendGroupRequestCopyToLocal(dst *model_struct.LocalGroupRequest, src *server_api_params.GroupRequest) {
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

func TransferToLocalUserInfo(apiData *server_api_params.UserInfo) *model_struct.LocalUser {
	var localNode model_struct.LocalUser
	copier.Copy(&localNode, apiData)
	t, _ := time.Parse("2006-01-02", apiData.BirthStr)
	localNode.BirthTime = t
	localNode.AppMangerLevel = 0
	return &localNode
}

func TransferToLocalFriendRequest(apiFriendList []*server_api_params.FriendRequest) []*model_struct.LocalFriendRequest {
	localFriendList := make([]*model_struct.LocalFriendRequest, 0)
	//operationID := utils.OperationIDGenerator()
	for _, v := range apiFriendList {
		var localFriendRequest model_struct.LocalFriendRequest
		//	log2.NewDebug(operationID, "local test api ", v)
		friendRequestCopyToLocal(&localFriendRequest, v)
		//	log2.NewDebug(operationID, "local test local  ", localFriendRequest)
		localFriendList = append(localFriendList, &localFriendRequest)
	}
	//	log2.NewDebug(operationID, "local test local all ", localFriendList)
	return localFriendList
}

func TransferToLocalBlack(apiBlackList []*server_api_params.PublicUserInfo, ownerUserID string) []*model_struct.LocalBlack {
	localBlackList := make([]*model_struct.LocalBlack, 0)
	for _, v := range apiBlackList {
		var localBlack model_struct.LocalBlack
		blackCopyToLocal(&localBlack, v, ownerUserID)
		localBlackList = append(localBlackList, &localBlack)
	}

	return localBlackList
}

func CheckFriendListDiff(a []*model_struct.LocalFriend, b []*model_struct.LocalFriend) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*model_struct.LocalFriend)
	for _, v := range a {
		mapA[v.FriendUserID] = v
	}
	mapB := make(map[string]*model_struct.LocalFriend)
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

func CheckFriendRequestDiff(a []*model_struct.LocalFriendRequest, b []*model_struct.LocalFriendRequest) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*model_struct.LocalFriendRequest)
	for _, v := range a {
		mapA[v.FromUserID+v.ToUserID] = v
	}
	mapB := make(map[string]*model_struct.LocalFriendRequest)
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
func CheckBlackListDiff(a []*model_struct.LocalBlack, b []*model_struct.LocalBlack) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*model_struct.LocalBlack)
	for _, v := range a {
		mapA[v.BlockUserID] = v
	}
	mapB := make(map[string]*model_struct.LocalBlack)
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

func CheckGroupInfoDiff(a []*model_struct.LocalGroup, b []*model_struct.LocalGroup) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*model_struct.LocalGroup)
	for _, v := range a {
		//fmt.Println("mapa   ", *v)
		mapA[v.GroupID] = v

	}
	mapB := make(map[string]*model_struct.LocalGroup)
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

func CheckGroupMemberDiff(a []*model_struct.LocalGroupMember, b []*model_struct.LocalGroupMember) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*model_struct.LocalGroupMember)
	for _, v := range a {
		mapA[v.GroupID+v.UserID] = v
	}
	mapB := make(map[string]*model_struct.LocalGroupMember)
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

func CheckDepartmentMemberDiff(a []*model_struct.LocalDepartmentMember, b []*model_struct.LocalDepartmentMember) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*model_struct.LocalDepartmentMember)
	for _, v := range a {
		mapA[v.DepartmentID+v.UserID] = v
	}
	mapB := make(map[string]*model_struct.LocalDepartmentMember)
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

func CheckDepartmentDiff(a []*model_struct.LocalDepartment, b []*model_struct.LocalDepartment) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*model_struct.LocalDepartment)
	for _, v := range a {
		mapA[v.DepartmentID] = v
	}
	mapB := make(map[string]*model_struct.LocalDepartment)
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

func CheckGroupRequestDiff(a []*model_struct.LocalGroupRequest, b []*model_struct.LocalGroupRequest) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*model_struct.LocalGroupRequest)
	for _, v := range a {
		mapA[v.GroupID+v.UserID] = v
	}
	mapB := make(map[string]*model_struct.LocalGroupRequest)
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

func CheckAdminGroupRequestDiff(a []*model_struct.LocalAdminGroupRequest, b []*model_struct.LocalAdminGroupRequest) (aInBNot, bInANot, sameA, sameB []int) {
	//to map, friendid_>friendinfo
	mapA := make(map[string]*model_struct.LocalAdminGroupRequest)
	for _, v := range a {
		mapA[v.GroupID+v.UserID] = v
		//	fmt.Println("mapA   ", v)
	}
	mapB := make(map[string]*model_struct.LocalAdminGroupRequest)
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
				//fmt.Println(v, ia)
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
func CheckReactionExtensionsDiff(onServer, onLocal []*server_api_params.SingleMessageExtensionResult) (aInBNot, bInANot, sameA, sameB []*server_api_params.SingleMessageExtensionResult) {
	mapA := make(map[string]*server_api_params.SingleMessageExtensionResult)
	mapB := make(map[string]*server_api_params.SingleMessageExtensionResult)
	for _, v := range onServer {
		mapA[v.ClientMsgID] = v
	}
	for _, v := range onLocal {
		mapB[v.ClientMsgID] = v
	}
	aInBNot = make([]*server_api_params.SingleMessageExtensionResult, 0)
	bInANot = make([]*server_api_params.SingleMessageExtensionResult, 0)
	sameA = make([]*server_api_params.SingleMessageExtensionResult, 0)
	sameB = make([]*server_api_params.SingleMessageExtensionResult, 0)
	for _, v := range onServer {
		ia, ok := mapB[v.ClientMsgID]
		if !ok {
			//in a, but not in b
			//fmt.Println("aInBNot", conversationsOnServer[i], ia)
			aInBNot = append(aInBNot, v)
		} else {
			//fmt.Println("test result is v", v)
			//fmt.Println("test result is ia", ia)
			if !cmp.Equal(v.ReactionExtensionList, ia.ReactionExtensionList) {
				fmt.Println(v, ia)
				// key of a and b is equal, but value different
				//fmt.Println("sameA", conversationsOnServer[i], ia)
				sameA = append(sameA, v)
			}
		}
	}

	for _, v := range onLocal {
		ib, ok := mapA[v.ClientMsgID]
		if !ok {
			//fmt.Println("bInANot", conversationsOnLocal[i], ib)
			bInANot = append(bInANot, v)
		} else {
			if !cmp.Equal(v.ReactionExtensionList, ib.ReactionExtensionList) {
				//	fmt.Println("sameB", conversationsOnLocal[i], ib)
				sameB = append(sameB, v)
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

func TransferToLocalAdminGroupRequest(apiData []*server_api_params.GroupRequest) []*model_struct.LocalAdminGroupRequest {
	local := make([]*model_struct.LocalAdminGroupRequest, 0)
	//operationID := utils.OperationIDGenerator()
	for _, v := range apiData {
		var node model_struct.LocalAdminGroupRequest
		//	log2.NewDebug(operationID, "local test api ", v)
		AdminGroupRequestCopyToLocal(&node, v)
		//		log2.NewDebug(operationID, "local test local  ", node)
		local = append(local, &node)
	}
	//	log2.NewDebug(operationID, "local test local all ", local)
	return local
}

func TransferToLocalDepartmentMember(apiData []*server_api_params.UserDepartmentMember) []*model_struct.LocalDepartmentMember {
	local := make([]*model_struct.LocalDepartmentMember, 0)
	for _, v := range apiData {
		var node model_struct.LocalDepartmentMember
		copier.Copy(&node, v.DepartmentMember)
		copier.Copy(&node, v.OrganizationUser)
		local = append(local, &node)
	}
	return local
}

func TransferToLocalDepartment(apiData []*server_api_params.Department) []*model_struct.LocalDepartment {
	local := make([]*model_struct.LocalDepartment, 0)
	for _, v := range apiData {
		var node model_struct.LocalDepartment
		copier.Copy(&node, v)
		local = append(local, &node)
	}
	return local
}

func TransferToLocalSendGroupRequest(apiData []*server_api_params.GroupRequest) []*model_struct.LocalGroupRequest {
	local := make([]*model_struct.LocalGroupRequest, 0)
	for _, v := range apiData {
		var node model_struct.LocalGroupRequest
		SendGroupRequestCopyToLocal(&node, v)
		local = append(local, &node)
	}
	return local
}

type tempConversation struct {
	RecvMsgOpt            int32
	ConversationID        string
	ConversationType      int32
	UserID                string
	GroupID               string
	IsPrivateChat         bool
	BurnDuration          int32
	IsPinned              bool
	UnreadCount           int32
	GroupAtType           int32
	IsNotInGroup          bool
	UpdateUnreadCountTime int64
	AttachedInfo          string
	Ex                    string
}

func ServerTransferToTempConversation(resp server_api_params.GetAllConversationsResp) []*tempConversation {
	var tempConversations []*tempConversation
	for _, serverConversation := range resp.Conversations {
		tempConversations = append(tempConversations, &tempConversation{
			RecvMsgOpt:            serverConversation.RecvMsgOpt,
			ConversationID:        serverConversation.ConversationID,
			ConversationType:      serverConversation.ConversationType,
			UserID:                serverConversation.UserID,
			GroupID:               serverConversation.GroupID,
			IsPrivateChat:         serverConversation.IsPrivateChat,
			BurnDuration:          serverConversation.BurnDuration,
			IsPinned:              serverConversation.IsPinned,
			GroupAtType:           serverConversation.GroupAtType,
			IsNotInGroup:          serverConversation.IsNotInGroup,
			AttachedInfo:          serverConversation.AttachedInfo,
			UnreadCount:           serverConversation.UnreadCount,
			UpdateUnreadCountTime: serverConversation.UpdateUnreadCountTime,
			Ex:                    serverConversation.Ex,
		})
	}
	return tempConversations
}

func LocalTransferToTempConversation(local []*model_struct.LocalConversation) []*tempConversation {
	var tempConversations []*tempConversation
	for _, localConversation := range local {
		tempConversations = append(tempConversations, &tempConversation{
			RecvMsgOpt:            localConversation.RecvMsgOpt,
			ConversationID:        localConversation.ConversationID,
			ConversationType:      localConversation.ConversationType,
			UserID:                localConversation.UserID,
			GroupID:               localConversation.GroupID,
			IsPrivateChat:         localConversation.IsPrivateChat,
			BurnDuration:          localConversation.BurnDuration,
			IsPinned:              localConversation.IsPinned,
			GroupAtType:           localConversation.GroupAtType,
			IsNotInGroup:          localConversation.IsNotInGroup,
			AttachedInfo:          localConversation.AttachedInfo,
			UnreadCount:           localConversation.UnreadCount,
			UpdateUnreadCountTime: localConversation.UpdateUnreadCountTime,
			Ex:                    localConversation.Ex,
		})
	}
	return tempConversations
}

func TransferToLocalConversation(resp server_api_params.GetAllConversationsResp) []*model_struct.LocalConversation {
	var localConversations []*model_struct.LocalConversation
	for _, serverConversation := range resp.Conversations {
		localConversations = append(localConversations, &model_struct.LocalConversation{
			RecvMsgOpt:            serverConversation.RecvMsgOpt,
			ConversationID:        serverConversation.ConversationID,
			ConversationType:      serverConversation.ConversationType,
			UserID:                serverConversation.UserID,
			GroupID:               serverConversation.GroupID,
			IsPrivateChat:         serverConversation.IsPrivateChat,
			BurnDuration:          serverConversation.BurnDuration,
			IsPinned:              serverConversation.IsPinned,
			GroupAtType:           serverConversation.GroupAtType,
			IsNotInGroup:          serverConversation.IsNotInGroup,
			AttachedInfo:          serverConversation.AttachedInfo,
			UnreadCount:           serverConversation.UnreadCount,
			UpdateUnreadCountTime: serverConversation.UpdateUnreadCountTime,
			Ex:                    serverConversation.Ex,
		})
	}
	return localConversations
}

func TransferToServerConversation(local []*model_struct.LocalConversation) server_api_params.GetAllConversationsResp {
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

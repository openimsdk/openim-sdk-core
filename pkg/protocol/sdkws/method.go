package sdkws

//import (
//	"go/constant"
//	"sort"
//	"strings"
//)
//
//func (x *RequestPagination) GetPage() (pageNumber, showNumber int32) {
//	if x != nil {
//		return x.PageNumber, x.ShowNumber
//	}
//	return
//}
//
//func (x *MsgData) GetNotificationConversationID() string {
//	switch x.SessionType {
//	case constant.SingleChatType:
//		l := []string{x.SendID, x.RecvID}
//		sort.Strings(l)
//		return "n_" + strings.Join(l, "_")
//	case constant.GroupChatType:
//		return "n_" + x.GroupID
//	case constant.SuperGroupChatType:
//		return "n_" + x.GroupID
//	case constant.NotificationChatType:
//		return "n_" + x.SendID + "_" + x.RecvID
//	}
//	return ""
//}
//
//func (x *MsgData) GetChatConversationIDByMsg() string {
//	switch x.SessionType {
//	case constant.SingleChatType:
//		l := []string{x.SendID, x.RecvID}
//		sort.Strings(l)
//		return "si_" + strings.Join(l, "_")
//	case constant.GroupChatType:
//		return "g_" + x.GroupID
//	case constant.SuperGroupChatType:
//		return "sg_" + x.GroupID
//	case constant.NotificationChatType:
//		return "sn_" + x.SendID + "_" + x.RecvID
//	}
//	return ""
//}
//
//func (x *MsgData) GenConversationUniqueKey() string {
//	switch x.SessionType {
//	case constant.SingleChatType, constant.NotificationChatType:
//		l := []string{x.SendID, x.RecvID}
//		sort.Strings(l)
//		return strings.Join(l, "_")
//	case constant.SuperGroupChatType:
//		return x.GroupID
//	}
//	return ""
//}
//
//func (x *MsgData) GetConversationIDByMsg() string {
//	options := Options(x.Options)
//	switch x.SessionType {
//	case constant.SingleChatType:
//		l := []string{x.SendID, x.RecvID}
//		sort.Strings(l)
//		if !options.IsNotNotification() {
//			return "n_" + strings.Join(l, "_")
//		}
//		return "si_" + strings.Join(l, "_") // single chat
//	case constant.GroupChatType:
//		if !options.IsNotNotification() {
//			return "n_" + x.GroupID // group chat
//		}
//		return "g_" + x.GroupID // group chat
//	case constant.SuperGroupChatType:
//		if !options.IsNotNotification() {
//			return "n_" + x.GroupID // super group chat
//		}
//		return "sg_" + x.GroupID // super group chat
//	case constant.NotificationChatType:
//		if !options.IsNotNotification() {
//			return "n_" + x.SendID + "_" + x.RecvID // super group chat
//		}
//		return "sn_" + x.SendID + "_" + x.RecvID // server notification chat
//	}
//	return ""
//}

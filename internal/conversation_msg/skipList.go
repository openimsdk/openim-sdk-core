package conversation_msg

//type SortConversationList struct {
//	list                  *skiplist.SkipList
//	pinnedConversationIDs map[string]struct{}
//}
//
//func NewActiveConversationList(pinnedIDs map[string]struct{}) *SortConversationList {
//	return &SortConversationList{
//		list: skiplist.New(skiplist.GreaterThan(func(a, b interface{}) int {
//			return compareConversations(a.(*msg.ActiveConversation), b.(*msg.ActiveConversation), pinnedIDs)
//		})),
//		pinnedConversationIDs: pinnedIDs,
//	}
//}
//
// Insert inserts or updates a conversation
//func (acl *SortConversationList) Insert(c *msg.ActiveConversation) {
//	acl.Delete(c.ConversationID)
//	acl.list.Set(c, c)
//}
//
// Update updates a conversation (same as Insert)
//func (acl *SortConversationList) Update(c *msg.ActiveConversation) {
//	acl.Insert(c)
//}
//
// Delete removes a conversation by its ID
//func (acl *SortConversationList) Delete(conversationID string) {
//	for e := acl.list.Front(); e != nil; e = e.Next() {
//		if e.Value.(*msg.ActiveConversation).ConversationID == conversationID {
//			acl.list.Remove(e.Key)
//			return
//		}
//	}
//}
//
// Init initializes the conversation list with the given list of conversations
//func (acl *SortConversationList) Init(conversations []*msg.ActiveConversation) {
//	acl.list = skiplist.New(skiplist.GreaterThan(func(a, b interface{}) int {
//		return compareConversations(a.(*msg.ActiveConversation), b.(*msg.ActiveConversation), acl.pinnedConversationIDs)
//	}))
//	for _, c := range conversations {
//		acl.Insert(c)
//	}
//}
//
// Top returns the top 'limit' conversations
//func (acl *SortConversationList) Top(limit int) []*msg.ActiveConversation {
//	res := make([]*msg.ActiveConversation, 0, limit)
//	for e := acl.list.Front(); e != nil && (limit <= 0 || len(res) < limit); e = e.Next() {
//		res = append(res, e.Value.(*msg.ActiveConversation))
//	}
//	return res
//}
//
// After returns the next 'n' conversations after the given conversation ID
//func (acl *SortConversationList) After(conversationID string, n int) []*msg.ActiveConversation {
//	var start *skiplist.Element
//	for e := acl.list.Front(); e != nil; e = e.Next() {
//		if e.Value.(*msg.ActiveConversation).ConversationID == conversationID {
//			start = e
//			break
//		}
//	}
//	if start == nil {
//		return nil
//	}
//	res := make([]*msg.ActiveConversation, 0, n)
//	for e := start.Next(); e != nil && (n <= 0 || len(res) < n); e = e.Next() {
//		res = append(res, e.Value.(*msg.ActiveConversation))
//	}
//	return res
//}
//
// All returns all conversations
//func (acl *SortConversationList) All() []*msg.ActiveConversation {
//	res := make([]*msg.ActiveConversation, 0, acl.list.Len())
//	for e := acl.list.Front(); e != nil; e = e.Next() {
//		res = append(res, e.Value.(*msg.ActiveConversation))
//	}
//	return res
//}
//
// Comparator for comparing two conversations
//func compareConversations(a, b *msg.ActiveConversation, pinned map[string]struct{}) int {
//	_, ap := pinned[a.ConversationID]
//	_, bp := pinned[b.ConversationID]
//	if ap != bp {
//		if ap {
//			return 1
//		}
//		return -1
//	}
//	atime := getEffectiveTime(a)
//	btime := getEffectiveTime(b)
//	if atime == btime {
//		return 0
//	}
//	if atime > btime {
//		return 1
//	}
//	return -1
//}
//
//func getEffectiveTime(c *msg.ActiveConversation) int64 {
//	if c.DraftTime > c.LatestMsgSendTime {
//		return c.DraftTime
//	}
//	return c.LatestMsgSendTime
//}

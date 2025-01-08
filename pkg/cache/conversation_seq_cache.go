package cache

import (
	"strings"

	"github.com/openimsdk/tools/utils/stringutil"
)

const (
	ViewHistory = iota
	ViewSearch
)

type ConversationSeqContextCache struct {
	*Cache[string, int64]
}

func NewConversationSeqContextCache() *ConversationSeqContextCache {
	return &ConversationSeqContextCache{Cache: NewCache[string, int64]()}
}

func (c ConversationSeqContextCache) Load(conversationID string, viewType int) (int64, bool) {
	return c.Cache.Load(c.getConversationViewTypeKey(conversationID, viewType))

}
func (c ConversationSeqContextCache) Delete(conversationID string, viewType int) {
	c.Cache.Delete(c.getConversationViewTypeKey(conversationID, viewType))

}

func (c ConversationSeqContextCache) Store(conversationID string, viewType int, thisEndSeq int64) {
	c.Cache.Store(c.getConversationViewTypeKey(conversationID, viewType), thisEndSeq)

}

func (c ConversationSeqContextCache) StoreWithFunc(conversationID string, viewType int, thisEndSeq int64, fn func(key string, value int64) bool) {

	c.Cache.StoreWithFunc(c.getConversationViewTypeKey(conversationID, viewType), thisEndSeq, fn)
}

func (c ConversationSeqContextCache) getConversationViewTypeKey(conversationID string, viewType int) string {
	return conversationID + "::viewType::" + stringutil.IntToString(viewType)
}

func (c ConversationSeqContextCache) DeleteByViewType(viewType int) {
	c.Cache.DeleteCon(func(key string, value int64) bool {
		return strings.Contains(key, "::viewType::"+stringutil.IntToString(viewType))
	})
}

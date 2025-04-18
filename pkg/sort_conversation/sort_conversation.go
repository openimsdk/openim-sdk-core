package sort_conversation

import (
	"context"
	"sort"
	"strings"

	"github.com/openimsdk/tools/utils/datautil"
)

type ConversationMetaData struct {
	ConversationID    string
	IsPinned          bool
	LatestMsgSendTime int64
	DraftTextTime     int64
}

type SortConversationList struct {
	conversations         []*ConversationMetaData
	pinnedConversationIDs map[string]struct{}
	conversationIndex     map[string]int
}

func NewSortConversationList(list []*ConversationMetaData, pinnedIDs []string) *SortConversationList {
	s := &SortConversationList{
		pinnedConversationIDs: datautil.SliceSetAny(pinnedIDs, func(e string) string {
			return e
		}),
		conversationIndex: make(map[string]int),
	}
	s.conversations = list
	sort.Slice(s.conversations, func(i, j int) bool {
		return s.compareConversations(s.conversations[i], s.conversations[j])
	})
	s.refreshIndex()
	return s
}

// InsertOrUpdate inserts or updates a conversation, updating its LatestMsgSendTime / DraftTime / pinned status (all fields are passed in)
func (l *SortConversationList) InsertOrUpdate(c *ConversationMetaData) {
	if c.IsPinned {
		l.pinnedConversationIDs[c.ConversationID] = struct{}{}
	} else {
		delete(l.pinnedConversationIDs, c.ConversationID)
	}
	if idx, exists := l.conversationIndex[c.ConversationID]; exists {
		l.conversations = append(l.conversations[:idx], l.conversations[idx+1:]...)
	}
	idx := l.findInsertIndex(c)
	l.conversations = append(l.conversations, nil)
	copy(l.conversations[idx+1:], l.conversations[idx:])
	l.conversations[idx] = c
	l.refreshIndex()
}

// Delete removes a conversation
func (l *SortConversationList) Delete(conversationID string) {
	idx, exists := l.conversationIndex[conversationID]
	if !exists {
		return
	}
	l.conversations = append(l.conversations[:idx], l.conversations[idx+1:]...)
	delete(l.pinnedConversationIDs, conversationID)
	l.refreshIndex()
}

// Top returns the top N conversations
func (l *SortConversationList) Top(limit int) []*ConversationMetaData {
	if limit > 0 && len(l.conversations) > limit {
		return l.conversations[:limit]
	}
	return l.conversations
}

// After gets the N conversations after a specific conversation
func (l *SortConversationList) After(conversationID string, n int) []*ConversationMetaData {
	idx, exists := l.conversationIndex[conversationID]
	if !exists {
		return nil
	}
	start := idx + 1
	end := start + n
	if end > len(l.conversations) {
		end = len(l.conversations)
	}
	return l.conversations[start:end]
}

// All returns all conversations
func (l *SortConversationList) All() []*ConversationMetaData {
	return l.conversations
}

// Internal functions

func (l *SortConversationList) findInsertIndex(c *ConversationMetaData) int {
	for i, exist := range l.conversations {
		if l.compareConversations(c, exist) {
			return i
		}
	}
	return len(l.conversations)
}

func (l *SortConversationList) refreshIndex() {
	l.conversationIndex = make(map[string]int)
	for i, c := range l.conversations {
		l.conversationIndex[c.ConversationID] = i
	}
}

func (l *SortConversationList) compareConversations(a, b *ConversationMetaData) bool {
	_, ap := l.pinnedConversationIDs[a.ConversationID]
	_, bp := l.pinnedConversationIDs[b.ConversationID]
	if ap != bp {
		return ap // pinned conversations come first
	}
	at, bt := l.getEffectiveTime(a), l.getEffectiveTime(b)
	return at > bt
}

func (l *SortConversationList) getEffectiveTime(c *ConversationMetaData) int64 {
	if c.DraftTextTime > c.LatestMsgSendTime {
		return c.DraftTextTime
	}
	return c.LatestMsgSendTime
}

func (l *SortConversationList) NewIterator() *SortConversationIterator {
	return &SortConversationIterator{
		list:   l,
		cursor: 0,
	}
}

type SortConversationIterator struct {
	list   *SortConversationList
	cursor int
}

// NextTop returns the next batch of size conversations
func (it *SortConversationIterator) NextTop(size int) []*ConversationMetaData {
	if it.cursor >= len(it.list.conversations) {
		return nil
	}
	end := it.cursor + size
	if end > len(it.list.conversations) {
		end = len(it.list.conversations)
	}
	batch := it.list.conversations[it.cursor:end]
	it.cursor = end
	return batch
}

type BatchHandler func(ctx context.Context, batchID int, needSyncSeqMap map[string][2]int64, isFirst bool)

type ConversationBatchProcessor struct {
	iter           *SortConversationIterator
	needSyncSeqMap map[string][2]int64
	isFirst        bool
	batchSize      int
	batchID        int
}

func NewConversationBatchProcessor(list *SortConversationList, needSyncSeqMap map[string][2]int64, batchSize int) *ConversationBatchProcessor {
	return &ConversationBatchProcessor{
		iter:           list.NewIterator(),
		needSyncSeqMap: needSyncSeqMap,
		batchSize:      batchSize,
		isFirst:        true,
		batchID:        1,
	}
}

func (p *ConversationBatchProcessor) Run(ctx context.Context, handler BatchHandler) {
	result := make(map[string][2]int64)

	for {
		batch := p.iter.NextTop(p.batchSize)
		if len(batch) == 0 {
			break
		}

		for _, conv := range batch {
			if v, ok := p.needSyncSeqMap[conv.ConversationID]; ok {
				result[conv.ConversationID] = v
			}
			notificationID := GetNotificationConversationIDByConversationID(conv.ConversationID)
			if v, ok := p.needSyncSeqMap[notificationID]; ok {
				result[notificationID] = v
			}

			if len(result) >= p.batchSize {
				p.emitResult(ctx, handler, result)
				result = make(map[string][2]int64)
			}
		}
	}

	if len(result) > 0 {
		p.emitResult(ctx, handler, result)
	}
}

func (p *ConversationBatchProcessor) emitResult(ctx context.Context, handler BatchHandler, result map[string][2]int64) {
	handler(ctx, p.batchID, result, p.isFirst)
	p.batchID++
	p.isFirst = false
}

func GetNotificationConversationIDByConversationID(conversationID string) string {
	l := strings.Split(conversationID, "_")
	if len(l) > 1 {
		l[0] = "n"
		return strings.Join(l, "_")
	}
	return ""
}

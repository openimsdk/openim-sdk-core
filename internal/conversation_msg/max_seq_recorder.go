package conversation_msg

type MaxSeqRecorder map[string]int64

func NewMaxSeqRecorder() MaxSeqRecorder {
	return make(map[string]int64)
}

func (m MaxSeqRecorder) Get(conversationID string) int64 {
	return m[conversationID]
}

func (m MaxSeqRecorder) Set(conversationID string, seq int64) {
	m[conversationID] = seq
}

func (m MaxSeqRecorder) Incr(conversationID string, num int64) {
	m[conversationID] += num
}

func (m MaxSeqRecorder) IsNewMsg(conversationID string, seq int64) bool {
	currentSeq := m[conversationID]
	return seq > currentSeq
}

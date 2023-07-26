// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conversation_msg

import "sync"

type MaxSeqRecorder struct {
	seqs map[string]int64
	lock sync.RWMutex
}

func NewMaxSeqRecorder() MaxSeqRecorder {
	m := make(map[string]int64)
	return MaxSeqRecorder{seqs: m}
}

func (m *MaxSeqRecorder) Get(conversationID string) int64 {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.seqs[conversationID]
}

func (m *MaxSeqRecorder) Set(conversationID string, seq int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.seqs[conversationID] = seq
}

func (m *MaxSeqRecorder) Incr(conversationID string, num int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.seqs[conversationID] += num
}

func (m *MaxSeqRecorder) IsNewMsg(conversationID string, seq int64) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	currentSeq := m.seqs[conversationID]
	return seq > currentSeq
}

// Copyright © 2026 OpenIM SDK. All rights reserved.
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

//go:build !js
// +build !js

package db

import (
	"context"
	"strings"
	"testing"
)

func TestConversationListSplitUsesVisibleOrderIndex(t *testing.T) {
	ctx := context.Background()
	database, err := NewDataBase(ctx, "conversation_index_user", t.TempDir(), 0)
	if err != nil {
		t.Fatalf("NewDataBase failed: %v", err)
	}
	t.Cleanup(func() {
		if err := database.Close(ctx); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	})

	type explainRow struct {
		Detail string
	}
	var rows []explainRow
	err = database.conn.WithContext(ctx).
		Raw("EXPLAIN QUERY PLAN SELECT * FROM local_conversations WHERE latest_msg_send_time > 0 ORDER BY "+conversationOrder+" LIMIT ?", 50).
		Scan(&rows).Error
	if err != nil {
		t.Fatalf("EXPLAIN QUERY PLAN failed: %v", err)
	}

	var details []string
	for _, row := range rows {
		details = append(details, row.Detail)
	}
	plan := strings.Join(details, "\n")
	if !strings.Contains(plan, "idx_local_conversations_visible_order") {
		t.Fatalf("query plan does not use visible order index:\n%s", plan)
	}
	if strings.Contains(plan, "USE TEMP B-TREE FOR ORDER BY") {
		t.Fatalf("query plan still sorts with a temp b-tree:\n%s", plan)
	}
}

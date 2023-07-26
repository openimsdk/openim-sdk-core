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

package syncer

import (
	"context"
	"reflect"

	"github.com/google/go-cmp/cmp"

	"github.com/OpenIMSDK/tools/log"
)

func New[T any, V comparable](
	insert func(ctx context.Context, value T) error,
	delete func(ctx context.Context, value T) error,
	update func(ctx context.Context, server T, local T) error,
	uuid func(value T) V,
	equal func(a T, b T) bool,
	notice func(ctx context.Context, state int, server, local T) error,
) *Syncer[T, V] {
	if insert == nil || update == nil || delete == nil || uuid == nil {
		panic("invalid params")
	}
	var t T
	tof := reflect.TypeOf(&t)
	for tof.Kind() == reflect.Ptr {
		tof = tof.Elem()
	}
	return &Syncer[T, V]{
		insert: insert,
		update: update,
		delete: delete,
		uuid:   uuid,
		equal:  equal,
		notice: notice,
		ts:     tof.String(),
	}
}

type Syncer[T any, V comparable] struct {
	insert func(ctx context.Context, server T) error
	update func(ctx context.Context, server T, local T) error
	delete func(ctx context.Context, local T) error
	notice func(ctx context.Context, state int, server, local T) error
	equal  func(server T, local T) bool
	uuid   func(value T) V
	ts     string
}

func (s *Syncer[T, V]) eq(server T, local T) bool {
	if s.equal != nil {
		return s.equal(server, local)
	}
	return cmp.Equal(server, local)
}

func (s *Syncer[T, V]) onNotice(ctx context.Context, state int, server, local T, fn func(ctx context.Context, state int, server, local T) error) error {
	if s.notice != nil {
		if err := s.notice(ctx, state, server, local); err != nil {
			return err
		}
	}
	if fn != nil {
		if err := fn(ctx, state, server, local); err != nil {
			return err
		}
	}
	return nil
}

func (s *Syncer[T, V]) Sync(ctx context.Context, serverData []T, localData []T, notice func(ctx context.Context, state int, server, local T) error, noDel ...bool) (err error) {
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "sync success", "type", s.ts)
		} else {
			log.ZError(ctx, "sync failed", err, "type", s.ts)
		}
	}()
	if len(serverData) == 0 && len(localData) == 0 {
		log.ZDebug(ctx, "sync both the server and client are empty", "type", s.ts)
		return nil
	}
	localMap := make(map[V]T)
	for i, item := range localData {
		localMap[s.uuid(item)] = localData[i]
	}
	for i := range serverData {
		server := serverData[i]
		id := s.uuid(server)
		local, ok := localMap[id]
		if !ok {
			if err := s.insert(ctx, server); err != nil {
				log.ZError(ctx, "sync insert failed", err, "type", s.ts, "server", server, "local", local)
				return err
			}
			if err := s.onNotice(ctx, Insert, server, local, notice); err != nil {
				log.ZError(ctx, "sync notice insert failed", err, "type", s.ts, "server", server, "local", local)
				return err
			}
			continue
		}
		delete(localMap, id)
		log.ZDebug(ctx, "syncer come", "type", s.ts, "server", server, "local", local, "isEq", s.eq(server, local))

		if s.eq(server, local) {
			if err := s.onNotice(ctx, Unchanged, local, server, notice); err != nil {
				log.ZError(ctx, "sync notice unchanged failed", err, "type", s.ts, "server", server, "local", local)
				return err
			}
			continue
		}
		if err := s.update(ctx, server, local); err != nil {
			log.ZError(ctx, "sync update failed", err, "type", s.ts, "server", server, "local", local)
			return err
		}
		if err := s.onNotice(ctx, Update, server, local, notice); err != nil {
			log.ZError(ctx, "sync notice update failed", err, "type", s.ts, "server", server, "local", local)
			return err
		}
	}
	if len(noDel) > 0 && noDel[0] {
		return nil
	}
	for id := range localMap {
		local := localMap[id]
		if err := s.delete(ctx, local); err != nil {
			log.ZError(ctx, "sync delete failed", err, "type", s.ts, "local", local)
			return err
		}
		var server T
		if err := s.onNotice(ctx, Delete, server, local, notice); err != nil {
			log.ZError(ctx, "sync notice delete failed", err, "type", s.ts, "local", local)
			return err
		}
	}
	return nil
}

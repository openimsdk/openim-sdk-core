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

// New creates a new Syncer instance with the provided synchronization functions.
// This function requires several callback functions to handle different aspects of data synchronization:
// - insert: A function to insert new data.
// - delete: A function to delete existing data.
// - update: A function to update existing data.
// - uuid: A function to generate a unique identifier for each data item.
// - equal: A function to check if two data items are equal.
// - notice: A function to handle notifications during the sync process.
// Panics if insert, delete, update, or uuid functions are not provided.
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

	// Determine the type of T and remove pointer indirection if necessary.
	var t T
	tof := reflect.TypeOf(&t)
	for tof.Kind() == reflect.Ptr {
		tof = tof.Elem()
	}

	// Return a new Syncer instance with the provided functions and the type as a string.
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

// Syncer is a struct that holds functions for synchronizing data.
// It includes functions for inserting, updating, deleting, and notifying,
// as well as functions for generating unique IDs and checking equality of data items.
type Syncer[T any, V comparable] struct {
	insert func(ctx context.Context, server T) error
	update func(ctx context.Context, server T, local T) error
	delete func(ctx context.Context, local T) error
	notice func(ctx context.Context, state int, server, local T) error
	equal  func(server T, local T) bool
	uuid   func(value T) V
	ts     string // Represents the type of T as a string.
}

// eq is a helper function to check equality of two data items.
// It uses the equal function if provided; otherwise, it falls back to the cmp.Equal function.
func (s *Syncer[T, V]) eq(server T, local T) bool {
	if s.equal != nil {
		return s.equal(server, local)
	}
	return cmp.Equal(server, local)
}

// onNotice is a helper function to handle notifications.
// It calls the Syncer's notice function and the provided notice function in sequence if they are not nil.
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

// Sync synchronizes server data with local data.
// Sync synchronizes the data between the server and local storage.
// It takes a context, two slices of data (serverData and localData),
// a notice function to handle notifications during the sync process,
// and a variadic parameter noDel to control deletion behavior.
func (s *Syncer[T, V]) Sync(ctx context.Context, serverData []T, localData []T, notice func(ctx context.Context, state int, server, local T) error, noDel ...bool) (err error) {
	defer func() {
		// Log the outcome of the synchronization process.
		if err == nil {
			log.ZDebug(ctx, "sync success", "type", s.ts)
		} else {
			log.ZError(ctx, "sync failed", err, "type", s.ts)
		}
	}()

	// If both server and local data are empty, log and return.
	if len(serverData) == 0 && len(localData) == 0 {
		log.ZDebug(ctx, "sync both the server and client are empty", "type", s.ts)
		return nil
	}

	// Convert local data into a map for easier lookup.
	localMap := make(map[V]T)
	for i, item := range localData {
		localMap[s.uuid(item)] = localData[i]
	}

	// Iterate through server data to sync with local data.
	for i := range serverData {
		server := serverData[i]
		id := s.uuid(server)
		local, ok := localMap[id]

		// If the item doesn't exist locally, insert it.
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

		// Remove the item from the local map as it's found.
		delete(localMap, id)
		log.ZDebug(ctx, "syncer come", "type", s.ts, "server", server, "local", local, "isEq", s.eq(server, local))

		// If the local and server items are equal, notify and continue.
		if s.eq(server, local) {
			if err := s.onNotice(ctx, Unchanged, local, server, notice); err != nil {
				log.ZError(ctx, "sync notice unchanged failed", err, "type", s.ts, "server", server, "local", local)
				return err
			}
			continue
		}

		// Update the local item with server data.
		if err := s.update(ctx, server, local); err != nil {
			log.ZError(ctx, "sync update failed", err, "type", s.ts, "server", server, "local", local)
			return err
		}
		if err := s.onNotice(ctx, Update, server, local, notice); err != nil {
			log.ZError(ctx, "sync notice update failed", err, "type", s.ts, "server", server, "local", local)
			return err
		}
	}

	// Check the noDel flag; if set, skip deletion.
	if len(noDel) > 0 && noDel[0] {
		return nil
	}

	// Delete any local items not present in server data.
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

// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package file

type PutFileCallback interface {
	Open(size int64)
	HashProgress(current, total int64)
	HashComplete(hash string, total int64)
	PutStart(current, total int64)
	PutProgress(save int64, current, total int64)
	PutComplete(total int64, putType int)
}

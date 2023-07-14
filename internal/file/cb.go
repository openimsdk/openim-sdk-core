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

package file

import "fmt"

type UploadFileCallback interface {
	Open(size int64)                                                    // 文件打开的大小
	PartSize(partSize int64, num int)                                   // 分片大小,数量
	HashPartProgress(index int, size int64, partHash string)            // 每块分片的hash值
	HashPartComplete(partsHash string, fileHash string)                 // 分块完成，服务端标记hash和文件最终hash
	UploadID(uploadID string)                                           // 上传ID
	UploadPartComplete(index int, partSize int64, partHash string)      // 上传分片进度
	UploadComplete(fileSize int64, streamSize int64, storageSize int64) // 整体进度
	Complete(size int64, url string, typ int)                           // 上传完成
}

type emptyUploadCallback struct{}

func (e emptyUploadCallback) Open(size int64) {
	fmt.Println("Callback Open:", size)
}

func (e emptyUploadCallback) PartSize(partSize int64, num int) {
	fmt.Println("Callback PartSize:", partSize, num)
}

func (e emptyUploadCallback) HashPartProgress(index int, size int64, partHash string) {
	//fmt.Println("Callback HashPartProgress:", index, size, partHash)
}

func (e emptyUploadCallback) HashPartComplete(partsHash string, fileHash string) {
	fmt.Println("Callback HashPartComplete:", partsHash, fileHash)
}

func (e emptyUploadCallback) UploadID(uploadID string) {
	fmt.Println("Callback UploadID:", uploadID)
}

func (e emptyUploadCallback) UploadPartComplete(index int, partSize int64, partHash string) {
	fmt.Println("Callback UploadPartComplete:", index, partSize, partHash)
}

func (e emptyUploadCallback) UploadComplete(fileSize int64, streamSize int64, storageSize int64) {
	fmt.Println("Callback UploadComplete:", fileSize, streamSize, storageSize)
}

func (e emptyUploadCallback) Complete(size int64, url string, typ int) {
	fmt.Println("Callback Complete:", size, url, typ)
}

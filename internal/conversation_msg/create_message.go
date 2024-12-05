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

import (
	"context"
	"errors"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/tools/log"

	"os"
	"path/filepath"
	"strings"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

func (c *Conversation) CreateTextMessage(ctx context.Context, text string) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Text)
	if err != nil {
		return nil, err
	}
	s.TextElem = &sdk_struct.TextElem{Content: text}
	return &s, nil
}
func (c *Conversation) CreateAdvancedTextMessage(ctx context.Context, text string, messageEntities []*sdk_struct.MessageEntity) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.AdvancedText)
	if err != nil {
		return nil, err
	}
	s.AdvancedTextElem = &sdk_struct.AdvancedTextElem{
		Text:              text,
		MessageEntityList: messageEntities,
	}
	return &s, nil
}

func (c *Conversation) CreateTextAtMessage(ctx context.Context, text string, userIDList []string, usersInfo []*sdk_struct.AtInfo, qs *sdk_struct.MsgStruct) (*sdk_struct.MsgStruct, error) {
	if text == "" {
		return nil, errors.New("text can not be empty")
	}
	if len(userIDList) > 10 {
		return nil, sdkerrs.ErrArgs
	}
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.AtText)
	if err != nil {
		return nil, err
	}
	//Avoid nested references
	if qs != nil {
		if qs.ContentType == constant.Quote {
			qs.ContentType = constant.Text
			qs.TextElem = &sdk_struct.TextElem{Content: qs.QuoteElem.Text}
		}
	}
	s.AtTextElem = &sdk_struct.AtTextElem{
		Text:         text,
		AtUserList:   userIDList,
		AtUsersInfo:  usersInfo,
		QuoteMessage: qs,
	}
	return &s, nil
}
func (c *Conversation) CreateLocationMessage(ctx context.Context, description string, longitude, latitude float64) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Location)
	if err != nil {
		return nil, err
	}
	s.LocationElem = &sdk_struct.LocationElem{
		Description: description,
		Longitude:   longitude,
		Latitude:    latitude,
	}
	return &s, nil
}

func (c *Conversation) CreateCustomMessage(ctx context.Context, data, extension string, description string) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Custom)
	if err != nil {
		return nil, err
	}
	s.CustomElem = &sdk_struct.CustomElem{
		Data:        data,
		Extension:   extension,
		Description: description,
	}
	return &s, nil

}
func (c *Conversation) CreateQuoteMessage(ctx context.Context, text string, qs *sdk_struct.MsgStruct) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Quote)
	if err != nil {
		return nil, err
	}
	//Avoid nested references
	if qs.ContentType == constant.Quote {
		qs.ContentType = constant.Text
		qs.TextElem = &sdk_struct.TextElem{Content: qs.QuoteElem.Text}
	}
	s.QuoteElem = &sdk_struct.QuoteElem{
		Text:         text,
		QuoteMessage: qs,
	}
	return &s, nil

}
func (c *Conversation) CreateAdvancedQuoteMessage(ctx context.Context, text string, qs *sdk_struct.MsgStruct, messageEntities []*sdk_struct.MessageEntity) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Quote)
	if err != nil {
		return nil, err
	}
	//Avoid nested references
	if qs.ContentType == constant.Quote {
		//qs.Content = qs.QuoteElem.Text
		qs.ContentType = constant.Text
		qs.TextElem = &sdk_struct.TextElem{Content: qs.QuoteElem.Text}
	}
	s.QuoteElem = &sdk_struct.QuoteElem{
		Text:              text,
		QuoteMessage:      qs,
		MessageEntityList: messageEntities,
	}
	return &s, nil
}

func (c *Conversation) CreateCardMessage(ctx context.Context, card *sdk_struct.CardElem) (*sdk_struct.MsgStruct,
	error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Card)
	if err != nil {
		return nil, err
	}
	s.CardElem = card
	return &s, nil
}

func (c *Conversation) CreateImageMessage(ctx context.Context, imageSourcePath string, sourcePicture, bigPicture, snapshotPicture *sdk_struct.PictureBaseInfo) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Picture)
	if err != nil {
		return nil, err
	}

	// Create by file path
	if sourcePicture != nil || bigPicture != nil || snapshotPicture != nil {
		dstFile := utils.FileTmpPath(imageSourcePath, c.DataDir) //a->b
		_, err := utils.CopyFile(imageSourcePath, dstFile)
		if err != nil {
			//log.Error(operationID, "open file failed: ", err, imageFullPath)
			return nil, err
		}

		imageInfo, err := getImageInfo(imageSourcePath)
		if err != nil {
			//log.Error(operationID, "getImageInfo err:", err.Error())
			return nil, err
		}

		s.PictureElem = &sdk_struct.PictureElem{
			SourcePath: imageSourcePath,
			SourcePicture: &sdk_struct.PictureBaseInfo{
				Width:  imageInfo.Width,
				Height: imageInfo.Height,
				Type:   imageInfo.Type,
			},
		}
	} else { // Create by URL
		s.PictureElem = &sdk_struct.PictureElem{
			SourcePath:      imageSourcePath,
			SourcePicture:   sourcePicture,
			BigPicture:      bigPicture,
			SnapshotPicture: snapshotPicture,
		}
	}

	return &s, nil
}

func (c *Conversation) CreateSoundMessage(ctx context.Context, soundPath string, duration int64, soundElem *sdk_struct.SoundBaseInfo) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}

	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Sound)
	if err != nil {
		return nil, err
	}

	// Create by file path
	if soundElem == nil {
		dstFile := utils.FileTmpPath(soundPath, c.DataDir) //a->b
		_, err := utils.CopyFile(soundPath, dstFile)
		if err != nil {
			//log.Error("internal", "open file failed: ", err, soundPath)
			return nil, err
		}

		fi, err := os.Stat(soundPath)
		if err != nil {
			//log.Error("internal", "getSoundInfo err:", err.Error(), s.SoundElem.SoundPath)
			return nil, err
		}

		s.SoundElem = &sdk_struct.SoundElem{
			SoundPath: soundPath,
			Duration:  duration,
			DataSize:  fi.Size(),
			SoundType: strings.Replace(filepath.Ext(fi.Name()), ".", "", 1),
		}
	} else { // Create by URL
		s.SoundElem = &sdk_struct.SoundElem{
			UUID:      soundElem.UUID,
			SoundPath: soundElem.SoundPath,
			SourceURL: soundElem.SourceURL,
			DataSize:  soundElem.DataSize,
			Duration:  soundElem.Duration,
			SoundType: soundElem.SoundType,
		}
	}
	return &s, nil
}

func (c *Conversation) CreateVideoMessage(ctx context.Context, videoSourcePath string, videoType string, duration int64, snapshotSourcePath string, videoElem *sdk_struct.VideoBaseInfo) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Video)
	if err != nil {
		return nil, err
	}

	// Create by file path
	if videoElem == nil {
		dstFile := utils.FileTmpPath(videoSourcePath, c.DataDir) //a->b
		written, err := utils.CopyFile(videoSourcePath, dstFile)
		if err != nil {
			//log.Error("internal", "open file failed: ", err, videoFullPath)
			return nil, err
		}

		log.ZDebug(ctx, "videoFullPath dstFile", "videoFullPath", videoSourcePath,
			"dstFile", dstFile, "written", written)

		dstFile = utils.FileTmpPath(snapshotSourcePath, c.DataDir) //a->b
		sWritten, err := utils.CopyFile(snapshotSourcePath, dstFile)
		if err != nil {
			//log.Error("internal", "open file failed: ", err, snapshotFullPath)
			return nil, err
		}

		log.ZDebug(ctx, "snapshotFullPath dstFile", "snapshotFullPath", snapshotSourcePath,
			"dstFile", dstFile, "sWritten", sWritten)

		s.VideoElem = &sdk_struct.VideoElem{
			VideoPath: videoSourcePath,
			VideoType: videoType,
			Duration:  duration,
		}

		if snapshotSourcePath == "" {
			s.VideoElem.SnapshotPath = ""
		} else {
			s.VideoElem.SnapshotPath = snapshotSourcePath
		}

		fi, err := os.Stat(s.VideoElem.VideoPath)
		if err != nil {
			//log.Error("internal", "get file Attributes error", err.Error())
			return nil, err
		}

		s.VideoElem.VideoSize = fi.Size()
		if snapshotSourcePath != "" {
			imageInfo, err := getImageInfo(s.VideoElem.SnapshotPath)
			if err != nil {
				log.ZError(ctx, "getImageInfo err:", err, "snapshotFullPath", snapshotSourcePath)
				return nil, err
			}

			s.VideoElem.SnapshotHeight = imageInfo.Height
			s.VideoElem.SnapshotWidth = imageInfo.Width
			s.VideoElem.SnapshotSize = imageInfo.Size
		}
	} else { // Create by URL
		s.VideoElem = &sdk_struct.VideoElem{
			VideoPath:      videoElem.VideoPath,
			VideoUUID:      videoElem.VideoUUID,
			VideoURL:       videoElem.VideoURL,
			VideoType:      videoElem.VideoType,
			VideoSize:      videoElem.VideoSize,
			Duration:       videoElem.Duration,
			SnapshotPath:   videoElem.SnapshotPath,
			SnapshotUUID:   videoElem.SnapshotUUID,
			SnapshotSize:   videoElem.SnapshotSize,
			SnapshotURL:    videoElem.SnapshotURL,
			SnapshotWidth:  videoElem.SnapshotWidth,
			SnapshotHeight: videoElem.SnapshotHeight,
			SnapshotType:   videoElem.SnapshotType,
		}
	}

	return &s, nil
}

func (c *Conversation) CreateFileMessage(ctx context.Context, fileSourcePath string, fileName string, fileElem *sdk_struct.FileBaseInfo) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.File)
	if err != nil {
		return nil, err
	}

	// Create by file path
	if fileElem == nil {
		dstFile := utils.FileTmpPath(fileSourcePath, c.DataDir)
		_, err := utils.CopyFile(fileSourcePath, dstFile)
		if err != nil {
			//log.Error("internal", "open file failed: ", err.Error(), fileFullPath)
			return nil, err

		}

		fi, err := os.Stat(fileSourcePath)
		if err != nil {
			//log.Error("internal", "get file Attributes error", err.Error())
			return nil, err
		}

		s.FileElem = &sdk_struct.FileElem{
			FilePath: fileSourcePath,
			FileName: fileName,
			FileSize: fi.Size(),
		}
	} else { // Create by URL
		s.FileElem = &sdk_struct.FileElem{
			FilePath:  fileElem.FilePath,
			UUID:      fileElem.UUID,
			SourceURL: fileElem.SourceURL,
			FileName:  fileElem.FileName,
			FileSize:  fileElem.FileSize,
			FileType:  fileElem.FileType,
		}
	}

	return &s, nil
}

func (c *Conversation) CreateMergerMessage(ctx context.Context, messages []*sdk_struct.MsgStruct, title string, summaries []string) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{MergeElem: &sdk_struct.MergeElem{}}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Merger)
	if err != nil {
		return nil, err
	}
	s.MergeElem.AbstractList = summaries
	s.MergeElem.Title = title
	s.MergeElem.MultiMessage = messages
	s.Content = utils.StructToJsonString(s.MergeElem)
	return &s, nil
}
func (c *Conversation) CreateFaceMessage(ctx context.Context, index int, data string) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{FaceElem: &sdk_struct.FaceElem{}}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Face)
	if err != nil {
		return nil, err
	}
	s.FaceElem.Data = data
	s.FaceElem.Index = index
	s.Content = utils.StructToJsonString(s.FaceElem)
	return &s, nil
}

func (c *Conversation) CreateForwardMessage(ctx context.Context, s *sdk_struct.MsgStruct) (*sdk_struct.MsgStruct, error) {
	if s.Status != constant.MsgStatusSendSuccess {
		log.ZError(ctx, "only send success message can be Forward",
			errors.New("only send success message can be Forward"))
		return nil, errors.New("only send success message can be Forward")
	}
	err := c.initBasicInfo(ctx, s, constant.UserMsgType, s.ContentType)
	if err != nil {
		return nil, err
	}
	//Forward message seq is set to 0
	s.Seq = 0
	s.Status = constant.MsgStatusSendSuccess
	return s, nil
}

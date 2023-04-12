package conversation_msg

import (
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"os"
)

func (c *Conversation) CreateTextMessage(ctx context.Context, text string) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Text)
	if err != nil {
		return nil, err
	}
	s.Content = text
	return &s, nil
}
func (c *Conversation) CreateAdvancedTextMessage(ctx context.Context, text string, messageEntitys []*sdk_struct.MessageEntity) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.AdvancedText)
	if err != nil {
		return nil, err
	}
	s.MessageEntityElem.Text = text
	s.MessageEntityElem.MessageEntityList = messageEntitys
	s.Content = utils.StructToJsonString(s.MessageEntityElem)
	return &s, nil
}
func (c *Conversation) CreateTextAtMessage(ctx context.Context, text string, userIDList []string, usersInfo []*sdk_struct.AtInfo, qs *sdk_struct.MsgStruct) (*sdk_struct.MsgStruct, error) {
	if text == "" {
		return nil, errors.New("text can not be empty")
	}
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.AtText)
	if err != nil {
		return nil, err
	}
	//Avoid nested references
	if qs.ContentType == constant.Quote {
		qs.Content = qs.QuoteElem.Text
		qs.ContentType = constant.Text
	}
	s.AtElem.Text = text
	s.AtElem.AtUserList = userIDList
	s.AtElem.AtUsersInfo = usersInfo
	s.AtElem.QuoteMessage = qs
	s.Content = utils.StructToJsonString(s.AtElem)
	return &s, nil
}
func (c *Conversation) CreateLocationMessage(ctx context.Context, description string, longitude, latitude float64) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Location)
	if err != nil {
		return nil, err
	}
	s.LocationElem.Description = description
	s.LocationElem.Longitude = longitude
	s.LocationElem.Latitude = latitude
	s.Content = utils.StructToJsonString(s.LocationElem)
	return &s, nil

}
func (c *Conversation) CreateCustomMessage(ctx context.Context, data, extension string, description string) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Custom)
	if err != nil {
		return nil, err
	}
	s.CustomElem.Data = data
	s.CustomElem.Extension = extension
	s.CustomElem.Description = description
	s.Content = utils.StructToJsonString(s.CustomElem)
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
		qs.Content = qs.QuoteElem.Text
		qs.ContentType = constant.Text
	}
	s.QuoteElem.Text = text
	s.QuoteElem.QuoteMessage = qs
	s.Content = utils.StructToJsonString(s.QuoteElem)
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
		qs.Content = qs.QuoteElem.Text
		qs.ContentType = constant.Text
	}
	s.QuoteElem.Text = text
	s.QuoteElem.MessageEntityList = messageEntities
	s.QuoteElem.QuoteMessage = qs
	s.Content = utils.StructToJsonString(s.QuoteElem)
	return &s, nil

}
func (c *Conversation) CreateCardMessage(ctx context.Context, cardInfo string) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Card)
	if err != nil {
		return nil, err
	}
	s.Content = cardInfo
	return &s, nil

}
func (c *Conversation) CreateVideoMessageFromFullPath(ctx context.Context, videoFullPath string, videoType string, duration int64, snapshotFullPath string) (*sdk_struct.MsgStruct, error) {
	dstFile := utils.FileTmpPath(videoFullPath, c.DataDir) //a->b
	written, err := utils.CopyFile(videoFullPath, dstFile)
	if err != nil {
		//log.Error("internal", "open file failed: ", err, videoFullPath)
		return nil, err
	}
	log.Info("internal", "videoFullPath dstFile", videoFullPath, dstFile, written)
	dstFile = utils.FileTmpPath(snapshotFullPath, c.DataDir) //a->b
	sWritten, err := utils.CopyFile(snapshotFullPath, dstFile)
	if err != nil {
		//log.Error("internal", "open file failed: ", err, snapshotFullPath)
		return nil, err
	}
	log.Info("internal", "snapshotFullPath dstFile", snapshotFullPath, dstFile, sWritten)

	s := sdk_struct.MsgStruct{}
	err = c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Video)
	if err != nil {
		return nil, err
	}
	s.VideoElem.VideoPath = videoFullPath
	s.VideoElem.VideoType = videoType
	s.VideoElem.Duration = duration
	if snapshotFullPath == "" {
		s.VideoElem.SnapshotPath = ""
	} else {
		s.VideoElem.SnapshotPath = snapshotFullPath
	}
	fi, err := os.Stat(s.VideoElem.VideoPath)
	if err != nil {
		//log.Error("internal", "get file Attributes error", err.Error())
		return nil, err
	}
	s.VideoElem.VideoSize = fi.Size()
	if snapshotFullPath != "" {
		imageInfo, err := getImageInfo(s.VideoElem.SnapshotPath)
		if err != nil {
			log.Error("internal", "get Image Attributes error", err.Error())
			return nil, err
		}
		s.VideoElem.SnapshotHeight = imageInfo.Height
		s.VideoElem.SnapshotWidth = imageInfo.Width
		s.VideoElem.SnapshotSize = imageInfo.Size
	}
	s.Content = utils.StructToJsonString(s.VideoElem)
	return &s, nil

}
func (c *Conversation) CreateFileMessageFromFullPath(ctx context.Context, fileFullPath string, fileName string) (*sdk_struct.MsgStruct, error) {
	dstFile := utils.FileTmpPath(fileFullPath, c.DataDir)
	_, err := utils.CopyFile(fileFullPath, dstFile)
	//log.Info(operationID, "copy file, ", fileFullPath, dstFile)
	if err != nil {
		//log.Error("internal", "open file failed: ", err.Error(), fileFullPath)
		return nil, err

	}
	s := sdk_struct.MsgStruct{}
	err = c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.File)
	if err != nil {
		return nil, err
	}
	s.FileElem.FilePath = fileFullPath
	fi, err := os.Stat(fileFullPath)
	if err != nil {
		//log.Error("internal", "get file Attributes error", err.Error())
		return nil, err
	}
	s.FileElem.FileSize = fi.Size()
	s.FileElem.FileName = fileName
	s.Content = utils.StructToJsonString(s.FileElem)
	return &s, nil
}
func (c *Conversation) CreateImageMessageFromFullPath(ctx context.Context, imageFullPath string) (*sdk_struct.MsgStruct, error) {
	dstFile := utils.FileTmpPath(imageFullPath, c.DataDir) //a->b
	_, err := utils.CopyFile(imageFullPath, dstFile)
	//log.Info(operationID, "copy file, ", imageFullPath, dstFile)
	if err != nil {
		//log.Error(operationID, "open file failed: ", err, imageFullPath)
		return nil, err
	}
	s := sdk_struct.MsgStruct{}
	err = c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Picture)
	if err != nil {
		return nil, err
	}
	s.PictureElem.SourcePath = imageFullPath
	//log.Info(operationID, "ImageMessage  path:", s.PictureElem.SourcePath)
	imageInfo, err := getImageInfo(s.PictureElem.SourcePath)
	if err != nil {
		//log.Error(operationID, "getImageInfo err:", err.Error())
		return nil, err
	}
	s.PictureElem.SourcePicture.Width = imageInfo.Width
	s.PictureElem.SourcePicture.Height = imageInfo.Height
	s.PictureElem.SourcePicture.Type = imageInfo.Type
	s.PictureElem.SourcePicture.Size = imageInfo.Size
	s.Content = utils.StructToJsonString(s.PictureElem)
	return &s, nil

}

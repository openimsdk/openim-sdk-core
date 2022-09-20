package indexdb

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

//1,使用wasm原生的方式，tinygo应用于go的嵌入式领域，支持的功能有限，甚至json序列化都无法支持
//2.函数命名遵从驼峰还是帕斯卡命名法需要确定一下
//3.提供的sql生成语句中，关于bool值需要特殊处理，create语句的设计的到bool值的我会在创建语句中单独说明，这是因为在原有的sqlite中并不支持bool，用整数1或者0替代，gorm对其做了转换。
//4.提供的sql生成语句中，字段名是下划线方式 例如：recv_id，但是接口转换的数据json tag字段的风格是recvID，类似的所有的字段需要做个map映射

func (i IndexDB) GetMessage(clientMsgID string) (*model_struct.LocalChatLog, error) {
	msg, err := Exec(clientMsgID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msg.(string); ok {
			result := model_struct.LocalChatLog{}
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return &result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i IndexDB) GetAllUnDeleteMessageSeqList() ([]uint32, error) {
	panic("implement me")
}

func (i IndexDB) UpdateColumnsMessageList(clientMsgIDList []string, args map[string]interface{}) error {
	panic("implement me")
}

func (i IndexDB) UpdateColumnsMessage(ClientMsgID string, args map[string]interface{}) error {
	panic("implement me")
}

func (i IndexDB) UpdateColumnsMessageController(ClientMsgID string, groupID string, sessionType int32, args map[string]interface{}) error {
	panic("implement me")
}

func (i IndexDB) UpdateMessage(c *model_struct.LocalChatLog) error {
	panic("implement me")
}

func (i IndexDB) UpdateMessageController(c *model_struct.LocalChatLog) error {
	panic("implement me")
}

func (i IndexDB) DeleteAllMessage() error {
	panic("implement me")
}

func (i IndexDB) UpdateMessageStatusBySourceID(sourceID string, status, sessionType int32) error {
	panic("implement me")
}

func (i IndexDB) UpdateMessageStatusBySourceIDController(sourceID string, status, sessionType int32) error {
	panic("implement me")
}

func (i IndexDB) UpdateMessageTimeAndStatus(clientMsgID string, serverMsgID string, sendTime int64, status int32) error {
	panic("implement me")
}

func (i IndexDB) UpdateMessageTimeAndStatusController(msg *sdk_struct.MsgStruct) error {
	panic("implement me")
}

func (i IndexDB) GetMessageList(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) GetMessageListController(sourceID string, sessionType, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) GetMessageListNoTime(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) GetMessageListNoTimeController(sourceID string, sessionType, count int, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) GetSendingMessageList() (result []*model_struct.LocalChatLog, err error) {
	panic("implement me")
}

func (i IndexDB) UpdateSingleMessageHasRead(sendID string, msgIDList []string) error {
	panic("implement me")
}

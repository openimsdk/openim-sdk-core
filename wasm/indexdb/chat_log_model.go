package indexdb

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
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
		result := model_struct.LocalChatLog{}
		err := utils.JsonStringToStruct(msg, &result)
		if err != nil {
			return nil, err
		}
		return &result, err

	}
}

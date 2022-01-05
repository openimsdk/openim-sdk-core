package open_im_sdk

import (
	"encoding/json"
	"runtime"
)

func (u *UserRelated) jsonUnmarshalAndArgsValidate(s string, args interface{}, callback Base) error {
	err := json.Unmarshal([]byte(s), args)
	if err != nil {
		if callback != nil {
			callback.OnError(ErrCodeConversation, wrap(err, "json Unmarshal failed").Error())
			runtime.Goexit()
		} else {
			return err
		}
	}
	err = u.validate.Struct(args)
	if err != nil {
		if callback != nil {
			callback.OnError(ErrCodeConversation, wrap(err, "args check failed").Error())
			runtime.Goexit()
		}
	}
	return err
}

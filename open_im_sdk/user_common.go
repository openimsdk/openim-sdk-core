package open_im_sdk

import (
	"encoding/json"
	"runtime"
)

func (u *UserRelated) jsonUnmarshalAndArgsValidate(s string, args interface{}, callback Base) error {
	err := json.Unmarshal([]byte(s), args)
	if err != nil {
		if callback != nil {
			callback.OnError(ErrCodeConversation, err.Error())
			runtime.Goexit()
		} else {
			return wrap(err, "json Unmarshal failed")
		}
	}
	err = u.validate.Struct(args)
	if err != nil {
		if callback != nil {
			callback.OnError(ErrCodeConversation, err.Error())
			runtime.Goexit()
		}
	}
	return wrap(err, "args check failed")
}

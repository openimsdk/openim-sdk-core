package sdkerrs

// UnknownCode 没有解析到code或解析的code=0

const NetworkError = 10000

// 通用错误码
const (
	NoError                      = 0     //无错误
	UnknownCode                  = 1000  //没有解析到code
	SdkInternalError             = 10500 //SDK内部错误
	ArgsError                    = 11001 //输入参数错误
	DuplicateKeyError            = 11003
	RecordNotFoundError          = 11004 //记录不存在
	ResourceLoadNotCompleteError = 11005 //资源初始化未完成

	UserIDNotFoundError  = 1101 //UserID不存在 或未注册
	GroupIDNotFoundError = 1201 //GroupID不存在
)

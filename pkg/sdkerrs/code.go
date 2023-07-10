package sdkerrs

// 通用错误码
const (
	NetworkError             = 10000
	NetworkTimeoutError      = 10001
	ArgsError                = 10002 //输入参数错误
	CtxDeadlineExceededError = 10003 //上下文超时

	ResourceLoadNotCompleteError = 10004 //资源初始化未完成
	UnknownCode                  = 10005 //没有解析到code
	SdkInternalError             = 10006 //SDK内部错误

	UserIDNotFoundError = 10100 //UserID不存在 或未注册
	LoginOutError       = 10101 //用户已经退出登录
	LoginRepeatError    = 10102 //用户重复登录

	//消息相关
	FileNotFoundError             = 10200 //记录不存在
	MsgDeCompressionError         = 10201 //消息解压失败
	MsgDecodeBinaryWsError        = 10202 //消息解码失败
	MsgBinaryTypeNotSupportError  = 10203 //消息类型不支持
	MsgRepeatError                = 10204 //消息重复发送
	MsgContentTypeNotSupportError = 10205 //消息类型不支持

	//会话相关
	NotSupportOptError = 10301 //不支持的操作

	//群组相关
	GroupIDNotFoundError = 10400 //GroupID不存在
	GroupTypeErr         = 10401 //群组类型错误

)

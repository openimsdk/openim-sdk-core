package sdkerrs

// 通用错误码
const (
	NetworkError                 = 10000
	NetworkTimeoutError          = 10001
	ArgsError                    = 11001 //输入参数错误
	CtxDeadlineExceededError     = 11002 //上下文超时
	DuplicateKeyError            = 11003
	RecordNotFoundError          = 11004 //记录不存在
	ResourceLoadNotCompleteError = 11005 //资源初始化未完成
	UnknownCode                  = 10006 //没有解析到code
	SdkInternalError             = 10500 //SDK内部错误

	UserIDNotFoundError  = 11101 //UserID不存在 或未注册
	GroupIDNotFoundError = 11201 //GroupID不存在
	TokenInvalidError    = 11502

	//消息相关
	MsgContentEmptyError          = 12001 //消息内容为空
	MsgDeCompressionError         = 12002 //消息解压失败
	MsgDecodeBinaryWsError        = 12003 //消息解码失败
	MsgTypeNotSupportError        = 12004 //消息类型不支持
	MsgRepeatError                = 12005 //消息重复发送
	MsgContentTypeNotSupportError = 12006 //消息类型不支持
	MsgNotFoundError              = 12007 //消息发送不成功

	//会话相关
	NotSupportOptError       = 13001 //不支持的操作
	NotResetGroupAtTypeError = 13002 //不需要重置会话
	NotFoundConversation     = 13003 //会话不存在

	//群组相关
	NotInGroupError = 14001 //不在群组中
	GroupTypeErr    = 14002 //群组类型错误

	LoginOutError    = 13001 //用户已经退出登录
	LoginRepeatError = 13002 //用户重复登录
)

package public

// 需要设置的上下文相关的一些全局 Key
const (
	CtxResponseKey = "response" // response
	CtxUserIDKey   = "userID"   // userID
	CtxUserKey     = "user"

	// 验证相关 Key
	CtxValidatorKey  = "ValidatorKey"
	CtxTranslatorKey = "TranslatorKey"
)

// Context 之外的其他全局 Key
const (
	// Session Key 用于 session 认证
	KeySessionUser = "user"
)

// requestlog 中使用
const (
	HeaderTraceID    = "com-header-rid"
	HeaderSpanID     = "com-header-spanid"
	ContextStartTime = "startExecTime"
	ContextTrace     = "trace"
)

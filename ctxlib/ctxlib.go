package ctxlib

type CtxKey string

var (
	UserInfoKey       CtxKey = "userInfo"
	RefererKey        CtxKey = "referer"
	XSaaSusRefererKey CtxKey = "xSaaSusReferer"
)

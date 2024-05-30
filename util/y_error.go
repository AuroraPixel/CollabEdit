package util

import "errors"

var (
	ErrUnexpectedCase      = errors.New("未知异常")
	ErrMethodUnimplemented = errors.New("方法没有被实现")
	ErrTypeConversion      = errors.New("类型转换错误")
	ErrParamUnimplemented  = errors.New("参数未实现")
)

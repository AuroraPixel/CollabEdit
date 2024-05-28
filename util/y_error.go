package util

import "errors"

// 定义错误类型
var (
	ErrUnexpectedCase      = errors.New("未知异常")
	ErrMethodUnimplemented = errors.New("方法没有被实现")
	ErrTypeConversion      = errors.New("类型转换错误")
)

package test

import (
	"CollabEdit/util"
	"fmt"
	"testing"
)

func TestEventHandler(t *testing.T) {
	handler := util.NewEventHandler()
	event1 := func(arg0 interface{}, arg1 interface{}) {
		fmt.Printf("监听事件 1: 参数0=%v, 参数1=%v\n", arg0, arg1)
	}
	event2 := func(arg0 interface{}, arg1 interface{}) {
		fmt.Printf("监听事件 2: 参数0=%v, 参数1=%v\n", arg0, arg1)
	}

	// 添加一个事件监听器
	handler.AddEvent(event1)
	handler.AddEvent(event2)

	// 事件发送
	handler.CallEvents("hello", 45)

	// 移除事件1
	handler.RemoveEvent(event1)

	// 事件发送
	handler.CallEvents("removeHello", 45)

	// 移除所有事件
	handler.RemoveAllEvent()
	handler.CallEvents("hello", 45)
}

package main

import (
	"CollabEdit/core"
	"fmt"
)

func main() {
	//示例
	eventBus := core.NewObservable()

	//定义观察者
	logger := func(args interface{}) {
		fmt.Println("日志系统:", args.(string))
	}
	notifier := func(args interface{}) {
		fmt.Println("通知:", args.(string))
	}

	//注册观察者
	eventBus.On("dataChanged", logger)
	eventBus.On("dataChanged", notifier)

	//发送事件
	eventBus.Emit("dataChanged", "Data has been updated.")

	//注销日志观察者
	eventBus.Off("dataChanged", logger)

	//发送事件
	eventBus.Emit("dataChanged", "Data has been updated again.")
}

package test

import (
	"CollabEdit/core"
	"testing"
)

func TestObservable(t *testing.T) {
	// 创建一个新的 observable
	eventBus := core.NewObservable()

	// 定义观察者
	var logMessage string
	logger := func(args interface{}) {
		logMessage = "日志系统:" + args.(string)
	}

	var notifyMessage string
	notifier := func(args interface{}) {
		notifyMessage = "通知:" + args.(string)
	}

	// 注册观察者
	eventBus.On("dataChanged", logger)
	eventBus.On("dataChanged", notifier)

	// 发送事件
	eventBus.Emit("dataChanged", "Data has been updated.")

	// 检查观察者是否接收到事件
	if logMessage != "日志系统:Data has been updated." {
		t.Errorf("期望日志信息为 '日志系统:Data has been updated.', 但得到 '%s'", logMessage)
	}
	if notifyMessage != "通知:Data has been updated." {
		t.Errorf("期望通知信息为 '通知:Data has been updated.', 但得到 '%s'", notifyMessage)
	}

	// 注销日志观察者
	eventBus.Off("dataChanged", logger)

	// 再次发送事件
	eventBus.Emit("dataChanged", "Data has been updated again.")

	// 检查日志观察者是否已注销，通知观察者是否仍然有效
	if logMessage != "日志系统:Data has been updated." {
		t.Errorf("期望日志信息保持为 '日志系统:Data has been updated.', 但得到 '%s'", logMessage)
	}
	if notifyMessage != "通知:Data has been updated again." {
		t.Errorf("期望通知信息为 '通知:Data has been updated again.', 但得到 '%s'", notifyMessage)
	}
}

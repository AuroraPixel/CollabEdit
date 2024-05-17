package util

import (
	"log"
	"reflect"
)

// EventHandler 通用事件处理器
type EventHandler struct {
	Events []func(arg0 interface{}, arg1 interface{})
}

// NewEventHandler 创建新的EventHandler实力
func NewEventHandler() *EventHandler {
	return &EventHandler{
		Events: make([]func(arg0 interface{}, arg1 interface{}), 0),
	}
}

// AddEvent 添加一个事件
func (eh *EventHandler) AddEvent(event func(arg0 interface{}, arg1 interface{})) {
	eh.Events = append(eh.Events, event)
}

// RemoveEvent 移除一个事件
func (eh *EventHandler) RemoveEvent(event func(arg0 interface{}, arg1 interface{})) {
	newEvents := make([]func(arg0 interface{}, arg1 interface{}), 0)
	for _, l := range eh.Events {
		// 通过反射比较函数指针
		lPtr := reflect.ValueOf(l).Pointer()
		eventPtr := reflect.ValueOf(event).Pointer()
		if lPtr != eventPtr {
			newEvents = append(newEvents, l)
		}
	}
	if len(newEvents) == len(eh.Events) {
		log.Fatalf("[CollabEdit] Tried to remove event handler that doesn't exist.")
	}
	eh.Events = newEvents
}

// RemoveAllEvent 移除所有事件
func (eh *EventHandler) RemoveAllEvent() {
	eh.Events = make([]func(arg0 interface{}, arg1 interface{}), 0)
}

// CallEvents 调用所有事件监听器
func (eh *EventHandler) CallEvents(arg0 interface{}, arg1 interface{}) {
	for _, event := range eh.Events {
		event(arg0, arg1)
	}
}

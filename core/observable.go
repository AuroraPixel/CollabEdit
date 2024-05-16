package core

import (
	"reflect"
	"sync"
)

// Observable 观察者结构
type Observable struct {
	observers map[string][]func(args interface{})
	mu        sync.Mutex
}

// NewObservable 初始化新观察者
func NewObservable() *Observable {
	return &Observable{
		observers: make(map[string][]func(args interface{})),
	}
}

// On 注册观察者
func (o *Observable) On(eventName string, f func(args interface{})) {
	//上锁
	o.mu.Lock()
	defer o.mu.Unlock()
	o.observers[eventName] = append(o.observers[eventName], f)
}

// Off 注销观察者
func (o *Observable) Off(eventName string, f func(args interface{})) {
	o.mu.Lock()
	defer o.mu.Unlock()
	observers := o.observers[eventName]
	for i, observer := range observers {
		if funcEqual(observer, f) {
			o.observers[eventName] = append(observers[:i], observers[i+1:]...)
			break
		}
	}
	if len(o.observers[eventName]) == 0 {
		delete(o.observers, eventName)
	}
}

// Emit 事件触发，通知所有注册的观察者
func (o *Observable) Emit(eventName string, args interface{}) {
	o.mu.Lock()
	observers := append([]func(args interface{}){}, o.observers[eventName]...)
	o.mu.Unlock()
	for _, observer := range observers {
		observer(args)
	}
}

// 比较两个函数是否相等
func funcEqual(a, b interface{}) bool {
	return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}

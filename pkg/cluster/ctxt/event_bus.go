/*
Copyright © 2020 Marvin

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package ctxt

import (
	"fmt"

	ev "github.com/asaskevich/EventBus"

	"go.uber.org/zap"
)

// EventBus 是任务事件的事件总线.
type EventBus struct {
	eventBus ev.Bus
}

// EventKind 代表事件类型.
type EventKind string

const (
	// 在任务将要执行时发出 EventTaskBegin.
	EventTaskBegin EventKind = "task_begin"
	// 任务完成执行时发出 EventTaskFinish.
	EventTaskFinish EventKind = "task_finish"
	// 当任务取得一些进展时，会发出 EventTaskProgress
	EventTaskProgress EventKind = "task_progress"
)

// NewEventBus 创建事件总线.
func NewEventBus() EventBus {
	return EventBus{
		eventBus: ev.New(),
	}
}

// PublishTaskBegin 发布一个 TaskBegin 事件。这只能由并行或串行调用
func (ev *EventBus) PublishTaskBegin(task fmt.Stringer) {
	zap.L().Debug("TaskBegin", zap.String("task", task.String()))
	ev.eventBus.Publish(string(EventTaskBegin), task)
}

// PublishTaskFinish 发布 TaskFinish 事件。这只能由并行或串行调用
func (ev *EventBus) PublishTaskFinish(task fmt.Stringer, err error) {
	zap.L().Debug("TaskFinish", zap.String("task", task.String()), zap.Error(err))
	ev.eventBus.Publish(string(EventTaskFinish), task, err)
}

// PublishTaskProgress 发布一个 TaskProgress 事件
func (ev *EventBus) PublishTaskProgress(task fmt.Stringer, progress string) {
	zap.L().Debug("TaskProgress", zap.String("task", task.String()), zap.String("progress", progress))
	ev.eventBus.Publish(string(EventTaskProgress), task, progress)
}

// Subscribe 订阅事件.
func (ev *EventBus) Subscribe(eventName EventKind, handler interface{}) error {
	err := ev.eventBus.Subscribe(string(eventName), handler)
	if err != nil {
		zap.L().Debug("TaskSubscribe", zap.String("error", err.Error()))
		return fmt.Errorf("TaskSubscribe appear error: %v", err)
	}
	return nil
}

// Unsubscribe 取消订阅事件.
func (ev *EventBus) Unsubscribe(eventName EventKind, handler interface{}) error {
	err := ev.eventBus.Unsubscribe(string(eventName), handler)
	if err != nil {
		zap.L().Debug("TaskUnsubscribe", zap.String("error", err.Error()))
		return fmt.Errorf("TaskUnsubscribe appear error: %v", err)
	}
	return nil
}

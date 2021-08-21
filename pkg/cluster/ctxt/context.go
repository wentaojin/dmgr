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
	"sync"

	"github.com/wentaojin/dmgr/pkg/cluster/executor"
	"github.com/wentaojin/dmgr/pkg/cluster/mock"
)

// 上下文用于在多个任务执行时共享状态。
// 使用互斥锁来防止某些字段的并发读/写
// 因为可以在并行任务中共享相同的上下文。
type Context struct {
	Ev EventBus

	Exec struct {
		sync.RWMutex
		Executors    map[string]executor.Executor
		Stdouts      map[string][]byte
		Stderrs      map[string][]byte
		CheckResults map[string][]interface{}
	}

	// 私钥/公钥用于通过用户访问远程服务器
	PrivateKeyPath string
	PublicKeyPath  string
}

// NewContext create a context instance.
func NewContext() *Context {
	return &Context{
		Ev: NewEventBus(),
		Exec: struct {
			sync.RWMutex
			Executors    map[string]executor.Executor
			Stdouts      map[string][]byte
			Stderrs      map[string][]byte
			CheckResults map[string][]interface{}
		}{
			Executors:    make(map[string]executor.Executor),
			Stdouts:      make(map[string][]byte),
			Stderrs:      make(map[string][]byte),
			CheckResults: make(map[string][]interface{}),
		},
	}
}

// SetExecutor set the executor.
func (ctx *Context) SetExecutor(host string, e executor.Executor) {
	ctx.Exec.Lock()
	ctx.Exec.Executors[host] = e
	ctx.Exec.Unlock()
}

// GetExecutor get the executor.
func (ctx *Context) GetExecutor(host string) (e executor.Executor, ok bool) {
	// Mock point for unit test
	if e := mock.On("FakeExecutor"); e != nil {
		return e.(executor.Executor), true
	}

	ctx.Exec.RLock()
	e, ok = ctx.Exec.Executors[host]
	ctx.Exec.RUnlock()
	return
}

// GetSSHKeySet implements the operation.ExecutorGetter interface.
func (ctx *Context) GetSSHKeySet() (privateKeyPath, publicKeyPath string) {
	return ctx.PrivateKeyPath, ctx.PublicKeyPath
}

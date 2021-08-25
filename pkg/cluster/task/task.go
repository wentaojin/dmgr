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
package task

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/wentaojin/dmgr/pkg/cluster/ctxt"

	"github.com/wentaojin/dmgr/pkg/dmgrutil"
	"go.uber.org/zap"

	"github.com/joomcode/errorx"
)

var (
	// 来自某逻辑区域错误
	errNS = errorx.NewNamespace("task")

	// ErrUnsupportedRollback 表示任务不支持回滚
	ErrUnsupportedRollback = errors.New("unsupported rollback")
	// ErrNoExecutor 表示无法获得 SSH executor 执行者
	ErrNoExecutor = errors.New("no executor")
	// ErrNoOutput 表示无法获得主机的输出
	ErrNoOutput = errors.New("no outputs available")
)

type (
	// 	Task 接口
	Task interface {
		fmt.Stringer
		Execute(ctx *ctxt.Context) error
		Rollback(ctx *ctxt.Context) error
	}

	// Serial 会以序列化的方式执行一组任务
	Serial struct {
		ignoreError       bool
		hideDetailDisplay bool
		inner             []Task
	}

	// Parallel 会以并行方式执行一组任务
	Parallel struct {
		ignoreError       bool
		hideDetailDisplay bool
		inner             []Task
	}
)

// Execute implements the Task interface
func (s *Serial) Execute(ctx *ctxt.Context) error {
	for _, t := range s.inner {
		if !s.hideDetailDisplay {
			dmgrutil.Logger.Info("Serial", zap.String("msg", t.String()))
		}
		err := t.Execute(ctx)
		if err != nil && !s.ignoreError {
			return err
		}
	}
	return nil
}

// Rollback implements the Task interface
func (s *Serial) Rollback(ctx *ctxt.Context) error {
	// Rollback in reverse order
	for i := len(s.inner) - 1; i >= 0; i-- {
		err := s.inner[i].Rollback(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// String implements the fmt.Stringer interface
func (s *Serial) String() string {
	var ss []string
	for _, t := range s.inner {
		ss = append(ss, t.String())
	}
	return strings.Join(ss, "\n")
}

// Execute implements the Task interface
func (p *Parallel) Execute(ctx *ctxt.Context) error {
	var firstError error
	var mu sync.Mutex
	wg := sync.WaitGroup{}
	for _, t := range p.inner {
		wg.Add(1)
		go func(t Task, logger *zap.Logger) {
			defer wg.Done()

			if !p.hideDetailDisplay {
				logger.Info("Parallel", zap.String("msg", t.String()))
			}
			err := t.Execute(ctx)
			if err != nil {
				mu.Lock()
				if firstError == nil {
					firstError = err
				}
				mu.Unlock()
			}
		}(t, dmgrutil.Logger)
	}
	wg.Wait()
	if p.ignoreError {
		return nil
	}
	return firstError
}

// Rollback implements the Task interface
func (p *Parallel) Rollback(ctx *ctxt.Context) error {
	var firstError error
	var mu sync.Mutex
	wg := sync.WaitGroup{}
	for _, t := range p.inner {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			err := t.Rollback(ctx)
			if err != nil {
				mu.Lock()
				if firstError == nil {
					firstError = err
				}
				mu.Unlock()
			}
		}(t)
	}
	wg.Wait()
	return firstError
}

// String implements the fmt.Stringer interface
func (p *Parallel) String() string {
	var ss []string
	for _, t := range p.inner {
		ss = append(ss, t.String())
	}
	return strings.Join(ss, "\n")
}

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
	"fmt"

	"github.com/wentaojin/dmgr/pkg/cluster/ctxt"

	"github.com/wentaojin/dmgr/pkg/dmgrutil"
	"go.uber.org/zap"

	"github.com/wentaojin/dmgr/pkg/cluster/module"
)

// EnableInstance 用于启动组件实例
type EnableInstance struct {
	host           string
	servicePort    uint64
	instanceName   string
	logDir         string
	serviceName    string // 服务名
	executeTimeout uint64 // 通过 SSH 连接时超时（以秒为单位）
	isEnable       bool
}

// Execute implements the Task interface
func (e *EnableInstance) Execute(ctx *ctxt.Context) error {
	exec, found := ctx.GetExecutor(e.host)
	if !found {
		return ErrNoExecutor
	}
	action := module.OperatorDisable
	if e.isEnable {
		action = module.OperatorEnable
	}
	if err := systemctl(exec, e.serviceName, action, e.executeTimeout); err != nil {
		return toFailedActionError(err, action, e.instanceName, e.serviceName, e.logDir)
	}
	dmgrutil.Logger.Info("Enable/Disable instance success",
		zap.String("instance", e.instanceName), zap.Bool("enabled", e.isEnable))
	return nil
}

// Rollback implements the Task interface
func (e *EnableInstance) Rollback(ctx *ctxt.Context) error {
	ctx.SetExecutor(e.host, nil)
	return nil
}

// String implements the fmt.Stringer interface
func (e *EnableInstance) String() string {
	return fmt.Sprintf("EnableInstance: host=%s, serviceName=%s, isEnable=%v", e.host, e.serviceName, e.isEnable)
}

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

	"go.uber.org/zap"

	"github.com/wentaojin/dmgr/pkg/dmgrutil"

	"github.com/wentaojin/dmgr/pkg/cluster/module"
)

// StartInstance 用于启动组件实例
type StartInstance struct {
	host           string
	servicePort    uint64
	instanceName   string
	logDir         string
	serviceName    string // 服务名
	executeTimeout uint64 // 通过 SSH 连接时超时（以秒为单位）
}

// Execute implements the Task interface
func (s *StartInstance) Execute(ctx *ctxt.Context) error {
	exec, found := ctx.GetExecutor(s.host)
	if !found {
		return ErrNoExecutor
	}
	if err := systemctl(exec, s.serviceName, module.OperatorStart, s.executeTimeout); err != nil {
		return toFailedActionError(err, module.OperatorStart, s.host, s.instanceName, s.serviceName, s.logDir)
	}

	// Check ready.
	if err := PortStarted(exec, s.servicePort, s.executeTimeout); err != nil {
		return toFailedActionError(err, module.OperatorStart, s.host, s.instanceName, s.serviceName, s.logDir)
	}
	dmgrutil.Logger.Info("Start host instance success",
		zap.String("host", s.host),
		zap.Uint64("service-port", s.servicePort),
		zap.String("instance", s.instanceName))
	return nil
}

// Rollback implements the Task interface
func (s *StartInstance) Rollback(ctx *ctxt.Context) error {
	ctx.SetExecutor(s.host, nil)
	return nil
}

// String implements the fmt.Stringer interface
func (s *StartInstance) String() string {
	return fmt.Sprintf("StartInstance: host=%s, serviceName=%s", s.host, s.serviceName)
}

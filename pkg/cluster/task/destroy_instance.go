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
	"path/filepath"
	"strings"

	"github.com/wentaojin/dmgr/pkg/cluster/ctxt"

	"github.com/pingcap/errors"

	"github.com/wentaojin/dmgr/pkg/cluster/module"
	"github.com/wentaojin/dmgr/pkg/dmgrutil"
	"go.uber.org/zap"
)

// DestroyInstance 用于启动组件实例
type DestroyInstance struct {
	host           string
	servicePort    uint64
	componentName  string
	instanceName   string
	deployDir      string
	dataDir        string
	logDir         string
	executeTimeout uint64 // 通过 SSH 连接时超时（以秒为单位）
}

// Execute implements the Task interface
func (s *DestroyInstance) Execute(ctx *ctxt.Context) error {
	dmgrutil.Logger.Info("Destroying component instance", zap.String("component", s.componentName), zap.String("instance", s.instanceName))

	exec, found := ctx.GetExecutor(s.host)
	if !found {
		return ErrNoExecutor
	}

	// would save parent dir
	delPaths := dmgrutil.NewStringSet()
	delPaths.Insert(fmt.Sprintf(`rm -rf %s`, filepath.Join(s.deployDir, s.instanceName)))
	delPaths.Insert(fmt.Sprintf(`rm -rf %s`, filepath.Join(s.dataDir, s.instanceName)))
	delPaths.Insert(fmt.Sprintf(`rm -rf %s`, filepath.Join(s.logDir, s.instanceName)))
	delPaths.Insert(fmt.Sprintf("/etc/systemd/system/%s-%d.service", s.componentName, s.servicePort))
	c := module.ShellModuleConfig{
		Command:  fmt.Sprintf("rm -rf %s;", strings.Join(delPaths.Slice(), " ")),
		Sudo:     true, // the .service files are in a directory owned by root
		Chdir:    "",
		UseShell: false,
	}
	shell := module.NewShellModule(c)
	_, stderr, err := shell.Execute(exec)

	if len(stderr) > 0 {
		dmgrutil.Logger.Error("Destroying component instance failed", zap.String("component", s.componentName), zap.String("instance", s.instanceName), zap.String("error", string(stderr)))
		return errors.Annotatef(err, "failed to destroy component instance: %s", zap.String("component", s.componentName), zap.String("instance", s.instanceName))
	}

	if err != nil {
		dmgrutil.Logger.Error("Destroying component instance failed", zap.String("component", s.componentName), zap.String("instance", s.instanceName), zap.Error(err))
		return errors.Annotatef(err, "failed to destroy component instance: %s", zap.String("component", s.componentName), zap.String("instance", s.instanceName))
	}

	dmgrutil.Logger.Info("Destroying component instance success", zap.String("component", s.componentName), zap.String("instance", s.instanceName))
	return nil
}

// Rollback implements the Task interface
func (s *DestroyInstance) Rollback(ctx *ctxt.Context) error {
	ctx.SetExecutor(s.host, nil)
	return nil
}

// String implements the fmt.Stringer interface
func (s *DestroyInstance) String() string {
	return fmt.Sprintf("DestroyInstance: host=%s, instanceName=%s", s.host, s.instanceName)
}

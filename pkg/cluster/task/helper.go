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
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wentaojin/dmgr/pkg/cluster/ctxt"

	"go.uber.org/zap"

	"github.com/pingcap/errors"

	"github.com/wentaojin/dmgr/pkg/dmgrutil"

	"github.com/wentaojin/dmgr/pkg/cluster/executor"
	"github.com/wentaojin/dmgr/pkg/cluster/module"
)

// PortStarted wait until a port is being listened
func PortStarted(e executor.Executor, port uint64, timeout uint64) error {
	c := module.WaitForConfig{
		Port:    int(port),
		State:   "started",
		Timeout: time.Second * time.Duration(timeout),
	}
	w := module.NewWaitFor(c)
	return w.Execute(e)
}

// PortStopped wait until a port is being released
func PortStopped(e executor.Executor, port uint64, timeout uint64) error {
	c := module.WaitForConfig{
		Port:    int(port),
		State:   "stopped",
		Timeout: time.Second * time.Duration(timeout),
	}
	w := module.NewWaitFor(c)
	return w.Execute(e)
}

// DeletePublicKey deletes the SSH public key from host
func DeletePublicKey(ctx *ctxt.Context, host string) error {
	e, exist := ctx.GetExecutor(host)
	if !exist {
		return ErrNoExecutor
	}
	dmgrutil.Logger.Info("Delete public key", zap.String("host", host))
	_, pubKeyPath := ctx.GetSSHKeySet()
	publicKey, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return errors.Trace(err)
	}

	pubKey := string(bytes.TrimSpace(publicKey))
	pubKey = strings.ReplaceAll(pubKey, "/", "\\/")
	pubKeysFile := executor.FindSSHAuthorizedKeysFile(e)

	// delete the public key with Linux `sed` toolkit
	c := module.ShellModuleConfig{
		Command:  fmt.Sprintf("sed -i '/%s/d' %s", pubKey, pubKeysFile),
		UseShell: false,
	}
	shell := module.NewShellModule(c)
	stdout, stderr, err := shell.Execute(e)

	if len(stdout) > 0 {
		dmgrutil.Logger.Info(string(stdout))
	}
	if len(stderr) > 0 {
		dmgrutil.Logger.Error("Delete public key failed", zap.String("host", host), zap.String("error", string(stderr)))
		return errors.Annotatef(err, "failed to delete public key", zap.String("host", host))
	}

	if err != nil {
		dmgrutil.Logger.Error("Delete public key failed", zap.String("host", host), zap.Error(err))
		return errors.Annotatef(err, "failed to delete pulblic key on: %s", host)
	}
	dmgrutil.Logger.Info("Delete public key success", zap.String("host", host))
	return nil
}

// 服务启动
func systemctl(executor executor.Executor, service string, action string, timeout uint64) error {
	c := module.SystemdModuleConfig{
		Unit:           service,
		ReloadDaemon:   true,
		Action:         action,
		ExecuteTimeout: time.Second * time.Duration(timeout),
	}
	systemd := module.NewSystemdModule(c)
	stdout, stderr, err := systemd.Execute(executor)

	if len(stdout) > 0 {
		dmgrutil.Logger.Warn("Systemctl", zap.String("Stdout", string(stdout)))
	}
	if len(stderr) > 0 && !bytes.Contains(stderr, []byte("Created symlink ")) && !bytes.Contains(stderr, []byte("Removed symlink ")) {
		dmgrutil.Logger.Error(string(stderr))
		return fmt.Errorf("host [%v] systemctl action [%s] service [%v] failed: %v", executor, action, service, string(stderr))
	}
	if len(stderr) > 0 && action == "stop" {
		// ignore "unit not loaded" error, as this means the unit is not
		// exist, and that's exactly what we want
		// NOTE: there will be a potential bug if the unit name is set
		// wrong and the real unit still remains started.
		if bytes.Contains(stderr, []byte(" not loaded.")) {
			dmgrutil.Logger.Warn(string(stderr))
			return nil // reset the error to avoid exiting
		}
		dmgrutil.Logger.Warn(string(stderr))
	}
	return err
}

// toFailedActionError formats the errror msg for failed action
func toFailedActionError(err error, action string, host, instance, service, logDir string) error {
	return errors.Annotatef(err,
		"failed to %s: %s %s %s, please check the instance's log(%s) for more detail.",
		action, host, instance, service, logDir,
	)
}

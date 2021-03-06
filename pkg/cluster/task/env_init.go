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
	"io/ioutil"
	"strings"

	"github.com/wentaojin/dmgr/pkg/cluster/executor"

	"github.com/joomcode/errorx"
	"github.com/wentaojin/dmgr/pkg/cluster/ctxt"
	"github.com/wentaojin/dmgr/pkg/cluster/module"
)

var (
	errNSEnvInit               = errNS.NewSubNamespace("env_init")
	errEnvInitSubCommandFailed = errNSEnvInit.NewType("sub_command_failed")
	// ErrEnvInitFailed is ErrEnvInitFailed
	ErrEnvInitFailed = errNSEnvInit.NewType("failed")
)

// EnvInit 用于初始化远程环境，例如：
// 1.创建集群管理用户
// 2.授权集群管理用户
type EnvInit struct {
	host           string
	clusterUser    string
	userGroup      string
	skipCreateUser bool
}

// Execute implements the Task interface
func (e *EnvInit) Execute(ctx *ctxt.Context) error {
	wrapError := func(err error) *errorx.Error {
		return ErrEnvInitFailed.Wrap(err, "Failed to initialize DM environment on remote host '%s'", e.host)
	}

	exec, found := ctx.GetExecutor(e.host)
	if !found {
		return wrapError(ErrNoExecutor)
	}

	if !e.skipCreateUser {
		um := module.NewUserModule(module.UserModuleConfig{
			Action: module.UserActionAdd,
			Name:   e.clusterUser,
			Group:  e.userGroup,
			Sudoer: true,
		})

		_, _, errx := um.Execute(exec)
		if errx != nil {
			return wrapError(errx)
		}
	}
	pubKey, err := ioutil.ReadFile(ctx.PublicKeyPath)
	if err != nil {
		return wrapError(err)
	}

	// clusterUser Authorize(PublicKeyPath)
	cmd := fmt.Sprintf(`su - %[1]s -c 'mkdir -p ~/.ssh && chmod 700 ~/.ssh'`, e.clusterUser)
	_, stderr, err := exec.Execute(cmd, true)
	if err != nil || len(stderr) > 0 {
		return wrapError(errEnvInitSubCommandFailed.
			Wrap(fmt.Errorf("error: %v, stderr: %v", err, string(stderr)), "Failed to create '~/.ssh' directory for user '%s'", e.clusterUser))
	}

	pk := strings.Fields(string(pubKey))
	sshAuthorizedKeys := executor.FindSSHAuthorizedKeysFile(exec)
	cmd = fmt.Sprintf(`su - %[1]s -c 'grep %[2]s %[3]s | wc -l'`,
		e.clusterUser, pk[1], sshAuthorizedKeys)
	stdout, stderr, err := exec.Execute(cmd, true)
	// 忽略 sshAuthorizedKeys 文件不存在错误
	if err != nil {
		return wrapError(errEnvInitSubCommandFailed.
			Wrap(err, "Failed to grep '~/.ssh' directory for user '%s'", e.clusterUser))
	}

	if strings.Replace(string(stdout), "\n", "", -1) == "0" {
		cmd = fmt.Sprintf(`su - %[1]s -c 'echo %[2]s >> %[3]s && chmod 600 %[3]s'`,
			e.clusterUser, strings.Replace(string(pubKey), "\n", "", -1), sshAuthorizedKeys)
		_, stderr, err = exec.Execute(cmd, true)
		if err != nil || len(stderr) > 0 {
			return wrapError(errEnvInitSubCommandFailed.
				Wrap(fmt.Errorf("error: %v, stderr: %v", err, string(stderr)), "Failed to write public keys to '%s' for user '%s'", sshAuthorizedKeys, e.clusterUser))
		}
	}

	return nil
}

// Rollback implements the Task interface
func (e *EnvInit) Rollback(ctx *ctxt.Context) error {
	return ErrUnsupportedRollback
}

// String implements the fmt.Stringer interface
func (e *EnvInit) String() string {
	return fmt.Sprintf("EnvInit: user=%s, host=%s", e.clusterUser, e.host)
}

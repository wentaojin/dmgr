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
package executor

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pingcap/errors"

	"github.com/joomcode/errorx"

	"github.com/wentaojin/dmgr/pkg/dmgrutil"

	"go.uber.org/zap"

	"github.com/appleboy/easyssh-proxy"
)

const (
	// 默认 SSH 执行命令超时时间，单位：秒
	DefaultExecuteTimeout = 60
	// 默认 SSH 连接超时时间，单位：秒
	DefaultConnectTimeout = 10
)

var (
	errNSSSH              = errNS.NewSubNamespace("ssh")
	ErrPropSSHCommand     = errorx.RegisterPrintableProperty("ssh_command")
	ErrPropSSHStdout      = errorx.RegisterPrintableProperty("ssh_stdout")
	ErrPropSSHStderr      = errorx.RegisterPrintableProperty("ssh_stderr")
	ErrSSHExecuteFailed   = errNSSSH.NewType("execute_failed")
	ErrSSHExecuteTimedout = errNSSSH.NewType("execute_timeout")
)

// EasySSHExecutor 实现 EasySSH Executor 作为 SSH 传输协议层
type (
	EasySSHExecutor struct {
		Config *easyssh.MakeConfig
		Locale string // 执行命令时使用的语言环境
		Sudo   bool   // 使用此执行程序运行的所有命令是否使用 sudo
	}

	// SSHConfig 是建立 SSH 连接所需的配置
	SSHConfig struct {
		Host           string        // SSH 服务器的主机名
		Port           int           // SSH 服务器的主机端口
		User           string        // SSH 服务器的用户名
		Password       string        // SSH 服务器的用户密码
		KeyFile        string        // SSH 私有密钥文件
		Passphrase     string        // SSH 私有密钥密码
		ConnectTimeout time.Duration // TCP 连接建立的最长时间
		ExecuteTimeout time.Duration // 命令完成的最长时间

	}
)

var _ Executor = &EasySSHExecutor{}

// 通过 SSH 执行运行命令，默认情况下它不调用任何特定的 shell
func (e *EasySSHExecutor) Execute(cmd string, sudo bool, execTimeout ...time.Duration) ([]byte, []byte, error) {
	// 尝试获取 root 权限
	if e.Sudo || sudo {
		cmd = fmt.Sprintf("sudo -H bash -c \"%s\"", cmd)
	}

	//设置一个基本的 PATH 以防在登录时为空
	cmd = fmt.Sprintf("PATH=$PATH:/usr/bin:/usr/sbin %s", cmd)

	if e.Locale != "" {
		cmd = fmt.Sprintf("export LANG=%s; %s", e.Locale, cmd)
	}

	// 在远程主机上运行命令
	// easyssh-proxy 中的默认超时时间为 60 秒
	if len(execTimeout) == 0 {
		execTimeout = append(execTimeout, time.Second*DefaultExecuteTimeout)
	}

	fmt.Println(e.Config)

	stdout, stderr, done, err := e.Config.Run(cmd, execTimeout...)
	dmgrutil.Logger.Info("SSHCommand",
		zap.String("host", e.Config.Server),
		zap.String("port", e.Config.Port),
		zap.String("cmd", cmd),
		zap.Error(err),
		zap.String("stdout", stdout),
		zap.String("stderr", stderr))
	if err != nil {
		sshErr := ErrSSHExecuteFailed.
			Wrap(err, "Failed to execute command over SSH for '%s@%s:%s'", e.Config.User, e.Config.Server, e.Config.Port).
			WithProperty(ErrPropSSHCommand, cmd).
			WithProperty(ErrPropSSHStdout, stdout).
			WithProperty(ErrPropSSHStderr, stderr)
		if len(stdout) > 0 || len(stderr) > 0 {
			output := strings.TrimSpace(strings.Join([]string{stdout, stderr}, "\n"))
			sshErr = sshErr.
				WithProperty(
					errorx.RegisterPrintableProperty(
						fmt.Sprintf("Command output on remote host %s", e.Config.Server)),
					output)
		}
		return []byte(stdout), []byte(stderr), sshErr
	}
	// 执行超时
	if !done {
		return []byte(stdout), []byte(stderr), ErrSSHExecuteTimedout.
			Wrap(err, "Execute command over SSH timedout for '%s@%s:%s'", e.Config.User, e.Config.Server, e.Config.Port).
			WithProperty(ErrPropSSHCommand, cmd).
			WithProperty(ErrPropSSHStdout, stdout).
			WithProperty(ErrPropSSHStderr, stderr)
	}

	return []byte(stdout), []byte(stderr), nil
}

// 通过 SCP 传输副本文件
// 此函数依赖于 `scp`（来自 OpenSSH 或其他 SSH 实现的工具）
// 该函数基于　easyssh.MakeConfig.Scp()　但支持复制从远程到本地的文件
func (e *EasySSHExecutor) Transfer(src, dst string, download bool, limit int) error {
	if !download {
		err := e.Config.Scp(src, dst)
		if err != nil {
			return errors.Annotatef(err, "failed to scp %s to %s@%s:%s", src, e.Config.User, e.Config.Server, dst)
		}
		return nil
	}

	// download file from remote
	session, client, err := e.Config.Connect()
	if err != nil {
		return err
	}
	defer client.Close()
	defer session.Close()

	return ScpDownload(session, client, src, dst, limit)
}

//初始化构建并初始化一个 EasySSHExecutor
func (e *EasySSHExecutor) initialize(config SSHConfig) {
	// 创建 easyssh 配置
	e.Config = &easyssh.MakeConfig{
		Server:  config.Host,
		Port:    strconv.Itoa(config.Port),
		User:    config.User,
		Timeout: config.ConnectTimeout,
	}

	// 有私钥，优先使用私钥认证
	if len(config.KeyFile) > 0 {
		e.Config.KeyPath = config.KeyFile
		e.Config.Passphrase = config.Passphrase
	} else if len(config.Password) > 0 {
		e.Config.Password = config.Password
	}
}

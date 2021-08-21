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
	"time"

	"github.com/wentaojin/dmgr/pkg/cluster/ctxt"

	"github.com/wentaojin/dmgr/pkg/cluster/executor"
)

// RootSSH 用于使用特定密钥建立到目标主机的 SSH 连接
type RootSSH struct {
	host           string // SSH 服务器的主机名
	port           uint64 // SSH 服务器的主机端口
	user           string // SSH 服务器的登录用户名
	password       string // SSH 服务器的登录用户密码
	keyFile        string // 私钥文件的路径
	passphrase     string // 私钥文件的密码
	connectTimeout uint64 // 通过 SSH 连接时超时（以秒为单位）
	executeTimeout uint64 // 以秒为单位的超时等待命令完成
}

// Execute implements the Task interface
func (s *RootSSH) Execute(ctx *ctxt.Context) error {
	sc := executor.SSHConfig{
		Host:           s.host,
		Port:           int(s.port),
		User:           s.user,
		Password:       s.password,
		KeyFile:        s.keyFile,
		Passphrase:     s.passphrase,
		ConnectTimeout: time.Duration(s.connectTimeout) * time.Second,
		ExecuteTimeout: time.Duration(s.executeTimeout) * time.Second,
	}

	e, err := executor.NewSSHExecutor(s.user != "root", sc)
	if err != nil {
		return err
	}

	ctx.SetExecutor(s.host, e)
	return nil
}

// Rollback implements the Task interface
func (s *RootSSH) Rollback(ctx *ctxt.Context) error {
	ctx.Exec.Lock()
	delete(ctx.Exec.Executors, s.host)
	ctx.Exec.Unlock()
	return nil
}

// String implements the fmt.Stringer interface
func (s *RootSSH) String() string {
	if len(s.keyFile) > 0 {
		return fmt.Sprintf("RootSSH: user=%s, host=%s, port=%d, key=%s", s.user, s.host, s.port, s.keyFile)
	}
	return fmt.Sprintf("RootSSH: user=%s, host=%s, port=%d", s.user, s.host, s.port)
}

// UserSSH 用于使用特定密钥建立到目标主机的 SSH 连接
type UserSSH struct {
	host           string
	port           uint64
	clusterUser    string
	connectTimeout uint64 // 通过 SSH 连接时超时（以秒为单位）
	executeTimeout uint64 // 以秒为单位的超时等待命令完成
}

// Execute implements the Task interface
func (s *UserSSH) Execute(ctx *ctxt.Context) error {
	sc := executor.SSHConfig{
		Host:           s.host,
		Port:           int(s.port),
		User:           s.clusterUser,
		ConnectTimeout: time.Duration(s.connectTimeout) * time.Second,
		ExecuteTimeout: time.Duration(s.executeTimeout) * time.Second,
	}

	e, err := executor.NewSSHExecutor(false, sc)
	if err != nil {
		return err
	}

	ctx.SetExecutor(s.host, e)
	return nil
}

// Rollback implements the Task interface
func (s *UserSSH) Rollback(ctx *ctxt.Context) error {
	ctx.Exec.Lock()
	delete(ctx.Exec.Executors, s.host)
	ctx.Exec.Unlock()
	return nil
}

// String implements the fmt.Stringer interface
func (s *UserSSH) String() string {
	return fmt.Sprintf("UserSSH: user=%s, host=%s", s.clusterUser, s.host)
}

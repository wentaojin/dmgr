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
package module

import (
	"fmt"

	"github.com/wentaojin/dmgr/pkg/cluster/executor"
)

const (
	defaultShell = "/bin/bash"

	// UserActionAdd 创建用户以及删除用户
	UserActionAdd = "add"
	UserActionDel = "del"
	// UserActionModify = "modify"

	// TODO: 在 RHEL/CentOS 中，命令在 /usr/sbin 中，但在其他一些发行版中，它们可能在其他位置，例如 /usr/bin，我们将来需要检查并找到命令的正确路径
	useraddCmd  = "/usr/sbin/useradd"
	userdelCmd  = "/usr/sbin/userdel"
	groupaddCmd = "/usr/sbin/groupadd"
	// usermodCmd = "/usr/sbin/usermod"
)

var (
	// 来自某逻辑区域错误
	errNSUser = errNS.NewSubNamespace("user")
	// 来自创建用户以及删除用户失败错误
	ErrUserAddFailed    = errNSUser.NewType("user_add_failed")
	ErrUserDeleteFailed = errNSUser.NewType("user_delete_failed")
)

// UserModuleConfig 是用于初始化 UserModule 的配置
type UserModuleConfig struct {
	Action string // 创建、删除或者更改用户
	Name   string // 用户名
	Group  string // 用户组名
	Home   string // 用户家目录
	Shell  string // 用户的登录 shell
	Sudoer bool   // 当为真时，用户将被添加到 sudoers 列表
}

// UserModule 是用于控制 systemd 单元的模块
type UserModule struct {
	config UserModuleConfig
	cmd    string // 待执行的 Shell 命令
}

// NewUserModule 基于给定的配置构建并返回一个 UserModule 对象
func NewUserModule(config UserModuleConfig) *UserModule {
	cmd := ""

	switch config.Action {
	case UserActionAdd:
		cmd = useraddCmd
		// 必须使用 -m，否则不会创建主目录。如果要指定 home 目录的路径，使用 -d 并指定路径
		// useradd -m -d /PATH/TO/FOLDER
		cmd += " -m"
		if config.Home != "" {
			cmd += " -d" + config.Home
		}

		//设置用户的登录 shell
		if config.Shell != "" {
			cmd = fmt.Sprintf("%s -s %s", cmd, config.Shell)
		} else {
			cmd = fmt.Sprintf("%s -s %s", cmd, defaultShell)
		}

		//设置用户组
		if config.Group == "" {
			config.Group = config.Name
		}

		// groupadd -f <group-name>
		groupAdd := fmt.Sprintf("%s -f %s", groupaddCmd, config.Group)

		// useradd -g <group-name> <user-name>
		cmd = fmt.Sprintf("%s -g %s %s", cmd, config.Group, config.Name)

		// chown privilege and group
		var chownCmd string
		if config.Home != "" {
			chownCmd = fmt.Sprintf("chown %s:%s %s", config.Name, config.Group, config.Home)
		} else {
			chownCmd = fmt.Sprintf("chown %s:%s %s", config.Name, config.Group, fmt.Sprintf("/home/%s", config.Name))
		}

		//防止用户名已被使用时出错
		cmd = fmt.Sprintf("id -u %s > /dev/null 2>&1 || (%s && %s && %s)", config.Name, groupAdd, cmd, chownCmd)

		// 将用户添加到 sudoers 列表
		if config.Sudoer {
			sudoLine := fmt.Sprintf("%s ALL=(ALL) NOPASSWD:ALL",
				config.Name)
			cmd = fmt.Sprintf("%s && %s",
				cmd,
				fmt.Sprintf("echo '%s' > /etc/sudoers.d/%s", sudoLine, config.Name))
		}

	case UserActionDel:
		cmd = fmt.Sprintf("%s -r %s", userdelCmd, config.Name)
		// prevent errors when user does not exist
		cmd = fmt.Sprintf("%s || [ $? -eq 6 ]", cmd)

		//	case UserActionModify:
		//		cmd = usermodCmd
	}

	return &UserModule{
		config: config,
		cmd:    cmd,
	}
}

// Execute 将命令传递给 executor 并返回其结果，executor 应该已经初始化了
func (mod *UserModule) Execute(exec executor.Executor) ([]byte, []byte, error) {
	a, b, err := exec.Execute(mod.cmd, true)
	if err != nil {
		switch mod.config.Action {
		case UserActionAdd:
			return a, b, ErrUserAddFailed.
				Wrap(err, "Failed to create new system user '%s' on remote host", mod.config.Name)
		case UserActionDel:
			return a, b, ErrUserDeleteFailed.
				Wrap(err, "Failed to delete system user '%s' on remote host", mod.config.Name)
		}
	}
	return a, b, nil
}

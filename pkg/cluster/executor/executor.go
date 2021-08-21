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
	"strings"
	"time"

	"github.com/joomcode/errorx"
)

var (
	// 来自某逻辑区域错误
	errNS = errorx.NewNamespace("executor")

	// SSH authorized_keys file
	defaultSSHAuthorizedKeys = "~/.ssh/authorized_keys"
)

// Executor 是 SSH executor 接口，所有任务都将通过 SSH executor 执行
type Executor interface {
	// 执行 run 命令，返回 stdout 和 stderr
	// 如果 cmd 超时不能退出，会返回 error，默认超时时间为 60 秒
	Execute(cmd string, sudo bool, timeout ...time.Duration) (stdout []byte, stderr []byte, err error)

	// 从或向目标传输副本文件
	Transfer(src, dst string, download bool, limit int) error
}

// 创建 SSH executor
func NewSSHExecutor(sudo bool, c SSHConfig) (Executor, error) {
	// set default values
	if c.Port <= 0 {
		c.Port = 22
	}

	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = time.Duration(DefaultConnectTimeout) * time.Second // 默认 SSH 连接超时时间
	}

	executor := &EasySSHExecutor{
		Locale: "C",
		Sudo:   sudo,
	}
	executor.initialize(c)

	return executor, nil
}

// FindSSHAuthorizedKeysFile 找到 SSH 授权密钥文件的正确路径
func FindSSHAuthorizedKeysFile(exec Executor) string {
	// 检测是否设置了授权密钥文件的自定义路径
	// NOTE: we do not yet support:
	//   - custom config for user (~/.ssh/config)
	//   - sshd started with custom config (other than /etc/ssh/sshd_config)
	//   - ssh server implementations other than OpenSSH (such as dropbear)
	sshAuthorizedKeys := defaultSSHAuthorizedKeys
	cmd := "grep -Ev '^\\s*#|^\\s*$' /etc/ssh/sshd_config"

	// 错误被忽略，因为有默认值
	stdout, _, _ := exec.Execute(cmd, true)
	for _, line := range strings.Split(string(stdout), "\n") {
		if !strings.Contains(line, "AuthorizedKeysFile") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			sshAuthorizedKeys = fields[1]
			break
		}
	}

	if !strings.HasPrefix(sshAuthorizedKeys, "/") && !strings.HasPrefix(sshAuthorizedKeys, "~") {
		sshAuthorizedKeys = fmt.Sprintf("~/%s", sshAuthorizedKeys)
	}
	return sshAuthorizedKeys
}

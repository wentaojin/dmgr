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
	"time"

	expect "github.com/google/goexpect"
	"github.com/wentaojin/dmgr/pkg/cluster/ctxt"
	"github.com/wentaojin/dmgr/pkg/dmgrutil"
)

var (
	sshKeyGenCMD = "ssh-keygen -t ed25519"
)

// SSHKeyCopy
// 用于生成 SSH 密钥
type SSHKeyGen struct {
	homeSshDir     string
	executeTimeout uint64
}

// Execute implements the Task interface
func (s *SSHKeyGen) Execute(ctx *ctxt.Context) error {
	ctx.Ev.PublishTaskProgress(s, "Generate SSH keys")

	// 存放于用户家目录，用于日常管理 SSH
	edHomePath := filepath.Join(s.homeSshDir, "id_ed25519")
	edHomePubPath := filepath.Join(s.homeSshDir, "id_ed25519.pub")

	// Skip ssh key generate
	if dmgrutil.IsExist(edHomePath) && dmgrutil.IsExist(edHomePubPath) {
		return nil
	}

	// 默认生成认证文件 HOME 家目录
	ctx.Ev.PublishTaskProgress(s, "Generate private and public key")
	e, _, err := expect.Spawn(sshKeyGenCMD, time.Second*time.Duration(s.executeTimeout))
	if err != nil {
		return err
	}
	defer e.Close()

	caser := []expect.Caser{
		&expect.BCase{R: "Enter", S: "\n"},
		&expect.BCase{R: "y/n", S: "y\n"},
		&expect.BCase{R: "fingerprint", S: "\n"},
	}

	for {
		output, _, _, err := e.ExpectSwitchCase(caser, time.Second*time.Duration(s.executeTimeout))

		if strings.Contains(output, "fingerprint") {
			break
		}
		if err != nil {
			e, _, _ = expect.Spawn(sshKeyGenCMD, time.Second*time.Duration(s.executeTimeout))
			continue
		}
	}
	return nil
}

// Rollback implements the Task interface
func (s *SSHKeyGen) Rollback(ctx *ctxt.Context) error {
	return ErrNoExecutor
}

// String implements the fmt.Stringer interface
func (s *SSHKeyGen) String() string {
	return fmt.Sprintf("SSHKeyGen: homePath=%s", s.homeSshDir)
}

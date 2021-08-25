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
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"

	"github.com/wentaojin/dmgr/pkg/cluster/executor"

	"github.com/wentaojin/dmgr/response"
	"github.com/xxjwxc/gowp/workpool"

	"github.com/wentaojin/dmgr/pkg/cluster/ctxt"

	"github.com/wentaojin/dmgr/pkg/dmgrutil"
)

// SSHKeyCopy
// 分发 SSH 密钥 ssh-copy-id
type SSHKeyCopy struct {
	homeSshDir     string
	clusterUser    string
	clusterSshDir  string
	hosts          []response.MachineRespStruct
	executeTimeout uint64
	workerThreads  int
}

// Execute implements the Task interface
func (s *SSHKeyCopy) Execute(ctx *ctxt.Context) error {
	// 存放于用户家目录，用于日常管理 SSH
	edHomePath := filepath.Join(s.homeSshDir, "id_ed25519")
	expandedHomePath, err := homedir.Expand(edHomePath)
	if err != nil {
		return err
	}

	edHomePubPath := filepath.Join(s.homeSshDir, "id_ed25519.pub")
	expandedHomePubPath, err := homedir.Expand(edHomePubPath)
	if err != nil {
		return err
	}

	// SSH 认证文件分发
	wp := workpool.New(s.workerThreads)
	for _, host := range s.hosts {
		server := host
		edFile := expandedHomePath
		timeout := s.executeTimeout
		wp.DoWait(func() error {
			isConnect, err := server.SshAuthTest(edFile)
			if err != nil {
				return err
			}
			if !isConnect {
				if err := server.SshCopyID(timeout); err != nil {
					return err
				}
			}
			return nil
		})
	}
	if err := wp.Wait(); err != nil {
		return err
	}
	if !wp.IsDone() {
		return fmt.Errorf("ssh key gen error")
	}

	// 存放于集群 SSH 目录，用于程序管理 SSH
	edSshPath := filepath.Join(s.clusterSshDir, "id_ed25519")
	edSshPubPath := filepath.Join(s.clusterSshDir, "id_ed25519.pub")

	// 本机 COPY 认证文件到集群管理目录
	currentUser, currentIP, err := dmgrutil.GetClientOutBoundIP()
	_, stdErr, err := executor.NewLocalExecutor(currentIP, currentUser, currentUser == "root").Execute(fmt.Sprintf("cp %v %v;cp %v %v", expandedHomePath, edSshPath, expandedHomePubPath, edSshPubPath), executor.DefaultExecuteTimeout)
	if err != nil || len(stdErr) != 0 {
		return fmt.Errorf("local copy err: [%v], stderr: [%v]", err, string(stdErr))
	}

	ctx.PrivateKeyPath = edSshPath
	ctx.PublicKeyPath = edSshPubPath
	return nil
}

// Rollback implements the Task interface
func (s *SSHKeyCopy) Rollback(ctx *ctxt.Context) error {
	return os.Remove(s.clusterSshDir)
}

// String implements the fmt.Stringer interface
func (s *SSHKeyCopy) String() string {
	return fmt.Sprintf("SSHKeyCopy: homePath=%s, clusterPath=%s", s.homeSshDir, s.clusterSshDir)
}

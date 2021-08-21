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
	"strings"

	"github.com/wentaojin/dmgr/pkg/cluster/ctxt"

	"github.com/pingcap/errors"
)

// Mkdir 用于在目标主机上创建目录
type Mkdir struct {
	user string
	host string
	dirs []string
}

// Execute implements the Task interface
func (m *Mkdir) Execute(ctx *ctxt.Context) error {
	exec, found := ctx.GetExecutor(m.host)
	if !found {
		return ErrNoExecutor
	}
	for _, dir := range m.dirs {
		if !strings.HasPrefix(dir, "/") {
			return fmt.Errorf("dir is a relative path: %s", dir)
		}
		if strings.Contains(dir, ",") {
			return fmt.Errorf("dir name contains invalid characters: %v", dir)
		}

		xs := strings.Split(dir, "/")

		// 递归创建目录
		// 目录 /a/b/c 将展平为
		// 		test -d /a || (mkdir /a && chown tidb:tidb /a)
		//		test -d /a/b || (mkdir /a/b && chown tidb:tidb /a/b)
		//		test -d /a/b/c || (mkdir /a/b/c && chown tidb:tidb /a/b/c)
		for i := 0; i < len(xs); i++ {
			if xs[i] == "" {
				continue
			}
			cmd := fmt.Sprintf(
				`test -d %[1]s || (mkdir -p %[1]s && chown %[2]s:$(id -g -n %[2]s) %[1]s)`,
				strings.Join(xs[:i+1], "/"),
				m.user,
			)
			_, _, err := exec.Execute(cmd, true) // use root to create the dir
			if err != nil {
				return errors.Trace(err)
			}
		}
	}

	return nil
}

// Rollback implements the Task interface
func (m *Mkdir) Rollback(ctx *ctxt.Context) error {
	return ErrUnsupportedRollback
}

// String implements the fmt.Stringer interface
func (m *Mkdir) String() string {
	return fmt.Sprintf("Mkdir: host=%s, directories='%s'", m.host, strings.Join(m.dirs, "','"))
}

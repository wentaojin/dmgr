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

	"github.com/wentaojin/dmgr/pkg/dmgrutil"

	"github.com/pingcap/errors"
)

// CopyComponent 用于复制一个组件特定版本相关的所有文件到 path 的目标目录
type CopyComponent struct {
	componentName  string
	clusterVersion string
	host           string
	srcPath        string
	dstPath        string
}

// Execute implements the Task interface
func (c *CopyComponent) Execute(ctx *ctxt.Context) error {
	exec, found := ctx.GetExecutor(c.host)
	if !found {
		return ErrNoExecutor
	}

	err := exec.Transfer(c.srcPath, c.dstPath, false, 0)
	if err != nil {
		return errors.Annotatef(err, "failed to scp %s to %s:%s", c.srcPath, c.host, c.dstPath)
	}

	if strings.ToLower(c.componentName) == dmgrutil.ComponentGrafana {
		baseDir := strings.Split(c.dstPath, "/")
		cmd := fmt.Sprintf(`tar --no-same-owner -zxf %s -C %s && rm %s`,
			c.dstPath,
			baseDir[len(baseDir)-1],
			c.dstPath)

		_, stderr, err := exec.Execute(cmd, false)
		if err != nil || len(stderr) != 0 {
			return errors.Annotatef(err, "stderr: %s", string(stderr))
		}
	}
	return nil
}

// Rollback implements the Task interface
func (c *CopyComponent) Rollback(ctx *ctxt.Context) error {
	return ErrUnsupportedRollback
}

// String implements the fmt.Stringer interface
func (c *CopyComponent) String() string {
	return fmt.Sprintf("CopyComponent: component=%s, version=%s, remote=%s:%s",
		c.componentName, c.clusterVersion, c.host, c.dstPath)
}

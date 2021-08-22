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

	"github.com/wentaojin/dmgr/pkg/dmgrutil"

	"github.com/wentaojin/dmgr/pkg/cluster/ctxt"

	"github.com/pingcap/errors"
)

// CopyFile will copy a local file to the target host
type CopyFile struct {
	src        string
	dst        string
	fileType   string
	remoteHost string
	download   bool
	limit      int
}

// Execute implements the Task interface
func (c *CopyFile) Execute(ctx *ctxt.Context) error {
	e, ok := ctx.GetExecutor(c.remoteHost)
	if !ok {
		return ErrNoExecutor
	}

	if c.fileType == dmgrutil.FileTypeSystemd {
		err := e.Transfer(c.src, c.dst, c.download, c.limit)
		if err != nil {
			return errors.Annotate(err, "failed to transfer file")
		}

		cmd := fmt.Sprintf(`cp %s %s && rm %s`,
			c.dst,
			dmgrutil.AbsClusterSystemdDir(),
			c.dst)
		_, stderr, err := e.Execute(cmd, true)
		if err != nil || len(stderr) != 0 {
			return errors.Annotatef(err, "stderr: %s", string(stderr))
		}
		return nil
	}
	err := e.Transfer(c.src, c.dst, c.download, c.limit)
	if err != nil {
		return errors.Annotate(err, "failed to transfer file")
	}
	return nil
}

// Rollback implements the Task interface
func (c *CopyFile) Rollback(ctx *ctxt.Context) error {
	return ErrUnsupportedRollback
}

// String implements the fmt.Stringer interface
func (c *CopyFile) String() string {
	if c.download {
		return fmt.Sprintf("CopyFile: remote=%s:%s, local=%s", c.remoteHost, c.src, c.dst)
	}
	return fmt.Sprintf("CopyFile: local=%s, remote=%s:%s", c.src, c.remoteHost, c.dst)
}

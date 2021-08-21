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
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/joomcode/errorx"
	"github.com/wentaojin/dmgr/pkg/dmgrutil"

	"go.uber.org/zap"
)

// Local execute the command at local host.
type Local struct {
	Host   string // local IP
	User   string // username
	Sudo   bool   // all commands run with this executor will be using sudo
	Locale string // the locale used when executing the command
}

func NewLocalExecutor(host string, user string, sudo bool) *Local {
	return &Local{
		Host:   host,
		User:   user,
		Sudo:   sudo,
		Locale: "",
	}
}

// Execute implements Executor interface.
func (l *Local) Execute(cmd string, execTimeout ...time.Duration) ([]byte, []byte, error) {
	// try to acquire root permission
	if l.Sudo {
		cmd = fmt.Sprintf("sudo -H -u root bash -c \"cd; %s\"", cmd)
	} else {
		cmd = fmt.Sprintf("sudo -H -u %s bash -c \"cd; %s\"", l.User, cmd)
	}

	// set a basic PATH in case it's empty on login
	cmd = fmt.Sprintf("PATH=$PATH:/usr/bin:/usr/sbin %s", cmd)

	if l.Locale != "" {
		cmd = fmt.Sprintf("export LANG=%s; %s", l.Locale, cmd)
	}

	// run command on remote host
	// default timeout is 60s in easyssh-proxy
	if len(execTimeout) == 0 {
		execTimeout = append(execTimeout, time.Duration(DefaultExecuteTimeout)*time.Second)
	}

	ctx := context.Background()
	if len(execTimeout) > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), execTimeout[0])
		defer cancel()
	}

	command := exec.CommandContext(ctx, "/bin/sh", "-c", cmd)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	command.Stdout = stdout
	command.Stderr = stderr

	err := command.Run()

	dmgrutil.Logger.Info("LocalCommand",
		zap.String("cmd", cmd),
		zap.Error(err),
		zap.String("stdout", stdout.String()),
		zap.String("stderr", stderr.String()))

	if err != nil {
		sshErr := ErrSSHExecuteFailed.
			Wrap(err, "Failed to execute command locally").
			WithProperty(ErrPropSSHCommand, cmd).
			WithProperty(ErrPropSSHStdout, stdout).
			WithProperty(ErrPropSSHStderr, stderr)
		if len(stdout.Bytes()) > 0 || len(stderr.Bytes()) > 0 {
			output := strings.TrimSpace(strings.Join([]string{stdout.String(), stderr.String()}, "\n"))
			sshErr = sshErr.
				WithProperty(
					errorx.RegisterPrintableProperty(
						fmt.Sprintf("Command output on remote host %s", l.Host)),
					output)
		}
		return stdout.Bytes(), stderr.Bytes(), sshErr
	}

	return stdout.Bytes(), stderr.Bytes(), err
}

// Transfer implements Executer interface.
func (l *Local) Transfer(src string, dst string, download bool, limit int) error {
	targetPath := filepath.Dir(dst)
	if err := dmgrutil.CreateDir(targetPath); err != nil {
		return err
	}

	cmd := ""
	user, err := user.Current()
	if err != nil {
		return err
	}
	if download || user.Username == l.User {
		cmd = fmt.Sprintf("cp %s %s", src, dst)
	} else {
		cmd = fmt.Sprintf("sudo -H -u root bash -c \"cp %[1]s %[2]s && chown %[3]s:$(id -g -n %[3]s) %[2]s\"", src, dst, l.User)
	}

	command := exec.Command("/bin/sh", "-c", cmd)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	command.Stdout = stdout
	command.Stderr = stderr

	err = command.Run()

	dmgrutil.Logger.Info("CPCommand",
		zap.String("cmd", cmd),
		zap.Error(err),
		zap.String("stdout", stdout.String()),
		zap.String("stderr", stderr.String()))

	if err != nil {
		sshErr := ErrSSHExecuteFailed.
			Wrap(err, "Failed to transfer file over local cp").
			WithProperty(ErrPropSSHCommand, cmd).
			WithProperty(ErrPropSSHStdout, stdout).
			WithProperty(ErrPropSSHStderr, stderr)
		if len(stdout.Bytes()) > 0 || len(stderr.Bytes()) > 0 {
			output := strings.TrimSpace(strings.Join([]string{stdout.String(), stderr.String()}, "\n"))
			sshErr = sshErr.
				WithProperty(
					errorx.RegisterPrintableProperty(
						fmt.Sprintf("Command output on remote host %s", l.Host)),
					output)
		}
		return sshErr
	}

	return err
}

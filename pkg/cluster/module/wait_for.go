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
	"bytes"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/wentaojin/dmgr/pkg/dmgrutil"

	"github.com/wentaojin/dmgr/pkg/cluster/executor"
)

// WaitForConfig is the configurations of WaitFor module.
type WaitForConfig struct {
	Port  int           // Port number to poll.
	Sleep time.Duration // Duration to sleep between checks, default 1 second.
	// Choices:
	// started
	// stopped
	// When checking a port started will ensure the port is open, stopped will check that it is closed
	State   string
	Timeout time.Duration // Maximum duration to wait for.
}

// WaitFor is the module used to wait for some condition.
type WaitFor struct {
	c WaitForConfig
}

// NewWaitFor create a WaitFor instance.
func NewWaitFor(c WaitForConfig) *WaitFor {
	if c.Sleep == 0 {
		c.Sleep = time.Duration(DefaultSystemdSleepTime) * time.Second
	}
	if c.Timeout == 0 {
		c.Timeout = time.Duration(DefaultSystemdExecuteTimeout) * time.Second
	}
	if c.State == "" {
		c.State = "started"
	}

	w := &WaitFor{
		c: c,
	}

	return w
}

// Execute the module return nil if successfully wait for the event.
func (w *WaitFor) Execute(e executor.Executor) (err error) {
	pattern := []byte(fmt.Sprintf(":%d ", w.c.Port))

	retryOpt := dmgrutil.RetryOption{
		Delay:   w.c.Sleep,
		Timeout: w.c.Timeout,
	}
	if err := dmgrutil.Retry(func() error {
		// only listing TCP ports
		stdout, _, err := e.Execute("ss -ltn", false)
		if err == nil {
			switch w.c.State {
			case "started":
				if bytes.Contains(stdout, pattern) {
					return nil
				}
			case "stopped":
				if !bytes.Contains(stdout, pattern) {
					return nil
				}
			}
			return errors.New("still waiting for port state to be satisfied")
		}
		return err
	}, retryOpt); err != nil {
		dmgrutil.Logger.Debug("retry error: %s", zap.Error(err))
		return fmt.Errorf("timed out waiting for port %d to be %s after %s", w.c.Port, w.c.State, w.c.Timeout)
	}
	return nil
}

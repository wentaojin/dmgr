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
package systemd

import (
	"fmt"
	"html/template"
	"os"
)

// Config represent the data to generate systemd config
type Config struct {
	ServiceName        string
	User               string
	DeployDir          string
	DisableSendSigkill bool
	// Takes one of no, on-success, on-failure, on-abnormal, on-watchdog, on-abort, or always.
	// The Template set as always if this is not setted.
	Restart string
}

// NewSystemdConfig returns a Config with given arguments
func NewSystemdConfig(service, user, deployDir string) *Config {
	return &Config{
		ServiceName: service,
		User:        user,
		DeployDir:   deployDir,
	}
}

func (c *Config) WithDisableSendSigkill(action bool) *Config {
	c.DisableSendSigkill = action
	return c
}

// WithRestart Takes one of no, on-success, on-failure, on-abnormal, on-watchdog, on-abort, or always.
//	// The Template set as always if this is not setted.
func (c *Config) WithRestart(action string) *Config {
	c.Restart = action
	return c
}

// ConfigToFile write config content to specific path
func (c *Config) ConfigToFile(tmplFile, outFile string) error {
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		return fmt.Errorf("SystemdConfig template new failed: %v", err)
	}

	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("SystemdConfig open file failed: %v", err)
	}
	if err := tmpl.Execute(f, c); err != nil {
		return fmt.Errorf("SystemdConfig template execute failed: %v", err)
	}
	return nil
}

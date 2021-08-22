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
package script

import (
	"fmt"
	"html/template"
	"os"
)

// DMWorkerScript represent the data to generate TiDB config
type DMWorkerScript struct {
	Name      string
	IP        string
	Port      uint64
	DeployDir string
	LogDir    string
	Endpoints []*DMMasterScript
}

// NewDMWorkerScript returns a DMWorkerScript with given arguments
func NewDMWorkerScript(name, ip, deployDir, logDir string) *DMWorkerScript {
	return &DMWorkerScript{
		Name:      name,
		IP:        ip,
		Port:      8262,
		DeployDir: deployDir,
		LogDir:    logDir,
	}
}

// WithPort set Port field of DMWorkerScript
func (c *DMWorkerScript) WithPort(port uint64) *DMWorkerScript {
	c.Port = port
	return c
}

// AppendEndpoints add new PDScript to Endpoints field
func (c *DMWorkerScript) AppendEndpoints(ends ...*DMMasterScript) *DMWorkerScript {
	c.Endpoints = append(c.Endpoints, ends...)
	return c
}

// ConfigToFile write config content to specific path
func (c *DMWorkerScript) ConfigToFile(tmplFile, outFile string) error {
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		return fmt.Errorf("DMWorkerScript template new failed: %v", err)
	}

	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("DMWorkerScript open file failed: %v", err)
	}
	if err := tmpl.Execute(f, c); err != nil {
		return fmt.Errorf("DMWorkerScript template execute failed: %v", err)
	}
	return nil
}

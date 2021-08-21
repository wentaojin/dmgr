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

// DMMasterScript 表示生成 dm-master config 的数据
type DMMasterScript struct {
	Name      string
	Scheme    string
	IP        string
	Port      uint64
	PeerPort  uint64
	DeployDir string
	DataDir   string
	LogDir    string
	Endpoints []*DMMasterScript
}

// NewDMMasterScript 返回带有给定参数的 DMMasterScript
func NewDMMasterScript(name, ip, deployDir, dataDir, logDir string) *DMMasterScript {
	return &DMMasterScript{
		Name:      name,
		Scheme:    "http",
		IP:        ip,
		Port:      8261,
		PeerPort:  8291,
		DeployDir: deployDir,
		DataDir:   dataDir,
		LogDir:    logDir,
	}
}

// WithScheme set Scheme field of NewDMMasterScript
func (c *DMMasterScript) WithScheme(scheme string) *DMMasterScript {
	c.Scheme = scheme
	return c
}

// WithPort set Port field of DMMasterScript
func (c *DMMasterScript) WithPort(port uint64) *DMMasterScript {
	c.Port = port
	return c
}

// WithPeerPort set PeerPort field of DMMasterScript
func (c *DMMasterScript) WithPeerPort(port uint64) *DMMasterScript {
	c.PeerPort = port
	return c
}

// AppendEndpoints add new DMMasterScript to Endpoints field
func (c *DMMasterScript) AppendEndpoints(ends ...*DMMasterScript) *DMMasterScript {
	c.Endpoints = append(c.Endpoints, ends...)
	return c
}

// ConfigToFile write config content to specific path
func (c *DMMasterScript) ConfigToFile(tmplFile, outFile string) error {
	tmpl, err := template.New("dm-master").Parse(tmplFile)
	if err != nil {
		return fmt.Errorf("DMMasterScript template new failed: %v", err)
	}

	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("DMMasterScript open file failed: %v", err)
	}
	if err := tmpl.Execute(f, c); err != nil {
		return fmt.Errorf("DMMasterScript template execute failed: %v", err)
	}
	return nil
}

// DMMasterScaleScript 表示在缩放时生成 dm-master 配置的数据
type DMMasterScaleScript struct {
	DMMasterScript
}

// NewDMMasterScaleScript return a new DMMasterScaleScript
func NewDMMasterScaleScript(name, ip, deployDir, dataDir, logDir string) *DMMasterScaleScript {
	return &DMMasterScaleScript{*NewDMMasterScript(name, ip, deployDir, dataDir, logDir)}
}

// WithScheme set Scheme field of DMMasterScaleScript
func (c *DMMasterScaleScript) WithScheme(scheme string) *DMMasterScaleScript {
	c.Scheme = scheme
	return c
}

// WithPort set Port field of DMMasterScript
func (c *DMMasterScaleScript) WithPort(port uint64) *DMMasterScaleScript {
	c.Port = port
	return c
}

// WithPeerPort set PeerPort field of DMMasterScript
func (c *DMMasterScaleScript) WithPeerPort(port uint64) *DMMasterScaleScript {
	c.PeerPort = port
	return c
}

// AppendEndpoints add new DMMasterScript to Endpoints field
func (c *DMMasterScaleScript) AppendEndpoints(ends ...*DMMasterScript) *DMMasterScaleScript {
	c.Endpoints = append(c.Endpoints, ends...)
	return c
}

// ConfigToFile write config content to specific path
func (c *DMMasterScaleScript) ConfigToFile(tmplFile, outFile string) error {
	tmpl, err := template.New("dm-master").Parse(tmplFile)
	if err != nil {
		return fmt.Errorf("DMMasterScript template new failed: %v", err)
	}

	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("DMMasterScript open file failed: %v", err)
	}
	if err := tmpl.Execute(f, c); err != nil {
		return fmt.Errorf("DMMasterScript template execute failed: %v", err)
	}
	return nil
}

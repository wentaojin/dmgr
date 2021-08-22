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

// AlertManagerScript represent the data to generate AlertManager start script
type AlertManagerScript struct {
	IP          string
	WebPort     uint64
	ClusterPort uint64
	DeployDir   string
	DataDir     string
	LogDir      string
	TLSEnabled  bool
	EndPoints   []*AlertManagerScript
}

// NewAlertManagerScript returns a AlertManagerScript with given arguments
func NewAlertManagerScript(ip, deployDir, dataDir, logDir string, enableTLS bool) *AlertManagerScript {
	return &AlertManagerScript{
		IP:          ip,
		WebPort:     9093,
		ClusterPort: 9094,
		DeployDir:   deployDir,
		DataDir:     dataDir,
		LogDir:      logDir,
		TLSEnabled:  enableTLS,
	}
}

// WithWebPort set WebPort field of AlertManagerScript
func (c *AlertManagerScript) WithWebPort(port uint64) *AlertManagerScript {
	c.WebPort = port
	return c
}

// WithClusterPort set WebPort field of AlertManagerScript
func (c *AlertManagerScript) WithClusterPort(port uint64) *AlertManagerScript {
	c.ClusterPort = port
	return c
}

// AppendEndpoints add new alert manager to Endpoints field
func (c *AlertManagerScript) AppendEndpoints(ends []*AlertManagerScript) *AlertManagerScript {
	c.EndPoints = append(c.EndPoints, ends...)
	return c
}

// ConfigToFile write config content to specific path
func (c *AlertManagerScript) ConfigToFile(tmplFile, outFile string) error {
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		return fmt.Errorf("AlertManagerScript template new failed: %v", err)
	}

	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("AlertManagerScript open file failed: %v", err)
	}
	if err := tmpl.Execute(f, c); err != nil {
		return fmt.Errorf("AlertManagerScript template execute failed: %v", err)
	}
	return nil
}

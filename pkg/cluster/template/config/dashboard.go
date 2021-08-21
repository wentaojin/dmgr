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
package config

import (
	"fmt"
	"os"
	"text/template"
)

// DashboardConfig represent the data to generate Dashboard config
type DashboardConfig struct {
	ClusterName string
	DeployDir   string
}

// NewDashboardConfig returns a DashboardConfig
func NewDashboardConfig(cluster, deployDir string) *DashboardConfig {
	return &DashboardConfig{
		ClusterName: cluster,
		DeployDir:   deployDir,
	}
}

// ConfigToFile write config content to specific path
func (c *DashboardConfig) ConfigToFile(tmplFile, outFile string) error {
	tmpl, err := template.New("dashboard").Parse(tmplFile)
	if err != nil {
		return fmt.Errorf("DashboardConfig template new failed: %v", err)
	}

	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("DashboardConfig open file failed: %v", err)
	}
	if err := tmpl.Execute(f, c); err != nil {
		return fmt.Errorf("DashboardConfig template execute failed: %v", err)
	}
	return nil
}

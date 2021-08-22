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

// GrafanaScript represent the data to generate Grafana config
type GrafanaScript struct {
	ClusterName string
	DeployDir   string
}

// NewGrafanaScript returns a GrafanaScript with given arguments
func NewGrafanaScript(cluster, deployDir string) *GrafanaScript {
	return &GrafanaScript{
		ClusterName: cluster,
		DeployDir:   deployDir,
	}
}

// ConfigToFile write config content to specific path
func (c *GrafanaScript) ConfigToFile(tmplFile, outFile string) error {
	tmpl, err := template.New("grafana").ParseFiles(tmplFile)
	if err != nil {
		return fmt.Errorf("GrafanaScript template new failed: %v", err)
	}

	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("GrafanaScript open file failed: %v", err)
	}
	if err := tmpl.Execute(f, c); err != nil {
		return fmt.Errorf("GrafanaScript template execute failed: %v", err)
	}
	return nil
}

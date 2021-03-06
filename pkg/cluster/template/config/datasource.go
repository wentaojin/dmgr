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

// DatasourceConfig represent the data to generate Datasource config
type DatasourceConfig struct {
	ClusterName string
	IP          string
	Port        string
}

// NewDatasourceConfig returns a DatasourceConfig
func NewDatasourceConfig(cluster, ip string) *DatasourceConfig {
	return &DatasourceConfig{
		ClusterName: cluster,
		IP:          ip,
		Port:        "9090",
	}
}

// WithPort set Port field of DatasourceConfig
func (c *DatasourceConfig) WithPort(port string) *DatasourceConfig {
	c.Port = port
	return c
}

// ConfigToFile write config content to specific path
func (c *DatasourceConfig) ConfigToFile(tmplFile, outFile string) error {
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		return fmt.Errorf("DatasourceConfig template new failed: %v", err)
	}

	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("DatasourceConfig open file failed: %v", err)
	}
	if err := tmpl.Execute(f, c); err != nil {
		return fmt.Errorf("DatasourceConfig template execute failed: %v", err)
	}
	return nil
}

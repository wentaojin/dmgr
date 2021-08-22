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
	"regexp"
)

// PrometheusScript represent the data to generate Prometheus config
type PrometheusScript struct {
	IP        string
	Port      uint64
	DeployDir string
	DataDir   string
	LogDir    string
	Retention string
}

// NewPrometheusScript returns a PrometheusScript with given arguments
func NewPrometheusScript(ip, deployDir, dataDir, logDir string) *PrometheusScript {
	return &PrometheusScript{
		IP:        ip,
		Port:      9090,
		DeployDir: deployDir,
		DataDir:   dataDir,
		LogDir:    logDir,
	}
}

// WithPort set Port field of PrometheusScript
func (c *PrometheusScript) WithPort(port uint64) *PrometheusScript {
	c.Port = port
	return c
}

// WithRetention set Retention field of PrometheusScript
func (c *PrometheusScript) WithRetention(retention string) *PrometheusScript {
	valid, _ := regexp.MatchString("^[1-9]\\d*d$", retention)
	if retention == "" || !valid {
		c.Retention = "30d"
	} else {
		c.Retention = retention
	}
	return c
}

// ConfigToFile write config content to specific path
func (c *PrometheusScript) ConfigToFile(tmplFile, outFile string) error {
	tmpl, err := template.New("prometheus").ParseFiles(tmplFile)
	if err != nil {
		return fmt.Errorf("PrometheusScript template new failed: %v", err)
	}

	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("PrometheusScript open file failed: %v", err)
	}
	if err := tmpl.Execute(f, c); err != nil {
		return fmt.Errorf("PrometheusScript template execute failed: %v", err)
	}
	return nil
}

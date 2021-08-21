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

// PrometheusConfig represent the data to generate Prometheus config
type PrometheusConfig struct {
	ClusterName string
	TLSEnabled  bool

	AlertmanagerAddrs []string

	GrafanaAddr string

	DMMasterAddrs []string
	DMWorkerAddrs []string

	RemoteConfig string
}

// NewPrometheusConfig returns a PrometheusConfig
func NewPrometheusConfig(clusterName, clusterVersion string, enableTLS bool) *PrometheusConfig {
	cfg := &PrometheusConfig{
		ClusterName: clusterName,
		TLSEnabled:  enableTLS,
	}
	return cfg
}

// AddAlertmanager add an alertmanager address
func (c *PrometheusConfig) AddAlertmanager(alertmanagerAddrs []string) *PrometheusConfig {
	c.AlertmanagerAddrs = alertmanagerAddrs
	return c
}

// AddGrafana add an kafka exporter address
func (c *PrometheusConfig) AddGrafana(grafanaAddr string) *PrometheusConfig {
	c.GrafanaAddr = grafanaAddr
	return c
}

// AddDMMaster add an dm-master address
func (c *PrometheusConfig) AddDMMaster(dmMasters []string) *PrometheusConfig {
	c.DMMasterAddrs = dmMasters
	return c
}

// AddDMWorker add an dm-worker address
func (c *PrometheusConfig) AddDMWorker(dmWorkers []string) *PrometheusConfig {
	c.DMWorkerAddrs = dmWorkers
	return c
}

// SetRemoteConfig set remote read/write config
func (c *PrometheusConfig) SetRemoteConfig(cfg string) *PrometheusConfig {
	c.RemoteConfig = cfg
	return c
}

// ConfigToFile write config content to specific path
func (c *PrometheusConfig) ConfigToFile(tmplFile, outFile string) error {
	tmpl, err := template.New("prometheus").Parse(tmplFile)
	if err != nil {
		return fmt.Errorf("PrometheusConfig template new failed: %v", err)
	}

	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("PrometheusConfig open file failed: %v", err)
	}
	if err := tmpl.Execute(f, c); err != nil {
		return fmt.Errorf("PrometheusConfig template execute failed: %v", err)
	}
	return nil
}

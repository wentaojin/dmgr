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

// GrafanaConfig represent the data to generate Grafana config
type GrafanaConfig struct {
	DeployDir       string
	IP              string
	Port            uint64
	Username        string // admin_user
	Password        string // admin_password
	AnonymousEnable bool   // anonymous enable
	RootURL         string // root_url
	Domain          string // domain
}

// NewGrafanaConfig returns a GrafanaConfig
func NewGrafanaConfig(ip, deployDir string) *GrafanaConfig {
	return &GrafanaConfig{
		DeployDir: deployDir,
		IP:        ip,
		Port:      3000,
	}
}

// WithPort set Port field of GrafanaConfig
func (c *GrafanaConfig) WithPort(port uint64) *GrafanaConfig {
	c.Port = port
	return c
}

// WithUsername sets username of admin user
func (c *GrafanaConfig) WithUsername(user string) *GrafanaConfig {
	c.Username = user
	return c
}

// WithPassword sets password of admin user
func (c *GrafanaConfig) WithPassword(passwd string) *GrafanaConfig {
	c.Password = passwd
	return c
}

// WithAnonymousenable sets anonymousEnable of anonymousEnable
func (c *GrafanaConfig) WithAnonymousenable(anonymousEnable bool) *GrafanaConfig {
	c.AnonymousEnable = anonymousEnable
	return c
}

// WithRootURL sets rootURL of root url
func (c *GrafanaConfig) WithRootURL(rootURL string) *GrafanaConfig {
	c.RootURL = rootURL
	return c
}

// WithDomain sets domain of server domain
func (c *GrafanaConfig) WithDomain(domain string) *GrafanaConfig {
	c.Domain = domain
	return c
}

// ConfigToFile write config content to specific path
func (c *GrafanaConfig) ConfigToFile(tmplFile, outFile string) error {
	tmpl, err := template.New("grafana").ParseFiles(tmplFile)
	if err != nil {
		return fmt.Errorf("GrafanaConfig template new failed: %v", err)
	}

	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("GrafanaConfig open file failed: %v", err)
	}
	if err := tmpl.Execute(f, c); err != nil {
		return fmt.Errorf("GrafanaConfig template execute failed: %v", err)
	}
	return nil
}

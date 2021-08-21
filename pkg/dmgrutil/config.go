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
package dmgrutil

import (
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ListenPort       string           `toml:"listen-port" json:"listen-port"`
	DbConfig         DbConfig         `toml:"db" json:"db"`
	LogConfig        LogConfig        `toml:"log" json:"log"`
	MiddlewareConfig MiddlewareConfig `toml:"middleware"`
}

type DbConfig struct {
	Host     string `toml:"host" json:"host"`
	Port     int    `toml:"port" json:"port"`
	User     string `toml:"user" json:"user"`
	Password string `toml:"password" json:"password"`
	DBName   string `toml:"db-name" json:"db-name"`
}

type LogConfig struct {
	LogLevel   string `toml:"log-level" json:"log-level"`
	LogFile    string `toml:"log-file" json:"log-file"`
	MaxSize    int    `toml:"max-size" json:"max-size"`
	MaxAge     int    `toml:"max-age" json:"max-age"`
	MaxBackups int    `toml:"max-backups" json:"max-backups"`
}

type MiddlewareConfig struct {
	MaxRateLimiter int64  `toml:"max-rate-limiter" json:"max-rate-limiter"`
	JwtRealm       string `toml:"jwt-realm" json:"jwt-realm"`
	JwtKey         string `toml:"jwt-key" json:"jwt-key"`
	JwtTimeout     int    `toml:"jwt-timeout" json:"jwt-timeout"`
	JwtMaxRefresh  int    `toml:"jwt-max-refresh" json:"jwt-max-refresh"`
}

// 配置文件读取
func ReadConfigFile(file string) (*Config, error) {
	cfg := &Config{}
	if err := cfg.configFromFile(file); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (c *Config) configFromFile(file string) error {
	if _, err := toml.DecodeFile(file, c); err != nil {
		return fmt.Errorf("failed decode toml config file %s: %v", file, err)
	}
	return nil
}

// 配置文件内容序列化
func (c *Config) String() string {
	cfg, err := json.Marshal(c)
	if err != nil {
		return "<nil>"
	}
	return string(cfg)
}

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
package request

// 任务上游数据源
type TaskSourceStruct struct {
	SourceName string `json:"source_name" form:"source_name" db:"source_name" binding:"required"`
}

// 任务 Source 配置信息
type TaskSourceConfStruct struct {
	EnableGtid             string `json:"enable_gtid" form:"enable_gtid" db:"enable_gtid" binding:"required"`
	RelayBinlogGtid        string `json:"relay_binlog_gtid" form:"relay_binlog_gtid" db:"relay_binlog_gtid"`
	EnableRelay            string `json:"enable_relay" form:"enable_relay" db:"enable_relay" binding:"required"`
	RelayBinlogName        string `json:"relay_binlog_name" form:"relay_binlog_name" db:"relay_binlog_name"`
	RelayDir               string `json:"relay_dir" form:"relay_dir" db:"relay_dir"`
	PurgeInterval          int    `json:"purge_interval" form:"purge_interval" db:"purge_interval"`
	PurgeExpires           int    `json:"purge_expires" form:"purge_expires" db:"purge_expires"`
	PurgeRemainSpace       int    `json:"purge_remain_space" form:"purge_remain_space" db:"purge_remain_space"`
	CheckerCheckEnable     string `json:"checker_check_enable" form:"checker_check_enable" db:"checker_check_enable"`
	CheckerBackoffRollback string `json:"checker_backoff_rollback" form:"checker_backoff_rollback" db:"checker_backoff_rollback"`
	CheckerBackoffMax      string `json:"checker_backoff_max" form:"checker_backoff_max" db:"checker_backoff_max"`
}

// 任务下游数据源
type TaskTargetStruct struct {
	TargetName       string `json:"target_name" form:"target_name" db:"target_name" binding:"required"`
	MaxAllowedPacket int    `json:"max_allowed_packet" form:"max_allowed_packet" db:"max_allowed_packet"`
	Session          string `json:"session" form:"session" db:"session"`
}

// 任务数据源请求
type TaskDatasourceStruct struct {
	Host     string `json:"host" form:"host" db:"host" binding:"required"`
	User     string `json:"user" form:"user" db:"user" binding:"required"`
	Password string `json:"password" form:"password" db:"password" binding:"required"`
	Port     int    `json:"port" form:"port" db:"port" binding:"required"`
}

// 任务数据源 SSL 认证请求
type TaskDatasourceSslStruct struct {
	SslCA   string `json:"ssl_ca" form:"ssl_ca" db:"ssl_ca"`
	SslCert string `json:"ssl_cert" form:"ssl_cert" db:"ssl_cert"`
	SslKey  string `json:"ssl_key" form:"ssl_key" db:"ssl_key"`
}

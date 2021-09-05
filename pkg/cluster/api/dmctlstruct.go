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
package api

// 任务请求 JSON 接口请求 PATH
const (
	RelayStatusPath = "relay_status.stage"
	WorkerNamePath  = "worker_name"
)

// 任务请求 JSON 字符变量
const (
	RelayStatusRunningStage = "Running"
)

// 任务 DM Worker 请求
type WorkerNameBodyStruct struct {
	WorkerName string `json:"worker_name"`
}

func NewWorkerNameBody(workerName string) *WorkerNameBodyStruct {
	return &WorkerNameBodyStruct{
		WorkerName: workerName,
	}
}

// 任务 Source 请求
type RelayStatusBodyStruct struct {
	WorkerName      string `json:"worker_name"`
	RelayBinlogName string `json:"relay_binlog_name"`
	RelayBinlogGtid string `json:"relay_binlog_gtid"`
	RelayDir        string `json:"relay_dir"`
	Purge           Purge  `json:"purge"`
}

type Purge struct {
	Interval    int `json:"interval"`
	Expires     int `json:"expires"`
	RemainSpace int `json:"remain_space"`
}

// todo: 待完善
func NewRelayStatusBody(respBody []byte) *RelayStatusBodyStruct {
	return &RelayStatusBodyStruct{
		WorkerName:      "",
		RelayBinlogName: "",
		RelayBinlogGtid: "",
		RelayDir:        "",
		Purge:           Purge{},
	}
}

// 任务请求
type TaskBodyStruct struct {
	RemoveMeta bool `json:"remove_meta"`
	Task       Task `json:"task"`
}

type Security struct {
	SslCaContent   string `json:"ssl_ca_content"`
	SslCertContent string `json:"ssl_cert_content"`
	SslKeyContent  string `json:"ssl_key_content"`
}

type TargetConfig struct {
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	User     string   `json:"user"`
	Password string   `json:"password"`
	Security Security `json:"security"`
}

type EventFilterRule struct {
	RuleName    string   `json:"rule_name"`
	IgnoreEvent []string `json:"ignore_event"`
	IgnoreSQL   []string `json:"ignore_sql"`
}

type Source struct {
	SourceName string `json:"source_name"`
	Schema     string `json:"schema"`
	Table      string `json:"table"`
}

type Target struct {
	Schema string `json:"schema"`
	Table  string `json:"table"`
}

type TableMigrateRule struct {
	Source          Source   `json:"source"`
	Target          Target   `json:"target"`
	EventFilterName []string `json:"event_filter_name"`
}

type FullMigrateConf struct {
	ExportThreads int    `json:"export_threads"`
	ImportThreads int    `json:"import_threads"`
	DataDir       string `json:"data_dir"`
}

type IncrMigrateConf struct {
	ReplThreads int `json:"repl_threads"`
	ReplBatch   int `json:"repl_batch"`
}

type SourceConf struct {
	SourceName string `json:"source_name"`
	BinlogName string `json:"binlog_name"`
	BinlogPos  int    `json:"binlog_pos"`
	BinlogGtid string `json:"binlog_gtid"`
}

type SourceConfig struct {
	FullMigrateConf FullMigrateConf `json:"full_migrate_conf"`
	IncrMigrateConf IncrMigrateConf `json:"incr_migrate_conf"`
	SourceConf      []SourceConf    `json:"source_conf"`
}

type Task struct {
	Name                      string             `json:"name"`
	TaskMode                  string             `json:"task_mode"`
	ShardMode                 string             `json:"shard_mode"`
	MetaSchema                string             `json:"meta_schema"`
	EnhanceOnlineSchemaChange bool               `json:"enhance_online_schema_change"`
	OnDuplication             string             `json:"on_duplication"`
	TargetConfig              TargetConfig       `json:"target_config"`
	EventFilterRule           []EventFilterRule  `json:"event_filter_rule"`
	TableMigrateRule          []TableMigrateRule `json:"table_migrate_rule"`
	SourceConfig              SourceConfig       `json:"source_config"`
}

func NewTaskBody(respByte []byte) *TaskBodyStruct {
	return &TaskBodyStruct{
		RemoveMeta: false,
		Task:       Task{},
	}
}

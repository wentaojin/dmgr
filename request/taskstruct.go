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

// 任务上游数据源请求
type TaskSourceReqStruct struct {
	SourceID               string `json:"source_id" form:"source_id" db:"source_id" binding:"required"`
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
	SourceHost             string `json:"source_host" form:"source_host" db:"source_host" binding:"required"`
	SourceUser             string `json:"source_user" form:"source_user" db:"source_user" binding:"required"`
	SourcePassword         string `json:"source_password" form:"source_password" db:"source_password" binding:"required"`
	SourcePort             int    `json:"source_port" form:"source_port" db:"source_port" binding:"required"`
	SourceSslCA            string `json:"source_ssl_ca" form:"source_ssl_ca" db:"source_ssl_ca"`
	SourceSslCert          string `json:"source_ssl_cert" form:"source_ssl_cert" db:"source_ssl_cert"`
	SourceSslKey           string `json:"source_ssl_key" form:"source_ssl_key" db:"source_ssl_key"`
}

// 任务下游数据源请求
type TaskTargetReqStruct struct {
	TargetID               string `json:"target_id" form:"target_id" db:"target_id" binding:"required"`
	TargetHost             string `json:"target_host" form:"target_host" db:"target_host" binding:"required"`
	TargetUser             string `json:"target_user" form:"target_user" db:"target_user" binding:"required"`
	TargetPassword         string `json:"target_password" form:"target_password" db:"target_password" binding:"required"`
	TargetPort             int    `json:"target_port" form:"target_port" db:"target_port" binding:"required"`
	TargetPacket           int    `json:"target_packet" form:"target_packet" db:"target_packet"`
	SqlMode                string `json:"sql_mode" form:"sql_mode" db:"sql_mode"`
	SkipUtf8Check          int    `json:"skip_utf_8_check" form:"skip_utf_8_check" db:"skip_utf_8_check"`
	ConstraintCheckInPlace int    `json:"constraint_check_in_place" form:"constraint_check_in_place" db:"constraint_check_in_place"`
	TargetSslCA            string `json:"target_ssl_ca" form:"target_ssl_ca" db:"target_ssl_ca"`
	TargetSslCert          string `json:"target_ssl_cert" form:"target_ssl_cert" db:"target_ssl_cert"`
	TargetSslKey           string `json:"target_ssl_key" form:"target_ssl_cert" db:"target_ssl_cert"`
}

// 任务同步元数据请求
type TaskMetaReqStruct struct {
	ClusterName         string `json:"cluster_name" form:"cluster_name" db:"cluster_name" binding:"required"`
	TaskName            string `json:"task_name" form:"task_name" db:"task_name" binding:"required"`
	TaskMode            string `json:"task_mode" form:"task_mode" db:"task_mode" binding:"required"`
	ShardMode           string `json:"shard_mode" form:"shard_mode" db:"shard_mode" binding:"required"`
	MetaSchema          string `json:"meta_schema" form:"meta_schema" db:"meta_schema" binding:"required"`
	Timezone            string `json:"timezone" form:"timezone" db:"timezone" binding:"required"`
	CaseSensitive       string `json:"case_sensitive" form:"case_sensitive" db:"case_sensitive" binding:"required"`
	OnlineDDL           string `json:"online_ddl" form:"online_ddl" db:"online_ddl" binding:"required"`
	OnlineDDLScheme     string `json:"online_ddl_scheme" form:"online_ddl_scheme" db:"online_ddl_scheme" binding:"required"`
	IgnoreCheckingItems string `json:"ignore_checking_items" form:"ignore_checking_items" db:"ignore_checking_items" binding:"required"`
	CleanDumpFile       string `json:"clean_dump_file" form:"clean_dump_file" db:"clean_dump_file" binding:"required"`
	TaskSourceID        string `json:"task_source_id" form:"task_source_id" db:"task_source_id" binding:"required"`
	TargetID            string `json:"target_id" form:"target_id" db:"target_id" binding:"required"`
}

// 任务同步路由规则请求
type TaskRouteRuleReqStruct struct {
	ClusterName   string `json:"cluster_name" form:"cluster_name" db:"cluster_name" binding:"required"`
	TaskName      string `json:"task_name" form:"task_name" db:"task_name" binding:"required"`
	RouteName     string `json:"route_name" form:"route_name" db:"route_name" binding:"required"`
	SchemaPattern string `json:"schema_pattern" form:"schema_pattern" db:"schema_pattern"`
	TablePattern  string `json:"table_pattern" form:"table_pattern" db:"table_pattern"`
	TargetSchema  string `json:"target_schema" form:"target_schema" db:"target_schema"`
	TargetTable   string `json:"target_table" form:"target_table" db:"target_table"`
	TaskSourceID  string `json:"task_source_id" form:"task_source_id" db:"task_source_id" binding:"required"`
}

// 任务同步过滤规则请求
type TaskFilterRuleReqStruct struct {
	ClusterName   string `json:"cluster_name" form:"cluster_name" db:"cluster_name" binding:"required"`
	TaskName      string `json:"task_name" form:"task_name" db:"task_name" binding:"required"`
	FilterName    string `json:"filter_name" form:"filter_name" db:"filter_name" binding:"required"`
	SchemaPattern string `json:"schema_pattern" form:"schema_pattern" db:"schema_pattern"`
	TablePattern  string `json:"table_pattern" form:"table_pattern" db:"table_pattern"`
	Events        string `json:"events" form:"events" db:"events"`
	SqlPattern    string `json:"sql_pattern" form:"sql_pattern" db:"sql_pattern"`
	Action        string `json:"action" form:"action" db:"action"`
	TaskSourceID  string `json:"task_source_id" form:"task_source_id" db:"task_source_id" binding:"required"`
}

// 任务同步行过滤规则请求
type TaskExpressionRuleReqStruct struct {
	ClusterName     string `json:"cluster_name" form:"cluster_name" db:"cluster_name" binding:"required"`
	TaskName        string `json:"task_name" form:"task_name" db:"task_name" binding:"required"`
	ExpressionName  string `json:"expression_name" form:"expression_name" db:"expression_name" binding:"required"`
	SchemaName      string `json:"schema_name" form:"schema_name" db:"schema_name"`
	TableName       string `json:"table_name" form:"table_name" db:"table_name"`
	InsertValueExpr string `json:"insert_value_expr" form:"insert_value_expr" db:"insert_value_expr"`
	TaskSourceID    string `json:"task_source_id" form:"task_source_id" db:"task_source_id" binding:"required"`
}

// 任务同步黑白名单规则请求
type TaskBlockAllowRuleReqStruct struct {
	ClusterName    string `json:"cluster_name" form:"cluster_name" db:"cluster_name" binding:"required"`
	TaskName       string `json:"task_name" form:"task_name" db:"task_name" binding:"required"`
	BlockAllowName string `json:"block_allow_name" form:"block_allow_name" db:"block_allow_name" binding:"required"`
	DoDBS          string `json:"do_dbs" form:"do_dbs" db:"do_dbs"`
	IgnoreDbs      string `json:"ignore_dbs" form:"ignore_dbs" db:"ignore_dbs"`
	DoTables       string `json:"do_tables" form:"do_tables" db:"do_tables"`
	IgnoreTables   string `json:"ignore_tables" form:"ignore_tables" db:"ignore_tables"`
	TaskSourceID   string `json:"task_source_id" form:"task_source_id" db:"task_source_id" binding:"required"`
}

// 任务同步导出请求
type TaskMydumperToolReqStruct struct {
	ClusterName   string `json:"cluster_name" form:"cluster_name" db:"cluster_name" binding:"required"`
	TaskName      string `json:"task_name" form:"task_name" db:"task_name" binding:"required"`
	RuleName      string `json:"rule_name" form:"rule_name" db:"rule_name" binding:"required"`
	Threads       int    `json:"threads" form:"threads" db:"threads"`
	ChunkFilesize int    `json:"chunk_filesize" form:"chunk_filesize" db:"chunk_filesize"`
	ExtraArgs     string `json:"extra_args" form:"extra_args" db:"extra_args"`
	TaskSourceID  string `json:"task_source_id" form:"task_source_id" db:"task_source_id" binding:"required"`
}

// 任务同步导入请求
type TaskLoaderToolReqStruct struct {
}

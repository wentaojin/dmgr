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
type TaskSourceCreateReqStruct struct {
	ClusterName string `db:"cluster_name"`
	TaskName    string `db:"task_name"`
	TaskSourceStruct
	TaskDatasourceStruct
	TaskDatasourceSslStruct
	Label string `json:"label" form:"label" db:"label"`
}

// 任务下游数据源创建请求
type TaskTargetCreateReqStruct struct {
	TaskTargetStruct
	TaskDatasourceStruct
	TaskDatasourceSslStruct
	Label string `json:"label" form:"label" db:"label"`
}

// 任务集群数据源映射请求
type TaskCLusterReqStruct struct {
	ClusterName string `json:"cluster_name" form:"cluster_name" db:"cluster_name" binding:"required"`
	TaskName    string `json:"task_name" form:"task_name" db:"task_name" binding:"required"`
	SourceName  string `json:"source_name" form:"source_name" db:"source_name" binding:"required"`
	TargetName  string `json:"target_name" form:"target_name" db:"target_name" binding:"required"`
	TaskSourceConfStruct
}

// 任务上游数据源删除请求
type TaskSourceDeleteReqStruct struct {
	ClusterName string `json:"cluster_name" form:"cluster_name" db:"cluster_name" binding:"required"`
	TaskName    string `json:"task_name" form:"task_name" db:"task_name" binding:"required"`
	SourceName  string `json:"source_name" form:"source_name" db:"source_name" binding:"required"`
}

// 任务上游数据源更新请求
type TaskSourceUpdateReqStruct struct {
	ClusterName string `json:"cluster_name" form:"cluster_name" binding:"required"`
	TaskName    string `json:"task_name" form:"task_name" binding:"required"`
	TaskSourceCreateReqStruct
	TaskSourceConfStruct
}

// 任务下游数据源删除请求
type TaskTargetDeleteReqStruct struct {
	ClusterName string `json:"cluster_name" form:"cluster_name" db:"cluster_name" binding:"required"`
	TaskName    string `json:"task_name" form:"task_name" db:"task_name" binding:"required"`
	TargetName  string `json:"target_name" form:"target_name" db:"target_name" binding:"required"`
}

// 任务下游数据源更新请求
type TaskTargetUpdateReqStruct struct {
	ClusterName string `json:"cluster_name" form:"cluster_name" binding:"required"`
	TaskName    string `json:"task_name" form:"task_name" binding:"required"`
	TaskTargetCreateReqStruct
}

// 任务同步元数据请求
type TaskMetaReqStruct struct {
	ClusterName         string `json:"cluster_name" form:"cluster_name" db:"cluster_name" binding:"required"`
	TaskName            string `json:"task_name" form:"task_name" db:"task_name" binding:"required"`
	ConfVersion         string `json:"conf_version" form:"conf_version" db:"conf_version"`
	TaskMode            string `json:"task_mode" form:"task_mode" db:"task_mode" binding:"required"`
	ShardMode           string `json:"shard_mode" form:"shard_mode" db:"shard_mode"`
	MetaSchema          string `json:"meta_schema" form:"meta_schema" db:"meta_schema"`
	Timezone            string `json:"timezone" form:"timezone" db:"timezone"`
	CaseSensitive       string `json:"case_sensitive" form:"case_sensitive" db:"case_sensitive"`
	OnlineDDL           string `json:"online_ddl" form:"online_ddl" db:"online_ddl" binding:"required"`
	IgnoreCheckingItems string `json:"ignore_checking_items" form:"ignore_checking_items" db:"ignore_checking_items"`
	CleanDumpFile       string `json:"clean_dump_file" form:"clean_dump_file" db:"clean_dump_file"`
	OnDuplication       string `json:"on_duplication" form:"on_duplication" db:"on_duplication"`
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
}

// 任务同步行过滤规则请求
type TaskExpressionRuleReqStruct struct {
	ClusterName     string `json:"cluster_name" form:"cluster_name" db:"cluster_name" binding:"required"`
	TaskName        string `json:"task_name" form:"task_name" db:"task_name" binding:"required"`
	ExpressionName  string `json:"expression_name" form:"expression_name" db:"expression_name" binding:"required"`
	SchemaName      string `json:"schema_name" form:"schema_name" db:"schema_name"`
	TableName       string `json:"table_name" form:"table_name" db:"table_name"`
	InsertValueExpr string `json:"insert_value_expr" form:"insert_value_expr" db:"insert_value_expr"`
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
}

// 任务同步导出请求
type TaskMydumperToolReqStruct struct {
	ClusterName   string `json:"cluster_name" form:"cluster_name" db:"cluster_name" binding:"required"`
	TaskName      string `json:"task_name" form:"task_name" db:"task_name" binding:"required"`
	RuleName      string `json:"rule_name" form:"rule_name" db:"rule_name" binding:"required"`
	Threads       int    `json:"threads" form:"threads" db:"threads"`
	ChunkFilesize int    `json:"chunk_filesize" form:"chunk_filesize" db:"chunk_filesize"`
	ExtraArgs     string `json:"extra_args" form:"extra_args" db:"extra_args"`
}

// 任务同步导入请求
type TaskLoaderToolReqStruct struct {
	ClusterName string `json:"cluster_name" form:"cluster_name" db:"cluster_name" binding:"required"`
	TaskName    string `json:"task_name" form:"task_name" db:"task_name" binding:"required"`
	RuleName    string `json:"rule_name" form:"rule_name" db:"rule_name" binding:"required"`
	PoolSize    int    `json:"pool_size" form:"pool_size" db:"pool_size"`
	DataDir     int    `json:"data_dir" form:"data_dir" db:"data_dir"`
}

// 任务增量同步请求
type TaskSyncerToolReqStruct struct {
	ClusterName      string `json:"cluster_name" form:"cluster_name" db:"cluster_name" binding:"required"`
	TaskName         string `json:"task_name" form:"task_name" db:"task_name" binding:"required"`
	RuleName         string `json:"rule_name" form:"rule_name" db:"rule_name" binding:"required"`
	WorkerCount      int    `json:"worker_count" form:"worker_count" db:"worker_count"`
	Batch            int    `json:"batch" form:"batch" db:"batch"`
	EnableAnsiQuotes string `json:"enable_ansi_quotes" form:"enable_ansi_quotes" db:"enable_ansi_quotes"`
	SafeMode         string `json:"safe_mode" form:"safe_mode" db:"safe_mode"`
}

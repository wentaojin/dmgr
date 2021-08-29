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
package service

// 任务同步管理数据表
const (
	TaskTables = `-- DM 上游数据源列表
CREATE TABLE IF NOT EXISTS task_source (
source_id varchar(30) NOT NULL COMMENT '源数据 ID',
enable_gtid varchar(10) NOT NULL DEFAULT 'true' COMMENT 'DM-Worker 是否使用 gtid',
relay_binlog_gtid varchar(255) DEFAULT NULL COMMENT '拉取上游 binlog 起始 gtid',
enable_relay varchar(10) NOT NULL DEFAULT 'false' COMMENT 'DM-Worker 是否使用 relay',
relay_binlog_name varchar(255) DEFAULT NULL COMMENT '拉取上游 binlog 起始文件名',
relay_dir varchar(255) NOT NULL DEFAULT './relay_log' COMMENT '存储 relay log 目录',
purge_interval int NOT NULL DEFAULT 3600 COMMENT '定期检查 relay log 是否过期的间隔时间',
purge_expires int NOT NULL DEFAULT 0 COMMENT 'relay log 的过期时间',
purge_remain_space int NOT NULL DEFAULT 15 COMMENT '设置最小的可用磁盘空间',
checker_check_enable varchar(10) NOT NULL DEFAULT 'true' COMMENT '启用自动重试功能',
checker_backoff_rollback varchar(10) NOT NULL DEFAULT '5m0s' COMMENT '如果指数回退策略的间隔大于该值，且任务处于正常状态，尝试减小间隔',
checker_backoff_max varchar(10) NOT NULL DEFAULT '5m0s' COMMENT '指数回退策略的间隔的最大值',
source_host varchar(255) NOT NULL COMMENT '源数据库用户',
source_user varchar(30) NOT NULL COMMENT '源数据库用户密码',
source_password varchar(255) NOT NULL COMMENT '源数据库密码',
source_port int NOT NULL COMMENT '源数据库端口',
source_ssl_ca varchar(255) DEFAULT NULL COMMENT '源数据库 SSL CA',
source_ssl_cert varchar(255) DEFAULT NULL COMMENT '源数据库 SSL CERT',
source_ssl_key varchar(255) DEFAULT NULL COMMENT '源数据库 SSL KEY',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (source_id),
UNIQUE INDEX idx_host_port (source_host,source_port)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_bin
COMMENT = '源数据列表';

-- DM 下游数据源列表
CREATE TABLE IF NOT EXISTS task_target (
target_id varchar(30) NOT NULL COMMENT '目标数据源 ID',
target_host varchar(255) NOT NULL COMMENT '目标数据库用户',
target_user varchar(30) NOT NULL COMMENT '目标数据库用户密码',
target_password varchar(255) NOT NULL COMMENT '目标数据库密码',
target_port int NOT NULL COMMENT '目标数据库端口',
target_packet bigint NOT NULL DEFAULT 67108864 COMMENT '目标数据库 max_allowed_packet',
sql_mode varchar(255) DEFAULT NULL COMMENT '目标数据库 SQL MODE',
skip_utf8_check int DEFAULT NULL COMMENT '目标数据库 tidb_skip_utf8_check',
constraint_check_in_place int DEFAULT NULL COMMENT '目标数据库 tidb_constraint_check_in_place',
target_ssl_key varchar(255) DEFAULT NULL COMMENT '目标数据库 SSL KEY',
target_ssl_ca varchar(255) DEFAULT NULL COMMENT '目标数据库 SSL CA',
target_ssl_cert varchar(255) DEFAULT NULL COMMENT '目标数据库 SSL CERT',
target_ssl_key varchar(255) DEFAULT NULL COMMENT '目标数据库 SSL KEY',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (target_id),
UNIQUE INDEX idx_host_port (target_host,target_port)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_bin
COMMENT = '目标数据源列表';

-- DM 同步任务数据列表
CREATE TABLE IF NOT EXISTS task_meta (
cluster_name varchar(30) NOT NULL COMMENT '集群名',
task_name varchar(30) NOT NULL COMMENT '任务名',
task_mode varchar(30) NOT NULL DEFAULT 'all' COMMENT '任务模式 "full" - "只进行全量数据迁移"、"incremental" - "Binlog 实时同步"、"all" - "全量 + Binlog 迁移"',
shard_mode varchar(30) NOT NULL DEFAULT NULL COMMENT '任务协调模式 ""、"pessimistic、"optimistic" 默认使用 ""',
meta_schema varchar(30) NOT NULL DEFAULT 'dm_meta' COMMENT '任务元数据库',
timezone varchar(30) NOT NULL DEFAULT 'Asia/Shanghai' COMMENT '时区',
case_sensitive varchar(30) NOT NULL DEFAULT 'false' COMMENT 'schema/table 是否大小写敏感',
online_ddl varchar(30) NOT NULL DEFAULT 'true' COMMENT '是否激活 online_ddl',
online_ddl_scheme varchar(30) NOT NULL DEFAULT 'pt' COMMENT 'online_ddl 模式，只支持 "gh-ost" 、"pt" 的自动处理',
ignore_checking_items varchar(125) NOT NULL DEFAULT NULL COMMENT '是否关闭任何检查项，默认""不关闭',
clean_dump_file varchar(30) NOT NULL DEFAULT 'true' COMMENT '是否清理 dump 阶段产生的文件',
task_source_id varchar(1024) NOT NULL COMMENT '源库实例名,格式 source1;source2 多个 source 分号分割',
target_id varchar(30) NOT NULL COMMENT '目标数据源 ID',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (cluster_name,task_name),
UNIQUE INDEX idx_host_port (target_host,target_port)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_bin
COMMENT = '任务元数据列表';

-- DM 同步任务路由规则列表
CREATE TABLE IF NOT EXISTS task_route_rule (
cluster_name varchar(30) NOT NULL COMMENT '集群名',
task_name varchar(30) NOT NULL COMMENT '任务名',
route_name varchar(30) NOT NULL COMMENT '路由名',
schema_pattern varchar(125) NOT NULL DEFAULT NULL COMMENT '源库名匹配规则',
table_pattern varchar(125) NOT NULL DEFAULT NULL COMMENT '源库表名匹配规则',
target_schema varchar(125) NOT NULL DEFAULT NULL COMMENT '目标库名称',
target_table varchar(125) NOT NULL DEFAULT NULL COMMENT '目标表名称',
task_source_id varchar(1024) NOT NULL DEFAULT NULL COMMENT '源库实例名,格式 source1;source2 多个 source 分号分割',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (cluster_name,task_name,route_name)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_bin
COMMENT = '任务路由规则列表';

-- DM 同步任务过滤规则列表
CREATE TABLE IF NOT EXISTS task_filter_rule (
cluster_name varchar(30) NOT NULL COMMENT '集群名',
task_name varchar(30) NOT NULL COMMENT '任务名',
filter_name varchar(30) NOT NULL COMMENT '过滤名',
schema_pattern varchar(125) NOT NULL DEFAULT NULL COMMENT '源库名匹配规则',
table_pattern varchar(125) NOT NULL DEFAULT NULL COMMENT '源库表名匹配规则',
events varchar(1024) NOT NULL DEFAULT NULL COMMENT '匹配上 schema-pattern 和 table-pattern 的库或者表的操作类型',
sql_pattern varchar(1024) NOT NULL DEFAULT NULL COMMENT '匹配上 schema-pattern 和 table-pattern 的库或者表的 sql 语句',
action varchar(125) NOT NULL DEFAULT NULL COMMENT '迁移（Do）还是忽略(Ignore)',
task_source_id varchar(1024) NOT NULL DEFAULT NULL COMMENT '源库实例名,格式 source1;source2 多个 source 分号分割',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (cluster_name,task_name,filter_name)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_bin
COMMENT = '任务数据过滤规则列表';

-- DM 同步任务行过滤规则列表
CREATE TABLE IF NOT EXISTS task_expression_rule (
cluster_name varchar(30) NOT NULL COMMENT '集群名',
task_name varchar(30) NOT NULL COMMENT '任务名',
expression_name varchar(30) NOT NULL COMMENT '过滤名',
schema_name varchar(125) NOT NULL DEFAULT NULL COMMENT '匹配的上游数据库库名，不支持通配符匹配或正则匹配',
table_name varchar(125) NOT NULL DEFAULT NULL COMMENT '匹配的上游表名，不支持通配符匹配或正则匹配',
insert_value_expr varchar(125) NOT NULL DEFAULT NULL COMMENT '匹配上 schema 和 table 的库表的操作类型',
task_source_id varchar(1024) NOT NULL DEFAULT NULL COMMENT '源库实例名,格式 source1;source2 格式 source1;source2 格式 source1;source2 多个 source 分号分割',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (cluster_name,task_name,expression_name)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_bin
COMMENT = '任务数据行过滤规则列表';

-- DM 同步任务黑白名单列表
CREATE TABLE IF NOT EXISTS task_block_allow_rule (
cluster_name varchar(30) NOT NULL COMMENT '集群名',
task_name varchar(30) NOT NULL COMMENT '任务名',
block_allow_name varchar(30) NOT NULL COMMENT '黑白规则名',
do_dbs varchar(125) NOT NULL DEFAULT NULL COMMENT '匹配的上游数据库库名迁移',
ignore_dbs varchar(125) NOT NULL DEFAULT NULL COMMENT '忽略匹配的上游数据库库名迁移',
do_tables varchar(1024) NOT NULL DEFAULT NULL COMMENT '匹配的上游数据库表名迁移，格式 dbName@tableName，多个表名以分号分割',
ignore_tables varchar(1024) NOT NULL DEFAULT NULL COMMENT '忽略匹配的上游数据库表名迁移，格式 dbName@tableName，多个表名以分号分割',    
task_source_id varchar(1024) NOT NULL DEFAULT NULL COMMENT '源库实例名,格式 source1;source2 格式 source1;source2 多个 source 分号分割',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (cluster_name,task_name,block_allow_name)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_bin
COMMENT = '任务黑白名单列表';

-- DM 同步任务导出列表
CREATE TABLE IF NOT EXISTS task_mydumper_tool (
cluster_name varchar(30) NOT NULL COMMENT '集群名',
task_name varchar(30) NOT NULL COMMENT '任务名',
rule_name varchar(125) NOT NULL DEFAULT 'global' COMMENT '工具配置规则名',
threads  int NOT NULL DEFAULT 4 COMMENT '数据导出并发',
chunk_filesize int NOT NULL DEFAULT 64 COMMENT '数据导出文件切分大小',
extra_ars varchar(1024) NOT NULL DEFAULT '--consistency none' COMMENT '数据导出其他参数配置',
task_source_id varchar(1024) NOT NULL DEFAULT NULL COMMENT '源库实例名,格式 source1;source2 格式 source1;source2 多个 source 分号分割',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (cluster_name,task_name,rule_name)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_bin
COMMENT = '任务导出列表';

-- DM 同步任务导入列表
CREATE TABLE IF NOT EXISTS task_loader_tool (
cluster_name varchar(30) NOT NULL COMMENT '集群名',
task_name varchar(30) NOT NULL COMMENT '任务名',
rule_name varchar(125) NOT NULL DEFAULT 'global' COMMENT '工具配置规则名',
pool_size  int NOT NULL DEFAULT 16 COMMENT '数据导出并发',
data_dir varchar(125) NOT NULL DEFAULT './dumped_data' COMMENT '数据目录',
task_source_id varchar(1024) NOT NULL DEFAULT NULL COMMENT '源库实例名,格式 source1;source2 格式 source1;source2 多个 source 分号分割',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (cluster_name,task_name,rule_name)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_bin
COMMENT = '任务导入列表';

-- DM 同步任务增量列表
CREATE TABLE IF NOT EXISTS task_syncer_tool (
cluster_name varchar(30) NOT NULL COMMENT '集群名',
task_name varchar(30) NOT NULL COMMENT '任务名',
rule_name varchar(125) NOT NULL DEFAULT 'global' COMMENT '工具配置规则名',
worker_count  int NOT NULL DEFAULT 16 COMMENT '数据同步并发',
batch int NOT NULL DEFAULT 100 COMMENT '数据同步 batch 大小',
enable_ansi_quotes varchar(10) NOT NULL DEFAULT 'true' COMMENT '目标库连接中 session 设置 sql-mode: "ANSI_QUOTES"，则需开启此项',
safe_mode varchar(10) NOT NULL DEFAULT 'false' COMMENT '数据同步模式是否开启 safe-mode',
task_source_id varchar(1024) NOT NULL DEFAULT NULL COMMENT '源库实例名,格式 source1;source2 格式 source1;source2 多个 source 分号分割',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (cluster_name,task_name,rule_name)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_bin
COMMENT = '任务导出列表';`
)

// 集群管理数据表
const (
	ClusterTables = `CREATE TABLE IF NOT EXISTS user (
id int NOT NULL AUTO_INCREMENT COMMENT '用户 ID',
username varchar(30) NOT NULL COMMENT '用户名',
password varchar(255) NOT NULL COMMENT '用户密码',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (id) ,
UNIQUE INDEX idx_username (username)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_bin
COMMENT = '用户列表';

CREATE TABLE IF NOT EXISTS machine (
id bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
ssh_host varchar(255) NOT NULL COMMENT 'SSH 主机',
ssh_user varchar(30) NOT NULL COMMENT 'SSH 用户',
ssh_password varchar(255) NOT NULL COMMENT 'SSH密码',
ssh_port int NOT NULL COMMENT 'SSH 端口',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (id) ,
UNIQUE INDEX idx_ssh_host_port (ssh_host, ssh_port)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COMMENT = '机器列表';

CREATE TABLE IF NOT EXISTS warehouse (
id bigint NOT NULL AUTO_INCREMENT COMMENT '元数据 ID',
cluster_version varchar(30) NOT NULL COMMENT '集群版本',
package_name varchar(255) NOT NULL COMMENT '离线包名',
package_path varchar(255) NOT NULL COMMENT '离线包存放路径',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (id) ,
UNIQUE INDEX idx_cluster_version (cluster_version),
UNIQUE INDEX idx_package_name (package_name)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_bin
COMMENT = '离线包镜像包仓库列表';

	CREATE TABLE IF NOT EXISTS cluster_meta (
id bigint NOT NULL AUTO_INCREMENT COMMENT '元数据 ID',
cluster_name varchar(255) NOT NULL COMMENT '集群名',
cluster_user varchar(30) NOT NULL COMMENT '集群用户',
cluster_version varchar(30) NOT NULL COMMENT '集群版本',
cluster_path  varchar(255) NOT NULL COMMENT '集群离线包解压路径',
cluster_status varchar(10) NOT NULL DEFAULT 'Offline' COMMENT '集群状态是否启动，Up 启动; Offline 离线',
admin_user  varchar(255) NOT NULL COMMENT 'grafana 用户名',
admin_password  varchar(255) NOT NULL COMMENT 'grafana 用户密码',
skip_create_user varchar(30) NOT NULL DEFAULT 'false' COMMENT '是否创建集群用户名, false 不跳过; true 跳过',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (id) ,
UNIQUE INDEX idx_cluster_name (cluster_name),
INDEX idx_cluster_version (cluster_version),
INDEX idx_cluster_status (cluster_status)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_bin
COMMENT = '集群元数据表';

CREATE TABLE IF NOT EXISTS cluster_topology (
id bigint NOT NULL AUTO_INCREMENT COMMENT '拓扑 ID',
cluster_name varchar(255) NOT NULL COMMENT '集群名',
component_name varchar(35) NOT NULL COMMENT '集群组件名',
instance_name varchar(35) NOT NULL COMMENT '部署实例名',
service_port int NOT NULL COMMENT '实例服务端口',
peer_port int NOT NULL DEFAULT 8291 COMMENT 'DM Master 状态端口',
cluster_port int NOT NULL DEFAULT 9094 COMMENT 'Alertmanager 集群端口',
deploy_dir varchar(255) NOT NULL COMMENT '实例部署目录',
data_dir varchar(255) NOT NULL COMMENT '实例数据目录',
log_dir varchar(255) NOT NULL COMMENT '实例日志目录',
machine_host varchar(255) NOT NULL COMMENT '实例部署主机',
hotfix varchar(255) NOT NULL DEFAULT 'Normal' COMMENT 'hotfix 状态 Normal 正常状态，Reload 配置变更状态，Patched 补丁状态',
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (id) ,
INDEX idx_cluster_name (cluster_name),
UNIQUE INDEX idx_cluster_instance_name (cluster_name,instance_name),
INDEX idx_machine_host (machine_host)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_bin
COMMENT = '集群部署拓扑列表';`
)

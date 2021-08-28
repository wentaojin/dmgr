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

const (
	Schema = `CREATE DATABASE IF NOT EXISTS dmgr`

	Tables = `CREATE TABLE IF NOT EXISTS user (
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

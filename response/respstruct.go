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
package response

import (
	"time"
)

// 用户登陆响应
type LoginRespStruct struct {
	Token     string    `json:"token"`     // jwt令牌
	ExpiresAt time.Time `json:"expiresAt"` // 过期时间, 秒
}

// 用户信息响应
type UserRespStruct struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

// 机器列表响应
type MachineRespStruct struct {
	SshHost     string `json:"ssh_host" db:"ssh_host"`
	SshPort     uint64 `json:"ssh_port" db:"ssh_port"`
	SshUser     string `json:"ssh_user" db:"ssh_user"`
	SshPassword string `json:"ssh_password" db:"ssh_password"`
}

// 集群离线包响应
type WarehouseRespStruct struct {
	ClusterVersion string `json:"cluster_version" db:"cluster_version"`
	PackageName    string `json:"package_name" db:"package_name"`
	PackagePath    string `json:"package_path" db:"package_path"`
}

// 集群拓扑响应
type ClusterTopologyRespStruct struct {
	ClusterName    string `json:"cluster_name" db:"cluster_name"`
	ClusterUser    string `json:"cluster_user" db:"cluster_user"`
	ClusterVersion string `json:"cluster_version" db:"cluster_version"`
	ClusterPath    string `json:"cluster_path" db:"cluster_path"`
	AdminUser      string `json:"admin_user" db:"admin_user"`
	AdminPassword  string `json:"admin_password" db:"admin_password"`
	SshUser        string `json:"ssh_user" db:"ssh_user"`
	SshPassword    string `json:"ssh_password" db:"ssh_password"`
	SshPort        uint64 `json:"ssh_port" db:"ssh_port"`
	ComponentName  string `json:"component_name" db:"component_name"`
	InstanceName   string `json:"instance_name" db:"instance_name"`
	MachineHost    string `json:"machine_host" db:"machine_host"`
	ServicePort    uint64 `json:"service_port" db:"service_port"`
	PeerPort       uint64 `json:"peer_port" db:"peer_port"`
	ClusterPort    uint64 `json:"cluster_port" db:"cluster_port"`
	DeployDir      string `json:"deploy_dir" db:"deploy_dir"`
	DataDir        string `json:"data_dir" db:"data_dir"`
	LogDir         string `json:"log_dir" db:"log_dir"`
}

// 集群状态响应
type ClusterMetaRespStruct struct {
	ClusterName    string `json:"cluster_name" db:"cluster_name"`
	ClusterUser    string `json:"cluster_user" db:"cluster_user"`
	ClusterVersion string `json:"cluster_version" db:"cluster_version"`
	ClusterPath    string `json:"cluster_path" db:"cluster_path"`
	ClusterStatus  string `json:"cluster_status" db:"cluster_status"`
	AdminUser      string `json:"admin_user" db:"admin_user"`
	AdminPassword  string `json:"admin_password" db:"admin_password"`
}

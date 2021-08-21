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

import "mime/multipart"

// 用户注册登录请求
type RegisterAndLoginReqStruct struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required,gte=5"`
}

// 修改密码请求
type ChangePwdReqStruct struct {
	OldPassword string `json:"oldPassword" form:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" form:"newPassword" binding:"required"`
}

// 机器新增请求
type MachineReqStruct struct {
	SshHost     string `json:"ssh_host" form:"ssh_host" binding:"required"`
	SshUser     string `json:"ssh_user" form:"ssh_user" binding:"required"`
	SshPassword string `json:"ssh_password" form:"ssh_password" binding:"required"`
	SshPort     uint64 `json:"ssh_port" form:"ssh_port" binding:"required"`
}

// 离线镜像包新增请求
type PackageReqStruct struct {
	ClusterVersion string                `json:"cluster_version" form:"cluster_version" binding:"required"`
	PackageName    string                `json:"package_name" form:"package_name" binding:"required"`
	PackagePath    string                `json:"package_path" form:"package_path" binding:"required"`
	File           *multipart.FileHeader `json:"file" form:"file" binding:"required"`
}

// 集群拓扑部署请求
type ClusterDeployReqStruct struct {
	ClusterMetaReqStruct
	ClusterTopology []TopologyReqStruct `json:"cluster_topology" form:"cluster_topology" binding:"required"`
}

// 集群元数据请求
type ClusterMetaReqStruct struct {
	ClusterName    string `json:"cluster_name" form:"cluster_name" binding:"required" db:"cluster_name"`
	ClusterUser    string `json:"cluster_user" form:"cluster_user" binding:"required" db:"cluster_user"`
	ClusterVersion string `json:"cluster_version" form:"cluster_version" binding:"required" db:"cluster_version"`
	ClusterPath    string `json:"cluster_path" form:"cluster_path" binding:"required" db:"cluster_path"`
	AdminUser      string `json:"admin_user" form:"admin_user" binding:"required" db:"admin_user"`
	AdminPassword  string `json:"admin_password" form:"admin_password" binding:"required" db:"admin_password"`
	SkipCreateUser string `json:"skip_create_user" form:"skip_create_user" binding:"validIsSkip" db:"skip_create_user"`
}

type TopologyReqStruct struct {
	ClusterName   string `json:"cluster_name" form:"cluster_name" db:"cluster_name"`
	ComponentName string `json:"component_name" form:"component_name" binding:"required" db:"component_name"`
	InstanceName  string `json:"instance_name" form:"instance_name" binding:"required" db:"instance_name"`
	MachineHost   string `json:"machine_host" form:"machine_host" binding:"required" db:"machine_host"`
	ServicePort   uint64 `json:"service_port" form:"service_port" binding:"required" db:"service_port"`
	PeerPort      uint64 `json:"peer_port" form:"peer_port" db:"peer_port"`
	ClusterPort   uint64 `json:"cluster_port" form:"cluster_port" db:"cluster_port"`
	DeployDir     string `json:"deploy_dir" form:"deploy_dir" binding:"required" db:"deploy_dir"`
	DataDir       string `json:"data_dir" form:"data_dir" db:"data_dir"`
	LogDir        string `json:"log_dir" form:"log_dir" db:"log_dir"`
}

// 集群状态请求
type ClusterStatusReqStruct struct {
	ClusterStatus string `json:"cluster_status" form:"cluster_status"`
}

// 集群启动、停止或销毁请求
type ClusterOperatorReqStruct struct {
	ClusterName   string   `json:"cluster_name" form:"cluster_name" binding:"required"`
	ComponentName []string `json:"component_name" form:"component_name"`
	InstanceName  []string `json:"instance_name" form:"instance_name"`
}

// 集群扩容请求
type CLusterScaleOutReqStruct struct {
	TopologyReqStruct
	SkipCreateUser string `json:"skip_create_user" form:"skip_create_user" binding:"validIsSkip"`
	AdminUser      string `json:"admin_user" form:"admin_user" `        // 扩容 grafana 才需使用
	AdminPassword  string `json:"admin_password" form:"admin_password"` // 扩容 grafana 才需使用
}

// 集群升级请求
type ClusterUpgradeReqStruct struct {
	ClusterName    string `json:"cluster_name" form:"cluster_name" binding:"required"`
	ClusterVersion string `json:"cluster_version" form:"cluster_version" binding:"required"`
}

// 集群补丁请求
type ClusterPatchReqStruct struct {
	ClusterName   string                `json:"cluster_name" form:"cluster_name" binding:"required"`
	ComponentName string                `json:"component_name" form:"component_name"`
	InstanceName  []string              `json:"instance_name" form:"instance_name"`
	Overwrite     string                `json:"overwrite" form:"overwrite" binding:"validIsSkip"`
	File          *multipart.FileHeader `json:"file" form:"file" binding:"required"`
}

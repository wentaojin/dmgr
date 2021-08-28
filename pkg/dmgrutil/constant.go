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

const (
	// Bool 字符值
	BoolTrue  = "true"
	BoolFalse = "false"

	// 默认 grafana 用户密码
	DefaultGrafanaUser     = "admin"
	DefaultGrafanaPassword = "admin"

	// 默认端口
	DefaultDmMasterPeerPort        = 8291
	DefaultAlertmanagerClusterPort = 9094

	// 存放离线镜像包软件
	// {root_dir}/soft/{package_name}
	DirSoft = "soft"
	// 存放已部署的集群
	// {root_dir}/cluster/{cluster_name}/{cluster_version}
	DirCluster = "cluster"
	// 集群部署相对目录名
	DirBin    = "bin"
	DirConf   = "conf"
	DirScript = "script"
	DirData   = "data"
	DirLog    = "log"
	// 补丁目录
	DirPatch = "patch"
	// 缓存目录
	DirCache = "cache"
	// SSH 存放目录
	DirSSH = "ssh"
	// grafana dashboard 目录
	DirGrafanaDashboard = "provisioning/dashboards"
	// grafana datasource 目录
	DirGrafanaDatasource = "provisioning/datasources"
	// systemd 进程存放目录
	DirSystemd = "/etc/systemd/system"

	// 组件名 -> 组件离线包名前缀
	ComponentDmMaster      = "dm-master"
	ComponentDmWorker      = "dm-worker"
	ComponentAlertmanager  = "alertmanager"
	ComponentGrafana       = "grafana"
	ComponentPrometheus    = "prometheus"
	ComponentGrafanaTarPKG = "grafana.tar.gz"

	// DM 集群压缩包内容
	TmplGrafanaINIYAML        = "template/grafana.ini.tmpl"
	TmplGrafanaDashboardYAML  = "template/dashboard.yml.tmpl"
	TmplGrafanaDatasourceYAML = "template/datasource.yml.tmpl"
	TmplGrafanaScript         = "template/run_grafana.sh.tmpl"

	TmplPrometheusYAML   = "template/prometheus.yml.tmpl"
	TmplPrometheusScript = "template/run_prometheus.sh.tmpl"

	TmplAlertmanagerScript  = "template/run_alertmanager.sh.tmpl"
	TmplDmMasterScript      = "template/run_dm-master.sh.tmpl"
	TmplDmMasterScaleScript = "template/run_dm-master-scale.sh.tmpl"
	TmplDmWorkerScript      = "template/run_dm-worker.sh.tmpl"
	TmplSystemdScript       = "template/systemd.service.tmpl"

	ConfDmWorkerRuleFile = "conf/dm_worker.rules.yml"
	ConfAlertmanagerFile = "conf/alertmanager.yml"
	ConfDmMasterFile     = "conf/dm-master.toml"
	ConfDmWorkerFile     = "conf/dm-worker.toml"

	// 集群状态
	ClusterUpStatus      = "Up"
	ClusterOfflineStatus = "Offline"

	// HOME SSH
	HomeSshDir = "~/.ssh"

	// 密钥生成并发
	RsaConcurrency = 10

	// 用于 copy file 类型区分
	FileTypeComponent = "component"
	FileTypeScript    = "script"
	FileTypeSystemd   = "systemd"
	FileTypeRule      = "rule"

	// 用于指定组件 reload 滚更
	ReloadAlertmanagerFile = "alertmanager.yml"
	ReloadDmMasterFile     = "dm-master.toml"
	ReloadDmWorkerFile     = "dm-worker.toml"

	// 组件 hotfix 状态
	NormalComponent  = "Normal"  // 未打补丁以及配置变更
	ReloadComponent  = "Reload"  // 存在配置变更（可能打过补丁）
	PatchedComponent = "Patched" //存在补丁状态（可能配置变更过）
)

var (
	// 集群组件启动顺序
	StartComponentOrder = []string{ComponentDmMaster, ComponentDmWorker, ComponentPrometheus, ComponentGrafana, ComponentAlertmanager}
	StopComponentOrder  = []string{ComponentAlertmanager, ComponentGrafana, ComponentPrometheus, ComponentDmWorker, ComponentDmMaster}
)

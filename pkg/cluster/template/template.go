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
package template

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wentaojin/dmgr/response"

	"github.com/wentaojin/dmgr/pkg/cluster/template/systemd"

	"github.com/wentaojin/dmgr/pkg/cluster/template/script"

	"github.com/wentaojin/dmgr/pkg/cluster/template/config"
	"github.com/wentaojin/dmgr/pkg/dmgrutil"
)

const (
	ClusterDeployStage   = "deploy"
	ClusterScaleOutStage = "scale-out"
)

// 用于集群部署配置文件、运行脚本等文件生成
type ClusterOperatorStage struct {
	DmMasterScripts     []*script.DMMasterScript
	AlertmanagerScripts []*script.AlertManagerScript
	AlertmanagerAddrs   []string
	DmMasterAddrs       []string
	DmWorkerAddrs       []string
	GrafanaAddr         string
}

// 获取生成集群部署配置文件、运行脚本等文件信息
func GetClusterFile(topo []response.ClusterTopologyRespStruct) ClusterOperatorStage {
	var cos ClusterOperatorStage
	for _, cluster := range topo {
		switch strings.ToLower(cluster.ComponentName) {
		case dmgrutil.ComponentGrafana:
			cos.GrafanaAddr = fmt.Sprintf("%s:%v", cluster.MachineHost, cluster.ServicePort)
		case dmgrutil.ComponentDmMaster:
			cos.DmMasterAddrs = append(cos.DmMasterAddrs, fmt.Sprintf("%s:%v", cluster.MachineHost, cluster.ServicePort))
			cos.DmMasterScripts = append(cos.DmMasterScripts, &script.DMMasterScript{
				Name:      cluster.InstanceName,
				Scheme:    "http",
				IP:        cluster.MachineHost,
				Port:      cluster.ServicePort,
				PeerPort:  cluster.PeerPort,
				DeployDir: dmgrutil.AbsClusterDeployDir(cluster.DeployDir, cluster.InstanceName),
				DataDir:   dmgrutil.AbsClusterDataDir(cluster.DeployDir, cluster.DataDir, cluster.InstanceName),
				LogDir:    dmgrutil.AbsClusterLogDir(cluster.DeployDir, cluster.LogDir, cluster.InstanceName),
			})
		case dmgrutil.ComponentDmWorker:
			cos.DmWorkerAddrs = append(cos.DmWorkerAddrs, fmt.Sprintf("%s:%v", cluster.MachineHost, cluster.ServicePort))
		case dmgrutil.ComponentAlertmanager:
			cos.AlertmanagerAddrs = append(cos.AlertmanagerAddrs, fmt.Sprintf("%s:%v", cluster.MachineHost, cluster.ServicePort))
			cos.AlertmanagerScripts = append(cos.AlertmanagerScripts, &script.AlertManagerScript{
				IP:          cluster.MachineHost,
				WebPort:     cluster.ServicePort,
				ClusterPort: cluster.ClusterPort,
				DeployDir:   dmgrutil.AbsClusterDeployDir(cluster.DeployDir, cluster.InstanceName),
				DataDir:     dmgrutil.AbsClusterDataDir(cluster.DeployDir, cluster.DataDir, cluster.InstanceName),
				LogDir:      dmgrutil.AbsClusterLogDir(cluster.DeployDir, cluster.LogDir, cluster.InstanceName),
				TLSEnabled:  false,
			})
		}

	}

	return cos
}

// 生成集群配置文件 - 集群阶段
func GenerateClusterFileWithStage(
	topo []response.ClusterTopologyRespStruct,
	dmMasterScripts []*script.DMMasterScript,
	alertmanagerScripts []*script.AlertManagerScript,
	alertmanagerAddrs []string,
	dmMasterAddrs []string,
	dmWorkerAddrs []string,
	grafanaAddr string,
	clusterStage string,
	adminUser, adminPassword string) error {

	for _, t := range topo {
		switch strings.ToLower(t.ComponentName) {
		case dmgrutil.ComponentDmMaster:
			if clusterStage == ClusterDeployStage {
				if err := script.NewDMMasterScript(t.InstanceName, t.MachineHost,
					dmgrutil.AbsClusterDeployDir(t.DeployDir, t.InstanceName),
					dmgrutil.AbsClusterDataDir(t.DeployDir, t.DataDir, t.InstanceName),
					dmgrutil.AbsClusterLogDir(t.DeployDir, t.LogDir, t.InstanceName)).
					AppendEndpoints(dmMasterScripts...).
					WithPort(t.ServicePort).
					WithPeerPort(t.PeerPort).
					WithScheme("http").
					ConfigToFile(
						filepath.Join(
							dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
							t.ClusterVersion,
							dmgrutil.TmplDmMasterScript),
						filepath.Join(
							dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
							fmt.Sprintf("run_dm-master-%s-%d.sh", t.MachineHost, t.ServicePort))); err != nil {
					return err
				}
			}
			if clusterStage == ClusterScaleOutStage {
				if err := script.NewDMMasterScaleScript(t.InstanceName, t.MachineHost,
					dmgrutil.AbsClusterDeployDir(t.DeployDir, t.InstanceName),
					dmgrutil.AbsClusterDataDir(t.DeployDir, t.DataDir, t.InstanceName),
					dmgrutil.AbsClusterLogDir(t.DeployDir, t.LogDir, t.InstanceName)).
					WithPort(t.ServicePort).
					AppendEndpoints(dmMasterScripts...).
					WithPeerPort(t.PeerPort).
					WithScheme("http").
					ConfigToFile(
						filepath.Join(
							dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
							t.ClusterVersion,
							dmgrutil.TmplDmMasterScaleScript),
						filepath.Join(
							dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
							fmt.Sprintf("run_dm-master-%s-%d.sh", t.MachineHost, t.ServicePort))); err != nil {
					return err
				}
			}
			if err := systemd.NewSystemdConfig(strings.ToLower(t.ComponentName), t.ClusterUser,
				dmgrutil.AbsClusterDeployDir(t.DeployDir, t.InstanceName)).ConfigToFile(
				filepath.Join(
					dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
					t.ClusterVersion,
					dmgrutil.TmplSystemdScript),
				filepath.Join(
					dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
					fmt.Sprintf("%s-%s-%v.service", strings.ToLower(t.ComponentName), t.MachineHost, t.ServicePort))); err != nil {
				return err
			}
		case dmgrutil.ComponentDmWorker:
			if err := script.NewDMWorkerScript(t.InstanceName, t.MachineHost,
				dmgrutil.AbsClusterDeployDir(t.DeployDir, t.InstanceName),
				dmgrutil.AbsClusterLogDir(t.DeployDir, t.LogDir, t.InstanceName)).
				WithPort(t.ServicePort).AppendEndpoints(dmMasterScripts...).ConfigToFile(
				filepath.Join(
					dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
					t.ClusterVersion,
					dmgrutil.TmplDmWorkerScript),
				filepath.Join(
					dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
					fmt.Sprintf("run_dm-worker-%s-%d.sh", t.MachineHost, t.ServicePort))); err != nil {
				return err
			}

			if err := systemd.NewSystemdConfig(strings.ToLower(t.ComponentName), t.ClusterUser,
				dmgrutil.AbsClusterDeployDir(t.DeployDir, t.InstanceName)).ConfigToFile(
				filepath.Join(
					dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
					t.ClusterVersion,
					dmgrutil.TmplSystemdScript),
				filepath.Join(
					dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
					fmt.Sprintf("%s-%s-%v.service", strings.ToLower(t.ComponentName), t.MachineHost, t.ServicePort))); err != nil {
				return err
			}
		case dmgrutil.ComponentGrafana:
			var grafanaUser, grafanaPassword string

			// 用于集群扩容阶段 -》 扩容阶段如果未指定 grafana 用户密码，则使用数据库中已有的 grafana 用户密码
			if adminUser == "" {
				grafanaUser = t.AdminUser
			} else {
				grafanaUser = adminUser
			}
			if adminPassword == "" {
				grafanaPassword = t.AdminPassword
			} else {
				grafanaUser = adminPassword
			}

			if err := config.NewGrafanaConfig(t.MachineHost, dmgrutil.AbsClusterDataDir(t.DeployDir, t.DataDir, t.InstanceName), dmgrutil.AbsClusterLogDir(t.DeployDir, t.LogDir, t.InstanceName)).
				WithPort(t.ServicePort).
				WithUsername(grafanaUser).
				WithPassword(grafanaPassword).
				WithAnonymousenable(false).
				WithRootURL("").
				WithDomain("").
				ConfigToFile(
					filepath.Join(
						dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
						t.ClusterVersion,
						dmgrutil.TmplGrafanaINIYAML),
					filepath.Join(
						dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
						"grafana.ini")); err != nil {
				return err
			}

			if err := config.NewDatasourceConfig(t.ClusterName, t.MachineHost).WithPort(t.ServicePort).ConfigToFile(
				filepath.Join(
					dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
					t.ClusterVersion,
					dmgrutil.TmplGrafanaDatasourceYAML),
				filepath.Join(
					dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
					"datasource.yml")); err != nil {
				return err
			}

			if err := config.NewDashboardConfig(t.ClusterName, t.DeployDir).ConfigToFile(
				filepath.Join(
					dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
					t.ClusterVersion,
					dmgrutil.TmplGrafanaDashboardYAML),
				filepath.Join(
					dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
					"dashboard.yml")); err != nil {
				return err
			}

			if err := script.NewGrafanaScript(t.ClusterName, dmgrutil.AbsClusterDeployDir(t.DeployDir, t.InstanceName)).ConfigToFile(
				filepath.Join(
					dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
					t.ClusterVersion,
					dmgrutil.TmplGrafanaScript),
				filepath.Join(
					dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
					"run_grafana.sh")); err != nil {
				return err

			}

			if err := systemd.NewSystemdConfig(strings.ToLower(t.ComponentName), t.ClusterUser,
				dmgrutil.AbsClusterDeployDir(t.DeployDir, t.InstanceName)).ConfigToFile(
				filepath.Join(
					dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
					t.ClusterVersion,
					dmgrutil.TmplSystemdScript),
				filepath.Join(
					dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
					fmt.Sprintf("%s-%s-%v.service", strings.ToLower(t.ComponentName), t.MachineHost, t.ServicePort))); err != nil {
				return err
			}
		case dmgrutil.ComponentPrometheus:
			if err := config.NewPrometheusConfig(t.ClusterName, t.ClusterVersion, false).
				AddAlertmanager(alertmanagerAddrs).
				AddDMMaster(dmMasterAddrs).
				AddDMWorker(dmWorkerAddrs).
				AddGrafana(grafanaAddr).
				SetRemoteConfig("").
				ConfigToFile(
					filepath.Join(
						dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
						t.ClusterVersion,
						dmgrutil.TmplPrometheusYAML),
					filepath.Join(
						dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
						"prometheus.yml")); err != nil {
				return err
			}

			if err := script.NewPrometheusScript(t.MachineHost,
				dmgrutil.AbsClusterDeployDir(t.DeployDir, t.InstanceName),
				dmgrutil.AbsClusterDataDir(t.DeployDir, t.DataDir, t.InstanceName),
				dmgrutil.AbsClusterLogDir(t.DeployDir, t.LogDir, t.InstanceName),
			).
				WithPort(t.ServicePort).
				WithRetention("").
				ConfigToFile(
					filepath.Join(
						dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
						t.ClusterVersion,
						dmgrutil.TmplPrometheusScript),
					filepath.Join(
						dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
						"run_prometheus.sh")); err != nil {
				return err
			}

			if err := systemd.NewSystemdConfig(strings.ToLower(t.ComponentName), t.ClusterUser,
				dmgrutil.AbsClusterDeployDir(t.DeployDir, t.InstanceName)).ConfigToFile(
				filepath.Join(
					dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
					t.ClusterVersion,
					dmgrutil.TmplSystemdScript),
				filepath.Join(
					dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
					fmt.Sprintf("%s-%s-%v.service", strings.ToLower(t.ComponentName), t.MachineHost, t.ServicePort))); err != nil {
				return err
			}
		case dmgrutil.ComponentAlertmanager:
			if err := script.NewAlertManagerScript(t.MachineHost,
				dmgrutil.AbsClusterDeployDir(t.DeployDir, t.InstanceName),
				dmgrutil.AbsClusterDataDir(t.DeployDir, t.DataDir, t.InstanceName),
				dmgrutil.AbsClusterLogDir(t.DeployDir, t.LogDir, t.InstanceName), false).
				AppendEndpoints(alertmanagerScripts).WithClusterPort(t.ClusterPort).WithWebPort(t.ServicePort).
				ConfigToFile(
					filepath.Join(
						dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
						t.ClusterVersion,
						dmgrutil.TmplAlertmanagerScript),
					filepath.Join(
						dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
						fmt.Sprintf("run_alertmanager-%s-%d.sh", t.MachineHost, t.ServicePort))); err != nil {
				return err
			}

			if err := systemd.NewSystemdConfig(strings.ToLower(t.ComponentName), t.ClusterUser,
				dmgrutil.AbsClusterDeployDir(t.DeployDir, t.InstanceName)).ConfigToFile(
				filepath.Join(
					dmgrutil.AbsClusterUntarDir(t.ClusterPath, t.ClusterName),
					t.ClusterVersion,
					dmgrutil.TmplSystemdScript),
				filepath.Join(
					dmgrutil.AbsClusterCacheDir(t.ClusterPath, t.ClusterName),
					fmt.Sprintf("%s-%s-%v.service", strings.ToLower(t.ComponentName), t.MachineHost, t.ServicePort))); err != nil {
				return err
			}
		default:
			return fmt.Errorf("component [%v] not exist, panic", t.ComponentName)
		}
	}
	return nil
}

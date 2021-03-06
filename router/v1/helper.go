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
package v1

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wentaojin/dmgr/pkg/cluster/api"
	"github.com/wentaojin/dmgr/service"

	"github.com/wentaojin/dmgr/request"

	"github.com/wentaojin/dmgr/pkg/cluster/executor"
	"github.com/wentaojin/dmgr/pkg/cluster/task"
	"github.com/wentaojin/dmgr/pkg/dmgrutil"
	"github.com/wentaojin/dmgr/response"
)

// 用于集群拓扑生成 - Deploy/Scale-out 阶段
func GenerateClusterTopology(clusterMeta interface{}, topoReq []request.TopologyReqStruct, machineList []response.MachineRespStruct) ([]response.ClusterTopologyRespStruct, error) {
	var clusterTopo []response.ClusterTopologyRespStruct
	for _, topo := range topoReq {
		for _, machine := range machineList {
			if topo.MachineHost == machine.SshHost {
				t := response.ClusterTopologyRespStruct{
					SshUser:       machine.SshUser,
					SshPassword:   machine.SshPassword,
					SshPort:       machine.SshPort,
					ComponentName: topo.ComponentName,
					InstanceName:  topo.InstanceName,
					MachineHost:   topo.MachineHost,
					ServicePort:   topo.ServicePort,
					PeerPort:      topo.PeerPort,
					ClusterPort:   topo.ClusterPort,
					DeployDir:     topo.DeployDir,
					DataDir:       topo.DataDir,
					LogDir:        topo.LogDir,
				}

				switch v := clusterMeta.(type) {
				case request.ClusterMetaReqStruct:
					t.ClusterName = v.ClusterName
					t.ClusterUser = v.ClusterUser
					t.ClusterVersion = v.ClusterVersion
					t.ClusterPath = v.ClusterPath
					t.AdminUser = v.AdminUser
					t.AdminPassword = v.AdminPassword
				case response.ClusterMetaRespStruct:
					t.ClusterName = v.ClusterName
					t.ClusterUser = v.ClusterUser
					t.ClusterVersion = v.ClusterVersion
					t.ClusterPath = v.ClusterPath
					t.AdminUser = v.AdminUser
					t.AdminPassword = v.AdminPassword
				default:
					return clusterTopo, fmt.Errorf("component [%v] instance [%v] host [%v] assert failed", topo.ComponentName, topo.InstanceName, topo.MachineHost)
				}
				clusterTopo = append(clusterTopo, t)
			}
		}
	}
	return clusterTopo, nil
}

// 用于 COPY 集群文件的任务
func CopyClusterFile(clusterTopo []response.ClusterTopologyRespStruct) []task.Task {
	var copyFileTasks []task.Task

	for _, cluster := range clusterTopo {
		copyFileTask := task.NewBuilder().
			UserSSH(
				cluster.MachineHost,
				cluster.SshPort,
				cluster.ClusterUser,
				executor.DefaultConnectTimeout,
				executor.DefaultExecuteTimeout,
			)
		componentName := strings.ToLower(cluster.ComponentName)
		if componentName == dmgrutil.ComponentGrafana {
			copyFileTask.CopyFile(
				cluster.ClusterName,
				filepath.Join(dmgrutil.AbsClusterCacheDir(cluster.ClusterPath, cluster.ClusterName), "grafana.ini"),
				filepath.Join(dmgrutil.AbsClusterConfDir(cluster.DeployDir, cluster.InstanceName), "grafana.ini"),
				dmgrutil.FileTypeComponent,
				cluster.MachineHost,
				false,
				0).
				CopyFile(
					cluster.ClusterName,
					filepath.Join(dmgrutil.AbsClusterCacheDir(cluster.ClusterPath, cluster.ClusterName), "dashboard.yml"),
					filepath.Join(dmgrutil.AbsClusterDataboardDir(cluster.DeployDir, cluster.InstanceName), "dashboard.yml"),
					dmgrutil.FileTypeComponent,
					cluster.MachineHost,
					false,
					0).
				CopyFile(
					cluster.ClusterName,
					filepath.Join(dmgrutil.AbsClusterCacheDir(cluster.ClusterPath, cluster.ClusterName), "datasource.yml"),
					filepath.Join(dmgrutil.AbsClusterDatasourceDir(cluster.DeployDir, cluster.InstanceName), "datasource.yml"),
					dmgrutil.FileTypeComponent,
					cluster.MachineHost,
					false,
					0).
				CopyFile(
					cluster.ClusterName,
					filepath.Join(dmgrutil.AbsClusterCacheDir(cluster.ClusterPath, cluster.ClusterName), "run_grafana.sh"),
					filepath.Join(dmgrutil.AbsClusterScriptDir(cluster.DeployDir, cluster.InstanceName), "run_grafana.sh"),
					dmgrutil.FileTypeScript,
					cluster.MachineHost,
					false,
					0)
		}

		if componentName == dmgrutil.ComponentPrometheus {
			copyFileTask.CopyFile(
				cluster.ClusterName,
				filepath.Join(dmgrutil.AbsClusterCacheDir(cluster.ClusterPath, cluster.ClusterName), "prometheus.yml"),
				filepath.Join(dmgrutil.AbsClusterConfDir(cluster.DeployDir, cluster.InstanceName), "prometheus.yml"),
				dmgrutil.FileTypeComponent,
				cluster.MachineHost,
				false,
				0).
				CopyFile(
					cluster.ClusterName,
					filepath.Join(dmgrutil.AbsUntarConfDir(cluster.ClusterPath, cluster.ClusterName, cluster.ClusterVersion, dmgrutil.ConfDmWorkerRuleFile)),
					filepath.Join(dmgrutil.AbsClusterConfDir(cluster.DeployDir, cluster.InstanceName), "dm_worker.rules.yml"),
					dmgrutil.FileTypeRule,
					cluster.MachineHost,
					false,
					0).
				CopyFile(
					cluster.ClusterName,
					filepath.Join(dmgrutil.AbsClusterCacheDir(cluster.ClusterPath, cluster.ClusterName), "run_prometheus.sh"),
					filepath.Join(dmgrutil.AbsClusterScriptDir(cluster.DeployDir, cluster.InstanceName), "run_prometheus.sh"),
					dmgrutil.FileTypeScript,
					cluster.MachineHost,
					false,
					0)
		}

		if componentName == dmgrutil.ComponentAlertmanager {
			copyFileTask.CopyFile(
				cluster.ClusterName,
				filepath.Join(dmgrutil.AbsUntarConfDir(cluster.ClusterPath, cluster.ClusterName, cluster.ClusterVersion, dmgrutil.ConfAlertmanagerFile)),
				filepath.Join(dmgrutil.AbsClusterConfDir(cluster.DeployDir, cluster.InstanceName), "alertmanager.yml"),
				dmgrutil.FileTypeComponent,
				cluster.MachineHost,
				false,
				0).
				CopyFile(
					cluster.ClusterName,
					filepath.Join(dmgrutil.AbsClusterCacheDir(cluster.ClusterPath, cluster.ClusterName), fmt.Sprintf("run_alertmanager-%s-%d.sh", cluster.MachineHost, cluster.ServicePort)),
					filepath.Join(dmgrutil.AbsClusterScriptDir(cluster.DeployDir, cluster.InstanceName), "run_alertmanager.sh"),
					dmgrutil.FileTypeScript,
					cluster.MachineHost,
					false,
					0)
		}

		if componentName == dmgrutil.ComponentDmMaster {
			copyFileTask.CopyFile(
				cluster.ClusterName,
				filepath.Join(dmgrutil.AbsUntarConfDir(cluster.ClusterPath, cluster.ClusterName, cluster.ClusterVersion, dmgrutil.ConfDmMasterFile)),
				filepath.Join(dmgrutil.AbsClusterConfDir(cluster.DeployDir, cluster.InstanceName), "dm-master.toml"),
				dmgrutil.FileTypeComponent,
				cluster.MachineHost,
				false,
				0).
				CopyFile(
					cluster.ClusterName,
					filepath.Join(dmgrutil.AbsClusterCacheDir(cluster.ClusterPath, cluster.ClusterName), fmt.Sprintf("run_dm-master-%s-%d.sh", cluster.MachineHost, cluster.ServicePort)),
					filepath.Join(dmgrutil.AbsClusterScriptDir(cluster.DeployDir, cluster.InstanceName), "run_dm-master.sh"),
					dmgrutil.FileTypeScript,
					cluster.MachineHost,
					false,
					0)
		}

		if componentName == dmgrutil.ComponentDmWorker {
			copyFileTask.CopyFile(
				cluster.ClusterName,
				filepath.Join(dmgrutil.AbsUntarConfDir(cluster.ClusterPath, cluster.ClusterName, cluster.ClusterVersion, dmgrutil.ConfDmWorkerFile)),
				filepath.Join(dmgrutil.AbsClusterConfDir(cluster.DeployDir, cluster.InstanceName), "dm-worker.toml"),
				dmgrutil.FileTypeComponent,
				cluster.MachineHost,
				false,
				0).
				CopyFile(
					cluster.ClusterName,
					filepath.Join(dmgrutil.AbsClusterCacheDir(cluster.ClusterPath, cluster.ClusterName), fmt.Sprintf("run_dm-worker-%s-%d.sh", cluster.MachineHost, cluster.ServicePort)),
					filepath.Join(dmgrutil.AbsClusterScriptDir(cluster.DeployDir, cluster.InstanceName), "run_dm-worker.sh"),
					dmgrutil.FileTypeScript,
					cluster.MachineHost,
					false,
					0)
		}

		copyFileTask.CopyFile(
			cluster.ClusterName,
			filepath.Join(dmgrutil.AbsClusterCacheDir(cluster.ClusterPath, cluster.ClusterName), fmt.Sprintf("%s-%s-%d.service", componentName, cluster.MachineHost, cluster.ServicePort)),
			filepath.Join(dmgrutil.AbsClusterTempSystemdDir(cluster.DeployDir), fmt.Sprintf("%s-%d.service", componentName, cluster.ServicePort)),
			dmgrutil.FileTypeSystemd,
			cluster.MachineHost,
			false,
			0)
		copyFileTasks = append(copyFileTasks, copyFileTask.BuildTask())
	}

	return copyFileTasks
}

// 用于初始化环境的任务
func EnvClusterUserInit(machineList []response.MachineRespStruct, clusterUser, skipCreateUser string) []task.Task {
	var envInitTasks []task.Task

	// 集群环境初始化以及组件名 Copy
	for _, machine := range machineList {
		envInitTask := task.NewBuilder().
			RootSSH(
				machine.SshHost,
				machine.SshPort,
				machine.SshUser,
				machine.SshPassword,
				"",
				"",
				executor.DefaultConnectTimeout,
				executor.DefaultExecuteTimeout,
			).
			EnvInit(
				machine.SshHost,
				clusterUser,
				"",
				dmgrutil.StringEqualFold(skipCreateUser, dmgrutil.BoolTrue),
			).BuildTask()

		envInitTasks = append(envInitTasks, envInitTask)
	}

	return envInitTasks
}

// 用于 COPY 集群组件的任务
func EnvClusterComponentInit(clusterTopo []response.ClusterTopologyRespStruct,
	clusterUntarDir string) []task.Task {
	var copyCompTasks []task.Task

	// 集群环境初始化以及组件名 Copy
	for _, cluster := range clusterTopo {
		copyCompTask := task.NewBuilder().
			UserSSH(
				cluster.MachineHost,
				cluster.SshPort,
				cluster.ClusterUser,
				executor.DefaultConnectTimeout,
				executor.DefaultExecuteTimeout,
			).
			Mkdir(cluster.ClusterUser, cluster.MachineHost, []string{
				dmgrutil.AbsClusterBinDir(cluster.DeployDir, cluster.InstanceName),
				dmgrutil.AbsClusterConfDir(cluster.DeployDir, cluster.InstanceName),
				dmgrutil.AbsClusterScriptDir(cluster.DeployDir, cluster.InstanceName),
				dmgrutil.AbsClusterDataDir(cluster.DeployDir, cluster.DataDir, cluster.InstanceName),
				dmgrutil.AbsClusterLogDir(cluster.DeployDir, cluster.LogDir, cluster.InstanceName)}...)

		switch strings.ToLower(cluster.ComponentName) {
		case dmgrutil.ComponentGrafana:
			copyCompTask.CopyComponent(
				cluster.ClusterName,
				cluster.ComponentName,
				cluster.ClusterVersion,
				dmgrutil.AbsClusterGrafanaComponent(cluster.ClusterPath, cluster.ClusterName, cluster.ClusterVersion, dmgrutil.ComponentGrafanaTarPKG),
				cluster.MachineHost,
				fmt.Sprintf("%s/%s", dmgrutil.AbsClusterDeployDir(cluster.DeployDir, cluster.InstanceName), dmgrutil.ComponentGrafanaTarPKG))
		default:
			copyCompTask.CopyComponent(
				cluster.ClusterName,
				cluster.ComponentName,
				cluster.ClusterVersion,
				filepath.Join(clusterUntarDir, dmgrutil.DirBin, strings.ToLower(cluster.ComponentName)),
				cluster.MachineHost,
				filepath.Join(dmgrutil.AbsClusterBinDir(cluster.DeployDir, cluster.InstanceName), strings.ToLower(cluster.ComponentName)),
			)
		}
		copyCompTasks = append(copyCompTasks, copyCompTask.BuildTask())
	}

	return copyCompTasks
}

// 用于 DM Master API 访问
func GetActiveDmMasterAddr(s *service.MysqlService, clusterName string) (string, error) {
	var (
		dmMasterAddr     []string
		activeMasterAddr string
	)

	dmMasters, err := s.GetClusterComponent(clusterName, dmgrutil.ComponentDmMaster)
	if err != nil {
		return activeMasterAddr, err
	}
	for _, dm := range dmMasters {
		dmMasterAddr = append(dmMasterAddr, fmt.Sprintf("%s:%d", dm.MachineHost, dm.ServicePort))
	}

	dmMasterClient := api.NewDMMasterClient(dmMasterAddr, api.DmMasterApiTimeout, nil)

	_, activeMasterAddr, err = dmMasterClient.GetLeader(api.DefaultRetryOpt)
	if err != nil {
		return activeMasterAddr, err
	}

	return dmMasterClient.GetURL(activeMasterAddr), nil
}

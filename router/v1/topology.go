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
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wentaojin/dmgr/pkg/cluster/template/script"

	"github.com/wentaojin/dmgr/pkg/cluster/ctxt"

	"github.com/wentaojin/dmgr/pkg/cluster/api"

	"github.com/wentaojin/dmgr/pkg/cluster/module"

	"github.com/wentaojin/dmgr/pkg/cluster/template"

	"github.com/wentaojin/dmgr/pkg/cluster/executor"

	"github.com/wentaojin/dmgr/pkg/cluster/task"

	"github.com/wentaojin/dmgr/pkg/dmgrutil"

	"github.com/gin-gonic/gin"
	"github.com/wentaojin/dmgr/request"
	"github.com/wentaojin/dmgr/response"
	"github.com/wentaojin/dmgr/service"
)

// 集群部署
func ClusterDeploy(c *gin.Context) {
	var req request.ClusterDeployReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	// 验证判断前端请求
	topo, uniqueHosts, instanceList := request.ValidDeployReqStructField(req)

	// 判断集群名是否冲突
	s := service.NewMysqlService()
	exist, err := s.ValidClusterNameIsExist(topo.ClusterName)
	if response.FailWithMsg(c, err) {
		return
	}
	if exist {
		if response.FailWithMsg(c, fmt.Errorf("cluster deploy failed: cluster_name [%v] exist", topo.ClusterName)) {
			return
		}
	}

	// 判断是否存在主机端口冲突
	dbHostPortArr, err := s.GetMachinePortArray(uniqueHosts)
	if response.FailWithMsg(c, err) {
		return
	}
	if response.FailWithMsg(c, request.ValidComponentPortConflict(topo.ClusterTopology, dbHostPortArr)) {
		return
	}

	// 判断实例名是否冲突【集群全局唯一】
	repeatInstList := dmgrutil.FilterRepeatElem(instanceList)
	if len(repeatInstList) != 0 {
		if response.FailWithMsg(c, fmt.Errorf("cluster topology component isn't global unique, exist conflict [%v]", repeatInstList)) {
			return
		}
	}

	// 判断对应版本离线安装包是否存在
	pkg, err := s.ValidClusterVersionPackageIsExist(topo.ClusterVersion)
	if response.FailWithMsg(c, err) {
		return
	}
	if dmgrutil.IsStructureEqual(pkg, response.WarehouseRespStruct{}) {
		if response.FailWithMsg(c, fmt.Errorf("cluster_version [%v] offline package not exist", topo.ClusterVersion)) {
			return
		}
	}

	// 解压离线镜像包到指定目录
	// {cluster_path}/cluster/{cluster_name}/{cluster_version}
	clusterNameDir := dmgrutil.AbsClusterUntarDir(topo.ClusterPath, topo.ClusterName)
	clusterUntarDir := filepath.Join(clusterNameDir, topo.ClusterVersion)

	if response.FailWithMsg(c, dmgrutil.UnCompressTarGz(filepath.Join(pkg.PackagePath, pkg.PackageName), clusterUntarDir)) {
		return
	}

	// 初始化组件配置文件、脚本等文件缓存目录以及 SSH 认证存放目录
	if response.FailWithMsg(c, dmgrutil.InitComponentCacheAndSSHDir(topo.ClusterPath, topo.ClusterName)) {
		return
	}

	// 集群拓扑查询生成
	// 1、根据拓扑中主机信息，获取主机 SSH 信息
	// 2、判断是否存在主机信息，存在则重新生成集群拓扑
	machineList, err := s.GetMachineList(uniqueHosts)
	if response.FailWithMsg(c, err) {
		return
	}
	if len(machineList) == 0 || len(machineList) != len(uniqueHosts) {
		if response.FailWithMsg(c, fmt.Errorf("cluster topology host [%v] isn't match machine list, please add host information first", uniqueHosts)) {
			return
		}
	}
	clusterTopo, err := GenerateClusterTopology(req.ClusterMetaReqStruct, req.ClusterTopology, machineList)
	if response.FailWithMsg(c, err) {
		return
	}

	// 集群部署
	// 集群环境初始化以及集群组件复制 COPY
	envInitTasks := EnvClusterUserInit(machineList, topo.ClusterUser, topo.SkipCreateUser)
	copyCompTasks := EnvClusterComponentInit(clusterTopo, clusterUntarDir)

	// 获取生成集群部署配置文件、运行脚本等文件信息
	cos := template.GetClusterFile(clusterTopo)

	// 生成以及 Copy 组件配置文件、运行脚本
	if response.FailWithMsg(c,
		template.GenerateClusterFileWithStage(
			clusterTopo,
			cos,
			template.ClusterDeployStage,
			topo.AdminUser,
			topo.AdminPassword)) {
		return
	}
	copyFileTasks := CopyClusterFile(clusterTopo)

	builder := task.NewBuilder().
		Serial("+ Generate SSH keys",
			task.NewBuilder().
				SSHKeyGen(dmgrutil.HomeSshDir, executor.DefaultExecuteTimeout).
				SSHKeyCopy(dmgrutil.HomeSshDir, dmgrutil.AbsClusterSSHDir(topo.ClusterPath, topo.ClusterName), machineList, executor.DefaultExecuteTimeout, dmgrutil.RsaConcurrency).BuildTask()).
		Parallel("+ Initialize target host environments", false, envInitTasks...).
		Parallel("+ Copy components", false, copyCompTasks...).
		Parallel("+ Copy files", false, copyFileTasks...).BuildTask()

	if response.FailWithMsg(c, builder.Execute(ctxt.NewContext())) {
		return
	}

	// 集群元数据以及集群拓扑更新
	if response.FailWithMsg(c, s.AddClusterMetaAndTopology(topo.ClusterMetaReqStruct, topo.ClusterTopology)) {
		return
	}

	// TODO: 清理缓存目录

	response.SuccessWithoutData(c)
}

// 集群启动
func ClusterStart(c *gin.Context) {
	var req request.ClusterOperatorReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	// 判断指定实例名是否在指定组件中 [组件操作以实例名为准，实例名全局唯一]
	s := service.NewMysqlService()
	instNames, err := s.FilterComponentInstance(req)
	if response.FailWithMsg(c, err) {
		return
	}

	// 根据集群名、实例名查询集群拓扑
	clusterTopos, err := s.GetClusterTopologyByInstanceName(req.ClusterName, instNames)
	if response.FailWithMsg(c, err) {
		return
	}

	// 按组件启动顺序启动
	for _, component := range dmgrutil.StartComponentOrder {
		for _, t := range clusterTopos {
			compName := strings.ToLower(t.ComponentName)
			if component == compName {
				startCompTask := task.NewBuilder().
					SSHKeySet(
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519"),
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519.pub")).
					UserSSH(
						t.MachineHost,
						t.SshPort,
						t.ClusterUser,
						executor.DefaultConnectTimeout,
						module.DefaultSystemdExecuteTimeout).
					StartInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
						fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort),
						module.DefaultSystemdExecuteTimeout).
					EnableInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
						fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort),
						module.DefaultSystemdExecuteTimeout, true).BuildTask()

				if response.FailWithMsg(c, startCompTask.Execute(ctxt.NewContext())) {
					return
				}
			}
		}
	}

	// 更新集群状态
	if response.FailWithMsg(c, s.UpdateClusterMetaStatus(req.ClusterName, dmgrutil.ClusterUpStatus)) {
		return
	}
	response.SuccessWithoutData(c)
}

// 集群停止
func ClusterStop(c *gin.Context) {
	var req request.ClusterOperatorReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	// 判断指定实例名是否在指定组件中 [组件操作以实例名为准，实例名全局唯一]
	s := service.NewMysqlService()
	instNames, err := s.FilterComponentInstance(req)
	if response.FailWithMsg(c, err) {
		return
	}

	// 根据集群名、实例名查询集群拓扑
	clusterTopos, err := s.GetClusterTopologyByInstanceName(req.ClusterName, instNames)
	if response.FailWithMsg(c, err) {
		return
	}

	// 按组件停止顺序停止
	for _, component := range dmgrutil.StopComponentOrder {
		for _, t := range clusterTopos {
			compName := strings.ToLower(t.ComponentName)
			if component == compName {
				startCompTask := task.NewBuilder().
					SSHKeySet(
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519"),
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519.pub")).
					UserSSH(
						t.MachineHost,
						t.SshPort,
						t.ClusterUser,
						executor.DefaultConnectTimeout,
						module.DefaultSystemdExecuteTimeout).
					StopInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
						fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort),
						module.DefaultSystemdExecuteTimeout).BuildTask()

				if response.FailWithMsg(c, startCompTask.Execute(ctxt.NewContext())) {
					return
				}
			}
		}
	}
	// 更新集群状态
	if response.FailWithMsg(c, s.UpdateClusterMetaStatus(req.ClusterName, dmgrutil.ClusterOfflineStatus)) {
		return
	}
	response.SuccessWithoutData(c)
}

// 集群扩容
func ClusterScaleOut(c *gin.Context) {
	var req request.CLusterScaleOutReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	// 验证判断前端请求
	topo := request.ValidScaleOutReqStructField(req)

	// 判断扩容集群是否存在
	s := service.NewMysqlService()
	exist, err := s.ValidClusterNameIsExist(topo.ClusterName)
	if response.FailWithMsg(c, err) {
		return
	}
	if !exist {
		if response.FailWithMsg(c, fmt.Errorf("scale out failed: cluster [%v] not exist", topo.ClusterName)) {
			return
		}
	}

	// 判断是否存在主机端口冲突
	dbHostPortArr, err := s.GetMachinePortArray([]string{req.MachineHost})
	if response.FailWithMsg(c, err) {
		return
	}
	if response.FailWithMsg(c, request.ValidComponentPortConflict([]request.TopologyReqStruct{req.TopologyReqStruct}, dbHostPortArr)) {
		return
	}

	// 判断是否实例名冲突【集群全局唯一】
	// 获取集群拓扑信息
	topoDB, err := s.GetClusterTopologyByClusterName(topo.ClusterName)
	if response.FailWithMsg(c, err) {
		return
	}
	for _, t := range topoDB {
		if topo.InstanceName == strings.ToLower(t.InstanceName) {
			if response.FailWithMsg(c, fmt.Errorf("scale out component [%s] instance [%s] conflict with the cluster exist component [%s] instance [%s]", topo.ComponentName, topo.InstanceName, t.ComponentName, t.InstanceName)) {
				return
			}
		}
	}

	// 获取扩容集群元信息
	clusterMeta, err := s.GetClusterMeta(req.ClusterName)
	if response.FailWithMsg(c, err) {
		return
	}

	// 判断解压离线镜像包指定目录是否存在扩容组件
	// {cluster_path}/cluster/{cluster_name}/{cluster_version}
	clusterNameDir := dmgrutil.AbsClusterUntarDir(clusterMeta.ClusterPath, clusterMeta.ClusterName)
	clusterUntarDir := filepath.Join(clusterNameDir, clusterMeta.ClusterVersion)

	if exist, _ := dmgrutil.PathExists(filepath.Join(clusterUntarDir, dmgrutil.DirBin, strings.ToLower(topo.ComponentName))); !exist {
		if response.FailWithMsg(c, fmt.Errorf("scale out failed: component [%v] bin not exist", strings.ToLower(topo.ComponentName))) {
			return
		}
	}

	// 集群拓扑查询生成
	// 1、根据拓扑中主机信息，获取主机 SSH 信息
	// 2、判断是否存在主机信息，存在则重新生成集群拓扑
	machineList, err := s.GetMachineList([]string{req.MachineHost})
	if response.FailWithMsg(c, err) {
		return
	}
	if len(machineList) == 0 || len(machineList) != 1 {
		if response.FailWithMsg(c, fmt.Errorf("cluster topology host [%v] isn't match machine list, please add host information first", req.MachineHost)) {
			return
		}
	}
	clusterTopo, err := GenerateClusterTopology(clusterMeta, []request.TopologyReqStruct{req.TopologyReqStruct}, machineList)
	if response.FailWithMsg(c, err) {
		return
	}

	// 集群环境初始化以及集群组件复制 COPY
	envInitTasks := EnvClusterUserInit(machineList, clusterMeta.ClusterUser, topo.SkipCreateUser)
	copyCompTasks := EnvClusterComponentInit(clusterTopo, clusterUntarDir)

	// 获取生成集群部署配置文件、运行脚本等文件信息【根据元数据库已有集群组件信息】
	cos := template.GetClusterFile(topoDB)

	// 生成以及 Copy 组件配置文件、运行脚本
	// 1、重新生成以及分发集群已存在的组件配置文件、运行脚本
	// 2、生成以及分发扩容组件配置文件、运行脚本
	switch strings.ToLower(topo.ComponentName) {
	case dmgrutil.ComponentDmMaster:
		cos.DmMasterAddrs = append(cos.DmMasterAddrs, fmt.Sprintf("%s:%v", topo.MachineHost, topo.ServicePort))
	case dmgrutil.ComponentDmWorker:
		cos.DmWorkerAddrs = append(cos.DmWorkerAddrs, fmt.Sprintf("%s:%v", topo.MachineHost, topo.ServicePort))
	case dmgrutil.ComponentGrafana:
		cos.GrafanaAddr = fmt.Sprintf("%s:%v", topo.MachineHost, topo.ServicePort)
	case dmgrutil.ComponentAlertmanager:
		cos.AlertmanagerAddrs = append(cos.AlertmanagerAddrs, fmt.Sprintf("%s:%v", topo.MachineHost, topo.ServicePort))
		cos.AlertmanagerScripts = append(cos.AlertmanagerScripts, &script.AlertManagerScript{
			IP:          topo.MachineHost,
			WebPort:     topo.ServicePort,
			ClusterPort: topo.ClusterPort,
			DeployDir:   dmgrutil.AbsClusterDeployDir(topo.DeployDir, topo.InstanceName),
			DataDir:     dmgrutil.AbsClusterDataDir(topo.DeployDir, topo.DataDir, topo.InstanceName),
			LogDir:      dmgrutil.AbsClusterLogDir(topo.DeployDir, topo.LogDir, topo.InstanceName),
			TLSEnabled:  false,
		})
	case dmgrutil.ComponentPrometheus:
		cos.PrometheusAddr = fmt.Sprintf("%s:%v", topo.MachineHost, topo.ServicePort)
	default:
		if response.FailWithMsg(c, fmt.Errorf("component [%v] not exist, panic", topo.ComponentName)) {
			return
		}
	}

	if response.FailWithMsg(c, template.GenerateClusterFileWithStage(topoDB,
		cos,
		template.ClusterDeployStage,
		"",
		"")) {
		return
	}

	if response.FailWithMsg(c, template.GenerateClusterFileWithStage(clusterTopo,
		cos,
		template.ClusterScaleOutStage,
		topo.AdminUser,
		topo.AdminPassword)) {
		return
	}

	// 集群已部署存在的组件配置文件以及脚本刷新 refresh
	refreshFileTasks := CopyClusterFile(topoDB)

	// 扩容组件配置文件以及脚本分发
	scaleFileTasks := CopyClusterFile(clusterTopo)

	// 扩容集群组件
	builder := task.NewBuilder().
		Serial("+ Generate SSH keys",
			task.NewBuilder().
				SSHKeyGen(dmgrutil.HomeSshDir, executor.DefaultExecuteTimeout).
				SSHKeyCopy(dmgrutil.HomeSshDir, dmgrutil.AbsClusterSSHDir(clusterTopo[0].ClusterPath, topo.ClusterName), machineList, executor.DefaultExecuteTimeout, dmgrutil.RsaConcurrency).BuildTask()).
		Parallel("+ Initialize target host environments", false, envInitTasks...).
		Parallel("+ Copy components", false, copyCompTasks...).
		Parallel("+ Refresh files", false, refreshFileTasks...).
		Parallel("+ Copy files", false, scaleFileTasks...).BuildTask()

	if response.FailWithMsg(c, builder.Execute(ctxt.NewContext())) {
		return
	}

	// 启动扩容组件
	for _, component := range dmgrutil.StartComponentOrder {
		for _, t := range clusterTopo {
			if component == strings.ToLower(t.ComponentName) {
				scaleOutCompTask := task.NewBuilder().
					SSHKeySet(
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519"),
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519.pub")).
					UserSSH(
						t.MachineHost,
						t.SshPort,
						t.ClusterUser,
						executor.DefaultConnectTimeout,
						module.DefaultSystemdExecuteTimeout).
					StartInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
						fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort),
						module.DefaultSystemdExecuteTimeout).
					EnableInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
						fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort),
						module.DefaultSystemdExecuteTimeout, true).BuildTask()

				if response.FailWithMsg(c, scaleOutCompTask.Execute(ctxt.NewContext())) {
					return
				}

				// 更新 grafana 用户密码
				// 如果存在其他多个 grafana，表记录只记录最后一个 grafana 用户密码，并且缩容 grafana 不会自动更新清理 admin_user、admin_password 字段记录
				switch component {
				case dmgrutil.ComponentGrafana:
					if response.FailWithMsg(c, s.UpdateGrafanaUserAndPassword(t.ClusterName, req.AdminUser, req.AdminPassword)) {
						return
					}
				}
			}
		}
	}

	// 刷新 prometheus 组件
	for _, t := range topoDB {
		if strings.ToLower(t.ComponentName) == dmgrutil.ComponentPrometheus {
			refreshCompTask := task.NewBuilder().
				SSHKeySet(
					filepath.Join(
						dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519"),
					filepath.Join(
						dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519.pub")).
				UserSSH(
					t.MachineHost,
					t.SshPort,
					t.ClusterUser,
					executor.DefaultConnectTimeout,
					module.DefaultSystemdExecuteTimeout).
				StopInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
					fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort),
					module.DefaultSystemdExecuteTimeout).
				StartInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
					fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort), module.DefaultSystemdExecuteTimeout).BuildTask()

			if response.FailWithMsg(c, refreshCompTask.Execute(ctxt.NewContext())) {
				return
			}
		}
	}

	// 更新集群拓扑
	if response.FailWithMsg(c, s.AddClusterTopology([]request.TopologyReqStruct{req.TopologyReqStruct})) {
		return
	}

	response.SuccessWithoutData(c)
}

// 集群缩容
func ClusterScaleIn(c *gin.Context) {
	var req request.ClusterOperatorReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	// 判断指定实例名是否在指定组件中 [组件操作以实例名为准，实例名全局唯一]
	s := service.NewMysqlService()
	instNames, err := s.FilterComponentInstance(req)
	if response.FailWithMsg(c, err) {
		return
	}

	// 获取缩容节点信息
	var delDmMasterAddrs []string
	delInsts, err := s.GetClusterComponentByInstance(req.ClusterName, instNames)
	if response.FailWithMsg(c, err) {
		return
	}
	for _, delInst := range delInsts {
		if delInst.ComponentName == dmgrutil.ComponentDmMaster {
			delDmMasterAddrs = append(delDmMasterAddrs, fmt.Sprintf("%s:%d", delInst.MachineHost, delInst.ServicePort))
		}
	}

	// 获取集群 DM Master 组件信息
	var dmMasterAddrs []string
	dmMasters, err := s.GetClusterComponent(req.ClusterName, dmgrutil.ComponentDmMaster)
	if response.FailWithMsg(c, err) {
		return
	}

	for _, dm := range dmMasters {
		dmMasterAddrs = append(dmMasterAddrs, fmt.Sprintf("%s:%d", dm.MachineHost, dm.ServicePort))
	}

	originDmMasters := dmgrutil.NewStringSet(dmMasterAddrs...)
	reqDelDmMasters := dmgrutil.NewStringSet(delDmMasterAddrs...)

	activeDmMasters := originDmMasters.Difference(reqDelDmMasters).Slice()
	if len(activeDmMasters) == 0 {
		if response.FailWithMsg(c, errors.New("cannot delete all dm-master servers")) {
			return
		}
	}
	dmMasterClient := api.NewDMMasterClient(activeDmMasters, api.DmMasterApiTimeout, nil)

	// 根据集群名、实例名查询集群拓扑
	clusterTopos, err := s.GetClusterTopologyByInstanceName(req.ClusterName, instNames)
	if response.FailWithMsg(c, err) {
		return
	}

	// 缩容组件
	// 注意：缩容组件 DestroyInstance 只会清理子目录，不会清理父目录
	// 比如：deployDir=/data/marvin/{instance_name}, 则清理执行命令 m -rf /data/marvin/{instance_name}，保留 /data/marvin/ 目录，防止误删除
	for _, component := range dmgrutil.StartComponentOrder {
		for _, t := range clusterTopos {
			if component == strings.ToLower(t.ComponentName) {
				scaleInCompTask := task.NewBuilder().
					SSHKeySet(
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519"),
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519.pub")).
					UserSSH(
						t.MachineHost,
						t.SshPort,
						t.ClusterUser,
						executor.DefaultConnectTimeout,
						module.DefaultSystemdExecuteTimeout).
					StopInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
						fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort),
						module.DefaultSystemdExecuteTimeout).
					DestroyInstance(t.MachineHost, t.ServicePort, t.ComponentName, t.InstanceName, t.DeployDir, t.DataDir, t.LogDir, executor.DefaultExecuteTimeout).BuildTask()

				if response.FailWithMsg(c, scaleInCompTask.Execute(ctxt.NewContext())) {
					return
				}

				switch component {
				case dmgrutil.ComponentDmMaster:
					err := dmMasterClient.OfflineMaster(t.InstanceName, nil)
					if response.FailWithMsg(c, err) {
						return
					}
				case dmgrutil.ComponentDmWorker:
					err := dmMasterClient.OfflineWorker(t.InstanceName, nil)
					if response.FailWithMsg(c, err) {
						return
					}
				}
			}
		}
	}

	// 清理元数据表
	if response.FailWithMsg(c, s.DelClusterTopologyByInstanceName(req.ClusterName, instNames)) {
		return
	}

	// 刷新 prometheus 组件
	if !dmgrutil.IsContainElem(req.ComponentName, dmgrutil.ComponentPrometheus) {
		// 获取集群拓扑信息
		topoDB, err := s.GetClusterTopologyByClusterName(req.ClusterName)
		if response.FailWithMsg(c, err) {
			return
		}
		// 获取扩容集群元信息
		clusterMeta, err := s.GetClusterMeta(req.ClusterName)
		if response.FailWithMsg(c, err) {
			return
		}

		// 获取生成集群部署配置文件、运行脚本等文件信息【根据元数据库已有集群组件信息】
		cos := template.GetClusterFile(topoDB)

		// 生成以及 Copy 组件配置文件、运行脚本
		if response.FailWithMsg(c, template.GenerateClusterFileWithStage(topoDB,
			cos,
			template.ClusterScaleOutStage,
			"",
			"")) {
			return
		}

		// 集群已部署存在的组件配置文件以及脚本刷新 refresh
		copyFileTasks := CopyClusterFile(topoDB)

		builder := task.NewBuilder().
			SSHKeySet(
				filepath.Join(
					dmgrutil.AbsClusterSSHDir(clusterMeta.ClusterPath, clusterMeta.ClusterName), "id_ed25519"),
				filepath.Join(
					dmgrutil.AbsClusterSSHDir(clusterMeta.ClusterPath, clusterMeta.ClusterName), "id_ed25519.pub")).
			Parallel("+ Copy files", false, copyFileTasks...).BuildTask()

		if response.FailWithMsg(c, builder.Execute(ctxt.NewContext())) {
			return
		}

		// 刷新 prometheus 组件
		for _, t := range topoDB {
			if strings.ToLower(t.ComponentName) == dmgrutil.ComponentPrometheus {
				refreshCompTask := task.NewBuilder().
					SSHKeySet(
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519"),
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519.pub")).
					UserSSH(
						t.MachineHost,
						t.SshPort,
						t.ClusterUser,
						executor.DefaultConnectTimeout,
						module.DefaultSystemdExecuteTimeout).
					StopInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
						fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort),
						module.DefaultSystemdExecuteTimeout).
					StartInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
						fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort), module.DefaultSystemdExecuteTimeout).BuildTask()

				if response.FailWithMsg(c, refreshCompTask.Execute(ctxt.NewContext())) {
					return
				}
			}
		}

	}
	response.SuccessWithoutData(c)
}

// 集群滚更
// Only 支持
// dm-master 组件配置文件 dm-master.toml
// dm-worker 组件配置文件 dm-worker.toml
// alertmanager 组件配置文件 alertmanager.yml 滚更
func CLusterReload(c *gin.Context) {
	var req request.ClusterPatchReqStruct
	if response.FailWithMsg(c, c.ShouldBind(&req)) {
		return
	}

	// 确保文件未缓存（例如在 iOS 设备上发生的情况）
	c.Header("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Cache-Control", "post-check=0, pre-check=0")
	c.Header("Pragma", "no-cache")

	s := service.NewMysqlService()
	// 获取集群元信息
	clusterMeta, err := s.GetClusterMeta(req.ClusterName)
	if response.FailWithMsg(c, err) {
		return
	}

	// 解压离线镜像包到指定目录
	// {cluster_path}/cluster/{cluster_name}/{cluster_version}
	clusterNameDir := dmgrutil.AbsClusterUntarDir(clusterMeta.ClusterPath, clusterMeta.ClusterName)
	// 新集群版本路径
	clusterUntarDir := filepath.Join(clusterNameDir, clusterMeta.ClusterVersion)
	// 目标目录是否存在
	pkgDir := filepath.Join(clusterUntarDir, dmgrutil.DirPatch)
	if exist, _ := dmgrutil.PathExists(pkgDir); !exist {
		if response.FailWithMsg(c, os.MkdirAll(pkgDir, 0750)) {
			return
		}
	}

	// 获取前端数组
	reqInstances := c.PostFormArray("instance_name")

	// 定义需要前端上传的文件字段名
	file, err := c.FormFile("file")
	if response.FailWithMsg(c, err) {
		return
	}

	// 判断是否符合 reload 组件要求
	switch strings.ToLower(req.ComponentName) {
	case dmgrutil.ComponentDmMaster:
		if file.Filename != dmgrutil.ReloadDmMasterFile {
			if response.FailWithMsg(c, fmt.Errorf("component [%v] reload file name is non-compliant, only support file name is dm-master.toml", dmgrutil.ComponentDmMaster)) {
				return
			}
		}
	case dmgrutil.ComponentDmWorker:
		if file.Filename != dmgrutil.ReloadDmWorkerFile {
			if response.FailWithMsg(c, fmt.Errorf("component [%v] reload file name is non-compliant, only support file name is dm-worker.toml", dmgrutil.ComponentDmMaster)) {
				return
			}
		}
	case dmgrutil.ComponentAlertmanager:
		if file.Filename != dmgrutil.ReloadAlertmanagerFile {
			if response.FailWithMsg(c, fmt.Errorf("component [%v] reload file name is non-compliant, only support file name is alertmanager.yml", dmgrutil.ComponentDmMaster)) {
				return
			}
		}
	default:
		if response.FailWithMsg(c, fmt.Errorf("only support dm-master、dm-worker and alertmanager component reload, others not support")) {
			return
		}
	}

	// 文件上传
	filePath := filepath.Join(pkgDir, file.Filename)
	if response.FailWithMsg(c, c.SaveUploadedFile(file, filePath)) {
		return
	}

	// 	本地运行
	// 文件解压是否覆盖
	var cmd string
	if req.Overwrite == dmgrutil.BoolTrue {
		// 覆盖
		cmd = fmt.Sprintf(`cp %s %s`, filepath.Join(pkgDir, file.Filename), filepath.Join(clusterUntarDir, dmgrutil.DirConf, file.Filename))
	}
	currentUser, currentIP, err := dmgrutil.GetClientOutBoundIP()
	if response.FailWithMsg(c, err) {
		return
	}
	_, stdErr, err := executor.NewLocalExecutor(currentIP, currentUser, currentUser == "root").Execute(cmd, executor.DefaultExecuteTimeout)
	if response.FailWithMsg(c, err) {
		return
	}
	if len(stdErr) != 0 {
		if response.FailWithMsg(c, fmt.Errorf("local host [%v] user [%v] running cmd [%v] failed: %v", currentIP, currentUser, cmd, string(stdErr))) {
			return
		}
	}

	// 判断指定实例名是否在指定组件中 [组件操作以实例名为准，实例名全局唯一]
	instNames, err := s.FilterComponentInstance(request.ClusterOperatorReqStruct{
		ClusterName:   req.ClusterName,
		ComponentName: []string{req.ComponentName},
		InstanceName:  reqInstances,
	})
	if response.FailWithMsg(c, err) {
		return
	}

	// 根据集群名、实例名查询集群拓扑
	clusterTopos, err := s.GetClusterTopologyByInstanceName(req.ClusterName, instNames)
	if response.FailWithMsg(c, err) {
		return
	}

	// 启停对应组件
	for _, component := range dmgrutil.StartComponentOrder {
		for _, t := range clusterTopos {
			if component == strings.ToLower(t.ComponentName) {
				reloadCompTask := task.NewBuilder().
					SSHKeySet(
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519"),
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519.pub")).
					UserSSH(
						t.MachineHost,
						t.SshPort,
						t.ClusterUser,
						executor.DefaultConnectTimeout,
						module.DefaultSystemdExecuteTimeout).
					CopyFile(
						t.ClusterName,
						filepath.Join(pkgDir, file.Filename),
						filepath.Join(dmgrutil.AbsClusterConfDir(t.DeployDir, t.InstanceName), file.Filename),
						dmgrutil.FileTypeComponent,
						t.MachineHost,
						false,
						0,
					).
					StopInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
						fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort),
						module.DefaultSystemdExecuteTimeout).
					StartInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
						fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort), module.DefaultSystemdExecuteTimeout).BuildTask()
				if response.FailWithMsg(c, reloadCompTask.Execute(ctxt.NewContext())) {
					return
				}
				// 元数据表更新
				if response.FailWithMsg(c, s.UpdateClusterHotFixStatus(req.ClusterName, t.InstanceName, dmgrutil.ReloadComponent)) {
					return
				}
			}
		}
	}

	response.SuccessWithoutData(c)
}

// 集群升级
func ClusterUpgrade(c *gin.Context) {
	var req request.ClusterUpgradeReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	// 判断对应版本离线安装包是否存在
	s := service.NewMysqlService()
	pkg, err := s.ValidClusterVersionPackageIsExist(req.ClusterVersion)
	if response.FailWithMsg(c, err) {
		return
	}
	if dmgrutil.IsStructureEqual(pkg, response.WarehouseRespStruct{}) {
		if response.FailWithMsg(c, fmt.Errorf("cluster_version [%v] offline package not exist", req.ClusterVersion)) {
			return
		}
	}

	// 获取集群元信息
	clusterMeta, err := s.GetClusterMeta(req.ClusterName)
	if response.FailWithMsg(c, err) {
		return
	}
	// 解压离线镜像包到指定目录
	// {cluster_path}/cluster/{cluster_name}/{cluster_version}
	clusterNameDir := dmgrutil.AbsClusterUntarDir(clusterMeta.ClusterPath, clusterMeta.ClusterName)

	// 创建新集群版本路径
	clusterUntarDir := filepath.Join(clusterNameDir, req.ClusterVersion)
	if response.FailWithMsg(c, dmgrutil.UnCompressTarGz(filepath.Join(pkg.PackagePath, pkg.PackageName), clusterUntarDir)) {
		return
	}

	// 集群拓扑查询生成 From 数据库
	clusterTopoDB, err := s.GetClusterTopologyByClusterName(clusterMeta.ClusterName)
	if response.FailWithMsg(c, err) {
		return
	}

	// 考虑升级，需要把数据库的 cluster_version 替换成当前升级集群版本
	var clusterTopos []response.ClusterTopologyRespStruct
	for _, c := range clusterTopoDB {
		c.ClusterVersion = req.ClusterVersion
		clusterTopos = append(clusterTopos, c)
	}

	// 获取生成集群部署配置文件、运行脚本等文件信息
	cos := template.GetClusterFile(clusterTopos)

	// 继承上个版本参数配置文件 dm-master.toml、dm-worker.toml 以及 alertmanager.yml，其他保持默认默认，不影响
	cmds := []string{
		fmt.Sprintf(`cp %v %v`,
			filepath.Join(clusterNameDir, clusterMeta.ClusterVersion, dmgrutil.DirConf, "*.toml"),
			filepath.Join(clusterNameDir, req.ClusterVersion, dmgrutil.DirConf)),
		fmt.Sprintf(`cp %v %v`,
			filepath.Join(clusterNameDir, clusterMeta.ClusterVersion, dmgrutil.DirConf, "alertmanager.yml"),
			filepath.Join(clusterNameDir, req.ClusterVersion, dmgrutil.DirConf)),
	}
	currentUser, currentIP, err := dmgrutil.GetClientOutBoundIP()
	if response.FailWithMsg(c, err) {
		return
	}
	for _, cmd := range cmds {
		_, stdErr, err := executor.NewLocalExecutor(currentIP, currentUser, currentUser == "root").Execute(cmd, executor.DefaultExecuteTimeout)
		if response.FailWithMsg(c, err) {
			return
		}
		if len(stdErr) != 0 {
			if response.FailWithMsg(c, fmt.Errorf("local host [%v] user [%v] running cmd [%v] failed: %v", currentIP, currentUser, cmd, string(stdErr))) {
				return
			}
		}
	}

	// 生成以及 Copy 组件配置文件、运行脚本
	if response.FailWithMsg(c,
		template.GenerateClusterFileWithStage(
			clusterTopos,
			cos,
			template.ClusterDeployStage,
			clusterMeta.AdminUser,
			clusterMeta.AdminPassword)) {
		return
	}
	copyFileTasks := CopyClusterFile(clusterTopos)
	copyFileTask := task.NewBuilder().
		SSHKeySet(
			filepath.Join(
				dmgrutil.AbsClusterSSHDir(clusterMeta.ClusterPath, clusterMeta.ClusterName), "id_ed25519"),
			filepath.Join(
				dmgrutil.AbsClusterSSHDir(clusterMeta.ClusterPath, clusterMeta.ClusterName), "id_ed25519.pub")).
		Parallel("+ Copy files", false, copyFileTasks...).BuildTask()

	if response.FailWithMsg(c, copyFileTask.Execute(ctxt.NewContext())) {
		return
	}

	// 升级对应组件
	for _, component := range dmgrutil.StartComponentOrder {
		for _, t := range clusterTopos {
			if component == strings.ToLower(t.ComponentName) {
				upgradeCompTask := task.NewBuilder().
					SSHKeySet(
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519"),
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519.pub")).
					UserSSH(
						t.MachineHost,
						t.SshPort,
						t.ClusterUser,
						executor.DefaultConnectTimeout,
						module.DefaultSystemdExecuteTimeout).
					StopInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
						fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort),
						module.DefaultSystemdExecuteTimeout)

				switch strings.ToLower(t.ComponentName) {
				case dmgrutil.ComponentGrafana:
					upgradeCompTask = upgradeCompTask.CopyComponent(
						t.ClusterName,
						t.ComponentName,
						req.ClusterVersion,
						dmgrutil.AbsClusterGrafanaComponent(t.ClusterPath, t.ClusterName, req.ClusterVersion, dmgrutil.ComponentGrafanaTarPKG),
						t.MachineHost,
						fmt.Sprintf("%s/%s", dmgrutil.AbsClusterDeployDir(t.DeployDir, t.InstanceName), dmgrutil.ComponentGrafanaTarPKG),
					)
				default:
					upgradeCompTask = upgradeCompTask.CopyComponent(
						t.ClusterName,
						t.ComponentName,
						req.ClusterVersion,
						filepath.Join(clusterUntarDir, dmgrutil.DirBin, strings.ToLower(t.ComponentName)),
						t.MachineHost,
						dmgrutil.AbsClusterBinDir(t.DeployDir, t.InstanceName),
					)
				}
				upgradeCompTask = upgradeCompTask.StartInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
					fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort), module.DefaultSystemdExecuteTimeout)

				if response.FailWithMsg(c, upgradeCompTask.BuildTask().Execute(ctxt.NewContext())) {
					return
				}
			}
		}
	}

	// 更新元数据集群版本信息
	if response.FailWithMsg(c, s.UpdateClusterVersion(req.ClusterName, req.ClusterVersion)) {
		return
	}
	response.SuccessWithoutData(c)
}

// 集群补丁
func ClusterPatch(c *gin.Context) {
	var req request.ClusterPatchReqStruct
	if response.FailWithMsg(c, c.ShouldBind(&req)) {
		return
	}

	// 确保文件未缓存（例如在 iOS 设备上发生的情况）
	c.Header("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Cache-Control", "post-check=0, pre-check=0")
	c.Header("Pragma", "no-cache")

	s := service.NewMysqlService()
	// 获取集群元信息
	clusterMeta, err := s.GetClusterMeta(req.ClusterName)
	if response.FailWithMsg(c, err) {
		return
	}

	// 解压离线镜像包到指定目录
	// {cluster_path}/cluster/{cluster_name}/{cluster_version}
	clusterNameDir := dmgrutil.AbsClusterUntarDir(clusterMeta.ClusterPath, clusterMeta.ClusterName)
	// 新集群版本路径
	clusterUntarDir := filepath.Join(clusterNameDir, clusterMeta.ClusterVersion)
	// 目标目录是否存在
	pkgDir := filepath.Join(clusterUntarDir, dmgrutil.DirPatch)
	if exist, _ := dmgrutil.PathExists(pkgDir); !exist {
		if response.FailWithMsg(c, os.MkdirAll(pkgDir, 0750)) {
			return
		}
	}

	// 获取前端数组
	reqInstances := c.PostFormArray("instance_name")

	// 判断指定实例名是否在指定组件中 [组件操作以实例名为准，实例名全局唯一]
	instNames, err := s.FilterComponentInstance(request.ClusterOperatorReqStruct{
		ClusterName:   req.ClusterName,
		ComponentName: []string{req.ComponentName},
		InstanceName:  reqInstances,
	})
	if response.FailWithMsg(c, err) {
		return
	}

	// 定义需要前端上传的文件字段名
	file, err := c.FormFile("file")
	if response.FailWithMsg(c, err) {
		return
	}

	// 文件上传
	filePath := filepath.Join(pkgDir, fmt.Sprintf("%v.tar.gz", req.ComponentName))
	if response.FailWithMsg(c, c.SaveUploadedFile(file, filePath)) {
		return
	}

	// 	本地运行
	// 文件解压是否覆盖
	var cmds []string
	if req.Overwrite == dmgrutil.BoolTrue {
		// 覆盖
		// 组件是 grafana 组件则不进行解压
		if strings.ToLower(req.ComponentName) == dmgrutil.ComponentGrafana {
			cmds = []string{
				fmt.Sprintf(`cp %s %s`, filepath.Join(pkgDir, dmgrutil.ComponentGrafanaTarPKG), filepath.Join(clusterUntarDir, dmgrutil.ComponentGrafanaTarPKG)),
			}
		} else {
			cmds = []string{
				fmt.Sprintf(`tar --no-same-owner -zxvf %v -C %v; rm -rf %v`, filePath, pkgDir, filePath),
				fmt.Sprintf(`cp %s %s`, filepath.Join(pkgDir, strings.ToLower(req.ComponentName)), filepath.Join(clusterUntarDir, strings.ToLower(req.ComponentName))),
			}
		}
	} else {
		// 补丁组件非 grafana 组件需要解压
		if strings.ToLower(req.ComponentName) != dmgrutil.ComponentGrafana {
			cmds = []string{fmt.Sprintf(`tar --no-same-owner -zxvf %v -C %v; rm -rf %v`, filePath, pkgDir, filePath)}
		}
	}

	for _, cmd := range cmds {
		currentUser, currentIP, err := dmgrutil.GetClientOutBoundIP()
		if response.FailWithMsg(c, err) {
			return
		}
		_, stdErr, err := executor.NewLocalExecutor(currentIP, currentUser, currentUser == "root").Execute(cmd, executor.DefaultExecuteTimeout)
		if response.FailWithMsg(c, err) {
			return
		}
		if len(stdErr) != 0 {
			if response.FailWithMsg(c, fmt.Errorf("local host [%v] user [%v] running cmd [%v] failed: %v", currentIP, currentUser, cmd, string(stdErr))) {
				return
			}
		}
	}

	// 根据集群名、实例名查询集群拓扑
	clusterTopos, err := s.GetClusterTopologyByInstanceName(req.ClusterName, instNames)
	if response.FailWithMsg(c, err) {
		return
	}

	// 集群组件补丁
	for _, component := range dmgrutil.StartComponentOrder {
		for _, t := range clusterTopos {
			if component == strings.ToLower(t.ComponentName) {
				patchCompTask := task.NewBuilder().
					SSHKeySet(
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519"),
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519.pub")).
					UserSSH(
						t.MachineHost,
						t.SshPort,
						t.ClusterUser,
						executor.DefaultConnectTimeout,
						module.DefaultSystemdExecuteTimeout).
					StopInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
						fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort),
						module.DefaultSystemdExecuteTimeout)

				switch strings.ToLower(t.ComponentName) {
				case dmgrutil.ComponentGrafana:
					patchCompTask = patchCompTask.CopyComponent(
						t.ClusterName,
						t.ComponentName,
						"patched",
						filepath.Join(pkgDir, dmgrutil.ComponentGrafanaTarPKG),
						t.MachineHost,
						dmgrutil.AbsClusterBinDir(t.DeployDir, t.InstanceName),
					)
				default:
					patchCompTask = patchCompTask.CopyComponent(
						t.ClusterName,
						t.ComponentName,
						"patched",
						filepath.Join(pkgDir, t.ComponentName),
						t.MachineHost,
						dmgrutil.AbsClusterBinDir(t.DeployDir, t.InstanceName),
					)
				}
				patchCompTask = patchCompTask.StartInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
					fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort), module.DefaultSystemdExecuteTimeout)

				if response.FailWithMsg(c, patchCompTask.BuildTask().Execute(ctxt.NewContext())) {
					return
				}
				// 元数据表更新
				if response.FailWithMsg(c, s.UpdateClusterHotFixStatus(req.ClusterName, t.InstanceName, dmgrutil.PatchedComponent)) {
					return
				}
			}
		}
	}
	response.SuccessWithoutData(c)
}

// 集群状态查询
func ClusterStatus(c *gin.Context) {
	var req request.ClusterStatusReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}
	// 创建服务
	s := service.NewMysqlService()
	clusterInfo, err := s.GetClusterStatus(req.ClusterStatus)
	if response.FailWithMsg(c, err) {
		return
	}
	var resp []response.ClusterMetaRespStruct
	if response.FailWithMsg(c, dmgrutil.Struct2StructByJson(clusterInfo, &resp)) {
		return
	}
	response.SuccessWithData(c, resp)
}

// 集群销毁
func ClusterDestroy(c *gin.Context) {
	var req request.ClusterOperatorReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	s := service.NewMysqlService()
	// 集群拓扑查询生成
	clusterTopos, err := s.GetClusterTopologyByClusterName(req.ClusterName)
	if response.FailWithMsg(c, err) {
		return
	}

	// 清理集群
	// 注意：清理集群 DestroyInstance 所有组件只会清理子目录，不会清理父目录
	// 比如：deployDir=/data/marvin/{instance_name}, 则清理执行命令 rm -rf /data/marvin/{instance_name}，保留 /data/marvin/ 目录，防止误删除
	for _, component := range dmgrutil.StopComponentOrder {
		for _, t := range clusterTopos {
			compName := strings.ToLower(t.ComponentName)
			if component == compName {
				destroyCompTask := task.NewBuilder().
					SSHKeySet(
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519"),
						filepath.Join(
							dmgrutil.AbsClusterSSHDir(t.ClusterPath, t.ClusterName), "id_ed25519.pub")).
					UserSSH(
						t.MachineHost,
						t.SshPort,
						t.ClusterUser,
						executor.DefaultConnectTimeout,
						module.DefaultSystemdExecuteTimeout).
					StopInstance(t.MachineHost, t.ServicePort, t.InstanceName, t.LogDir,
						fmt.Sprintf("%s-%d.service", t.ComponentName, t.ServicePort),
						module.DefaultSystemdExecuteTimeout).
					DestroyInstance(t.MachineHost, t.ServicePort, t.ComponentName, t.InstanceName, t.DeployDir, t.DataDir, t.LogDir, executor.DefaultExecuteTimeout).BuildTask()

				if response.FailWithMsg(c, destroyCompTask.Execute(ctxt.NewContext())) {
					return
				}
			}
		}
	}

	// 清理元数据信息
	if response.FailWithMsg(c, s.DestroyClusterMetaAndTopology(req.ClusterName)) {
		return
	}

	response.SuccessWithoutData(c)
}

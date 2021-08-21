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

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/wentaojin/dmgr/pkg/dmgrutil"

	"github.com/gin-gonic/gin/binding"
)

// 初始化 gin 参数验证器
func InitGinValidator() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		// 自定义验证方法
		if err := v.RegisterValidation("validIsSkip", validTopologyReqStructSkipCreateUser); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("init gin validator falied")
}

// 验证 TopologyReqStruct 是否跳过创建集群管理用户，验证是否时 true/false
func validTopologyReqStructSkipCreateUser(fl validator.FieldLevel) bool {
	if fl.Field().String() == dmgrutil.BoolTrue || fl.Field().String() == dmgrutil.BoolFalse {
		return true
	}
	return false
}

// 验证集群部署请求字段
func ValidDeployReqStructField(clusterTopo ClusterDeployReqStruct) (ClusterDeployReqStruct, []string) {
	var (
		clusterTopos []TopologyReqStruct
		machineHosts []string
	)

	for _, topo := range clusterTopo.ClusterTopology {
		if topo.PeerPort == 0 {
			topo.PeerPort = dmgrutil.DefaultDmMasterPeerPort
		}
		if topo.ClusterPort == 0 {
			topo.ClusterPort = dmgrutil.DefaultAlertmanagerClusterPort
		}
		if topo.DataDir == "" {
			topo.DataDir = topo.DeployDir
		}
		if topo.LogDir == "" {
			topo.LogDir = topo.DeployDir
		}
		topo.ClusterName = clusterTopo.ClusterMetaReqStruct.ClusterName
		clusterTopos = append(clusterTopos, topo)
		machineHosts = append(machineHosts, topo.MachineHost)
	}

	if clusterTopo.AdminUser == "" {
		clusterTopo.AdminUser = dmgrutil.DefaultGrafanaUser
	}
	if clusterTopo.AdminPassword == "" {
		clusterTopo.AdminPassword = dmgrutil.DefaultGrafanaPassword
	}

	clusterTopo.ClusterTopology = clusterTopos
	uniqueHosts := dmgrutil.NewStringSet(machineHosts...)

	return clusterTopo, uniqueHosts.Slice()
}

// 验证集群扩容请求字段
func ValidScaleOutReqStructField(clusterTopo CLusterScaleOutReqStruct) CLusterScaleOutReqStruct {
	if clusterTopo.PeerPort == 0 {
		clusterTopo.PeerPort = dmgrutil.DefaultDmMasterPeerPort
	}
	if clusterTopo.ClusterPort == 0 {
		clusterTopo.ClusterPort = dmgrutil.DefaultAlertmanagerClusterPort
	}
	if clusterTopo.DataDir == "" {
		clusterTopo.DataDir = clusterTopo.DeployDir
	}
	if clusterTopo.LogDir == "" {
		clusterTopo.LogDir = clusterTopo.DeployDir
	}
	return clusterTopo
}

// 验证集群组件端口是否冲突
func ValidComponentPortConflict(clusterTopo []TopologyReqStruct, dbHostPortArr []string) error {
	var (
		topologyHostPorts []string
		repeatHostPorts   []string
	)

	// 集群部署拓扑中，dm-master 存有 peer-port、alertmanager 存有 cluster-port，其他组件只有服务端口，并未实际占用，只记录
	for _, topo := range clusterTopo {
		if topo.ComponentName == dmgrutil.ComponentDmMaster {
			topologyHostPorts = append(topologyHostPorts, fmt.Sprintf("%s:%v", topo.MachineHost, topo.PeerPort))
		}
		if topo.ComponentName == dmgrutil.ComponentAlertmanager {
			topologyHostPorts = append(topologyHostPorts, fmt.Sprintf("%s:%v", topo.MachineHost, topo.ClusterPort))
		}
		topologyHostPorts = append(topologyHostPorts, fmt.Sprintf("%s:%v", topo.MachineHost, topo.ServicePort))
	}

	repeatHostPorts = dmgrutil.FilterRepeatElem(topologyHostPorts)
	if len(dmgrutil.FilterRepeatElem(topologyHostPorts)) != 0 {
		return fmt.Errorf("cluster deployment topology exist host port conflict: [%v]", repeatHostPorts)
	}

	dbHostPortArr = append(dbHostPortArr, topologyHostPorts...)
	repeatHostPorts = dmgrutil.FilterRepeatElem(dbHostPortArr)
	if len(dmgrutil.FilterRepeatElem(topologyHostPorts)) != 0 {
		return fmt.Errorf("there is a host port conflict between the cluster deployment topology and the existing cluster topology: [%v]", repeatHostPorts)
	}
	return nil
}

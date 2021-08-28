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

import (
	"fmt"

	"github.com/wentaojin/dmgr/pkg/dmgrutil"

	"github.com/jmoiron/sqlx"
	"github.com/wentaojin/dmgr/request"
)

func (s *MysqlService) AddClusterTopology(topo []request.TopologyReqStruct) error {
	if _, err := s.Engine.NamedExec(`INSERT INTO cluster_topology (
cluster_name, 
component_name, 
instance_name,
machine_host,
service_port,
peer_port,
cluster_port,
deploy_dir,
data_dir,
log_dir) VALUES (
:cluster_name, 
:component_name, 
:instance_name,
:machine_host,
:service_port,
:peer_port,
:cluster_port,
:deploy_dir,
:data_dir,
:log_dir)`, topo); err != nil {
		return err
	}
	return nil
}

func (s *MysqlService) GetClusterInstanceName(componentName []string) ([]string, error) {
	var comps []string
	query, args, err := sqlx.In(`SELECT instance_name FROM cluster_topology WHERE component_name IN (?)`, componentName)
	if err != nil {
		return comps, err
	}

	query = s.Engine.Rebind(query)
	if err := s.Engine.Select(&comps, query, args...); err != nil {
		return comps, err
	}
	return comps, nil
}

func (s *MysqlService) GetClusterComponent(clusterName, componentName string) ([]request.TopologyReqStruct, error) {
	var req []request.TopologyReqStruct
	if err := s.Engine.Select(&req, `SELECT cluster_name, 
component_name, 
instance_name,
machine_host,
service_port,
peer_port,
cluster_port,
deploy_dir,
data_dir,
log_dir FROM cluster_topology WHERE cluster_name = ? AND component_name = ?`, clusterName, componentName); err != nil {
		return req, err
	}
	return req, nil
}

func (s *MysqlService) GetClusterComponentByInstance(clusterName string, instanceNames []string) ([]request.TopologyReqStruct, error) {
	var req []request.TopologyReqStruct

	query, args, err := sqlx.In(`SELECT cluster_name, 
component_name, 
instance_name,
machine_host,
service_port,
peer_port,
cluster_port,
deploy_dir,
data_dir,
log_dir FROM cluster_topology WHERE cluster_name = ? AND instance_name IN (?)`, clusterName, instanceNames)
	if err != nil {
		return req, err
	}

	query = s.Engine.Rebind(query)
	if err := s.Engine.Select(&req, query, args...); err != nil {
		return req, err
	}

	return req, nil
}

func (s *MysqlService) DelClusterTopologyByInstanceName(clusterName string, instanceName []string) error {
	query, args, err := sqlx.In(`DELETE FROM cluster_topology WHERE cluster_name = ? AND instance_name IN (?)`, clusterName, instanceName)
	if err != nil {
		return err
	}

	query = s.Engine.Rebind(query)
	if _, err := s.Engine.Exec(query, args...); err != nil {
		return err
	}
	return nil
}

func (s *MysqlService) GetMachinePortArray(machineHosts []string) ([]string, error) {
	var (
		req       []request.TopologyReqStruct
		hostPorts []string
	)
	query, args, err := sqlx.In(`SELECT cluster_name,
component_name, 
instance_name,
machine_host,
service_port,
peer_port,
cluster_port,
deploy_dir,
data_dir,
log_dir FROM cluster_topology WHERE machine_host IN (?)`, machineHosts)
	if err != nil {
		return hostPorts, err
	}

	query = s.Engine.Rebind(query)
	if err := s.Engine.Select(&req, query, args...); err != nil {
		return hostPorts, err
	}

	// 集群部署拓扑中，dm-master 存有 peer-port、alertmanager 存有 cluster-port，其他组件只有服务端口，并未实际占用，只记录
	for _, topo := range req {
		if topo.ComponentName == dmgrutil.ComponentDmMaster {
			hostPorts = append(hostPorts, fmt.Sprintf("%s:%v", topo.MachineHost, topo.PeerPort))
		}
		if topo.ComponentName == dmgrutil.ComponentAlertmanager {
			hostPorts = append(hostPorts, fmt.Sprintf("%s:%v", topo.MachineHost, topo.ClusterPort))
		}
		hostPorts = append(hostPorts, fmt.Sprintf("%s:%v", topo.MachineHost, topo.ServicePort))
	}

	return hostPorts, nil
}

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

	"github.com/jmoiron/sqlx"
	"github.com/wentaojin/dmgr/pkg/dmgrutil"
	"github.com/wentaojin/dmgr/request"
	"github.com/wentaojin/dmgr/response"
)

// 过滤组件实例名
func (s *MysqlService) FilterComponentInstance(req request.ClusterOperatorReqStruct) ([]string, error) {
	var instNames []string
	// 判断指定实例名是否在指定组件中 [组件操作以实例名为准，实例名全局唯一]
	if len(req.ComponentName) != 0 {
		instNameByDB, err := s.GetClusterInstanceName(req.ComponentName)
		if err != nil {
			return instNames, err
		}
		if len(req.InstanceName) != 0 {
			originInsts := dmgrutil.NewStringSet(instNameByDB...)
			requestInsts := dmgrutil.NewStringSet(req.InstanceName...)
			diffInsts := requestInsts.Difference(originInsts).Slice()
			if len(diffInsts) != 0 {
				return instNames, fmt.Errorf("request component instance [%v] not exists in the component", diffInsts)
			}
			instNames = req.InstanceName
		} else {
			instNames = instNameByDB
		}
	}
	return instNames, nil
}

// 集群拓扑查询 - 集群名
func (s *MysqlService) GetClusterTopologyByClusterName(clusterName string) ([]response.ClusterTopologyRespStruct, error) {
	var clusterTopo []response.ClusterTopologyRespStruct
	if err := s.Engine.Select(&clusterTopo, `SELECT
	meta.cluster_name,
	meta.cluster_user,
	meta.cluster_version,
	meta.cluster_path,
	meta.admin_user,
	meta.admin_password,
	topo.component_name,
	topo.instance_name,
	topo.machine_host,
	topo.service_port,
	topo.peer_port,
	topo.cluster_port,
	topo.deploy_dir,
	topo.data_dir,
	topo.log_dir,
	mh.ssh_user,
	mh.ssh_password,
	mh.ssh_port 
FROM
    cluster_meta meta INNER JOIN cluster_topology topo ON meta.cluster_name = topo.cluster_name AND meta.cluster_name = ? 
	LEFT JOIN machine mh ON topo.machine_host = mh.ssh_host`, clusterName); err != nil {
		return clusterTopo, err
	}

	return clusterTopo, nil
}

// 集群拓扑查询 - 集群名/实例名
func (s *MysqlService) GetClusterTopologyByInstanceName(clusterName string, instanceNames []string) ([]response.ClusterTopologyRespStruct, error) {
	var (
		clusterTopo []response.ClusterTopologyRespStruct
		query       string
		args        []interface{}
		err         error
	)
	if len(instanceNames) == 0 {
		query, args, err = sqlx.In(`SELECT
	meta.cluster_name,
	meta.cluster_user,
	meta.cluster_version,
	meta.cluster_path,
	meta.admin_user,
	meta.admin_password,
	topo.component_name,
	topo.instance_name,
	topo.machine_host,
	topo.service_port,
	topo.peer_port,
	topo.cluster_port,
	topo.deploy_dir,
	topo.data_dir,
	topo.log_dir,
	mh.ssh_user,
	mh.ssh_password,
	mh.ssh_port 
FROM
    cluster_meta meta INNER JOIN cluster_topology topo ON meta.cluster_name = topo.cluster_name LEFT JOIN machine mh ON topo.machine_host = mh.ssh_host 
WHERE meta.cluster_name = ?`, clusterName)
	} else {
		query, args, err = sqlx.In(`SELECT
	meta.cluster_name,
	meta.cluster_user,
	meta.cluster_version,
	meta.cluster_path,
	meta.admin_user,
	meta.admin_password,
	topo.component_name,
	topo.instance_name,
	topo.machine_host,
	topo.service_port,
	topo.peer_port,
	topo.cluster_port,
	topo.deploy_dir,
	topo.data_dir,
	topo.log_dir,
	mh.ssh_user,
	mh.ssh_password,
	mh.ssh_port 
FROM
    cluster_meta meta INNER JOIN cluster_topology topo ON meta.cluster_name = topo.cluster_name LEFT JOIN machine mh ON topo.machine_host = mh.ssh_host
WHERE meta.cluster_name = ? AND topo.instance_name IN (?)`, clusterName, instanceNames)
	}
	if err != nil {
		return clusterTopo, err
	}
	query = s.Engine.Rebind(query)
	if err := s.Engine.Select(&clusterTopo, query, args...); err != nil {
		return clusterTopo, err
	}

	return clusterTopo, nil
}

// 集群元信息以及拓扑清理
func (s *MysqlService) DestroyClusterMetaAndTopology(clusterName string) error {
	tx, err := s.Engine.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM cluster_topology WHERE cluster_name = ?`, clusterName)
	_, err = tx.Exec(`DELETE FROM cluster_meta WHERE cluster_name = ?`, clusterName)

	if err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// 集群拓扑以及元数据新增
func (s *MysqlService) AddClusterMetaAndTopology(clusterMeta request.ClusterMetaReqStruct, topo []request.TopologyReqStruct) error {
	tx, err := s.Engine.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.NamedExec(`INSERT INTO cluster_meta (
cluster_name, 
cluster_user, 
cluster_version,
cluster_path,
admin_user,
admin_password,
skip_create_user) VALUES (
:cluster_name, 
:cluster_user, 
:cluster_version,
:cluster_path,
:admin_user,
:admin_password,
:skip_create_user)`, clusterMeta)
	_, err = tx.NamedExec(`INSERT INTO cluster_topology (
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
:log_dir)`, topo)

	if err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

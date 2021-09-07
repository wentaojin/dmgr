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
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/wentaojin/dmgr/pkg/dmgrutil"
	"github.com/wentaojin/dmgr/request"
	"github.com/wentaojin/dmgr/response"
)

// 获取数据源配置
func (s *MysqlService) GetSourceConf(clusterName, taskName, sourceName string) (response.TaskSourceConfRespStruct, error) {
	var resp response.TaskSourceConfRespStruct
	if err := s.Engine.Get(&resp, `SELECT 
s.cluster_name,
s.task_name,
s.source_name,
c.host,
c.user,
c.password,
c.port,
c.ssl_ca,
c.ssl_cert,
c.ssl_key,
c.label,
s.enable_gtid,
s.relay_binlog_gtid,
s.enable_relay,
s.relay_binlog_name,
s.relay_dir,
s.purge_interval,
s.purge_expires,
s.purge_remain_space,
s.checker_check_enable,
s.checker_backoff_rollback,
s.checker_backoff_max
FROM
task_source s,
task_cluster c
WHERE 
s.source_name=c.source_name
AND s.cluster_name = ?
AND s.task_name = ?
AND s.source_name = ?`, clusterName, taskName, sourceName); err != nil {
		return resp, err
	}
	return resp, nil
}

// 上游数据源信息新增
func (s *MysqlService) AddTaskSource(source request.TaskSourceCreateReqStruct) error {
	if _, err := s.Engine.NamedExec(`INSERT INTO task_source (
source_name,
host,
user,
password,
port,
ssl_ca,
ssl_cert,
ssl_key,
label) VALUES (
:source_name,
:host,
:user,
:password,
:port,
:ssl_ca,
:ssl_cert,
:ssl_key,
:label       
)`, source); err != nil {
		return err
	}
	return nil
}

// 获取上游数据源信息
func (s *MysqlService) GetTaskSourceBySourceName(sourceName string) (response.TaskSourceRespStruct, error) {
	var resp response.TaskSourceRespStruct
	if err := s.Engine.Get(&resp, `SELECT * FROM task_source WHERE source_name = ?`, sourceName); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *MysqlService) GetTaskSourceByLabel(label string) ([]response.TaskSourceRespStruct, error) {
	var resp []response.TaskSourceRespStruct
	if err := s.Engine.Select(&resp, `SELECT * FROM task_source WHERE label = ?`, label); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *MysqlService) GetTaskSourceALL() ([]response.TaskSourceRespStruct, error) {
	var resp []response.TaskSourceRespStruct
	if err := s.Engine.Select(&resp, `SELECT * FROM task_source`); err != nil {
		return resp, err
	}
	return resp, nil
}

// 更新上游数据源信息
func (s *MysqlService) UpdateTaskSource(reqStruct request.TaskSourceUpdateReqStruct) error {
	return Transact(s.Engine, func(tx *sqlx.Tx) error {
		if _, err := tx.NamedExec(`UPDATE task_source SET
host = :host,
user = :user,
password = :password,
port = :port,
ssl_ca = :ssl_ca,
ssl_cert = :ssl_cert,
ssl_key = :ssl_key,
label = :label WHERE source_name = :source_name`, reqStruct.TaskSourceCreateReqStruct); err != nil {
			return err
		}

		reqStruct.TaskSourceCreateReqStruct.ClusterName = reqStruct.ClusterName
		reqStruct.TaskSourceCreateReqStruct.TaskName = reqStruct.TaskName

		if _, err := tx.NamedExec(`UPDATE task_cluster SET
enable_gtid=:enable_gtid,
relay_binlog_gtid=:relay_binlog_gtid,
enable_relay=:enable_relay,
relay_binlog_name=:relay_binlog_name,
relay_dir=:relay_dir,
purge_interval=:purge_interval,
purge_expires=:purge_expires,
purge_remain_space=:purge_remain_space,
checker_check_enable=:checker_check_enable,
checker_backoff_rollback=:checker_backoff_rollback,
checker_backoff_max=:checker_backoff_max,
target_name=:target_name WHERE cluster_name = :cluster_name AND task_name = :task_name AND source_name = :source_name`, reqStruct.TaskSourceCreateReqStruct); err != nil {
			return err
		}
		return nil
	})
}

// 获取任务集群信息
func (s *MysqlService) GetTaskClusterBySourceName(sourceName string) ([]response.TaskClusterRespStruct, error) {
	var taskMeta []response.TaskClusterRespStruct
	if err := s.Engine.Select(&taskMeta, `SELECT * FROM task_cluster WHERE source_name = ?`, sourceName); err != nil {
		return taskMeta, err
	}
	return taskMeta, nil
}

// 创建任务集群信息
func (s *MysqlService) AddTaskCluster(reqStruct request.TaskCLusterReqStruct) error {
	var taskSources []request.TaskCLusterReqStruct
	sources := strings.Split(reqStruct.SourceName, dmgrutil.TaskSourceDelimiter)
	for _, source := range sources {
		reqStruct.TaskName = source
		taskSources = append(taskSources, reqStruct)
	}

	for _, source := range taskSources {
		if _, err := s.Engine.NamedExec(`INSERT INTO task_cluster (
cluster_name,
task_name,
source_name,
enable_gtid,
relay_binlog_gtid,
enable_relay,
relay_binlog_name,
relay_dir,
purge_interval,
purge_expires,
purge_remain_space,
checker_check_enable,
checker_backoff_rollback,
checker_backoff_max,
target_name) VALUES (
:cluster_name,
:task_name,
:source_name,
:enable_gtid,
:relay_binlog_gtid,
:enable_relay,
:relay_binlog_name,
:relay_dir,
:purge_interval,
:purge_expires,
:purge_remain_space,
:checker_check_enable,
:checker_backoff_rollback,
:checker_backoff_max,
:target_name)`, source); err != nil {
			return err
		}
	}

	return nil
}

// 删除上游数据源信息
func (s *MysqlService) DeleteTaskSource(reqStruct request.TaskSourceDeleteReqStruct) error {
	return Transact(s.Engine, func(tx *sqlx.Tx) error {
		if _, err := tx.Exec(`DELETE FROM task_cluster WHERE cluster_name = ? AND task_name = ? AND source_name = ?`,
			reqStruct.ClusterName, reqStruct.TaskName, reqStruct.SourceName); err != nil {
			return err
		}
		if _, err := tx.Exec(`DELETE FROM task_source WHERE source_name = ?`, reqStruct.SourceName); err != nil {
			return err
		}
		return nil
	})
}

// 下游数据源信息新增
func (s *MysqlService) AddTaskTarget(target request.TaskTargetCreateReqStruct) error {
	if _, err := s.Engine.NamedExec(`INSERT INTO task_target (
target_id,
host,
user,
password,
port,
target_packet,
sql_mode,
skip_utf8_check,
constraint_check_in_place,
ssl_ca,
ssl_cert,
ssl_key,
label
) VALUES (
:target_id,
:host,
:user,
:password,
:port,
:target_packet,
:sql_mode,
:skip_utf8_check,
:constraint_check_in_place,
:ssl_ca,
:ssl_cert,
:ssl_key,
:label)`, target); err != nil {
		return err
	}
	return nil
}

// 获取下游数据源信息
func (s *MysqlService) GetTaskTargetByTargetName(targetName string) (response.TaskTargetRespStruct, error) {
	var resp response.TaskTargetRespStruct
	if err := s.Engine.Get(&resp, `SELECT * FROM task_target WHERE target_name = ?`, targetName); err != nil {
		return resp, err
	}
	return resp, nil
}

// 获取任务集群信息
func (s *MysqlService) GetTaskClusterByTargetName(targetName string) ([]response.TaskClusterRespStruct, error) {
	var taskMeta []response.TaskClusterRespStruct
	if err := s.Engine.Select(&taskMeta, `SELECT * FROM task_cluster WHERE target_name = ?`, targetName); err != nil {
		return taskMeta, err
	}
	return taskMeta, nil
}

// 删除下游数据源信息
func (s *MysqlService) DeleteTaskTarget(reqStruct request.TaskTargetDeleteReqStruct) error {
	return Transact(s.Engine, func(tx *sqlx.Tx) error {
		if _, err := tx.Exec(`DELETE FROM task_cluster WHERE cluster_name = ? AND task_name = ? AND target_name = ?`,
			reqStruct.ClusterName, reqStruct.TaskName, reqStruct.TargetName); err != nil {
			return err
		}
		if _, err := tx.Exec(`DELETE FROM task_target WHERE target_name = ?`, reqStruct.TargetName); err != nil {
			return err
		}
		return nil
	})
}

// 更新下游数据源信息
func (s *MysqlService) UpdateTaskTarget(reqStruct request.TaskTargetUpdateReqStruct) error {
	return Transact(s.Engine, func(tx *sqlx.Tx) error {
		if _, err := tx.NamedExec(`UPDATE task_target SET
host = :host,
user = :user,
password = :password,
port = :port,
max_allowed_packet = :max_allowed_packet,
session=:session,
ssl_ca = :ssl_ca,
ssl_cert = :ssl_cert,
ssl_key = :ssl_key,
label = :label WHERE target_name = :target_name`, reqStruct.TaskTargetCreateReqStruct); err != nil {
			return err
		}
		return nil
	})
}

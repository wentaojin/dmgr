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
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/wentaojin/dmgr/pkg/cluster/api"
	"github.com/wentaojin/dmgr/pkg/dmgrutil"
	"github.com/wentaojin/dmgr/request"
	"github.com/wentaojin/dmgr/response"
	"github.com/wentaojin/dmgr/service"
)

// 上游数据源创建
func TaskSourceCreate(c *gin.Context) {
	var req request.TaskSourceCreateReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	s := service.NewMysqlService()
	if response.FailWithMsg(c, s.AddTaskSource(req)) {
		return
	}
	response.SuccessWithoutData(c)
}

// 下游数据源创建
func TaskTargetCreate(c *gin.Context) {
	var req request.TaskTargetCreateReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	s := service.NewMysqlService()
	if response.FailWithMsg(c, s.AddTaskTarget(req)) {
		return
	}
	response.SuccessWithoutData(c)
}

// 任务集群数据关系映射创建
func TaskClusterCreate(c *gin.Context) {
	var req request.TaskCLusterReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	s := service.NewMysqlService()
	if response.FailWithMsg(c, s.AddTaskCluster(req)) {
		return
	}
	response.SuccessWithoutData(c)
}

// 上游数据源删除
func TaskSourceDelete(c *gin.Context) {
	var req request.TaskSourceDeleteReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	// 1. 获取 source name
	// 2. 判断 source name 引用关系
	// 	- 未被引用直接删除
	//	- 存在多个引用不能删除
	//  - 符合条件引用删除
	s := service.NewMysqlService()
	sourceInfo, err := s.GetTaskSourceBySourceName(req.SourceName)
	if response.FailWithMsg(c, err) {
		return
	}
	if dmgrutil.IsStructureEqual(sourceInfo, response.TaskSourceRespStruct{}) {
		if response.FailWithMsg(c, fmt.Errorf("task source cannot exist")) {
			return
		}
	}

	taskCluster, err := s.GetTaskClusterBySourceName(req.SourceName)
	if response.FailWithMsg(c, err) {
		return
	}

	if len(taskCluster) > 1 {
		if response.FailWithMsg(c, fmt.Errorf("source name [%v] is referenced by the multiple [%v] cluster task, cannot be deleted", req.SourceName, len(taskCluster))) {
			return
		}
	}

	if dmgrutil.IsStructureEqual(taskCluster, []response.TaskClusterRespStruct{}) {
		if response.FailWithMsg(c, s.DeleteTaskSource(req)) {
			return
		}
	}

	// 1. 获取 source 状态【状态判断，删除前任务停止】
	// 2. 停止 source 任务
	// 3. 删除 dm-master source 任务
	// 4. 清理数据库元数据信息
	dmMasterUrl, err := GetActiveDmMasterAddr(s, req.ClusterName)
	if response.FailWithMsg(c, err) {
		return
	}

	respByte, err := api.GetSourceStatusBySourceName(dmMasterUrl, req.SourceName)
	if response.FailWithMsg(c, err) {
		return
	}

	relayStatusStage := gjson.GetBytes(respByte, api.RelayStatusPath).String()
	workerName := gjson.GetBytes(respByte, api.WorkerNamePath).String()
	if relayStatusStage == api.RelayStatusRunningStage {
		jsonWK, err := dmgrutil.Struct2Json(api.NewWorkerNameBody(workerName))
		if response.FailWithMsg(c, err) {
			return
		}
		_, err = api.StopSourceBySourceName(dmMasterUrl, req.SourceName, strings.NewReader(jsonWK))
		if response.FailWithMsg(c, err) {
			return
		}
	}

	_, err = api.DeleteSourceBySourceName(dmMasterUrl, req.SourceName)
	if response.FailWithMsg(c, err) {
		return
	}

	if response.FailWithMsg(c, s.DeleteTaskSource(req)) {
		return
	}
	response.SuccessWithoutData(c)
}

// 上游数据源修改
// 禁止修改 source_name
func TaskSourceUpdate(c *gin.Context) {
	var req request.TaskSourceUpdateReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	// 1. 获取 source name
	// 2. 查询 source name 引用关系
	// 	- 根据引用关系判断是否可修改
	s := service.NewMysqlService()
	sourceInfo, err := s.GetTaskSourceBySourceName(req.SourceName)
	if response.FailWithMsg(c, err) {
		return
	}
	if dmgrutil.IsStructureEqual(sourceInfo, response.TaskSourceRespStruct{}) {
		if response.FailWithMsg(c, fmt.Errorf("task source cannot exist")) {
			return
		}
	}

	taskCluster, err := s.GetTaskClusterBySourceName(req.SourceName)
	if response.FailWithMsg(c, err) {
		return
	}

	if !dmgrutil.IsStructureEqual(req.TaskSourceCreateReqStruct, sourceInfo.TaskSourceCreateReqStruct) && len(taskCluster) > 1 {
		if response.FailWithMsg(c, fmt.Errorf("task source [%v] cannot update, has be referenced by the multiple [%v] cluster task, cannot be update", req.SourceName, len(taskCluster))) {
			return
		}
	}

	if dmgrutil.IsStructureEqual(taskCluster, []response.TaskClusterRespStruct{}) {
		if response.FailWithMsg(c, s.UpdateTaskSource(req)) {
			return
		}
	}

	for _, task := range taskCluster {
		if task.ClusterName == req.ClusterName && task.TaskName == req.TaskName && task.SourceName == req.SourceName {
			dmMasterUrl, err := GetActiveDmMasterAddr(s, req.ClusterName)
			if response.FailWithMsg(c, err) {
				return
			}

			respByte, err := api.GetSourceStatusBySourceName(dmMasterUrl, req.SourceName)
			if response.FailWithMsg(c, err) {
				return
			}

			relayStatusStage := gjson.GetBytes(respByte, api.RelayStatusPath).String()
			workerName := gjson.GetBytes(respByte, api.WorkerNamePath).String()
			if relayStatusStage == api.RelayStatusRunningStage {
				jsonWK, err := dmgrutil.Struct2Json(api.NewWorkerNameBody(workerName))
				if response.FailWithMsg(c, err) {
					return
				}
				_, err = api.StopSourceBySourceName(dmMasterUrl, req.SourceName, strings.NewReader(jsonWK))
				if response.FailWithMsg(c, err) {
					return
				}
			}

			// todo: 待完善
			// 1. 启动 source (待补充完善)
			jsonSRC, err := dmgrutil.Struct2Json(api.NewRelayStatusBody(respByte))
			if response.FailWithMsg(c, err) {
				return
			}
			_, err = api.StartSourceBySourceName(dmMasterUrl, req.SourceName, strings.NewReader(jsonSRC))
			if response.FailWithMsg(c, err) {
				return
			}

			if response.FailWithMsg(c, s.UpdateTaskSource(req)) {
				return
			}
		}
	}
	response.SuccessWithoutData(c)
}

// 下游数据源删除
func TaskTargetDelete(c *gin.Context) {
	var req request.TaskTargetDeleteReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	// 1. 获取 target name
	// 2. 判断 target name 引用关系
	// 	- 未被引用直接删除
	//	- 存在多个引用不能删除
	//  - 符合条件引用删除
	s := service.NewMysqlService()
	targetInfo, err := s.GetTaskSourceBySourceName(req.TargetName)
	if response.FailWithMsg(c, err) {
		return
	}
	if dmgrutil.IsStructureEqual(targetInfo, response.TaskTargetRespStruct{}) {
		if response.FailWithMsg(c, fmt.Errorf("task target cannot exist")) {
			return
		}
	}

	taskCluster, err := s.GetTaskClusterByTargetName(req.TargetName)
	if response.FailWithMsg(c, err) {
		return
	}

	if len(taskCluster) > 1 {
		if response.FailWithMsg(c, fmt.Errorf("target name [%v] is referenced by the multiple [%v] cluster task, cannot be deleted", req.TargetName, len(taskCluster))) {
			return
		}
	}

	if dmgrutil.IsStructureEqual(taskCluster, []response.TaskClusterRespStruct{}) {
		if response.FailWithMsg(c, s.DeleteTaskTarget(req)) {
			return
		}
	}

	dmMasterUrl, err := GetActiveDmMasterAddr(s, req.ClusterName)
	if response.FailWithMsg(c, err) {
		return
	}

	// todo: 待完善
	// 1. 停止同步任务（待补充）
	// 2. 删除同步任务
	_, err = api.DeleteTaskByTaskName(dmMasterUrl, req.TaskName)
	if response.FailWithMsg(c, err) {
		return
	}

	if response.FailWithMsg(c, s.DeleteTaskTarget(req)) {
		return
	}
	response.SuccessWithoutData(c)
}

// 下游数据源修改
// 禁止修改 target_name
func TaskTargetUpdate(c *gin.Context) {
	var req request.TaskTargetUpdateReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	// 1. 获取 source name
	// 2. 查询 source name 引用关系
	// 	- 根据引用关系判断是否可修改
	s := service.NewMysqlService()
	targetInfo, err := s.GetTaskTargetByTargetName(req.TargetName)
	if response.FailWithMsg(c, err) {
		return
	}
	if dmgrutil.IsStructureEqual(targetInfo, response.TaskTargetRespStruct{}) {
		if response.FailWithMsg(c, fmt.Errorf("task source cannot exist")) {
			return
		}
	}

	taskCluster, err := s.GetTaskClusterByTargetName(req.TargetName)
	if response.FailWithMsg(c, err) {
		return
	}

	if !dmgrutil.IsStructureEqual(req.TaskTargetCreateReqStruct, targetInfo.TaskTargetCreateReqStruct) && len(taskCluster) > 1 {
		if response.FailWithMsg(c, fmt.Errorf("task target [%v] cannot update, has be referenced by the multiple [%v] cluster task, cannot be update", req.TargetName, len(taskCluster))) {
			return
		}
	}

	if dmgrutil.IsStructureEqual(taskCluster, []response.TaskClusterRespStruct{}) {
		if response.FailWithMsg(c, s.UpdateTaskTarget(req)) {
			return
		}
	}

	for _, task := range taskCluster {
		if task.ClusterName == req.ClusterName && task.TaskName == req.TaskName && task.TargetName == req.TargetName {
			dmMasterUrl, err := GetActiveDmMasterAddr(s, req.ClusterName)
			if response.FailWithMsg(c, err) {
				return
			}

			// todo: 待完善
			// 1. 停止任务同步 (待补充完善)
			// 2. 重新启动任务同步

			_, err = api.StartTaskMigration(dmMasterUrl, strings.NewReader("task struct"))
			if response.FailWithMsg(c, err) {
				return
			}

			if response.FailWithMsg(c, s.UpdateTaskTarget(req)) {
				return
			}
		}
	}
	response.SuccessWithoutData(c)
}

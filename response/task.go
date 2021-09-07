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
package response

import (
	"time"

	"github.com/wentaojin/dmgr/request"
)

// 查询任务 source 信息响应
type TaskSourceConfRespStruct struct {
	request.TaskCLusterReqStruct
	request.TaskDatasourceStruct
	request.TaskDatasourceSslStruct
	Label string `json:"label" form:"label" db:"label"`
}

// 查询上游数据源信息响应
type TaskSourceRespStruct struct {
	request.TaskSourceCreateReqStruct
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
}

// 查询下游数据源响应
type TaskTargetRespStruct struct {
	request.TaskTargetCreateReqStruct
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
}

// 查询任务集群映射响应
type TaskClusterRespStruct struct {
	request.TaskCLusterReqStruct
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
}

// 查询任务同步元数据响应
type TaskMetaRespStruct struct {
	request.TaskMetaReqStruct
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
}

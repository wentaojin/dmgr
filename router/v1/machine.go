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
	"github.com/gin-gonic/gin"
	"github.com/wentaojin/dmgr/request"
	"github.com/wentaojin/dmgr/response"
	"github.com/wentaojin/dmgr/service"
)

// 新增机器
func AddMachine(c *gin.Context) {
	var req request.MachineReqStruct
	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	s := service.NewMysqlService()
	if response.FailWithMsg(c, s.AddMachine(req)) {
		return
	}
	response.SuccessWithoutData(c)
}

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

	"github.com/wentaojin/dmgr/response"

	"github.com/wentaojin/dmgr/pkg/dmgrutil"

	"github.com/gin-gonic/gin"
	"github.com/wentaojin/dmgr/request"
	"github.com/wentaojin/dmgr/service"
)

// 获取当前用户
func GetCurrentUser(c *gin.Context) response.UserRespStruct {
	user, exists := c.Get("user")
	var newUser = response.UserRespStruct{}
	if !exists {
		return newUser
	}
	u, _ := user.(*response.UserRespStruct)
	// 创建服务
	s := service.NewMysqlService()
	newUser, _ = s.GetUserPasswordByUsername(u.Username)
	return newUser
}

// 获取当前用户信息
func GetUserInfo(c *gin.Context) {
	user := GetCurrentUser(c)
	// 转为 UserInfoResponseStruct, 隐藏部分字段
	var resp response.UserRespStruct
	if response.FailWithMsg(c, dmgrutil.Struct2StructByJson(user, &resp)) {
		return
	}
	response.SuccessWithData(c, resp)
}

// 修改密码
func ChangePwd(c *gin.Context) {
	// 请求 json 绑定
	var req request.ChangePwdReqStruct

	if response.FailWithMsg(c, c.ShouldBindJSON(&req)) {
		return
	}

	// 获取当前用户
	user := GetCurrentUser(c)

	s := service.NewMysqlService()
	newUser, err := s.GetUserPasswordByUsername(user.Username)
	if response.FailWithMsg(c, err) {
		return
	}
	// 密码校验
	eCryptPwd, err := dmgrutil.AesEcryptCode([]byte(req.OldPassword))
	if response.FailWithMsg(c, err) {
		return
	}
	if ok := dmgrutil.CompareAesPwd(eCryptPwd, newUser.Password); !ok {
		if response.FailWithMsg(c, fmt.Errorf("origin password error")) {
			return
		}
	}
	// 密码更新
	eCryptPwd, err = dmgrutil.AesEcryptCode([]byte(req.NewPassword))
	if response.FailWithMsg(c, err) {
		return
	}

	err = s.UpdateUserPasswordByUsername(user.Username, eCryptPwd)
	if response.FailWithMsg(c, err) {
		return
	}

	response.SuccessWithoutData(c)
}

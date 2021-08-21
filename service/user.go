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

	"github.com/pingcap/errors"
	"github.com/wentaojin/dmgr/pkg/dmgrutil"
	"github.com/wentaojin/dmgr/response"
)

// 登录校验
func (s *MysqlService) LoginCheck(user *response.UserRespStruct) (*response.UserRespStruct, error) {
	var u response.UserRespStruct
	// 查询用户
	if err := s.Engine.Get(&u, "SELECT username, password FROM user WHERE username = ?", user.Username); err != nil {
		return nil, fmt.Errorf("failed get username [%s] by db: %v", user.Username, err)
	}

	// 校验密码
	decryptPWD, err := dmgrutil.AesDeCryptCode(u.Password)
	if err != nil {
		return nil, fmt.Errorf("password aes desCrypt failed: %v", err)
	}
	if ok := dmgrutil.CompareAesPwd(user.Password, string(decryptPWD)); !ok {
		return nil, errors.New(response.LoginCheckErrorMsg)
	}

	return &u, err
}

// 获取单个用户密码
func (s *MysqlService) GetUserPasswordByUsername(username string) (response.UserRespStruct, error) {
	var user response.UserRespStruct
	if err := s.Engine.Get(&user, "SELECT id, username, password FROM user WHERE username = ?", username); err != nil {
		return user, fmt.Errorf("failed get username [%s] by db: %v", user.Username, err)
	}
	return user, nil
}

// 修改用户密码
func (s *MysqlService) UpdateUserPasswordByUsername(username, password string) error {
	if _, err := s.Engine.Exec("UPDATE user set password = ? WHERE username = ?", password, username); err != nil {
		return fmt.Errorf("failed update username [%s] password [%s] by db: %v", username, password, err)
	}
	return nil
}

// 初始化管理用户
func (s *MysqlService) initUserTableData(username, password string) error {
	var userCount int
	if err := s.Engine.Get(&userCount, "SELECT count(1) FROM user WHERE username = ?", username); err != nil {
		return fmt.Errorf("failed get username [%s] counts by db: %v", username, err)
	}
	if userCount == 0 {
		if _, err := s.Engine.NamedExec("INSERT INTO user (username, password) VALUES (:username, :password)", map[string]interface{}{
			"username": username,
			"password": password,
		}); err != nil {
			return err
		}
	}
	return nil
}

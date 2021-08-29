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

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/wentaojin/dmgr/pkg/dmgrutil"
)

var Engine *sqlx.DB

func NewMySQLEngineDB(cfg *dmgrutil.DbConfig) (err error) {
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	Engine, err = sqlx.Connect("mysql", DSN)
	if err != nil {
		return errors.Errorf("connect server failed, err: %v\n", err)
	}
	Engine.SetMaxOpenConns(200)
	Engine.SetMaxIdleConns(10)
	return nil
}

func SyncMysqlEngineDB() error {
	if _, err := Engine.Exec(ClusterTables, TaskTables); err != nil {
		return errors.Errorf("create table struct failed, err: %v\n", err)
	}
	if err := initMysqlEngineData(); err != nil {
		return err
	}
	return nil
}

func initMysqlEngineData() error {
	// 1. 初始化用户
	encryptSuperPwd, err := dmgrutil.AesEcryptCode([]byte("admin"))
	if err != nil {
		return fmt.Errorf("falied make encrypt super admin user password: %v", err)
	}

	if err := NewMysqlService().initUserTableData("admin", encryptSuperPwd); err != nil {
		return err
	}
	return nil
}

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
package main

import (
	"flag"
	"log"

	"github.com/wentaojin/dmgr/service"

	"github.com/wentaojin/dmgr/router"

	"go.uber.org/zap"

	"github.com/wentaojin/dmgr/pkg/dmgrutil"
)

var (
	config  = flag.String("config", "config.toml", "specify the configuration file, default is config.toml")
	version = flag.Bool("version", false, "view DMGR version info")
)

func main() {
	flag.Parse()

	// 1. 获取程序版本以及读取配置文件
	dmgrutil.GetAppVersion(*version)
	cfg, err := dmgrutil.ReadConfigFile(*config)
	if err != nil {
		log.Fatalf("failed read config file %s: %v", *config, err)
	}

	// 2. 初始化日志记录器
	if err := dmgrutil.InitLogger(&cfg.LogConfig); err != nil {
		panic(err)
	}
	dmgrutil.RecordAppVersion("dmgr", dmgrutil.Logger, cfg)

	// 3. 创建数据库连接以及初始化表结构
	if err := service.NewMySQLEngineDB(&cfg.DbConfig); err != nil {
		dmgrutil.Logger.Fatal("mysql open error", zap.Error(err))
	}
	if err := service.SyncMysqlEngineDB(); err != nil {
		dmgrutil.Logger.Fatal("mysql sync error", zap.Error(err))
	}

	// 4. 程序运行
	if err := router.Run(cfg); err != nil {
		dmgrutil.Logger.Fatal("server run error", zap.Error(err))
	}
}

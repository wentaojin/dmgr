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
package router

import (
	"fmt"

	"github.com/wentaojin/dmgr/response"

	"github.com/wentaojin/dmgr/request"

	"github.com/wentaojin/dmgr/pkg/dmgrutil"
	"github.com/wentaojin/dmgr/pkg/middleware"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// 程序运行
func Run(file *dmgrutil.Config) error {
	r := registerRouter(&file.MiddlewareConfig)
	if err := r.Run(file.ListenPort); err != nil {
		return err
	}
	return nil
}

// 请求路由注册
func registerRouter(file *dmgrutil.MiddlewareConfig) *gin.Engine {
	r := gin.New()

	// 添加日志中间件
	r.Use(dmgrutil.GinLogger(dmgrutil.Logger))

	// Recover Panic
	r.Use(response.GinRecovery(dmgrutil.Logger, true))

	// 注册自定义验证器
	if err := request.InitGinValidator(); err != nil {
		dmgrutil.Logger.DPanic("middleware init-validator", zap.String("success", "init middleware validator success"))
	}
	dmgrutil.Logger.Info("middleware init-validator", zap.String("success", "init middleware validator success"))

	// 添加速率访问中间件
	r.Use(middleware.RateLimiter(file.MaxRateLimiter))
	dmgrutil.Logger.Info("middleware rate-limiter", zap.String("success", "init middleware rate limiter success"))

	// 添加跨域中间件, 让请求支持跨域
	r.Use(middleware.Cors())
	dmgrutil.Logger.Info("middleware cross-domain", zap.String("success", "init gin request cross-domain success"))

	// 初始化 jwt auth 中间件
	authMiddleware, err := middleware.JwtAuth(file)
	if err != nil {
		panic(fmt.Sprintf("init jwt auth middleware: %v", err))
	}
	dmgrutil.Logger.Info("middleware jwt-auth", zap.String("success", "init jwt auth middleware success"))

	// 方便统一添加路由前缀
	v1Group := r.Group("v1")
	InitBaseRouter(v1Group, authMiddleware)       // 注册基础路由, 不会鉴权
	InitUserRouter(v1Group, authMiddleware)       // 注册用户路由
	InitMachineRouter(v1Group, authMiddleware)    // 注册机器路由
	InitWarehouseRouter(v1Group, authMiddleware)  // 注册离线包路由
	InitClusterRouter(v1Group, authMiddleware)    // 注册集群管理路由
	InitDatasourceRouter(v1Group, authMiddleware) // 注册任务数据源路由
	InitTaskRouter(v1Group, authMiddleware)       // 注册任务管理路由

	return r
}

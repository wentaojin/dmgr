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
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "github.com/wentaojin/dmgr/router/v1"
)

// 基础路由
func InitBaseRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/logout", authMiddleware.LogoutHandler)
	r.POST("/refreshToken", authMiddleware.RefreshHandler)
	return r
}

// 用户路由
func InitUserRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("/user").Use(authMiddleware.MiddlewareFunc())
	{
		router.POST("/info", v1.GetUserInfo)
		router.PUT("/changePwd", v1.ChangePwd)
	}
	return router
}

// 机器路由
func InitMachineRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("/machine").Use(authMiddleware.MiddlewareFunc())
	{
		router.POST("/add", v1.AddMachine)
	}
	return router
}

// 离线包路由
func InitWarehouseRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("/warehouse").Use(authMiddleware.MiddlewareFunc())
	{
		router.POST("/upload", v1.FileChunkUpload)
	}
	return router
}

// 集群路由
func InitClusterRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("/cluster").Use(authMiddleware.MiddlewareFunc())
	{
		router.POST("/deploy", v1.ClusterDeploy)
		router.POST("/start", v1.ClusterStart)
		router.POST("/stop", v1.ClusterStop)
		router.POST("/scale-out", v1.ClusterScaleOut)
		router.POST("/scale-in", v1.ClusterScaleIn)
		router.POST("/reload", v1.CLusterReload)
		router.POST("/upgrade", v1.ClusterUpgrade)
		router.POST("/destroy", v1.ClusterDestroy)
		router.POST("/patch", v1.ClusterPatch)
		router.POST("/status", v1.ClusterStatus)
	}
	return router
}

// 数据源路由
func InitDatasourceRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("/datasource").Use(authMiddleware.MiddlewareFunc())
	{
		router.POST("/source-create", v1.TaskSourceCreate)
		router.POST("/source-delete", v1.TaskSourceDelete)
		router.PATCH("/source-update", v1.TaskSourceUpdate)
		router.POST("/target-create", v1.TaskTargetCreate)
		router.POST("/target-delete", v1.TaskTargetDelete)
		router.PATCH("/target-update", v1.TaskTargetUpdate)
		router.POST("/task-cluster", v1.TaskClusterCreate)
	}
	return router
}

// 任务路由
func InitTaskRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("/task")
	initGP := router.Group("/init").Use(authMiddleware.MiddlewareFunc())
	{
		initGP.POST("/route-create", v1.TaskRouteCreate)
	}

	return router
}

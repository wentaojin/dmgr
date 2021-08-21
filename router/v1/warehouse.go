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
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/wentaojin/dmgr/service"

	"github.com/wentaojin/dmgr/pkg/dmgrutil"
	"github.com/wentaojin/dmgr/request"

	"github.com/wentaojin/dmgr/response"

	"github.com/gin-gonic/gin"
)

// 上传离线包
func FileChunkUpload(c *gin.Context) {
	var req request.PackageReqStruct
	if response.FailWithMsg(c, c.ShouldBind(&req)) {
		return
	}

	// 确保文件未缓存（例如在 iOS 设备上发生的情况）
	c.Header("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Cache-Control", "post-check=0, pre-check=0")
	c.Header("Pragma", "no-cache")

	// 目标目录是否存在
	pkgDir := filepath.Join(req.PackagePath, dmgrutil.DirSoft)
	if exist, _ := dmgrutil.PathExists(pkgDir); !exist {
		if response.FailWithMsg(c, os.MkdirAll(pkgDir, 0750)) {
			return
		}
	}

	// 定义需要前端上传的文件字段名
	file, err := c.FormFile("file")
	if response.FailWithMsg(c, err) {
		return
	}

	// 文件上传
	filePath := filepath.Join(pkgDir, req.PackageName)
	if response.FailWithMsg(c, c.SaveUploadedFile(file, filePath)) {
		return
	}

	// 更新元数据仓库
	s := service.NewMysqlService()

	req.PackagePath = pkgDir
	if response.FailWithMsg(c, s.AddPackage(req)) {
		return
	}
	response.SuccessWithoutData(c)
}

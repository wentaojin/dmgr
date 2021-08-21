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
package middleware

import (
	"fmt"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/wentaojin/dmgr/response"

	"github.com/wentaojin/dmgr/request"
	"github.com/wentaojin/dmgr/service"

	"github.com/gin-gonic/gin"
	"github.com/wentaojin/dmgr/pkg/dmgrutil"
)

func JwtAuth(file *dmgrutil.MiddlewareConfig) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           file.JwtRealm,                                      // jwt标识 - 中间件名称
		Key:             []byte(file.JwtKey),                                // 服务端密钥
		Timeout:         time.Hour * time.Duration(file.JwtTimeout),         // token过期时间
		MaxRefresh:      time.Hour * time.Duration(file.JwtMaxRefresh),      // token更新时间
		PayloadFunc:     payloadFunc,                                        // 有效载荷处理 - 登录时调用，可将载荷添加到token中
		IdentityHandler: identityHandler,                                    // 解析 Claims - 验证登录状态
		Authenticator:   login,                                              // 验证登录 - 校验 token 正确性并处理登录逻辑
		Authorizator:    authorizator,                                       // 用户登录成功之后鉴权处理 - 判断用户是否有权限访问
		Unauthorized:    unauthorized,                                       // 登录失败处理
		LoginResponse:   loginResponse,                                      // 登录成功后的响应
		LogoutResponse:  logoutResponse,                                     // 登出后的响应
		TokenLookup:     "header: Authorization, query: token, cookie: jwt", // 自动在这几个地方寻找请求中的token
		TokenHeadName:   "Bearer",                                           // header名称
		// TimeFunc提供当前时间。您可以覆盖它以使用其他时间值。这对于测试或服务器使用不同于令牌的时区很有用
		TimeFunc: time.Now,
	})
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*response.UserRespStruct); ok {
		return jwt.MapClaims{
			jwt.IdentityKey: v.Username,
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	// 此处返回值类型与 payloadFunc和 authorizator 的 data 类型必须一致, 否则会导致授权失败还不容易找到原因
	return &response.UserRespStruct{
		Username: claims[jwt.IdentityKey].(string),
	}
}

func login(c *gin.Context) (interface{}, error) {
	var (
		req request.RegisterAndLoginReqStruct
	)
	// 请求json绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	u := &response.UserRespStruct{
		Username: req.Username,
		Password: req.Password,
	}

	// 创建服务
	s := service.NewMysqlService()
	// 登录校验
	user, err := s.LoginCheck(u)
	if err != nil {
		return nil, err
	}
	// 返回用户登录信息, payloadFunc/authorizator 函数会使用到
	return user, nil
}

func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*response.UserRespStruct); ok {
		// 将用户保存到 context, api调用时取数据方便
		c.Set("user", v)
		return true
	}
	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	response.FailWithMsg(c, fmt.Errorf(message))
}

func loginResponse(c *gin.Context, code int, token string, expires time.Time) {
	response.SuccessWithData(c, &response.LoginRespStruct{
		Token:     token,
		ExpiresAt: expires,
	})

}

func logoutResponse(c *gin.Context, code int) {
	response.SuccessWithoutData(c)
}

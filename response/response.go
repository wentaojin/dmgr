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
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"google.golang.org/grpc/codes"

	expect "github.com/google/goexpect"

	"golang.org/x/crypto/ssh"

	"github.com/wentaojin/dmgr/pkg/dmgrutil"
	"go.uber.org/zap"
	kh "golang.org/x/crypto/ssh/knownhosts"

	"github.com/gin-gonic/gin"
)

// http 请求响应封装
type Resp struct {
	Code int         `json:"code"` // 错误代码代码
	Data interface{} `json:"data"` // 数据内容
	Msg  string      `json:"msg"`  // 消息提示
}

func Result(code int, msg string, data interface{}) *Resp {
	// 结果以panic异常的形式抛出, 交由异常处理中间件处理
	return &Resp{
		Code: code,
		Data: data,
		Msg:  msg,
	}
}

func SuccessWithoutData(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, Result(Ok, CustomError[Ok], map[string]interface{}{}))
}

func SuccessWithData(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, Result(Ok, CustomError[Ok], data))
}

func FailWithMsg(c *gin.Context, err error) bool {
	if err != nil {
		dmgrutil.Logger.Error("gin api request error", zap.Error(err))
		c.AbortWithStatusJSON(200, Result(NotOk, c.Error(err).Error(), map[string]interface{}{}))
		return true //表示有错误，调用者应该返回
	}
	return false // 没有错误，调用者可以继续
}

func FailWithCode(c *gin.Context, code int) {
	// 查找给定的错误码存在对应的错误信息, 默认使用 NotOk
	msg := CustomError[NotOk]
	if val, ok := CustomError[code]; ok {
		msg = val
	}
	c.AbortWithStatusJSON(http.StatusOK, Result(code, msg, map[string]interface{}{}))
}

// 写入 json 响应返回值
func SuccessWithJSON(c *gin.Context, code int, resp interface{}) {
	// 调用 gin 写入 json
	c.JSON(code, resp)
}

// GinRecovery recover 掉项目可能出现的 panic，并使用 zap 记录相关日志
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					if FailWithMsg(c, fmt.Errorf(
						"HTTP Request [%v], URL Request Path [%v], Error: [%v]", string(httpRequest), c.Request.URL.Path, err)) {
						return
					}
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				FailWithCode(c, InternalServerError)
			}
		}()
		c.Next()
	}
}

// 测试 SSH 认证连通性
func (m *MachineRespStruct) SshAuthTest(edPrivatePath string) (bool, error) {
	key, err := ioutil.ReadFile(edPrivatePath)
	if err != nil {
		return false, fmt.Errorf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return false, fmt.Errorf("unable to parse private key: %v", err)
	}

	knowHostsPath := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	if !dmgrutil.IsExist(knowHostsPath) {
		knowHostsFile, err := os.Create(knowHostsPath)
		if err != nil {
			return false, err
		}
		defer knowHostsFile.Close()
	}

	hostKeyCallback, err := kh.New(knowHostsPath)
	if err != nil {
		return false, fmt.Errorf("could not create hostkeycallback function: %v", err)
	}

	config := &ssh.ClientConfig{
		User: m.SshUser,
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
	}
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", m.SshHost, m.SshPort), config)
	if err != nil {
		// 当 know_hosts 文件为空时，忽略错误
		// ssh: handshake failed: knownhosts: key is unknown
		return false, nil
	} else {
		client.Close()
	}
	return true, nil
}

func (m *MachineRespStruct) SshCopyID(executeTimeout uint64) error {
	e, _, err := expect.Spawn(fmt.Sprintf("ssh-copy-id %s@%s -p %d", m.SshUser, m.SshHost, m.SshPort), time.Second*time.Duration(executeTimeout))
	if err != nil {
		return err
	}
	defer e.Close()

	caser := []expect.Caser{
		&expect.BCase{R: "password", T: func() (tag expect.Tag, status *expect.Status) {
			_ = e.Send(m.SshPassword + "\n")
			return expect.OKTag, expect.NewStatus(codes.OK, "")
		}},
		&expect.BCase{R: "yes/no", S: "yes\n"},
	}

	for {
		output, _, _, err := e.ExpectSwitchCase(caser, time.Second*time.Duration(executeTimeout))

		if strings.Contains(output, "added") || strings.Contains(output, "exist") {
			break
		}
		if err != nil {
			if strings.Contains(output, "known_hosts") {
				cmd := fmt.Sprintf("sed -i '/%s/d' %s", m.SshHost, filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
				if _, execError := exec.Command("bash", "-c", cmd).Output(); execError != nil {
					return fmt.Errorf("sed file [%v] failed: %v", filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"), err)
				}
				e, _, _ = expect.Spawn(fmt.Sprintf("ssh-copy-id %s@%s -p %d", m.SshUser, m.SshHost, m.SshPort), time.Second*time.Duration(executeTimeout))
				continue
			}
		}
	}
	return nil
}

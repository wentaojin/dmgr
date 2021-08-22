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
package timeout

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/wentaojin/dmgr/pkg/cluster/executor"

	"github.com/wentaojin/dmgr/response"

	"github.com/gin-gonic/gin"
)

// thanks to https://github.com/justlazydog/gin-timeout/
type Option func(*Timeout)

type Timeout struct {
	timeout time.Duration
	code    int
	msg     string
}

func WithTimeout(timeout time.Duration) Option {
	return func(t *Timeout) {
		t.timeout = timeout
	}
}

func WithResponseCode(code int) Option {
	return func(t *Timeout) {
		t.code = code
	}
}

func WithResponseMsg(msg string) Option {
	return func(t *Timeout) {
		t.msg = msg
	}
}

func New(opts ...Option) gin.HandlerFunc {
	t := &Timeout{
		timeout: executor.DefaultConnectTimeout * time.Second,
		code:    response.RequestTimeout,
		msg:     response.CustomError[response.RequestTimeout],
	}

	for _, opt := range opts {
		opt(t)
	}

	return func(c *gin.Context) {
		ctx, cancelCtx := context.WithTimeout(c.Request.Context(), t.timeout)
		defer cancelCtx()
		c.Request = c.Request.WithContext(ctx)

		tw := &router.timeoutWriter{
			ResponseWriter: c.Writer,
			h:              make(http.Header),
		}
		c.Writer = tw

		done := make(chan struct{})
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			close(done)
		}()

		select {
		case p := <-panicChan:
			panic(p)
		case <-ctx.Done():
			tw.mu.Lock()
			defer tw.mu.Unlock()
			tw.ResponseWriter.WriteHeader(t.code)
			io.WriteString(tw.ResponseWriter, t.msg)
			tw.timedOut = true
			c.Abort()
		case <-done:
			tw.mu.Lock()
			defer tw.mu.Unlock()
			dst := c.Writer.Header()
			for k, vv := range tw.h {
				dst[k] = vv
			}
			if !tw.wroteHeader {
				tw.code = http.StatusOK
			}
			tw.ResponseWriter.WriteHeader(tw.code)
			tw.ResponseWriter.Write(tw.wbuf.Bytes())
		}
	}
}

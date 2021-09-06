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
package api

import (
	"bytes"
	"fmt"
	"io"
)

// 任务管理请求接口列表
const (
	sourceAPI = "/api/v1/sources"
	taskAPI   = "/api/v1/tasks"
)

// 获取所有数据源信息
func GetSourcesALL(dmMasterHttpUrl string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", dmMasterHttpUrl, sourceAPI)
	client := NewHTTPClient(DmMasterApiTimeout, nil)
	return client.Get(url)
}

// 获取数据源状态
func GetSourceStatusBySourceName(dmMasterHttpUrl, sourceName string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s/status", dmMasterHttpUrl, sourceAPI, sourceName)
	client := NewHTTPClient(DmMasterApiTimeout, nil)
	return client.Get(url)
}

// 删除某个上游
func DeleteSourceBySourceName(dmMasterHttpUrl, sourceName string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", dmMasterHttpUrl, sourceAPI, sourceName)
	client := NewHTTPClient(DmMasterApiTimeout, nil)
	body, statusCode, err := client.Delete(url, nil)

	if statusCode == 400 {
		return body, fmt.Errorf("source name [%v] to delete failed: %v", sourceName, string(body))
	}
	if bytes.Contains(body, []byte("not exists")) {
		return body, fmt.Errorf("source name [%v] to delete does not exist, ignore delete", sourceName)
	}
	if err != nil {
		return body, err
	}
	return body, nil
}

// 启动某个数据源同步任务
func StartSourceBySourceName(dmMasterHttpUrl, sourceName string, body io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s/start-relay", dmMasterHttpUrl, sourceAPI, sourceName)
	client := NewHTTPClient(DmMasterApiTimeout, nil)
	resp, statusCode, err := client.Patch(url, body)

	if statusCode == 400 {
		return resp, fmt.Errorf("source name [%v] to patch failed: %v", sourceName, string(resp))
	}

	if bytes.Contains(resp, []byte("not exists")) {
		return resp, fmt.Errorf("source name [%v] to patch does not exist, ignore start", sourceName)
	}
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// 停止某个数据源同步任务
func StopSourceBySourceName(dmMasterHttpUrl, sourceName string, body io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s/stop-relay", dmMasterHttpUrl, sourceAPI, sourceName)
	client := NewHTTPClient(DmMasterApiTimeout, nil)
	resp, statusCode, err := client.Patch(url, body)

	if statusCode == 400 {
		return resp, fmt.Errorf("source name [%v] to patch failed: %v", sourceName, string(resp))
	}

	if bytes.Contains(resp, []byte("not exists")) {
		return resp, fmt.Errorf("source name [%v] to patch does not exist, ignore stop", sourceName)
	}
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// 创建某个数据源
func CreateSource(dmMasterHttpUrl string, body io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", dmMasterHttpUrl, sourceAPI)
	client := NewHTTPClient(DmMasterApiTimeout, nil)
	return client.Post(url, body)
}

// 任务创建及启动
func StartTaskMigration(dmMasterHttpUrl string, body io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", dmMasterHttpUrl, taskAPI)
	client := NewHTTPClient(DmMasterApiTimeout, nil)
	return client.Post(url, body)
}

// 获取所有任务信息
func GetTaskALL(dmMasterHttpUrl string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", dmMasterHttpUrl, taskAPI)
	client := NewHTTPClient(DmMasterApiTimeout, nil)
	return client.Get(url)
}

// 获取任务状态
func GetTaskStatusByTaskName(dmMasterHttpUrl, taskName string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s/status", dmMasterHttpUrl, taskAPI, taskName)
	client := NewHTTPClient(DmMasterApiTimeout, nil)
	return client.Get(url)
}

// 删除某个任务
func DeleteTaskByTaskName(dmMasterHttpUrl, taskName string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", dmMasterHttpUrl, taskAPI, taskName)
	client := NewHTTPClient(DmMasterApiTimeout, nil)
	body, statusCode, err := client.Delete(url, nil)

	if statusCode == 400 {
		return body, fmt.Errorf("source name [%v] to delete failed: %v", taskName, string(body))
	}
	if bytes.Contains(body, []byte("not exists")) {
		return body, fmt.Errorf("source name [%v] to delete does not exist, ignore delete", taskName)
	}

	if err != nil {
		return body, err
	}
	return body, nil
}

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
package dmgrutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// 初始化组件配置文件、脚本等文件缓存目录以及 SSH 认证存放目录
func InitComponentCacheAndSSHDir(clusterPath, clusterName string) error {
	if exist, _ := PathExists(AbsClusterCacheDir(clusterPath, clusterName)); !exist {
		if err := os.MkdirAll(AbsClusterCacheDir(clusterPath, clusterName), 0750); err != nil {
			return err
		}
	}
	// 初始化 SSH 存放认证目录
	if exist, _ := PathExists(AbsClusterSSHDir(clusterPath, clusterName)); !exist {
		if err := os.MkdirAll(AbsClusterSSHDir(clusterPath, clusterName), 0750); err != nil {
			return err
		}
	}
	return nil
}

// 集群解压存放目录
// {cluster_path}/cluster/{cluster_name}/{cluster_version}
func AbsClusterUntarDir(clusterPath, clusterName string) string {
	return filepath.Join(clusterPath, DirCluster, clusterName)
}

// 集群组件文件名
// {cluster_path}/cluster/{cluster_name}/{cluster_version}/grafana.tar.gz
func AbsClusterGrafanaComponent(clusterPath, clusterName, clusterVersion, grafanaPKG string) string {
	return fmt.Sprintf("%s/%s", filepath.Join(AbsClusterUntarDir(clusterPath, clusterName), clusterVersion), grafanaPKG)
}

// 集群目录
func AbsClusterDeployDir(deployDir, instanceName string) string {
	return filepath.Join(deployDir, instanceName)
}

// 集群部署 Bin 目录
func AbsClusterBinDir(deployDir, instanceName string) string {
	return filepath.Join(deployDir, instanceName, DirBin)
}

// 集群部署 Conf 目录
func AbsClusterConfDir(deployDir, instanceName string) string {
	return filepath.Join(deployDir, instanceName, DirConf)
}

// 集群部署 Dashboard 目录
func AbsClusterDataboardDir(deployDir, instanceName string) string {
	return filepath.Join(deployDir, instanceName, DirGrafanaDashboard)
}

// 集群部署 Datasource 目录
func AbsClusterDatasourceDir(deployDir, instanceName string) string {
	return filepath.Join(deployDir, instanceName, DirGrafanaDatasource)
}

// 集群部署 Datasource 目录
func AbsClusterTempSystemdDir(deployDir string) string {
	return filepath.Join(deployDir)
}

func AbsClusterSystemdDir() string {
	return filepath.Join(DirSystemd)
}

// 集群部署 Script 目录
func AbsClusterScriptDir(deployDir, instanceName string) string {
	return filepath.Join(deployDir, instanceName, DirScript)
}

// 集群部署 Data 目录
func AbsClusterDataDir(deployDir, dataDir, instanceName string) string {
	if deployDir == dataDir || dataDir == "" {
		return filepath.Join(deployDir, instanceName, DirData)
	}
	return filepath.Join(dataDir, instanceName)
}

// 集群部署 Log 目录
func AbsClusterLogDir(deployDir, logDir, instanceName string) string {
	if deployDir == logDir || logDir == "" {
		return filepath.Join(deployDir, instanceName, DirLog)
	}
	return filepath.Join(logDir, instanceName)
}

// 集群模板文件缓存 Cache 目录
func AbsClusterCacheDir(clusterPath, clusterName string) string {
	return filepath.Join(clusterPath, DirCluster, clusterName, DirCache)
}

// 集群模板文件 SSH 认证存放目录
func AbsClusterSSHDir(clusterPath, clusterName string) string {
	return filepath.Join(clusterPath, DirCluster, clusterName, DirSSH)
}

// 集群压缩包文件位置
func AbsUntarConfDir(clusterPath, clusterName, clusterVersion, fileName string) string {
	return filepath.Join(clusterPath, DirCluster, clusterName, clusterVersion, fileName)
}

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
	"path/filepath"
)

// 集群解压存放目录
// {cluster_path}/cluster/{cluster_name}/{cluster_version}
func AbsClusterUntarDir(clusterPath, clusterName string) string {
	return filepath.Join(clusterPath, DirCluster, clusterName)
}

// 集群组件文件名
// {cluster_path}/cluster/{cluster_name}/{cluster_version}/{componentName}.tar.gz
func AbsClusterComponent(clusterPath, clusterName, clusterVersion, componentName string) string {
	return filepath.Join(AbsClusterUntarDir(clusterPath, clusterName), clusterVersion, fmt.Sprintf("%s.tar.gz", componentName))
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
func AbsClusterSystemdDir() string {
	return filepath.Join(DirSystemd)
}

// 集群部署 Script 目录
func AbsClusterScriptDir(deployDir, instanceName string) string {
	return filepath.Join(deployDir, instanceName, DirScript)
}

// 集群部署 Data 目录
func AbsClusterDataDir(dataDir, instanceName string) string {
	return filepath.Join(dataDir, instanceName, DirData)
}

// 集群部署 Log 目录
func AbsClusterLogDir(logDir, instanceName string) string {
	return filepath.Join(logDir, instanceName, DirLog)
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
func AbsUntarConfDir(clusterPath, clusterVersion, fileName string) string {
	return filepath.Join(clusterPath, clusterVersion, fileName)
}

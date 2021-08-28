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
package service

import (
	"fmt"

	"github.com/wentaojin/dmgr/response"
)

func (s *MysqlService) GetClusterStatus(clusterStatus string) ([]response.ClusterMetaRespStruct, error) {
	var cm []response.ClusterMetaRespStruct
	if clusterStatus == "" {
		if err := s.Engine.Select(&cm, `SELECT cluster_name,cluster_user,cluster_version,cluster_path,cluster_status,admin_user,admin_password FROM cluster_meta`); err != nil {
			return cm, err
		}
		return cm, nil
	}
	if err := s.Engine.Select(&cm, `SELECT cluster_name,cluster_user,cluster_version,cluster_path,cluster_status,admin_user,admin_password FROM cluster_meta WHERE cluster_status = ?`, clusterStatus); err != nil {
		return cm, err
	}
	return cm, nil
}

func (s *MysqlService) UpdateClusterMetaStatus(clusterName, clusterStatus string) error {
	if _, err := s.Engine.Exec(`UPDATE cluster_meta SET cluster_status = ? WHERE cluster_name = ?`, clusterStatus, clusterName); err != nil {
		return err
	}
	return nil
}

func (s *MysqlService) ValidClusterNameIsExist(clusterName string) (bool, error) {
	if clusterName != "" {
		var ct int
		if err := s.Engine.Get(&ct, "SELECT count(1) FROM cluster_meta WHERE cluster_name = ?", clusterName); err != nil {
			return false, err
		}
		if ct == 1 {
			return true, nil
		}
		return false, nil
	}
	return false, fmt.Errorf("cluster_name cannot be null")
}

func (s *MysqlService) GetClusterMeta(clusterName string) (response.ClusterMetaRespStruct, error) {
	var cm response.ClusterMetaRespStruct
	if err := s.Engine.Get(&cm, `SELECT cluster_name,cluster_user,cluster_version,cluster_path,cluster_status,admin_user,admin_password FROM cluster_meta WHERE cluster_name = ?`, clusterName); err != nil {
		return cm, err
	}
	return cm, nil
}

func (s *MysqlService) UpdateClusterVersion(clusterName, clusterVersion string) error {
	if _, err := s.Engine.Exec(`UPDATE cluster_meta SET cluster_version = ? WHERE cluster_name = ?`, clusterVersion, clusterName); err != nil {
		return err
	}
	return nil
}

func (s *MysqlService) UpdateGrafanaUserAndPassword(clusterName, adminUser, adminPassword string) error {
	if _, err := s.Engine.Exec(`UPDATE cluster_meta SET admin_user = ? AND admin_password = ? WHERE cluster_name = ?`, adminUser, adminPassword, clusterName); err != nil {
		return err
	}
	return nil
}

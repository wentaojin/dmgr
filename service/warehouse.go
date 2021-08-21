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
	"github.com/wentaojin/dmgr/request"
	"github.com/wentaojin/dmgr/response"
)

func (s *MysqlService) ValidClusterVersionPackageIsExist(clusterVersion string) (response.WarehouseRespStruct, error) {
	var pkg response.WarehouseRespStruct
	if err := s.Engine.Get(&pkg, "SELECT cluster_version,package_name,package_path FROM warehouse WHERE cluster_version = ?", clusterVersion); err != nil {
		return pkg, err
	}
	return pkg, nil
}

func (s *MysqlService) AddPackage(pkg request.PackageReqStruct) error {
	if _, err := s.Engine.NamedExec(`INSERT INTO warehouse (cluster_version, package_name, package_path) values (:cluster_version, :package_name, :package_path)`,
		map[string]interface{}{
			"cluster_version": pkg.ClusterVersion,
			"package_name":    pkg.PackageName,
			"package_path":    pkg.PackagePath,
		}); err != nil {
		return err
	}
	return nil
}

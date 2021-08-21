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
	"github.com/jmoiron/sqlx"
	"github.com/wentaojin/dmgr/request"
	"github.com/wentaojin/dmgr/response"
)

func (s *MysqlService) AddMachine(machine request.MachineReqStruct) error {
	if _, err := s.Engine.NamedExec(`INSERT INTO machine (ssh_host, ssh_user, ssh_port, ssh_password) values (:ssh_host, :ssh_user, :ssh_port, :ssh_password)`,
		map[string]interface{}{
			"ssh_host":     machine.SshHost,
			"ssh_user":     machine.SshUser,
			"ssh_port":     machine.SshPort,
			"ssh_password": machine.SshPassword,
		}); err != nil {
		return err
	}
	return nil
}

func (s *MysqlService) GetMachineList(machineHosts []string) ([]response.MachineRespStruct, error) {
	var machineList []response.MachineRespStruct
	query, args, err := sqlx.In(`SELECT
	ssh_host,
	ssh_user,
	ssh_password,
	ssh_port 
FROM
    machine WHERE ssh_host IN (?)`, machineHosts)
	if err != nil {
		return machineList, err
	}
	query = s.Engine.Rebind(query)
	if err := s.Engine.Select(&machineList, query, args...); err != nil {
		return machineList, err
	}
	return machineList, nil
}

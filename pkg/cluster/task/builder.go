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
package task

import (
	"github.com/wentaojin/dmgr/response"
)

// Builder 任务 build
type Builder struct {
	tasks []Task
}

// NewBuilder 返回一个 *Builder 实例
func NewBuilder() *Builder {
	return &Builder{}
}

// RootSSH 将 RootSSH 任务附加到当前任务集合
func (b *Builder) RootSSH(
	host string,
	port uint64,
	user, password, keyFile, passphrase string,
	connectTimeout, executeTimeout uint64) *Builder {
	b.tasks = append(b.tasks, &RootSSH{
		host:           host,
		port:           port,
		user:           user,
		password:       password,
		keyFile:        keyFile,
		passphrase:     passphrase,
		connectTimeout: connectTimeout,
		executeTimeout: executeTimeout,
	})
	return b
}

// UserSSH 将 UserSSH 任务附加到当前任务集合
func (b *Builder) UserSSH(
	host string, port uint64, clusterUser string, connectTimeout, executeTimeout uint64) *Builder {
	b.tasks = append(b.tasks, &UserSSH{
		host:           host,
		port:           port,
		clusterUser:    clusterUser,
		connectTimeout: connectTimeout,
		executeTimeout: executeTimeout,
	})
	return b
}

// StartInstance 将 StartInstance 任务附加到当前任务集合
func (b *Builder) StartInstance(
	host string,
	servicePort uint64,
	instanceName string,
	logDir string,
	serviceName string,
	executeTimeout uint64) *Builder {
	b.tasks = append(b.tasks, &StartInstance{
		host:           host,
		servicePort:    servicePort,
		instanceName:   instanceName,
		logDir:         logDir,
		serviceName:    serviceName,
		executeTimeout: executeTimeout,
	})
	return b
}

// StopInstance 将 StopInstance 任务附加到当前任务集合
func (b *Builder) StopInstance(
	host string,
	servicePort uint64,
	instanceName string,
	logDir string,
	serviceName string,
	executeTimeout uint64) *Builder {
	b.tasks = append(b.tasks, &StopInstance{
		host:           host,
		servicePort:    servicePort,
		instanceName:   instanceName,
		logDir:         logDir,
		serviceName:    serviceName,
		executeTimeout: executeTimeout,
	})
	return b
}

// DestroyInstance 将 DestroyInstance 任务附加到当前任务集合
func (b *Builder) DestroyInstance(
	host string,
	servicePort uint64,
	componentName string,
	instanceName string,
	deployDir string,
	dataDir string,
	logDir string,
	executeTimeout uint64) *Builder {
	b.tasks = append(b.tasks, &DestroyInstance{
		host:           host,
		servicePort:    servicePort,
		componentName:  componentName,
		instanceName:   instanceName,
		deployDir:      deployDir,
		dataDir:        dataDir,
		logDir:         logDir,
		executeTimeout: executeTimeout,
	})
	return b
}

// EnableInstance 将 StartInstance 任务附加到当前任务集合
func (b *Builder) EnableInstance(
	host string,
	servicePort uint64,
	instanceName string,
	logDir string,
	serviceName string,
	executeTimeout uint64, isEnable bool) *Builder {
	b.tasks = append(b.tasks, &EnableInstance{
		host:           host,
		servicePort:    servicePort,
		instanceName:   instanceName,
		logDir:         logDir,
		serviceName:    serviceName,
		executeTimeout: executeTimeout,
		isEnable:       isEnable,
	})
	return b
}

// CopyComponent 将 CopyComponent 任务附加到当前任务集合
func (b *Builder) CopyComponent(clusterName, componentName string,
	clusterVersion string,
	srcPath, dstHost, dstPath string,
) *Builder {
	b.tasks = append(b.tasks, &CopyComponent{
		clusterName:    clusterName,
		componentName:  componentName,
		clusterVersion: clusterVersion,
		srcPath:        srcPath,
		host:           dstHost,
		dstPath:        dstPath,
	})
	return b
}

// CopyFile 将 CopyFile 任务附加到当前任务集合
func (b *Builder) CopyFile(src, dst, fileType, remoteHost string, download bool, limit int) *Builder {
	b.tasks = append(b.tasks, &CopyFile{
		src:        src,
		dst:        dst,
		fileType:   fileType,
		remoteHost: remoteHost,
		download:   download,
		limit:      limit,
	})
	return b
}

// SSHKeyGen 将 SSHKeyGen 任务附加到当前任务集合
func (b *Builder) SSHKeyGen(homeSshDir string,
	executeTimeout uint64) *Builder {
	b.tasks = append(b.tasks, &SSHKeyGen{
		homeSshDir:     homeSshDir,
		executeTimeout: executeTimeout,
	})
	return b
}

// SSHKeyCopy 将 SSHKeyCopy 任务附加到当前任务集合
func (b *Builder) SSHKeyCopy(homeSshDir, clusterSshDir string,
	hosts []response.MachineRespStruct,
	executeTimeout uint64,
	workerThreads int) *Builder {
	b.tasks = append(b.tasks, &SSHKeyCopy{
		homeSshDir:     homeSshDir,
		clusterSshDir:  clusterSshDir,
		hosts:          hosts,
		executeTimeout: executeTimeout,
		workerThreads:  workerThreads,
	})
	return b
}

// SSHKeySet 将 SSHKeySet 任务附加到当前任务集合
func (b *Builder) SSHKeySet(privKeyPath, pubKeyPath string) *Builder {
	b.tasks = append(b.tasks, &SSHKeySet{
		privateKeyPath: privKeyPath,
		publicKeyPath:  pubKeyPath,
	})
	return b
}

// EnvInit 将 EnvInit 任务附加到当前任务集合
func (b *Builder) EnvInit(host, clusterUser string, userGroup string, skipCreateUser bool) *Builder {
	b.tasks = append(b.tasks, &EnvInit{
		host:           host,
		clusterUser:    clusterUser,
		userGroup:      userGroup,
		skipCreateUser: skipCreateUser,
	})
	return b
}

// Mkdir 将 Mkdir 任务附加到当前任务集合
func (b *Builder) Mkdir(user, host string, dirs ...string) *Builder {
	b.tasks = append(b.tasks, &Mkdir{
		user: user,
		host: host,
		dirs: dirs,
	})
	return b
}

// Serial 将任务附加到队列的尾部
func (b *Builder) Serial(prefix string, tasks ...Task) *Builder {
	if len(tasks) > 0 {
		b.tasks = append(b.tasks, tasks...)
	}
	return b
}

// Parallel 将并行任务附加到当前任务集合
func (b *Builder) Parallel(prefix string, ignoreError bool, tasks ...Task) *Builder {
	if len(tasks) > 0 {
		b.tasks = append(b.tasks, &Parallel{ignoreError: ignoreError, hideDetailDisplay: false, inner: tasks})
	}
	return b
}

// Build 返回一个任务，其中包含由先前操作附加的所有任务
func (b *Builder) BuildTask() Task {
	// Serial handles event internally. So the following 3 lines are commented out.
	// if len(b.tasks) == 1 {
	//  return b.tasks[0]
	// }
	return &Serial{ignoreError: false, hideDetailDisplay: false, inner: b.tasks}
}

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
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	perrs "github.com/pingcap/errors"
	"github.com/wentaojin/dmgr/pkg/cluster/api/dmpb"
	"github.com/wentaojin/dmgr/pkg/dmgrutil"

	"github.com/golang/protobuf/jsonpb"
	"go.uber.org/zap"
)

var (
	dmMembersURI = "apis/v1alpha1/members"

	DefaultRetryOpt = &dmgrutil.RetryOption{
		Delay:   time.Second * 5,
		Timeout: time.Second * 60,
	}

	// DM Master 请求接口耗时
	DmMasterApiTimeout = 10 * time.Second
)

// DMMasterClient is an HTTP client of the dm-master server
type DMMasterClient struct {
	addrs      []string
	tlsEnabled bool
	httpClient *HTTPClient
}

// NewDMMasterClient returns a new PDClient
func NewDMMasterClient(addrs []string, timeout time.Duration, tlsConfig *tls.Config) *DMMasterClient {
	enableTLS := false
	if tlsConfig != nil {
		enableTLS = true
	}

	return &DMMasterClient{
		addrs:      addrs,
		tlsEnabled: enableTLS,
		httpClient: NewHTTPClient(timeout, tlsConfig),
	}
}

// GetURL builds the the client URL of DMClient
func (dm *DMMasterClient) GetURL(addr string) string {
	httpPrefix := "http"
	if dm.tlsEnabled {
		httpPrefix = "https"
	}
	return fmt.Sprintf("%s://%s", httpPrefix, addr)
}

func (dm *DMMasterClient) getEndpoints(cmd string) (endpoints []string) {
	for _, addr := range dm.addrs {
		endpoint := fmt.Sprintf("%s/%s", dm.GetURL(addr), cmd)
		endpoints = append(endpoints, endpoint)
	}

	return
}

func (dm *DMMasterClient) getMember(endpoints []string) (*dmpb.ListMemberResponse, error) {
	resp := &dmpb.ListMemberResponse{}
	_, err := tryURLs(endpoints, func(endpoint string) ([]byte, error) {
		body, err := dm.httpClient.Get(endpoint)
		if err != nil {
			return body, err
		}

		err = jsonpb.Unmarshal(strings.NewReader(string(body)), resp)

		if err != nil {
			return body, err
		}

		if !resp.Result {
			return body, errors.New("dm-master get members failed: " + resp.Msg)
		}

		return body, nil
	})
	return resp, err
}

func (dm *DMMasterClient) deleteMember(endpoints []string) (*dmpb.OfflineMemberResponse, error) {
	resp := &dmpb.OfflineMemberResponse{}
	_, err := tryURLs(endpoints, func(endpoint string) ([]byte, error) {
		body, statusCode, err := dm.httpClient.Delete(endpoint, nil)

		if statusCode == 404 || bytes.Contains(body, []byte("not exists")) {
			dmgrutil.Logger.Debug("member to offline does not exist, ignore.")
			return body, nil
		}
		if err != nil {
			return body, err
		}

		err = jsonpb.Unmarshal(strings.NewReader(string(body)), resp)

		if err != nil {
			return body, err
		}

		if !resp.Result {
			return body, errors.New("dm-master offline member failed: " + resp.Msg)
		}

		return body, nil
	})
	return resp, err
}

// GetMaster returns the dm master leader
// returns isFound, isActive, isLeader, error
func (dm *DMMasterClient) GetMaster(name string) (isFound bool, isActive bool, isLeader bool, err error) {
	query := "?leader=true&master=true&names=" + name
	endpoints := dm.getEndpoints(dmMembersURI + query)
	memberResp, err := dm.getMember(endpoints)

	if err != nil {
		dmgrutil.Logger.Error("get dm master status failed", zap.Error(err))
		return false, false, false, err
	}

	for _, member := range memberResp.GetMembers() {
		if leader := member.GetLeader(); leader != nil {
			if leader.GetName() == name {
				isFound = true
				isLeader = true
			}
		} else if masters := member.GetMaster(); masters != nil {
			for _, master := range masters.GetMasters() {
				if master.GetName() == name {
					isFound = true
					isActive = master.GetAlive()
				}
			}
		}
	}
	return
}

// GetWorker returns the dm worker status
// returns (worker stage, error). If worker stage is "", that means this worker is in cluster
func (dm *DMMasterClient) GetWorker(name string) (string, error) {
	query := "?worker=true&names=" + name
	endpoints := dm.getEndpoints(dmMembersURI + query)
	memberResp, err := dm.getMember(endpoints)

	if err != nil {
		dmgrutil.Logger.Error("get dm worker status failed", zap.Error(err))
		return "", err
	}

	stage := ""
	for _, member := range memberResp.Members {
		if workers := member.GetWorker(); workers != nil {
			for _, worker := range workers.GetWorkers() {
				if worker.GetName() == name {
					stage = worker.GetStage()
				}
			}
		}
	}
	if len(stage) > 0 {
		stage = strings.ToUpper(stage[0:1]) + stage[1:]
	}

	return stage, nil
}

// GetLeader gets leader of dm cluster
func (dm *DMMasterClient) GetLeader(retryOpt *dmgrutil.RetryOption) (string, string, error) {
	query := "?leader=true"
	endpoints := dm.getEndpoints(dmMembersURI + query)

	if retryOpt == nil {
		retryOpt = DefaultRetryOpt
	}

	var (
		memberResp *dmpb.ListMemberResponse
		err        error
	)

	if err := dmgrutil.Retry(func() error {
		memberResp, err = dm.getMember(endpoints)
		return err
	}, *retryOpt); err != nil {
		return "", "", err
	}

	leaderName := ""
	leaderAddr := ""
	for _, member := range memberResp.Members {
		if leader := member.GetLeader(); leader != nil {
			leaderName = leader.GetName()
			leaderAddr = leader.GetAddr()
		}
	}
	return leaderName, leaderAddr, nil
}

// GetRegisteredMembers gets all registerer members of dm cluster
func (dm *DMMasterClient) GetRegisteredMembers() ([]string, []string, error) {
	query := "?master=true&worker=true"
	endpoints := dm.getEndpoints(dmMembersURI + query)
	memberResp, err := dm.getMember(endpoints)

	var (
		registeredMasters []string
		registeredWorkers []string
	)

	if err != nil {
		dmgrutil.Logger.Error("get dm master status failed", zap.Error(err))
		return registeredMasters, registeredWorkers, err
	}

	for _, member := range memberResp.Members {
		if masters := member.GetMaster(); masters != nil {
			for _, master := range masters.GetMasters() {
				registeredMasters = append(registeredMasters, master.Name)
			}
		} else if workers := member.GetWorker(); workers != nil {
			for _, worker := range workers.GetWorkers() {
				registeredWorkers = append(registeredWorkers, worker.Name)
			}
		}
	}

	return registeredMasters, registeredWorkers, nil
}

// EvictDMMasterLeader evicts the dm master leader
func (dm *DMMasterClient) EvictDMMasterLeader(retryOpt *dmgrutil.RetryOption) error {
	return nil
}

// OfflineMember offlines the member of dm cluster
func (dm *DMMasterClient) OfflineMember(query string, retryOpt *dmgrutil.RetryOption) error {
	endpoints := dm.getEndpoints(dmMembersURI + query)

	if retryOpt == nil {
		retryOpt = DefaultRetryOpt
	}

	if err := dmgrutil.Retry(func() error {
		_, err := dm.deleteMember(endpoints)
		return err
	}, *retryOpt); err != nil {
		return fmt.Errorf("error offline member %s, %v", query, err)
	}
	return nil
}

// OfflineWorker offlines the dm worker
func (dm *DMMasterClient) OfflineWorker(name string, retryOpt *dmgrutil.RetryOption) error {
	query := "/worker/" + name
	return dm.OfflineMember(query, retryOpt)
}

// OfflineMaster offlines the dm master
func (dm *DMMasterClient) OfflineMaster(name string, retryOpt *dmgrutil.RetryOption) error {
	query := "/master/" + name
	return dm.OfflineMember(query, retryOpt)
}

func tryURLs(endpoints []string, f func(endpoint string) ([]byte, error)) ([]byte, error) {
	if len(endpoints) == 0 {
		return nil, errors.New("no endpoint available")
	}
	var err error
	var bytes []byte
	for _, endpoint := range endpoints {
		var u *url.URL
		u, err = url.Parse(endpoint)

		if err != nil {
			return bytes, perrs.AddStack(err)
		}

		endpoint = u.String()

		bytes, err = f(endpoint)
		if err != nil {
			continue
		}
		return bytes, nil
	}
	if len(endpoints) > 1 && err != nil {
		err = perrs.Errorf("no endpoint available, the last err was: %s", err)
	}
	return bytes, err
}

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
package main

import (
	"fmt"

	"github.com/tidwall/gjson"
)

func main() {
	jsonStr := `{
  "source_name": "mysql-replica-01",
  "worker_name": "worker-1",
  "enable_relay": false,
  "relay_status": {
    "master_binlog": "(mysql-bin.000001, 1979)",
    "master_binlog_gtid": "e9a1fc22-ec08-11e9-b2ac-0242ac110003:1-7849",
    "relay_dir": "./sub_dir",
    "relay_binlog_gtid": "e9a1fc22-ec08-11e9-b2ac-0242ac110003:1-7849",
    "relay_catch_up_master": true,
    "stage": "Running"
  }
}`
	workName := gjson.Get(jsonStr, "worker_name|@join")
	fmt.Println(workName)
}

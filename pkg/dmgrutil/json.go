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
	"encoding/json"
	"fmt"
	"reflect"
)

// json interface转为结构体
func JsonI2Struct(str interface{}, obj interface{}) error {
	// 将json interface转为string
	jsonStr, _ := str.(string)
	if err := Json2Struct(jsonStr, obj); err != nil {
		return err
	}
	return nil
}

// json转为结构体
func Json2Struct(str string, obj interface{}) error {
	// 将json转为结构体
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		return fmt.Errorf("[Json2Struct] convert falied: %v", err)
	}
	return nil
}

// 结构体转为json
func Struct2Json(obj interface{}) (string, error) {
	str, err := json.Marshal(obj)
	if err != nil {
		return string(str), fmt.Errorf("[Struct2Json] convert falied: %v", err)
	}
	return string(str), nil
}

// 结构体判断是否 null
func StructIsEmpty(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// 结构体转结构体, json为中间桥梁, struct2必须以指针方式传递, 否则可能获取到空数据
func Struct2StructByJson(struct1 interface{}, struct2 interface{}) error {
	// 转换为响应结构体, 隐藏部分字段
	jsonStr, err := Struct2Json(struct1)
	if err != nil {
		return err
	}
	if err := Json2Struct(jsonStr, struct2); err != nil {
		return err
	}
	return nil
}

// 两结构体比对不同的字段, 不同时将取struct1中的字段返回, json为中间桥梁, update必须以指针方式传递, 否则可能获取到空数据
func CompareDifferenceStructByJson(oldStruct interface{}, newStruct interface{}, update interface{}) error {
	// 通过json先将其转为map集合
	m1 := make(map[string]interface{}, 0)
	m2 := make(map[string]interface{}, 0)
	m3 := make(map[string]interface{}, 0)
	if err := Struct2StructByJson(newStruct, &m1); err != nil {
		return err
	}
	if err := Struct2StructByJson(oldStruct, &m2); err != nil {
		return err
	}
	for k1, v1 := range m1 {
		for k2, v2 := range m2 {
			switch v1.(type) {
			// 复杂结构不做对比
			case map[string]interface{}:
				continue
			}
			// key相同, 值不同
			if k1 == k2 && v1 != v2 {
				m3[k1] = v1
				break
			}
		}
	}
	if err := Struct2StructByJson(m3, &update); err != nil {
		return err
	}
	return nil
}

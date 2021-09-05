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
	"archive/tar"
	"compress/gzip"
	"io"
	"net"
	"os"
	"os/user"
	"path"
	"reflect"
	"strings"
)

// PathExists 判断文件或目录是否已经存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 检查路径是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// 判断结构体是否相等
func IsStructureEqual(newStruct, originStruct interface{}) bool {
	return reflect.DeepEqual(newStruct, originStruct)
}

// 创建目录
func CreateDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, 0755)
		}
		return err
	}
	return nil
}

// 解压 TarGZ 文件
// 解压前需要先检查文件是否存在
func UnCompressTarGz(srcFilePath string, destDirPath string) error {
	// 判断并创建目标目录
	if exist, _ := PathExists(destDirPath); !exist {
		if err := os.MkdirAll(destDirPath, 0750); err != nil {
			return err
		}
	}
	fr, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer fr.Close()

	// Gzip/Tar reader
	gr, err := gzip.NewReader(fr)
	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}
		// 检查是否是文件
		if hdr.Typeflag != tar.TypeDir {
			// 从压缩文件内获取文件
			// 创建文件前创建目录
			if exist, _ := PathExists(destDirPath + "/" + path.Dir(hdr.Name)); !exist {
				if err := os.MkdirAll(destDirPath+"/"+path.Dir(hdr.Name), os.ModePerm); err != nil {
					return err
				}
			}

			// 写数据到文件
			fw, err := os.Create(destDirPath + "/" + hdr.Name)
			if err != nil {
				return err
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 判断字符是否相等
func StringEqualFold(str1, str2 string) bool {
	if strings.EqualFold(str1, str2) {
		return true
	}
	return false
}

// 筛选重复元素
func FilterRepeatElem(str []string) []string {
	var repeatElem []string
	num := make(map[string]bool)
	for _, v := range str {
		if !num[v] {
			num[v] = true
		} else {
			repeatElem = append(repeatElem, v)
		}
	}
	return repeatElem
}

// 判断数组是否存在某个元素
func IsContainElem(array interface{}, value interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}

// 获取本机客户端地址
func GetClientOutBoundIP() (username, ip string, err error) {
	u, err := user.Current()
	if err != nil {
		return
	}
	username = u.Username

	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}

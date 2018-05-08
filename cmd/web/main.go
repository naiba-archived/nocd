/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package main

import (
	"os"
	"path/filepath"

	"github.com/naiba/nocd"
	"github.com/naiba/nocd/router"
)

func init() {
	// initial global settings
	nocd.InitSysConfig("conf/app.ini")
	unzipAssets("resource/", "24", []string{"resource"}, RestoreAssets)
}

func main() {
	router.Start()
}

// 释放资源文件
func unzipAssets(path, ver string, dirs []string, call func(s1, s2 string) error) {
	if _, err := os.Stat(path); err == nil {
		if _, err := os.Stat(path + ver + ".ver"); os.IsNotExist(err) {
			nocd.Logger().Infoln("[" + ver + "]: Delete Old Assets.")
			os.RemoveAll(path)
		} else {
			nocd.Logger().Infoln("[" + ver + "]: Assets File Exists.")
			return
		}
	}
	nocd.Logger().Infoln("[" + ver + "]: Unpkg Assets.")
	isSuccess := true
	for _, dir := range dirs {
		// 解压dir目录到当前目录
		if err := call("./", dir); err != nil {
			isSuccess = false
			break
		}
	}
	if !isSuccess {
		for _, dir := range dirs {
			os.RemoveAll(filepath.Join("./", dir))
		}
	}
}

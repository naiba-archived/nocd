/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package main

import (
	"os"
	"path/filepath"

	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"

	"git.cm/naiba/gocd"
	"git.cm/naiba/gocd/router"
)

func init() {
	// initial global settings
	gocd.Initial("conf/app.ini")
	// initial sentry dsn
	hook, err := logrus_sentry.NewSentryHook(gocd.Conf.Section("third_party").Key("sentry_dsn").String(), []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})
	if err == nil {
		gocd.Log.Hooks.Add(hook)
	} else {
		gocd.Log.Panicln(err)
	}
	unzipAssets("resource/", "2", []string{"resource"}, RestoreAssets)
}

func main() {
	router.Start()
}

// 释放资源文件
func unzipAssets(path, ver string, dirs []string, call func(s1, s2 string) error) {
	if _, err := os.Stat(path); err == nil {
		if _, err := os.Stat(path + ver + ".ver"); os.IsNotExist(err) {
			gocd.Log.Info("[" + ver + "]: Delete Old Assets.")
			os.RemoveAll(path)
		} else {
			gocd.Log.Info("[" + ver + "]: Assets File Exists.")
			return
		}
	}
	gocd.Log.Info("[" + ver + "]: Unpkg Assets.")
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

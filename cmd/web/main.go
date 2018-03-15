/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package main

import (
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
}

func main() {
	router.Start()
}

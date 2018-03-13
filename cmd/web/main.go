/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package main

import (
	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
	"git.cm/naiba/gocd"
	"git.cm/naiba/gocd/web"
)

func init() {
	// init sentry dsn
	var webDSN = "xxxxxx"
	hook, err := logrus_sentry.NewSentryHook(webDSN, []logrus.Level{
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
	web.Start()
}

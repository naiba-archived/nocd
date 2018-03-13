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
	var webDSN = "https://b2be1d09de6a4765aa1bf2f02c58d156:f1598a174c9441648e09b7d88e29d7a6@sentry.io/301979"
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

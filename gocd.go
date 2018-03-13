/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gocd

import (
	log "github.com/sirupsen/logrus"
)

var Log *log.Logger

func init() {
	if Log == nil {
		Log = log.New()
		Log.SetLevel(log.DebugLevel)
	}
}

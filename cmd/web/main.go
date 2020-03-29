/*
 * Copyright (c) 2017 - 2020, 奶爸<1@5.nu>
 * All rights reserved.
 */

package main

import (
	"github.com/naiba/nocd"
	"github.com/naiba/nocd/router"
)

func init() {
	// initial global settings
	nocd.InitSysConfig("conf/app.ini")
}

func main() {
	router.Start()
}

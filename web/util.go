/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package web

import (
	"github.com/gin-gonic/gin"
	"git.cm/naiba/gocd"
)

func setCookie(c *gin.Context, key string, val string) {
	c.SetCookie(key, val, 60*60*24*365*1.5, "/", gocd.Conf.Section("gocd").Key("domain").String(), false, false)
}

func commonData(c *gin.Context, data gin.H) gin.H {
	isLogin := c.MustGet(CtxIsLogin).(bool)
	data["isLogin"] = isLogin
	if isLogin {
		data["user"] = c.MustGet(CtxUser)
	}
	return data
}

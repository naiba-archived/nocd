/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"git.cm/naiba/gocd/router/admin"
	"git.cm/naiba/gocd/utils/mgin"
	"github.com/gin-gonic/gin"
)

func serveAdmin(r *gin.Engine) {
	ra := r.Group("/admin")
	ra.Use(mgin.FilterMiddleware(mgin.FilterOption{Admin: true}))
	{
		ra.GET("/", admin.Index)
		ra.GET("/user/", admin.User(userService))
		ra.GET("/user/:id/:col/:act", admin.UserToggle(userService))
	}
}

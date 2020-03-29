/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/naiba/nocd/router/admin"
	"github.com/naiba/nocd/utils/mgin"
)

func serveAdmin(r *gin.Engine) {
	ra := r.Group("/admin")
	ra.Use(mgin.FilterMiddleware(mgin.FilterOption{Admin: true}))
	{
		ra.GET("/", admin.Index)
		ra.GET("/user/", admin.User(userService))
		ra.GET("/running/", admin.Running(pipelogService))
		ra.GET("/user/:id/:col/:act", admin.UserToggle(userService))
	}
}

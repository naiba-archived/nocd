/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import "github.com/gin-gonic/gin"

func serveSttings(r *gin.Engine) {
	settings := r.Group("/settings")
	settings.Use(filterMiddleware(filterOption{User: true}))
	{
		settings.GET("/", func(c *gin.Context) {
			c.HTML(200, "settings/index", commonData(c, c.GetBool(CtxIsLogin), gin.H{}))
		})
	}
}

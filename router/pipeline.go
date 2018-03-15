/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import "github.com/gin-gonic/gin"

func servePipeline(r *gin.Engine) {
	pipeline := r.Group("/pipeline")
	pipeline.Use(filterMiddleware(filterOption{User: true}))
	{
		pipeline.GET("/", func(c *gin.Context) {
			c.HTML(200, "pipeline/index", commonData(c, c.GetBool(CtxIsLogin), gin.H{}))
		})
	}
}

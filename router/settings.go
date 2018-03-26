/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"git.cm/naiba/gocd/utils/mgin"
	"github.com/gin-gonic/gin"
	"net/http"
)

func serveSettings(r *gin.Engine) {
	settings := r.Group("/settings")
	settings.Use(mgin.FilterMiddleware(mgin.FilterOption{User: true}))
	{
		settings.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "settings/index", mgin.CommonData(c, c.GetBool(mgin.CtxIsLogin), gin.H{}))
		})
	}
}

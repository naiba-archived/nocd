/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package web

import (
	"github.com/gin-gonic/gin"
)

func Start() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Static("/static", "data/static")
	router.LoadHTMLGlob("data/template/**/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "page/index", nil)
	})

	router.Run(":8000")
}

/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package admin

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/naiba/nocd/utils/mgin"
)

//Index 管理面板首页
func Index(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	c.HTML(http.StatusOK, "admin/index", mgin.CommonData(c, false, gin.H{
		"memory": m,
	}))
}

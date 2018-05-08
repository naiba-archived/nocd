/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/naiba/nocd"
	"github.com/naiba/nocd/utils/mgin"
	"net/http"
	"strconv"
)

//Running 管理部署中的任务
func Running(ps nocd.PipeLogService) func(c *gin.Context) {
	return func(c *gin.Context) {
		page := c.Query("page")
		var pageInt int64
		pageInt, _ = strconv.ParseInt(page, 10, 64)
		if pageInt < 0 {
			c.String(http.StatusForbidden, "GG")
			return
		}
		if pageInt == 0 {
			pageInt = 1
		}

		logs, num := ps.Logs(nocd.PipeLogStatusRunning, pageInt-1, 20)
		for i, l := range logs {
			ps.Pipeline(&l)
			logs[i] = l
		}

		c.HTML(http.StatusOK, "admin/running", mgin.CommonData(c, false, gin.H{
			"logs":        logs,
			"allPage":     num,
			"currentPage": pageInt,
		}))
	}
}

/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package admin

import (
	"git.cm/naiba/gocd/utils/mgin"
	"github.com/gin-gonic/gin"
	"net/http"
)

//Index 管理面板首页
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/index", mgin.CommonData(c, false, gin.H{}))
}

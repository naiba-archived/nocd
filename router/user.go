/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naiba/nocd"
	"github.com/naiba/nocd/utils/mgin"
)

func serveUser(r *gin.Engine) {
	server := r.Group("/user")
	server.Use(mgin.FilterMiddleware(mgin.FilterOption{User: true}))
	{
		server.Any("/transfer", accountTransfer)
	}
}

type accountTransferForm struct {
	Name string `form:"name" binding:"required"`
}

func accountTransfer(c *gin.Context) {
	var req accountTransferForm
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusForbidden, "数据不规范，请检查后重新填写："+err.Error())
		return
	}
	origin := c.MustGet(mgin.CtxUser).(*nocd.User)
	if origin.GName == req.Name {
		c.String(http.StatusForbidden, "目标账户不可与当前账户相同")
		return
	}
	dist, err := userService.UserByGName(req.Name)
	if err != nil {
		c.String(http.StatusForbidden, "目标用户未找到，是否从未登录过平台？")
		return
	}
	if err = userService.Transfer(origin, dist); err != nil {
		c.String(http.StatusForbidden, "数据库错误："+err.Error())
		return
	}
}

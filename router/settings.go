/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package router

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naiba/nocd"
	"github.com/naiba/nocd/utils/mgin"
)

func serveSettings(r *gin.Engine) {
	settings := r.Group("/settings")
	settings.Use(mgin.FilterMiddleware(mgin.FilterOption{User: true}))
	{
		settings.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "settings/index", mgin.CommonData(c, true, gin.H{}))
		})
		settings.POST("/", saveSetting)
	}
}

func saveSetting(c *gin.Context) {
	var uf nocd.User
	if err := c.Bind(&uf); err != nil {
		c.String(http.StatusForbidden, "输入不符合规范："+err.Error())
		return
	}
	var kv map[string]string
	err := json.Unmarshal([]byte(uf.RequestBody), &kv)
	if err != nil {
		c.String(http.StatusForbidden, "输入不符合规范：Body解析错误"+err.Error())
		return
	}
	if uf.RequestType < nocd.RequestTypeJSON || uf.RequestType > nocd.RequestTypeForm || uf.RequestMethod < nocd.RequestMethodGet || uf.RequestMethod > nocd.RequestMethodPost {
		c.String(http.StatusForbidden, "输入不符合规范：类型不存在")
		return
	}

	u := c.MustGet(mgin.CtxUser).(*nocd.User)
	u.WebhookURL = uf.WebhookURL
	u.RequestMethod = uf.RequestMethod
	u.RequestBody = uf.RequestBody
	u.RequestType = uf.RequestType
	u.VerifySSL = uf.VerifySSL
	u.PushSuccess = uf.PushSuccess
	if err := userService.Update(u); err != nil {
		nocd.Logger().Errorln(err)
		c.String(http.StatusInternalServerError, "数据库错误")
		return
	}
}

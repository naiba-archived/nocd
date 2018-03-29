/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"git.cm/naiba/gocd"
	"git.cm/naiba/gocd/utils/ftqq"
	"git.cm/naiba/gocd/utils/mgin"
	"github.com/gin-gonic/gin"
	"net/http"
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
	var uf gocd.User
	if err := c.Bind(&uf); err != nil {
		c.String(http.StatusForbidden, "输入不符合规范："+err.Error())
		return
	}
	u := c.MustGet(mgin.CtxUser).(*gocd.User)
	if u.Sckey != uf.Sckey {
		resp := ftqq.SendMessage(uf.Sckey, "[GoCD - "+gocd.Conf.Section("gocd").Key("domain").String()+"]", "Server酱推送绑定成功。")
		if resp.Errno != 0 {
			gocd.Logger().Errorln(resp.Error)
			c.String(http.StatusForbidden, "SCKEY验证失败："+resp.Errmsg)
			return
		}
	}
	u.Sckey = uf.Sckey
	u.PushSuccess = uf.PushSuccess
	if err := userService.Update(u); err != nil {
		gocd.Logger().Errorln(err)
		c.String(http.StatusInternalServerError, "数据库错误")
		return
	}
}

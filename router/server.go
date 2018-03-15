/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"github.com/gin-gonic/gin"
	"git.cm/naiba/gocd"
	"git.cm/naiba/gocd/ssh"
)

func ServeServer(r *gin.Engine) {
	server := r.Group("/server")
	server.Use(filterMiddleware(filterOption{User: true}))
	{
		server.GET("/", func(c *gin.Context) {
			c.HTML(200, "server/index", commonData(c, c.GetBool(CtxIsLogin), gin.H{
				"servers": serverService.GetServersByUser(c.MustGet(CtxUser).(*gocd.User)),
			}))
		})
		server.POST("/", addServer)
	}
}

func addServer(c *gin.Context) {
	var s gocd.Server
	if err := c.Bind(&s); err != nil {
		c.String(400, "数据不规范，请检查后重新填写"+err.Error())
		return
	}
	user := c.MustGet(CtxUser).(*gocd.User)
	if err := ssh.CheckLogin(s.Address, s.Port, user.PrivateKey, s.Login); err != nil {
		c.String(400, err.Error())
		return
	}
	s.UserID = user.ID
	if err := serverService.CreateServer(&s); err == nil {
		c.String(200, "")
	} else {
		gocd.Log.Error(err)
		c.String(500, "数据库错误")
	}
}

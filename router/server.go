/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/naiba/nocd"
	"github.com/naiba/nocd/utils/mgin"
	"github.com/naiba/nocd/utils/ssh"
	"net/http"
)

func serveServer(r *gin.Engine) {
	server := r.Group("/server")
	server.Use(mgin.FilterMiddleware(mgin.FilterOption{User: true}))
	{
		server.Any("/", serverHandler)
	}
}

func serverHandler(c *gin.Context) {
	method := c.Request.Method
	if method == http.MethodGet {
		c.HTML(http.StatusOK, "server/index", mgin.CommonData(c, c.GetBool(mgin.CtxIsLogin), gin.H{
			"servers": serverService.GetServersByUser(c.MustGet(mgin.CtxUser).(*nocd.User)),
		}))
	} else {
		var s nocd.Server
		user := c.MustGet(mgin.CtxUser).(*nocd.User)
		if err := c.Bind(&s); err != nil {
			c.String(http.StatusForbidden, "数据不规范，请检查后重新填写"+err.Error())
			return
		}
		if method == http.MethodPost {
			if err := ssh.CheckLogin(s.Address, s.Port, user.PrivateKey, s.Login); err != nil {
				c.String(http.StatusForbidden, err.Error())
				return
			}
			s.UserID = user.ID
			if err := serverService.CreateServer(&s); err != nil {
				nocd.Logger().Errorln(err)
				c.String(http.StatusInternalServerError, "数据库错误")
			}
		} else {
			if s.ID < 1 {
				c.String(http.StatusForbidden, "ID错误")
				return
			}
			// 用户鉴权
			server, err := serverService.GetServersByUserAndSid(user, s.ID)
			if err != nil {
				c.String(http.StatusForbidden, "ID错误")
				return
			}
			if method == http.MethodPatch {
				if err := ssh.CheckLogin(s.Address, s.Port, user.PrivateKey, s.Login); err != nil {
					c.String(http.StatusForbidden, err.Error())
					return
				}
				// 更新数据
				server.Name = s.Name
				server.Address = s.Address
				server.Login = s.Login
				server.Port = s.Port
				if err := serverService.UpdateServer(&server); err != nil {
					nocd.Logger().Errorln(err)
					c.String(http.StatusInternalServerError, "数据库错误")
				}
			} else if method == http.MethodDelete {
				if serverService.DeleteServer(s.ID) != nil {
					nocd.Logger().Errorln(err)
					c.String(http.StatusInternalServerError, "数据库错误")
					return
				}
			} else {
				c.String(http.StatusForbidden, "非法请求")
			}
		}
	}
}

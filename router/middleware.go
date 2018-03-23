/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"github.com/gin-gonic/gin"
)

//CtxIsLogin 用户是否登录
const CtxIsLogin = "login"

//CtxUser 用户Key
const CtxUser = "user"

func authMiddleware(c *gin.Context) {
	uid, err := c.Cookie("uid")
	login := false
	if err == nil {
		token, err := c.Cookie("token")
		if err == nil && len(uid) > 0 && len(token) > 0 {
			u, err := userService.VerifyUser(uid, token)
			login = err == nil
			if login {
				c.Set(CtxUser, u)
			}
		}
	}
	c.Set(CtxIsLogin, login)
	c.Header("X-Owner", "Naiba(1@5.nu)")
}

type filterOption struct {
	User  bool
	Guest bool
}

func filterMiddleware(o filterOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		if o.Guest && c.MustGet(CtxIsLogin).(bool) {
			jsAlertAndRedirect("限制已登录用户访问", "/", c)
		}
		if o.User && !c.MustGet(CtxIsLogin).(bool) {
			jsAlertAndRedirect("需要登录", "/", c)
		}
	}
}

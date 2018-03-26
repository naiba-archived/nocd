/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package mgin

import (
	"git.cm/naiba/gocd"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"git.cm/naiba/com"
)

//CtxIsLogin 用户是否登录
const CtxIsLogin = "login"

//CtxUser 用户Key
const CtxUser = "user"

//AuthMiddleware 身份验证中间件
func AuthMiddleware(userService gocd.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := c.Cookie("uid")
		login := false
		if err == nil {
			token, err := c.Cookie("token")
			if err == nil && len(uid) > 0 && len(token) > 0 {
				u, err := userService.Verify(uid, token)
				login = err == nil
				if login {
					// 账户已被锁定
					if u.IsBlocked {
						u.Token = com.MD5("blocked" + time.Now().String())
						userService.Update(u)
						c.String(http.StatusForbidden, "您的账户已被锁定")
						c.Abort()
						return
					}
					c.Set(CtxUser, u)
				}
			}
		}
		c.Set(CtxIsLogin, login)
		c.Header("X-Owner", "Naiba(1@5.nu)")
	}
}

//FilterOption 权限控制设置
type FilterOption struct {
	User  bool
	Guest bool
	Admin bool
}

//FilterMiddleware 权限控制
func FilterMiddleware(o FilterOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		if o.Guest && c.MustGet(CtxIsLogin).(bool) {
			AlertAndRedirect("限制已登录用户访问", "/", c)
		}
		if (o.User || o.Admin) && !c.MustGet(CtxIsLogin).(bool) {
			AlertAndRedirect("需要登录", "/", c)
		}
		if o.Admin && !c.MustGet(CtxUser).(*gocd.User).IsAdmin {
			AlertAndRedirect("需要管理员权限", "/", c)
		}
	}
}

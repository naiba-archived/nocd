/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package admin

import (
	"git.cm/naiba/gocd"
	"git.cm/naiba/gocd/utils/mgin"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//User 用户管理
func User(us gocd.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin/user", mgin.CommonData(c, false, gin.H{
			"users": us.Users(),
		}))
	}
}

//UserToggle 用户状态管理
func UserToggle(us gocd.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		col := c.Param("col")
		act := c.Param("act")
		uid, err := strconv.ParseInt(id, 10, 64)
		if err != nil || uid < 1 {
			c.String(http.StatusForbidden, "ID有误，请重试")
			return
		}
		u, err := us.UserByGID(uid)
		if err != nil {
			c.String(http.StatusInternalServerError, "获取用户错误:"+err.Error())
			return
		}
		switch col {
		case "admin":
			u.IsAdmin = act == "on"
			err = us.Update(u)
			break
		case "block":
			u.IsBlocked = act == "on"
			err = us.Update(u)
			break
		}
		if err != nil {
			c.String(http.StatusInternalServerError, "更新用户错误:"+err.Error())
		}
	}
}

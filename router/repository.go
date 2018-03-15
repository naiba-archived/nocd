/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"github.com/gin-gonic/gin"
	"git.cm/naiba/gocd"
	"git.cm/naiba/com"
	"fmt"
	"time"
)

func serveRepository(r *gin.Engine) {
	repo := r.Group("/repository")
	repo.Use(filterMiddleware(filterOption{User: true}))
	{
		repo.GET("/", func(c *gin.Context) {
			user := c.MustGet(CtxUser).(*gocd.User)
			c.HTML(200, "repository/index", commonData(c, c.GetBool(CtxIsLogin), gin.H{
				"repos": repoService.GetRepoByUser(user),
			}))
		})
		repo.POST("/", addRepo)
	}
}

func addRepo(c *gin.Context) {
	var repo gocd.Repository
	if err := c.Bind(&repo); err != nil {
		c.String(400, "数据不规范，请检查后重新填写"+err.Error())
		return
	}
	user := c.MustGet(CtxUser).(*gocd.User)
	repo.UserID = user.ID
	repo.Secret = com.MD5(fmt.Sprintf("%d%s%s%d", user.ID, repo.Name, user.GLogin, time.Now().UnixNano()))
	if err := repoService.CreateRepo(&repo); err != nil {
		gocd.Log.Error(err)
		c.String(500, "数据库错误")
	} else {
		c.String(200, "")
	}
}

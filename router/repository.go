/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/naiba/com"
	"github.com/naiba/nocd"
	"github.com/naiba/nocd/utils/mgin"
	"net/http"
	"time"
)

func serveRepository(r *gin.Engine) {
	repo := r.Group("/repository")
	repo.Use(mgin.FilterMiddleware(mgin.FilterOption{User: true}))
	{
		repo.Any("/", repoHandler)
	}
}

func repoHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		user := c.MustGet(mgin.CtxUser).(*nocd.User)
		c.HTML(http.StatusOK, "repository/index", mgin.CommonData(c, c.GetBool(mgin.CtxIsLogin), gin.H{
			"repos":     repoService.GetRepoByUser(user),
			"platforms": nocd.RepoPlatforms,
			"events":    nocd.RepoEvents,
			"servers":   serverService.GetServersByUser(user),
		}))
	} else {
		// 通用数据校验
		var repo nocd.Repository
		if err := c.Bind(&repo); err != nil {
			c.String(http.StatusForbidden, "数据不规范，请检查后重新填写"+err.Error())
			return
		}
		user := c.MustGet(mgin.CtxUser).(*nocd.User)
		if c.Request.Method == http.MethodPost {
			// 添加
			repo.UserID = user.ID
			repo.Secret = com.MD5(fmt.Sprintf("%d%s%s%d", user.ID, repo.Name, user.GLogin, time.Now().UnixNano()))
			if repoService.Create(&repo) != nil {
				c.String(http.StatusInternalServerError, "数据库错误")
			}
		} else {
			// 对 repo 的操作权限
			mRepo, err := repoService.GetRepoByUserAndID(user, repo.ID)
			if err != nil {
				c.String(http.StatusForbidden, "不具备操作权限")
				return
			}
			if c.Request.Method == http.MethodPatch {
				mRepo.Name = repo.Name
				mRepo.Platform = repo.Platform
				if repoService.Update(&mRepo) != nil {
					c.String(http.StatusInternalServerError, "数据库错误")
				}
			} else if c.Request.Method == http.MethodDelete {
				if repoService.Delete(mRepo.ID) != nil {
					c.String(http.StatusInternalServerError, "数据库错误")
				}
			}
		}
	}
}

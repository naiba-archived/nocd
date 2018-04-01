/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"encoding/json"
	"git.cm/naiba/gocd"
	"git.cm/naiba/gocd/utils/mgin"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func servePipeline(r *gin.Engine) {
	pipeline := r.Group("/pipeline")
	pipeline.Use(mgin.FilterMiddleware(mgin.FilterOption{User: true}))
	{
		pipeline.Any("/", pipelineX)
	}
	pipelog := r.Group("/pipelog")
	pipelog.Use(mgin.FilterMiddleware(mgin.FilterOption{User: true}))
	{
		pipelog.GET("/", pipeLogs)
		pipelog.GET("/:id", viewLog)
	}
}

func pipeLogs(c *gin.Context) {
	user := c.MustGet(mgin.CtxUser).(*gocd.User)

	page := c.Query("page")
	var pageInt int64
	pageInt, _ = strconv.ParseInt(page, 10, 64)
	if pageInt < 0 {
		c.String(http.StatusForbidden, "GG")
		return
	}
	if pageInt == 0 {
		pageInt = 1
	}

	logs, num := pipelogService.UserLogs(user.ID, pageInt-1, 20)
	for i, l := range logs {
		pipelogService.Pipeline(&l)
		logs[i] = l
	}

	c.HTML(http.StatusOK, "pipelog/index", mgin.CommonData(c, false, gin.H{
		"logs":        logs,
		"allPage":     num,
		"currentPage": pageInt,
	}))
}

func viewLog(c *gin.Context) {
	lid := c.Param("id")
	user := c.MustGet(mgin.CtxUser).(*gocd.User)
	u64lid, err := strconv.ParseUint(lid, 10, 64)
	if err != nil {
		c.String(http.StatusInternalServerError, "非法ID")
		return
	}
	log, err := pipelogService.GetByUID(user.ID, uint(u64lid))
	if err != nil {
		c.String(http.StatusForbidden, "您无权查看此Log")
		return
	}
	c.HTML(http.StatusOK, "pipelog/log", mgin.CommonData(c, false, gin.H{
		"log": log,
	}))
}

func pipelineX(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "pipeline/index", mgin.CommonData(c, true, gin.H{}))
	} else {
		// 通用数据校验
		var pl gocd.Pipeline
		if err := c.Bind(&pl); err != nil {
			c.String(http.StatusForbidden, "填写数据不规范，请重新输入。"+err.Error())
			return
		}
		tmp, err := json.Marshal(pl.EventsSlice)
		if err != nil {
			gocd.Logger().Errorln(err)
			c.String(http.StatusInternalServerError, "序列化失败，请重试。"+err.Error())
			return
		}
		pl.Events = string(tmp)
		user := c.MustGet(mgin.CtxUser).(*gocd.User)
		repo, err := repoService.GetRepoByUserAndID(user, pl.RepositoryID)
		if err != nil {
			c.String(http.StatusForbidden, "这个项目不属于您，您无权操作。")
			return
		}
		if !validEvents(pl.EventsSlice, repo.Platform) {
			c.String(http.StatusForbidden, "非法的监控事件。")
			return
		}
		// 校验对于 Server 的操作权限
		_, err = serverService.GetServersByUserAndSid(user, pl.ServerID)
		if err != nil {
			gocd.Logger().Debug(err)
			c.String(http.StatusForbidden, "这个服务器不属于您，您无权操作。")
			return
		}
		if c.Request.Method == http.MethodPost {
			pl.UserID = user.ID
			if err = pipelineService.Create(&pl); err != nil {
				gocd.Logger().Errorln(err)
				c.String(http.StatusInternalServerError, "数据库错误。")
			}
		} else {
			// 校验对于 Pipeline 的操作权限
			pip, err := pipelineService.UserPipeline(user.ID, pl.ID)
			if err != nil {
				c.String(http.StatusForbidden, "您无权操作此 Pipeline")
				return
			}
			if c.Request.Method == http.MethodPatch {
				pip.Name = pl.Name
				pip.Events = pl.Events
				pip.Shell = pl.Shell
				pip.ServerID = pl.ServerID
				pip.Branch = pl.Branch
				if pipelineService.Update(&pip) != nil {
					c.String(http.StatusInternalServerError, "数据库错误。")
				}
			} else if c.Request.Method == http.MethodDelete {
				if pipelineService.Delete(pip.ID) != nil {
					c.String(http.StatusInternalServerError, "数据库错误。")
				}
			} else {
				c.String(http.StatusForbidden, "非法访问")
			}
		}

	}
}

func validEvents(events []string, platform int) bool {
	for _, event := range events {
		if _, has := gocd.RepoEvents[platform][event]; !has {
			return false
		}
	}
	return true
}

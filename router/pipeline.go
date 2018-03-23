/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"git.cm/naiba/gocd"
	"encoding/json"
	"strconv"
)

func servePipeline(r *gin.Engine) {
	pipeline := r.Group("/pipeline")
	pipeline.Use(filterMiddleware(filterOption{User: true}))
	{
		pipeline.Any("/", pipelineX)
	}
	pipelog := r.Group("/pipelog")
	pipelog.Use(filterMiddleware(filterOption{User: true}))
	{
		pipelog.GET("/", pipeLogs)
		pipelog.GET("/:id", viewLog)
	}
}

func pipeLogs(c *gin.Context) {
	user := c.MustGet(CtxUser).(*gocd.User)
	logs := pipelogService.UserLogs(user.ID)
	for i, l := range logs {
		pipelogService.Pipeline(&l)
		logs[i] = l
	}
	c.HTML(http.StatusOK, "pipelog/index", commonData(c, false, gin.H{
		"logs": logs,
	}))
}

func viewLog(c *gin.Context) {
	lid := c.Param("id")
	user := c.MustGet(CtxUser).(*gocd.User)
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
	c.HTML(http.StatusOK, "pipelog/log", commonData(c, false, gin.H{
		"log": log,
	}))
}

func pipelineX(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "pipeline/index", commonData(c, c.GetBool(CtxIsLogin), gin.H{
		}))
	} else {
		// 通用数据校验
		var pl gocd.Pipeline
		if err := c.Bind(&pl); err != nil {
			c.String(http.StatusForbidden, "填写数据不规范，请重新输入。"+err.Error())
			return
		}
		tmp, err := json.Marshal(pl.EventsSlice)
		if err != nil {
			gocd.Log.Error(err)
			c.String(http.StatusInternalServerError, "序列化失败，请重试。"+err.Error())
			return
		}
		pl.Events = string(tmp)
		user := c.MustGet(CtxUser).(*gocd.User)
		repo, err := repoService.GetRepoByUserAndID(user, pl.RepositoryID)
		if err != nil {
			gocd.Log.Debug(err)
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
			gocd.Log.Debug(err)
			c.String(http.StatusForbidden, "这个服务器不属于您，您无权操作。")
			return
		}
		if c.Request.Method == http.MethodPost {
			pl.UserID = user.ID
			gocd.Log.Error(pl.RepositoryID, pl.Repository)
			if err = pipelineService.Create(&pl); err != nil {
				gocd.Log.Error(err)
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
				if pipelineService.Delete(pl.ID) != nil {
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

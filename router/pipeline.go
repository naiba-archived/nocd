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
)

func servePipeline(r *gin.Engine) {
	pipeline := r.Group("/pipeline")
	pipeline.Use(filterMiddleware(filterOption{User: true}))
	{
		pipeline.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "pipeline/index", commonData(c, c.GetBool(CtxIsLogin), gin.H{
			}))
		})
		pipeline.POST("/", addPipeline)
	}
	pipelog := r.Group("/pipelog")
	pipelog.Use(filterMiddleware(filterOption{User: true}))
	{
		pipelog.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "pipelog/index", commonData(c, c.GetBool(CtxIsLogin), gin.H{
			}))
		})
	}
}

func addPipeline(c *gin.Context) {
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
	_, err = serverService.GetServersByUserAndSid(user, pl.ServerID)
	if err != nil {
		gocd.Log.Debug(err)
		c.String(http.StatusForbidden, "这个服务器不属于您，您无权操作。")
		return
	}
	pl.UserID = user.ID
	gocd.Log.Error(pl.RepositoryID, pl.Repository)
	if err = pipelineService.CreatePipeline(&pl); err != nil {
		gocd.Log.Error(err)
		c.String(http.StatusInternalServerError, "数据库错误。")
		return
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

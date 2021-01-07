/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/naiba/nocd"
	"github.com/naiba/nocd/utils/mgin"
)

func servePipeline(r *gin.Engine) {
	pipeline := r.Group("/pipeline")
	pipeline.Use(mgin.FilterMiddleware(mgin.FilterOption{User: true}))
	{
		pipeline.Any("/", pipelineHandler)
		pipeline.Any("/webhook", pipelineWebhook)
	}
	pipelog := r.Group("/pipelog")
	pipelog.Use(mgin.FilterMiddleware(mgin.FilterOption{User: true}))
	{
		pipelog.GET("/", pipeLogs)
		pipelog.GET("/:id", viewLog)
	}
}

func pipelineWebhook(c *gin.Context) {
	user := c.MustGet(mgin.CtxUser).(*nocd.User)
	var wb nocd.Webhook
	if err := c.ShouldBind(&wb); err != nil {
		c.String(http.StatusForbidden, "数据不规范，请检查后重新填写"+err.Error())
		return
	}
	err := wb.Send("测试消息", &nocd.Pipeline{}, &nocd.PipeLog{})
	if err != nil {
		c.String(http.StatusForbidden, "测试消息发送失败"+err.Error())
		return
	}
	if wb.RequestType < nocd.WebhookRequestTypeJSON || wb.RequestType > nocd.WebhookRequestTypeForm || wb.RequestMethod < nocd.WebhookRequestMethodGET || wb.RequestMethod > nocd.WebhookRequestMethodPOST {
		c.String(http.StatusForbidden, "输入不符合规范：类型不存在")
		return
	}
	// 鉴权
	_, err = pipelineService.UserPipeline(user.ID, wb.PipelineID)
	if err != nil {
		c.String(http.StatusForbidden, "Pipeline 不存在："+err.Error())
		return
	}

	var webhook nocd.Webhook
	if c.Request.Method == http.MethodPost {
		webhook.PipelineID = wb.PipelineID
		webhook.URL = wb.URL
		webhook.RequestMethod = wb.RequestMethod
		webhook.RequestType = wb.RequestType
		webhook.RequestBody = wb.RequestBody
		webhook.VerifySSL = wb.VerifySSL
		webhook.PushSuccess = wb.PushSuccess
		webhook.Enable = wb.Enable
		err = webhookService.Create(&webhook)
	}

	if c.Request.Method == http.MethodDelete || c.Request.Method == http.MethodPatch {
		err = db.First(&webhook, "id = ? AND pipeline_id = ?", wb.ID, wb.PipelineID).Error
		if err == nil {
			if c.Request.Method == http.MethodPatch {
				webhook.URL = wb.URL
				webhook.RequestMethod = wb.RequestMethod
				webhook.RequestType = wb.RequestType
				webhook.RequestBody = wb.RequestBody
				webhook.VerifySSL = wb.VerifySSL
				webhook.PushSuccess = wb.PushSuccess
				webhook.Enable = wb.Enable
				err = db.Save(&webhook).Error
			} else {
				err = db.Delete(&webhook).Error
			}
		}
	}

	if err != nil {
		c.String(http.StatusForbidden, "数据库错误："+err.Error())
		return
	}
}

func pipeLogs(c *gin.Context) {
	user := c.MustGet(mgin.CtxUser).(*nocd.User)

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

	var repoIDs []uint
	for i, l := range logs {
		pipelogService.Pipeline(&l)
		logs[i] = l
		repoIDs = append(repoIDs, l.Pipeline.RepositoryID)
	}

	var repos []nocd.Repository
	db.Find(&repos, "id in (?)", repoIDs)
	repoIndex := make(map[uint]nocd.Repository)
	for i := 0; i < len(repos); i++ {
		repoIndex[repos[i].ID] = repos[i]
	}

	for i := 0; i < len(logs); i++ {
		logs[i].Pipeline.Repository = repoIndex[logs[i].Pipeline.RepositoryID]
	}

	c.HTML(http.StatusOK, "pipelog/index", mgin.CommonData(c, false, gin.H{
		"logs":        logs,
		"allPage":     num,
		"currentPage": pageInt,
	}))
}

func viewLog(c *gin.Context) {
	lid := c.Param("id")
	user := c.MustGet(mgin.CtxUser).(*nocd.User)
	u64lid, err := strconv.ParseUint(lid, 10, 64)
	if err != nil {
		c.String(http.StatusInternalServerError, "非法ID")
		return
	}
	var log nocd.PipeLog
	if !user.IsAdmin {
		log, err = pipelogService.GetByUID(user.ID, uint(u64lid))
		if err != nil {
			c.String(http.StatusForbidden, "您无权查看此Log")
			return
		}
	} else {
		log, err = pipelogService.GetByID(uint(u64lid))
		if err != nil {
			c.String(http.StatusForbidden, "Log不存在")
			return
		}
	}
	isAjax := c.Query("ajax") != ""
	lineNumberStr := c.Query("line")
	actStr := c.Query("act")
	lineNumber, _ := strconv.Atoi(lineNumberStr)
	if isAjax {
		nocd.RunningLogsLock.RLock()
		defer nocd.RunningLogsLock.RUnlock()
		run, has := nocd.RunningLogs[log.ID]
		switch actStr {
		case "view":
			if has {
				if lineNumber == 0 {
					c.JSON(http.StatusOK, map[string]string{
						"end":  "false",
						"log":  strings.Join(run.RunningLog, "\n"),
						"line": strconv.Itoa(len(run.RunningLog)),
					})
				} else {
					if lineNumber > len(run.RunningLog)-1 {
						c.JSON(http.StatusOK, map[string]string{
							"end":  "false",
							"log":  "",
							"line": lineNumberStr,
						})
					} else {
						c.JSON(http.StatusOK, map[string]string{
							"end":  "false",
							"log":  strings.Join(run.RunningLog[lineNumber:], "\n"),
							"line": strconv.Itoa(len(run.RunningLog)),
						})
					}
				}
			} else {
				c.JSON(http.StatusOK, map[string]string{
					"end":  "true",
					"log":  "00:00:00#部署已结束。\n",
					"line": "0",
				})
			}
			break
		case "stop":
			if has {
				run.Log.Status = nocd.PipeLogStatusHumanStopped
				run.Finish <- true
			} else {
				if log.Status != nocd.PipeLogStatusRunning {
					c.String(http.StatusOK, "部署已经停止，无需再次停止。")
				} else {
					log.Status = nocd.PipeLogStatusHumanStopped
					log.StoppedAt = time.Now()
					err := pipelogService.Update(&log)
					if err != nil {
						nocd.Logger().Error(err)
					}
				}
			}
			c.String(http.StatusOK, "success")
			break
		}
		return
	}
	c.HTML(http.StatusOK, "pipelog/log", mgin.CommonData(c, false, gin.H{
		"log":       log,
		"fromAdmin": c.Query("admin") == "true",
	}))
}

func pipelineHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "pipeline/index", mgin.CommonData(c, true, gin.H{}))
	} else {
		// 通用数据校验
		var pl nocd.Pipeline
		if err := c.Bind(&pl); err != nil {
			c.String(http.StatusForbidden, "填写数据不规范，请重新输入。"+err.Error())
			return
		}
		tmp, err := json.Marshal(pl.EventsSlice)
		if err != nil {
			nocd.Logger().Errorln(err)
			c.String(http.StatusInternalServerError, "序列化失败，请重试。"+err.Error())
			return
		}
		pl.Events = string(tmp)
		user := c.MustGet(mgin.CtxUser).(*nocd.User)
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
			nocd.Logger().Debug(err)
			c.String(http.StatusForbidden, "这个服务器不属于您，您无权操作。")
			return
		}
		if c.Request.Method == http.MethodPost {
			pl.UserID = user.ID
			if err = pipelineService.Create(&pl); err != nil {
				nocd.Logger().Errorln(err)
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
		if _, has := nocd.RepoEvents[platform][event]; !has {
			return false
		}
	}
	return true
}

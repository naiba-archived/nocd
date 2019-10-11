/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	client "github.com/gogs/go-gogs-client"
	"gopkg.in/go-playground/webhooks.v3"
	"gopkg.in/go-playground/webhooks.v3/bitbucket"
	"gopkg.in/go-playground/webhooks.v3/github"
	"gopkg.in/go-playground/webhooks.v3/gitlab"
	"gopkg.in/go-playground/webhooks.v3/gogs"

	"time"

	"github.com/naiba/nocd"
	"github.com/naiba/nocd/utils/ftqq"
	"github.com/naiba/nocd/utils/ssh"
)

var webHookSQLIndex map[string]string

func init() {
	webHookSQLIndex = make(map[string]string)
	// github
	webHookSQLIndex["github.PushPayload"] = string(github.PushEvent)
	// bitBucket
	webHookSQLIndex["bitbucket.PullRequestMergedPayload"] = string(bitbucket.PullRequestMergedEvent)
	// gitlab
	webHookSQLIndex["gitlab.PushEventPayload"] = string(gitlab.PushEvents)
	// gogs
	webHookSQLIndex["gogs.PushPayload"] = string(gogs.PushEvent)
}

func serveWebHook(r *gin.Engine) {
	hook := r.Group("/webhook")
	{
		hook.POST("/:id", webHook)
	}
}

func webHook(c *gin.Context) {
	rid := c.Param("id")
	id, err := strconv.ParseUint(rid, 10, 64)
	if err != nil {
		c.String(http.StatusForbidden, "ID转换错误："+err.Error())
		return
	}
	repo, err := repoService.GetRepoByID(uint(id))
	if err != nil {
		c.String(http.StatusForbidden, "项目不存在："+err.Error())
		return
	}
	// 设置监听事件
	var hook webhooks.Webhook
	webhooks.DefaultLog = webhooks.NewLogger(nocd.Debug)
	switch repo.Platform {
	case nocd.RepoPlatGitHub:
		gh := github.New(&github.Config{Secret: repo.Secret})
		gh.RegisterEvents(dispatchWebHook(repo.ID), github.PushEvent)
		hook = gh
		break
	case nocd.RepoPlatBitBucket:
		bb := bitbucket.New(&bitbucket.Config{UUID: repo.Secret})
		bb.RegisterEvents(dispatchWebHook(repo.ID), bitbucket.PullRequestMergedEvent)
		hook = bb
		break
	case nocd.RepoPlatGitlab:
		gl := gitlab.New(&gitlab.Config{Secret: repo.Secret})
		gl.RegisterEvents(dispatchWebHook(repo.ID), gitlab.PushEvents)
		hook = gl
		break
	case nocd.RepoPlatGogs:
		gs := gogs.New(&gogs.Config{Secret: repo.Secret})
		gs.RegisterEvents(dispatchWebHook(repo.ID), gogs.PushEvent)
		hook = gs
		break
	default:
		c.String(http.StatusInternalServerError, "服务器错误，不支持的托管平台："+strconv.Itoa(repo.Platform))
		return
	}
	hook.ParsePayload(c.Writer, c.Request)
}

func dispatchWebHook(id uint) webhooks.ProcessPayloadFunc {
	return func(payload interface{}, header webhooks.Header) {
		payloadType := reflect.TypeOf(payload).String()
		p, has := webHookSQLIndex[payloadType]
		if !has {
			return
		}
		who, branch := parsePayloadInfo(payload)
		ps, err := pipelineService.GetPipelinesByRidAndEventAndBranch(id, p, branch)
		if err != nil {
			return
		}
		for _, p := range ps {
			if err := pipelineService.Server(&p); err == nil {
				pipelineService.User(&p)
				go deploy(p, who)
			} else {
				nocd.Logger().Errorln(err)
			}
		}
	}
}

func parsePayloadInfo(payload interface{}) (string, string) {
	var who = "unknown"
	var branch = "unknown"
	switch payload.(type) {
	case github.PushPayload:
		p := payload.(github.PushPayload).Pusher
		who = p.Name + "(" + p.Email + ")"
		branch = payload.(github.PushPayload).Ref
		break

	case bitbucket.PullRequestMergedPayload:
		p := payload.(bitbucket.PullRequestMergedPayload)
		who = p.Actor.Username
		branch = p.PullRequest.Destination.Branch.Name
		break

	case gitlab.PushEventPayload:
		p := payload.(gitlab.PushEventPayload)
		who = p.UserName + "(" + p.UserEmail + ")"
		branch = p.Ref[11:]
		break

	case client.PushPayload:
		p := payload.(client.PushPayload)
		who = p.Pusher.UserName + "(" + p.Pusher.Email + ")"
		branch = p.Ref[11:]
		break
	}
	return who, branch
}

func deploy(pipeline nocd.Pipeline, who string) {
	var deployLog nocd.PipeLog
	deployLog.PipelineID = pipeline.ID
	deployLog.StartedAt = time.Now()
	deployLog.Log = ""
	deployLog.Pusher = who
	deployLog.Status = nocd.PipeLogStatusRunning
	// 更新运行中
	pipelogService.Create(&deployLog)
	// 进行部署
	ssh.Deploy(pipeline, &deployLog)
	// 部署完成
	pipelogService.Update(&deployLog)

	if (deployLog.Status == nocd.PipeLogStatusSuccess && !pipeline.User.PushSuccess) || len(pipeline.User.Sckey) < 1 {
		return
	}

	status := ""
	switch deployLog.Status {
	case nocd.PipeLogStatusSuccess:
		status = "交付成功"
		break
	case nocd.PipeLogStatusErrorShellExec:
		status = "Shell错误"
		break
	case nocd.PipeLogStatusErrorServerConn:
		status = "服务器连接错误"
		break
	case nocd.PipeLogStatusHumanStopped:
		status = "人工停止"
		break
	case nocd.PipeLogStatusErrorTimeout:
		status = "执行超时"
		break
	default:
		status = "未知错误"
	}

	ftqq.SendMessage(pipeline.User.Sckey, "[NoCD]"+pipeline.Name+"-"+status, "# 部署日志\r\n```\r\n"+deployLog.Log+"\r\n```")
}

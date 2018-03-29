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
	client "github.com/gogits/go-gogs-client"
	"github.com/naiba/webhooks"
	"github.com/naiba/webhooks/bitbucket"
	"github.com/naiba/webhooks/github"
	"github.com/naiba/webhooks/gitlab"
	"github.com/naiba/webhooks/gogs"

	"git.cm/naiba/gocd"
	"git.cm/naiba/gocd/utils/ftqq"
	"git.cm/naiba/gocd/utils/ssh"
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
	webhooks.DefaultLog = webhooks.NewLogger(gocd.Debug)
	switch repo.Platform {
	case gocd.RepoPlatGitHub:
		gh := github.New(&github.Config{Secret: repo.Secret})
		gh.RegisterEvents(dispatchWebHook(repo.ID), github.PushEvent)
		hook = gh
		break
	case gocd.RepoPlatBitBucket:
		bb := bitbucket.New(&bitbucket.Config{UUID: repo.Secret})
		bb.RegisterEvents(dispatchWebHook(repo.ID), bitbucket.PullRequestMergedEvent)
		hook = bb
		break
	case gocd.RepoPlatGitlab:
		gl := gitlab.New(&gitlab.Config{Secret: repo.Secret})
		gl.RegisterEvents(dispatchWebHook(repo.ID), gitlab.PushEvents)
		hook = gl
		break
	case gocd.RepoPlatGogs:
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
				gocd.Logger().Errorln(err)
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
		branch = payload.(github.PushPayload).Ref[11:]
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

func deploy(pipeline gocd.Pipeline, who string) {
	deployLog := ssh.Deploy(pipeline, who)
	pipelogService.Create(&deployLog)

	if (deployLog.Status == gocd.PipeLogStatusSuccess && !pipeline.User.PushSuccess) || len(pipeline.User.Sckey) < 1 {
		return
	}

	status := ""
	switch deployLog.Status {
	case gocd.PipeLogStatusSuccess:
		status = "交付成功"
		break
	case gocd.PipeLogStatusErrorShellExec:
		status = "Shell错误"
		break
	case gocd.PipeLogStatusErrorServerConn:
		status = "服务器连接错误"
		break
	default:
		status = "未知错误"
	}

	ftqq.SendMessage(pipeline.User.Sckey, "[GoCD]"+pipeline.Name+"-"+status, "部署日志：\n\n```\n"+deployLog.Log+"\n```")
}

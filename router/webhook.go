/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package router

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	client "github.com/gogs/go-gogs-client"
	"gopkg.in/go-playground/webhooks.v5/bitbucket"
	"gopkg.in/go-playground/webhooks.v5/github"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
	"gopkg.in/go-playground/webhooks.v5/gogs"

	"time"

	"github.com/naiba/nocd"
	"github.com/naiba/nocd/utils/ssh"
)

var webhookSQLIndex map[string]string

func init() {
	webhookSQLIndex = make(map[string]string)
	// github
	webhookSQLIndex["github.PushPayload"] = string(github.PushEvent)
	// bitBucket
	webhookSQLIndex["bitbucket.RepoPushPayload"] = string(bitbucket.RepoPushEvent)
	// gitlab
	webhookSQLIndex["gitlab.PushEventPayload"] = string(gitlab.PushEvents)
	// gogs
	webhookSQLIndex["gogs.PushPayload"] = string(gogs.PushEvent)
}

func serveWebHook(r *gin.Engine) {
	hook := r.Group("/webhook")
	{
		hook.POST("/:id", webhook)
	}
}

func webhook(c *gin.Context) {
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

	switch repo.Platform {
	case nocd.RepoPlatGitHub:
		gh, err := github.New(github.Options.Secret(repo.Secret))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		payload, err := gh.Parse(c.Request, github.PushEvent)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		proccessPayload(repo.ID, payload)

	case nocd.RepoPlatBitBucket:
		bb, err := bitbucket.New(bitbucket.Options.UUID(repo.Secret))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		payload, err := bb.Parse(c.Request, bitbucket.RepoPushEvent)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		proccessPayload(repo.ID, payload)

	case nocd.RepoPlatGitlab:
		gl, err := gitlab.New(gitlab.Options.Secret(repo.Secret))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		payload, err := gl.Parse(c.Request, gitlab.PushEvents)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		proccessPayload(repo.ID, payload)

	case nocd.RepoPlatGogs:
		gs, err := gogs.New(gogs.Options.Secret(repo.Secret))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		payload, err := gs.Parse(c.Request, gogs.PushEvent)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		proccessPayload(repo.ID, payload)

	default:
		c.String(http.StatusInternalServerError, "服务器错误，不支持的托管平台："+strconv.Itoa(repo.Platform))
	}
}

func proccessPayload(repoID uint, payload interface{}) {
	payloadType := reflect.TypeOf(payload).String()
	p, has := webhookSQLIndex[payloadType]
	if !has {
		return
	}
	who, branch := parsePayloadInfo(payload)
	ps, err := pipelineService.GetPipelinesByRidAndEventAndBranch(repoID, p, branch)
	if err != nil {
		return
	}
	for _, p := range ps {
		if err := pipelineService.Server(&p); err == nil {
			pipelineService.User(&p)
			pipelineService.Webhooks(&p)
			go deploy(p, who)
		} else {
			nocd.Logger().Errorln(err)
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
		if strings.HasPrefix(branch, "refs/heads/") {
			branch = branch[11:]
		}
		break

	case bitbucket.PullRequestMergedPayload:
		p := payload.(bitbucket.PullRequestMergedPayload)
		who = p.Actor.DisplayName
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

	status := ""
	switch deployLog.Status {
	case nocd.PipeLogStatusSuccess:
		status = "交付成功"
		break
	case nocd.PipeLogStatusErrorExec:
		status = "执行错误"
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

	nocd.Logger().Infoln("Pipeline 交付完成", deployLog.Pipeline.Name)
	var wg sync.WaitGroup
	for i := 0; i < len(pipeline.Webhook); i++ {
		wg.Add(1)
		go procWebhook(&wg, pipeline.Webhook[i], status, &pipeline, &deployLog)
	}
	wg.Wait()
	nocd.Logger().Infoln("Pipeline 推送完成", deployLog.Pipeline.Name)
}

func procWebhook(wg *sync.WaitGroup, w nocd.Webhook, status string, pipeline *nocd.Pipeline, deployLog *nocd.PipeLog) {
	defer wg.Done()
	var pushSuccess bool

	// 检查 Webhook 状态
	if w.Enable == nil || !*w.Enable {
		return
	}
	if w.PushSuccess != nil && *w.PushSuccess {
		pushSuccess = true
	}
	if !pushSuccess && deployLog.Status == nocd.PipeLogStatusSuccess {
		return
	}

	nocd.Logger().Infoln("Webhook 触发", w.ID, w.URL)
	if err := w.Send(status, pipeline, deployLog); err != nil {
		nocd.Logger().Error(err)
	}
}

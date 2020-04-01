/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package router

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	client "github.com/gogs/go-gogs-client"
	"gopkg.in/go-playground/webhooks.v3"
	"gopkg.in/go-playground/webhooks.v3/bitbucket"
	"gopkg.in/go-playground/webhooks.v3/github"
	"gopkg.in/go-playground/webhooks.v3/gitlab"
	"gopkg.in/go-playground/webhooks.v3/gogs"

	"time"

	"github.com/naiba/nocd"
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
				pipelineService.Webhooks(&p)
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
		if strings.HasPrefix(branch, "refs/heads/") {
			branch = branch[11:]
		}
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
	var verifySSL, pushSuccess bool

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
	if w.VerifySSL != nil && *w.VerifySSL {
		verifySSL = true
	}

	nocd.Logger().Infoln("Webhook 触发", w.ID, w.URL)
	var err error
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: verifySSL},
	}
	client := &http.Client{Transport: transCfg, Timeout: time.Minute * 10}
	var reqURL *url.URL
	reqURL, err = url.Parse(w.URL)
	var data map[string]string
	if err == nil {
		err = json.Unmarshal([]byte(w.RequestBody), &data)
	}
	var resp *http.Response
	if err == nil {
		if w.RequestMethod == nocd.WebhookRequestMethodGET {
			// GET 请求的 Webhook
			for k, v := range data {
				reqURL.Query().Set(k, replaceParamsInString(v, status, pipeline, deployLog))
			}
			resp, err = client.Get(reqURL.String())
		} else {
			// POST 请求的 Webhook
			if w.RequestType == nocd.WebhookRequestTypeForm {
				params := url.Values{}
				for k, v := range data {
					params.Add(k, replaceParamsInString(v, status, pipeline, deployLog))
				}
				resp, err = client.PostForm(reqURL.String(), params)
			} else {
				for k, v := range data {
					data[k] = replaceParamsInString(v, status, pipeline, deployLog)
				}
				var jsonValue []byte
				jsonValue, err = json.Marshal(data)
				if err == nil {
					resp, err = client.Post(reqURL.String(), "application/json", bytes.NewBuffer(jsonValue))
				}
			}
		}
	}
	if err != nil {
		nocd.Logger().Error(err)
		return
	}
	if resp != nil && (resp.StatusCode < 200 || resp.StatusCode > 299) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			nocd.Logger().Error(err)
		}
		nocd.Logger().Error(string(body))
	}
}

func replaceParamsInString(str string, status string, pipeline *nocd.Pipeline, pipelog *nocd.PipeLog) string {
	var dist string
	dist = strings.ReplaceAll(str, "#Pusher#", pipelog.Pusher)
	dist = strings.ReplaceAll(str, "#Log#", pipelog.Log)
	dist = strings.ReplaceAll(str, "#Status#", status)
	dist = strings.ReplaceAll(str, "#PipelineName#", pipeline.Name)
	dist = strings.ReplaceAll(str, "#PipelineID#", fmt.Sprintf("%d", pipeline.ID))
	dist = strings.ReplaceAll(str, "#StartedAt#", pipelog.StartedAt.String())
	dist = strings.ReplaceAll(str, "#StoppedAt#", pipelog.StoppedAt.String())
	return dist
}

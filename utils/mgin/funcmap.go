/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package mgin

import (
	"encoding/json"
	"fmt"
	"html/template"
	"runtime"
	"strings"
	"time"

	"github.com/naiba/nocd"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//FuncMap 自定义模板函数
func FuncMap(pipelineService nocd.PipelineService, pipelogService nocd.PipeLogService,
	webhookService nocd.WebhookService) template.FuncMap {
	return template.FuncMap{
		"unsafe": func(raw string) template.HTML {
			return template.HTML(raw)
		},
		"RepoPipelines": func(rid uint) []nocd.Pipeline {
			return pipelineService.RepoPipelines(&nocd.Repository{ID: rid})
		},
		"UserPipelines": func(uid uint) []nocd.Pipeline {
			var u nocd.User
			u.ID = uid
			return pipelineService.UserPipelines(&u)
		},
		"PipelineWebhooks": func(pid uint) []nocd.Webhook {
			return webhookService.PipelineWebhooks(&nocd.Pipeline{ID: pid})
		},
		"LastServerLog": func(rid uint) nocd.PipeLog {
			return pipelogService.LastServerLog(rid)
		},
		"JSON": func(obj interface{}) string {
			b, _ := json.Marshal(obj)
			return string(b)
		},
		"LastPipelineLog": func(pid uint) nocd.PipeLog {
			return pipelogService.LastPipelineLog(pid)
		},
		"TimeDiff": func(t1, t2 time.Time) string {
			if t2.IsZero() {
				return "Running"
			}
			sec := t2.Sub(t1).Seconds()
			if sec < 60 {
				return fmt.Sprintf(" %.0f s", sec)
			}
			if sec < 60*60 {
				return fmt.Sprintf(" %.0f minute", sec/60)
			}
			if sec < 60*60*24 {
				return fmt.Sprintf(" %.0f H", sec/60/60)
			}
			if sec < 60*60*24*30 {
				return fmt.Sprintf(" %.0f day", sec/60/60/24)
			}
			if sec < 60*60*24*30*12 {
				return fmt.Sprintf(" %.0f month", sec/60/60/24/30)
			}
			return fmt.Sprintf(" %.0f year", sec/60/60/24/30/12)
		},
		"Now": func() time.Time {
			return time.Now().In(nocd.Loc)
		},
		"TimeFormat": func(t time.Time) string {
			return t.In(nocd.Loc).Format("2006-01-02 15:04:05")
		},
		"HasPrefix": strings.HasPrefix,
		"MathSub": func(o, n int64) int64 {
			return o - n
		},
		"MathAdd": func(o, n int64) int64 {
			return o + n
		},
		"NumGoroutine": runtime.NumGoroutine,
		"T": func(localizer *i18n.Localizer, key string, data interface{}) string {
			return localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID:    key,
				TemplateData: data,
			})
		},
	}
}

/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"git.cm/naiba/gocd"
	"time"
	"fmt"
)

func setFuncMap(r *gin.Engine) {
	r.SetFuncMap(template.FuncMap{
		"RepoPipelines": func(rid uint) []gocd.Pipeline {
			return pipelineService.RepoPipelines(&gocd.Repository{ID: rid})
		},
		"LastServerLog": func(rid uint) gocd.PipeLog {
			return pipelogService.LastServerLog(rid)
		},
		"LastPipelineLog": func(pid uint) gocd.PipeLog {
			return pipelogService.LastPipelineLog(pid)
		},
		"TimeDiff": func(t1, t2 time.Time) string {
			sec := t2.Sub(t1).Seconds()
			if sec < 60 {
				return fmt.Sprintf(" %.0f 秒", sec)
			}
			if sec < 60*60 {
				return fmt.Sprintf(" %.0f 分钟", sec/60)
			}
			if sec < 60*60*24 {
				return fmt.Sprintf(" %.0f 分钟", sec/60/60)
			}
			if sec < 60*60*24*30 {
				return fmt.Sprintf(" %.0f 天", sec/60/60/24)
			}
			if sec < 60*60*24*30*12 {
				return fmt.Sprintf(" %.0f 个月", sec/60/60/24/30)
			}
			return fmt.Sprintf(" %.0f 年", sec/60/60/24/30/12)
		},
		"Now": func() time.Time {
			return time.Now()
		},
		"TimeFormat": func(t time.Time) string {
			return t.In(gocd.Loc).Format("2006-01-02 15:04:05")
		},
	})
}

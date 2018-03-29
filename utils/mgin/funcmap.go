/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package mgin

import (
	"fmt"
	"git.cm/naiba/gocd"
	"html/template"
	"strings"
	"time"
	"strconv"
)

//Pagination 分页
type Pagination struct {
	No      int64
	Current bool
	Text    string
}

//FuncMap 自定义模板函数
func FuncMap(pipelineService gocd.PipelineService, pipelogService gocd.PipeLogService) template.FuncMap {
	return template.FuncMap{
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
		"HasPrefix": strings.HasPrefix,
		"Pagination": func(all, current int64) []Pagination {
			mMap := make([]Pagination, 0)
			var i, num int64
			num = current
			if all-current < 11 {
				num = all - 11
			} else if current >= 6 {
				num = current - 5
			} else if current < 6 {
				num = 1
			}
			for i = num; i <= all; i++ {
				if i-num > 10 && i != all {
					mMap = append(mMap, Pagination{No: i, Current: i == current, Text: "..."})
					break
				}
				mMap = append(mMap, Pagination{No: i, Current: i == current, Text: strconv.FormatInt(i, 10)})
			}
			return mMap
		},
		"MathSub": func(o, n int64) int64 {
			return o - n
		},
		"MathAdd": func(o, n int64) int64 {
			return o + n
		},
	}
}

/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"git.cm/naiba/gocd"
)

func setFuncMap(r *gin.Engine) {
	r.SetFuncMap(template.FuncMap{
		"RepoPipelines": func(rid uint) []gocd.Pipeline {
			return pipelineService.RepoPipelines(&gocd.Repository{ID: rid})
		},
	})
}

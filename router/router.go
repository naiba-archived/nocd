/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/utrack/gin-csrf"
	"git.cm/naiba/gocd"
	"git.cm/naiba/gocd/sqlite3"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"git.cm/naiba/com"
	"net/http"
	"strings"
)

var userService gocd.UserService
var serverService gocd.ServerService
var repoService gocd.RepositoryService
var pipelineService gocd.PipelineService
var pipelogService gocd.PipeLogService

var oauthConf *oauth2.Config

func init() {
	// 地址验证
	binding.Validator.RegisterValidation("address", func(v *validator.Validate, topStruct reflect.Value, currentStruct reflect.Value, field reflect.Value, fieldtype reflect.Type, fieldKind reflect.Kind, param string) bool {
		return com.IsDomain(field.String()) || com.IsIPv4(field.String())
	})
}

func Start() {
	initService()
	initOauthConf()

	r := initEngine()

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "page/index", commonData(c, true, gin.H{
		}))
	})

	serveOauth2(r)
	servePipeline(r)
	serveServer(r)
	serveRepository(r)
	serveSttings(r)
	serveWebHook(r)

	r.Run(":8000")
}

func initOauthConf() {
	// init github oauth2
	oauthConf = &oauth2.Config{
		ClientID:     gocd.Conf.Section("third_party").Key("github_oauth2_client_id").String(),
		ClientSecret: gocd.Conf.Section("third_party").Key("github_oauth2_client_secret").String(),
		Scopes:       []string{},
		Endpoint:     githuboauth.Endpoint,
	}
}

func initEngine() *gin.Engine {
	// init router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	setFuncMap(r)
	// csrf protection
	r.Use(sessions.Sessions("gocd_session", sessions.NewCookieStore([]byte(gocd.Conf.Section("gocd").Key("cookie_key_pair").String()))))
	r.Use(csrf.Middleware(csrf.Options{
		Secret: gocd.Conf.Section("gocd").Key("cookie_key_pair").String(),
		ErrorFunc: func(c *gin.Context) {
			if !strings.HasPrefix(c.Request.URL.Path, "/webhook") {
				gocd.Log.Debug(c.Request.URL.Path)
				c.String(http.StatusForbidden, "CSRF Token 验证失败")
				c.Abort()
			}
		},
	}))
	r.Static("/static", "resource/static")
	r.LoadHTMLGlob("resource/template/**/*")
	r.Use(authMiddleware)
	return r
}

func initService() {
	// init service
	db, err := gorm.Open("sqlite3", "conf/app.db?_loc="+gocd.Conf.Section("gocd").Key("loc").String())
	if err != nil {
		gocd.Log.Panicln(err)
	}
	if gocd.Debug {
		db.Debug()
		db.LogMode(gocd.Debug)
	}
	db.AutoMigrate(gocd.User{}, gocd.Server{}, gocd.Repository{}, gocd.Pipeline{}, gocd.PipeLog{})
	// user service
	sus := sqlite3.UserService{DB: db,}
	userService = &sus
	// server service
	ss := sqlite3.ServerService{DB: db}
	serverService = &ss
	// repo service
	rs := sqlite3.RepositoryService{DB: db}
	repoService = &rs
	// pipeline service
	ps := sqlite3.PipelineService{DB: db}
	pipelineService = &ps
	// pipelog service
	pl := sqlite3.PipeLogService{DB: db}
	pipelogService = &pl
}

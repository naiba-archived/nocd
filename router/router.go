/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package router

import (
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/getsentry/raven-go"
	"github.com/gin-contrib/sentry"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	csrf "github.com/utrack/gin-csrf"

	// sqlite支持
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/naiba/nocd"
	"github.com/naiba/nocd/sqlite3"
	"github.com/naiba/nocd/utils/mgin"
)

var userService nocd.UserService
var serverService nocd.ServerService
var repoService nocd.RepositoryService
var pipelineService nocd.PipelineService
var pipelogService nocd.PipeLogService
var webhookService nocd.WebhookService
var db *gorm.DB

var oauthConf *oauth2.Config

//Start 运行Web
func Start() {
	initService()
	initOauthConf()

	r := initEngine()

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "page/index", mgin.CommonData(c, true, gin.H{}))
	})

	serveOauth2(r)
	servePipeline(r)
	serveServer(r)
	serveRepository(r)
	serveWebHook(r)
	serveAdmin(r)
	serveUser(r)

	r.Run(nocd.Conf.Section("nocd").Key("web_listen").String())
}

func initOauthConf() {
	// init github oauth2
	oauthConf = &oauth2.Config{
		ClientID:     nocd.Conf.Section("third_party").Key("github_oauth2_client_id").String(),
		ClientSecret: nocd.Conf.Section("third_party").Key("github_oauth2_client_secret").String(),
		Scopes:       []string{},
		Endpoint:     github.Endpoint,
	}
}

func initEngine() *gin.Engine {
	// 初始化Sentry
	raven.SetDSN(nocd.Conf.Section("third_party").Key("sentry_dsn").String())
	// init router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(sentry.Recovery(raven.DefaultClient, false))
	r.Use(gin.Recovery())
	r.Use(mgin.AuthMiddleware(userService))
	r.SetFuncMap(mgin.FuncMap(pipelineService, pipelogService, webhookService))
	// csrf protection
	r.Use(sessions.Sessions("nocd_session", cookie.NewStore([]byte(nocd.Conf.Section("nocd").Key("cookie_key_pair").String()))))
	r.Use(csrf.Middleware(csrf.Options{
		Secret: nocd.Conf.Section("nocd").Key("cookie_key_pair").String(),
		ErrorFunc: func(c *gin.Context) {
			if !strings.HasPrefix(c.Request.URL.Path, "/webhook/") {
				nocd.Logger().Infoln(c.Request.URL.Path)
				c.String(http.StatusForbidden, "CSRF Token 验证失败")
				c.Abort()
			}
		},
	}))
	r.Static("/static", "resource/static")
	r.LoadHTMLGlob("resource/template/**/*")
	return r
}

func initService() {
	var err error
	// init service
	db, err = gorm.Open("sqlite3", "conf/app.db?_loc="+nocd.Conf.Section("nocd").Key("loc").String())
	if err != nil {
		nocd.Logger().Panicln(err)
	}
	if nocd.Debug {
		db.Debug()
		db.LogMode(nocd.Debug)
	}
	db.AutoMigrate(nocd.User{}, nocd.Server{}, nocd.Repository{}, nocd.Pipeline{}, nocd.PipeLog{}, nocd.Webhook{})

	upgradeV002(db)
	nocd.InitStats(db)

	userService = &sqlite3.UserService{DB: db}
	serverService = &sqlite3.ServerService{DB: db}
	repoService = &sqlite3.RepositoryService{DB: db}
	pipelineService = &sqlite3.PipelineService{DB: db}
	pipelogService = &sqlite3.PipeLogService{DB: db}
	webhookService = &sqlite3.WebhookService{DB: db}
}

func upgradeV002(db *gorm.DB) {
	// 赋予管理员权限
	var admin nocd.User
	db.Where("id = 1").First(&admin)
	if !admin.IsAdmin {
		admin.IsAdmin = true
		db.Save(admin)
	}
	// 给空昵称用户添加昵称
	var emptyName []nocd.User
	db.Where("g_name = ''").Find(&emptyName)
	for _, u := range emptyName {
		u.GName = u.GLogin
		db.Save(u)
	}
}

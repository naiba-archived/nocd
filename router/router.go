/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package router

import (
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/getsentry/raven-go"
	"github.com/gin-contrib/sentry"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	csrf "github.com/utrack/gin-csrf"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/text/language"

	// sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/naiba/nocd"
	"github.com/naiba/nocd/sqlite3"
	"github.com/naiba/nocd/utils/mgin"
)

var (
	userService     nocd.UserService
	serverService   nocd.ServerService
	repoService     nocd.RepositoryService
	pipelineService nocd.PipelineService
	pipelogService  nocd.PipeLogService
	webhookService  nocd.WebhookService
	i18nBundle      *i18n.Bundle
	oauthConf       *oauth2.Config
	db              *gorm.DB
)

const localizerKey = "i18n"

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
	r.Use(func(c *gin.Context) {
		var localizer *i18n.Localizer
		if strings.Contains(c.Request.Header.Get("accept-language"), "zh") {
			localizer = i18n.NewLocalizer(i18nBundle, "zh")
		} else {
			localizer = i18n.NewLocalizer(i18nBundle, "en")
		}
		c.Set(localizerKey, localizer)
	})
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

	nocd.InitStats(db)

	i18nBundle = i18n.NewBundle(language.English)
	i18nBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	i18nBundle.MustLoadMessageFile("resource/i18n/zh.toml")
	i18nBundle.MustLoadMessageFile("resource/i18n/en.toml")

	userService = &sqlite3.UserService{DB: db}
	serverService = &sqlite3.ServerService{DB: db}
	repoService = &sqlite3.RepositoryService{DB: db}
	pipelineService = &sqlite3.PipelineService{DB: db}
	pipelogService = &sqlite3.PipeLogService{DB: db}
	webhookService = &sqlite3.WebhookService{DB: db}
}

/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package web

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
)

var userService gocd.UserService
var oauthConf *oauth2.Config

func Start() {
	initService()

	r := initEngine()
	r.Use(authMiddleware)

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "page/index", commonData(c, gin.H{
			"csrf_token": csrf.GetToken(c),
		}))
	})

	initOauthConf()
	serveOauth2(r)

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
	// csrf protection
	r.Use(sessions.Sessions("gocd_session", sessions.NewCookieStore([]byte(gocd.Conf.Section("gocd").Key("cookie_key_pair").String()))))
	r.Use(csrf.Middleware(csrf.Options{
		Secret: gocd.Conf.Section("gocd").Key("cookie_key_pair").String(),
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF Token 验证失败")
			c.Abort()
		},
	}))
	r.Static("/static", "resource/static")
	r.LoadHTMLGlob("resource/template/**/*")
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
	db.AutoMigrate(gocd.User{})
	sus := sqlite3.UserService{DB: db,}
	userService = &sus
}

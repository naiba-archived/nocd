/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/naiba/com"
	"github.com/naiba/nocd"
	"github.com/naiba/nocd/utils/mgin"
	"golang.org/x/oauth2"
)

func serveOauth2(r *gin.Engine) {
	oauth2router := r.Group("/oauth2")
	oauth2router.Use(mgin.FilterMiddleware(mgin.FilterOption{Guest: true}))
	{
		oauth2router.POST("/login", func(c *gin.Context) {
			session := sessions.Default(c)
			oauthToken := com.RandomString(18)
			session.Set("oauth_token", oauthToken)
			session.Save()
			c.Redirect(http.StatusMovedPermanently, oauthConf.AuthCodeURL(oauthToken, oauth2.AccessTypeOnline))
		})

		oauth2router.GET("/callback", func(c *gin.Context) {
			type oauthCallback struct {
				State       string `form:"state"`
				RedirectURI string `form:"redirect_uri"`
				Code        string `form:"code"`
			}
			var call oauthCallback
			if err := c.ShouldBindQuery(&call); err != nil {
				c.String(http.StatusForbidden, "回调参数有误")
				return
			}
			// delete oauth_token
			session := sessions.Default(c)
			if session.Get("oauth_token").(string) != call.State {
				c.String(http.StatusForbidden, "登录未授权，请从首页重新登录")
				return
			}
			session.Delete("oauth_token")
			token, err := oauthConf.Exchange(context.Background(), call.Code)
			if err != nil {
				c.String(http.StatusForbidden, "回调验证失败")
				return
			}
			client := github.NewClient(oauthConf.Client(context.Background(), token))
			user, _, err := client.Users.Get(context.Background(), "")
			if err != nil {
				c.String(http.StatusInternalServerError, "GitHub通信失败，请重试")
				return
			}

			// 检测入库
			u, err := userService.UserByGID(user.GetID())
			if err != nil {
				// 首次登录
				if err == gorm.ErrRecordNotFound {
					if err != nil {
						nocd.Logger().Errorln(err)
						c.String(http.StatusInternalServerError, "生成私钥失败，请再次常试")
						return
					}
					u = new(nocd.User)
					u.GID = uint(user.GetID())
					u.GLogin = user.GetLogin()
					if len(user.GetName()) > 0 {
						u.GName = user.GetName()
					} else {
						u.GName = u.GLogin
					}
					if userService.Create(u) != nil {
						nocd.Logger().Errorln(err)
						c.String(http.StatusInternalServerError, "数据库错误")
						return
					}
					// 首位用户赋管理员权限
					if u.ID == 1 {
						u.IsAdmin = true
						userService.Update(u)
					}
				} else {
					nocd.Logger().Errorln(err)
					c.String(http.StatusInternalServerError, "数据库错误")
					return
				}
			}
			// 更新token
			u.Token = com.MD5(fmt.Sprintf("%d%d%s%d", u.ID, u.GID, u.GLogin, time.Now().UnixNano()))
			if userService.Update(u) != nil {
				nocd.Logger().Errorln(err)
				c.String(http.StatusInternalServerError, "数据库错误")
				return
			}
			mgin.SetCookie(c, "uid", fmt.Sprintf("%d", u.ID))
			mgin.SetCookie(c, "token", u.Token)
			c.Redirect(http.StatusMovedPermanently, "/")
		})
	}
}

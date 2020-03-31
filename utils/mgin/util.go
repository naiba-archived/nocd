/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package mgin

import (
	"html/template"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"

	"github.com/naiba/nocd"
)

//SetCookie 设置Cookie
func SetCookie(c *gin.Context, key string, val string) {
	c.SetCookie(key, val, 60*60*24*365*1.5, "/", "", false, false)
}

//CommonData 公共参数
func CommonData(c *gin.Context, csrfToken bool, data gin.H) gin.H {
	data["stat"] = nocd.GetStats()
	data["domain"] = nocd.Conf.Section("nocd").Key("domain").String()
	data["router"] = c.Request.RequestURI
	data["GA_id"] = nocd.Conf.Section("third_party").Key("google_analysis").String()
	isLogin := c.GetBool(CtxIsLogin)
	data["isLogin"] = isLogin
	if isLogin {
		data["user"] = c.MustGet(CtxUser)
	}
	if csrfToken {
		data["csrf_token"] = template.HTML(`<input type="hidden" name="_csrf" value="` + csrf.GetToken(c) + `">`)
	}
	return data
}

//AlertAndRedirect 弹窗并跳转
func AlertAndRedirect(msg, url string, c *gin.Context) {
	c.Writer.WriteString(`
<script>
alert('` + msg + `');window.location.href='` + url + `'
</script>
`)
	c.Abort()
}

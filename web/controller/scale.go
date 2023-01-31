package controller

import (
	"net/http"
	"psyWeb/utils"
	"psyWeb/web/models"
	"psyWeb/web/views"
)

func HandleScale(w http.ResponseWriter, r *http.Request) {
	// 验证用户是否已登录
	if ok, err := utils.IsLogged(w, r); !ok || (err != nil) {
		return
	}
	user := models.User{}
	if err := views.ParseJson(w, r, &user); err != nil {
		return
	}
	cookie, err := r.Cookie("PhoneNumber")
	if err != nil {
		return
	}
	user.PhoneNumber = cookie.Value
	if err := user.UpdateScaleResult(); err != nil {
		return
	}
	views.RenderHtmlPage(w, "upload.html", nil)
}

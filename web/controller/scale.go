package controller

import (
	"net/http"
	"psyWeb/utils"
	"psyWeb/web/models"
	"psyWeb/web/views"
)

func HandleScale(w http.ResponseWriter, r *http.Request) {
	if !verifyUserInformation(w, r) {
		return
	}
	user := models.User{}
	err := views.ParseJson(w, r, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		views.RenderHtmlPage(w, "failed.html", "服务器解析JSON错误")
		return
	}
	user.PhoneNumber, err = utils.GetPhoneNumberFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		views.RenderHtmlPage(w, "failed.html", "Cookie错误，请检查浏览器Cookie设置")
		return
	}
	if err := user.UpdateScaleResult(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		views.RenderHtmlPage(w, "failed.html", "服务器内部数据库错误")
		return
	}
	views.RenderHtmlPage(w, "success.html", nil)
}

package controller

import (
	"net/http"
	"psyWeb/utils"
	"psyWeb/web/views"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	views.RenderHtmlPage(w, "login.html", nil)
}

func HandleStaffJmp(w http.ResponseWriter, r *http.Request) {
	views.RenderHtmlPage(w, "staff.html", nil)
}

func HandleScaleSASJmp(w http.ResponseWriter, r *http.Request) {
	// 验证用户是否已登录
	if ok, err := utils.IsLogged(w, r); !ok || (err != nil) {
		return
	}
	views.RenderHtmlPage(w, "ess.html", nil)
}

func HandleScaleESSJmp(w http.ResponseWriter, r *http.Request) {
	// 验证用户是否已登录
	if ok, err := utils.IsLogged(w, r); !ok || (err != nil) {
		return
	}
	views.RenderHtmlPage(w, "isi.html", nil)
}

func HandleScaleISIJmp(w http.ResponseWriter, r *http.Request) {
	// 验证用户是否已登录
	if ok, err := utils.IsLogged(w, r); !ok || (err != nil) {
		return
	}
	views.RenderHtmlPage(w, "sds.html", nil)
}

func HandleScaleSDSJmp(w http.ResponseWriter, r *http.Request) {
	// 验证用户是否已登录
	if ok, err := utils.IsLogged(w, r); !ok || (err != nil) {
		return
	}
	views.RenderHtmlPage(w, "upload.html", nil)
}

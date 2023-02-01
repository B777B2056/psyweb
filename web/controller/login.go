package controller

import (
	"log"
	"net/http"
	"psyWeb/utils"
	"psyWeb/web/models"
	"psyWeb/web/views"
)

func HandleVerificationCode(w http.ResponseWriter, r *http.Request) {
	models.SendVerificationCodeToUserPhone(r.URL.Query().Get("phone_number"))
}

func HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		PhoneNumber:      r.FormValue("phone_number"),
		VerificationCode: r.FormValue("verification_code"),
	}
	result := user.Check()
	if !result {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// 登录成功，创建Token
	err := utils.AuthenticateUserLogin(w, user.PhoneNumber, user.VerificationCode)
	if err != nil {
		log.Printf("Authenticate User Login, err=%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		views.RenderHtmlPage(w, "failed.html", "服务器内部错误：Token生成失败")
		return
	}
	// 制作html
	views.RenderHtmlPage(w, "nav.html", nil)
}

func verifyUserInformation(w http.ResponseWriter, r *http.Request) bool {
	// 验证用户是否已登录
	if ok, err := utils.IsLogged(w, r); !ok || (err != nil) {
		w.WriteHeader(http.StatusUnauthorized)
		views.RenderHtmlPage(w, "failed.html", "用户校验失败，请重新登录！")
		return false
	}
	return true
}

package controller

import (
	"log"
	"net/http"
	"psyWeb/utils"
	"psyWeb/web/models"
	"psyWeb/web/views"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	views.RenderHtmlPage(w, "login.html", nil)
}

func HandleScaleTestJmp(w http.ResponseWriter, r *http.Request) {
	if !verifyUserInformation(w, r) {
		return
	}
	views.RenderHtmlPage(w, "sas.html", nil)
}

func HandleEEGFileJmp(w http.ResponseWriter, r *http.Request) {
	if !verifyUserInformation(w, r) {
		return
	}
	views.RenderHtmlPage(w, "upload.html", nil)
}

func HandleQueryReportJmp(w http.ResponseWriter, r *http.Request) {
	if !verifyUserInformation(w, r) {
		return
	}
	phone_number, err := utils.GetPhoneNumberFromCookie(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusMethodNotAllowed)
		views.RenderHtmlPage(w, "failed.html", "Cookie错误，请检查浏览器Cookie设置")
		return
	}
	status, _ := models.QueryUserStatus(phone_number)
	if utils.Done != status {
		views.RenderHtmlPage(w, "failed.html", "报告未出，请耐心等待！")
		return
	}
	result, err := models.GetDiagnosticResult(phone_number)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		views.RenderHtmlPage(w, "failed.html", "服务器内部数据库错误")
		return
	}
	views.RenderHtmlPage(w, "index.html", &result)
}

func HandleScaleSASJmp(w http.ResponseWriter, r *http.Request) {
	if !verifyUserInformation(w, r) {
		return
	}
	views.RenderHtmlPage(w, "ess.html", nil)
}

func HandleScaleESSJmp(w http.ResponseWriter, r *http.Request) {
	if !verifyUserInformation(w, r) {
		return
	}
	views.RenderHtmlPage(w, "isi.html", nil)
}

func HandleScaleISIJmp(w http.ResponseWriter, r *http.Request) {
	if !verifyUserInformation(w, r) {
		return
	}
	views.RenderHtmlPage(w, "sds.html", nil)
}

func HandleScaleSDSJmp(w http.ResponseWriter, r *http.Request) {
	if !verifyUserInformation(w, r) {
		return
	}
	views.RenderHtmlPage(w, "success.html", "提交成功！")
}

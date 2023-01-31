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
	result, status := user.Check()
	if !result {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// 登录成功，创建Token
	err := utils.AuthenticateUserLogin(w, user.PhoneNumber, user.VerificationCode)
	if err != nil {
		log.Printf("Authenticate User Login, err=%s", err)
		return
	}
	// 制作html
	switch status {
	case utils.NotExist:
		http.Redirect(w, r, "/", http.StatusNotFound)
	case utils.New:
		views.RenderHtmlPage(w, "sas.html", nil)
	case utils.WaitForReport:
		views.RenderHtmlPage(w, "failed.html", "结果未出，请耐心等待")
	case utils.Done:
		result, err := models.GetDiagnosticResult(user.PhoneNumber)
		if err != nil {
			log.Println(err)
			return
		}
		views.RenderHtmlPage(w, "index.html", &result)
	}
}

func HandleStaffLogin(w http.ResponseWriter, r *http.Request) {
	staff := models.StaffUser{
		Id:       r.FormValue("StaffName"),
		Password: r.FormValue("Password"),
	}
	result := staff.IsPassVerification()
	if !result {
		return
	}
	err := utils.AuthenticateStaffLogin(w, staff.Id, staff.Password)
	if err != nil {
		log.Printf("Authenticate Staff Login, err=%s", err)
		return
	}
	views.RenderHtmlPage(w, "upload.html", nil)
}

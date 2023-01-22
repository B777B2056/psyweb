package controller

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"psyWeb/configuration"
	"psyWeb/utils"
	"psyWeb/web/models"
	"psyWeb/web/views"
	"strings"
)

func HandleResource(response http.ResponseWriter, path string) {
	log.Printf("path: %s", path)
	views.RenderStaticResources(response, path)
}

func HandleHome(response http.ResponseWriter, root_path string) {
	path := root_path + "html/login.html"
	HandleResource(response, path)
}

func HandleVerificationCode(response http.ResponseWriter, request *http.Request) {
	var user models.User
	if err := views.ParseJson(response, request, &user); err != nil {
		return
	}
	user.SendVerificationCodeToUserPhone()
}

func HandleUserLogin(response http.ResponseWriter, request *http.Request) {
	var user models.User
	if err := views.ParseJson(response, request, &user); err != nil {
		return
	}
	views.RenderJson(response, (&models.UserVerification{}).Check(user))
}

func HandleStaffLogin(response http.ResponseWriter, request *http.Request) {
	var staff models.StaffUser
	if err := views.ParseJson(response, request, &staff); err != nil {
		return
	}
	result := staff.IsPassVerification()
	views.RenderJson(response, &result)
}

func HandleScale(response http.ResponseWriter, request *http.Request) {
	var user models.User
	if err := views.ParseJson(response, request, &user); err != nil {
		return
	}
	user.UpdateScaleResult()
}

func handleOneFileUpload(request *http.Request, html_input_id string) string {
	// 处理上传文件
	request.ParseMultipartForm(32 << 20)
	file, handler, err := request.FormFile(html_input_id)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer file.Close()
	work_dir, _ := os.Getwd()
	eeg_data_path := work_dir + "/" + configuration.GetConfigInstance().EEGDataPath + handler.Filename
	f, err := os.OpenFile(eeg_data_path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer f.Close()
	io.Copy(f, file)
	return strings.TrimSuffix(eeg_data_path, path.Ext(handler.Filename))
}

func HandleEEGUpload(response http.ResponseWriter, request *http.Request) {
	pathNoSuffix := handleOneFileUpload(request, "set_file")
	// 校验是否上传的同一被试的数据
	if pathNoSuffix == handleOneFileUpload(request, "fdt_file") {
		// 通知深度学习模型进程，跑结果
		dl := utils.GetDeepLearningInstance()
		if err := dl.Do(pathNoSuffix); err != nil {
			log.Println(err)
		}
	} else {
		log.Println("Wraning: Upload different subject data.")
	}
}

func HandleGetReport(response http.ResponseWriter, request *http.Request) {
	phone_number := request.URL.Query().Get("phone_number")
	report := models.UserReport{}
	err := report.GetResult(phone_number)
	if err != nil {
		log.Println(err)
		return
	}
	views.RenderJson(response, report)
}

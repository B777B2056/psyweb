package controller

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"psyWeb/configuration"
	"psyWeb/utils"
	"psyWeb/web/views"
	"strings"
)

func handleOneFileUpload(r *http.Request, html_input_id string) string {
	// 处理上传文件
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile(html_input_id)
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

func HandleEEGUpload(w http.ResponseWriter, r *http.Request) {
	// 验证用户是否已登录
	if ok, err := utils.IsLogged(w, r); !ok || (err != nil) {
		return
	}
	pathNoSuffix := handleOneFileUpload(r, "set_file")
	// 校验是否上传的同一被试的数据
	if pathNoSuffix == handleOneFileUpload(r, "fdt_file") {
		views.RenderHtmlPage(w, "success.html", "请耐心等待诊断结果")
		// 通知深度学习模型进程，跑结果
		dl := utils.GetDeepLearningInstance()
		if err := dl.Do(pathNoSuffix); err != nil {
			log.Println(err)
		}
	} else {
		log.Println("Wraning: Upload different subject data.")
	}
}

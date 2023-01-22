package router

import (
	"net/http"
	"psyWeb/configuration"
	"psyWeb/web/controller"
	"regexp"
	"strings"
)

type PsyWebServer struct {
	rootPath string
}

func (svr *PsyWebServer) Init() {
	svr.rootPath = configuration.GetConfigInstance().ViewRootPath
}

func (svr *PsyWebServer) pathMapping(path string) string {
	if path[0] == '/' {
		path = path[1:]
	}
	file_path := svr.rootPath
	if match, _ := regexp.MatchString("p([0-9]*)v([0-9]*)", path); match {
		path = path[strings.IndexByte(path, '/'):]
	}
	if strings.HasSuffix(path, "html") {
		file_path = file_path + "html"
	}
	if path[0] != '/' {
		file_path += "/"
	}
	return file_path + path
}

func (svr *PsyWebServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	switch url := request.URL.Path; url {
	case "/":
		// 渲染login页面
		controller.HandleHome(response, svr.rootPath)
	case "/login/staff":
		// 工作人员登录
		controller.HandleStaffLogin(response, request)
	case "/login/get_verification_code":
		// 获取验证码
		controller.HandleVerificationCode(response, request)
	case "/login/verification":
		// 登录
		controller.HandleUserLogin(response, request)
	case "/scale":
		// 量表
		controller.HandleScale(response, request)
	case "/eeg_data":
		// EEG数据文件上传
		controller.HandleEEGUpload(response, request)
	default:
		matched, _ := regexp.MatchString("/report*", url)
		if matched {
			// 获取诊断报告:
			controller.HandleGetReport(response, request)
		} else {
			// 处理资源请求
			controller.HandleResource(response, svr.pathMapping(url))
		}
	}
}

package router

import (
	"net/http"
	"psyWeb/web/controller"
)

type PsyWebServer struct {
}

func (svr *PsyWebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch url := r.URL.Path; url {
	case "/":
		controller.HandleHome(w, r)
	case "/get_verification_code":
		controller.HandleVerificationCode(w, r)
	case "/login":
		controller.HandleUserLogin(w, r)
	case "/scale":
		controller.HandleScale(w, r)
	case "/eeg_data_upload":
		controller.HandleEEGUpload(w, r)
	case "/jmp_to_scale":
		controller.HandleScaleTestJmp(w, r)
	case "/jmp_to_eeg":
		controller.HandleEEGFileJmp(w, r)
	case "/jmp_to_report":
		controller.HandleQueryReportJmp(w, r)
	case "/jmp_to_scale_sas":
		controller.HandleScaleSASJmp(w, r)
	case "/jmp_to_scale_ess":
		controller.HandleScaleESSJmp(w, r)
	case "/jmp_to_scale_isi":
		controller.HandleScaleISIJmp(w, r)
	case "/jmp_to_scale_sds":
		controller.HandleScaleSDSJmp(w, r)
	default:
		controller.HandleStaticResource(w, r)
	}
}

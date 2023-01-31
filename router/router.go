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
	case "/login_staff":
		controller.HandleStaffLogin(w, r)
	case "/get_verification_code":
		controller.HandleVerificationCode(w, r)
	case "/login_user":
		controller.HandleUserLogin(w, r)
	case "/scale":
		controller.HandleScale(w, r)
	case "/eeg_data":
		controller.HandleEEGUpload(w, r)
	case "/staff":
		controller.HandleStaffJmp(w, r)
	case "/scale_sas":
		controller.HandleScaleSASJmp(w, r)
	case "/scale_ess":
		controller.HandleScaleESSJmp(w, r)
	case "/scale_isi":
		controller.HandleScaleISIJmp(w, r)
	case "/scale_sds":
		controller.HandleScaleSDSJmp(w, r)
	default:
		controller.HandleStaticResource(w, r)
	}
}

package views

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func ParseJson(response http.ResponseWriter, request *http.Request, v interface{}) error {
	content, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		goto RET
	}
	err = json.Unmarshal(content, v)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		goto RET
	}
RET:
	return err
}

func RenderJson(response http.ResponseWriter, data interface{}) {
	if bytes, err := json.Marshal(data); err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
	} else {
		response.Write(bytes)
	}
}

func RenderStaticResources(response http.ResponseWriter, path string) {
	if strings.HasSuffix(path, "js") {
		response.Header().Set("Content-Type", "text/javascript")
	} else if strings.HasSuffix(path, "css") {
		response.Header().Set("Content-Type", "text/css")
	} else if strings.HasSuffix(path, "html") {
		response.Header().Set("Content-Type", "text/html")
	} else {
		log.Printf("Content-Type %s not supported", response.Header().Get("Content-Type"))
	}

	if data, err := os.ReadFile(path); err != nil {
		response.WriteHeader(http.StatusNotFound)
	} else {
		response.Write(data)
	}
}

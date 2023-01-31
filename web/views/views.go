package views

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"psyWeb/configuration"
	"strings"
	"text/template"
)

func ParseJson(w http.ResponseWriter, r *http.Request, v interface{}) error {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		goto RET
	}
	err = json.Unmarshal(content, v)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		goto RET
	}
RET:
	return err
}

func RenderJson(w http.ResponseWriter, data interface{}) {
	if bytes, err := json.Marshal(data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(bytes)
	}
}

func pathMapping(path string) string {
	if path[0] == '/' {
		path = path[1:]
	}
	file_path := configuration.GetConfigInstance().ViewRootPath
	if path[0] != '/' {
		file_path += "/"
	}
	return file_path + path
}

func RenderStaticResources(w http.ResponseWriter, path string) {
	path = pathMapping(path)
	log.Printf("path: %s", path)
	if strings.HasSuffix(path, "js") {
		w.Header().Set("Content-Type", "text/javascript")
	} else if strings.HasSuffix(path, "css") {
		w.Header().Set("Content-Type", "text/css")
	} else {

	}
	if data, err := os.ReadFile(path); err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Write(data)
	}
}

func RenderHtmlPage(w http.ResponseWriter, path string, data interface{}) {
	t := template.Must(template.ParseFiles(configuration.GetConfigInstance().ViewRootPath + "html/" + path))
	t.Execute(w, data)
}

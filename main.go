package main

import (
	"log"
	"net/http"
	"psyWeb/configuration"
	"psyWeb/router"
	"psyWeb/utils"
	"psyWeb/web/models"
	"regexp"
)

func initDataBase() {
	db := utils.GetPsyWebDataBaseInstance()
	if err := db.ConnectToSQL(); err != nil {
		log.Fatal(err)
	}
}

func initDeepLearning() {
	// 创建深度学习模型运行子进程
	dl := utils.GetDeepLearningInstance()
	dl.RegistMsgHandler(func(msg []byte) {
		matched, _ := regexp.MatchString("[0-9]{11}", string(msg))
		if matched && (len(msg) == 11) {
			user := models.User{PhoneNumber: string(msg)}
			user.UpdateEEGResult()
		} else {
			log.Printf("recv msg: %s, len=%d", msg, len(msg))
		}
	})
	if err := dl.Start(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDataBase()
	initDeepLearning()
	defer utils.GetDeepLearningInstance().Stop()
	// 服务器运行
	svr := &router.PsyWebServer{}
	if err := http.ListenAndServe(configuration.GetConfigInstance().Port, svr); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

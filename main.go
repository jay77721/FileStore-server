package main

import (
	"filestore-server/handler"
	"filestore-server/redis"
	"fmt"
	"net/http"
	"os"
)

func main() {
	dir, _ := os.Getwd()
	redis.InitRedis()
	fmt.Println("Current working directory:", dir)

	// 提供 /static/ 下的静态文件
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileHandler)
	http.HandleFunc("/file/query", handler.FileQueryHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)

	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	http.HandleFunc("/file/upload/chunk", handler.UploadChunkHandler)
	http.HandleFunc("/file/upload/status", handler.UploadStatusHandler)
	http.HandleFunc("/file/upload/merge", handler.MergeChunkHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start server,err:%s", err.Error())
	}

}

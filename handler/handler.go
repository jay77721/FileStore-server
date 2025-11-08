package handler

import (
	"encoding/json"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("UploadHandler called")
	switch r.Method {
	case "GET":
		//返回上传页面（index.html）
		http.ServeFile(w, r, "./static/view/index.html")
	case "POST":
		//解析上传的文件
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get data", http.StatusBadRequest)
			return
		}
		defer file.Close()

		//创建上传目录
		os.MkdirAll("./uploads", 0755)
		dstPath := filepath.Join("./uploads", header.Filename)

		loc, _ := time.LoadLocation("Asia/Shanghai")
		now := time.Now().In(loc)

		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			Location: dstPath,
			UploadAt: now,
		}

		dst, err := os.Create(fileMeta.Location)
		if err != nil {
			http.Error(w, "Failed to create file", http.StatusBadRequest)
			return
		}
		defer dst.Close()

		fileMeta.FileSize, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to upload file", http.StatusBadRequest)
			return
		}

		dst.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(dst)
		//meta.UpdateFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)

		//打印文件hash值（测试）
		fmt.Fprintf(w, "上传成功！文件SHA1: %s", fileMeta.FileSha1)

		//http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}

}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished")
}

// GetFileMetaHandler: 获取文件元信息
func GetFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	filehash := r.Form.Get("filehash")
	//fMeta := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(filehash)

	file, err := os.Open(fMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+fMeta.FileName+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(data)
}

// FileMetaUpdateHandler:更新元信息接口（重命名）
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method == "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	//meta.UpdateFileMeta(curFileMeta)
	_ = meta.UpdateFileMetaDB(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// FileDeleteHandler: 删除文件及元信息
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	meta.RemoveFileMeta(fileSha1)

	fileMeta := meta.GetFileMeta(fileSha1)
	os.RemoveAll(fileMeta.Location)

	w.WriteHeader(http.StatusOK)
}

// FileQueryHandler:返回所有文件元信息列表
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("only support GET method"))
		return
	}

	//获取内存中所有文件元信息
	fileMetas := meta.GetAllFileMeta()

	//转成JSON
	data, err := json.Marshal(fileMetas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("JSON Marshal fail"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

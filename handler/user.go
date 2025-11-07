package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"net/http"
	"time"
)

const (
	pwd_salt = "*#890"
)

// SignupHandler : 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SignupHandler called") // <- 确认函数是否触发
	if r.Method == "GET" {
		http.ServeFile(w, r, "./static/view/signup.html")
		return
	}

	//if err := r.ParseForm(); err != nil {
	//	w.Write([]byte("parse form error"))
	//	return
	//}

	//username := r.Form.Get("username")
	//password := r.Form.Get("password")
	//err := r.ParseMultipartForm(10 << 20) // 最大 10MB
	//if err != nil {
	//	fmt.Println("ParseMultipartForm error:", err)
	//	w.Write([]byte("fail"))
	//	return
	//}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		fmt.Println("ParseMultipartForm error:", err)
		w.Write([]byte("fail"))
		return
	}
	usernameArr := r.MultipartForm.Value["username"]
	passwordArr := r.MultipartForm.Value["password"]
	if len(usernameArr) == 0 || len(passwordArr) == 0 {
		fmt.Println("username or password empty")
		w.Write([]byte("fail"))
		return
	}
	//username := r.FormValue("username")
	//password := r.FormValue("password")
	username := usernameArr[0]
	password := passwordArr[0]
	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("invalid parameter"))
		return
	}

	fmt.Println("username:", username)
	fmt.Println("password:", password)
	//fmt.Println("Received username:", r.FormValue("username"))
	//fmt.Println("Received password:", r.FormValue("password"))

	enc_passwd := util.Sha1([]byte(password + pwd_salt))
	fmt.Println("enc_passwd:", enc_passwd)

	suc := dblayer.UserSignup(username, enc_passwd)
	if suc {
		fmt.Println("UserSignup success")
		w.Write([]byte("success"))
	} else {
		fmt.Println("UserSignup failed!")
		w.Write([]byte("fail"))
	}
}

// SigninHandler :登录接口
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		fmt.Println("ParseMultipartForm error:", err)
		w.Write([]byte("fail"))
		return
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if len(username) == 0 || len(password) == 0 {
		fmt.Println("username or password empty")
		w.Write([]byte("fail"))
		return
	}

	encPasswd := util.Sha1([]byte(password + pwd_salt))

	//1.校验用户名和密码
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	if !pwdChecked {
		w.Write([]byte("fail"))
		return
	}

	//2.生成访问凭证（token）
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("fail"))
		return
	}

	//3.登录成功后重定向到首页
	w.Write([]byte("success"))
	return
}

func GenToken(username string) string {
	// 40位字符 ： md5(username + timestamp + token_salt) + timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

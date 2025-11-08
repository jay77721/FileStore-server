package handler

import (
	dblayer "filestore-server/db"
	mydb "filestore-server/db/mysql"
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

// SignInHandler :登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
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
	//w.Write([]byte("success"))
	//return
	resp := util.RespMsg{
		Code: 0,
		Msg:  "ok",
		Data: struct {
			Location string `json:"Location"`
			Username string `json:"Username"`
			Token    string `json:"Token"`
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.JSONBytes())
}

// UserInfoHandler :查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//1.解析请求参数
	if err := r.ParseForm(); err != nil {
		fmt.Println("ParseForm error:", err)
		w.Write([]byte("fail"))
		return
	}
	username := r.FormValue("username")
	token := r.FormValue("token")
	//2.验证token是否有效
	isValidToken := isTokenValid(username, token)
	if !isValidToken {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//3.查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	//4.组装并响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "ok",
		Data: user,
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.JSONBytes())
}

func GenToken(username string) string {
	// 40位字符 ： md5(username + timestamp + token_salt) + timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// isTokenValid : token是否有效
func isTokenValid(username string, token string) bool {
	//TODO:判断token的时效性 ，是否过期
	//TODO：从数据库表tbl_user_token查询username对应的token信息
	//TODO：对比连个token是否相同
	stmt, err := mydb.DBConn().Prepare(
		"select user_token,expired_at from tbl_user_token where user_name=? limit 1")
	if err != nil {
		fmt.Println("err:", err)
		return false
	}
	defer stmt.Close()

	var expired_at time.Time
	var user_token string

	err = stmt.QueryRow(username).Scan(&user_token, &expired_at)
	if err != nil {
		fmt.Println("Query err:", err)
		return false
	}

	if user_token != token {
		fmt.Println("Token mismath for user ", username)
		return false
	}

	if expired_at.Before(time.Now()) {
		fmt.Println("Token expired for user :", username, "expired at: ", expired_at)
	}
	//Token 有效
	return true
}

//ALTER TABLE tbl_user_token ADD COLUMN expired_at DATETIME;
//ALTER TABLE tbl_user_token ADD COLUMN update_at DATETIME;

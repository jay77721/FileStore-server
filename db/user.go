package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
	"time"
)

// UserSignup：通过用户名及密码完成user表的注册操作
func UserSignup(username string, password string) bool {
	//db := mydb.DBConn()
	//fmt.Println("DB connection:", db)

	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to insert ,err" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Exec failed,err" + err.Error())
		return false
	}
	//if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
	//	return true
	//}
	//return false
	rowsAffected, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("RowsAffected failed:", err)
		return false
	}

	fmt.Println("Rows affected:", rowsAffected)
	if rowsAffected == 0 {
		fmt.Println("User already exists, signup failed")
		return false
	}

	return true
}

func UserSignin(username string, encpwd string) bool {
	//db := mydb.DBConn()
	//fmt.Println("DB connection:", db)
	//

	//db := mydb.DBConn()
	//var version, dbName string
	//err := db.QueryRow("SELECT VERSION()").Scan(&version)
	//if err != nil {
	//	fmt.Println("Version query error:", err)
	//} else {
	//	fmt.Println("MySQL version:", version)
	//}
	//
	//err = db.QueryRow("SELECT DATABASE()").Scan(&dbName)
	//if err != nil {
	//	fmt.Println("Database query error:", err)
	//} else {
	//	fmt.Println("Current database:", dbName)
	//}
	db := mydb.DBConn()
	rows, err := db.Query("SELECT user_name,user_pwd FROM tbl_user")
	if err != nil {
		fmt.Println("Query error:", err)
	}
	for rows.Next() {
		var uname, pwd string
		rows.Scan(&uname, &pwd)
		fmt.Printf("DB row: [%s] [%s]\n", uname, pwd)
	}
	defer rows.Close()

	stmt, err := mydb.DBConn().Prepare("SELECT user_pwd FROM tbl_user WHERE user_name=? limit 1")
	if err != nil {
		fmt.Println("Failed to prepare statement ,err" + err.Error())
		return false
	}
	defer stmt.Close()

	//测试
	fmt.Printf("Login attempt, username:[%s], password:[%s]\n", username, encpwd)

	var storedPwd string
	err = stmt.QueryRow(username).Scan(&storedPwd)
	if err != nil {
		fmt.Println("User does not exist, username:", username)
		return false
	}

	//比较密码是否一致
	if storedPwd != encpwd {
		fmt.Println("Password mismathch for user,username :" + username)
		return false
	}
	fmt.Println("User login success:" + username)
	return true
}

// UpdateToken:刷新用户登录的token
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token(`user_name`,`user_token`,update_at,expired_at) values (?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	updateAt := time.Now()
	expireAt := updateAt.Add(24 * time.Hour) //token 有效期为24h

	_, err = stmt.Exec(username, token, updateAt, expireAt)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

// GetUserInfo
func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mydb.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}

	defer stmt.Close()
	//执行查询操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

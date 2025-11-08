package db

import (
	"database/sql"
	mydb "filestore-server/db/mysql"
	"fmt"
	"time"
)

// OnFileUploadFinished:文件上传完成，保持meta
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string, uploadAt time.Time) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file(`file_sha1`,`file_name`,`file_size`,`file_addr`,`create_at`,status)" +
			" values(?,?,?,?,?,1)")
	if err != nil {
		fmt.Println("failed to prepared statement,err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr, uploadAt)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Printf("file with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false
}

type TableFile struct {
	FileSha1 string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
	CreateAt sql.NullTime
}

// GetFileMeta 从MySQL获取文件元信息
func GetFileMeta(filehash string) (fileMeta *TableFile, err error) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1, file_name, file_size, file_addr FROM tbl_file  where file_sha1= ? and status=1 limit 1")
	if err != nil {
		fmt.Println("failed to prepare statement,err:" + err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileSha1, &tfile.FileName, &tfile.FileSize, &tfile.FileAddr)
	if err != nil {
		fmt.Println("failed to query row,err:" + err.Error())
		return nil, err
	}
	return &tfile, nil
}

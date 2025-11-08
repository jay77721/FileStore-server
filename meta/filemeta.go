package meta

import (
	mydb "filestore-server/db"
	"time"
)

// FileMeta: 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt time.Time
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta:新增/更新文件元信息
func UpdateFileMeta(fMeta FileMeta) {
	fileMetas[fMeta.FileSha1] = fMeta
}

// UpdateFileMetaDB:新增/更新文件元到MySQL中
func UpdateFileMetaDB(fMeta FileMeta) bool {
	return mydb.OnFileUploadFinished(
		fMeta.FileSha1, fMeta.FileName, fMeta.FileSize, fMeta.Location, fMeta.UploadAt)
}

// GetFileMeta:通过sha1值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFileMetaDB:从mysql获取文件元信息
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tfile, err := mydb.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileSha1,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
		UploadAt: tfile.CreateAt.Time,
	}
	return fmeta, nil
}

// RemoveFileMeta：删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}

func GetAllFileMeta() []FileMeta {
	all := make([]FileMeta, 0)
	for _, file := range fileMetas {
		all = append(all, file)
	}
	return all
}

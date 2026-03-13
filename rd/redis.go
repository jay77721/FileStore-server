package rd

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var (
	RDB *redis.Client
	ctx = context.Background()
)

// 初始化Redis连接池
func InitRedis() {

	RDB = redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379",
		Password:     "",
		DB:           0,
		PoolSize:     20, // 连接池大小
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Redis connect failed:", err)
	}

	log.Println("Redis connected")
}

// 记录chunk上传成功
func AddChunk(filehash string, index int) {

	key := "chunk:" + filehash

	RDB.SAdd(ctx, key, index)
}

// 获取已上传chunk
func GetUploadedChunks(filehash string) ([]string, error) {

	key := "chunk:" + filehash

	return RDB.SMembers(ctx, key).Result()
}

// 删除chunk记录
func ClearChunks(filehash string) {

	key := "chunk:" + filehash

	RDB.Del(ctx, key)
}

// 记录文件hash
func SetFileHash(hash string, location string) {

	key := "file:" + hash

	RDB.Set(ctx, key, location, 0)
}

// 查询文件hash
func GetFileHash(hash string) (string, error) {

	key := "file:" + hash

	return RDB.Get(ctx, key).Result()
}

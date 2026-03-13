package util

import (
	"context"
	"filestore-server/redis"
)

var ctx = context.Background()

// 记录chunk上传成功
func AddChunk(filehash string, index int) error {

	key := "chunk:" + filehash

	return redis.RDB.SAdd(ctx, key, index).Err()
}

// 获取已上传chunk
func GetUploadedChunks(filehash string) ([]string, error) {

	key := "chunk:" + filehash

	return redis.RDB.SMembers(ctx, key).Result()
}

// 判断chunk是否存在
func ChunkExists(filehash string, index int) bool {

	key := "chunk:" + filehash

	res, _ := redis.RDB.SIsMember(ctx, key, index).Result()

	return res
}

// 删除chunk记录
func ClearChunks(filehash string) {

	key := "chunk:" + filehash

	redis.RDB.Del(ctx, key)
}

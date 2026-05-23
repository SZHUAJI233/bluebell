package redis

import (
	"strconv"
	"time"
	"web/models"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func CreatePost(postID, communityID uint64) error {
	now := float64(time.Now().Unix())
	pipeline := rdb.Pipeline()
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  now,
		Member: postID,
	})

	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  now,
		Member: postID,
	})

	communityKey := getRedisKey(KeyPostCommunitySetPrefix + strconv.Itoa(int(communityID)))
	pipeline.SAdd(communityKey, postID)

	_, err := pipeline.Exec()
	return err
}

func getIDsFormKey(key string, page int64, size int64) ([]string, error) {
	// 确定查询索引的起始点
	start := (page - 1) * size
	end := start + size - 1
	// 从redis中查询
	return rdb.ZRevRange(key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis中获取id
	// 根据用户携带的order参数
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFormKey(key, p.Page, p.Size)
}

func GetCommunityPostIDsInOrder(p *models.ParamCommunityPostList) ([]string, error) {

	// 社区的key
	communityKey := getRedisKey(KeyPostCommunitySetPrefix + strconv.Itoa(int(p.CommunityID)))

	// 使用zinterstore 把分区的帖子set与贴子分数的zset 生成一个新的zset
	// 缓存key 减少zintstore执行的次数

	// 从redis中获取id
	// 根据用户携带的order参数
	orderKey := getRedisKey(KeyPostTimeZSet)
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
		key = orderKey + strconv.Itoa(int(p.CommunityID))
	}

	if rdb.Exists(key).Val() < 1 {
		pipeline := rdb.Pipeline()
		// 不存在,计算set和zset交集
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, communityKey, orderKey)
		pipeline.Expire(key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			zap.L().Error("ZInterStore failed: ", zap.Error(err))
			return nil, err
		}
	}

	// 存在，直接查询
	return getIDsFormKey(key, p.Page, p.Size)
}

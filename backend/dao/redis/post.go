package redis

import "web/models"

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis中获取id
	// 根据用户携带的order参数
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 确定查询索引的起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// 从redis中查询
	return rdb.ZRevRange(key, start, end).Result()
}

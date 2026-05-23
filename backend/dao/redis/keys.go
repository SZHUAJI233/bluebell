package redis

const (
	KeyPrefix                 = "bluebell:"
	KeyPostTimeZSet           = "post:time"       // zset 帖子及发帖时间
	KeyPostScoreZSet          = "post:score"      // zset 帖子及投票分数
	KeyPostVotedZSetPrefix    = "post:voted:"     // zset 记录用户及投票类型
	KeyPostCommunitySetPrefix = "post:community:" // zset 记录社区及其帖子

	KeyJwtTokenSet = "jwt:token"
)

// 给redis key 加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}

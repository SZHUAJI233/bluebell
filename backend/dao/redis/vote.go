package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// 投票功能
// 1. 用户投票的数据

/*
投票情况：
direction = 1
	1. d=0 >> d=1	--> 更新分数和投票记录	差值的绝对值	 1	+432
	2. d=-1 >> d=1	--> 更新分数和投票记录					2	+432*2
direction = 0
	1. d=1 >> d=0	--> 更新分数和投票记录					1	-432
	2. d=-1 >> d=0	--> 更新分数和投票记录					1	+432
direction = -1
	1. d=0 >> d=-1	--> 更新分数和投票记录					1	-432
	2. d=1 >> d=-1	--> 更新分数和投票记录					2	-432*2
*/

/* 投票限制：
帖子发表之日七一个星期内允许投票。超过不允许投票
1. 到期之后将redis中保存的赞成票及反对存储到mysql中
2. 到期之后删除对应的zset
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每票的分数
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
)

func CreatePost(postID string) error {
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
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, vote float64) error {
	// 1. 判断投票的限制
	postTime, err := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Result()
	if err == redis.Nil {
		zap.L().Error("rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID) postID不存在: ", zap.Error(err))
		return ErrVoteTimeExpire
	} else if err != nil {
		zap.L().Error("rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Result() failed: ", zap.Error(err))
		return err
	}
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 2. 更新帖子分数
	// 查询当前用户对当前帖子的原来的投票分数
	oldVote, err := rdb.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Result()
	if err == redis.Nil {
		oldVote = 0
	} else if err != nil {
		zap.L().Error("ZScore(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Result() failed: ", zap.Error(err))
		return err
	}

	var symbol int
	if vote > oldVote {
		symbol = 1
	} else {
		symbol = -1
	}

	// 计算新旧投票差值	(如何新旧差值相同，则最后增加票数也为0)
	diff := math.Abs(oldVote - vote)

	// 事务
	pipeline := rdb.TxPipeline()
	// 更新票数
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), float64(symbol*scorePerVote)*diff, postID)

	// 3. 记录用户为该帖子投票
	if vote == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  vote,
			Member: userID,
		})
	}
	_, err = pipeline.Exec()
	return err
}

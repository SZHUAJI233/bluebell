package server

import (
	"strconv"
	"web/dao/redis"
	"web/models"

	"go.uber.org/zap"
)

// 为帖子投票
func PostVote(userID uint64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", int64(userID)),
		zap.String("postID", p.PostID),
		zap.Int8("vote", p.Vote))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Vote))
}

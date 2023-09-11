package vote

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

const (
	oneWeekInSecond   = 7 * 86400 * time.Second
	voteScore         = 432
	timeKey           = "time:"
	votedKey          = "voted:"
	scoreKey          = "score:"
	articleFieldVotes = "votes"
)

func ArticleVote(client *redis.Client, user string, article string) {
	cutoff := float64(time.Now().Unix()) - oneWeekInSecond.Seconds()
	ctx := context.Background()
	if client.ZScore(ctx, timeKey, article).Val() < cutoff {
		return
	}
	articleID := strings.Split(article, ":")[1]
	if client.SAdd(ctx, votedKey+articleID, user).Val() == 1 {
		client.ZIncrBy(ctx, scoreKey, voteScore, article)
		client.HIncrBy(ctx, article, articleFieldVotes, 1)
	}
	return
}

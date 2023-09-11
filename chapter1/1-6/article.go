package article

import (
	"context"
	"strings"
	"time"

	"github.com/XdpCs/redis-in-action/common"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func (c *Client) ArticleVote(user string, article string) {
	cutoff := float64(time.Now().Unix()) - common.OneWeekInSecond.Seconds()
	ctx := context.Background()
	if c.client.ZScore(ctx, common.TimeKey, article).Val() < cutoff {
		return
	}
	articleID := strings.Split(article, ":")[1]
	if c.client.SAdd(ctx, common.VotedKey+articleID, user).Val() == 1 {
		c.client.ZIncrBy(ctx, common.ScoreKey, common.VoteScore, article)
		c.client.HIncrBy(ctx, article, common.ArticleFieldVotes, 1)
	}
	return
}

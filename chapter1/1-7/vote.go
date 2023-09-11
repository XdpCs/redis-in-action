package vote

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
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
	articleKey        = "article:"
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

func PostArticle(client *redis.Client, user, title, link string) string {
	ctx := context.Background()
	articleID := strconv.FormatInt(client.Incr(ctx, articleKey).Val(), 10)
	voted := votedKey + articleID
	_ = client.SAdd(ctx, voted, user)
	_ = client.Expire(ctx, voted, oneWeekInSecond)
	now := time.Now().Unix()
	article := articleKey + articleID
	_ = client.HMSet(ctx, article, map[string]interface{}{
		"title":  title,
		"link":   link,
		"poster": user,
		"time":   now,
		"votes":  1,
	})
	_ = client.ZAdd(ctx, scoreKey, redis.Z{
		Score:  float64(now + voteScore),
		Member: article,
	})
	_ = client.ZAdd(ctx, timeKey, redis.Z{
		Score:  float64(now),
		Member: article,
	})
	return article
}

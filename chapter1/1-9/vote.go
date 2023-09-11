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
	groupKey          = "group:"
	articleFieldVotes = "votes"
	articleKey        = "article:"
	articlesPerPage   = 25
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

func GetArticles(client *redis.Client, page int, order string) []map[string]string {
	ctx := context.Background()
	start := (page - 1) * articlesPerPage
	end := start + articlesPerPage - 1
	ids := client.ZRevRange(ctx, order, int64(start), int64(end)).Val()
	var articles []map[string]string
	for _, id := range ids {
		articleData := client.HGetAll(ctx, id).Val()
		articleData["id"] = id
		articles = append(articles, articleData)
	}
	return articles
}

func AddRemoveGroups(client *redis.Client, articleID string, toAdd []string, toRemove []string) {
	ctx := context.Background()
	article := articleKey + articleID
	for _, group := range toAdd {
		_ = client.SAdd(ctx, groupKey+group, article)
	}
	for _, group := range toRemove {
		_ = client.SRem(ctx, groupKey+group, article)
	}
}

package article

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/XdpCs/redis-in-action/common"
	"github.com/redis/go-redis/v9"
)

func (c *Client) ArticleVote(user string, article string, voteType string) {
	cutoff := float64(time.Now().Unix()) - common.OneWeekInSecond.Seconds()
	ctx := context.Background()
	if c.client.ZScore(ctx, common.TimeKey, article).Val() < cutoff {
		return
	}
	articleID := strings.Split(article, ":")[1]
	if c.client.SAdd(ctx, common.VotedKey+articleID, user).Val() == 1 {
		switch voteType {
		case common.AffirmativeVotes:
			c.client.ZIncrBy(ctx, common.ScoreKey, common.VoteScore, article)
			c.client.HIncrBy(ctx, article, common.AffirmativeVotes, 1)
		case common.OpposingVotes:
			c.client.ZIncrBy(ctx, common.ScoreKey, -common.VoteScore, article)
			c.client.HIncrBy(ctx, article, common.OpposingVotes, 1)
		}

	}
	return
}

func (c *Client) PostArticle(user, title, link string) string {
	ctx := context.Background()
	articleID := strconv.FormatInt(c.client.Incr(ctx, common.ArticleKey).Val(), 10)
	voted := common.VotedKey + articleID
	_ = c.client.SAdd(ctx, voted, user)
	_ = c.client.Expire(ctx, voted, common.OneWeekInSecond)
	now := time.Now().Unix()
	article := common.ArticleKey + articleID
	_ = c.client.HMSet(ctx, article, map[string]interface{}{
		"title":                 title,
		"link":                  link,
		"poster":                user,
		"time":                  now,
		common.AffirmativeVotes: 1,
		common.OpposingVotes:    0,
	})
	_ = c.client.ZAdd(ctx, common.ScoreKey, redis.Z{
		Score:  float64(now + common.VoteScore),
		Member: article,
	})
	_ = c.client.ZAdd(ctx, common.TimeKey, redis.Z{
		Score:  float64(now),
		Member: article,
	})
	return articleID
}

func (c *Client) GetArticles(page int, order string) []map[string]string {
	if order == "" {
		order = common.ScoreKey
	}
	ctx := context.Background()
	start := (page - 1) * common.ArticlesPerPage
	end := start + common.ArticlesPerPage - 1

	ids := c.client.ZRevRange(ctx, order, int64(start), int64(end)).Val()
	var articles []map[string]string
	for _, id := range ids {
		articleData := c.client.HGetAll(ctx, id).Val()
		articleData["id"] = id
		articles = append(articles, articleData)
	}
	return articles
}

func (c *Client) AddRemoveGroups(articleID string, toAdd, toRemove []string) {
	ctx := context.Background()
	article := common.ArticleKey + articleID
	for _, group := range toAdd {
		_ = c.client.SAdd(ctx, common.GroupKey+group, article)
	}
	for _, group := range toRemove {
		_ = c.client.SRem(ctx, common.GroupKey+group, article)
	}
}

func (c *Client) GetGroupArticles(group string, page int, order string) []map[string]string {
	if order == "" {
		order = common.ScoreKey
	}
	ctx := context.Background()
	key := order + group
	if c.client.Exists(ctx, key).Val() == 0 {
		_ = c.client.ZInterStore(ctx, key, &redis.ZStore{
			Aggregate: "MAX",
			Keys:      []string{common.GroupKey + group, order},
		})
		_ = c.client.Expire(ctx, key, 60*time.Second)
	}
	return c.GetArticles(page, key)
}

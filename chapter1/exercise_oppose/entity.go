package article

import "github.com/redis/go-redis/v9"

type Client struct {
	client *redis.Client
}

func NewClient(client *redis.Client) *Client {
	return &Client{client: client}
}

type Article interface {
	ArticleVote(string, string, string)
	PostArticle(string, string, string) string
	GetArticles(int, string) []map[string]string
	AddRemoveGroups(string, []string, []string)
	GetGroupArticles(string, int, string) []map[string]string
}

package article

import (
	"context"
	"github.com/XdpCs/redis-in-action/util"
	"strconv"
	"testing"

	"github.com/XdpCs/redis-in-action/common"
	"github.com/redis/go-redis/v9"
)

func TestClient_PostArticle(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer teardown(t, client, []string{common.ArticleKey, common.VotedKey + "1",
		common.VotedKey + "2", common.ArticleKey + "1", common.ArticleKey + "2",
		common.ScoreKey, common.TimeKey})
	c := NewClient(client)
	oneArticleID := c.PostArticle("XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	twoArticleID := c.PostArticle("XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	cases := []struct {
		Name     string
		Actual   string
		Expected string
	}{
		{
			Name:     "OnePostArticle",
			Actual:   "1",
			Expected: oneArticleID,
		},
		{
			Name:     "TwoPostArticle",
			Actual:   "2",
			Expected: twoArticleID,
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if c.Actual != c.Expected {
				t.Fatalf("actual:[%v] expected:[%v]", c.Actual, c.Expected)
			}
		})
	}

}

func TestClient_ArticleVote(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	c := NewClient(client)
	defer teardown(t, client, []string{common.ArticleKey, common.VotedKey + "1",
		common.ArticleKey + "1", common.ScoreKey, common.TimeKey})
	ctx := context.Background()
	oneArticleID := c.PostArticle("XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	c.ArticleVote("XdpCs1", common.ArticleKey+oneArticleID, common.AffirmativeVotes)
	c.ArticleVote("XdpCs2", common.ArticleKey+oneArticleID, common.OpposingVotes)
	all := c.client.HGetAll(ctx, common.ArticleKey+oneArticleID).Val()
	initOneArticle := c.client.ZScore(ctx, common.ScoreKey, common.ArticleKey+oneArticleID).Val()
	if all[common.AffirmativeVotes] != "2" {
		t.Fatalf("actual:[%v] expected:[%v]", all[common.AffirmativeVotes], "1")
	}
	if all[common.OpposingVotes] != "1" {
		t.Fatalf("actual:[%v] expected:[%v]", all[common.OpposingVotes], "1")
	}
	finalOneArticle := c.client.ZScore(ctx, common.ScoreKey, common.ArticleKey+oneArticleID).Val()
	if util.IsEqual(initOneArticle, finalOneArticle) {
		t.Fatalf("actual:[%v] expected:[%v]", initOneArticle, finalOneArticle)
	}
}

func TestClient_GetArticles(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	c := NewClient(client)
	defer teardown(t, client, []string{common.ArticleKey, common.VotedKey + "1",
		common.VotedKey + "2", common.ArticleKey + "1", common.ArticleKey + "2",
		common.ScoreKey, common.TimeKey})
	_ = c.PostArticle("XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	_ = c.PostArticle("XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	articles := c.GetArticles(0, common.ScoreKey)
	compareArticles(t, articles)
}

func TestClient_AddRemoveGroups(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	c := NewClient(client)
	defer teardown(t, client, []string{common.ArticleKey, common.VotedKey + "1",
		common.VotedKey + "2", common.VotedKey + "3", common.ArticleKey + "1",
		common.ArticleKey + "2", common.ArticleKey + "3", common.ScoreKey,
		common.TimeKey, common.GroupKey + "new-group", common.ScoreKey + "new-group"})
	oneArticleID := c.PostArticle("XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	twoArticleID := c.PostArticle("XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	threeArticleID := c.PostArticle("XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	c.AddRemoveGroups(oneArticleID, []string{"new-group"}, []string{})
	c.AddRemoveGroups(twoArticleID, []string{"new-group"}, []string{})
	c.AddRemoveGroups(threeArticleID, []string{"new-group"}, []string{})
	articleIDs := client.SMembers(context.Background(), common.GroupKey+"new-group").Val()
	for i, ID := range articleIDs {
		if ID != common.ArticleKey+strconv.FormatInt(int64(i+1), 10) {
			t.Fatalf("actual:[%v] expected:[%v]", ID, common.GroupKey+strconv.FormatInt(int64(i+1), 10))
		}
	}
	c.AddRemoveGroups(oneArticleID, []string{}, []string{"new-group"})
	c.AddRemoveGroups(twoArticleID, []string{}, []string{"new-group"})
	articleIDs = client.SMembers(context.Background(), common.GroupKey+"new-group").Val()
	for _, ID := range articleIDs {
		if ID != common.ArticleKey+"3" {
			t.Fatalf("actual:[%v] expected:[%v]", ID, common.ArticleKey+"3")
		}
	}
}

func TestClient_GetGroupArticles(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	c := NewClient(client)
	defer teardown(t, client, []string{common.ArticleKey, common.VotedKey + "1",
		common.VotedKey + "2", common.VotedKey + "3", common.ArticleKey + "1",
		common.ArticleKey + "2", common.ArticleKey + "3", common.ScoreKey,
		common.TimeKey, common.GroupKey + "new-group", common.ScoreKey + "new-group"})
	oneArticleID := c.PostArticle("XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	twoArticleID := c.PostArticle("XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	threeArticleID := c.PostArticle("XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	c.AddRemoveGroups(oneArticleID, []string{"new-group"}, []string{})
	c.AddRemoveGroups(twoArticleID, []string{"new-group"}, []string{})
	c.AddRemoveGroups(threeArticleID, []string{"new-group"}, []string{})
	articles := c.GetGroupArticles("new-group", 0, common.ScoreKey)
	compareArticles(t, articles)
}

func teardown(t *testing.T, rc *redis.Client, keys []string) {
	t.Helper()

	for _, key := range keys {
		if err := rc.Del(context.Background(), key).Err(); err != nil {
			t.Fatal(err)
		}
	}

	if err := rc.Close(); err != nil {
		t.Fatal(err)
	}
}

func compareArticles(t *testing.T, articles []map[string]string) {
	t.Helper()
	for i, article := range articles {
		if article["id"] != common.ArticleKey+strconv.FormatInt(int64(len(articles)-i), 10) {
			t.Fatalf("actual:[%v] expected:[%v]", article["id"], common.ArticleKey+strconv.FormatInt(int64(2-i), 10))
		}
		if article["title"] != "Xdp's girlfriend" {
			t.Fatalf("actual:[%v] expected:[%v]", article["title"], "Xdp's girlfriend")
		}
		if article["link"] != "https://github.com/XdpCs" {
			t.Fatalf("actual:[%v] expected:[%v]", article["link"], "https://github.com/XdpCs")
		}
		if article[common.AffirmativeVotes] != "1" {
			t.Fatalf("actual:[%v] expected:[%v]", article["votes"], "1")
		}
		if article["poster"] != "XdpCs" {
			t.Fatalf("actual:[%v] expected:[%v]", article["poster"], "XdpCs")
		}
	}
}

package vote

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"testing"
)

func TestPostArticle(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	oneArticleID := PostArticle(client, "XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	twoArticleID := PostArticle(client, "XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	defer teardown(t, client, []string{articleKey, votedKey + "1", votedKey + "2", oneArticleID, twoArticleID, scoreKey, timeKey})
	cases := []struct {
		Name     string
		Actual   string
		Expected string
	}{
		{
			Name:     "OnePostArticle",
			Actual:   articleKey + "1",
			Expected: oneArticleID,
		},
		{
			Name:     "TwoPostArticle",
			Actual:   articleKey + "2",
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

func TestArticleVote(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx := context.Background()
	oneArticleID := PostArticle(client, "XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	twoArticleID := PostArticle(client, "XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	defer teardown(t, client, []string{articleKey, votedKey + "1", votedKey + "2", votedKey + "3", votedKey + "4", oneArticleID, twoArticleID, scoreKey, timeKey})
	ArticleVote(client, "XdpCs1", oneArticleID)
	ArticleVote(client, "XdpCs2", oneArticleID)
	all := client.HGetAll(ctx, oneArticleID).Val()
	if all["votes"] != "3" {
		t.Fatalf("actual:[%v] expected:[%v]", all["votes"], "3")
	}
	ArticleVote(client, "XdpCs1", twoArticleID)
	ArticleVote(client, "XdpCs2", twoArticleID)
	all = client.HGetAll(ctx, twoArticleID).Val()
	if all["votes"] != "3" {
		t.Fatalf("actual:[%v] expected:[%v]", all["votes"], "3")
	}
}

func TestGetArticles(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	oneArticleID := PostArticle(client, "XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	twoArticleID := PostArticle(client, "XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	defer teardown(t, client, []string{articleKey, votedKey + "1", votedKey + "2", oneArticleID, twoArticleID, scoreKey, timeKey})
	articles := GetArticles(client, 0, scoreKey)
	compareArticles(t, articles)
}

func TestAddRemoveGroups(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	oneArticleID := PostArticle(client, "XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	twoArticleID := PostArticle(client, "XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	threeArticleID := PostArticle(client, "XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	defer teardown(t, client, []string{articleKey, votedKey + "1", votedKey + "2", votedKey + "3", oneArticleID, twoArticleID, threeArticleID, scoreKey, timeKey, groupKey + "new-group"})
	AddRemoveGroups(client, "1", []string{"new-group"}, []string{})
	AddRemoveGroups(client, "2", []string{"new-group"}, []string{})
	AddRemoveGroups(client, "3", []string{"new-group"}, []string{})
	articleIDs := client.SMembers(context.Background(), groupKey+"new-group").Val()
	for i, ID := range articleIDs {
		if ID != articleKey+strconv.FormatInt(int64(i+1), 10) {
			t.Fatalf("actual:[%v] expected:[%v]", ID, articleKey+strconv.FormatInt(int64(i+1), 10))
		}
	}
	AddRemoveGroups(client, "1", []string{}, []string{"new-group"})
	AddRemoveGroups(client, "2", []string{}, []string{"new-group"})
	articleIDs = client.SMembers(context.Background(), groupKey+"new-group").Val()
	for _, ID := range articleIDs {
		if ID != articleKey+"3" {
			t.Fatalf("actual:[%v] expected:[%v]", ID, articleKey+"3")
		}
	}
}

func TestGetGroupArticles(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	oneArticleID := PostArticle(client, "XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	twoArticleID := PostArticle(client, "XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	threeArticleID := PostArticle(client, "XdpCs", "Xdp's girlfriend", "https://github.com/XdpCs")
	defer teardown(t, client, []string{articleKey, votedKey + "1", votedKey + "2", votedKey + "3", oneArticleID, twoArticleID, threeArticleID, scoreKey, timeKey, groupKey + "new-group"})
	AddRemoveGroups(client, "1", []string{"new-group"}, []string{})
	AddRemoveGroups(client, "2", []string{"new-group"}, []string{})
	AddRemoveGroups(client, "3", []string{"new-group"}, []string{})
	articles := GetGroupArticles(client, "new-group", 0, scoreKey)
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
		if article["id"] != articleKey+strconv.FormatInt(int64(len(articles)-i), 10) {
			t.Fatalf("actual:[%v] expected:[%v]", article["id"], articleKey+strconv.FormatInt(int64(2-i), 10))
		}
		if article["title"] != "Xdp's girlfriend" {
			t.Fatalf("actual:[%v] expected:[%v]", article["title"], "Xdp's girlfriend")
		}
		if article["link"] != "https://github.com/XdpCs" {
			t.Fatalf("actual:[%v] expected:[%v]", article["link"], "https://github.com/XdpCs")
		}
		if article["votes"] != "1" {
			t.Fatalf("actual:[%v] expected:[%v]", article["votes"], "1")
		}
		if article["poster"] != "XdpCs" {
			t.Fatalf("actual:[%v] expected:[%v]", article["poster"], "XdpCs")
		}
	}
}

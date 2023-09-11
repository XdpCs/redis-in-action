package common

import "time"

const (
	OneWeekInSecond   = 7 * 86400 * time.Second
	VoteScore         = 432
	TimeKey           = "time:"
	VotedKey          = "voted:"
	ScoreKey          = "score:"
	GroupKey          = "group:"
	ArticleFieldVotes = "votes"
	ArticleKey        = "article:"
	ArticlesPerPage   = 25
	AffirmativeVotes  = "affirmative_votes"
	OpposingVotes     = "opposing_votes"
)

func ReverseVoteType(voteType string) string {
	if voteType == AffirmativeVotes {
		return OpposingVotes
	}
	return AffirmativeVotes
}

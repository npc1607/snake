package game

import "sort"

const leaderboardLimit = 10

// ScoreEntry records one completed game score.
type ScoreEntry struct {
	Rank  int
	Score int
}

func (e *Engine) recordScore() {
	if e.score <= 0 || e.recorded {
		return
	}

	e.scores = append(e.scores, e.score)
	sort.Slice(e.scores, func(i int, j int) bool {
		return e.scores[i] > e.scores[j]
	})

	if len(e.scores) > leaderboardLimit {
		e.scores = e.scores[:leaderboardLimit]
	}

	e.recorded = true
}

func scoreEntries(scores []int) []ScoreEntry {
	entries := make([]ScoreEntry, len(scores))
	for index, score := range scores {
		entries[index] = ScoreEntry{
			Rank:  index + 1,
			Score: score,
		}
	}

	return entries
}

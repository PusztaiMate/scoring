package main

import "strconv"

type PlayerScore struct {
	scores []float64
	id     int
	name   string
}

func NewPlayerScore(initial float64, name string, id int) *PlayerScore {
	return &PlayerScore{scores: []float64{initial}, id: id, name: name}
}

func (ps *PlayerScore) GetLatest() float64 {
	return ps.scores[len(ps.scores)-1]
}

func (ps *PlayerScore) AddScore(score float64) {
	ps.scores = append(ps.scores, score)
}

type PlayerScoreDb struct {
	initialScore float64
	scores       map[int]*PlayerScore
	idNamePairs  map[string]PlayerName
}

func NewPlayerScoreDb(initialScore float64, idNamePairs map[string]PlayerName) *PlayerScoreDb {
	scores := make(map[int]*PlayerScore)
	return &PlayerScoreDb{
		initialScore: initialScore,
		scores:       scores,
		idNamePairs:  idNamePairs,
	}
}

func (db *PlayerScoreDb) GetLatestScoreForPlayer(playerId int) float64 {
	if score, ok := db.scores[playerId]; ok {
		return score.GetLatest()
	}
	name := db.idNamePairs[strconv.Itoa(playerId)].String()
	newPlayerScore := NewPlayerScore(db.initialScore, name, playerId)
	db.scores[playerId] = newPlayerScore
	return newPlayerScore.GetLatest()
}

func (db *PlayerScoreDb) UpdateScoreForPlayer(playerId int, diff float64) {
	if score, ok := db.scores[playerId]; ok {
		newScore := score.GetLatest() + diff
		score.AddScore(newScore)
		return
	}

	name := db.idNamePairs[strconv.Itoa(playerId)].String()
	newPlayerScore := NewPlayerScore(db.initialScore+diff, name, playerId)
	db.scores[playerId] = newPlayerScore
}

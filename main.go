package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// First create a type alias
type UWHTrainingDateTime time.Time

// Implement Marshaler and Unmarshaler interface
func (j *UWHTrainingDateTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	// we get someting like this here: "2022-09-20T06:43:00.308692"
	// let's only extract the date for now, hopefully no 2 games on a single night
	date := strings.Split(s, "T")[0]

	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return err
	}
	*j = UWHTrainingDateTime(t)
	return nil
}

func (j UWHTrainingDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

type TraininData struct {
	BlackScorers    map[string]int      `json:"black_scorers"`
	BlackScore      int                 `json:"black_score"`
	BlackTeam       []int               `json:"black_team"`
	Date            UWHTrainingDateTime `json:"date"`
	Id              int                 `json:"id"`
	Players         []int               `json:"players"`
	SwimmingPlayers []int               `json:"swimming_players"`
	WhiteScore      int                 `json:"white_score"`
	WhiteTeam       []int               `json:"white_team"`
	WhiteScorers    map[string]int      `json:"white_scorers"`
}

type PlayerName struct {
	Fname string `json:"fname"`
	Lname string `json:"lname"`
}

func (p PlayerName) String() string {
	return fmt.Sprintf("%s %s", p.Lname, p.Fname)
}

func getAllTrainings() ([]TraininData, error) {
	response, err := http.Get("https://www.vizalattihoki.hu/api/v1/trainings/")
	if err != nil {
		return nil, err
	}

	var trainings []TraininData
	defer response.Body.Close()
	json.NewDecoder(response.Body).Decode(&trainings)

	return trainings, nil
}

func getIdNamePairs() (map[string]PlayerName, error) {
	response, err := http.Get("https://www.vizalattihoki.hu/api/v1/players/id-name-pairs/")
	if err != nil {
		return nil, err
	}

	var idNameMap map[string]PlayerName
	defer response.Body.Close()
	json.NewDecoder(response.Body).Decode(&idNameMap)

	return idNameMap, nil
}

func pointCalculator1(training TraininData, playerScoreDb *PlayerScoreDb) float64 {
	if training.BlackScore == 0 && training.WhiteScore == 0 {
		// we have no score data from this match
		t := time.Time(training.Date)
		s := t.Format("2006-01-02")
		log.Printf("skipping training from %s, number of players is %d", s, len(training.Players))
		return 0.0
	}

	// calculates avarage: (2 strongest player strength) + ((rest avarage) * 4)
	whiteTeamStrength := CalculateTeamStrength2(training.WhiteTeam, playerScoreDb)
	blackTeamStrength := CalculateTeamStrength2(training.BlackTeam, playerScoreDb)

	scoreDiff := training.WhiteScore - training.BlackScore
	teamStrengthDiff := whiteTeamStrength - blackTeamStrength

	scoreDiffModifier := 0.1
	teamStrengthDiffModifier := 0.1

	diffForWhite := scoreDiffModifier*float64(scoreDiff) - teamStrengthDiffModifier*teamStrengthDiff

	return diffForWhite
}

func main() {
	pairs, err := getIdNamePairs()
	if err != nil {
		log.Fatal(err)
	}

	playerScoreDb := NewPlayerScoreDb(5.0, pairs)
	trainings, err := getAllTrainings()
	if err != nil {
		log.Fatal(err)
	}

	// lets assume for now, that the trainings are in cronological order (probably true)

	lastXTraining := 100
	for _, training := range trainings[len(trainings)-lastXTraining:] {
		// update white team score
		diffForWhite := pointCalculator1(training, playerScoreDb)
		if diffForWhite == 0.0 {
			continue
		}

		for _, playerId := range training.WhiteTeam {
			playerScoreDb.UpdateScoreForPlayer(playerId, diffForWhite)
		}

		// update black team score, note the negative diffForWhite
		for _, playerId := range training.BlackTeam {
			playerScoreDb.UpdateScoreForPlayer(playerId, -diffForWhite)
		}
	}

	for playerId, score := range playerScoreDb.scores {
		if len(score.scores) < 5 {
			continue
		}
		name := pairs[strconv.Itoa(playerId)]
		fmt.Printf("%s pontja %.2f\n", name, score.GetLatest())
	}
}

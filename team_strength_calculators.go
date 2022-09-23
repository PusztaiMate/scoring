package main

import "sort"

func sumFloatSlice(s []float64, from, to int) float64 {
	sum := 0.0
	for i := from; i < to; i++ {
		sum += s[i]
	}
	return sum
}

func CalculateTeamStrength1(members []int, scoreDb *PlayerScoreDb) float64 {
	strengthValues := make([]float64, len(members))
	for i := 0; i < len(members); i++ {
		strengthValues[i] = scoreDb.GetLatestScoreForPlayer(members[i])
	}
	sort.Float64s(strengthValues)
	strongestPlayer := strengthValues[len(strengthValues)-1]
	secondStrongestPlayer := strengthValues[len(strengthValues)-2]

	sum := 0.0
	for i := 0; i < len(strengthValues)-2; i++ {
		sum += strengthValues[i]
	}
	teamStrength := (sum/float64(len(strengthValues)-2))*4 + strongestPlayer + secondStrongestPlayer
	return teamStrength
}

func CalculateTeamStrength2(members []int, scoreDb *PlayerScoreDb) float64 {
	strengthValues := make([]float64, len(members))
	for i, id := range members {
		strengthValues[i] = scoreDb.GetLatestScoreForPlayer(id)
	}
	sort.Float64s(strengthValues)

	strongestPlayer := strengthValues[len(strengthValues)-1]
	secondStrongestPlayer := strengthValues[len(strengthValues)-2]
	avgOfRest := sumFloatSlice(strengthValues, 0, len(strengthValues)-2) / float64((len(strengthValues) - 2))

	return 2*(strongestPlayer+secondStrongestPlayer) + avgOfRest*4
}

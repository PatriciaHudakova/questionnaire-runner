package rating

import (
	"errors"
	"skillQuiz/pkg"
	"skillQuiz/pkg/db"
)

func mockCalculateImmediateRating(answers []pkg.Question) (string, int) {
	return "40", 0
}

func mockCalculateAverageRating(db db.IDatabase, currentRating string, currPositives int) (string, error) {
	return "60", nil
}

func mockCalculateAverageRatingErr(db db.IDatabase, currentRating string, currPositives int) (string, error) {
	return "", errors.New("test error")
}

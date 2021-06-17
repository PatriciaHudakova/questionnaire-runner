package rating

import (
	"fmt"
	"skillQuiz/pkg"
	"skillQuiz/pkg/db"
)

type CurrentRun func(answers []pkg.Question) (string, int)
type AverageRun func(db db.IDatabase, currentRating string, positives int) (string, error)

// PrintRatings is a wrapper function to calculate and print current & average ratings
func PrintRatings(currentRunFunc CurrentRun, averageRunFunc AverageRun, db db.IDatabase, answers []pkg.Question) error {
	// Based on user input, calculate the current rating
	currentRating, yesses := currentRunFunc(answers)
	fmt.Printf("Your rating is: %s/100\n", currentRating)

	// Using the persisted data from the database and new user input, print the new average
	averageRating, err := averageRunFunc(db, currentRating, yesses)
	if err != nil {
		return fmt.Errorf("something went wrong calculating your average score: %v", err)
	}
	fmt.Printf("The average rating is: %s/100", averageRating)

	return nil
}

// CalculateImmediateRating calculates the rating of the current run
func CalculateImmediateRating(answers []pkg.Question) (string, int) {
	count := 0

	// Base Case
	if len(answers) <= 0 {
		return "0", 0
	}

	// Iterate through the responses and accumulate total points for the run
	for _, question := range answers {
		count = count + question.Value
	}

	// Calculate a percentage like rating (using integers would result in 0)
	rating := 100 * (float64(count) / float64(len(answers)))

	// Round the rating to 0 decimal places for consistency
	return fmt.Sprintf("%.0f", rating), count
}

// CalculateAverageRating calculates the average rating of all runs
func CalculateAverageRating(db db.IDatabase, currentRating string, currPositives int) (string, error) {
	// Retrieve all rows from the averages table
	rows, err := db.GetAllRows()
	if err != nil {
		return "", err
	}

	// If there are no entries, current rating becomes the average
	if db.IsEmpty(rows) {
		if err := db.CreateANewEntry(5, currPositives); err != nil {
			return "", fmt.Errorf("unable to persist current average: %v", err)
		}
		return currentRating, nil
	}

	// If not, retrieve params from db and calculate a new rating, then update
	dbQuestions, dbPositives, err := db.GetPersistedParams()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve params: %v", err)
	}

	totalQuestions := dbQuestions + 5
	totalPositives := dbPositives + currPositives

	// Calculate the new average
	newAverage := 100 * (float64(totalPositives) / float64(totalQuestions))

	// Update the database with new totals
	if err = db.UpdateDatabaseParams(totalQuestions, totalPositives); err != nil {
		return "", fmt.Errorf("unable to update average: %v", err)
	}

	return fmt.Sprintf("%.0f", newAverage), nil
}

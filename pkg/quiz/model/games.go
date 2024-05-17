package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/margulan-kalykul/JustQuiz/pkg/quiz/validator"
)

type Game struct {
	Id			string		`json:"id"`
	Finished	string		`json:"finished"`
	Player		int			`json:"player"`
	Quiz		int			`json:"quiz"`
}

type GameModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (g GameModel) GetAll(player, quiz int, from, to string, filters Filters) ([]*Game, Metadata, error) {
	// Retrieve all gamees from the database
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, finished, player, quiz
		FROM games
		WHERE (player = $1 OR $1 = 0)
		AND (quiz = $2 OR $2 = 0)
		AND (finished > $2 OR $2 = '')
		AND (finished < $3 OR $3 = '')
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5;
		`,
		filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our four placeholder parameter values in a slice.
	args := []interface{}{player, quiz, from, to, filters.limit(), filters.offset()}

	// Use QueryContext to execute the query. This returns a sql.Rows result set containing
	// the result.
	rows, err := g.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}	

	// Importantly, defer a call to rows.Close() to ensure that the result set is closed
	// before GetAll returns.
	defer func() {
		if err := rows.Close(); err != nil {
			g.ErrorLog.Println(err)
		}
	}()

	// Declare a totalRecords variable
	totalRecords := 0

	var games []*Game
	for rows.Next() {
		var game Game
		err := rows.Scan(&totalRecords, &game.Id, &game.Finished, &game.Player, &game.Quiz)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the game struct to the slice
		games = append(games, &game)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	
	// Generate a Metadata struct, passing in the total record count and pagination parameters
	// from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of the movies and metadata.
	return games, metadata, nil
}

func (g GameModel) Insert(game *Game) error {
	// Create a new game in the database
	query := `
		INSERT INTO games(player, quiz) 
		VALUES ($1, $2)
		RETURNING id, player, quiz;
		`
	args := []interface{}{game.Player, game.Quiz}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return g.DB.QueryRowContext(ctx, query, args...).Scan(&game.Id, &game.Player, &game.Quiz)
}

func (g GameModel) Get(id int) (*Game, error) {
	// Invalid id. Return an error if the ID is less than 1.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Retrieve a game with its ID
	query := `
		SELECT id, player, quiz
		FROM games
		WHERE id = $1;
		`
	var game Game
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := g.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&game.Id, &game.Player, &game.Quiz)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive game with id: %v, %w", id, err)
	}
	return &game, nil
}

// func (g GameModel) Answer(game *Game) (*Game, error) {
// 	// Invalid id. Return an error if the ID is less than 1.
// 	if id < 1 {
// 		return nil, ErrRecordNotFound
// 	}

// 	// Answer a game with its ID
// 	// Get answers
// 	query := `
// 		SELECT id, answers
// 		FROM quizes
// 		WHERE id = $1;
// 		`
// 	var quiz Quiz
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	row := g.DB.QueryRowContext(ctx, query, id)
// 	err := row.Scan(&quiz.Id, (*pq.StringArray)(&quiz.Answers))
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot retrive quiz with id: %v, %w", id, err)
// 	}
	
// 	// Compare them
// 	if quiz.Answers == 
// }

// Can't update the alredy finished game
// func (g GameModel) Update(game *Game) error {
// 	// Update game name and score
// 	query := `
// 		UPDATE games
// 		SET category = $1, reward = $2, questions = $3, answers = $4
// 		WHERE id = $5
// 		RETURNING id;
// 		`
// 		// pq.Array(game.Questions), pq.Array(game.Answers)
// 	args := []interface{}{game.Category, game.Reward, pq.Array(game.Questions), pq.Array(game.Answers), game.Id}
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	return g.DB.QueryRowContext(ctx, query, args...).Scan(&game.Id)
// }

func (g GameModel) Delete(id int) error {
	// Invalid id. Return an error if the ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}

	// Delete game by id
	query := `
		DELETE FROM games
		WHERE id = $1;
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := g.DB.ExecContext(ctx, query, id)
	return err
}

func ValidateGame(v *validator.Validator, game *Game) {
	v.Check(game.Finished != "", "finished", "must be provided")
}
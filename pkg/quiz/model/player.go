package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/margulan-kalykul/JustQuiz/pkg/quiz/validator"
)

type Player struct {
	Id			string 	`json:"id"`
	Name		string 	`json:"name"`
	Joined		string 	`json:"joined"`
	LastUpdate	string 	`json:"last_update"`
	Score		int		`json:"score"`
}
type PlayerModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m PlayerModel) GetAll(name string, from, to int, filters Filters) ([]*Player, Metadata, error) {
	// Retrieve all players from the database
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, name, joined, last_update, score
		FROM players
		WHERE (LOWER(name) = LOWER($1) OR $1 = '')
		AND (score >= $2 OR $2 = 0)
		AND (score <= $3 OR $3 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5
		`,
		filters.sortColumn(), filters.sortDirection())

	// query := `
	// 	SELECT id, name, joined, last_update, score
	// 	FROM players
	// 	`

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our four placeholder parameter values in a slice.
	args := []interface{}{name, from, to, filters.limit(), filters.offset()}

	// log.Println(query, title, from, to, filters.limit(), filters.offset())
	// Use QueryContext to execute the query. This returns a sql.Rows result set containing
	// the result.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}	

	// Importantly, defer a call to rows.Close() to ensure that the result set is closed
	// before GetAll returns.
	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	// Declare a totalRecords variable
	totalRecords := 0

	var players []*Player
	for rows.Next() {
		var player Player
		err := rows.Scan(&totalRecords, &player.Id, &player.Name, &player.Joined, &player.LastUpdate, &player.Score)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the player struct to the slice
		players = append(players, &player)
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
	return players, metadata, nil
}

func (p PlayerModel) Insert(player *Player) error {
	// Create a new player in the database
	query := `
		INSERT INTO players(name) 
		VALUES ($1) 
		RETURNING id, joined, last_update, score;
		`
	args := []interface{}{player.Name}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&player.Id, &player.Joined, &player.LastUpdate, &player.Score)
}

func (p PlayerModel) Get(id int) (*Player, error) {
	// Invalid id. Return an error if the ID is less than 1.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Retrieve a player with its ID
	query := `
		SELECT id, name, joined, last_update, score
		FROM players
		WHERE id = $1;
		`
	var player Player
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := p.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&player.Id, &player.Name, &player.Joined, &player.LastUpdate, &player.Score)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive player with id: %v, %w", id, err)
	}
	return &player, nil
}

func (p PlayerModel) Update(player *Player) error {
	// Update player name and score
	query := `
		UPDATE players
		SET name = $1, score = $2, last_update = $3
		WHERE id = $4
		RETURNING last_update;
		`
	args := []interface{}{player.Name, player.Score, time.Now(), player.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&player.LastUpdate)
}

func (p PlayerModel) Delete(id int) error {
	// Invalid id. Return an error if the ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}

	// Delete player by id
	query := `
		DELETE FROM players
		WHERE id = $1;
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, id)
	return err
}

func ValidatePlayer(v *validator.Validator, player *Player) {
	// Check if the name field is empty.
	v.Check(player.Name != "", "name", "must be provided")
	// Check if the title name is no more than 100 characters.
	v.Check(len(player.Name) <= 100, "name", "must not be more than 100 bytes long")
	// Check if the score value is not negative.
	v.Check(player.Score >= 0, "score", "must not be negative")
}

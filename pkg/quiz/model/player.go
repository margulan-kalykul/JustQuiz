package model

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
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
	// Invalid id
	if id < 1 {
		return nil, errors.New("Invalid player id")
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
		return nil, err
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
	// Invalid id
	if id < 1 {
		return errors.New("Not Found")
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

package model

import (
	"database/sql"
	"log"
)

type Quiz struct {
	Id         string `json:"id"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	Category   string `json:"category"`
	Difficulty string `json:"difficulty"`
	Q1         string `json:"q1"`
	A1         string `json:"a1"`
	Q2         string `json:"q2"`
	A2         string `json:"a2"`
}

type QuizModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// func (q QuizModel) GetQuizes() ([]Quiz, error) {
// 	query := `
// 		SELECT id, category, dificulty, q1, a1, q2, a2
// 		FROM quizes;
// 		`
// 	var quizes []Quiz
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	row, err := q.DB.QueryContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = row.Scan(quizes)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return quizes, nil
// }

// func (q QuizModel) GetQuizById(id int) (*Quiz, error) {
// 	// Invalid id
// 	if id < 1 {
// 		return nil, errors.New("Not Found")
// 	}

// 	query := `
// 		SELECT id, category, dificulty, q1, a1, q2, a2
// 		FROM quizes
// 		WHERE id = $1;
// 		`
// 	var quizes []Quiz
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	row := q.DB.QueryRowContext(ctx, query, id)
// 	err := row.Scan()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &player, nil
// }

package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/margulan-kalykul/JustQuiz/pkg/quiz/validator"
)

type Quiz struct {
	Id			string		`json:"id"`
	Category	string		`json:"category"`
	Reward		int			`json:"reward"`
	Questions	[]string	`json:"questions"`
	Answers		[]string	`json:"answers"`
}

type QuizModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (q QuizModel) GetAll(category string, from, to int, filters Filters) ([]*Quiz, Metadata, error) {
	// Retrieve all quizes from the database
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, category, reward, questions, answers
		FROM quizes
		WHERE (LOWER(category) = LOWER($1) OR $1 = '')
		AND (reward >= $2 OR $2 = 0)
		AND (reward <= $3 OR $3 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5;
		`,
		filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our four placeholder parameter values in a slice.
	args := []interface{}{category, from, to, filters.limit(), filters.offset()}

	// Use QueryContext to execute the query. This returns a sql.Rows result set containing
	// the result.
	rows, err := q.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}	

	// Importantly, defer a call to rows.Close() to ensure that the result set is closed
	// before GetAll returns.
	defer func() {
		if err := rows.Close(); err != nil {
			q.ErrorLog.Println(err)
		}
	}()

	// Declare a totalRecords variable
	totalRecords := 0

	var quizes []*Quiz
	for rows.Next() {
		var quiz Quiz
		err := rows.Scan(&totalRecords, &quiz.Id, &quiz.Category, &quiz.Reward, (*pq.StringArray)(&quiz.Questions), (*pq.StringArray)(&quiz.Answers))
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the quiz struct to the slice
		quizes = append(quizes, &quiz)
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
	return quizes, metadata, nil
}

func (q QuizModel) Insert(quiz *Quiz) error {
	// Create a new quiz in the database
	query := `
		INSERT INTO quizes(category, reward, questions, answers) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, category, reward;
		`
	args := []interface{}{quiz.Category, quiz.Reward, pq.Array(quiz.Questions), pq.Array(quiz.Answers)}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return q.DB.QueryRowContext(ctx, query, args...).Scan(&quiz.Id, &quiz.Category, &quiz.Reward)
}

func (q QuizModel) Get(id int) (*Quiz, error) {
	// Invalid id. Return an error if the ID is less than 1.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Retrieve a quiz with its ID
	query := `
		SELECT id, category, reward, questions, answers
		FROM quizes
		WHERE id = $1;
		`
	var quiz Quiz
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := q.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&quiz.Id, &quiz.Category, &quiz.Reward, (*pq.StringArray)(&quiz.Questions), (*pq.StringArray)(&quiz.Answers))
	if err != nil {
		return nil, fmt.Errorf("cannot retrive quiz with id: %v, %w", id, err)
	}
	return &quiz, nil
}

func (q QuizModel) Update(quiz *Quiz) error {
	// Update quiz name and score
	query := `
		UPDATE quizes
		SET category = $1, reward = $2, questions = $3, answers = $4
		WHERE id = $5
		RETURNING id;
		`
		// pq.Array(quiz.Questions), pq.Array(quiz.Answers)
	args := []interface{}{quiz.Category, quiz.Reward, pq.Array(quiz.Questions), pq.Array(quiz.Answers), quiz.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return q.DB.QueryRowContext(ctx, query, args...).Scan(&quiz.Id)
}

func (q QuizModel) Delete(id int) error {
	// Invalid id. Return an error if the ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}

	// Delete quiz by id
	query := `
		DELETE FROM quizes
		WHERE id = $1;
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := q.DB.ExecContext(ctx, query, id)
	return err
}

func ValidateQuiz(v *validator.Validator, quiz *Quiz) {
	// Check if the category field is empty.
	v.Check(quiz.Category != "", "category", "must be provided")
	// Check if the category is no more than 100 characters.
	v.Check(len(quiz.Category) <= 100, "category", "must not be more than 100 bytes long")
	// Check if the reward value is not negative.
	v.Check(quiz.Reward >= 0, "reward", "must not be negative")
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

package main

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/margulan-kalykul/JustQuiz/pkg/quiz/model"
	"github.com/margulan-kalykul/JustQuiz/pkg/quiz/validator"
)

func (app *application) createGameHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Player int `json:"player"`
		Quiz   int `json:"quiz"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	game := &model.Game{
		Player: input.Player,
		Quiz:   input.Quiz,
	}

	err = app.models.Games.Insert(game)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"game": game}, nil)
}

func (app *application) getGamesList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Player       int
		Quiz         int
		FinishedFrom string
		FinishedTo   string
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	// Use our helpers to extract the name and score value range query string values, falling back to the
	// defaults of an empty string and an empty slice, respectively, if they are not provided
	// by the client.
	input.Player = app.readInt(qs, "player", 0, v)
	input.Quiz = app.readInt(qs, "quiz", 0, v)
	input.FinishedFrom = app.readStrings(qs, "finisedFrom", "1980-01-01 00:00:00+06")
	input.FinishedTo = app.readStrings(qs, "finishedTo", "1980-01-01 00:00:00+06")

	// Ge the page and page_size query string value as integers. Notice that we set the default
	// page value to 1 and default page_size to 20, and that we pass the validator instance
	// as the final argument.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply an ascending sort on game ID).
	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	// Add the supported sort value for this endpoint to the sort safelist.
	// name of the column in the database.
	input.Filters.SortSafeList = []string{
		// ascending sort values
		"id", "finished", "player", "quiz",
		// descending sort values
		"-id", "-finished", "-player", "-quiz",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	games, metadata, err := app.models.Games.GetAll(input.Player, input.Quiz, input.FinishedFrom, input.FinishedTo, input.Filters)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"games": games, "metadata": metadata}, nil)
}

func (app *application) getGameHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	game, err := app.models.Games.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"game": game}, nil)
}

func (app *application) answerGameHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	quiz, err := app.models.Quizes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Player  *string   `json:"playerId"`
		Answers *[]string `json:"answers"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	playerAnswers := *input.Answers
	if !reflect.DeepEqual(quiz.Answers, playerAnswers) {
		app.writeJSON(w, http.StatusOK, envelope{"result": "Answers are incorrect"}, nil)
		return
	}

	// Give player a score
	playerId, _ := strconv.Atoi(*input.Player)
	player, err := app.models.Players.Get(playerId)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	player.Score += quiz.Reward
	err = app.models.Players.Update(player)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Create new games record
	quizId, _ := strconv.Atoi(quiz.Id)
	game := &model.Game{
		Player: playerId,
		Quiz:   quizId,
	}
	err = app.models.Games.Insert(game)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"result": "Answers are correct"}, nil)
}

// func (app *application) updateGameHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := app.readIDParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}

// 	game, err := app.models.Games.Get(id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, model.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	var input struct {
// 		Finished	*string		`json:"finished"`
// 		Player		*string		`json:"player"`
// 		Quiz		*string		`json:"quiz"`
// 	}

// 	err = app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	// Check fileds
// 	if input.Player != nil {
// 		game.Player = *input.Player
// 	}
// 	if input.Quiz != nil {
// 		game.Quiz = *input.Quiz
// 	}
// 	if input.Finished != nil {
// 		game.Finished = *input.Finished
// 	}

// 	v := validator.New()

// 	if model.ValidateGame(v, game); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}

// 	err = app.models.Games.Update(game)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, model.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	app.writeJSON(w, http.StatusOK, envelope{"game": game}, nil)
// }

func (app *application) deleteGameHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Games.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
}

package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/margulan-kalykul/JustQuiz/pkg/quiz/model"
	"github.com/margulan-kalykul/JustQuiz/pkg/quiz/validator"
)

// func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "status: available")
// 	fmt.Fprintf(w, "environment: %s\n", app.config.env)
// }
	

func (app *application) createPlayerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name	string	`json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	player := &model.Player{
		Name: input.Name,
	}

	err = app.models.Players.Insert(player)
	if err != nil {
		// app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		// return
		app.serverErrorResponse(w, r, err)
		return
	}

	// app.respondWithJSON(w, http.StatusCreated, player)
	app.writeJSON(w, http.StatusCreated, envelope{"player": player}, nil)
}

func (app *application) getPlayersList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name		string
		ScoreFrom	int
		ScoreTo		int
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	// Use our helpers to extract the name and score value range query string values, falling back to the
	// defaults of an empty string and an empty slice, respectively, if they are not provided
	// by the client.
	input.Name = app.readStrings(qs, "name", "")
	input.ScoreFrom = app.readInt(qs, "scoreFrom", 0, v)
	input.ScoreTo = app.readInt(qs, "scoreTo", 0, v)

	// Ge the page and page_size query string value as integers. Notice that we set the default
	// page value to 1 and default page_size to 20, and that we pass the validator instance
	// as the final argument.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply an ascending sort on player ID).
	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	// Add the supported sort value for this endpoint to the sort safelist.
	// name of the column in the database.
	input.Filters.SortSafeList = []string{
		// ascending sort values
		"id", "name", "score", "joined",
		// descending sort values
		"-id", "-name", "-score", "-joined",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	players, metadata, err := app.models.Players.GetAll(input.Name, input.ScoreFrom, input.ScoreTo, input.Filters)


	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"players": players, "metadata": metadata}, nil)
}

func (app *application) getPlayerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	player, err := app.models.Players.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"player": player}, nil)
}

// func (app *application) getPlayerQuizes(w http.ResponseWriter, r *http.Request) {
// 	id, err := app.readIDParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}

// 	player, err := app.models.Players.Get(id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, model.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	app.writeJSON(w, http.StatusOK, envelope{"player": player}, nil)
// }

func (app *application) updatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	player, err := app.models.Players.Get(id)
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
		Name  *string `json:"name"`
		Score *int    `json:"score"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Check fileds
	if input.Name != nil {
		player.Name = *input.Name
	}
	if input.Score != nil {
		player.Score = *input.Score
	}

	v := validator.New()

	if model.ValidatePlayer(v, player); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

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

	app.writeJSON(w, http.StatusOK, envelope{"player": player}, nil)
}

func (app *application) deletePlayerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Players.Delete(id)
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

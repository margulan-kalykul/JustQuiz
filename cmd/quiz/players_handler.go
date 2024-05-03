package main

import (
	"errors"
	"log"
	"net/http"
	"github.com/margulan-kalykul/JustQuiz/pkg/quiz/model"
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
	// TODO: implement filtering
	players, err := app.models.Players.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"players": players}, nil)
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
